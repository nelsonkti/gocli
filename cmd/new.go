package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/nelsonkti/gocli/util/project"
	"github.com/nelsonkti/gocli/util/xprintf"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var branch string

func init() {

	Cmd.AddCommand(newProject)

	newProject.Flags().StringVarP(&path, "path", "p", "", PathTips)
	newProject.Flags().StringVarP(&branch, "branch", "b", "master", "请输入 [分支] 名称")
}

const layout = "https://github.com/nelsonkti/gin-framework.git"

var newProject = &cobra.Command{
	Use:   "new",
	Short: "create a new project",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			xprintf.Red("new: unexpected arguments")
		}
		projectName := args[0]
		var dir = "./"
		if len(args) > 2 {
			dir = args[1]
		}

		if len(args) > 3 && path == "" {
			dir = args[1]
		}
		if len(args) > 3 && path == "" {
			dir = args[1]
		}
		if path != dir {
			dir = path
		}
		to := filepath.Join(dir, projectName)
		if _, err := os.Stat(to); !os.IsNotExist(err) {
			fmt.Printf("🚫 %s already exists\n", projectName)
			prompt := &survey.Confirm{
				Message: "📂 Do you want to override the folder ?",
				Help:    "Delete the existing folder and create the project.",
			}
			var override bool
			e := survey.AskOne(prompt, &override)
			if e != nil {
				fmt.Println(e)
				return
			}
			if !override {
				fmt.Println("Operation cancelled.")
				return
			}
			err := os.RemoveAll(to)
			if err != nil {
				fmt.Println("Failed to delete existing folder:", err)
				return
			}
		}
		fmt.Printf("🚀 Creating Project %s, layout repo is %s, please wait a moment.\n\n", projectName, layout)
		pro := project.NewRepo(layout, branch)

		err := pro.CopyTo(cmd.Context(), to, projectName, []string{})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(args)
	},
}
