package rg

import "testing"

func TestHistogram(t *testing.T) {
	tests := []struct {
		input        []float64
		bin          int
		expectedHist []float64
		expectedEdge []float64
		expectedMean []float64
	}{
		{
			[]float64{0.0, 1.0, 2.0, 2.5, 3.0, 3.0, 3.0},
			3,
			[]float64{2.0, 1.0, 4.0},
			[]float64{0.0, 1.0, 2.0, 3.0},
			[]float64{0.5, 1.5, 2.5},
		},
	}

	for _, tt := range tests {
		actualHist, actualEdge, actualMean := histogram(tt.input, tt.bin)
		testEqualFloatSlice(t, tt.expectedHist, actualHist)
		testEqualFloatSlice(t, tt.expectedEdge, actualEdge)
		testEqualFloatSlice(t, tt.expectedMean, actualMean)
	}
}

func testEqualFloatSlice(t *testing.T, s1, s2 []float64) {
	if len(s1) != len(s2) {
		t.Errorf("Slices of different lenghts (%d and %d) are compared", len(s1), len(s2))
	}

	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			t.Errorf("The element of the slices (%f and %f) are different.", s1[i], s2[i])
		}
	}
}
