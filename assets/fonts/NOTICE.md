# Font licences

Six variable WOFF2 fonts back the theme editor's font pair picker
(`internal/preset/catalog.go`'s `fontCatalog`) and its `@font-face`
declarations in `fonts.css`. All six are self-hosted here — no npm package,
no Google Fonts/CDN request, no build step — so both gsxui's own theme
preview and a project that vendors this directory get the same fonts with
zero extra dependencies.

**Every font below is licensed under the SIL Open Font License, Version
1.1 (OFL-1.1)**, a free, redistributable licence that explicitly permits
bundling, embedding, and redistributing the font files (including inside a
larger software project like gsxui) at no cost, with attribution. The
full licence text is reproduced once at the bottom of this file — it is
identical for all six fonts; only the copyright line differs, listed per
font below.

| Font | Catalog name | Copyright | Source |
|---|---|---|---|
| Inter Variable | `inter` | Copyright 2016 The Inter Project Authors | https://github.com/rsms/inter |
| Geist Variable | `geist` | Copyright 2024 The Geist Project Authors | https://github.com/vercel/geist-font |
| Figtree Variable | `figtree` | Copyright 2022 The Figtree Project Authors | https://github.com/erikdkennedy/figtree |
| JetBrains Mono Variable | `jetbrains-mono` | Copyright 2020 The JetBrains Mono Project Authors | https://github.com/JetBrains/JetBrainsMono |
| Noto Sans Variable | `noto-sans` | Copyright 2022 The Noto Project Authors | https://github.com/notofonts/latin-greek-cyrillic |
| Playfair Display Variable | `playfair-display` | Copyright 2017 The Playfair Display Project Authors, with Reserved Font Name "Playfair Display" | https://github.com/clauseggers/Playfair-Display |

## Provenance and subsetting

Each `.woff2` file here is the upstream OFL-licensed variable font's own
**"latin" Unicode-range subset, normal style only** (no latin-ext,
cyrillic, greek, vietnamese, etc.; no italic), re-hosted as a static binary
asset rather than pulled from an npm package or CDN at build or run time.
The subset files and this licence text were sourced from the published
`@fontsource-variable/*` npm packages (themselves OFL-licensed
redistributions of each font's own upstream release) purely as a
convenient, already-subsetted source of the binary font data — gsxui does
not depend on those packages at runtime or build time; the `.woff2` files
are checked into this repository like any other generated/vendored asset
(the same pattern `assets/css/animate.css` already uses for tw-animate-css,
per that file's own header comment).

Each font is a single variable-weight file (`font-weight: 1 999` in
`fonts.css`), so one file covers every weight the picker's chosen font can
render — no per-weight duplication.

## Deliberately not included

- **Italic styles** — the picker only ever sets `--font-sans`/
  `--font-heading` to a normal-weight family name; italic text still
  renders (browsers synthesize an oblique when no italic face is
  registered), just not with each font's true italic hinting. Kept out to
  hold the byte budget down.
- **Non-latin subsets** (latin-ext, cyrillic, cyrillic-ext, greek,
  vietnamese, devanagari, etc.) — same reasoning. A future pass can add
  them if non-Latin-script coverage becomes a real requirement; nothing
  about `fonts.css`'s `@font-face` shape needs to change to add more
  `unicode-range`-scoped `src` rules per family later.

## Byte budget

All six `.woff2` files together are about 208 KiB (`figtree` 20,156 B,
`geist` 29,400 B, `inter` 48,256 B, `jetbrains-mono` 40,404 B, `noto-sans`
35,820 B, `playfair-display` 38,404 B — 212,440 B total). Not embedded in
any Go binary; these are static files served/copied alongside `assets/css`.

## Full licence text (OFL-1.1)

```
-----------------------------------------------------------
SIL OPEN FONT LICENSE Version 1.1 - 26 February 2007
-----------------------------------------------------------

PREAMBLE
The goals of the Open Font License (OFL) are to stimulate worldwide
development of collaborative font projects, to support the font creation
efforts of academic and linguistic communities, and to provide a free and
open framework in which fonts may be shared and improved in partnership
with others.

The OFL allows the licensed fonts to be used, studied, modified and
redistributed freely as long as they are not sold by themselves. The
fonts, including any derivative works, can be bundled, embedded,
redistributed and/or sold with any software provided that any reserved
names are not used by derivative works. The fonts and derivatives,
however, cannot be released under any other type of license. The
requirement for fonts to remain under this license does not apply
to any document created using the fonts or their derivatives.

DEFINITIONS
"Font Software" refers to the set of files released by the Copyright
Holder(s) under this license and clearly marked as such. This may
include source files, build scripts and documentation.

"Reserved Font Name" refers to any names specified as such after the
copyright statement(s).

"Original Version" refers to the collection of Font Software components as
distributed by the Copyright Holder(s).

"Modified Version" refers to any derivative made by adding to, deleting,
or substituting -- in part or in whole -- any of the components of the
Original Version, by changing formats or by porting the Font Software to a
new environment.

"Author" refers to any designer, engineer, programmer, technical
writer or other person who contributed to the Font Software.

PERMISSION & CONDITIONS
Permission is hereby granted, free of charge, to any person obtaining
a copy of the Font Software, to use, study, copy, merge, embed, modify,
redistribute, and sell modified and unmodified copies of the Font
Software, subject to the following conditions:

1) Neither the Font Software nor any of its individual components,
in Original or Modified Versions, may be sold by itself.

2) Original or Modified Versions of the Font Software may be bundled,
redistributed and/or sold with any software, provided that each copy
contains the above copyright notice and this license. These can be
included either as stand-alone text files, human-readable headers or
in the appropriate machine-readable metadata fields within text or
binary files as long as those fields can be easily viewed by the user.

3) No Modified Version of the Font Software may use the Reserved Font
Name(s) unless explicit written permission is granted by the corresponding
Copyright Holder. This restriction only applies to the primary font name as
presented to the users.

4) The name(s) of the Copyright Holder(s) or the Author(s) of the Font
Software shall not be used to promote, endorse or advertise any
Modified Version, except to acknowledge the contribution(s) of the
Copyright Holder(s) and the Author(s) or with their explicit written
permission.

5) The Font Software, modified or unmodified, in part or in whole,
must be distributed entirely under this license, and must not be
distributed under any other license. The requirement for fonts to
remain under this license does not apply to any document created
using the Font Software.

TERMINATION
This license becomes null and void if any of the above conditions are
not met.

DISCLAIMER
THE FONT SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO ANY WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT
OF COPYRIGHT, PATENT, TRADEMARK, OR OTHER RIGHT. IN NO EVENT SHALL THE
COPYRIGHT HOLDER BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
INCLUDING ANY GENERAL, SPECIAL, INDIRECT, INCIDENTAL, OR CONSEQUENTIAL
DAMAGES, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF THE USE OR INABILITY TO USE THE FONT SOFTWARE OR FROM
OTHER DEALINGS IN THE FONT SOFTWARE.
```
