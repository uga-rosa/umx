package fs

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ReadLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }

	lines := make([]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func GetFiles(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
    if err != nil {
        return nil, err
    }

	paths := make([]string, 0)

	for _, file := range files {
		if !file.IsDir() {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}

	return paths, nil
}

func GetFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func ReadJson(file string) (interface{}, error) {
    bytes, err := ioutil.ReadFile(file)
    if err != nil {
        return nil, err
    }

    var jsonObj interface{}
    err = json.Unmarshal(bytes, jsonObj)
    if err != nil {
        return nil, err
    }

    return jsonObj, nil
}

func WriteJson(file string, obj interface{}) error {
    bytes, err := json.Marshal(obj)
    if err != nil {
        return err
    }

    err = ioutil.WriteFile(file, bytes, 0644)

    return err
}

func FileExist(file string) bool {
    _, err := os.Stat(file)
    return err == nil
}
