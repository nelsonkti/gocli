package cmd

import (
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"gocli/util/helper"
	"gocli/util/template"
	"gocli/util/xfile"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ProtoServerService 定义结构体来存储service信息
type ProtoServerService struct {
	PbPkgName string
	Name      string
	Comment   string
	Methods   []Method
}

type Method struct {
	Name     string
	Comment  string
	Request  string
	Response string
}

func generateRpcServer(fileP string) error {
	// 读取文件内容
	data, err := os.ReadFile(fileP)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}

	protoServices, err := RpcServerDecoder(string(data))
	if err != nil {
		return fmt.Errorf("protobuf Error decoding")
	}
	for _, protoService := range protoServices {
		if protoService.Name == "" {
			continue
		}
		err := generateServerRpcHandler(fileP, protoService)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func generateServerRpcHandler(fileP string, protoService ProtoServerService) error {
	var structInfo protobufStructInfo
	namespace := filepath.Dir(fileP)
	structInfo.Namespace = namespace
	structInfo.Package = filepath.Base(namespace)
	structInfo.PbPkgName = structInfo.Package
	structInfo.ModName = xfile.GetModPath(RelativeSymbol)
	structInfo.StructName = protoService.Name
	structInfo.PbPkgName = protoService.PbPkgName
	structInfo.ProtoService = protoService

	fileOutputPath := namespace + "/"
	newOutPutDir := strings.ReplaceAll(fileOutputPath, "proto", RPCOutPutDir)

	xfile.MkdirAll(newOutPutDir)

	box := packr.New(tmplPath, tmplPath)
	tmpl, _ := box.FindString(rpcTemplateFile)
	err := template.WriteFile(newOutPutDir+helper.ToSnakeCase(protoService.Name)+".go", tmpl, structInfo)
	if err != nil {
		return err
	}

	return nil
}

func RpcServerDecoder(protoContent string) ([]ProtoServerService, error) {
	var services []ProtoServerService

	// 正则表达式来匹配service和rpc方法
	packageRegex := regexp.MustCompile(`option\s+go_package\s*=\s*"([^"]*)";`)
	serviceRegex := regexp.MustCompile(`service\s+(\w+)\s+\{`)
	methodRegex := regexp.MustCompile(`rpc\s+(\w+)\s*\(\s*(\w+)\s*\)\s*returns\s*\(\s*(\w+)\s*\)`)

	// 按行分割Proto内容
	var currentService *ProtoServerService
	lines := strings.Split(protoContent, "\n")
	for k, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 匹配package定义
		packageMatches := packageRegex.FindStringSubmatch(line)
		if len(packageMatches) > 0 {
			pbPkgName := strings.ReplaceAll(packageMatches[1], ".;", "")
			// 初始化 currentService 并设置包名
			if currentService == nil {
				currentService = &ProtoServerService{PbPkgName: pbPkgName}
			} else {
				currentService.PbPkgName = pbPkgName
			}
			continue
		}

		// 匹配service定义
		serviceMatches := serviceRegex.FindStringSubmatch(line)
		if len(serviceMatches) > 0 {
			if currentService != nil {
				services = append(services, *currentService)
			}

			// 增加服务注释
			comment := ""
			commentKey := k - 1
			if commentKey > 0 && lines[commentKey] != "" && strings.Contains(lines[commentKey], "//") {
				comment = strings.ReplaceAll(lines[commentKey], "//", "// "+serviceMatches[1]+"Handler")
			}
			currentService = &ProtoServerService{
				PbPkgName: currentService.PbPkgName, // 保持相同包名
				Name:      capitalize(serviceMatches[1]),
				Comment:   comment,
			}
			continue
		}

		// 匹配rpc方法定义
		methodMatches := methodRegex.FindStringSubmatch(line)
		if len(methodMatches) > 0 && currentService != nil {
			// 增加服务注释
			methodComment := ""
			commentKey := k - 1
			if commentKey > 0 && lines[commentKey] != "" && strings.Contains(lines[commentKey], "//") {
				methodComment = strings.ReplaceAll(lines[commentKey], "//", "// "+methodMatches[1])
				methodComment = strings.TrimLeft(methodComment, " ")
			}
			method := Method{
				Name:     methodMatches[1],
				Comment:  methodComment,
				Request:  methodMatches[2],
				Response: methodMatches[3],
			}
			currentService.Methods = append(currentService.Methods, method)
		}
	}

	// 添加最后一个service
	if currentService != nil {
		services = append(services, *currentService)
	}

	return services, nil
}
