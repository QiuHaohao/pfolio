package metamodel

import (
	"fmt"
	"sort"

	"github.com/qiuhaohao/pfolio/internal/db"
)

const (
	TypeSyfe = "Syfe"
)

type Metamodel interface {
	GetModelEntries() (db.ModelEntries, error)
}

type metamodel struct {
	client
	commonInfo db.CommonMetamodelInfo
}

func New(mmInfo db.Metamodel) (Metamodel, error) {
	mm := metamodel{
		client:     nil,
		commonInfo: mmInfo.CommonMetamodelInfo,
	}
	switch mmInfo.MetamodelType {
	case TypeSyfe:
		mm.client = NewSyfeClient(mmInfo.SyfeMetamodelInfo)
	default:
		return nil, newErrUnknownMetamodelType(mmInfo.MetamodelType)
	}

	return mm, nil
}

type errUnknownMetamodelType struct {
	MetamodelType string
}

func newErrUnknownMetamodelType(metamodelType string) *errUnknownMetamodelType {
	return &errUnknownMetamodelType{MetamodelType: metamodelType}
}

func (e errUnknownMetamodelType) Error() string {
	return fmt.Sprintf("unknown metamodel type: %s", e.MetamodelType)
}

func (m metamodel) GetModelEntries() (db.ModelEntries, error) {
	alloc, err := m.client.GetAllocation()
	if err != nil {
		return nil, err
	}

	if m.commonInfo.TruncateAfterNthLargest > 0 {
		alloc = truncateAfterNthLargest(alloc, int(m.commonInfo.TruncateAfterNthLargest))
	}

	if m.commonInfo.NeedMapSymbols {
		alloc, err = mapSymbols(alloc, m.commonInfo.SymbolMapping, m.commonInfo.RequireMappingForAllSymbols)
		if err != nil {
			return nil, err
		}
	}

	return alloc.ToModelEntries(m.commonInfo.EquivalenceMapping), nil
}

func NewSyfeMetamodel(syfeInfo db.SyfeMetamodelInfo, commonInfo db.CommonMetamodelInfo) Metamodel {
	return metamodel{
		client:     NewSyfeClient(syfeInfo),
		commonInfo: commonInfo,
	}
}

type client interface {
	GetAllocation() (allocation, error)
}

type allocation map[string]float64

func (al allocation) ToModelEntries(equivalenceMapping map[string][]string) db.ModelEntries {
	res := make(db.ModelEntries, 0)
	for symbol, weight := range al {
		res = append(res, db.ModelEntry{
			InstrumentIdentifier:  symbol,
			Weight:                weight,
			EquivalentInstruments: equivalenceMapping[symbol],
		})
	}

	return res
}

type SymbolMapping map[string]string

func truncateAfterNthLargest(alloc allocation, n int) allocation {
	if n >= len(alloc) {
		return alloc
	}

	type item struct {
		Symbol string
		Weight float64
	}

	items := make([]item, 0)
	for s, weight := range alloc {
		items = append(items, item{
			Symbol: s,
			Weight: weight,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Weight > items[j].Weight || (items[i].Weight == items[j].Weight &&
			items[i].Symbol < items[j].Symbol)
	})

	items = items[:n]
	res := make(allocation)
	for _, i := range items {
		res[i.Symbol] = i.Weight
	}

	return res
}

func mapSymbols(alloc allocation, symbolMapping SymbolMapping, requireMappingForAllSymbols bool) (allocation, error) {
	res := make(allocation)
	for s, weight := range alloc {
		mappedSymbol, ok := symbolMapping[s]
		if !ok {
			if requireMappingForAllSymbols {
				return nil, NewSymbolMissingFromMappingError(s)
			} else {
				res[s] = res[s] + weight
				continue
			}
		}
		res[mappedSymbol] = res[mappedSymbol] + weight
	}
	return res, nil
}
