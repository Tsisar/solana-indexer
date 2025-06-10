package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type IDL struct {
	Types  []TypeDef  `json:"types"`
	Events []EventDef `json:"events"`
}

type TypeDef struct {
	Name string   `json:"name"`
	Type StructTy `json:"type"`
}

type StructTy struct {
	Kind   string     `json:"kind"`
	Fields []FieldDef `json:"fields"`
}

type FieldDef struct {
	Name string          `json:"name"`
	Ty   json.RawMessage `json:"type"`
}

type EventDef struct {
	Name string `json:"name"`
}

// snakeToCamel converts snake_case to CamelCase
func snakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + strings.ToLower(p[1:])
	}
	return strings.Join(parts, "")
}

// coreType maps IDL types to Go types
func coreType(raw json.RawMessage) (string, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		switch s {
		case "pubkey":
			return "solana.PublicKey", nil
		case "u8":
			return "uint8", nil
		case "i8":
			return "int8", nil
		case "u16":
			return "uint16", nil
		case "i16":
			return "int16", nil
		case "u32":
			return "uint32", nil
		case "i32":
			return "int32", nil
		case "u64":
			return "uint64", nil
		case "i64":
			return "int64", nil
		case "bool":
			return "bool", nil
		case "string":
			return "string", nil
		case "u128", "i128":
			return "[16]byte", nil
		default:
			return "", fmt.Errorf("unknown basic type: %s", s)
		}
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(raw, &obj); err != nil {
		return "", err
	}
	if arr, ok := obj["array"]; ok {
		items := arr.([]interface{})
		et, err := coreType(json.RawMessage(fmt.Sprintf(`"%s"`, items[0].(string))))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("[%d]%s", int(items[1].(float64)), et), nil
	}
	if def, ok := obj["defined"]; ok {
		return def.(map[string]interface{})["name"].(string), nil
	}
	if vec, ok := obj["vec"]; ok {
		switch v := vec.(type) {
		case string:
			et, err := coreType(json.RawMessage(fmt.Sprintf(`"%s"`, v)))
			if err != nil {
				return "", err
			}
			return "[]" + et, nil
		case map[string]interface{}:
			inner := v["defined"].(map[string]interface{})["name"].(string)
			return "[]" + inner, nil
		}
	}
	return "", fmt.Errorf("unsupported type: %v", obj)
}

// subgraphType maps IDL types to Go types
func subgraphType(raw json.RawMessage) (string, error) {
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		switch s {
		case "pubkey":
			return "solana.PublicKey", nil
		case "u8":
			return "big.Int", nil
		case "i8":
			return "big.Int", nil
		case "u16":
			return "big.Int", nil
		case "i16":
			return "big.Int", nil
		case "u32":
			return "big.Int", nil
		case "i32":
			return "big.Int", nil
		case "u64":
			return "big.Int", nil
		case "i64":
			return "big.Int", nil
		case "bool":
			return "bool", nil
		case "string":
			return "string", nil
		case "u128", "i128":
			return "[16]byte", nil
		default:
			return "", fmt.Errorf("unknown basic type: %s", s)
		}
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(raw, &obj); err != nil {
		return "", err
	}
	if arr, ok := obj["array"]; ok {
		items := arr.([]interface{})
		et, err := coreType(json.RawMessage(fmt.Sprintf(`"%s"`, items[0].(string))))
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("[%d]%s", int(items[1].(float64)), et), nil
	}
	if def, ok := obj["defined"]; ok {
		return def.(map[string]interface{})["name"].(string), nil
	}
	if vec, ok := obj["vec"]; ok {
		switch v := vec.(type) {
		case string:
			et, err := coreType(json.RawMessage(fmt.Sprintf(`"%s"`, v)))
			if err != nil {
				return "", err
			}
			return "[]" + et, nil
		case map[string]interface{}:
			inner := v["defined"].(map[string]interface{})["name"].(string)
			return "[]" + inner, nil
		}
	}
	return "", fmt.Errorf("unsupported type: %v", obj)
}

