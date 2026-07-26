package manifest

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"regexp"
	"strings"
	"testing"
)

type expectedOperation struct {
	Anchor string `json:"anchor"`
	Group  string `json:"group"`
	Label  string `json:"label"`
	Method string `json:"method"`
	Name   string `json:"name"`
	Path   string `json:"path"`
}

type expectedOperations struct {
	Operations []expectedOperation `json:"operations"`
}

type sourceOperation struct {
	Anchor string
	Group  string
	Label  string
	Method string
	Name   string
	Path   string
}

func TestSourceToManifest(t *testing.T) {
	bytes, err := os.ReadFile("sources/reference.html")
	if err != nil {
		t.Fatal(err)
	}
	var expected expectedOperations
	if err := json.Unmarshal(mustRead(t, "expected-operations.json"), &expected); err != nil {
		t.Fatal(err)
	}
	parsed, err := parseReference(string(bytes), expected.Operations)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != len(expected.Operations) {
		t.Fatalf("source parser returned %d operations, want %d", len(parsed), len(expected.Operations))
	}
	for index, want := range expected.Operations {
		got := parsed[index]
		if got != (sourceOperation{want.Anchor, want.Group, want.Label, want.Method, want.Name, want.Path}) {
			t.Errorf("operation %d: got %#v, want %#v", index, got, want)
		}
	}
	manifest, err := LoadManifest()
	if err != nil {
		t.Fatal(err)
	}
	byAnchor := make(map[string]Operation, len(manifest.Operations))
	for _, operation := range manifest.Operations {
		byAnchor[operation.Anchor] = operation
	}
	for _, parsedOperation := range parsed {
		operation, ok := byAnchor[parsedOperation.Anchor]
		if !ok {
			t.Errorf("manifest missing source anchor %q", parsedOperation.Anchor)
			continue
		}
		if operation.Group != parsedOperation.Group || operation.Stability != Stability(parsedOperation.Label) || operation.Method != parsedOperation.Method || operation.Path != parsedOperation.Path {
			t.Errorf("manifest/source mismatch for %q: %#v vs %#v", parsedOperation.Anchor, operation, parsedOperation)
		}
	}
}

func parseReference(source string, expected []expectedOperation) ([]sourceOperation, error) {
	firstOperation := strings.Index(source, `<h2 id="`)
	if firstOperation < 0 {
		return nil, fmt.Errorf("reference has no operation headings")
	}
	groups := parseGroups(source[:firstOperation])
	headingRE := regexp.MustCompile(`(?is)<h2\s+id="([^"]+)"[^>]*>(.*?)</h2>`)
	headings := headingRE.FindAllStringSubmatchIndex(source, -1)
	sections := make(map[string]string, len(headings))
	for index, heading := range headings {
		end := len(source)
		if index+1 < len(headings) {
			end = headings[index+1][0]
		}
		sections[source[heading[2]:heading[3]]] = source[heading[0]:end]
	}
	urlRE := regexp.MustCompile(`(?is)<h3[^>]*>\s*URL\s*</h3>.*?\b(GET|POST|PUT|PATCH|DELETE|HEAD)\s+https?://[^<\s]+(/helix/[^<\s]+)`)
	pillRE := regexp.MustCompile(`(?is)class="[^"]*pill-(new|beta)`)
	parsed := make([]sourceOperation, 0, len(expected))
	for _, operation := range expected {
		section, ok := sections[operation.Anchor]
		if !ok {
			return nil, fmt.Errorf("source missing anchor %q", operation.Anchor)
		}
		url := urlRE.FindStringSubmatch(section)
		if len(url) != 3 {
			return nil, fmt.Errorf("source missing URL for %q", operation.Anchor)
		}
		label := "stable"
		intro := section
		if authorization := strings.Index(strings.ToLower(section), "<h3>authorization</h3>"); authorization >= 0 {
			intro = section[:authorization]
		}
		if pill := pillRE.FindStringSubmatch(intro); len(pill) == 2 {
			label = strings.ToUpper(pill[1])
		}
		parsed = append(parsed, sourceOperation{
			Anchor: operation.Anchor,
			Group:  groups[operation.Anchor],
			Label:  label,
			Method: strings.ToUpper(url[1]),
			Name:   cleanText(sectionHeading(section)),
			Path:   html.UnescapeString(url[2]),
		})
	}
	return parsed, nil
}

func parseGroups(source string) map[string]string {
	groups := make(map[string]string)
	rowRE := regexp.MustCompile(`(?is)<tr[^>]*>(.*?)</tr>`)
	linkRE := regexp.MustCompile(`(?is)<a[^>]+href="#([^"]+)"[^>]*>.*?</a>`)
	cellRE := regexp.MustCompile(`(?is)<td[^>]*>(.*?)</td>`)
	for _, row := range rowRE.FindAllStringSubmatch(source, -1) {
		link := linkRE.FindStringSubmatch(row[1])
		cells := cellRE.FindStringSubmatch(row[1])
		if len(link) == 2 && len(cells) == 2 {
			groups[link[1]] = cleanText(cells[1])
		}
	}
	return groups
}

func sectionHeading(section string) string {
	match := regexp.MustCompile(`(?is)<h2[^>]*>(.*?)</h2>`).FindStringSubmatch(section)
	if len(match) != 2 {
		return ""
	}
	return match[1]
}

func cleanText(value string) string {
	value = regexp.MustCompile(`(?is)<[^>]+>`).ReplaceAllString(value, " ")
	return strings.Join(strings.Fields(html.UnescapeString(value)), " ")
}

func mustRead(t *testing.T, name string) []byte {
	t.Helper()
	bytes, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}
