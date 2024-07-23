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
)

var (
	fileName = ""
)

const (
	FileNameTips = "请输入 [表名]"
)

const (
	repositoryType         = "repository"
	repositoryTemplateFile = "repository.tmpl"
	fileIOCClassType       = "model"
	DIRepositoryFilepath   = "internal/container/repository/repository_container.go"
)

func init() {

	Cmd.AddCommand(repositoryCommand)

	repositoryCommand.Flags().StringVarP(&name, "name", "n", "", NameTips)
	repositoryCommand.Flags().StringVarP(&fileName, "file_name", "f", "", FileNameTips)
	repositoryCommand.Flags().StringVarP(&path, "path", "p", "", PathTips)
}

var repositoryCommand = &cobra.Command{
	Use:   "make:repository",
	Short: "create a new repository file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		zhName = RepositoryZhName

		createRepository()
	},
}

func createRepository() {
	fmt.Println(xprintf.Blue("creating repository file ..."))
	newPath := getDirPath(RepositoryDirPath, repositoryType, path)
	fileNames := stringToSplit(fileName)
	isOverwrite := checkFileExists(newPath, fileNames, repositoryType)

	for _, f := range fileNames {
		filePath := getGenerateFilePath(newPath, f, repositoryType)
		exists, _ := xfile.PathExists(filePath)
		if exists && isOverwrite == false {
			continue
		}
		generateRepository(f, filePath, newPath)
	}
	fmt.Println(xprintf.Green("create a new 【repository】 file success\n"))
}

// generateRepository 生成 repository
func generateRepository(fileName string, outputFilePath, newPath string) {
	var structInfo repositoryStructInfo
	structInfo.ModName = xfile.GetModPath(RelativeSymbol)
	structInfo.Package = xfile.GetPackageName(newPath)
	structName := helper.SnakeToCamel(fileName)
	structInfo.StructName = structName + helper.Capitalize(repositoryType)
	var modelNamespace = newPath

	// 检查是否存在
	if strings.Contains(newPath, repositoryType) {
		modelNamespace = strings.Replace(newPath, repositoryType, fileIOCClassType, -1)
	}

	// 去掉路径
	if strings.Contains(modelNamespace, RelativeSymbol) {
		modelNamespace = strings.Replace(modelNamespace, RelativeSymbol, "", -1)
	}
	structInfo.IOCNamespace = modelNamespace

	structInfo.IOCPackage = xfile.GetPackageName(modelNamespace)
	structInfo.IOCStructName = structName + helper.Capitalize(fileIOCClassType)
	structInfo.StructComment = name + zhName

	xfile.MkdirAll(newPath)

	box := packr.New(tmplPath, tmplPath)
	tmpl, _ := box.FindString(repositoryTemplateFile)
	err := template.WriteFile(outputFilePath, tmpl, structInfo)
	if err != nil {
		fmt.Println(xprintf.Red(fmt.Sprintf("生成 %s 失败: %+v", outputFilePath, err)))
		return
	}

	err = diRepository(DIRepositoryFilepath, structName, newPath, structInfo)
	if err != nil {
		fmt.Println(xprintf.Red(fmt.Sprintf("更新 %s 失败: %+v", DIServiceFilepath, err)))
		return
	}

	fmt.Println(xprintf.Green("CREATED ") + outputFilePath)
}

// diRepository 依赖注入仓库
func diRepository(filePath, newService, newServicePath string, structInfo repositoryStructInfo) error {
	// 读取文件内容
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return nil
	}
	defer file.Close()

	packageName := filepath.Base(newServicePath)
	// 新的导入语句和结构体字段
	newImport := fmt.Sprintf("\t\"%s/%s\"\n", structInfo.ModName, strings.ReplaceAll(newServicePath, RelativeSymbol, ""))
	newModImport := fmt.Sprintf("\t\"%s/%s\"\n", structInfo.ModName, structInfo.IOCNamespace)

	originalContainerField := fmt.Sprintf("%sRepository *%s.%sRepository", capitalize(newService), packageName, capitalize(newService))
	containerField := fmt.Sprintf("\t%s  // %s\n", originalContainerField, name)

	registerModelField := fmt.Sprintf("%s.New%s(db)", structInfo.IOCPackage, structInfo.IOCStructName)
	originalRegisterField := fmt.Sprintf("%sRepository: %s.New%sRepository(%s, log)", capitalize(newService), packageName, capitalize(newService), registerModelField)
	registerField := fmt.Sprintf("\t\t%s, // %s\n", originalRegisterField, name)

	// 读取文件内容并处理
	var updatedContent strings.Builder
	scanner := bufio.NewScanner(file)
	importFound := false
	containerFound := false
	registerContainerFound := false

	importExists := false
	importModelExists := false
	containerExists := false
	registerContainerExists := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, newServicePath) {
			importExists = true
		}

		if strings.Contains(line, structInfo.IOCNamespace) {
			importModelExists = true
		}

		if strings.Contains(line, "import (") {
			importFound = true
		}
		if importFound && strings.TrimSpace(line) == ")" {
			if !importExists {
				updatedContent.WriteString(newImport)
			}
			if !importModelExists {
				updatedContent.WriteString(newModImport)
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
