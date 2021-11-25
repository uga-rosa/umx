package rg

import (
	"fmt"
	"math"
	"sort"

	"github.com/Arafatk/glot"
	"github.com/mattn/go-jsonpointer"
	"github.com/uga-rosa/umx/internal/adsp"
	"github.com/uga-rosa/umx/internal/fs"
	"github.com/urfave/cli/v2"
)

const bins = 100

func Cmd(c *cli.Context) error {
	input := c.String("input")
	number := c.Int("number")

	ads, non, err := DivideByAds(input, number)
	if err != nil {
		return err
	}

	err = fs.WriteJson(fmt.Sprintf("umx/rg%d.json", number), map[string][]float64{"ads": ads, "non": non})
    if err != nil {
        return err
    }

	adsHist, adsEdge, adsMean := histogram(ads, bins)
	nonHist, nonEdge, nonMean := histogram(non, bins)

	xr := []float64{
		math.Min(adsEdge[0], nonEdge[0]),
		math.Max(adsEdge[bins], nonEdge[bins]),
	}

	drawHist(adsMean, adsHist, xr, newOpt("ads", number))
	drawHist(nonMean, nonHist, xr, newOpt("non", number))

	return nil
}

func DivideByAds(input string, number int) ([]float64, []float64, error) {
	jsonObj, err := fs.ReadJson(input)
	if err != nil {
		return nil, nil, err
	}

	rgs, err := parseJson(jsonObj)
	if err != nil {
		return nil, nil, err
	}

	adsOnSteps, err := getAdsOnSteps(input, number)
	if err != nil {
		return nil, nil, err
	}

	ads := make([]float64, len(rgs))
	non := make([]float64, len(rgs))

	for name, rg := range rgs {
		adsOnStep, ok := adsOnSteps[name]
		if !ok {
			return nil, nil, fmt.Errorf("%q is not included in %q", name, adsp.AdsFileName(number))
		}

		a, n := divideByAds(rg, adsOnStep)
		ads = mergeSlice(ads, a)
		non = mergeSlice(non, n)
	}

	return ads, non, nil
}

func parseJson(obj interface{}) (map[string][]float64, error) {
	rgsInterface, err := jsonpointer.Get(obj, "/rg")
	if err != nil {
		return nil, err
	}
	rgs, ok := rgsInterface.(map[string][]float64)
	if !ok {
		return nil, fmt.Errorf("/rg is not map[string][]float64")
	}

	return rgs, nil
}

func getAdsOnSteps(input string, number int) (map[string][]bool, error) {
	adsOnSteps := make(map[string][]bool)
	adsFileName := adsp.AdsFileName(number)

	var ok bool
	if fs.FileExist(adsFileName) {
		adsStepJsonObj, err := fs.ReadJson(adsFileName)
		if err != nil {
			return nil, err
		}
		adsOnSteps, ok = adsStepJsonObj.(map[string][]bool)
		if !ok {
			return nil, fmt.Errorf("%q is not map[string][]bool", adsFileName)
		}
	} else {
		adsOnSteps, _, err := adsp.CalcAds(input, number)
		if err != nil {
			return nil, err
		}

		err = adsp.WriteAdsJson(adsOnSteps, number)
		if err != nil {
			return nil, err
		}
	}
	return adsOnSteps, nil
}

func divideByAds(rg []float64, adsOnStep []bool) ([]float64, []float64) {
	a := make([]float64, 0, len(rg))
	n := make([]float64, 0, len(rg))

	for i, r := range rg {
		if adsOnStep[i] {
			a = append(a, r)
		} else {
			n = append(n, r)
		}
	}

	return a, n
}

func mergeSlice(dst, src []float64) []float64 {
	for _, s := range src {
		dst = append(dst, s)
	}
	return dst
}

// return hist, edge, mean
func histogram(data []float64, bin int) ([]float64, []float64, []float64) {
	sort.Float64s(data)
	min := data[0]
	max := data[len(data)-1]
	width := (max - min) / float64(bin)

	hist := make([]float64, bin)
	edge := make([]float64, bin+1)
	mean := make([]float64, bin)

	edge[0] = min
	for i := 1; i < bin; i++ {
		edge[i] = min + width*float64(i)
	}
	edge[bin] = max

	for i := 0; i < bin; i++ {
		mean[i] = (edge[i] + edge[i+1]) / 2
	}

	countEdge := 0
	countData := 0
	for countData < len(data) {
		if data[countData] <= edge[countEdge+1] {
			hist[countEdge] += 1
			countData += 1
		} else {
			countEdge += 1
		}
	}

	return hist, edge, mean
}

type option struct {
	filename string
	title    string
	kind     string
}

func newOpt(kind string, number int) *option {
	return &option{
		fmt.Sprintf("Histogram_RG_PEG_Ads_%d.png", number),
		fmt.Sprintf("Histogram of Radius of Gyration of %s", kind),
		kind,
	}
}

func drawHist(x []float64, y []float64, xr []float64, opt *option) {
	plot, _ := glot.NewPlot(2, false, false)
	plot.AddPointGroup(opt.kind, "lines", [][]float64{x, y})
	plot.SetTitle(opt.title)
	plot.Cmd(fmt.Sprintf("set xrange [%.2f:%.2f]", xr[0], xr[1]))
	plot.SetXLabel("{/Arial=30 Radius of Gyration / nm}")
	plot.SetYLabel("{/Arial=30 Number of degrees}")
	plot.SavePlot(opt.filename)
}
