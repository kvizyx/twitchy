import { createHash } from "node:crypto";
import { closeSync, constants, fstatSync, openSync, readFileSync } from "node:fs";
import { spawnSync } from "node:child_process";

const expected = {
  "sources/receipt.json": "637790fac19c7cfebf9aaa94d735b34d4bcf203050594832cd5dc66f3a4f1e37",
  "sources/reference.html": "3c5ed44ffac743027b96897a07cca4e4b1d4225dc8fa941a2aec6acc7795e926",
  "sources/guide.html": "d4053e1645c82622053e1280abd4a11960caa177a9453b8ce3d639accee9b2e3",
  "sources/scopes.html": "123a411f410edf057c2eb7aa142b7f08e4335c13964164a2a0a71396a615d55e",
  "sources/oauth.html": "d20843e754c6a22b295d634df99718e6f62d23210e4c6fbdd8cf6e750f609a59",
  "sources/refresh.html": "a0f61ea5226dadb4c31b1f4b915540c246ae05df2d48ad601fcb1e3eeddb9eac",
  "sources/validate.html": "26dbaa4684586580ff0281725505e4491960e9745933be3152d66f4048b75f2a",
  "sources/revoke.html": "2264dd8fe2b79107906a484dc39f2436db8c2b33a9b4fa767aa4ced4108cd3cf",
  "sources/authentication.html": "f361ed78d1e54a392042f554753439a6916cf9e213023b6a5fd1f00eb11750e5",
  "expected-operations.json": "16f361ee282409ab4b5a94ab4a00a585a7a57a76f70855e7c73d93ef8b51615f",
  "public-descriptor.json": "eb230fb36e9f98f06a71a34cc4d2b05c9a39098c51bc032c8f0b6f919c459f73",
  "core-oauth-descriptor.json": "5830f46d453ec3490eb323982fc3a9075c2a16bf8177410dd6e3bc350342ebed",
};
const base = "helix/internal/manifest";
for (const [name, digest] of Object.entries(expected)) {
  const fd = openSync(`${base}/${name}`, constants.O_RDONLY | constants.O_NOFOLLOW);
  const stat = fstatSync(fd);
  if (!stat.isFile()) throw new Error(`${name}: not a regular file`);
  const bytes = readFileSync(fd);
  closeSync(fd);
  const actual = createHash("sha256").update(bytes).digest("hex");
  if (actual !== digest) throw new Error(`${name}: ${actual} != ${digest}`);
}

const operations = JSON.parse(readFileSync(`${base}/public-descriptor.json`, "utf8")).operations;
const core = JSON.parse(readFileSync(`${base}/core-oauth-descriptor.json`, "utf8"));
const doc = (pkg) => {
  const result = spawnSync("go", ["doc", "-all", pkg], {
    encoding: "utf8",
    env: { ...process.env, GOTOOLCHAIN: "local", GOPROXY: "off", GOSUMDB: "off" },
  });
  if (result.status !== 0) throw new Error(result.stderr || `go doc ${pkg} failed`);
  return result.stdout.replaceAll(/\s+/g, " ");
};
const helixDoc = doc("github.com/kvizyx/twitchy/helix");
const oauthDoc = doc("github.com/kvizyx/twitchy/oauth");
for (const operation of operations) {
  for (const required of [operation.service_type, operation.method, operation.request_type, operation.data_type]) {
    if (!helixDoc.includes(required)) throw new Error(`${operation.anchor}: missing ${required}`);
  }
  if (operation.pager_signature && !helixDoc.includes(`${operation.method}Pager`)) {
    throw new Error(`${operation.anchor}: missing pager`);
  }
}
for (const declaration of [
  ...Object.keys(core.core.interfaces),
  ...Object.keys(core.core.enums),
  ...Object.keys(core.core.structs),
  ...core.core.sentinels,
  ...Object.keys(core.oauth.structs),
  ...Object.keys(core.oauth.enums),
  ...core.oauth.errorTypes,
]) {
  if (!(helixDoc.includes(declaration) || oauthDoc.includes(declaration))) {
    throw new Error(`missing core/OAuth declaration ${declaration}`);
  }
}
for (const declaration of [...core.oauth.functions, ...core.oauth.methods]) {
  const name = declaration.match(/[.(]([A-Z][A-Za-z0-9]*)\(/)?.[1];
  if (name && !oauthDoc.includes(name)) throw new Error(`missing OAuth callable ${name}`);
}
const counts = operations.reduce((out, operation) => {
  out[operation.stability] = (out[operation.stability] ?? 0) + 1;
  return out;
}, {});
if (operations.length !== 149 || counts.stable !== 127 || counts.NEW !== 10 || counts.BETA !== 12) {
  throw new Error(`wrong operation partition: ${JSON.stringify(counts)}`);
}
console.log("independent verifier: 149 operations; public contracts and frozen inputs OK");
