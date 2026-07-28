# Task 7 report — transactional preset apply and crash recovery

Commit: `feat: apply presets transactionally`

## Scope delivered

- Added `gsxui apply --preset <file|code|-> [--only theme|style] [--yes]
  [--overwrite]` for the standalone Button pilot.
  - Full apply replaces the durable preset, theme, an installed Button when
    style changes, and managed configuration.
  - `--only theme` preserves style and Button source.
  - `--only style` preserves theme/radius and replaces an installed Button.
  - Installed styled components come from the registry's exact file/directory
    shape plus exact configured filesystem targets, not managed hashes.
  - Maia refuses every installed styled component other than standalone Button
    before prompting or mutation. Nova preserves non-Button components.
  - Modified managed targets and pre-managed/unmanaged Button source require
    `--overwrite`; apply never performs a three-way merge.
  - The summary pins changed axes, exact files, source replacement, conflicts,
    and the commit/stash recommendation. `--yes` skips only confirmation.
- Routed init, add, and apply source/config writes through one journaled
  transaction. Config remains the final artifact and never hashes itself.
- Added next-invocation recovery before init, mutating add, and apply planning.
  Read-only `add --diff` remains non-mutating.
- Added exact Nova → Maia → Nova Button e2e coverage. Each style is compared
  with its generated registry source and built before later Dialog work can
  mask the result.

## Transaction and recovery model

Every changed artifact is fully planned before the first project mutation.
The journal records exact target/staged/backup paths, previous
existence/hash/mode, intended hash/mode, created-directory ownership, phase,
cursor, and rollback source phase/cursor. Staged files are explicitly chmodded
before fsync so restrictive umasks cannot falsify the recorded mode.

Missing project directories are prebuilt under
`.gsxui-transaction-<nonce>/directories/<index>` with an exact nonce ownership
marker, then atomically renamed into the project. Recovery validates canonical
parent-first ordering, portable uniqueness, structural target ancestry, and
every published nonce marker before changing a target. Omitted, substituted,
or unowned directories therefore fail closed without mutation.

Forward checkpoints are:

```text
prepared
directory-create-pending i
directory-created i+1
backup-pending i
backed-up i
replace-pending i
replaced i+1
committed n
generation-pending n
finalizing n
```

Rollback checkpoints are:

```text
rollback-pending n ... 0
restored-generation-pending 0
rollback-directories-pending d ... 0
cleanup
```

Rollback validates the complete physical phase before its first mutation and
rechecks each target immediately before removing or restoring it. A reverse
cursor makes replacement removal and backup restoration idempotently
resumable. Transaction-created directories are atomically moved back under the
transaction root in reverse depth order. Cleanup revalidates each moved
directory contains only its nonce marker and uses non-recursive removal, so an
unowned concurrent file fails closed rather than being deleted.

A generation interruption leaves `generation-pending`. Recovery restores
source/config, durably publishes `restored-generation-pending`, and invokes GSX
once before any new planning. An interruption of restored-source generation
retries from that phase. Successful generation moves to resumable finalization;
the journal is retained until committed cleanup is durable.

## TDD evidence

### Initial transaction, init/add routing, and apply

The first transaction tests failed at compile time because the transaction API
did not exist. After the initial journal implementation, the injected matrix
covered 26 rename boundaries and 18 persisted checkpoints for
preset/theme/Button/config.

Init/add routing tests then failed because generation ran only once on failure,
replaced bytes remained, and init did not run GSX. Apply behavior tests failed
with `unknown command "apply"`. They became green only after init/add used the
transaction and the declared apply planner/command existed.

### Adversarial hardening RED

The independent review's probes observed all of these failures before their
focused fixes:

```text
recovery generation calls = 0, want 1
recover error = <nil>, want post-interruption edit refusal
terminal transaction target = "old preset\n", want committed bytes
DATA LOSS: recovery deleted user directory
recovery could not resume interrupted rollback:
generation-pending transaction target target.txt changed after interruption
existing target mode = 600, want 644
```

The corresponding authored tests now pin:

- post-interruption edits to existing and newly created targets;
- phase-inconsistent, portable-aliased, reserved, and symlinked journals;
- intended/previous hashes and modes plus exact staged/backup slots;
- resumable terminal finalization after directory-sync failure;
- recovery-time GSX regeneration;
- unowned, substituted, and omitted owned directories;
- restrictive umask behavior for existing and absent targets;
- every directory-create phase and the atomic publish-before-checkpoint window;
- replacement-removal, backup-restoration, restored-phase, directory-removal,
  and cleanup interruptions; and
- every rollback entry/directory cursor.

The final generation-failure matrix pins 22 exact checkpoints across two
entries and two new directories. The final independent disposable-copy review
also reran the original six findings and the later directory/rollback/mode
probes. Its verdict was PASS with no Critical or Important findings.

## Verification

A temporary `go.work` selected this worktree and exact GSX core
`ef72f5eba066d7e87adf7dcadc2db62d00f22efe`.

```text
$ go test ./internal/cli \
    -run 'Test(Transaction|Recover|Apply|Init|Add|E2E)' -count=1
ok github.com/gsxhq/gsxui/internal/cli 101.609s

$ go test ./internal/cli -race -count=1
ok github.com/gsxhq/gsxui/internal/cli 102.106s

$ gopls check -severity=hint \
    internal/cli/transaction.go internal/cli/transaction_test.go \
    internal/cli/apply.go internal/cli/apply_test.go \
    internal/cli/init.go internal/cli/init_test.go \
    internal/cli/add.go internal/cli/add_test.go \
    internal/cli/e2e_test.go internal/cli/run.go
# no output; exit 0

$ make ci
# exit 0
# 182 generated files up to date
# all Go tests and vet pass
# 292 Playwright tests pass
# gofmt clean

$ git diff --check
# no output; exit 0
```

The required wildcard `gopls check -severity=hint internal/cli/*.go` exits 0
and reports one pre-existing `bufio.Scanner` final-error hint in unchanged
`internal/cli/module.go`; every Task 7 file is clean.

## External side effects and concerns

Init validates and plans all artifact bytes before running its three dependency
commands. The commands run in this explicit order before the filesystem
transaction:

```text
go get github.com/gsxhq/gsx@latest
go get github.com/jackielii/tailwind-merge-go@latest
go get -tool github.com/gsxhq/gsx/cmd/gsx@latest
```

Those Go module/network side effects are intentionally not claimed as
transactional because the CLI cannot faithfully restore arbitrary external
effects. All gsxui-owned source/config changes and GSX generation are covered
by the transaction and recovery journal.

No reusable GSX or gsxui improvement beyond the approved transaction design
was exposed, so dogfood scope was not expanded.
