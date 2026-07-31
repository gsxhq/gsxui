package cli

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
)

func TestTransactionRenameFailureMatrixRestoresExactTree(t *testing.T) {
	root, artifacts := transactionFixture(t)
	originalFS := transactionFS
	var renameCount int
	transactionFS.rename = func(oldPath, newPath string) error {
		renameCount++
		return os.Rename(oldPath, newPath)
	}
	if err := executeArtifactTransaction(root, artifacts, nil, nil); err != nil {
		t.Fatalf("count transaction renames: %v", err)
	}
	transactionFS = originalFS
	if renameCount == 0 {
		t.Fatal("transaction performed no rename durability boundaries")
	}

	for failIndex := range renameCount {
		t.Run(fmt.Sprintf("rename_%02d", failIndex), func(t *testing.T) {
			root, artifacts := transactionFixture(t)
			before := treeDigest(t, root)
			originalFS := transactionFS
			t.Cleanup(func() { transactionFS = originalFS })
			var index int
			transactionFS.rename = func(oldPath, newPath string) error {
				current := index
				index++
				if current == failIndex {
					return fmt.Errorf("injected rename failure %d", failIndex)
				}
				return os.Rename(oldPath, newPath)
			}

			err := executeArtifactTransaction(root, artifacts, nil, nil)
			if err == nil || !strings.Contains(err.Error(), "injected rename failure") {
				t.Fatalf("transaction error = %v, want injected rename failure", err)
			}
			if got := treeDigest(t, root); got != before {
				t.Fatalf(
					"rename failure %d did not restore exact tree:\n before %s\n  after %s",
					failIndex,
					before,
					got,
				)
			}
			assertNoTransactionArtifacts(t, root)
		})
	}
}

func TestTransactionJournalPublishSyncFailureRestoresExactTree(t *testing.T) {
	root, artifacts := transactionFixture(t)
	before := treeDigest(t, root)
	originalFS := transactionFS
	t.Cleanup(func() { transactionFS = originalFS })
	injected := false
	resolvedRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	transactionFS.syncDirectory = func(directory string) error {
		if !injected && filepath.Clean(directory) == filepath.Clean(resolvedRoot) {
			if _, err := os.Stat(filepath.Join(root, transactionJournalName)); err == nil {
				injected = true
				return errors.New("injected journal directory sync failure")
			}
		}
		return syncOSDirectory(directory)
	}

	err = executeArtifactTransaction(root, artifacts, nil, nil)
	if err == nil || !strings.Contains(err.Error(), "injected journal directory sync failure") {
		t.Fatalf("transaction error = %v, want injected journal sync failure", err)
	}
	if got := treeDigest(t, root); got != before {
		t.Fatalf("journal sync failure did not restore exact tree:\n before %s\n  after %s", before, got)
	}
	assertNoTransactionArtifacts(t, root)
}

func TestRecoverTransactionRollsBackEveryPersistedPhase(t *testing.T) {
	root, artifacts := transactionFixture(t)
	originalHook := transactionCheckpointHook
	var checkpoints []transactionCheckpoint
	transactionCheckpointHook = func(checkpoint transactionCheckpoint) {
		checkpoints = append(checkpoints, checkpoint)
	}
	if err := executeArtifactTransaction(root, artifacts, nil, nil); err != nil {
		t.Fatalf("collect transaction checkpoints: %v", err)
	}
	transactionCheckpointHook = originalHook

	want := []transactionCheckpoint{{phase: transactionPhasePrepared, cursor: 0}}
	for index := range artifacts {
		want = append(want,
			transactionCheckpoint{phase: transactionPhaseBackupPending, cursor: index},
			transactionCheckpoint{phase: transactionPhaseBackedUp, cursor: index},
			transactionCheckpoint{phase: transactionPhaseReplacePending, cursor: index},
			transactionCheckpoint{phase: transactionPhaseReplaced, cursor: index + 1},
		)
	}
	want = append(want, transactionCheckpoint{
		phase:  transactionPhaseCommitted,
		cursor: len(artifacts),
	})
	want = append(want, transactionCheckpoint{
		phase:  transactionPhaseFinalizing,
		cursor: len(artifacts),
	})
	if !reflect.DeepEqual(checkpoints, want) {
		t.Fatalf("transaction checkpoints:\n got: %#v\nwant: %#v", checkpoints, want)
	}
	committedRoot, committedArtifacts := transactionFixture(t)
	if err := executeArtifactTransaction(committedRoot, committedArtifacts, nil, nil); err != nil {
		t.Fatalf("prepare committed tree expectation: %v", err)
	}
	committed := treeDigest(t, committedRoot)

	for failIndex, checkpoint := range checkpoints {
		t.Run(fmt.Sprintf("%02d_%s_%d", failIndex, checkpoint.phase, checkpoint.cursor), func(t *testing.T) {
			root, artifacts := transactionFixture(t)
			before := treeDigest(t, root)
			originalHook := transactionCheckpointHook
			t.Cleanup(func() { transactionCheckpointHook = originalHook })
			var index int
			transactionCheckpointHook = func(transactionCheckpoint) {
				current := index
				index++
				if current == failIndex {
					panic("injected process interruption")
				}
			}

			func() {
				defer func() {
					if recovered := recover(); recovered != "injected process interruption" {
						t.Fatalf("transaction panic = %v, want injected interruption", recovered)
					}
				}()
				_ = executeArtifactTransaction(root, artifacts, nil, nil)
			}()
			transactionCheckpointHook = originalHook

			if err := recoverArtifactTransaction(root); err != nil {
				t.Fatalf("recover transaction at %s/%d: %v", checkpoint.phase, checkpoint.cursor, err)
			}
			wantTree := before
			if checkpoint.phase == transactionPhaseFinalizing {
				wantTree = committed
			}
			if got := treeDigest(t, root); got != wantTree {
				t.Fatalf(
					"recovery at %s/%d produced the wrong exact tree:\n want %s\n  got %s",
					checkpoint.phase,
					checkpoint.cursor,
					wantTree,
					got,
				)
			}
			assertNoTransactionArtifacts(t, root)
		})
	}
}

