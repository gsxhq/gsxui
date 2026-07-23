// Package tabs holds the site's example gsx components for ui/tabs.
package tabs

import "github.com/gsxhq/gsxui/ui"

// Basic renders a 3-tab group with "account" selected on first paint —
// selected is server-resolved by comparing each part's value to the root's
// (see ui/tabs/tabs.gsx).
component Basic() {
	<ui.Tabs value="account">
		<ui.TabsList>
			<ui.TabsTrigger value="account" selected>Account</ui.TabsTrigger>
			<ui.TabsTrigger value="password">Password</ui.TabsTrigger>
			<ui.TabsTrigger value="team">Team</ui.TabsTrigger>
		</ui.TabsList>
		<ui.TabsContent value="account" selected>Update your account details.</ui.TabsContent>
		<ui.TabsContent value="password">Change your password.</ui.TabsContent>
		<ui.TabsContent value="team">Manage your team members.</ui.TabsContent>
	</ui.Tabs>
}
