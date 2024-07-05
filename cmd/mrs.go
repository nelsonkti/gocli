package cmd

import (
	"github.com/spf13/cobra"
)

func init() {

	Cmd.AddCommand(mrsCommand)

	mrsCommand.Flags().StringVarP(&name, "name", "n", "", NameTips)
	mrsCommand.Flags().StringVarP(&database, "database", "d", "", DatabaseTips)
	mrsCommand.Flags().StringVarP(&fileName, "file_name", "f", "", TableNameTips)
	mrsCommand.Flags().StringVarP(&path, "path", "p", "", PathTips)
}

var mrsCommand = &cobra.Command{
	Use:   "make:mrs",
	Short: "create a new model file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		zhName = ModelZhName
		createModel()

		zhName = RepositoryZhName
		createRepository()

		zhName = ServiceZhName
		createService()
	},
}
