# Third-party attribution

The components under `ui/` are ports of **shadcn/ui**
([shadcn-ui/ui](https://github.com/shadcn-ui/ui), MIT © shadcn) from React/JSX
to gsx: part names, `data-slot` contracts, and Tailwind class strings follow
the shadcn v4 "new-york" registry. The event-delegation JS architecture is
original to gsxui.

Icon markup embedded in components is from [Lucide](https://lucide.dev)
(ISC License, © Lucide Contributors — https://lucide.dev/license). `ui/icon`'s
generated files (`icon_data.go`, `icon_defs.go`) are produced by
`internal/lucidegen` from
[lucide-icons/lucide@4e2efb9350fac7dbcb94caae9ccac5654bdcc050](https://github.com/lucide-icons/lucide/tree/4e2efb9350fac7dbcb94caae9ccac5654bdcc050)
(1,748 icons); regenerate with `make icons` against a fresh checkout to pick
up a newer Lucide release.

gsxui's server-side shadcn ports (component structure, `data-slot` parity,
the generated-registry CLI) were developed with
[templUI](https://github.com/templui/templui) (MIT, © templUI contributors)
consulted as architectural precedent for porting shadcn/ui to a Go templating
language — structure only; no templUI code was copied.

When `gsxui add` vendors component source into your project, the vendored
files remain subject to these notices alongside gsxui's own MIT license.

---

## shadcn/ui license (MIT)

MIT License

Copyright (c) 2023 shadcn

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
