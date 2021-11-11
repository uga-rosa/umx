package fs

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/uga-rosa/umx/internal/set"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ReadLines(filename string) []string {
	file, err := os.Open(filename)
	check(err)
	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func GetFiles(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	check(err)
	paths := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}
	return paths
}

func ReadAUSZ(filename string) (set.SetF, float64) {
	ausF := make(set.SetF)
	ausS := make(set.SetS)
	lines := ReadLines(filename)
	for i := 2; i < len(lines)-1; i++ {
		if strings.TrimSpace(lines[i][10:15]) == "AUS" {
			ausS.Add(strings.TrimSpace(lines[i][36:44]))
		}
	}
	for s := range ausS {
		f, err := strconv.ParseFloat(s, 64)
		check(err)
		ausF.Add(f)
	}
	last := strings.TrimSpace(lines[len(lines)-1])
	lastSplit := regexp.MustCompile(`\s+`).Split(last, -1)
	boxZ, err := strconv.ParseFloat(lastSplit[2], 64)
	check(err)
	return ausF, boxZ
}

func ReadPEGO(file string) [][]float64 {
	lines := ReadLines(file)
	datas := make([][]float64, 0, len(lines))
	for _, line := range lines {
		if !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "@") {
			line = strings.TrimSpace(line)
			lineSplit := regexp.MustCompile(`\s+`).Split(line, -1)
			lineSplitF := make([]float64, 0, len(lineSplit))
			for i := 1; i < len(lineSplit); i++ {
				f, err := strconv.ParseFloat(lineSplit[i], 64)
				check(err)
				lineSplitF = append(lineSplitF, f)
			}
			datas = append(datas, lineSplitF)
		}
	}
	return datas
}
