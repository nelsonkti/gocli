package format

import (
	"gopkg.in/yaml.v3"
)

type Yaml struct {
}

func (y *Yaml) Load(content []byte, config *map[string]interface{}) error {
	// 解析YAML
	err := yaml.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	return nil
}
