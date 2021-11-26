package hist

import (
	"fmt"
	"math"

	"github.com/uga-rosa/umx/internal/fs"
	"github.com/uga-rosa/umx/internal/rg"
	"github.com/urfave/cli/v2"
)

func Cmd(c *cli.Context) error {
    input := c.String("input")
    number := c.Int("number")

    jsonRg := fs.JsonRg{}

    err := fs.ReadJson(input, &jsonRg)
    if err != nil {
        return err
    }

    ads, ok := jsonRg["ads"]
    if !ok {
        return fmt.Errorf(".ads is null")
    }
    non, ok := jsonRg["non"]
    if !ok {
        return fmt.Errorf(".non is null")
    }

    adsHist, adsEdge, adsMean := rg.Histogram(ads, rg.Bins)
    nonHist, nonEdge, nonMean := rg.Histogram(non, rg.Bins)

    xr := []float64{
        math.Min(adsEdge[0], nonEdge[rg.Bins]),
        math.Max(adsEdge[rg.Bins], nonEdge[rg.Bins]),
    }

	rg.DrawHist(adsMean, adsHist, xr, rg.NewOpt("ads", number))
	rg.DrawHist(nonMean, nonHist, xr, rg.NewOpt("non", number))

    return nil
}
