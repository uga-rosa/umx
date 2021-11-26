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

func FileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

// Json

type JsonStruct interface {
	IsJsonStruct()
}

type JsonPP struct {
	Aus  []float64              `json:"aus"`
	Boxz float64                `json:"boxz"`
	Pego map[string][][]float64 `json:"pego"`
	Rg   JsonRg                 `json:"rg"`
}

func (jp *JsonPP) IsJsonStruct() {}

type JsonRg map[string][]float64

func (jr JsonRg) IsJsonStruct() {}

type JsonAds map[string][]bool

func (ja JsonAds) IsJsonStruct() {}

func ReadJson(file string, jsonStruct JsonStruct) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, jsonStruct)
	if err != nil {
		return err
	}

	return nil
}

func WriteJson(file string, jsonStruct JsonStruct) error {
	bytes, err := json.Marshal(jsonStruct)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file, bytes, 0644)

	return err
}
