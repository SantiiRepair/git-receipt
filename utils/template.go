package utils

import (
	"fmt"
	"os"
)

func LoadTemplate(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error loading template: %v", err)
	}
	return string(content), nil
}
