# Compact Preset Transport Design

**Date:** 2026-07-29

## Goal

Make ordinary gsxui preset codes short enough to paste into commands and URLs
while preserving exact, self-contained transport for custom themes and decoding
every existing `gsxui:v1:` code.

## Decision

Use a hybrid transport selected from the preset value:

- An exact built-in catalogue preset encodes as `gsxui:p1:<base62>`.
- Any preset with a custom theme value or custom radius continues to encode as
  the existing canonical `gsxui:v1:<base64url-json>` form.
- Both Go and browser decoders accept both forms.
- Existing full codes remain canonical full codes. They decode successfully but
  are not rewritten or rejected merely because the same value now has a compact
  representation.

This keeps custom imports lossless and offline while making the common path
compact. It does not introduce server-side preset storage or network lookup.

## Compact Schema

`p1` packs four catalogue indexes into one unsigned integer:

| Field | Bits | Capacity | Initial values |
| --- | ---: | ---: | --- |
| style | 4 | 16 | `nova`, `maia` |
| base color | 4 | 16 | `neutral`, `stone`, `zinc`, `mauve`, `olive`, `mist`, `taupe` |
| theme | 5 | 32 | the seven base themes followed by the 17 accent themes |
| radius | 3 | 8 | `none`, `small`, `medium`, `large` |

Fields are packed least-significant first in the table order and encoded with
the base62 alphabet
`0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`.

The arrays above are transport ABI, not views over incidental catalogue order:

- existing values never reorder or disappear;
- new values append only;
- new fields require a new transport version;
- unused indexes are invalid;
- a decoded base-color/theme pair must also be a valid catalogue combination;
- leading-zero and otherwise non-canonical base62 spellings are rejected.

The initial payload uses 16 bits and therefore at most three base62 characters.
A full compact code is at most 12 characters including `gsxui:p1:`.

## Encoding

Go remains the canonical preset authority:

1. Validate the preset.
2. Call `MatchPalette`.
3. If base color, theme, and radius all match catalogue choices, look up the
   style and selection in the immutable transport arrays and emit `p1`.
4. Otherwise emit the existing canonical JSON and `v1` base64url payload.

Browser encoding performs the same exact-value match using the schema catalogue
and emits the same bytes. It must not trust UI selection state alone: imported
JSON/CSS can make selection metadata stale unless it is derived from the preset.

Hover previews never affect either transport. Only committed resolved state is
encoded.

## Decoding and Compatibility

`DecodeShare` dispatches on the transport segment:

- `p1`: strict base62 decode, index validation, catalogue resolution, then
  canonical re-encoding to reject alternate spellings.
- `v1`: the existing strict UTF-8, base64url, JSON validation, and canonical
  JSON checks remain unchanged.
- any other segment: reject as an unsupported transport version.

The browser schema exposes both transport prefixes and the ordered compact
catalogue values. Go and browser golden tests pin identical compact codes for
every catalogue combination.

The `InputResolver`, CLI `init --preset`, CLI `apply --preset`, raw files, stdin,
and HTTPS inputs keep their current interface. They gain `p1` through
`DecodeShare`; no command syntax changes.

## Editor UX

The share code, share URL, and generated CLI commands use compact transport
automatically for catalogue presets. Custom values retain the long full code,
making the length itself an honest signal that the payload is self-contained.

Preset JSON export remains the complete canonical document in both cases. The
transport optimization never changes installed `gsxui.preset.json`.

## Errors

Errors identify the failing transport layer without exposing partial state:

- empty payload;
- invalid base62 character;
- integer overflow;
- unused catalogue index;
- invalid base-color/theme combination;
- non-canonical compact spelling;
- unsupported transport version.

Decoding is side-effect free. CLI artifact planning begins only after successful
decode and validation.

## Tests

Go tests pin:

- exact compact codes for defaults and non-default boundary combinations;
- every catalogue combination round-trips through `p1`;
- custom theme and custom radius values remain `v1`;
- all existing `v1` tests remain green;
- invalid and non-canonical compact codes are rejected;
- input resolution accepts compact codes directly, from stdin, files, and HTTPS.

Browser tests pin:

- Go and JavaScript emit identical compact codes;
- built-in editor choices produce compact commands and URLs;
- custom JSON and CSS imports produce full codes and round-trip exactly;
- old full URLs still load;
- previews never alter committed transport.

## Rejected Alternatives

**Replace full transport with compact indexes.** Rejected because arbitrary
custom tokens could not round-trip.

**Compress every JSON document.** Rejected as the common path would remain much
larger than a catalogue code and browser compression would make the currently
synchronous state API asynchronous.

**Store presets server-side and share an ID.** Rejected because codes would
depend on service availability, retention, and trust. Both chosen transports
remain offline and self-contained within their stated semantics.
