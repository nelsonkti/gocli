package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const Version = "1.3"

func init() {
	Cmd.AddCommand(versionCommand)
}

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "显示版本号",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
