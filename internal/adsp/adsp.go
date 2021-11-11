package adsp

import (
	"fmt"
	"math"
	"path/filepath"

	"github.com/Arafatk/glot"
	"github.com/uga-rosa/umx/internal/fs"
	"github.com/uga-rosa/umx/internal/set"
	"github.com/urfave/cli/v2"
)

const AdsRegion = 0.7

func Cmd(c *cli.Context) error {
	input := c.String("input")
	directory := c.String("directory")
	number := c.Int("number")

	aus, boxz := fs.ReadAUSZ(input)
	pegOFiles := fs.GetFiles(directory)
	adsOnSteps := make(map[string][]bool)
	var numSteps int
	for i, file := range pegOFiles {
		pegO := fs.ReadPEGO(file)
		if i == 0 {
			numSteps = len(pegO)
		} else if numSteps != len(pegO) {
			panic("Data with different number of steps has been loaded.")
		}
		adsOnStep := CalcAds(aus, pegO, boxz, number)
		pro := float64(trueCount(adsOnStep)) / float64(len(adsOnStep)) * 100
		name := getFileNameWithoutExt(file)
		fmt.Println(name, pro, "%")
		adsOnSteps[name] = adsOnStep
	}
	x, y := numberOfAds(adsOnSteps, numSteps)
	drawLine(
		fmt.Sprintf("Number_of_adsorption_%d.png", number),
		fmt.Sprintf("Time transition of adsorption number (adsorption condition %d)", number),
		x,
		y,
		numSteps, len(adsOnSteps),
	)
	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func CalcAds(ausZ set.SetF, pegO [][]float64, boxZ float64, number int) []bool {
	step := make([]bool, 0, len(pegO))
	for _, pegOOnStep := range pegO {
		adsFlag := false
		counter := 0 // 吸着条件を満たしたO原子の数
		for _, o := range pegOOnStep {
			var lengthPerA []float64 // O原子との距離
			for a := range ausZ {
				l := math.Abs(o - a)
				// PBC
				if l > boxZ/2 {
					l = boxZ - l
				}
				lengthPerA = append(lengthPerA, l)
			}
			if min(lengthPerA) <= AdsRegion {
				counter += 1
				if counter >= number {
					adsFlag = true
					break
				}
			}
		}
		step = append(step, adsFlag)
	}
	return step
}

func min(s []float64) float64 {
	m := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] < m {
			m = s[i]
		}
	}
	return m
}

func trueCount(s []bool) int {
	c := 0
	for _, b := range s {
		if b {
			c += 1
		}
	}
	return c
}

func numberOfAds(datas map[string][]bool, numSteps int) ([]int, []int) {
	x := make([]int, numSteps)
	for i := 0; i < numSteps; i++ {
		x[i] = i
	}

	y := make([]int, numSteps)
	for _, ads := range datas {
		for i, b := range ads {
			if b {
				y[i] += 1
			}
		}
	}

	return x, y
}

func drawLine(filename, title string, x []int, y []int, xmax, ymax int) {
	plot, _ := glot.NewPlot(2, false, false)
	plot.AddPointGroup("Oxygen atoms in PEG", "lines", [][]int{x, y})
	plot.SetTitle(title)
	plot.SetXLabel("Steps")
	plot.SetYLabel("Number of adsorptions")
	plot.SetXrange(0, xmax)
	plot.SetYrange(0, ymax)
	plot.SavePlot(filename)
}

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
