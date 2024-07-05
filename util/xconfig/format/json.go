package format

import (
	"encoding/json"
)

type Json struct {
}

func (j *Json) Load(content []byte, config *map[string]interface{}) error {
	err := json.Unmarshal(content, &config)
	if err != nil {
		return err
	}
	return nil
}
