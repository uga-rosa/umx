package adsp

import (
	"fmt"
	"math"

	"github.com/Arafatk/glot"
	"github.com/mattn/natural"
	"github.com/uga-rosa/umx/internal/fs"
	"github.com/urfave/cli/v2"
)

const AdsRegion = 0.7

func Cmd(c *cli.Context) error {
	input := c.String("input")
	number := c.Int("number")

	adsOnSteps, numOfStep, err := CalcAds(input, number)
	if err != nil {
		return err
	}

	err = WriteAdsJson(adsOnSteps, number)

	x, y := numberOfAds(numOfStep, adsOnSteps)

	xr := []int{0, numOfStep}
	yr := []int{0, len(adsOnSteps)}
	opt := &option{
		fmt.Sprintf("Number_of_adsorption_%d.png", number),
		fmt.Sprintf("{/Arial=14 Time transition of adsorption number (adsorption condition %d)}", number),
	}

	drawLine(x, y, xr, yr, opt)

	return nil
}

func CalcAds(input string, number int) (fs.JsonAds, int, error) {
    jsonObj := &fs.JsonPP{}
	err := fs.ReadJson(input, jsonObj)
	if err != nil {
		return nil, 0, err
	}

    aus := jsonObj.Aus
    boxz := jsonObj.Boxz
    pegos := jsonObj.Pego

	pros := make([]float64, 0, len(pegos))
	adsOnSteps := make(fs.JsonAds)
	numOfStep := 0

	for _, name := range sortedKeys(pegos) {
        pego := pegos[name]

		if numOfStep == 0 {
			numOfStep = len(pego)
		}

		adsOnStep := calcAds(aus, pego, boxz, number)
		adsOnSteps[name] = adsOnStep

		pro := float64(trueCount(adsOnStep)) / float64(len(adsOnStep)) * 100
		fmt.Println(name, pro, "%")

		pros = append(pros, pro)
	}
	fmt.Println("mean: ", mean(pros))

	return adsOnSteps, numOfStep, nil
}

func WriteAdsJson(adsOnSteps fs.JsonAds, number int) error {
	adsFileName := AdsFileName(number)
    err := fs.WriteJson(adsFileName, adsOnSteps)
	if err != nil {
		return err
	}
	return nil
}

func AdsFileName(number int) string {
	return fmt.Sprintf("umx/adsStep%d.json", number)
}

func calcAds(aus []float64, pego [][]float64, boxZ float64, number int) []bool {
	step := make([]bool, 0, len(pego))
	for _, pegoOnStep := range pego {
		isAds := false
		counter := 0 // 吸着条件を満たしたO原子の数
		for _, o := range pegoOnStep {
			var lengthPegoAus []float64 // O原子との距離
			for _, a := range aus {
				l := math.Abs(o - a)
				if l > boxZ/2 { // PBC
					l = boxZ - l
				}
				lengthPegoAus = append(lengthPegoAus, l)
			}
			if min(lengthPegoAus) <= AdsRegion {
				counter += 1
				if counter >= number {
					isAds = true
					break
				}
			}
		}
		step = append(step, isAds)
	}
	return step
}

func sortedKeys(m map[string][][]float64) []string {
    s := make([]string, len(m))
    index := 0
    for key := range m {
        s[index] = key
        index++
    }
    natural.Sort(s)
    return s
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

func mean(s []float64) float64 {
	var sum float64
	for _, v := range s {
		sum += v
	}
	return sum / float64(len(s))
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

func numberOfAds(numSteps int, datas map[string][]bool) ([]int, []int) {
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

type option struct {
	filename string
	title    string
}

func drawLine(x, y []int, xr, yr []int, opt *option) {
	plot, _ := glot.NewPlot(2, false, false)
	plot.AddPointGroup("Oxygen atoms in PEG", "lines", [][]int{x, y})
	plot.SetTitle(opt.title)
	plot.SetXLabel("{/Arial=18 Steps}")
	plot.SetYLabel("{/Arial=18 Number of adsorptions}")
	plot.SetXrange(xr[0], xr[1])
	plot.SetYrange(yr[0], yr[1])
	plot.SavePlot(opt.filename)
}