func TestTransactionPostCommitFailureRestoresAndReportsRecoveryFailure(t *testing.T) {
	root, artifacts := transactionFixture(t)
	before := treeDigest(t, root)
	postCommitErr := errors.New("injected generation failure")
	postRollbackErr := errors.New("injected restoration generation failure")

	err := executeArtifactTransaction(
		root,
		artifacts,
		func() error { return postCommitErr },
		func() error { return postRollbackErr },
	)
	if !errors.Is(err, postCommitErr) || !errors.Is(err, postRollbackErr) {
		t.Fatalf("transaction error = %v, want both generation failures", err)
	}
	if got := treeDigest(t, root); got != before {
		t.Fatalf("post-commit failure did not restore exact tree:\n before %s\n  after %s", before, got)
	}
	assertNoTransactionArtifacts(t, root)
}

func TestRecoverTransactionRejectsSymlinkedJournal(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "journal.json")
	if err := os.WriteFile(outside, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, transactionJournalName)); err != nil {
		t.Fatal(err)
	}
	before := treeDigest(t, root)

	err := recoverArtifactTransaction(root)
	if err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("recover error = %v, want symlink rejection", err)
	}
	if got := treeDigest(t, root); got != before {
		t.Fatalf("symlinked-journal rejection changed tree:\n before %s\n  after %s", before, got)
	}
}

func TestTransactionRejectsPortableReservedNamespaceAncestorsBeforeMutation(t *testing.T) {
	for _, relative := range []string{
		".GSXUI-TRANSACTION.JSON/target",
		".Gsxui-Transaction-Attack/target",
	} {
		t.Run(relative, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "keep.txt"), []byte("keep\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			before := treeDigest(t, root)

			err := executeArtifactTransaction(root, []artifact{{
				RelativePath: relative,
				Content:      []byte("unexpected\n"),
				Managed:      true,
			}}, nil, nil)
			if err == nil || (!strings.Contains(err.Error(), "portable") &&
				!strings.Contains(err.Error(), "reserved")) {
				t.Fatalf("transaction error = %v, want portable reserved-namespace rejection", err)
			}
			if got := treeDigest(t, root); got != before {
				t.Fatalf("reserved namespace rejection changed tree:\n before %s\n  after %s", before, got)
			}
			assertNoTransactionArtifacts(t, root)
		})
	}
}

