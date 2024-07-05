package file

import (
	"errors"
	"mashang/util/xconfig/format"
	"os"
	"path/filepath"
)

type Config struct {
	Path        string
	content     []byte
	fileExtName string
	fileLoader  *format.Format
}

func NewConfig(path string) *Config {
	return &Config{
		Path:       path,
		fileLoader: format.NewFileFormat(),
	}
}

func (f *Config) Load() (map[string]interface{}, error) {
	// 读取本地文件逻辑
	f.fileExt(f.Path)
	if f.fileLoader.FileFormat[f.fileExtName] == nil {
		return nil, errors.New("不支持该文件类型")
	}

	err := f.readFile()
	if err != nil {
		return nil, err
	}
	config := make(map[string]interface{})

	err = f.fileLoader.FileFormat[f.fileExtName].Load(f.content, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (f *Config) fileExt(filePath string) {
	ext := filepath.Ext(filePath)
	if ext == "" {
		f.fileExtName = ext
	}
	f.fileExtName = ext[1:]
}

func (f *Config) readFile() error {
	var err error
	f.content, err = os.ReadFile(f.Path)
	return err
}
