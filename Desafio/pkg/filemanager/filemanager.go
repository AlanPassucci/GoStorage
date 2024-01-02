package filemanager

import (
	"encoding/json"
	"os"
)

func LoadDataFromJSON(path string, destiny any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, &destiny); err != nil {
		return err
	}

	return nil
}
