package metamodel

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapSymbols(t *testing.T) {
	type args struct {
		alloc                       allocation
		symbolMapping               SymbolMapping
		requireMappingForAllSymbols bool
	}
	tests := []struct {
		name     string
		args     args
		want     allocation
		checkErr assert.ErrorAssertionFunc
	}{
		{
			name: "map single correctly",
			args: args{
				alloc: map[string]float64{
					"ABC": 100,
				},
				symbolMapping: map[string]string{
					"ABC": "QQQ",
				},
				requireMappingForAllSymbols: true,
			},
			want: map[string]float64{
				"QQQ": 100,
			},
			checkErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return true
			},
		},
		{
			name: "map duplication correctly",
			args: args{
				alloc: map[string]float64{
					"ABC1": 100,
					"ABC2": 100,
				},
				symbolMapping: map[string]string{
					"ABC1": "QQQ",
					"ABC2": "QQQ",
				},
				requireMappingForAllSymbols: true,
			},
			want: map[string]float64{
				"QQQ": 200,
			},
			checkErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return true
			},
		},
		{
			name: "symbol missing from map",
			args: args{
				alloc: map[string]float64{
					"ABC1": 100,
					"ABC2": 100,
				},
				symbolMapping: map[string]string{
					"ABC1": "QQQ",
				},
				requireMappingForAllSymbols: true,
			},
			checkErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Equal(t, symbolMissingFromMappingError{missingSymbol: "ABC2"}, err)
				return true
			},
		},
		{
			name: "symbol missing from map, but not required",
			args: args{
				alloc: map[string]float64{
					"ABC1": 100,
					"ABC2": 100,
				},
				symbolMapping: map[string]string{
					"ABC1": "QQQ",
				},
				requireMappingForAllSymbols: false,
			},
			want: map[string]float64{
				"QQQ":  100,
				"ABC2": 100,
			},
			checkErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.Nil(t, err)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mapSymbols(tt.args.alloc, tt.args.symbolMapping, tt.args.requireMappingForAllSymbols)
			tt.checkErr(t, err, fmt.Sprintf("Validate()"))

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapSymbols() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTruncateAfterNthLargest(t *testing.T) {
	type args struct {
		alloc allocation
		n     int
	}
	tests := []struct {
		name string
		args args
		want allocation
	}{
		{
			name: "take exactly N",
			args: args{
				alloc: allocation{
					"QQQ": 100,
				},
				n: 1,
			},
			want: allocation{
				"QQQ": 100,
			},
		},
		{
			name: "take more",
			args: args{
				alloc: allocation{
					"QQQ": 100,
				},
				n: 100,
			},
			want: allocation{
				"QQQ": 100,
			},
		},
		{
			name: "take less",
			args: args{
				alloc: allocation{
					"A": 1,
					"B": 2,
				},
				n: 1,
			},
			want: allocation{
				"B": 2,
			},
		},
		{
			name: "take less with equal",
			args: args{
				alloc: allocation{
					"A": 1,
					"B": 1,
				},
				n: 1,
			},
			want: allocation{
				"A": 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, truncateAfterNthLargest(tt.args.alloc, tt.args.n), "truncateAfterNthLargest(%v, %v)", tt.args.alloc, tt.args.n)
		})
	}
}
