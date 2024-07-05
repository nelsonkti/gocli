package cmd

import (
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/spf13/cobra"
	"gocli/util/helper"
	"gocli/util/template"
	"gocli/util/xfile"
	"gocli/util/xprintf"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var (
	database             = ""
	tableName            = ""
	ModelNamespace       = ""
	ModelNamespacePoint  = ""
	ModelStructNameModel = ""
)

const (
	NameTips      = "请输入 [模型] 中文名称"
	TableNameTips = "请输入 [表名]"
	PathTips      = "请输入 [路径]"
	DatabaseTips  = "请输入 [database connection]"
)

const (
	ModelType         = "model"
	modelTemplateFile = "model.tmpl"
)

func init() {

	Cmd.AddCommand(modelCommand)

	modelCommand.Flags().StringVarP(&name, "name", "n", "", NameTips)
	modelCommand.Flags().StringVarP(&database, "database", "d", "", DatabaseTips)
	modelCommand.Flags().StringVarP(&fileName, "file_name", "f", "", TableNameTips)
	modelCommand.Flags().StringVarP(&path, "path", "p", "", PathTips)
}

var modelCommand = &cobra.Command{
	Use:   "make:model",
	Short: "create a new model file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		zhName = ModelZhName
		createModel()
	},
}

func createModel() {
	fmt.Println(xprintf.Blue("creating model file ..."))
	if fileName == "" {
		fmt.Println(xprintf.Red(TableNameTips))
		return
	}

	if database == "" {
		fmt.Println(xprintf.Red(DatabaseTips))
		return
	}

	newPath := getDirPath(ModelDirPath, ModelType, path)

	db, err := gorm.Open(mysql.Open(database))
	if err != nil {
		fmt.Println(xprintf.Red(fmt.Errorf("cannot establish db connection: %w", err).Error()))
		return
	}
	// 生成实例
	g := gen.NewGenerator(gen.Config{
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     false,
		FieldWithIndexTag: false,
		FieldWithTypeTag:  true,
	})

	// 设置目标 db
	g.UseDB(db)

	// 自定义字段的数据类型
	dataMap := map[string]func(columnType gorm.ColumnType) (dataType string){
		"tinyint":   func(columnType gorm.ColumnType) (dataType string) { return "int64" },
		"smallint":  func(columnType gorm.ColumnType) (dataType string) { return "int64" },
		"mediumint": func(columnType gorm.ColumnType) (dataType string) { return "int64" },
		"bigint":    func(columnType gorm.ColumnType) (dataType string) { return "int64" },
		"int":       func(columnType gorm.ColumnType) (dataType string) { return "int64" },
	}
	g.WithDataTypeMap(dataMap)

	tableNames := stringToSplit(fileName)

	isOverwrite := checkFileExists(newPath, tableNames, ModelType)

	for _, tab := range tableNames {
		filePath := getGenerateFilePath(newPath, tab, ModelType)
		exists, _ := xfile.PathExists(filePath)
		if exists && isOverwrite == false {
			continue
		}
		generateModel(g, tab, filePath, newPath)
	}
	fmt.Println(xprintf.Green("create a new 【model】 file success\n"))
}

func generateModel(g *gen.Generator, tab, outputPath, newPath string) {
	// 自定义生成模板
	model := g.GenerateModel(tab)

	xfile.MkdirAll(newPath)

	databaseName, err := helper.GetDatabaseName(database)
	if err != nil {
		databaseName = "default"
	}

	model.QueryStructName = databaseName
	model.S = xfile.GetModPath(RelativeSymbol)

	model.TableComment = model.TableComment + " " + zhName
	if name == "" {
		name = model.TableComment
	}

	model.StructInfo.Package = xfile.GetPackageName(newPath)

	box := packr.New(tmplPath, tmplPath)
	tmpl, _ := box.FindString(modelTemplateFile)
	err = template.WriteFile(outputPath, tmpl, model)
	if err != nil {
		fmt.Println(xprintf.Red(fmt.Sprintf("生成 %s 失败: %+v", outputPath, err)))
		return
	}
	fmt.Println(xprintf.Green("CREATED ") + outputPath)
}
