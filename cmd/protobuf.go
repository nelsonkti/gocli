package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gocli/util/xprintf"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {

	Cmd.AddCommand(protobufCommand)

	protobufCommand.Flags().StringVarP(&path, "path", "p", "proto", PathTips)
}

var protobufCommand = &cobra.Command{
	Use:   "make:pd",
	Short: "Generate Protobuf files",
	Long:  `This command generates Go files from Protobuf definitions.`,
	Run: func(cmd *cobra.Command, args []string) {

		scanProtobuf(path)

		fmt.Println(xprintf.Blue("Generate Protobuf files successfully!"))
	},
}

func scanProtobuf(path string) {
	protoDir := path
	var protoFiles []string

	// Recursively find all .proto files
	err := filepath.Walk(protoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".proto") {
			protoFiles = append(protoFiles, path)
		}
		return nil
	})

	if err != nil {
		xprintf.Red(fmt.Sprintf("Failed to scan proto directory: %v\n", err))
		return
	}

	for _, protoFile := range protoFiles {
		if err := generateProtobuf(protoFile); err != nil {
			xprintf.Red(fmt.Sprintf("Failed to generate Protobuf for %s: %v\n", protoFile, err))
		}
	}
}

func generateProtobuf(protoFile string) error {
	// Set the paths for protoc-gen-go and protoc-gen-go-grpc plugins
	genGoPath := os.Getenv("GOPATH") + "/bin/protoc-gen-go"
	genGoGrpcPath := os.Getenv("GOPATH") + "/bin/protoc-gen-go-grpc"

	// Prepare the command
	outDir := filepath.Join(filepath.Dir(protoFile), "pb")
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	cmd := exec.Command("protoc",
		fmt.Sprintf("--plugin=protoc-gen-go=%s", genGoPath),
		fmt.Sprintf("--plugin=protoc-gen-go-grpc=%s", genGoGrpcPath),
		fmt.Sprintf("--go_out=%s", outDir),
		fmt.Sprintf("--go-grpc_out=%s", outDir),
		protoFile,
	)

	// Capture the output
	var out, stderr strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Execute the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%v: %v", err, stderr.String())
	}

	fmt.Println(out.String())
	return nil
}
