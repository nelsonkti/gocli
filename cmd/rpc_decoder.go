package cmd

import (
	"regexp"
	"strings"
)

// ProtoService 定义结构体来存储service信息
type ProtoService struct {
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

func RpcDecoder(protoContent string) ([]ProtoService, error) {
	var services []ProtoService

	// 正则表达式来匹配service和rpc方法
	packageRegex := regexp.MustCompile(`option\s+go_package\s*=\s*"([^"]*)";`)
	serviceRegex := regexp.MustCompile(`service\s+(\w+)\s+\{`)
	methodRegex := regexp.MustCompile(`rpc\s+(\w+)\s*\(\s*(\w+)\s*\)\s*returns\s*\(\s*(\w+)\s*\)`)

	// 按行分割Proto内容
	var currentService *ProtoService
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
				currentService = &ProtoService{PbPkgName: pbPkgName}
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
			currentService = &ProtoService{
				PbPkgName: currentService.PbPkgName, // 保持相同包名
				Name:      serviceMatches[1],
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
