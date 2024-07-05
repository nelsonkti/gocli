package cmd

import (
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/cobra"
	"gocli/util/helper"
	"gocli/util/template"
	"gocli/util/xfile"
	"gocli/util/xprintf"
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
	newPath := getDirPath(RepositoryDirPath, path)
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
	fmt.Println(xprintf.Green("CREATED ") + outputFilePath)
}
