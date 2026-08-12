package pages

import "github.com/gsxhq/gsxui/ui"

type themePickerChoice struct {
	Title       string
	Swatch      string
	Value       string
	SwatchStyle string
}

func themePickerSelectedTitle(choices []themePickerChoice, selected string) string {
	for _, choice := range choices {
		if choice.Value == selected {
			return choice.Title
		}
	}
	return "Custom"
}

func themePickerSelectedSwatchStyle(choices []themePickerChoice, selected string) string {
	for _, choice := range choices {
		if choice.Value == selected {
			return themePickerSwatchStyle(choice)
		}
	}
	return ""
}

func themePickerSwatchStyle(choice themePickerChoice) string {
	if choice.Swatch != "" {
		return "background-color: " + choice.Swatch
	}
	return choice.SwatchStyle
}

component ThemePicker(name, label, selected string, choices []themePickerChoice) {
	<ui.Popover data-theme-picker={name} class="min-w-[180px]">
		<ui.PopoverTrigger
			data-theme-picker-trigger
			class="flex min-h-10 w-full items-center justify-between gap-3 rounded-md border border-input bg-transparent px-3 py-1.5 text-left text-sm shadow-xs outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
		>
			<span class="flex min-w-0 flex-col gap-0.5"><span class="text-xs text-muted-foreground">{ label }</span><span data-theme-selection-value class="truncate font-medium">{ themePickerSelectedTitle(choices, selected) }</span></span>
			<span data-theme-selection-swatch class="size-4 shrink-0 rounded-full border border-border bg-foreground/40" style={ themePickerSelectedSwatchStyle(choices, selected) }></span>
		</ui.PopoverTrigger>
		<ui.PopoverContent class="max-h-80 w-72 overflow-y-auto p-3">
			<fieldset class="grid gap-1">
				<legend class="sr-only">{ label }</legend>
				{ for _, choice := range choices {
					<label
						data-theme-choice={name}
						class="flex cursor-pointer items-center gap-2 rounded-md px-2 py-1.5 text-sm hover:bg-accent"
					>
						<ui.Radio
							name={"theme-picker-" + name}
							value={choice.Value}
							aria-label={choice.Title}
							checked={choice.Value == selected}
						/>
						<span data-theme-choice-swatch class="size-4 shrink-0 rounded-full border border-border bg-foreground/40" style={ themePickerSwatchStyle(choice) }></span>
						<span data-theme-choice-label>{ choice.Title }</span>
					</label>
				} }
			</fieldset>
		</ui.PopoverContent>
	</ui.Popover>
}
