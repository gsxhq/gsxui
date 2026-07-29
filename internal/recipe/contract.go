package recipe

import (
	"bytes"
	"encoding/json"
	"slices"
)

// ContractVersion is bumped whenever the emitted schema changes shape.
const ContractVersion = 1

// Contract is the serialized recipe model consumed by the theme editor. These
// types are marshalled, so their fields are exported.
type Contract struct {
	Version    int                                     `json:"version"`
	Components map[string]ContractComponent            `json:"components"`
	Styles     map[string]map[string]ContractUtilities `json:"styles"`
}

type ContractComponent struct {
	Base       bool                         `json:"base"`
	Dimensions map[string]ContractDimension `json:"dimensions"`
}

type ContractDimension struct {
	Default string   `json:"default"`
	Values  []string `json:"values"`
}

type ContractUtilities struct {
	Base       []string                       `json:"base,omitempty"`
	Dimensions map[string]map[string][]string `json:"dimensions"`
}

func BuildContract(components map[string]Shape, styles map[string]map[string]Resolved) Contract {
	contract := Contract{
		Version:    ContractVersion,
		Components: make(map[string]ContractComponent, len(components)),
		Styles:     make(map[string]map[string]ContractUtilities, len(styles)),
	}
	for name, shape := range components {
		dimensions := make(map[string]ContractDimension, len(shape.Dimensions))
		for _, dimension := range shape.Dimensions {
			dimensions[dimension.Name] = ContractDimension{
				Default: dimension.Default,
				Values:  slices.Clone(dimension.Values),
			}
		}
		contract.Components[name] = ContractComponent{Base: shape.Base, Dimensions: dimensions}
	}
	for style, resolvedComponents := range styles {
		entry := make(map[string]ContractUtilities, len(resolvedComponents))
		for name, resolved := range resolvedComponents {
			values := make(map[string]map[string][]string, len(resolved.Shape.Dimensions))
			for _, dimension := range resolved.Shape.Dimensions {
				byValue := make(map[string][]string, len(dimension.Values))
				for _, value := range dimension.Values {
					byValue[value] = resolved.Utilities(dimension.Name, value)
				}
				values[dimension.Name] = byValue
			}
			entry[name] = ContractUtilities{Base: slices.Clone(resolved.Base), Dimensions: values}
		}
		contract.Styles[style] = entry
	}
	return contract
}

// MarshalIndent renders the contract deterministically. encoding/json sorts map
// keys, so identical inputs always produce identical bytes.
func (c Contract) MarshalIndent() ([]byte, error) {
	out, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(bytes.TrimRight(out, "\n"), '\n'), nil
}
