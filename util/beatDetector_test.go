package util

import (
	"testing"
)

func TestCalcAbsoluteMeanAndVariance(t *testing.T) {
	tests := []struct {
		name     string
		set      []float64
		wantMean float64
		wantVar  float64
	}{
		{
			name:     "idk",
			set:      []float64{1, 2, 3},
			wantMean: 2,
			wantVar:  (1 + 1) / 3.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMean, gotVar := CalcMeanVariance(tt.set)
			if gotMean != tt.wantMean {
				t.Errorf("CalcMeanVariance() gotMean = %v, want %v", gotMean, tt.wantMean)
			}
			if gotVar != tt.wantVar {
				t.Errorf("CalcMeanVariance() gotVar = %v, want %v", gotVar, tt.wantVar)
			}
		})
	}
}

func TestCalcMultiSetStats(t *testing.T) {
	type args struct {
		n         float64
		means     []float64
		variances []float64
	}
	tests := []struct {
		name         string
		args         args
		wantMean     float64
		wantVariance float64
	}{
		{
			name: "1",
			args: args{
				n:5,
				means:[]float64{.1,.02,3,4,5},
				variances:[]float64{1,2,3,4,5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMean, gotVariance := CalcMultiSetStats(tt.args.n, tt.args.means, tt.args.variances)
			if gotMean != tt.wantMean {
				t.Errorf("CalcMultiSetStats() gotMean = %v, want %v", gotMean, tt.wantMean)
			}
			if gotVariance != tt.wantVariance {
				t.Errorf("CalcMultiSetStats() gotVariance = %v, want %v", gotVariance, tt.wantVariance)
			}
		})
	}
}
