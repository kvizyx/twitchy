package conformance

import (
	"encoding/json"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type documentationContract struct {
	Schema int           `json:"schema"`
	Rows   []contractRow `json:"rows"`
}

type contractRow struct {
	ID        string `json:"id"`
	Kind      string `json:"kind"`
	Package   string `json:"package"`
	Symbol    string `json:"symbol"`
	Anchor    string `json:"anchor"`
	Stability string `json:"stability"`
}

type packageDocumentation struct {
	doc     string
	symbols map[string]struct{}
}

type descriptor struct {
	Operations []struct {
		Anchor      string `json:"anchor"`
		ServiceType string `json:"service_type"`
		Method      string `json:"method"`
		RequestType string `json:"request_type"`
		DataType    string `json:"data_type"`
		Stability   string `json:"stability"`
	} `json:"operations"`
}

func TestDocumentationContract(t *testing.T) {
	root := repositoryRoot(t)
	contract := readDocumentationContract(t, filepath.Join(root, "helix", "internal", "conformance", "public_contract.json"))
	if contract.Schema != 1 {
		t.Fatalf("contract schema = %d, want 1", contract.Schema)
	}
	docs := map[string]packageDocumentation{
		"helix": loadPackageDocumentation(t, filepath.Join(root, "helix"), "helix"),
		"oauth": loadPackageDocumentation(t, filepath.Join(root, "oauth"), "oauth"),
	}
	for _, row := range contract.Rows {
		documentation, ok := docs[row.Package]
		if !ok {
			t.Fatalf("%s: unknown package %q", row.ID, row.Package)
		}
		if !strings.Contains(documentation.doc, row.Anchor) {
			t.Fatalf("%s: missing prose anchor %q in %s docs", row.ID, row.Anchor, row.Package)
		}
		switch row.Kind {
		case "docs-anchor":
		case "exported-symbol":
			if _, ok := documentation.symbols[row.Symbol]; !ok {
				t.Fatalf("%s: exported symbol %q is missing from go/doc", row.ID, row.Symbol)
			}
		case "descriptor-surface":
			verifyDescriptorSurface(t, root, documentation, row)
		default:
			t.Fatalf("%s: unknown contract kind %q", row.ID, row.Kind)
		}
	}
}

func repositoryRoot(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate documentation contract test")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", ".."))
}

func readDocumentationContract(t *testing.T, filename string) documentationContract {
	t.Helper()
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	var contract documentationContract
	if err := json.Unmarshal(data, &contract); err != nil {
		t.Fatal(err)
	}
	if len(contract.Rows) == 0 {
		t.Fatal("documentation contract has no rows")
	}
	return contract
}

func loadPackageDocumentation(t *testing.T, directory, packageName string) packageDocumentation {
	t.Helper()
	fileset := token.NewFileSet()
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	files := make(map[string]*ast.File)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		filename := filepath.Join(directory, entry.Name())
		parsed, parseErr := parser.ParseFile(fileset, filename, nil, parser.ParseComments)
		if parseErr != nil {
			t.Fatal(parseErr)
		}
		if parsed.Name.Name != packageName {
			t.Fatalf("file %s belongs to package %s, want %s", filename, parsed.Name.Name, packageName)
		}
		files[filename] = parsed
	}
	if len(files) == 0 {
		t.Fatalf("package %s not found in %s", packageName, directory)
	}
	parsed := &ast.Package{Name: packageName, Files: files}
	packageDoc := doc.New(parsed, "github.com/kvizyx/twitchy/"+packageName, doc.AllDecls)
	symbols := make(map[string]struct{})
	for _, value := range packageDoc.Consts {
		for _, name := range value.Names {
			symbols[name] = struct{}{}
		}
	}
	for _, value := range packageDoc.Vars {
		for _, name := range value.Names {
			symbols[name] = struct{}{}
		}
	}
	for _, function := range packageDoc.Funcs {
		symbols[function.Name] = struct{}{}
	}
	for _, typ := range packageDoc.Types {
		symbols[typ.Name] = struct{}{}
		for _, function := range typ.Funcs {
			symbols[function.Name] = struct{}{}
		}
		for _, method := range typ.Methods {
			symbols[typ.Name+"."+method.Name] = struct{}{}
		}
	}
	return packageDocumentation{doc: packageDoc.Doc, symbols: symbols}
}

func verifyDescriptorSurface(t *testing.T, root string, documentation packageDocumentation, row contractRow) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "helix", "internal", "manifest", "public-descriptor.json"))
	if err != nil {
		t.Fatal(err)
	}
	var descriptorFile descriptor
	if err := json.Unmarshal(data, &descriptorFile); err != nil {
		t.Fatal(err)
	}
	wantExperimental := row.Stability == "experimental"
	matched := 0
	for _, operation := range descriptorFile.Operations {
		isExperimental := operation.Stability != "stable"
		if isExperimental != wantExperimental {
			continue
		}
		matched++
		methodSymbol := operation.ServiceType + "." + operation.Method
		if _, ok := documentation.symbols[methodSymbol]; !ok {
			t.Fatalf("%s: operation %s is missing from go/doc", row.ID, methodSymbol)
		}
		if _, ok := documentation.symbols[operation.ServiceType]; !ok {
			t.Fatalf("%s: service type %s is missing from go/doc", row.ID, operation.ServiceType)
		}
		if _, ok := documentation.symbols[operation.RequestType]; !ok {
			t.Fatalf("%s: request type %s is missing from go/doc", row.ID, operation.RequestType)
		}
		if _, ok := documentation.symbols[operation.DataType]; !ok {
			serviceType, ok := serviceTypes[operation.ServiceType]
			if !ok {
				t.Fatalf("%s: service type %s is missing from compiled API", row.ID, operation.ServiceType)
			}
			method, ok := serviceType.MethodByName(operation.Method)
			if !ok || method.Type.NumOut() != 2 {
				t.Fatalf("%s: cannot resolve return type for %s", row.ID, methodSymbol)
			}
			actualDataType := responseDataType(method.Type.Out(0))
			if _, ok := documentation.symbols[actualDataType]; !ok {
				t.Fatalf("%s: descriptor data type %s and compiled data type %s are missing from go/doc", row.ID, operation.DataType, actualDataType)
			}
		}
	}
	want := 127
	if wantExperimental {
		want = 22
	}
	if matched != want {
		t.Fatalf("%s: descriptor rows = %d, want %d", row.ID, matched, want)
	}
}

func responseDataType(responseType reflect.Type) string {
	name := responseType.String()
	start := strings.Index(name, "Response[")
	end := strings.LastIndex(name, "]")
	if start < 0 || end <= start {
		return ""
	}
	dataType := name[start+len("Response[") : end]
	if dot := strings.LastIndex(dataType, "."); dot >= 0 {
		dataType = dataType[dot+1:]
	}
	return dataType
}
