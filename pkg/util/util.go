package util

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// ExpandPath expands ~ to home directory.
func ExpandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func ToJsonString(v interface{}) string {
	data, err := json.Marshal(v)
	if nil != err {
		return ""
	}
	return string(data)
}

func ToJsonIndent(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "\t")
	if nil != err {
		return ""
	}
	return string(data)
}

func Exists(dir string) (bool, error) {
	if dir == "" {
		return false, nil
	}
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
