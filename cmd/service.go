package cmd

import (
	"bufio"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/cobra"
	"gocli/util/helper"
	"gocli/util/template"
	"gocli/util/xfile"
	"gocli/util/xprintf"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

const (
	serviceType         = "service"
	serviceTemplateFile = "service.tmpl"
	DIFilepath          = "internal/container/service/service_container.go"
)

func init() {

	Cmd.AddCommand(serviceCommand)

	serviceCommand.Flags().StringVarP(&name, "name", "n", "", NameTips)
	serviceCommand.Flags().StringVarP(&fileName, "file_name", "f", "", FileNameTips)
	serviceCommand.Flags().StringVarP(&path, "path", "p", "", PathTips)
}

var serviceCommand = &cobra.Command{
	Use:   "make:service",
	Short: "create a new service file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		zhName = ServiceZhName

		createService()
	},
}

func createService() {
	fmt.Println(xprintf.Blue("creating service file ..."))
	newPath := getDirPath(ServiceDirPath, path)

	fileNames := stringToSplit(fileName)
	isOverwrite := checkFileExists(newPath, fileNames, serviceType)

	for _, f := range fileNames {
		filePath := getGenerateFilePath(newPath, f, serviceType)
		exists, _ := xfile.PathExists(filePath)
		if exists && isOverwrite == false {
			continue
		}
		generateService(f, filePath, newPath)
	}

	fmt.Println(xprintf.Green("create a new 【service】 file success\n"))

}

// generateRepository 生成 repository
func generateService(fileName string, outputFilePath, newPath string) {
	var structInfo serviceStructInfo
	structInfo.ModName = xfile.GetModPath(RelativeSymbol)
	structInfo.Package = xfile.GetPackageName(newPath)
	structName := helper.SnakeToCamel(fileName)
	structInfo.StructName = structName + helper.Capitalize(serviceType)
	structInfo.StructComment = name + zhName
	namespacePath := structInfo.ModName + strings.Replace(newPath, RelativeSymbol, "/", -1)

	xfile.MkdirAll(newPath)

	box := packr.New(tmplPath, tmplPath)
	tmpl, _ := box.FindString(serviceTemplateFile)
	err := template.WriteFile(outputFilePath, tmpl, structInfo)
	if err != nil {
		fmt.Println(xprintf.Red(fmt.Sprintf("生成 %s 失败: %+v", outputFilePath, err)))
		return
	}

	err = diService(DIFilepath, structName, namespacePath)
	if err != nil {
		fmt.Println(xprintf.Red(fmt.Sprintf("更新 %s 失败: %+v", DIFilepath, err)))
		return
	}

	fmt.Println(xprintf.Green("CREATED ") + outputFilePath)
}

// diService 依赖注入服务
func diService(filePath, newService, newServicePath string) error {
	// 读取文件内容
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	defer file.Close()

	packageName := filepath.Base(newServicePath)
	// 新的导入语句和结构体字段
	newImport := fmt.Sprintf("\t\"%s\"\n", newServicePath)

	originalContainerField := fmt.Sprintf("%sService %s.%sServiceImpl", capitalize(newService), packageName, capitalize(newService))
	containerField := fmt.Sprintf("\t%s  // %s\n", originalContainerField, name)

	originalRegisterField := fmt.Sprintf("%sService: %s.New%sService(svc)", capitalize(newService), packageName, capitalize(newService))
	registerField := fmt.Sprintf("\t\t%s, // %s\n", originalRegisterField, name)

	// 读取文件内容并处理
	var updatedContent strings.Builder
	scanner := bufio.NewScanner(file)
	importFound := false
	containerFound := false
	registerContainerFound := false

	importExists := false
	containerExists := false
	registerContainerExists := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, newServicePath) {
			importExists = true
		}

		if strings.Contains(line, "import (") {
			importFound = true
		}
		if importFound && strings.TrimSpace(line) == ")" {
			if !importExists {
				updatedContent.WriteString(newImport)
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

// isStrEqual 判断两个字符串是否相等
func isStrEqual(oldStr, str string) bool {
	oldStr = strings.ReplaceAll(oldStr, " ", "")
	str = strings.ReplaceAll(str, " ", "")
	return strings.Contains(oldStr, str)
}

// capitalize 首字母大写
func capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