// idlTypeMap returns the StructTy for a given event name.
func idlTypeMap(name string, types []TypeDef) StructTy {
	for _, t := range types {
		if t.Name == name {
			return t.Type
		}
	}
	return StructTy{}
}

func main() {
	idlDir := "idl"
	eventsDir := "generator/core/events"
	mapingDir := "generator/subgraph/maping"

	allEvents, allEventToFunc := processIdlDirectory(idlDir, eventsDir, mapingDir)
	generateEventRegistry(eventsDir, allEvents)
	generateMapperRegistry(mapingDir, allEventToFunc)
}

func processIdlDirectory(idlDir, eventsDir, mapingDir string) ([]string, []string) {
	var allEvents, allEventToFunc []string

	entries, err := os.ReadDir(idlDir)
	if err != nil {
		fmt.Printf("failed to read directory %s: %v", idlDir, err)
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}

		path := filepath.Join(idlDir, e.Name())
		idl := readIdlFile(path)
		idlName := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))

		evs := generateEventStructs(eventsDir, idlName, idl)
		allEvents = append(allEvents, evs...)

		funcMap := generateEventMappers(mapingDir, idlName, idl)
		allEventToFunc = append(allEventToFunc, funcMap...)

		generateSubgraphEvents("subgraph/events", idlName, idl)
	}

	return allEvents, allEventToFunc
}

func readIdlFile(path string) IDL {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("read %s error: %v", path, err)
		return IDL{}
	}

	var idl IDL
	if err := json.Unmarshal(data, &idl); err != nil {
		fmt.Printf("unmarshal %s error: %v", path, err)
		return IDL{}
	}

	return idl
}

func generateEventStructs(eventsDir, idlName string, idl IDL) []string {
	var b strings.Builder
	var allEvents []string

	b.WriteString("package events\n\n")
	b.WriteString("// Code generated by generate_events.go; DO NOT EDIT.\n\n")
	b.WriteString("import \"github.com/gagliardetto/solana-go\"\n\n")

	for _, ev := range idl.Events {
		b.WriteString(fmt.Sprintf("// %s event struct\n", ev.Name))
		b.WriteString(fmt.Sprintf("type %s struct {\n", ev.Name))

		structDef := idlTypeMap(ev.Name, idl.Types)
		for _, f := range structDef.Fields {
			gt, err := coreType(f.Ty)
			if err != nil {
				fmt.Printf("type %s field %s error: %v", ev.Name, f.Name, err)
				continue
			}
			fieldName := snakeToCamel(f.Name)
			b.WriteString(fmt.Sprintf("    %s %s `borsh:\"%s\"`\n", fieldName, gt, f.Name))
		}
		b.WriteString("}\n\n")
		allEvents = append(allEvents, ev.Name)
	}

	writeToFile(filepath.Join(eventsDir, idlName+".go"), b.String())
	return allEvents
}

func generateEventMappers(mapingDir, idlName string, idl IDL) []string {
	var m strings.Builder
	var mappings []string

	m.WriteString("package maping\n\nimport (\n")
	m.WriteString("\t\"context\"\n\t\"encoding/json\"\n\t\"fmt\"\n")
	m.WriteString("\t\"github.com/Tsisar/extended-log-go/log\"\n")
	m.WriteString("\t\"github.com/Tsisar/solana-indexer/storage/model/core\"\n")
	m.WriteString("\t\"github.com/Tsisar/solana-indexer/subgraph/events\"\n")
	m.WriteString("\t\"gorm.io/gorm\"\n)\n\n")

	for _, ev := range idl.Events {
		m.WriteString(fmt.Sprintf("func map%s(ctx context.Context, db *gorm.DB, event core.Event) error {\n", ev.Name))
		m.WriteString(fmt.Sprintf("\tlog.Infof(\"[mapping] %s: %%s\", event.TransactionSignature)\n", ev.Name))
		m.WriteString(fmt.Sprintf("\tvar ev events.%s\n", ev.Name))
		m.WriteString("\tif err := json.Unmarshal(event.JsonEv, &ev); err != nil {\n")
		m.WriteString(fmt.Sprintf("\t\treturn fmt.Errorf(\"failed to decode %s: %%w\", err)\n", ev.Name))
		m.WriteString("\t}\n\t// TODO: implement mapping logic\n\treturn nil\n}\n\n")

		mappings = append(mappings, fmt.Sprintf("\t\"%s\": map%s,", ev.Name, ev.Name))
	}

	writeToFile(filepath.Join(mapingDir, idlName+".go"), m.String())
	return mappings
}

