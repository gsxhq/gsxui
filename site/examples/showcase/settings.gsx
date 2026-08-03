package showcase

import "github.com/gsxhq/gsxui/ui"

// SettingsCard composes Switch, NativeSelect, Slider and Separator into a
// small preferences panel.
component SettingsCard() {
	<ui.Card>
		<ui.CardHeader>
			<ui.CardTitle>Preferences</ui.CardTitle>
			<ui.CardDescription>Manage how the workspace behaves.</ui.CardDescription>
		</ui.CardHeader>
		<ui.CardContent class="flex flex-col gap-4">
			<div class="flex items-center justify-between gap-4">
				<ui.Label for="home-showcase-notifications">Notifications</ui.Label>
				<ui.Switch id="home-showcase-notifications" checked/>
			</div>
			<div class="flex items-center justify-between gap-4">
				<ui.Label for="home-showcase-autosave">Autosave</ui.Label>
				<ui.Switch id="home-showcase-autosave"/>
			</div>
			<ui.Separator/>
			<div class="grid gap-2">
				<ui.Label for="home-showcase-theme">Theme</ui.Label>
				<ui.NativeSelect id="home-showcase-theme" name="home-showcase-theme">
					<ui.NativeSelectOption value="system" selected>System</ui.NativeSelectOption>
					<ui.NativeSelectOption value="light">Light</ui.NativeSelectOption>
					<ui.NativeSelectOption value="dark">Dark</ui.NativeSelectOption>
				</ui.NativeSelect>
			</div>
			<ui.Separator/>
			<div class="grid gap-2">
				<ui.Label for="home-showcase-density">Density</ui.Label>
				<ui.Slider id="home-showcase-density" value={40} min={0} max={100} step={10} aria-label="Density"/>
			</div>
		</ui.CardContent>
	</ui.Card>
}
