package stylecontract

import (
	"reflect"
	"testing"
)

func TestCompositeContractsDeclareEverySlotAndPresentationAxis(t *testing.T) {
	orientation := []Axis{{Attribute: "data-orientation", Values: []string{"horizontal", "vertical"}}}
	want := []Component{
		{
			Name:         "Calendar",
			RegistryName: "calendar",
			Slots: []Slot{
				{Name: "calendar", Axes: []Axis{
					{Attribute: "data-caption-layout", Values: []string{"dropdown", "label"}},
					{Attribute: "data-gsxui-calendar-mode", Values: []string{"multiple", "range", "single"}},
					{Attribute: "data-gsxui-calendar-show-outside-days", Values: []string{"false", "true"}},
				}},
				{Name: "calendar-months"},
				{Name: "calendar-nav"},
				{Name: "calendar-nav-button", Axes: []Axis{
					{Attribute: "data-variant", Values: []string{"ghost"}},
					{Attribute: "data-size", Values: []string{"icon"}},
					{Attribute: "aria-disabled", Values: []string{"true"}},
				}},
				{Name: "calendar-previous"},
				{Name: "calendar-next"},
				{Name: "calendar-month-caption"},
				{Name: "calendar-dropdowns"},
				{Name: "calendar-dropdown-root"},
				{Name: "calendar-caption", Axes: []Axis{{Attribute: "data-caption-layout", Values: []string{"dropdown", "label"}}}},
				{Name: "calendar-grid", Axes: []Axis{{Attribute: "aria-multiselectable", Values: []string{"true"}}}},
				{Name: "calendar-weekdays"},
				{Name: "calendar-weekday"},
				{Name: "calendar-week"},
				{Name: "calendar-day", Axes: []Axis{
					{Attribute: "data-outside"},
					{Attribute: "data-hidden"},
					{Attribute: "data-today"},
					{Attribute: "data-focused"},
					{Attribute: "data-disabled"},
					{Attribute: "data-selected", Values: []string{"true"}},
					{Attribute: "aria-selected", Values: []string{"false", "true"}},
				}},
				{Name: "calendar-day-button", Axes: []Axis{
					{Attribute: "data-variant", Values: []string{"ghost"}},
					{Attribute: "data-size", Values: []string{"icon"}},
					{Attribute: "data-selected-single", Values: []string{"false", "true"}},
					{Attribute: "data-range-start", Values: []string{"false", "true"}},
					{Attribute: "data-range-middle", Values: []string{"false", "true"}},
					{Attribute: "data-range-end", Values: []string{"false", "true"}},
					{Attribute: "disabled"},
					{Attribute: "aria-disabled", Values: []string{"true"}},
					{Attribute: "aria-hidden", Values: []string{"true"}},
				}},
			},
		},
		{
			Name:         "Carousel",
			RegistryName: "carousel",
			Slots: []Slot{
				{Name: "carousel", Axes: orientation},
				{Name: "carousel-content", Axes: orientation},
				{Name: "carousel-track", Axes: orientation},
				{Name: "carousel-item", Axes: orientation},
				{Name: "carousel-previous", Axes: []Axis{
					{Attribute: "data-orientation", Values: []string{"horizontal", "vertical"}},
					{Attribute: "data-variant", Values: []string{"outline"}},
					{Attribute: "data-size", Values: []string{"icon"}},
					{Attribute: "disabled"},
				}},
				{Name: "carousel-next", Axes: []Axis{
					{Attribute: "data-orientation", Values: []string{"horizontal", "vertical"}},
					{Attribute: "data-variant", Values: []string{"outline"}},
					{Attribute: "data-size", Values: []string{"icon"}},
					{Attribute: "disabled"},
				}},
				{Name: "carousel-control-label"},
			},
		},
		{
			Name:         "Resizable",
			RegistryName: "resizable",
			Slots: []Slot{
				{Name: "resizable-panel-group", Axes: []Axis{{Attribute: "aria-orientation", Values: []string{"horizontal", "vertical"}}}},
				{Name: "resizable-panel"},
				{Name: "resizable-handle", Axes: []Axis{{Attribute: "aria-orientation", Values: []string{"horizontal", "vertical"}}}},
				{Name: "resizable-handle-grip"},
			},
		},
		{
			Name:         "ScrollArea",
			RegistryName: "scroll-area",
			Slots:        []Slot{{Name: "scroll-area", Axes: orientation}},
		},
	}

	if !reflect.DeepEqual(compositeContracts, want) {
		t.Fatalf("compositeContracts = %#v\nwant %#v", compositeContracts, want)
	}
}
