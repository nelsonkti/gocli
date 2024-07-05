package xfile

import (
	"fmt"
	"gocli/util/xstring"
	"os"
	"path/filepath"
	"strings"
)

var defaultPath = "./"

func GetPath(file string) string {
	i := strings.LastIndex(file, string(os.PathSeparator))

	if i > 0 {
		file := file[:i+1]

		if xstring.IsUpper(file) {
			file = xstring.Camel2Case(file)
		}
		if file[0:1] != string(os.PathSeparator) && file[0:1] != "." {
			file = fmt.Sprintf(".%s%s", string(os.PathSeparator), file)
		}
		if strings.Contains(file, "//") {
			file = strings.Replace(file, "//", "/", -1)
		}
		return file
	}

	return defaultPath
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetPackageName(path string) string {
	return filepath.Base(path)
}

func MkdirAll(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
}
