package cmd

import (
	"bufio"
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

const RpcClientIocFile = "internal/container/grpc/grpc.go"

// ProtoClientService 定义结构体来存储service信息
type ProtoClientService struct {
	PbPkgName string
	Name      string
	Comment   string
}

func generateRpcClient(fileP string) error {
	// 读取文件内容
	data, err := os.ReadFile(fileP)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}

	protoServices, err := RpcClientDecoder(string(data))
	if err != nil {
		return fmt.Errorf("protobuf Error decoding")
	}
	err = generateClientRpcHandler(fileP, protoServices)
	if err != nil {
		return err
	}
	return nil
}

func generateClientRpcHandler(fileP string, protoService []ProtoClientService) error {
	if len(protoService) == 0 {
		return nil
	}
	var newProtoService []ProtoClientService

	for _, data := range protoService {
		if data.Name == "" {
			continue
		}
		newProtoService = append(newProtoService, data)
	}

	var structInfo protobufStructInfo
	namespace := filepath.Dir(fileP)
	structInfo.Namespace = namespace
	structInfo.Package = filepath.Base(namespace)
	structInfo.PbPkgName = structInfo.Package
	structInfo.ModName = xfile.GetModPath(RelativeSymbol)
	structInfo.StructName = capitalize(filepath.Base(namespace))
	structInfo.ProtoClientService = newProtoService

	fileOutputPath := namespace + "/"
	newOutPutDir := strings.ReplaceAll(fileOutputPath, "proto", RPCClientOutPutDir)

	xfile.MkdirAll(newOutPutDir)

	box := packr.New(tmplPath, tmplPath)
	tmpl, _ := box.FindString(rpcClientTemplateFile)
	outFilePath := newOutPutDir + helper.ToSnakeCase(structInfo.Package) + ".go"
	err := template.WriteFile(outFilePath, tmpl, structInfo)
	if err != nil {
		return err
	}

	err = GrpcClientIoc(RpcClientIocFile, structInfo.StructName, structInfo.ModName, newOutPutDir)
	return err
}

func RpcClientDecoder(protoContent string) ([]ProtoClientService, error) {
	var services []ProtoClientService

	// 正则表达式来匹配service和rpc方法
	packageRegex := regexp.MustCompile(`option\s+go_package\s*=\s*"([^"]*)";`)
	serviceRegex := regexp.MustCompile(`service\s+(\w+)\s+\{`)

	// 按行分割Proto内容
	var currentService *ProtoClientService
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
				currentService = &ProtoClientService{PbPkgName: pbPkgName}
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
			currentService = &ProtoClientService{
				PbPkgName: currentService.PbPkgName, // 保持相同包名
				Name:      capitalize(serviceMatches[1]),
				Comment:   comment,
			}
			continue
		}
	}

	// 添加最后一个service
	if currentService != nil {
		services = append(services, *currentService)
	}

	return services, nil
}

func GrpcClientIoc(filePath, clientName, modName, clientPath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	defer file.Close()

	namespace := modName + "/" + clientPath
	namespaceLen := len(namespace)
	if namespaceLen > 0 && namespace[namespaceLen-1:] == SlashSymbol {
		namespace = namespace[:namespaceLen-1]
	}
	cPackage := filepath.Base(clientPath)
	structField := capitalize(clientName)
	fmt.Println(cPackage, structField, namespace)

	// 新的导入语句和结构体字段
	newModImport := fmt.Sprintf("\t\"%s\"\n", namespace)
	newRpcNamespace := "mashang/pkg/grpc"
	newRpcNamespaceWirte := fmt.Sprintf("\t\"%s\"\n", newRpcNamespace)

	originalContainerField := fmt.Sprintf("%s %s.Client", structField, cPackage)
	containerField := fmt.Sprintf("\t%s %s.Client\n", structField, cPackage)

	originalRegisterField := fmt.Sprintf("%s: %s.New%s(grpc.NewClient(nil, ctx)),", structField, cPackage, structField)
	registerField := fmt.Sprintf("\t\t%s: %s.New%s(grpc.NewClient(nil, ctx)),\n", structField, cPackage, structField)

	var updatedContent strings.Builder
	scanner := bufio.NewScanner(file)
	importFound := false
	containerFound := false
	registerContainerFound := false

	importExists := false
	importRpcExists := false
	containerExists := false
	registerContainerExists := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, cPackage) {
			importExists = true
		}

		if strings.Contains(line, newRpcNamespace) {
			importRpcExists = true
		}

		if strings.Contains(line, "import (") {
			importFound = true
		}
		if importFound && strings.TrimSpace(line) == ")" {
			if !importExists {
				updatedContent.WriteString(newModImport)
			}

			if !importRpcExists {
				updatedContent.WriteString(newRpcNamespaceWirte)
			}

			importFound = false
		}

		if isStrEqual(line, originalContainerField) {
			containerExists = true
		}

		if strings.Contains(line, "type Container struct {") {
			containerFound = true
		}
		if containerFound && strings.TrimSpace(line) == "}" {
			if !containerExists {
				updatedContent.WriteString(containerField)
			}
			containerFound = false
		}

		if isStrEqual(line, originalRegisterField) {
			registerContainerExists = true
		}

		if strings.Contains(line, "return &Container{") {
			registerContainerFound = true
		}
		if registerContainerFound && strings.TrimSpace(line) == "}" {
			if !registerContainerExists {
				updatedContent.WriteString(registerField)
			}
			registerContainerFound = false
		}

		updatedContent.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return nil
	}

	// 将更新后的内容写回原始文件
	err = os.WriteFile(filePath, []byte(updatedContent.String()), 0644)
	if err != nil {
		fmt.Printf("Error writing updated file: %v\n", err)
		return nil
	}

	return nil
}
