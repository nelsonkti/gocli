package cmd

import (
	"bufio"
	"fmt"
	"gocli/util/xfile"
	"gocli/util/xprintf"
	"os"
	"strings"
)

const (
	GoFileFormat    = ".go"
	RelativeSymbol  = "./"
	UnderlineSymbol = "_"
	SlashSymbol     = "/"
	LineBreakSymbol = "\n"
)

const (
	ModelZhName      = "模型"
	RepositoryZhName = "仓库"
	ServiceZhName    = "服务"

	ModelDirPath      = "./internal/model"
	RepositoryDirPath = "./internal/repository"
	ServiceDirPath    = "./internal/service"
)

type repositoryStructInfo struct {
	ModName       string // mod 名称
	Package       string // 包名
	StructName    string // 结构体名称
	StructComment string // 结构体注释
	IOCStructName string // 模型结构体名称
	IOCNamespace  string // 模型命名空间
	IOCPackage    string // 模型包名
}

type serviceStructInfo struct {
	ModName       string // mod 名称
	Package       string // 包名
	StructName    string // 结构体名称
	StructComment string // 结构体注释
}

// checkFileExists 检查文件是否存在
func checkFileExists(path string, tableNames []string, types string) bool {
	var fileInfo string
	for _, tab := range tableNames {
		filePath := getGenerateFilePath(path, tab, types)
		exists, err := xfile.PathExists(filePath)
		if err != nil {
			panic(xprintf.Red(err.Error()))
		}

		if exists {
			fileInfo += filePath + LineBreakSymbol
		}
	}

	var isOverwrite = true
	if fileInfo != "" {
		fmt.Println(xprintf.Yellow("以下文件已存在："))
		fmt.Println(fileInfo)
		fmt.Printf(xprintf.Blue("是否覆盖？(y/n): "))

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if strings.ToLower(input) != "y" {
			isOverwrite = false
		}
	}

	return isOverwrite
}

// getGenerateFilePath 获取生成文件路径
func getGenerateFilePath(path, fileName, types string) string {
	return path + SlashSymbol + fileName + UnderlineSymbol + types + GoFileFormat
}

// getTemplatePath 获取模板路径
func getTemplatePath(fileName string) string {
	return TmplPath + fileName
}

// stringToSplit 字符串转数组
func stringToSplit(str string) []string {
	var split []string
	if strings.Contains(str, ",") {
		split = strings.Split(str, ",")
	} else {
		split = append(split, str)
	}
	return split
}

func getDirPath(dir string, path string) string {
	if path == "" {
		return dir
	}
	return dir + SlashSymbol + path
}
