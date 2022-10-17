package db

import (
	"fmt"
	"github.com/QiuHaohao/pfolio/internal/clock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDatabase_AddModel(t *testing.T) {
	now := time.Now()
	clock.SetNowFn(func() time.Time {
		return now
	})

	type fields struct {
		Models map[string]Model
	}
	type args struct {
		name                   string
		entries                []ModelEntry
		isDerivedFromMetaModel bool
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		expectedRes *Database
		expectedErr error
	}{
		{
			name: "add to empty",
			fields: struct{ Models map[string]Model }{
				Models: map[string]Model{},
			},
			args: struct {
				name                   string
				entries                []ModelEntry
				isDerivedFromMetaModel bool
			}{
				name: "portfolio1",
				entries: []ModelEntry{
					{
						InstrumentIdentifier: "VOO",
						Weight:               1,
					},
				},
				isDerivedFromMetaModel: true,
			},
			expectedRes: &Database{Models: map[string]Model{
				"portfolio1": {
					Entries: []ModelEntry{
						{
							InstrumentIdentifier: "VOO",
							Weight:               1,
						},
					},
					IsDerivedFromMetaModel: true,
					CreateTime:             now,
					UpdateTime:             now,
				},
			}},
			expectedErr: nil,
		},
		{
			name: "add to nonempty",
			fields: struct{ Models map[string]Model }{
				Models: map[string]Model{
					"portfolio1": {
						Entries: []ModelEntry{
							{
								InstrumentIdentifier: "VOO",
								Weight:               1,
							},
						},
						IsDerivedFromMetaModel: true,
						CreateTime:             now,
						UpdateTime:             now,
					},
				},
			},
			args: struct {
				name                   string
				entries                []ModelEntry
				isDerivedFromMetaModel bool
			}{
				name: "portfolio2",
				entries: []ModelEntry{
					{
						InstrumentIdentifier: "QQQ",
						Weight:               2,
					},
				},
				isDerivedFromMetaModel: false,
			},
			expectedRes: &Database{Models: map[string]Model{
				"portfolio1": {
					Entries: []ModelEntry{
						{
							InstrumentIdentifier: "VOO",
							Weight:               1,
						},
					},
					IsDerivedFromMetaModel: true,
					CreateTime:             now,
					UpdateTime:             now,
				},
				"portfolio2": {
					Entries: []ModelEntry{
						{
							InstrumentIdentifier: "QQQ",
							Weight:               2,
						},
					},
					IsDerivedFromMetaModel: false,
					CreateTime:             now,
					UpdateTime:             now,
				},
			}},
			expectedErr: nil,
		},
		{
			name: "duplicate",
			fields: struct{ Models map[string]Model }{
				Models: map[string]Model{
					"portfolio1": {
						Entries: []ModelEntry{
							{
								InstrumentIdentifier: "VOO",
								Weight:               1,
							},
						},
						IsDerivedFromMetaModel: true,
						CreateTime:             now,
						UpdateTime:             now,
					},
				},
			},
			args: struct {
				name                   string
				entries                []ModelEntry
				isDerivedFromMetaModel bool
			}{
				name: "portfolio1",
				entries: []ModelEntry{
					{
						InstrumentIdentifier: "QQQ",
						Weight:               2,
					},
				},
				isDerivedFromMetaModel: false,
			},
			expectedRes: &Database{Models: map[string]Model{
				"portfolio1": {
					Entries: []ModelEntry{
						{
							InstrumentIdentifier: "VOO",
							Weight:               1,
						},
					},
					IsDerivedFromMetaModel: true,
					CreateTime:             now,
					UpdateTime:             now,
				},
			}},
			expectedErr: ErrDuplicatedModelName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Database{
				Models: tt.fields.Models,
			}
			err := d.AddModel(tt.args.name, tt.args.entries, tt.args.isDerivedFromMetaModel)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedRes, d)
		})
	}
}

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
					EquivalentInstruments: []InstrumentIdentifier{},
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
					EquivalentInstruments: []InstrumentIdentifier{"CSPX"},
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
					EquivalentInstruments: []InstrumentIdentifier{"CSPX", "CSPX"},
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
					EquivalentInstruments: []InstrumentIdentifier{},
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
					EquivalentInstruments: []InstrumentIdentifier{},
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
					EquivalentInstruments: []InstrumentIdentifier{"VOO"},
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
					EquivalentInstruments: []InstrumentIdentifier{"VOO"},
				},
				{
					InstrumentIdentifier:  "ABC",
					Weight:                1,
					EquivalentInstruments: []InstrumentIdentifier{"VOO"},
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
