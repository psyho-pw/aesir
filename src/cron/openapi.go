package cron

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type OpenApiResponse struct {
	Response Response `json:"response"`
}

type Response struct {
	Body   OpenApiBody   `json:"body"`
	Header OpenApiHeader `json:"header"`
}

type OpenApiBody struct {
	Items      OpenApiItems `json:"items"`
	NumOfRows  int          `json:"numOfRows"`
	PageNo     int          `json:"pageNo"`
	TotalCount int          `json:"totalCount"`
}

type OpenApiItems struct {
	Item []OpenApiItem `json:"item"`
}

type OpenApiItem struct {
	DateKind  string `json:"dateKind"`
	DateName  string `json:"dateName"`
	IsHoliday string `json:"isHoliday"`
	LocDate   int    `json:"locDate"`
	Seq       int    `json:"seq"`
}

type OpenApiHeader struct {
	ResultCode string `json:"resultCode"`
	ResultMsg  string `json:"resultMsg"`
}

func (oi *OpenApiItems) UnmarshalJSON(data []byte) error {
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return err
	}

	itemData, ok := rawData["item"]
	if !ok {
		return fmt.Errorf("'item' field not found")
	}

	switch itemData.(type) {
	case []interface{}:
		var items []OpenApiItem
		if err := mapstructure.Decode(itemData, &items); err != nil {
			return err
		}
		oi.Item = items
	case map[string]interface{}:
		var item OpenApiItem
		if err := mapstructure.Decode(itemData, &item); err != nil {
			return err
		}
		oi.Item = []OpenApiItem{item}
	default:
		return fmt.Errorf("unsupported 'item' data type")
	}

	return nil
}
