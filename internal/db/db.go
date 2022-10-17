package db

import (
	"errors"
	"github.com/QiuHaohao/pfolio/internal/clock"
	"github.com/QiuHaohao/pfolio/internal/config"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

var database *Database

func Get() *Database {
	return database
}

var (
	ErrDuplicatedModelName = errors.New("model name already exists")
)

func Load() {
	database = new(Database)

	f, err := os.Open(viper.GetString(config.KeyDB))
	if os.IsNotExist(err) {
		database = &Database{
			Models: make(map[string]Model),
		}
		return
	} else if err != nil {
		log.Fatal(err)
	}

	d := yaml.NewDecoder(f)
	if err = d.Decode(database); err != nil {
		log.Fatal(err)
	}

	f.Close()
}

func Persist() {
	f, err := os.CreateTemp("", "")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.NewEncoder(f).Encode(database)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	os.Rename(f.Name(), viper.GetString(config.KeyDB))
}

type Database struct {
	Models map[string]Model `yaml:"models"`
}

func (d *Database) AddModel(name string, entries []ModelEntry, isDerivedFromMetaModel bool) error {
	if _, ok := d.Models[name]; ok {
		return ErrDuplicatedModelName
	}

	d.SetModel(name, NewModel(entries, isDerivedFromMetaModel))
	return nil
}

func (d *Database) SetModel(name string, model Model) {
	d.Models[name] = model
}

func (d *Database) GetModel(name string) (Model, bool) {
	m, ok := d.Models[name]
	return m, ok
}

func (d *Database) CheckIsNewModelName(name string) error {
	if _, ok := d.GetModel(name); ok {
		return ErrDuplicatedModelName
	}
	return nil
}

func NewModel(entries []ModelEntry, isDerivedFromMetaModel bool) Model {
	return Model{
		Entries:                entries,
		IsDerivedFromMetaModel: isDerivedFromMetaModel,
		CreateTime:             clock.Now(),
		UpdateTime:             clock.Now(),
	}
}

type Model struct {
	Entries                ModelEntries `yaml:"entries"`
	IsDerivedFromMetaModel bool         `yaml:"is_derived_from_meta_model"`
	CreateTime             time.Time    `yaml:"create_time"`
	UpdateTime             time.Time    `yaml:"update_time"`
}

type InstrumentIdentifier string

type ModelEntries []ModelEntry

var (
	ErrDuplicatedInstrumentIdentifierInModel = errors.New("duplicated instrument identifiers found in model")
	ErrNegativeWeight                        = errors.New("negative weight found")
	ErrDuplicatedEquivalentInstrument        = errors.New("equivalent instrument found to be associated with multiple instruments in model")
	ErrEquivalentInstrumentInModel           = errors.New("equivalent instrument found to be in model")
)

func (es ModelEntries) Validate() error {
	// identifiers must be unique
	iiMap := make(map[InstrumentIdentifier]struct{})
	for _, e := range es {
		if _, ok := iiMap[e.InstrumentIdentifier]; ok {
			return ErrDuplicatedInstrumentIdentifierInModel
		}
		iiMap[e.InstrumentIdentifier] = struct{}{}
	}
	// no negative weight
	for _, e := range es {
		if e.Weight < 0 {
			return ErrNegativeWeight
		}
	}

	// equivalent instruments must be associated to only one instrument in the model
	eiMap := make(map[InstrumentIdentifier]struct{})
	for _, e := range es {
		entryEIs := make([]InstrumentIdentifier, 0)
		for _, ei := range e.EquivalentInstruments {
			if _, ok := eiMap[ei]; ok {
				return ErrDuplicatedEquivalentInstrument
			}
			entryEIs = append(entryEIs, ei)
		}

		for _, ei := range entryEIs {
			eiMap[ei] = struct{}{}
		}
	}

	// equivalent instruments must not be in the model
	for ei := range eiMap {
		if _, ok := iiMap[ei]; ok {
			return ErrEquivalentInstrumentInModel
		}
	}

	return nil
}

type ModelEntry struct {
	InstrumentIdentifier  InstrumentIdentifier   `yaml:"instrument_identifier"`
	Weight                float64                `yaml:"weight"`
	EquivalentInstruments []InstrumentIdentifier `yaml:"equivalent_instruments"`
}