func TestRecoverTransactionRejectsPortableReservedJournalTargetWithoutMutation(t *testing.T) {
	root := t.TempDir()
	transactionDirectory := filepath.Join(root, ".gsxui-transaction-attack")
	if err := os.MkdirAll(filepath.Join(transactionDirectory, "staged"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(transactionDirectory, "backup"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(transactionDirectory, "staged", "000000"),
		[]byte("malicious staged bytes\n"),
		0o600,
	); err != nil {
		t.Fatal(err)
	}
	journal := []byte(`{
  "version": 1,
  "directory": ".gsxui-transaction-attack",
  "phase": "prepared",
  "cursor": 0,
  "entries": [
    {
      "target": ".GSXUI-TRANSACTION.JSON/target",
      "staged": ".gsxui-transaction-attack/staged/000000",
      "backup": ".gsxui-transaction-attack/backup/000000",
      "previousExists": false
    }
  ]
}
`)
	if err := os.WriteFile(filepath.Join(root, transactionJournalName), journal, 0o600); err != nil {
		t.Fatal(err)
	}
	before := treeDigest(t, root)

	err := recoverArtifactTransaction(root)
	if err == nil || !strings.Contains(err.Error(), "portable") {
		t.Fatalf("recover error = %v, want portable reserved-target rejection", err)
	}
	if got := treeDigest(t, root); got != before {
		t.Fatalf("invalid recovery journal changed tree:\n before %s\n  after %s", before, got)
	}
}

func TestRecoverTransactionRefusesPostInterruptionEditToExistingTarget(t *testing.T) {
	root, artifacts := transactionFixture(t)
	interruptTransactionAtCheckpoint(t, root, artifacts, transactionCheckpoint{
		phase:  transactionPhaseReplaced,
		cursor: 1,
	})
	target := filepath.Join(root, "gsxui.preset.json")
	if err := os.WriteFile(target, []byte("post-interruption user edit\n"), 0o640); err != nil {
		t.Fatal(err)
	}
	beforeRecovery := treeDigest(t, root)

	err := recoverArtifactTransaction(root)
	if err == nil || !strings.Contains(err.Error(), "changed after interruption") {
		t.Fatalf("recover error = %v, want post-interruption edit refusal", err)
	}
	if got := treeDigest(t, root); got != beforeRecovery {
		t.Fatalf("recovery changed post-interruption edit:\n before %s\n  after %s", beforeRecovery, got)
	}
	if _, err := os.Stat(filepath.Join(root, transactionJournalName)); err != nil {
		t.Fatalf("refused recovery removed journal: %v", err)
	}
}

func TestRecoverTransactionRefusesPostInterruptionEditToNewTargetAndParents(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "keep.txt"), []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	artifacts := []artifact{{
		RelativePath: "new/nested/button.gsx",
		Content:      []byte("transaction bytes\n"),
		Managed:      true,
	}}
	interruptTransactionAtCheckpoint(t, root, artifacts, transactionCheckpoint{
		phase:  transactionPhaseReplaced,
		cursor: 1,
	})
	target := filepath.Join(root, "new", "nested", "button.gsx")
	if err := os.WriteFile(target, []byte("post-interruption user edit\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	beforeRecovery := treeDigest(t, root)

	err := recoverArtifactTransaction(root)
	if err == nil || !strings.Contains(err.Error(), "changed after interruption") {
		t.Fatalf("recover error = %v, want post-interruption edit refusal", err)
	}
	if got := treeDigest(t, root); got != beforeRecovery {
		t.Fatalf("recovery deleted post-interruption target or parents:\n before %s\n  after %s", beforeRecovery, got)
	}
}

func TestRecoverTransactionRejectsPhaseInconsistentPreparedJournal(t *testing.T) {
	root := t.TempDir()
	victim := filepath.Join(root, "victim.txt")
	const victimContent = "victim bytes\n"
	if err := os.WriteFile(victim, []byte(victimContent), 0o644); err != nil {
		t.Fatal(err)
	}
	transactionDirectory := filepath.Join(root, ".gsxui-transaction-attack")
	if err := os.MkdirAll(filepath.Join(transactionDirectory, "staged"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(transactionDirectory, "backup"), 0o700); err != nil {
		t.Fatal(err)
	}
	victimSum := sha256.Sum256([]byte(victimContent))
	journal := fmt.Sprintf(`{
  "version": 1,
  "directory": ".gsxui-transaction-attack",
  "phase": "prepared",
  "cursor": 0,
  "entries": [
    {
      "target": "victim.txt",
      "staged": ".gsxui-transaction-attack/staged/000000",
      "backup": ".gsxui-transaction-attack/backup/000000",
      "previousExists": false,
      "previousMode": 420,
      "intendedHash": "%x",
      "intendedMode": 420
    }
  ]
}
`, victimSum)
	if err := os.WriteFile(filepath.Join(root, transactionJournalName), []byte(journal), 0o600); err != nil {
		t.Fatal(err)
	}
	before := treeDigest(t, root)

	err := recoverArtifactTransaction(root)
	if err == nil || !strings.Contains(err.Error(), "prepared") {
		t.Fatalf("recover error = %v, want phase-inconsistent journal rejection", err)
	}
	if got := treeDigest(t, root); got != before {
		t.Fatalf("phase-inconsistent journal deleted victim:\n before %s\n  after %s", before, got)
	}
}

func TestRecoverTransactionRejectsPortableAliasedJournalTargets(t *testing.T) {
	root := t.TempDir()
	transactionDirectory := filepath.Join(root, ".gsxui-transaction-attack")
	for _, directory := range []string{
		filepath.Join(transactionDirectory, "staged"),
		filepath.Join(transactionDirectory, "backup"),
	} {
		if err := os.MkdirAll(directory, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	journal := []byte(`{
  "version": 1,
  "directory": ".gsxui-transaction-attack",
  "phase": "prepared",
  "cursor": 0,
  "entries": [
    {
      "target": "UI/Button.gsx",
      "staged": ".gsxui-transaction-attack/staged/000000",
      "backup": ".gsxui-transaction-attack/backup/000000",
      "previousExists": false,
      "intendedHash": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
      "intendedMode": 420
    },
    {
      "target": "ui/button.gsx",
      "staged": ".gsxui-transaction-attack/staged/000001",
      "backup": ".gsxui-transaction-attack/backup/000001",
      "previousExists": false,
      "intendedHash": "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
      "intendedMode": 420
    }
  ]
}
`)
	if err := os.WriteFile(filepath.Join(root, transactionJournalName), journal, 0o600); err != nil {
		t.Fatal(err)
	}
	before := treeDigest(t, root)

	err := recoverArtifactTransaction(root)
	if err == nil || !strings.Contains(err.Error(), "portable") {
		t.Fatalf("recover error = %v, want portable target-alias rejection", err)
	}
	if got := treeDigest(t, root); got != before {
		t.Fatalf("portable-alias journal changed tree:\n before %s\n  after %s", before, got)
	}
}

func TestTransactionFinalizationSyncFailureRemainsRecoverable(t *testing.T) {
	root, artifacts := transactionFixture(t)
	originalFS := transactionFS
	originalHook := transactionCheckpointHook
	t.Cleanup(func() {
		transactionFS = originalFS
		transactionCheckpointHook = originalHook
	})
	resolvedRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	injected := false
	finalizationStarted := false
	transactionCheckpointHook = func(checkpoint transactionCheckpoint) {
		if checkpoint.phase == transactionPhaseFinalizing {
			finalizationStarted = true
		}
	}
	transactionFS.syncDirectory = func(directory string) error {
		if finalizationStarted && !injected &&
			filepath.Clean(directory) == filepath.Clean(resolvedRoot) {
			if _, journalErr := os.Stat(filepath.Join(root, transactionJournalName)); os.IsNotExist(journalErr) {
				injected = true
				return errors.New("injected final journal removal sync failure")
			}
		}
		return syncOSDirectory(directory)
	}

	err = executeArtifactTransaction(root, artifacts, nil, nil)
	if err == nil || !strings.Contains(err.Error(), "injected final journal removal sync failure") {
		t.Fatalf("transaction error = %v, want finalization sync failure", err)
	}
	if got := string(readProjectFile(t, root, "gsxui.preset.json")); got != "new preset\n" {
		t.Fatalf("terminal transaction target = %q, want committed bytes", got)
	}
	if _, err := os.Stat(filepath.Join(root, transactionJournalName)); err != nil {
		t.Fatalf("finalization failure did not retain a recovery journal: %v", err)
	}

	transactionFS = originalFS
	if err := recoverArtifactTransaction(root); err != nil {
		t.Fatal(err)
	}
	if got := string(readProjectFile(t, root, "gsxui.preset.json")); got != "new preset\n" {
		t.Fatalf("resumed finalization rolled back committed bytes: %q", got)
	}
	assertNoTransactionArtifacts(t, root)
}

func TestRecoverTransactionRejectsUnownedCreatedDirectoryWithoutMutation(t *testing.T) {
	root := t.TempDir()
	userDirectory := filepath.Join(root, "user-empty-directory")
	if err := os.Mkdir(userDirectory, 0o755); err != nil {
		t.Fatal(err)
	}
	transactionDirectory := filepath.Join(root, ".gsxui-transaction-attack")
	for _, directory := range []string{
		filepath.Join(transactionDirectory, "staged"),
		filepath.Join(transactionDirectory, "backup"),
	} {
		if err := os.MkdirAll(directory, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	journal := []byte(`{
  "version": 1,
  "directory": ".gsxui-transaction-attack",
  "phase": "prepared",
  "cursor": 0,
  "createdDirectories": ["user-empty-directory"],
  "entries": []
}
`)
	if err := os.WriteFile(filepath.Join(root, transactionJournalName), journal, 0o600); err != nil {
		t.Fatal(err)
	}
	before := treeDigest(t, root)

	err := recoverArtifactTransaction(root)
	if err == nil || !strings.Contains(err.Error(), "created director") {
		t.Fatalf("recover error = %v, want unowned created-directory refusal", err)
	}
	if got := treeDigest(t, root); got != before {
		t.Fatalf("unowned created-directory recovery changed tree:\n before %s\n  after %s", before, got)
	}
}

func TestRecoverTransactionRejectsOmittedOwnedDirectoryWithoutMutation(t *testing.T) {
	for _, createdDirectories := range [][]string{
		{"new"},
		{"new/nested"},
	} {
		t.Run(strings.Join(createdDirectories, "_"), func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "keep.txt"), []byte("keep\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			artifacts := []artifact{{
				RelativePath: "new/nested/file.txt",
				Content:      []byte("transaction bytes\n"),
				Managed:      true,
			}}
			interruptTransactionAtCheckpoint(t, root, artifacts, transactionCheckpoint{
				phase:  transactionPhaseCommitted,
				cursor: 1,
			})
			journalPath := filepath.Join(root, transactionJournalName)
			content, err := os.ReadFile(journalPath)
			if err != nil {
				t.Fatal(err)
			}
			journal, err := parseTransactionJournal(content)
			if err != nil {
				t.Fatal(err)
			}
			journal.CreatedDirectories = createdDirectories
			content, err = json.MarshalIndent(journal, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			content = append(content, '\n')
			if err := os.WriteFile(journalPath, content, 0o600); err != nil {
				t.Fatal(err)
			}
			before := treeDigest(t, root)

			err = recoverArtifactTransaction(root)
			if err == nil || !strings.Contains(err.Error(), "ownership marker") {
				t.Fatalf("recover error = %v, want omitted ownership-marker refusal", err)
			}
			if got := treeDigest(t, root); got != before {
				t.Fatalf(
					"omitted owned-directory recovery changed tree:\n before %s\n  after %s",
					before,
					got,
				)
			}
		})
	}
}

func TestRecoverTransactionResumesCrashDuringRollback(t *testing.T) {
	root, artifacts := transactionFixture(t)
	before := treeDigest(t, root)
	originalFS := transactionFS
	originalRunner := runCommand
	t.Cleanup(func() {
		transactionFS = originalFS
		runCommand = originalRunner
	})
	transactionFS.rename = func(oldPath, newPath string) error {
		if strings.Contains(filepath.ToSlash(oldPath), "/backup/") {
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
			panic("injected rollback interruption")
		}
		return os.Rename(oldPath, newPath)
	}
	func() {
		defer func() {
			if recovered := recover(); recovered != "injected rollback interruption" {
				t.Fatalf("rollback panic = %v", recovered)
			}
		}()
		_ = executeArtifactTransaction(
			root,
			artifacts,
			func() error { return errors.New("injected generation failure") },
			func() error { return nil },
		)
	}()

	transactionFS = originalFS
	runCommand = func(string, string, ...string) error { return nil }
	if err := recoverArtifactTransaction(root); err != nil {
		t.Fatalf("resume interrupted rollback: %v", err)
	}
	if got := treeDigest(t, root); got != before {
		t.Fatalf("resumed rollback did not restore exact tree:\n before %s\n  after %s", before, got)
	}
	assertNoTransactionArtifacts(t, root)
}

func TestTransactionPreservesModesDespiteRestrictiveUmask(t *testing.T) {
	root := t.TempDir()
	existing := filepath.Join(root, "existing.txt")
	if err := os.WriteFile(existing, []byte("old\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	previousUmask := syscall.Umask(0o077)
	t.Cleanup(func() { syscall.Umask(previousUmask) })

	if err := executeArtifactTransaction(root, []artifact{
		{
			RelativePath: "existing.txt",
			Content:      []byte("new\n"),
			Managed:      true,
		},
		{
			RelativePath: "absent.txt",
			Content:      []byte("new absent\n"),
			Managed:      true,
		},
	}, nil, nil); err != nil {
		t.Fatal(err)
	}
	for _, relative := range []string{"existing.txt", "absent.txt"} {
		info, err := os.Stat(filepath.Join(root, relative))
		if err != nil {
			t.Fatal(err)
		}
		if got := info.Mode().Perm(); got != 0o644 {
			t.Fatalf("%s mode = %o, want 644", relative, got)
		}
	}
}

func TestRecoverTransactionRollsBackEveryCreatedDirectoryPhase(t *testing.T) {
	fixture := func(t *testing.T) (string, []artifact) {
		t.Helper()
		root := t.TempDir()
		if err := os.WriteFile(filepath.Join(root, "keep.txt"), []byte("keep\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		return root, []artifact{{
			RelativePath: "new/nested/button.gsx",
			Content:      []byte("button\n"),
			Managed:      true,
		}}
	}
	root, artifacts := fixture(t)
	originalHook := transactionCheckpointHook
	var checkpoints []transactionCheckpoint
	transactionCheckpointHook = func(checkpoint transactionCheckpoint) {
		checkpoints = append(checkpoints, checkpoint)
	}
	if err := executeArtifactTransaction(root, artifacts, nil, nil); err != nil {
		t.Fatal(err)
	}
	transactionCheckpointHook = originalHook
	wantDirectoryCheckpoints := []transactionCheckpoint{
		{phase: transactionPhaseDirectoryCreatePending, cursor: 0},
		{phase: transactionPhaseDirectoryCreated, cursor: 1},
		{phase: transactionPhaseDirectoryCreatePending, cursor: 1},
		{phase: transactionPhaseDirectoryCreated, cursor: 2},
	}
	var gotDirectoryCheckpoints []transactionCheckpoint
	for _, checkpoint := range checkpoints {
		if checkpoint.phase == transactionPhaseDirectoryCreatePending ||
			checkpoint.phase == transactionPhaseDirectoryCreated {
			gotDirectoryCheckpoints = append(gotDirectoryCheckpoints, checkpoint)
		}
	}
	if !reflect.DeepEqual(gotDirectoryCheckpoints, wantDirectoryCheckpoints) {
		t.Fatalf(
			"directory checkpoints:\n got: %#v\nwant: %#v",
			gotDirectoryCheckpoints,
			wantDirectoryCheckpoints,
		)
	}
	committed := treeDigest(t, root)

	for index, checkpoint := range checkpoints {
		t.Run(fmt.Sprintf("%02d_%s_%d", index, checkpoint.phase, checkpoint.cursor), func(t *testing.T) {
			root, artifacts := fixture(t)
			before := treeDigest(t, root)
			interruptTransactionAtCheckpoint(t, root, artifacts, checkpoint)
			if err := recoverArtifactTransaction(root); err != nil {
				t.Fatalf("recover at %s/%d: %v", checkpoint.phase, checkpoint.cursor, err)
			}
			wantTree := before
			if checkpoint.phase == transactionPhaseFinalizing {
				wantTree = committed
			}
			if got := treeDigest(t, root); got != wantTree {
				t.Fatalf(
					"recovery at %s/%d:\n want %s\n  got %s",
					checkpoint.phase,
					checkpoint.cursor,
					wantTree,
					got,
				)
			}
			assertNoTransactionArtifacts(t, root)
		})
	}
}

func TestRecoverTransactionAfterOwnedDirectoryPublishBeforeCheckpoint(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "keep.txt"), []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	before := treeDigest(t, root)
	originalFS := transactionFS
	t.Cleanup(func() { transactionFS = originalFS })
	transactionFS.rename = func(oldPath, newPath string) error {
		if strings.Contains(filepath.ToSlash(oldPath), "/directories/") {
			if err := os.Rename(oldPath, newPath); err != nil {
				return err
			}
			panic("injected directory publication interruption")
		}
		return os.Rename(oldPath, newPath)
	}
	func() {
		defer func() {
			if recovered := recover(); recovered != "injected directory publication interruption" {
				t.Fatalf("directory publication panic = %v", recovered)
			}
		}()
		_ = executeArtifactTransaction(root, []artifact{{
			RelativePath: "new/nested/button.gsx",
			Content:      []byte("button\n"),
			Managed:      true,
		}}, nil, nil)
	}()

	transactionFS = originalFS
	if err := recoverArtifactTransaction(root); err != nil {
		t.Fatalf("recover published owned directory: %v", err)
	}
	if got := treeDigest(t, root); got != before {
		t.Fatalf(
			"owned-directory publication recovery:\n before %s\n  after %s",
			before,
			got,
		)
	}
	assertNoTransactionArtifacts(t, root)
}

func TestRecoverTransactionResumesEveryRollbackInterruption(t *testing.T) {
	points := []string{
		"replacement-removed",
		"backup-restored",
		"restored-published",
		"directory-removed",
		"cleanup",
	}
	for _, point := range points {
		t.Run(point, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "existing.txt"), []byte("old\n"), 0o640); err != nil {
				t.Fatal(err)
			}
			artifacts := []artifact{
				{
					RelativePath: "existing.txt",
					Content:      []byte("new\n"),
					Managed:      true,
				},
				{
					RelativePath: "new/nested/button.gsx",
					Content:      []byte("button\n"),
					Managed:      true,
				},
			}
			before := treeDigest(t, root)
			originalRollbackHook := transactionRollbackHook
			originalCheckpointHook := transactionCheckpointHook
			originalCleanupHook := transactionCleanupHook
			originalRunner := runCommand
			t.Cleanup(func() {
				transactionRollbackHook = originalRollbackHook
				transactionCheckpointHook = originalCheckpointHook
				transactionCleanupHook = originalCleanupHook
				runCommand = originalRunner
			})
			var generationCalls int
			runCommand = func(_ string, name string, args ...string) error {
				if name == "go" && len(args) >= 2 &&
					args[0] == "tool" && args[1] == "gsx" {
					generationCalls++
				}
				return nil
			}
			transactionRollbackHook = func(got transactionRollbackPoint) {
				if got.kind == point {
					panic("injected rollback interruption")
				}
			}
			transactionCheckpointHook = func(checkpoint transactionCheckpoint) {
				if point == "restored-published" &&
					checkpoint.phase == transactionPhaseRestoredGenerationPending {
					panic("injected rollback interruption")
				}
			}
			transactionCleanupHook = func() {
				if point == "cleanup" {
					panic("injected rollback interruption")
				}
			}

			func() {
				defer func() {
					if recovered := recover(); recovered != "injected rollback interruption" {
						t.Fatalf("rollback panic = %v", recovered)
					}
				}()
				_ = executeArtifactTransaction(
					root,
					artifacts,
					func() error { return errors.New("injected generation failure") },
					func() error {
						generationCalls++
						return nil
					},
				)
			}()
			transactionRollbackHook = originalRollbackHook
			transactionCheckpointHook = originalCheckpointHook
			transactionCleanupHook = originalCleanupHook

			if err := recoverArtifactTransaction(root); err != nil {
				t.Fatalf("resume %s interruption: %v", point, err)
			}
			if generationCalls != 1 {
				t.Fatalf("%s generation calls = %d, want 1", point, generationCalls)
			}
			if got := treeDigest(t, root); got != before {
				t.Fatalf(
					"%s recovery did not restore exact tree:\n before %s\n  after %s",
					point,
					before,
					got,
				)
			}
			assertNoTransactionArtifacts(t, root)
		})
	}
}

func TestTransactionGenerationRollbackPhaseMatrix(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "existing.txt"), []byte("old\n"), 0o640); err != nil {
		t.Fatal(err)
	}
	artifacts := []artifact{
		{
			RelativePath: "existing.txt",
			Content:      []byte("new\n"),
			Managed:      true,
		},
		{
			RelativePath: "new/nested/button.gsx",
			Content:      []byte("button\n"),
			Managed:      true,
		},
	}
	originalHook := transactionCheckpointHook
	t.Cleanup(func() { transactionCheckpointHook = originalHook })
	var got []transactionCheckpoint
	transactionCheckpointHook = func(checkpoint transactionCheckpoint) {
		got = append(got, checkpoint)
	}
	generationErr := errors.New("injected generation failure")
	err := executeArtifactTransaction(
		root,
		artifacts,
		func() error { return generationErr },
		func() error { return nil },
	)
	if !errors.Is(err, generationErr) {
		t.Fatalf("transaction error = %v, want generation failure", err)
	}
	want := []transactionCheckpoint{
		{phase: transactionPhasePrepared, cursor: 0},
		{phase: transactionPhaseDirectoryCreatePending, cursor: 0},
		{phase: transactionPhaseDirectoryCreated, cursor: 1},
		{phase: transactionPhaseDirectoryCreatePending, cursor: 1},
		{phase: transactionPhaseDirectoryCreated, cursor: 2},
		{phase: transactionPhaseBackupPending, cursor: 0},
		{phase: transactionPhaseBackedUp, cursor: 0},
		{phase: transactionPhaseReplacePending, cursor: 0},
		{phase: transactionPhaseReplaced, cursor: 1},
		{phase: transactionPhaseBackupPending, cursor: 1},
		{phase: transactionPhaseBackedUp, cursor: 1},
		{phase: transactionPhaseReplacePending, cursor: 1},
		{phase: transactionPhaseReplaced, cursor: 2},
		{phase: transactionPhaseCommitted, cursor: 2},
		{phase: transactionPhaseGenerationPending, cursor: 2},
		{phase: transactionPhaseRollbackPending, cursor: 2},
		{phase: transactionPhaseRollbackPending, cursor: 1},
		{phase: transactionPhaseRollbackPending, cursor: 0},
		{phase: transactionPhaseRestoredGenerationPending, cursor: 0},
		{phase: transactionPhaseRollbackDirectoriesPending, cursor: 2},
		{phase: transactionPhaseRollbackDirectoriesPending, cursor: 1},
		{phase: transactionPhaseRollbackDirectoriesPending, cursor: 0},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("generation rollback checkpoints:\n got: %#v\nwant: %#v", got, want)
	}
	assertNoTransactionArtifacts(t, root)
}

func TestRecoverTransactionResumesEveryRollbackCheckpoint(t *testing.T) {
	checkpoints := []transactionCheckpoint{
		{phase: transactionPhaseRollbackPending, cursor: 2},
		{phase: transactionPhaseRollbackPending, cursor: 1},
		{phase: transactionPhaseRollbackPending, cursor: 0},
		{phase: transactionPhaseRestoredGenerationPending, cursor: 0},
		{phase: transactionPhaseRollbackDirectoriesPending, cursor: 2},
		{phase: transactionPhaseRollbackDirectoriesPending, cursor: 1},
		{phase: transactionPhaseRollbackDirectoriesPending, cursor: 0},
	}
	for _, checkpoint := range checkpoints {
		t.Run(fmt.Sprintf("%s_%d", checkpoint.phase, checkpoint.cursor), func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "existing.txt"), []byte("old\n"), 0o640); err != nil {
				t.Fatal(err)
			}
			artifacts := []artifact{
				{
					RelativePath: "existing.txt",
					Content:      []byte("new\n"),
					Managed:      true,
				},
				{
					RelativePath: "new/nested/button.gsx",
					Content:      []byte("button\n"),
					Managed:      true,
				},
			}
			before := treeDigest(t, root)
			originalHook := transactionCheckpointHook
			originalRunner := runCommand
			t.Cleanup(func() {
				transactionCheckpointHook = originalHook
				runCommand = originalRunner
			})
			var generationCalls int
			runCommand = func(_ string, name string, args ...string) error {
				if name == "go" && len(args) >= 2 &&
					args[0] == "tool" && args[1] == "gsx" {
					generationCalls++
				}
				return nil
			}
			transactionCheckpointHook = func(got transactionCheckpoint) {
				if got == checkpoint {
					panic("injected rollback checkpoint interruption")
				}
			}
			func() {
				defer func() {
					if recovered := recover(); recovered != "injected rollback checkpoint interruption" {
						t.Fatalf("rollback checkpoint panic = %v", recovered)
					}
				}()
				_ = executeArtifactTransaction(
					root,
					artifacts,
					func() error { return errors.New("injected generation failure") },
					func() error {
						generationCalls++
						return nil
					},
				)
			}()
			transactionCheckpointHook = originalHook

			if err := recoverArtifactTransaction(root); err != nil {
				t.Fatalf("recover checkpoint %s/%d: %v", checkpoint.phase, checkpoint.cursor, err)
			}
			if generationCalls != 1 {
				t.Fatalf(
					"checkpoint %s/%d generation calls = %d, want 1",
					checkpoint.phase,
					checkpoint.cursor,
					generationCalls,
				)
			}
			if got := treeDigest(t, root); got != before {
				t.Fatalf(
					"checkpoint %s/%d recovery:\n before %s\n  after %s",
					checkpoint.phase,
					checkpoint.cursor,
					before,
					got,
				)
			}
			assertNoTransactionArtifacts(t, root)
		})
	}
}

