package syfe

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

const endpoint = "https://api.syfe.com/portfolios/data"

type PortfolioDataResponseMeta map[string]PortfolioDataResponseMetaEntry
type PortfolioDataResponseMetaEntry struct {
	Key             string `json:"key,omitempty"`
	PlaceHolderRisk int    `json:"placeHolderRisk,omitempty"`
}

type Symbol string
type Allocation map[Symbol]float64
type SymbolMapping map[Symbol]Symbol

type PortfolioData struct {
	EtfAllocation map[string]Allocation `json:"etfAllocation"`
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

func GetAllocation(portfolioType string) (Allocation, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("types", portfolioType)
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
	portfolioDataRaw, ok := dr[portfolioType]
	if !ok {
		return nil, errors.New("portfolio type not found in response")
	}

	if err = convert(portfolioDataRaw, portfolioData); err != nil {
		return nil, err
	}

	meta := make(PortfolioDataResponseMeta)
	metaRaw, ok := dr["meta"]
	if !ok {
		return nil, errors.New("portfolio meta not found in response")
	}

	if err = convert(metaRaw, &meta); err != nil {
		return nil, err
	}

	return portfolioData.EtfAllocation[strconv.Itoa(meta[portfolioType].PlaceHolderRisk)], nil
}
