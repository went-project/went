package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func CreateConfigFile(projectName string, req interface{}) error {
	configPath := filepath.Join(projectName, "wentconfig.json")

	jsonData, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, jsonData, 0755)

}

func ReadConfigFile() (map[string]interface{}, error) {
	configPath := "wentconfig.json"

	jsonData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var configData map[string]interface{}
	err = json.Unmarshal(jsonData, &configData)
	if err != nil {
		return nil, err
	}

	return configData, nil
}

func CreateFileWithContent(filePath string, content string) error {
	err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, []byte(content), 0644)
}
