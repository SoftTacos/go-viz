package util

import (
	"reflect"
	"testing"
)

func TestGroupFrequencies(t *testing.T) {
	tests := []struct {
		name       string
		n int
		f []float64
		wantGroups []float64
	}{
		{
			name:"ez",
			n:3,
			f:[]float64{
				1,1,1,
				2,2,2,
				1,5,9,
			},
			wantGroups: []float64{1,2,5},
		},
		{
			name:"1 extra, round down",
			n:3,
			f:[]float64{
				1,1,1,
				2,2,2,2,
				3,3,3,
			},
			wantGroups: []float64{1,2,3},
		},
		{
			name:"2 extra, round up",
			n:3,
			f:[]float64{
				1,1,1,1,
				2,2,2,
				3,3,3,3,
			},
			wantGroups: []float64{1,2,3},
		},
		{
			name:"nothing",
			n:3,
			f:[]float64{
			},
			wantGroups: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotGroups := GroupFrequencies(tt.n, tt.f); !reflect.DeepEqual(gotGroups, tt.wantGroups) {
				t.Errorf("GroupFrequencies() = %v, want %v", gotGroups, tt.wantGroups)
			}
		})
	}
}
