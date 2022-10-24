package view

import (
	"github.com/qiuhaohao/pfolio/internal/db"
	"github.com/qiuhaohao/pfolio/internal/metamodel"
)

type SyfeMetamodelEditView struct {
	CommonMetamodelInfo db.CommonMetamodelInfo `yaml:"common_metamodel_info"`
	SyfeMetamodelInfo   db.SyfeMetamodelInfo   `yaml:"syfe_metamodel_info"`
}

func GetDefaultInitialSyfeMetamodelEditView() SyfeMetamodelEditView {
	return SyfeMetamodelEditView{
		CommonMetamodelInfo: db.CommonMetamodelInfo{
			TruncateAfterNthLargest:     20,
			NeedMapSymbols:              true,
			RequireMappingForAllSymbols: true,
			SymbolMapping: map[string]string{
				"CICT":    "C38U",
				"CAREIT":  "A17U",
				"MLT":     "M44U",
				"MINT":    "ME8U",
				"KDCREIT": "AJBU",
				"SUN":     "T82U",
				"FLT":     "BUOU",
				"MAGIC":   "RW0U",
				"KREIT":   "K71U",
				"FCT":     "J69U",
				"CLCT":    "AU8U",
				"CLAS":    "HMN",
				"EREIT":   "J91U",
				"PREIT":   "C2PU",
				"AAREIT":  "O5RU",
				"CDREIT":  "J85",
				"LREIT":   "JYEU",
				"MPACT":   "N2IU",
				"SPHREIT": "SK6U",
				"CINT":    "CY6U",
				"FEHT":    "Q5T",
			},
		},
		SyfeMetamodelInfo: db.SyfeMetamodelInfo{
			Type: metamodel.SyfePortfolioTypeREIT,
		},
	}
}

func (v SyfeMetamodelEditView) Validate() error {
	if err := v.CommonMetamodelInfo.Validate(); err != nil {
		return err
	}
	if err := v.SyfeMetamodelInfo.Validate(); err != nil {
		return err
	}

	return nil
}

type ModelEditView db.ModelEntries

func NewModelEditViewFromDB(model db.Model) ModelEditView {
	return ModelEditView(model.Entries)
}

func GetDefaultInitialModelEditView() ModelEditView {
	return ModelEditView{
		{
			InstrumentIdentifier: "TLT",
			Weight:               4,
		},
		{
			InstrumentIdentifier: "VOO",
			Weight:               6,
			EquivalentInstruments: []string{
				"CSPX",
			},
		},
	}
}

func (v ModelEditView) Validate() error {
	return db.ModelEntries(v).Validate()
}
