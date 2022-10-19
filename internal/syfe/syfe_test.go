package syfe

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapSymbols(t *testing.T) {
	type args struct {
		allocation    Allocation
		symbolMapping SymbolMapping
	}
	tests := []struct {
		name     string
		args     args
		want     Allocation
		checkErr assert.ErrorAssertionFunc
	}{
		{
			name: "map single correctly",
			args: args{
				allocation: map[Symbol]float64{
					"ABC": 100,
				},
				symbolMapping: map[Symbol]Symbol{
					"ABC": "QQQ",
				},
			},
			want: map[Symbol]float64{
				"QQQ": 100,
			},
			checkErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "map duplication correctly",
			args: args{
				allocation: map[Symbol]float64{
					"ABC1": 100,
					"ABC2": 100,
				},
				symbolMapping: map[Symbol]Symbol{
					"ABC1": "QQQ",
					"ABC2": "QQQ",
				},
			},
			want: map[Symbol]float64{
				"QQQ": 200,
			},
			checkErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "symbol missing from map",
			args: args{
				allocation: map[Symbol]float64{
					"ABC1": 100,
					"ABC2": 100,
				},
				symbolMapping: map[Symbol]Symbol{
					"ABC1": "QQQ",
				},
			},
			checkErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == ErrSymbolMissingFromMapping
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapSymbols(tt.args.allocation, tt.args.symbolMapping)
			tt.checkErr(t, err, fmt.Sprintf("Validate()"))

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapSymbols() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTruncateAfterNthLargest(t *testing.T) {
	type args struct {
		allocation Allocation
		n          int
	}
	tests := []struct {
		name string
		args args
		want Allocation
	}{
		{
			name: "take exactly N",
			args: args{
				allocation: Allocation{
					"QQQ": 100,
				},
				n: 1,
			},
			want: Allocation{
				"QQQ": 100,
			},
		},
		{
			name: "take more",
			args: args{
				allocation: Allocation{
					"QQQ": 100,
				},
				n: 100,
			},
			want: Allocation{
				"QQQ": 100,
			},
		},
		{
			name: "take less",
			args: args{
				allocation: Allocation{
					"A": 1,
					"B": 2,
				},
				n: 1,
			},
			want: Allocation{
				"B": 2,
			},
		},
		{
			name: "take less with equal",
			args: args{
				allocation: Allocation{
					"A": 1,
					"B": 1,
				},
				n: 1,
			},
			want: Allocation{
				"A": 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, TruncateAfterNthLargest(tt.args.allocation, tt.args.n), "TruncateAfterNthLargest(%v, %v)", tt.args.allocation, tt.args.n)
		})
	}
}