func generateEventRegistry(eventsDir string, eventNames []string) {
	var r strings.Builder

	r.WriteString("package events\n\n// Code generated by generate_events.go; DO NOT EDIT.\n\n")
	r.WriteString("import \"github.com/near/borsh-go\"\n\ntype EventDecoder func([]byte) (any, error)\n\n")
	r.WriteString("func decode[T any](data []byte) (any, error) {\n")
	r.WriteString("    var out T\n    err := borsh.Deserialize(&out, data)\n    return out, err\n}\n\n")
	r.WriteString("var Registry = map[string]EventDecoder{\n")
	for _, name := range eventNames {
		r.WriteString(fmt.Sprintf("    \"%s\": decode[%s],\n", name, name))
	}
	r.WriteString("}\n")

	writeToFile(filepath.Join(eventsDir, "registry.go"), r.String())
}

func generateMapperRegistry(mapingDir string, mappings []string) {
	var r strings.Builder

	r.WriteString(`package maping

import (
	"context"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/storage/model/core"
	"gorm.io/gorm"
)

type EventMapper func(ctx context.Context, db *gorm.DB, event core.Event) error

var registry = map[string]EventMapper{
`)
	r.WriteString(strings.Join(mappings, "\n"))
	r.WriteString(`
}

func mapEvents(ctx context.Context, db *gorm.DB, event core.Event) error {
	if handler, ok := registry[event.Name]; ok {
		return handler(ctx, db, event)
	}
	log.Warnf("No mapping implemented for event: %s", event.Name)
	return nil
}
`)
	writeToFile(filepath.Join(mapingDir, "registry.go"), r.String())
}

func generateSubgraphEvents(subgraphEventsDir, idlName string, idl IDL) {
	var b strings.Builder
	b.WriteString("package events\n\n")
	b.WriteString("// Code generated by generate_events.go; DO NOT EDIT.\n\n")
	b.WriteString("import \"github.com/gagliardetto/solana-go\"\n")
	b.WriteString("import \"math/big\"\n\n")

	for _, ev := range idl.Events {
		b.WriteString(fmt.Sprintf("// %s event struct\n", ev.Name))
		b.WriteString(fmt.Sprintf("type %s struct {\n", ev.Name))
		structDef := idlTypeMap(ev.Name, idl.Types)
		for _, f := range structDef.Fields {
			gt, err := subgraphType(f.Ty)
			if err != nil {
				logErrorf("type %s field %s error: %v", ev.Name, f.Name, err)
				continue
			}
			fieldName := snakeToCamel(f.Name)
			b.WriteString(fmt.Sprintf("    %s %s\n", fieldName, gt))
		}
		b.WriteString("}\n\n")
	}

	outFile := filepath.Join(subgraphEventsDir, idlName+".go")
	writeToFile(outFile, b.String())
}

func writeToFile(path, content string) {
	fmt.Printf("Generating %s...\n", path)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		logErrorf("write %s error: %v", path, err)
	}
}

func logErrorf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
}

func logFatalf(format string, a ...any) {
	logErrorf(format, a...)
	os.Exit(1)
}
