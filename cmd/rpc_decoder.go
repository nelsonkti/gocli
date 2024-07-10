package cmd

import (
	"regexp"
	"strings"
)

// ProtoService 定义结构体来存储service信息
type ProtoService struct {
	PbPkgName string
	Name      string
	Methods   []Method
}

type Method struct {
	Name     string
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
	for _, line := range lines {
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
			currentService = &ProtoService{
				PbPkgName: currentService.PbPkgName, // 保持相同包名
				Name:      serviceMatches[1],
			}
			continue
		}

		// 匹配rpc方法定义
		methodMatches := methodRegex.FindStringSubmatch(line)
		if len(methodMatches) > 0 && currentService != nil {
			method := Method{
				Name:     methodMatches[1],
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
