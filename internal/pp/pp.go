package pp

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/mattn/go-jsonpointer"
	"github.com/uga-rosa/umx/internal/fs"
	"github.com/urfave/cli/v2"
)

func Cmd(c *cli.Context) error {
	argInput := c.String("input")
	argOutput := c.String("output")
	argPego := c.String("pego")
	argRg := c.String("rg")

	var jsonObj interface{}
	json.Unmarshal([]byte("{}"), &jsonObj)

	aus, boxz, err := ReadAUSZ(argInput)
	if err != nil {
		return err
	}
	jsonpointer.Set(jsonObj, "/aus", aus)
	jsonpointer.Set(jsonObj, "/boxz", boxz)

	pegoFiles, err := fs.GetFiles(argPego)
	if err != nil {
		return err
	}

	pegos := make(map[string][][]float64)
	for _, path := range pegoFiles {
		pego, err := ReadPEGO(path)
		if err != nil {
			return err
		}
		name := fs.GetFileNameWithoutExt(path)
		pegos[name] = pego
	}
	jsonpointer.Set(jsonObj, "/pego", pegos)

	rgFiles, err := fs.GetFiles(argRg)
	if err != nil {
		return err
	}

	rgs := make(map[string][]float64)
	for _, path := range rgFiles {
		r, err := ReadRg(path)
		if err != nil {
			return err
		}
		name := fs.GetFileNameWithoutExt(path)
		rgs[name] = r
	}
	jsonpointer.Set(jsonObj, "/rg", rgs)

	err = fs.WriteJson(argOutput, jsonObj)

	return err
}

func ReadAUSZ(filename string) ([]float64, float64, error) {
	ausSet := make(Set)
	aus := make([]float64, 0)

	lines, err := fs.ReadLines(filename)
	if err != nil {
		return nil, 0, err
	}

	for i := 2; i < len(lines)-1; i++ {
		z := strings.TrimSpace(lines[i][36:44])
		if strings.TrimSpace(lines[i][10:15]) == "AUS" && !ausSet.Contains(z) {
			f, err := strconv.ParseFloat(z, 64)
			if err != nil {
				return nil, 0, err
			}
			ausSet.Add(z)
			aus = append(aus, f)
		}
	}

	last := strings.TrimSpace(lines[len(lines)-1])
	lastSplit := regexp.MustCompile(`\s+`).Split(last, -1)
	boxZ, err := strconv.ParseFloat(lastSplit[2], 64)
	if err != nil {
		return nil, 0, err
	}

	return aus, boxZ, nil
}

func ReadPEGO(file string) ([][]float64, error) {
	lines, err := fs.ReadLines(file)
	if err != nil {
		return nil, err
	}

	datas := make([][]float64, 0, len(lines))

	for _, line := range lines {
		if !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "@") {
			line = strings.TrimSpace(line)
			lineSplit := regexp.MustCompile(`\s+`).Split(line, -1)
			lineSplitF := make([]float64, 0, len(lineSplit))
			for i := 1; i < len(lineSplit); i++ {
				f, err := strconv.ParseFloat(lineSplit[i], 64)
				if err != nil {
					return nil, err
				}
				lineSplitF = append(lineSplitF, f)
			}
			datas = append(datas, lineSplitF)
		}
	}
	return datas, nil
}

func ReadRg(filename string) ([]float64, error) {
	lines, err := fs.ReadLines(filename)
	if err != nil {
		return nil, err
	}

	rgs := make([]float64, 0, len(lines))
	for _, line := range lines {
		if !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "@") {
			line = strings.TrimSpace(line)
			lineSplit := regexp.MustCompile(`\s+`).Split(line, -1)
			rg, err := strconv.ParseFloat(lineSplit[1], 64)
			if err != nil {
				return nil, err
			}
			rgs = append(rgs, rg)
		}
	}
	return rgs, nil
}

type Set map[string]struct{}

func (s *Set) Add(str string) {
	(*s)[str] = struct{}{}
}

func (s *Set) Contains(str string) bool {
	_, ok := (*s)[str]
	return ok
}
