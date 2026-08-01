# RTL examples — batch 2 (collapsible..input-otp)

One RTL example added per component, registered as `Example{Name: "rtl", Title: "RTL"}`,
each wrapped in `<div dir="rtl" lang="ar">`, content translated to Arabic mirroring
shadcn's own `*-rtl.tsx` demos in `apps/v4/examples/base/`. None of the 12 use
`Isolated` or `PreviewRTL`: no example in this batch is registered `Isolated: true` in
its own `.go` file, so Rtl follows suit (grepped `Isolated`/`PreviewRTL` across
`site/examples/*.go` before starting — only `sidebar.go` uses either, and sidebar isn't
in this batch).

## collapsible
Order-status panel (order #4189, shipped status, shipping address, items), same shape
as Basic (title+chevron trigger row, one visible line, `CollapsibleContent` with two
more), same Arabic strings as shadcn's `collapsible-rtl.tsx`. No API gap.

## combobox
Six-category single-select (technology/design/business/marketing/education/health),
same Arabic strings as shadcn's `combobox-rtl.tsx`. **Deviation**: shadcn's own
`combobox-rtl.tsx` is a multi-select chips combobox (`ComboboxChips`/`ComboboxChip`/
`ComboboxValue`); this directory's Basic only composes the single-select shape
(`Combobox`/`ComboboxInput`/`ComboboxContent`/`ComboboxList`/`ComboboxEmpty`/
`ComboboxItem`), so Rtl keeps that shape with the same category list and translations
rather than inventing the chips API.

## command
Inline palette (Suggestions/Settings groups) plus the `CommandDialog` variant, same
Arabic strings as shadcn's `command-rtl.tsx` for the inline palette; the
`CommandDialog` half (not present in the shadcn RTL demo) was translated by hand
following the same vocabulary ("لوحة الأوامر" / "فتح الإعدادات"). No API gap.

## context-menu
Right-click zone with Back/Forward/Reload/Delete, same Arabic strings as shadcn's
`context-menu-rtl.tsx`. **Deviation**: shadcn's own `context-menu-rtl.tsx` adds
`ContextMenuSub` navigation/more-tools submenus plus checkbox and radio-group sections.
This directory's Basic already drops Sub/checkbox/radio (see its own doc comment
referencing `docs/jsx-parity.md`'s context-menu GAP entry), so Rtl keeps Basic's flat
shape instead of reintroducing that composition here.

## dialog
Confirm-delete-account dialog (title/description/cancel/continue), translated to
Arabic by hand. **Deviation**: shadcn's own `dialog-rtl.tsx` is an edit-profile form
using `Field`/`Input`/`Label`, which this directory's Basic doesn't compose (Basic is
purely the confirm-dialog shape); kept Basic's shape translated rather than
introducing a form dialog example for the first time here.

## dropdown-menu
Flat menu (My Account label, Profile/Billing/Settings items, shortcut), using
`dropdown-menu-rtl.tsx`'s own account/profile/billing/settings/logout vocabulary.
**Deviation**: shadcn's own `dropdown-menu-rtl.tsx` nests `DropdownMenuSub` submenus
and checkbox/radio-group sections. This directory's Basic is the flat/no-sub shape (see
`dropdown/basic.gsx`; the sub/checkbox/radio variants live in sibling example files,
not Basic), so Rtl mirrors Basic's flat shape instead.

## empty
"No projects yet" empty state (icon media, title, description, one action button),
using shadcn `empty-rtl.tsx`'s own Arabic project-themed copy. **Deviation**: shadcn's
own `empty-rtl.tsx` swaps in a folder-code icon, a second outline button, and a
"learn more" link with a rotated trailing arrow icon; this directory's Basic only
composes one icon (`icon.Inbox`) and one `EmptyContent` button, so Rtl keeps that
shape (reusing `icon.Inbox`, dropping the extra button and link) rather than
introducing new composition.

## field
Profile form (name field, "or" separator, horizontal field, disabled responsive field,
bio textarea) — Basic's own shape, translated to Arabic by hand. **Deviation**:
shadcn's own `field-rtl.tsx` is an entirely different payment-checkout form using
`Select` and `Checkbox`, neither of which this directory's Basic composes (Basic has
no Select/Checkbox import); rather than invent that composition, Rtl keeps Basic's own
profile-form structure and translates its English strings ("Name" → "الاسم", "Bio" →
"نبذة", etc.) since the shadcn source content doesn't map onto Basic's shape.

## hover-card
`@nextjs` link trigger + avatar/name/description/joined-date card, Basic's own shape,
translated to Arabic by hand (kept the `@nextjs` handle and `@vercel` mention
untranslated, as handles/mentions). **Deviation**: shadcn's own `hover-card-rtl.tsx` is
a product-price card demoing four physical-side and two logical-side trigger
placements; this directory's Basic is a single profile-card trigger, so Rtl keeps that
shape rather than introducing the multi-placement grid.

## input
A single `Input type="email"`, wrapped in `dir="rtl"`, placeholder translated by hand
("انت@مثال.com"). **Deviation**: shadcn's own `input-rtl.tsx` wraps the input in a
`Field`/`FieldLabel`/`FieldDescription` API-key form; this directory's Basic is the
bare single `Input` with no Field wrapper, so Rtl keeps that minimal shape.

## input-group
Search group (leading icon), email group (trailing send button), disabled group
(block-start/block-end addons), the variant/size button row, and the invalid-value
group — Basic's own shape, translated to Arabic by hand. **Deviation**: shadcn's own
`input-group-rtl.tsx` uses `Spinner` and `InputGroupTextarea` addons this directory's
Basic doesn't compose; kept Basic's own five-block shape instead.

## input-otp
6-digit/2-group OTP layout unchanged from Basic, wrapped in an Arabic `dir="rtl"` label
("رمز التحقق") and surrounding div — **the OTP widget itself is deliberately NOT forced
into RTL**: `ui/input-otp.gsx` pins `dir="ltr"` on both `InputOTP`'s own root and
`InputOTPGroup` by design (its own ADAPT doc comments explain this mirrors shadcn's own
RTL behavior for OTP digit entry — digits always read left-to-right regardless of
document direction). Confirmed via `grep -rn 'dir="ltr"' ui/input-otp.gsx` before
writing this file. No API gap otherwise.

## Verification
- `make generate` — clean (486 up to date plus the 12 new `rtl.x.go` files generated).
- `make highlight` — clean (161 blocks, up from 149 in batch 1).
- `go tool gsx fmt -l` on all 12 new `.gsx` files — empty output, no reformatting
  needed.
- `go build ./...` — clean.
- `go test ./...` — all green, no test updates needed this batch (unlike batch 1,
  `site/pages/pages_test.go` required no changes — none of these 12 components'
  existing TOC/count assertions needed touching).
