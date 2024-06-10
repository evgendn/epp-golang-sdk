package eppoclient

import (
	"encoding/json"
	"fmt"
	"time"
)

type ufcResponse struct {
	Flags map[string]flagConfiguration `json:"flags"`
}

type flagConfiguration struct {
	Key           string               `json:"key"`
	Enabled       bool                 `json:"enabled"`
	VariationType variationType        `json:"variationType"`
	Variations    map[string]variation `json:"variations"`
	Allocations   []allocation         `json:"allocations"`
	TotalShards   int64                `json:"totalShards"`
}

type variation struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type variationType int

const (
	stringVariation variationType = iota
	integerVariation
	numericVariation
	booleanVariation
	jsonVariation
)

func (v variationType) MarshalJSON() ([]byte, error) {
	switch v {
	case stringVariation:
		return json.Marshal("STRING")
	case integerVariation:
		return json.Marshal("INTEGER")
	case numericVariation:
		return json.Marshal("NUMERIC")
	case booleanVariation:
		return json.Marshal("BOOLEAN")
	case jsonVariation:
		return json.Marshal("JSON")
	default:
		return nil, fmt.Errorf("unsupported variation type: %d", v)
	}
}

func (v *variationType) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch value {
	case "STRING":
		*v = stringVariation
	case "INTEGER":
		*v = integerVariation
	case "NUMERIC":
		*v = numericVariation
	case "BOOLEAN":
		*v = booleanVariation
	case "JSON":
		*v = jsonVariation
	default:
		return fmt.Errorf("unknown variation type: %s", value)
	}

	return nil
}

func (ty variationType) valueToAssignmentValue(value interface{}) (interface{}, error) {
	switch ty {
	case stringVariation:
		s := value.(string)
		return s, nil
	case integerVariation:
		f64 := value.(float64)
		i64 := int64(f64)
		if f64 == float64(i64) {
			return i64, nil
		} else {
			return nil, fmt.Errorf("failed to convert number to integer")
		}
	case numericVariation:
		number := value.(float64)
		return number, nil
	case booleanVariation:
		v := value.(bool)
		return v, nil
	case jsonVariation:
		v := value.(string)
		var result interface{}
		err := json.Unmarshal([]byte(v), &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unexpected variation type: %v", ty)
	}
}

type allocation struct {
	Key     string    `json:"key"`
	Rules   []rule    `json:"rules"`
	StartAt time.Time `json:"startAt"`
	EndAt   time.Time `json:"endAt"`
	Splits  []split   `json:"splits"`
	DoLog   *bool     `json:"doLog"`
}

type rule struct {
	Conditions []condition `json:"conditions"`
}

type condition struct {
	Operator  string      `json:"operator"`
	Attribute string      `json:"attribute"`
	Value     interface{} `json:"value"`
}

type split struct {
	Shards       []shard           `json:"shards"`
	VariationKey string            `json:"variationKey"`
	ExtraLogging map[string]string `json:"extraLogging"`
}

type shard struct {
	Salt   string       `json:"salt"`
	Ranges []shardRange `json:"ranges"`
}

type shardRange struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}
