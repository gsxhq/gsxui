package recipe

import (
	"bytes"
	"encoding/json"
	"slices"
)

// ContractVersion identifies the shape of the emitted schema. A consumer
// reads it to decide whether it understands the artifact, and should reject a
// version it does not recognize. It is bumped whenever the schema changes
// shape.
const ContractVersion = 1

// Contract is the serialized recipe model consumed by the theme editor. These
// types are marshalled, so their fields are exported.
type Contract struct {
	Version    int                                     `json:"version"`
	Components map[string]ContractComponent            `json:"components"`
	Styles     map[string]map[string]ContractUtilities `json:"styles"`
}

type ContractComponent struct {
	Slots map[string]ContractSlot `json:"slots"`
}

type ContractSlot struct {
	Base       bool                         `json:"base"`
	Dimensions map[string]ContractDimension `json:"dimensions"`
}

type ContractDimension struct {
	Default string   `json:"default"`
	Values  []string `json:"values"`
}

type ContractUtilities struct {
	Slots map[string]ContractSlotUtilities `json:"slots"`
}

type ContractSlotUtilities struct {
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
		slots := make(map[string]ContractSlot, len(shape.Slots))
		for _, slot := range shape.Slots {
			dimensions := make(map[string]ContractDimension, len(slot.Dimensions))
			for _, dimension := range slot.Dimensions {
				dimensions[dimension.Name] = ContractDimension{
					Default: dimension.Default,
					Values:  slices.Clone(dimension.Values),
				}
			}
			slots[slot.Name] = ContractSlot{Base: slot.Base, Dimensions: dimensions}
		}
		contract.Components[name] = ContractComponent{Slots: slots}
	}
	for style, resolvedComponents := range styles {
		entry := make(map[string]ContractUtilities, len(resolvedComponents))
		for name, resolved := range resolvedComponents {
			slots := make(map[string]ContractSlotUtilities, len(resolved.Shape.Slots))
			for _, slot := range resolved.Shape.Slots {
				values := make(map[string]map[string][]string, len(slot.Dimensions))
				for _, dimension := range slot.Dimensions {
					byValue := make(map[string][]string, len(dimension.Values))
					for _, value := range dimension.Values {
						byValue[value] = resolved.Utilities(slot.Name, dimension.Name, value)
					}
					values[dimension.Name] = byValue
				}
				slots[slot.Name] = ContractSlotUtilities{
					Base:       resolved.BaseUtilities(slot.Name),
					Dimensions: values,
				}
			}
			entry[name] = ContractUtilities{Slots: slots}
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
