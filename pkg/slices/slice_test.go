package slices_test

import (
	"reflect"
	"testing"

	"notify/pkg/slices"
)

func TestReverse(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int
	}{
		{
			name: "Test 1",
			in:   []int{1, 2, 3, 4, 5},
			want: []int{5, 4, 3, 2, 1},
		},
		{
			name: "Test 2",
			in:   []int{10, 20, 30, 40},
			want: []int{40, 30, 20, 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := slices.Reverse(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}
