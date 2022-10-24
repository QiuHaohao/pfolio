package metamodel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/qiuhaohao/pfolio/internal/db"
)

const syfeEndpoint = "https://api.syfe.com/portfolios/data"

const (
	SyfePortfolioTypeGlobalEquity100      = "GLOBAL_EQUITY_100"
	SyfePortfolioTypeGlobalRiskManaged    = "GLOBAL_RISK_MANAGED"
	SyfePortfolioTypeREITRiskManaged      = "REIT_RISK_MANAGED"
	SyfePortfolioTypeREIT                 = "REIT"
	SyfePortfolioTypeCashPlus             = "CASH_PLUS"
	SyfePortfolioTypeCoreGrowth           = "CORE_GROWTH"
	SyfePortfolioTypeCoreBalanced         = "CORE_BALANCED"
	SyfePortfolioTypeCoreDefensive        = "CORE_DEFENSIVE"
	SyfePortfolioTypeESGAndCleanEnergy    = "ESG_AND_CLEAN_ENERGY"
	SyfePortfolioTypeDisruptiveTechnology = "DISRUPTIVE_TECHNOLOGY"
	SyfePortfolioTypeGlobalIncome         = "GLOBAL_INCOME"
	SyfePortfolioTypeChinaGrowth          = "CHINA_GROWTH"
	SyfePortfolioTypeHealthcareInnovation = "HEALTHCARE_INNOVATION"
)

var (
	ErrSyfePortfolioTypeNotFoundInResponse = errors.New("portfolio type not found in response")
	ErrSyfePortfolioMetaNotFoundInResponse = errors.New("portfolio meta not found in response")
)

type SyfePortfolioDataResponseMeta map[string]SyfePortfolioDataResponseMetaEntry
type SyfePortfolioDataResponseMetaEntry struct {
	Key             string `json:"key,omitempty"`
	PlaceHolderRisk int    `json:"placeHolderRisk,omitempty"`
}

type PortfolioData struct {
	EtfAllocation map[string]allocation `json:"etfAllocation"`
}

type syfeClient struct {
	portfolioType string
}

func NewSyfeClient(syfeInfo db.SyfeMetamodelInfo) *syfeClient {
	return &syfeClient{portfolioType: syfeInfo.Type}
}

func (c syfeClient) GetAllocation() (allocation, error) {
	req, err := http.NewRequest(http.MethodGet, syfeEndpoint, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("types", c.portfolioType)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type dataResponse map[string]interface{}
	dr := make(dataResponse)

	if err = json.Unmarshal(raw, &dr); err != nil {
		return nil, err
	}

	portfolioData := new(PortfolioData)
	portfolioDataRaw, ok := dr[c.portfolioType]
	if !ok {
		return nil, ErrSyfePortfolioTypeNotFoundInResponse
	}

	if err = convert(portfolioDataRaw, portfolioData); err != nil {
		return nil, err
	}

	meta := make(SyfePortfolioDataResponseMeta)
	metaRaw, ok := dr["meta"]
	if !ok {
		return nil, ErrSyfePortfolioMetaNotFoundInResponse
	}

	if err = convert(metaRaw, &meta); err != nil {
		return nil, err
	}

	return portfolioData.EtfAllocation[strconv.Itoa(meta[c.portfolioType].PlaceHolderRisk)], nil
}

func convert(src, dst interface{}) error {
	raw, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(raw, dst)
	if err != nil {
		return err
	}

	return nil
}

func NewSymbolMissingFromMappingError(missingSymbol string) error {
	return symbolMissingFromMappingError{
		missingSymbol: missingSymbol,
	}
}

type symbolMissingFromMappingError struct {
	missingSymbol string
}

func (s symbolMissingFromMappingError) Error() string {
	return fmt.Sprintf("symbol %s missing from mapping", s.missingSymbol)
}
