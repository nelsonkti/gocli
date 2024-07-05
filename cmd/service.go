package cmd

import (
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/cobra"
	"gocli/util/helper"
	"gocli/util/template"
	"gocli/util/xfile"
	"gocli/util/xprintf"
)

const (
	serviceType         = "service"
	serviceTemplateFile = "service.tmpl"
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

	xfile.MkdirAll(newPath)

	box := packr.New(tmplPath, tmplPath)
	tmpl, _ := box.FindString(serviceTemplateFile)
	err := template.WriteFile(outputFilePath, tmpl, structInfo)
	if err != nil {
		fmt.Println(xprintf.Red(fmt.Sprintf("生成 %s 失败: %+v", outputFilePath, err)))
		return
	}
	fmt.Println(xprintf.Green("CREATED ") + outputFilePath)
}
