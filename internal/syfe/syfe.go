package syfe

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"
)

const endpoint = "https://api.syfe.com/portfolios/data"

var (
	ErrPortfolioTypeNotFoundInResponse = errors.New("portfolio type not found in response")
	ErrPortfolioMetaNotFoundInResponse = errors.New("portfolio meta not found in response")
	ErrSymbolMissingFromMapping        = errors.New("symbol missing from mapping")
)

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

func TruncateAfterNthLargest(allocation Allocation, n int) Allocation {
	if n >= len(allocation) {
		return allocation
	}

	type item struct {
		Symbol Symbol
		Weight float64
	}

	items := make([]item, 0)
	for s, weight := range allocation {
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
	res := make(Allocation)
	for _, i := range items {
		res[i.Symbol] = i.Weight
	}

	return res
}

func MapSymbols(allocation Allocation, symbolMapping SymbolMapping) (Allocation, error) {
	res := make(Allocation)
	for s, weight := range allocation {
		mappedSymbol, ok := symbolMapping[s]
		if !ok {
			return nil, ErrSymbolMissingFromMapping
		}
		res[mappedSymbol] = res[mappedSymbol] + weight
	}
	return res, nil
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
		return nil, ErrPortfolioTypeNotFoundInResponse
	}

	if err = convert(portfolioDataRaw, portfolioData); err != nil {
		return nil, err
	}

	meta := make(PortfolioDataResponseMeta)
	metaRaw, ok := dr["meta"]
	if !ok {
		return nil, ErrPortfolioMetaNotFoundInResponse
	}

	if err = convert(metaRaw, &meta); err != nil {
		return nil, err
	}

	return portfolioData.EtfAllocation[strconv.Itoa(meta[portfolioType].PlaceHolderRisk)], nil
}
