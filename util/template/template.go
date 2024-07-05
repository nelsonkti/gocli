package template

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"
)

func Generate(tplPath, templateName string, outputPath string, data any) {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer outputFile.Close()
	tmpl, err := template.New(templateName).Funcs(template.FuncMap{"lower": strings.ToLower}).ParseFiles(tplPath)
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	// 执行模板，生成代码
	if err := tmpl.Execute(outputFile, data); err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}
}

func WriteFile(path, tmpl string, args interface{}) (err error) {
	data, err := ParseTmpl(tmpl, args)
	if err != nil {
		return
	}

	return os.WriteFile(path, data, 0755)
}

func ParseTmpl(tmpl string, args interface{}) ([]byte, error) {
	tmp, err := template.New("").Parse(tmpl)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err = tmp.Execute(&buf, args); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil

}
