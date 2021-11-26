package rg

import (
	"fmt"
	"math"
	"sort"

	"github.com/Arafatk/glot"
	"github.com/uga-rosa/umx/internal/adsp"
	"github.com/uga-rosa/umx/internal/fs"
	"github.com/urfave/cli/v2"
)

const Bins = 100

func Cmd(c *cli.Context) error {
	input := c.String("input")
	number := c.Int("number")

	ads, non, err := DivideByAds(input, number)
	if err != nil {
		return err
	}

	err = fs.WriteJson(fmt.Sprintf("umx/rg%d.json", number), &fs.JsonRg{"ads": ads, "non": non})
    if err != nil {
        return err
    }

	adsHist, adsEdge, adsMean := Histogram(ads, Bins)
	nonHist, nonEdge, nonMean := Histogram(non, Bins)

	xr := []float64{
		math.Min(adsEdge[0], nonEdge[0]),
		math.Max(adsEdge[Bins], nonEdge[Bins]),
	}

	DrawHist(adsMean, adsHist, xr, NewOpt("ads", number))
	DrawHist(nonMean, nonHist, xr, NewOpt("non", number))

	return nil
}

func DivideByAds(inputPP string, number int) ([]float64, []float64, error) {
    jsonObj := &fs.JsonPP{}
	err := fs.ReadJson(inputPP, jsonObj)
	if err != nil {
		return nil, nil, err
	}

    rgs := jsonObj.Rg

	adsFileName := adsp.AdsFileName(number)
    adsOnSteps := make(fs.JsonAds)

	err = fs.ReadJson(adsFileName, &adsOnSteps)
	if err != nil {
		return nil, nil, err
	}

	ads := make([]float64, 0)
	non := make([]float64, 0)

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
func Histogram(data []float64, bin int) ([]float64, []float64, []float64) {
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

type Option struct {
	filename string
	title    string
	kind     string
}

func NewOpt(kind string, number int) *Option {
	return &Option{
		fmt.Sprintf("Histogram_Rg_%s_%d.png", kind, number),
		fmt.Sprintf("{/Arial=18 Histogram of Radius of Gyration of %s}", kind),
		kind,
	}
}

func DrawHist(x []float64, y []float64, xr []float64, opt *Option) {
	plot, _ := glot.NewPlot(2, false, false)
	plot.AddPointGroup(opt.kind, "lines", [][]float64{x, y})
	plot.SetTitle(opt.title)
	plot.Cmd(fmt.Sprintf("set xrange [%.2f:%.2f]", xr[0], xr[1]))
	plot.SetXLabel("{/Arial=18 Radius of Gyration / nm}")
	plot.SetYLabel("{/Arial=18 Number of degrees}")
	plot.SavePlot(opt.filename)
}
