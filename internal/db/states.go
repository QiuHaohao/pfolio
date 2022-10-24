package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/qiuhaohao/pfolio/internal/clock"
)

type States interface {
	GetMetamodel(name string) (Metamodel, bool)
	SetMetamodel(name string, mm Metamodel)
	SetModel(name string, model Model)
	GetModel(name string) (Model, bool)
	GetModels() map[string]Model
	RemoveModel(name string)
}

func newInitialStates() *states {
	return &states{
		Models:     make(map[string]Model),
		Metamodels: make(map[string]Metamodel),
	}
}

type states struct {
	Models     map[string]Model     `yaml:"models"`
	Metamodels map[string]Metamodel `yaml:"metamodels"`
}

func (d *states) GetMetamodel(name string) (Metamodel, bool) {
	m, ok := d.Metamodels[name]
	return m, ok
}

func (d *states) SetMetamodel(name string, mm Metamodel) {
	d.Metamodels[name] = mm
}

func (d *states) SetModel(name string, model Model) {
	d.Models[name] = model
}

func (d *states) GetModel(name string) (Model, bool) {
	m, ok := d.Models[name]
	return m, ok
}

func (d *states) GetModels() map[string]Model {
	return d.Models
}

func (d *states) RemoveModel(name string) {
	delete(d.Models, name)
}

type Metamodel struct {
	MetamodelType       string              `yaml:"metamodel_type"`
	CommonMetamodelInfo CommonMetamodelInfo `yaml:"common_metamodel_info"`
	SyfeMetamodelInfo   SyfeMetamodelInfo   `yaml:"syfe_metamodel_info"`
	CreateTime          time.Time           `yaml:"create_time"`
	UpdateTime          time.Time           `yaml:"update_time"`
}

type CommonMetamodelInfo struct {
	TruncateAfterNthLargest     uint                `yaml:"truncate_after_nth_largest_weight"`
	NeedMapSymbols              bool                `yaml:"need_map_symbols"`
	RequireMappingForAllSymbols bool                `yaml:"require_mapping_for_all_symbols"`
	SymbolMapping               map[string]string   `yaml:"symbol_mapping"`
	EquivalenceMapping          map[string][]string `yaml:"equivalence_mapping"`
}

func (mi CommonMetamodelInfo) Validate() error {
	if mi.TruncateAfterNthLargest >= TruncateAfterNthLargestLimit {
		return ErrTruncateAfterNthLargestTooLarge
	}
	return nil
}

type SyfeMetamodels map[string]SyfeMetamodel

type SyfeMetamodel struct {
	Info       SyfeMetamodelInfo `yaml:"info"`
	CreateTime time.Time         `yaml:"create_time"`
	UpdateTime time.Time         `yaml:"update_time"`
}

func NewSyfeMetamodel(info SyfeMetamodelInfo) SyfeMetamodel {
	return SyfeMetamodel{
		Info:       info,
		CreateTime: clock.Now(),
		UpdateTime: clock.Now(),
	}
}

const TruncateAfterNthLargestLimit = 100

var (
	ErrTruncateAfterNthLargestTooLarge = fmt.Errorf(
		"truncate_after_nth_largest_weight is too large, must be smaller than %d",
		TruncateAfterNthLargestLimit)
)

type SyfeMetamodelInfo struct {
	Type string `yaml:"type"`
}

func (mi SyfeMetamodelInfo) Validate() error {
	return nil
}

type Model struct {
	Entries                ModelEntries `yaml:"entries"`
	IsDerivedFromMetaModel bool         `yaml:"is_derived_from_meta_model"`
	CreateTime             time.Time    `yaml:"create_time"`
	UpdateTime             time.Time    `yaml:"update_time"`
}

type ModelEntries []ModelEntry

var (
	ErrDuplicatedInstrumentIdentifierInModel = errors.New("duplicated instrument identifiers found in model")
	ErrNegativeWeight                        = errors.New("negative weight found")
	ErrDuplicatedEquivalentInstrument        = errors.New("equivalent instrument found to be associated with multiple instruments in model")
	ErrEquivalentInstrumentInModel           = errors.New("equivalent instrument found to be in model")
)

func (es ModelEntries) TotalWeight() float64 {
	totalWeight := float64(0)
	for _, e := range es {
		totalWeight += e.Weight
	}

	return totalWeight
}

func (es ModelEntries) Validate() error {
	// identifiers must be unique
	iiMap := make(map[string]struct{})
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
	eiMap := make(map[string]struct{})
	for _, e := range es {
		entryEIs := make([]string, 0)
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
	InstrumentIdentifier  string   `yaml:"instrument_identifier"`
	Weight                float64  `yaml:"weight"`
	EquivalentInstruments []string `yaml:"equivalent_instruments"`
}

func (e ModelEntry) GetStringsEquivalentInstruments() (ss []string) {
	for _, ei := range e.EquivalentInstruments {
		ss = append(ss, string(ei))
	}
	return
}
