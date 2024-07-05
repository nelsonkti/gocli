package format

const (
	YAML = "yaml"
	JSON = "json"
)

type Format struct {
	FileFormat map[string]fileType
}

type fileType interface {
	Load(content []byte, config *map[string]interface{}) error
}

func NewFileFormat() *Format {
	fileTypeMap := make(map[string]fileType)
	fileTypeMap[YAML] = &Yaml{}
	fileTypeMap[JSON] = &Json{}
	return &Format{
		FileFormat: fileTypeMap,
	}
}
