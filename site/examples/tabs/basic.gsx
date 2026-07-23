// Package tabs holds the site's example gsx components for ui/tabs.
package tabs

import uitabs "github.com/gsxhq/gsxui/ui/tabs"

// Basic renders a 3-tab group with "account" selected on first paint —
// selected is server-resolved by comparing each part's value to the root's
// (see ui/tabs/tabs.gsx).
component Basic() {
	<uitabs.Tabs value="account">
		<uitabs.TabsList>
			<uitabs.TabsTrigger value="account" selected>Account</uitabs.TabsTrigger>
			<uitabs.TabsTrigger value="password">Password</uitabs.TabsTrigger>
			<uitabs.TabsTrigger value="team">Team</uitabs.TabsTrigger>
		</uitabs.TabsList>
		<uitabs.TabsContent value="account" selected>Update your account details.</uitabs.TabsContent>
		<uitabs.TabsContent value="password">Change your password.</uitabs.TabsContent>
		<uitabs.TabsContent value="team">Manage your team members.</uitabs.TabsContent>
	</uitabs.Tabs>
}
