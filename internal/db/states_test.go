package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModelEntries_Validate(t *testing.T) {
	tests := []struct {
		name     string
		es       ModelEntries
		checkErr assert.ErrorAssertionFunc
	}{
		{
			name: "empty is valid",
			es:   ModelEntries{},
			checkErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "single entry positive weight",
			es: ModelEntries{
				{
					InstrumentIdentifier:  "VOO",
					Weight:                1,
					EquivalentInstruments: []string{},
				},
			},
			checkErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "single entry positive weight with equivalent instrument",
			es: ModelEntries{
				{
					InstrumentIdentifier:  "VOO",
					Weight:                1,
					EquivalentInstruments: []string{"CSPX"},
				},
			},
			checkErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "duplicated equivalent instrument in same entry",
			es: ModelEntries{
				{
					InstrumentIdentifier:  "VOO",
					Weight:                1,
					EquivalentInstruments: []string{"CSPX", "CSPX"},
				},
			},
			checkErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "zero weight",
			es: ModelEntries{
				{
					InstrumentIdentifier:  "VOO",
					Weight:                0,
					EquivalentInstruments: []string{},
				},
			},
			checkErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "negative weight",
			es: ModelEntries{
				{
					InstrumentIdentifier:  "VOO",
					Weight:                0,
					EquivalentInstruments: []string{},
				},
			},
			checkErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == ErrNegativeWeight
			},
		},
		{
			name: "equivalent instrument in model",
			es: ModelEntries{
				{
					InstrumentIdentifier:  "VOO",
					Weight:                1,
					EquivalentInstruments: []string{"VOO"},
				},
			},
			checkErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == ErrEquivalentInstrumentInModel
			},
		},
		{
			name: "equivalent instrument associated to multiple instruments in model",
			es: ModelEntries{
				{
					InstrumentIdentifier:  "VOO",
					Weight:                1,
					EquivalentInstruments: []string{"VOO"},
				},
				{
					InstrumentIdentifier:  "ABC",
					Weight:                1,
					EquivalentInstruments: []string{"VOO"},
				},
			},
			checkErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == ErrDuplicatedEquivalentInstrument
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.checkErr(t, tt.es.Validate(), fmt.Sprintf("Validate()"))
		})
	}
}
