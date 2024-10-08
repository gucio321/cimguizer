package cimguizer

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Cimguizer struct {
	Lines []string
}

type Funcs map[string][]Func

type Func struct {
	Args  string `json:"args"`
	ArgsT []ArgT `json:"argsT"`
	Name  string `json:"ov_cimguiname"`
	Ret   string `json:"ret"`
}

type ArgT struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type StructAndEnums struct {
	Structs map[string]Struct `json:"structs"`
	Enums   map[string][]Enum `json:"enums"`
}

type Struct struct{}

type Enum struct {
	Name  string `json:"name"`
	Value string `json:"calc_value"`
}

func Parse(data []byte) (*Cimguizer, error) {
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		// Remove comments
		if strings.HasPrefix(line, "//") && strings.HasPrefix(line, "#include") {
			lines = append(lines[:i], lines[i+1:]...)
		}
	}

	return &Cimguizer{Lines: lines}, nil
}

func (c *Cimguizer) Funcs() (string, error) {
	result := make(Funcs)
	for _, line := range c.Lines {
		// 1.0 skip non-func lines
		if !strings.HasPrefix(line, "CIMGUI_API") {
			continue
		}

		fn := Func{}

		line = strings.TrimPrefix(line, "CIMGUI_API ")
		retType := strings.Split(line, " ")[0]
		fn.Ret = retType

		line = strings.TrimPrefix(line, retType+" ")
		name := strings.Split(line, "(")[0]
		fn.Name = name
		args := strings.TrimSuffix(strings.TrimPrefix(line, name), ";")
		fn.Args = args

		// get argsT
		fn.ArgsT = make([]ArgT, 0)
		argsT := strings.TrimPrefix(strings.TrimSuffix(args, ")"), "(")
		for _, argT := range strings.Split(argsT, ",") {
			if argT == "" {
				continue
			}

			typeName := strings.Split(argT, " ")
			fn.ArgsT = append(fn.ArgsT, ArgT{Name: typeName[1], Type: typeName[0]})
		}

		// save
		if _, ok := result[name]; ok {
			result[name] = append(result[name], fn)
			continue
		}

		result[name] = make([]Func, 1)
		result[name][0] = fn
	}

	resultStr, err := json.MarshalIndent(result, "\t", " ")
	return string(resultStr), err
}

func (c *Cimguizer) StructAndEnums() (string, error) {
	result := &StructAndEnums{}
	// 1.0 structs
	result.Structs = make(map[string]Struct)

	// 2.0 enums
	result.Enums = make(map[string][]Enum)

	for i, line := range c.Lines {
		if !strings.HasPrefix(line, "typedef enum") {
			continue
		}

		data := line
		// the typedef should be foolowed by {
		for n := i + 1; !strings.Contains(data, ";"); n++ {
			data += c.Lines[n]
		}

		// remove duplicated spaces
		data = strings.Join(strings.Fields(data), " ")

		// remove everything from ; to the end of the line
		data = strings.Split(data, ";")[0]
		data = strings.TrimPrefix(data, "typedef enum {")
		data = strings.TrimPrefix(data, " ")
		data = strings.ReplaceAll(data, " ", "")
		name := strings.Split(data, "}")[1]
		values := strings.Split(strings.Split(data, "}")[0], ",")

		enums := make([]Enum, 0)
		for i, value := range values {
			enums = append(enums, Enum{Name: value, Value: fmt.Sprintf("%d", i)})
		}

		result.Enums[name] = enums
	}

	resultStr, err := json.MarshalIndent(result, "\t", " ")
	return string(resultStr), err
}