func interruptTransactionAtCheckpoint(
	t *testing.T,
	root string,
	artifacts []artifact,
	want transactionCheckpoint,
) {
	t.Helper()
	originalHook := transactionCheckpointHook
	transactionCheckpointHook = func(checkpoint transactionCheckpoint) {
		if checkpoint == want {
			panic("injected transaction interruption")
		}
	}
	defer func() {
		transactionCheckpointHook = originalHook
		if recovered := recover(); recovered != "injected transaction interruption" {
			t.Fatalf("transaction panic = %v", recovered)
		}
	}()
	_ = executeArtifactTransaction(root, artifacts, nil, nil)
}

func transactionFixture(t *testing.T) (string, []artifact) {
	t.Helper()
	root := t.TempDir()
	files := map[string]struct {
		content string
		mode    os.FileMode
	}{
		"gsxui.preset.json":      {content: "old preset\n", mode: 0o640},
		"web/gsxui/theme.css":    {content: "old theme\n", mode: 0o644},
		"ui/button.gsx":          {content: "old button\n", mode: 0o600},
		"gsxui.json":             {content: "old config\n", mode: 0o644},
		"unmanaged/keep.txt":     {content: "user bytes\n", mode: 0o604},
		"unmanaged/nested/a.txt": {content: "nested user bytes\n", mode: 0o644},
	}
	for relative, file := range files {
		target := filepath.Join(root, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, []byte(file.content), file.mode); err != nil {
			t.Fatal(err)
		}
	}
	return root, []artifact{
		{RelativePath: "gsxui.preset.json", Content: []byte("new preset\n"), Managed: true},
		{RelativePath: "web/gsxui/theme.css", Content: []byte("new theme\n"), Managed: true},
		{RelativePath: "ui/button.gsx", Content: []byte("new button\n"), Managed: true},
		{RelativePath: "gsxui.json", Content: []byte("new config\n")},
	}
}

func assertNoTransactionArtifacts(t *testing.T, root string) {
	t.Helper()
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == transactionJournalName ||
			strings.HasPrefix(info.Name(), transactionDirectoryPrefix) ||
			strings.HasPrefix(info.Name(), transactionOwnerPrefix) {
			t.Errorf(
				"transaction artifact remains after handled operation: %s",
				path,
			)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
