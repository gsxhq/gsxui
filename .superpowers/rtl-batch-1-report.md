# RTL examples — batch 1 (accordion..checkbox)

One RTL example added per component, registered as `Example{Name: "rtl", Title: "RTL"}`,
each wrapped in `<div dir="rtl" lang="ar">` (or on the outermost native element where
that's more natural), content translated to Arabic mirroring shadcn's own `*-rtl.tsx`
demos in `apps/v4/examples/base/`.

## accordion
Three FAQ question/answer pairs (password reset, subscription plan, payment methods),
first item open — same shape and same Arabic strings as shadcn's `accordion-rtl.tsx`.
No API gap.

## alert
Two alerts: payment-success (CircleCheck icon) and new-feature (Info icon), same
Arabic strings as shadcn's `alert-rtl.tsx`. No API gap.

## alert-dialog
Destructive-confirm flow (title/description/cancel/continue), same Arabic strings as
shadcn's first dialog. **Deviation**: shadcn's `alert-dialog-rtl.tsx` also renders a
second, smaller dialog using an `AlertDialogMedia` icon slot and a `size="sm"` content
variant. gsxui's `AlertDialogContent` (ui/alert-dialog.gsx) has neither a size prop nor
a media slot, so that second dialog was dropped rather than inventing new API surface;
kept the single dialog shape Basic already establishes.

## aspect-ratio
Figure/figcaption with a 16/9 box (bordered muted placeholder, matching Basic's own
substitution for a real `<Image>`) and an Arabic caption ("منظر طبيعي جميل"). No API gap
beyond what Basic already carries.

## avatar
A loaded avatar plus an overlapping group ending in a "+٣" count tile (Arabic-Indic
numerals, matching shadcn's `ar` translation). **Deviation**: shadcn's `avatar-rtl.tsx`
uses `AvatarBadge`, `AvatarGroup`, and `AvatarGroupCount` components; gsxui's avatar.gsx
only exports `Avatar`/`AvatarImage`/`AvatarFallback`. The status-badge avatar was
dropped, and the group's count indicator was approximated with a fourth `Avatar` whose
fallback text is "+٣" (following the same pattern avatar/group.gsx already uses for
overlapping avatars).

## badge
Default/secondary/destructive/outline badges plus an icon-leading ("verified" +
BadgeCheck) and icon-trailing ("bookmark" + Bookmark) badge, same Arabic strings as
shadcn's `badge-rtl.tsx`. No API gap.

## breadcrumb
Home link → ellipsis → Components link → current page, translated to Arabic.
**Deviation**: shadcn's `breadcrumb-rtl.tsx` replaces the ellipsis segment with a full
`DropdownMenu` (including RTL-aware `align`). gsxui's breadcrumb example set only
demonstrates the static `BreadcrumbEllipsis` (see Basic) — no example composes
dropdown-menu into breadcrumb — so this example keeps that established shape rather
than introducing that composition for the first time here.

## button
Outline / destructive("حذف") / outline-with-trailing-arrow (`rtl:rotate-180`) /
icon-only / disabled-secondary-with-spinner, same Arabic strings as shadcn's
`button-rtl.tsx`. No API gap.

## button-group
Icon-only back-arrow group (`rtl:rotate-180`) + an outline archive/report pair + a
quantity-stepper group (borrowed from Basic, kept in English since shadcn's own
button-group-rtl demo doesn't cover it). **Deviation**: shadcn's `button-group-rtl.tsx`
nests a `DropdownMenu` — including a radio-group submenu — inside a `ButtonGroup`.
gsxui's button-group example set doesn't compose dropdown-menu into a button group
anywhere (see Basic), so this example keeps the two already-established group shapes
(icon pair, quantity stepper) instead of introducing that composition here.

## card
Login card: title/description, a `CardAction` sign-up link, email/password form
(Input + Label, "forgot password?" link), stacked footer buttons (login, "login with
Google") — same Arabic strings as shadcn's `card-rtl.tsx`. No API gap.

## carousel
Five slides labeled with Arabic-Indic numerals (١..٥), wrapped in `dir="rtl"`.
**Deviation**: shadcn's `carousel-rtl.tsx` passes `dir`/`opts.direction` straight into
its embla-backed `Carousel`; gsxui's `ui.Carousel` (ui/carousel.gsx) has no direction
option — the same wrapper-`div[dir=rtl]` approach Basic and pagination/rtl.gsx already
rely on for logical-property flip is used instead, with no client-side embla direction
to set.

## checkbox
Four Field-composed checkboxes: plain, checked-with-description, disabled, and
title+description-wrapped-in-label — same Arabic strings as shadcn's
`checkbox-rtl.tsx`, using gsxui's `Field`/`FieldContent`/`FieldLabel`/
`FieldDescription`/`FieldTitle`/`FieldGroup` (checkbox/states.gsx's sibling API). No
API gap. Note: `data-disabled` was written as `data-disabled="true"` rather than a bare
boolean attribute — the style-contract test enforces enumerated values for that axis on
the `field` slot, matching the convention already used in field/basic.gsx and
dropdown/menubar/contextmenu examples.

## Verification
- `make generate`, `make highlight` (149 blocks now, was 137) — both clean.
- `go tool gsx fmt -w` run twice per new file, then `-l` confirmed empty (idempotent).
- `go build ./...` clean.
- `go test ./...` all green after two fixes:
  - `site/examples/checkbox/rtl.gsx`: `data-disabled` → `data-disabled="true"` (style
    contract requires declared values, not a bare boolean).
  - `site/pages/pages_test.go`: added the `example-rtl`/"RTL" TOC entry expected for
    `/components/button` now that it has three examples.
