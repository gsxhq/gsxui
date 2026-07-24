package selectbox

import "github.com/gsxhq/gsxui/ui"

// Scrollable is the long-list stress test: five SelectGroups spanning 27
// world timezones, separated by SelectSeparator. The content's max-h-96 cap
// (the base class) makes the whole viewport scroll natively — this is the
// demo for the "drop the scroll up/down buttons, rely on overflow-y-auto"
// decision (the custom listbox always anchors popper-equivalent, so Radix's
// item-aligned scroll affordances are unnecessary).
component Scrollable() {
	<ui.Select name="timezone">
		<ui.SelectTrigger class="w-[280px]">
			<ui.SelectValue placeholder="Select a timezone"/>
		</ui.SelectTrigger>
		<ui.SelectContent>
			<ui.SelectGroup>
				<ui.SelectLabel>North America</ui.SelectLabel>
				<ui.SelectItem value="est">Eastern Standard Time (EST)</ui.SelectItem>
				<ui.SelectItem value="cst">Central Standard Time (CST)</ui.SelectItem>
				<ui.SelectItem value="mst">Mountain Standard Time (MST)</ui.SelectItem>
				<ui.SelectItem value="pst">Pacific Standard Time (PST)</ui.SelectItem>
				<ui.SelectItem value="akst">Alaska Standard Time (AKST)</ui.SelectItem>
				<ui.SelectItem value="hst">Hawaii Standard Time (HST)</ui.SelectItem>
			</ui.SelectGroup>
			<ui.SelectSeparator/>
			<ui.SelectGroup>
				<ui.SelectLabel>Europe &amp; Africa</ui.SelectLabel>
				<ui.SelectItem value="gmt">Greenwich Mean Time (GMT)</ui.SelectItem>
				<ui.SelectItem value="cet">Central European Time (CET)</ui.SelectItem>
				<ui.SelectItem value="eet">Eastern European Time (EET)</ui.SelectItem>
				<ui.SelectItem value="west">Western European Summer Time (WEST)</ui.SelectItem>
				<ui.SelectItem value="cat">Central Africa Time (CAT)</ui.SelectItem>
				<ui.SelectItem value="eat">East Africa Time (EAT)</ui.SelectItem>
			</ui.SelectGroup>
			<ui.SelectSeparator/>
			<ui.SelectGroup>
				<ui.SelectLabel>Asia</ui.SelectLabel>
				<ui.SelectItem value="msk">Moscow Time (MSK)</ui.SelectItem>
				<ui.SelectItem value="ist">India Standard Time (IST)</ui.SelectItem>
				<ui.SelectItem value="cst_china">China Standard Time (CST)</ui.SelectItem>
				<ui.SelectItem value="jst">Japan Standard Time (JST)</ui.SelectItem>
				<ui.SelectItem value="kst">Korea Standard Time (KST)</ui.SelectItem>
				<ui.SelectItem value="ist_indonesia">Indonesia Central Standard Time (WITA)</ui.SelectItem>
			</ui.SelectGroup>
			<ui.SelectSeparator/>
			<ui.SelectGroup>
				<ui.SelectLabel>Australia &amp; Pacific</ui.SelectLabel>
				<ui.SelectItem value="awst">Australian Western Standard Time (AWST)</ui.SelectItem>
				<ui.SelectItem value="acst">Australian Central Standard Time (ACST)</ui.SelectItem>
				<ui.SelectItem value="aest">Australian Eastern Standard Time (AEST)</ui.SelectItem>
				<ui.SelectItem value="nzst">New Zealand Standard Time (NZST)</ui.SelectItem>
				<ui.SelectItem value="fjt">Fiji Time (FJT)</ui.SelectItem>
			</ui.SelectGroup>
			<ui.SelectSeparator/>
			<ui.SelectGroup>
				<ui.SelectLabel>South America</ui.SelectLabel>
				<ui.SelectItem value="art">Argentina Time (ART)</ui.SelectItem>
				<ui.SelectItem value="bot">Bolivia Time (BOT)</ui.SelectItem>
				<ui.SelectItem value="brt">Brasilia Time (BRT)</ui.SelectItem>
				<ui.SelectItem value="clt">Chile Standard Time (CLT)</ui.SelectItem>
			</ui.SelectGroup>
		</ui.SelectContent>
	</ui.Select>
}
