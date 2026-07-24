#!/usr/bin/env python3
"""Splice actual rendered class strings into pinned `want` literals.

Usage: python3 scripts/repin.py 'TestFoo|TestBar' ui/foo_test.go ui/bar_test.go
Run from the repo root AFTER `go tool gsx generate`. Idempotent: with no
failing pins it replaces nothing.

WARNING: only the FIRST class="..." on each failing line gets spliced. Each
want-literal is replaced globally across ALL listed files, so byte-identical
literals cross-contaminate between tests/components — always diff before committing.
"""
import pathlib, re, subprocess, sys

run_expr, files = sys.argv[1], sys.argv[2:]
r = subprocess.run(["go", "test", "./ui", "-run", run_expr],
                   capture_output=True, text=True)
out = r.stdout + r.stderr
gots = re.findall(r"got: (<[^\n]*)", out)
wants = re.findall(r"want: (<[^\n]*)", out)
assert len(gots) == len(wants), (len(gots), len(wants))
pairs = []
for g, w in zip(gots, wants):
    gm, wm = re.search(r'class="([^"]*)"', g), re.search(r'class="([^"]*)"', w)
    if gm and wm and gm.group(1) != wm.group(1):
        pairs.append((wm.group(1), gm.group(1)))
for f in files:
    p = pathlib.Path(f)
    s = p.read_text()
    n = 0
    for wc, gc in pairs:
        if wc in s:
            s = s.replace(wc, gc)
            n += 1
    p.write_text(s)
    print(f, "replaced", n)
