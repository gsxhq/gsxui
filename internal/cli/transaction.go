package cli

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

const (
	transactionJournalName     = ".gsxui-transaction.json"
	transactionDirectoryPrefix = ".gsxui-transaction-"
	transactionOwnerPrefix     = ".gsxui-transaction-owner-"
	transactionJournalVersion  = 1
)

type transactionPhase string

const (
	transactionPhasePrepared                   transactionPhase = "prepared"
	transactionPhaseDirectoryCreatePending     transactionPhase = "directory-create-pending"
	transactionPhaseDirectoryCreated           transactionPhase = "directory-created"
	transactionPhaseBackupPending              transactionPhase = "backup-pending"
	transactionPhaseBackedUp                   transactionPhase = "backed-up"
	transactionPhaseReplacePending             transactionPhase = "replace-pending"
	transactionPhaseReplaced                   transactionPhase = "replaced"
	transactionPhaseCommitted                  transactionPhase = "committed"
	transactionPhaseGenerationPending          transactionPhase = "generation-pending"
	transactionPhaseRollbackPending            transactionPhase = "rollback-pending"
	transactionPhaseRestoredGenerationPending  transactionPhase = "restored-generation-pending"
	transactionPhaseRollbackDirectoriesPending transactionPhase = "rollback-directories-pending"
	transactionPhaseFinalizing                 transactionPhase = "finalizing"
)

type transactionCheckpoint struct {
	phase  transactionPhase
	cursor int
}

type transactionFilesystem struct {
	rename        func(string, string) error
	syncDirectory func(string) error
}

var transactionFS = transactionFilesystem{
	rename:        os.Rename,
	syncDirectory: syncOSDirectory,
}

// transactionCheckpointHook is an injected process-interruption seam. Tests
// panic at each durable journal checkpoint; production leaves it as a no-op.
var transactionCheckpointHook = func(transactionCheckpoint) {}

type transactionRollbackPoint struct {
	kind  string
	index int
}

var transactionRollbackHook = func(transactionRollbackPoint) {}
var transactionCleanupHook = func() {}

type transactionJournal struct {
	Version            int                       `json:"version"`
	Directory          string                    `json:"directory"`
	Phase              transactionPhase          `json:"phase"`
	Cursor             int                       `json:"cursor"`
	CreatedDirectories []string                  `json:"createdDirectories,omitempty"`
	Entries            []transactionJournalEntry `json:"entries"`
	RollbackPhase      transactionPhase          `json:"rollbackPhase,omitempty"`
	RollbackCursor     int                       `json:"rollbackCursor,omitempty"`
}

type transactionJournalEntry struct {
	Target         string      `json:"target"`
	Staged         string      `json:"staged"`
	Backup         string      `json:"backup"`
	PreviousExists bool        `json:"previousExists"`
	PreviousHash   string      `json:"previousHash,omitempty"`
	PreviousMode   fs.FileMode `json:"previousMode,omitempty"`
	IntendedHash   string      `json:"intendedHash"`
	IntendedMode   fs.FileMode `json:"intendedMode"`
}

type artifactTransaction struct {
	root        string
	journalPath string
	journal     transactionJournal
}

func executeArtifactTransaction(
	root string,
	artifacts []artifact,
	postCommit func() error,
	postRollback func() error,
) error {
	if err := recoverArtifactTransaction(root); err != nil {
		return err
	}

	changed, err := changedArtifacts(root, artifacts)
	if err != nil {
		return err
	}
	if len(changed) == 0 && postCommit == nil {
		return nil
	}

	transaction, err := prepareArtifactTransaction(root, changed)
	if err != nil {
		return err
	}
	if err := transaction.commit(); err != nil {
		if rollbackErr := transaction.rollback(); rollbackErr != nil {
			return errors.Join(err, fmt.Errorf("roll back filesystem transaction: %w", rollbackErr))
		}
		return err
	}
	if postCommit != nil {
		if err := transaction.persist(
			transactionPhaseGenerationPending,
			len(transaction.journal.Entries),
		); err != nil {
			if rollbackErr := transaction.rollbackAndCleanup(); rollbackErr != nil {
				return errors.Join(err, fmt.Errorf("roll back filesystem transaction: %w", rollbackErr))
			}
			return err
		}
		if err := postCommit(); err != nil {
			rollbackErr := transaction.rollbackToRestoredGeneration()
			var restoreErr error
			if rollbackErr == nil && postRollback != nil {
				restoreErr = postRollback()
			}
			var cleanupErr error
			if rollbackErr == nil {
				cleanupErr = transaction.rollbackDirectories()
				if cleanupErr == nil {
					cleanupErr = cleanupTransaction(
						transaction.root,
						transaction.journalPath,
						transaction.journal,
					)
				}
			}
			return errors.Join(
				err,
				wrapIfNotNil("roll back filesystem transaction", rollbackErr),
				wrapIfNotNil("regenerate restored sources", restoreErr),
				wrapIfNotNil("clean rolled-back filesystem transaction", cleanupErr),
			)
		}
	}
	finalizeErr := transaction.finalize()
	if finalizeErr == nil || transaction.journal.Phase == transactionPhaseFinalizing {
		return finalizeErr
	}
	if postCommit == nil {
		rollbackErr := transaction.rollbackAndCleanup()
		return errors.Join(
			finalizeErr,
			wrapIfNotNil("roll back filesystem transaction", rollbackErr),
		)
	}
	rollbackErr := transaction.rollbackToRestoredGeneration()
	var restoreErr error
	if rollbackErr == nil && postRollback != nil {
		restoreErr = postRollback()
	}
	var cleanupErr error
	if rollbackErr == nil {
		cleanupErr = transaction.rollbackDirectories()
		if cleanupErr == nil {
			cleanupErr = cleanupTransaction(
				transaction.root,
				transaction.journalPath,
				transaction.journal,
			)
		}
	}
	return errors.Join(
		finalizeErr,
		wrapIfNotNil("roll back filesystem transaction", rollbackErr),
		wrapIfNotNil("regenerate restored sources", restoreErr),
		wrapIfNotNil("clean rolled-back filesystem transaction", cleanupErr),
	)
}

func wrapIfNotNil(message string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

func changedArtifacts(root string, artifacts []artifact) ([]artifact, error) {
	changed := make([]artifact, 0, len(artifacts))
	for _, planned := range artifacts {
		target, err := artifactPath(root, planned.RelativePath)
		if err != nil {
			return nil, err
		}
		existing, err := os.ReadFile(target)
		switch {
		case err == nil && bytes.Equal(existing, planned.Content):
			continue
		case err == nil:
		case os.IsNotExist(err):
		default:
			return nil, fmt.Errorf("read artifact %s: %w", target, err)
		}
		changed = append(changed, planned)
	}
	return changed, nil
}

func prepareArtifactTransaction(root string, artifacts []artifact) (*artifactTransaction, error) {
	if err := validateTransactionArtifactOrder(artifacts); err != nil {
		return nil, err
	}
	if err := validateArtifactPlan(root, DefaultConfig(), artifacts, true); err != nil {
		return nil, err
	}
	if err := validateTransactionTargetNamespaces(artifacts); err != nil {
		return nil, err
	}

	rootAbsolute, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve transaction root: %w", err)
	}
	rootResolved, err := filepath.EvalSymlinks(rootAbsolute)
	if err != nil {
		return nil, fmt.Errorf("resolve transaction root: %w", err)
	}
	journalPath, err := artifactPath(rootResolved, transactionJournalName)
	if err != nil {
		return nil, err
	}
	if _, err := os.Lstat(journalPath); err == nil {
		return nil, fmt.Errorf("transaction journal already exists")
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("inspect transaction journal: %w", err)
	}

	transactionDirectory, err := os.MkdirTemp(rootResolved, transactionDirectoryPrefix)
	if err != nil {
		return nil, fmt.Errorf("create transaction directory: %w", err)
	}
	cleanup := true
	defer func() {
		if cleanup {
			if _, journalErr := os.Lstat(journalPath); os.IsNotExist(journalErr) {
				_ = os.RemoveAll(transactionDirectory)
			}
		}
	}()
	if err := os.Chmod(transactionDirectory, 0o700); err != nil {
		return nil, fmt.Errorf("secure transaction directory: %w", err)
	}
	directoryRelative := filepath.Base(transactionDirectory)
	if err := validateTransactionDirectory(directoryRelative); err != nil {
		return nil, err
	}
	for _, name := range []string{"staged", "backup", "directories", "removed"} {
		if err := os.Mkdir(filepath.Join(transactionDirectory, name), 0o700); err != nil {
			return nil, fmt.Errorf("create transaction %s directory: %w", name, err)
		}
	}
	if err := syncDirectory(transactionDirectory); err != nil {
		return nil, fmt.Errorf("sync transaction directory: %w", err)
	}
	if err := syncDirectory(rootResolved); err != nil {
		return nil, fmt.Errorf("sync transaction root: %w", err)
	}

	journal := transactionJournal{
		Version:   transactionJournalVersion,
		Directory: directoryRelative,
		Phase:     transactionPhasePrepared,
		Entries:   make([]transactionJournalEntry, 0, len(artifacts)),
	}
	directories := make(map[string]bool)
	for index, planned := range artifacts {
		target, err := artifactPath(rootResolved, planned.RelativePath)
		if err != nil {
			return nil, err
		}
		info, err := os.Lstat(target)
		previousExists := err == nil
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("inspect transaction target %s: %w", target, err)
		}
		mode := fs.FileMode(0o644)
		previousHash := ""
		if previousExists {
			if !info.Mode().IsRegular() {
				return nil, fmt.Errorf("transaction target %s is not a regular file", target)
			}
			mode = info.Mode().Perm()
			content, err := os.ReadFile(target)
			if err != nil {
				return nil, fmt.Errorf("read transaction target %s: %w", target, err)
			}
			previousHash = contentHash(content)
		}
		for _, relative := range missingParentDirectories(rootResolved, planned.RelativePath) {
			directories[relative] = true
		}

		stagedRelative := path.Join(directoryRelative, "staged", fmt.Sprintf("%06d", index))
		backupRelative := path.Join(directoryRelative, "backup", fmt.Sprintf("%06d", index))
		stagedPath := filepath.Join(rootResolved, filepath.FromSlash(stagedRelative))
		if err := writeDurableFile(stagedPath, planned.Content, mode); err != nil {
			return nil, fmt.Errorf("stage transaction target %s: %w", planned.RelativePath, err)
		}
		journal.Entries = append(journal.Entries, transactionJournalEntry{
			Target:         planned.RelativePath,
			Staged:         stagedRelative,
			Backup:         backupRelative,
			PreviousExists: previousExists,
			PreviousHash:   previousHash,
			PreviousMode:   mode,
			IntendedHash:   contentHash(planned.Content),
			IntendedMode:   mode,
		})
	}
	journal.CreatedDirectories = make([]string, 0, len(directories))
	for relative := range directories {
		journal.CreatedDirectories = append(journal.CreatedDirectories, relative)
	}
	slices.SortFunc(journal.CreatedDirectories, func(left, right string) int {
		leftDepth := strings.Count(left, "/")
		rightDepth := strings.Count(right, "/")
		if leftDepth != rightDepth {
			return leftDepth - rightDepth
		}
		return strings.Compare(left, right)
	})
	markerName := transactionOwnerPrefix +
		strings.TrimPrefix(directoryRelative, transactionDirectoryPrefix)
	for index := range journal.CreatedDirectories {
		stagedDirectory := filepath.Join(
			transactionDirectory,
			"directories",
			fmt.Sprintf("%06d", index),
		)
		if err := os.Mkdir(stagedDirectory, 0o755); err != nil {
			return nil, fmt.Errorf("stage transaction-created directory: %w", err)
		}
		if err := os.Chmod(stagedDirectory, 0o755); err != nil {
			return nil, fmt.Errorf("set staged directory mode: %w", err)
		}
		if err := writeDurableFile(
			filepath.Join(stagedDirectory, markerName),
			[]byte(directoryRelative+"\n"),
			0o600,
		); err != nil {
			return nil, fmt.Errorf("stage transaction directory ownership marker: %w", err)
		}
		if err := syncDirectory(stagedDirectory); err != nil {
			return nil, fmt.Errorf("sync staged transaction-created directory: %w", err)
		}
	}
	if err := syncDirectory(filepath.Join(transactionDirectory, "staged")); err != nil {
		return nil, fmt.Errorf("sync staged transaction files: %w", err)
	}
	if err := syncDirectory(filepath.Join(transactionDirectory, "directories")); err != nil {
		return nil, fmt.Errorf("sync staged transaction directories: %w", err)
	}

	transaction := &artifactTransaction{
		root:        rootResolved,
		journalPath: journalPath,
		journal:     journal,
	}
	if err := transaction.persist(transactionPhasePrepared, 0); err != nil {
		publishErr := fmt.Errorf("publish transaction journal: %w", err)
		if _, journalErr := os.Lstat(journalPath); journalErr == nil {
			cleanup = false
			if rollbackErr := transaction.rollback(); rollbackErr != nil {
				return nil, errors.Join(
					publishErr,
					fmt.Errorf("roll back published transaction journal: %w", rollbackErr),
				)
			}
		}
		return nil, publishErr
	}
	cleanup = false
	return transaction, nil
}

func validateTransactionArtifactOrder(artifacts []artifact) error {
	for index, planned := range artifacts {
		if planned.RelativePath == "gsxui.json" && index != len(artifacts)-1 {
			return fmt.Errorf("gsxui.json must be the final transaction artifact")
		}
	}
	return nil
}

func validateTransactionTargetNamespaces(artifacts []artifact) error {
	journalIdentity := portableArtifactIdentity(transactionJournalName)
	for _, planned := range artifacts {
		identity := portableArtifactIdentity(planned.RelativePath)
		if portableArtifactIdentitiesEqual(identity, journalIdentity) ||
			portableArtifactIdentityIsAncestor(journalIdentity, identity) {
			return fmt.Errorf(
				"artifact path %q is a portable alias of or descendant under the transaction journal",
				planned.RelativePath,
			)
		}
		if len(identity) > 0 && strings.HasPrefix(identity[0], transactionDirectoryPrefix) {
			return fmt.Errorf("artifact path %q uses the reserved transaction directory namespace", planned.RelativePath)
		}
		for _, segment := range identity {
			if strings.HasPrefix(segment, transactionOwnerPrefix) {
				return fmt.Errorf("artifact path %q uses the reserved transaction owner namespace", planned.RelativePath)
			}
		}
	}
	return nil
}

func missingParentDirectories(root, relative string) []string {
	parent := path.Dir(relative)
	if parent == "." {
		return nil
	}
	segments := strings.Split(parent, "/")
	missing := make([]string, 0, len(segments))
	for index := range segments {
		candidate := strings.Join(segments[:index+1], "/")
		if _, err := os.Lstat(filepath.Join(root, filepath.FromSlash(candidate))); os.IsNotExist(err) {
			missing = append(missing, candidate)
		}
	}
	return missing
}

func writeDurableFile(target string, content []byte, mode fs.FileMode) error {
	file, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode.Perm())
	if err != nil {
		return err
	}
	closeNeeded := true
	defer func() {
		if closeNeeded {
			_ = file.Close()
		}
	}()
	if err := file.Chmod(mode.Perm()); err != nil {
		return err
	}
	written, err := file.Write(content)
	if err != nil {
		return err
	}
	if written != len(content) {
		return fmt.Errorf("wrote %d of %d bytes", written, len(content))
	}
	if err := file.Sync(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	closeNeeded = false
	return nil
}

func (transaction *artifactTransaction) persist(phase transactionPhase, cursor int) error {
	next := transaction.journal
	next.Phase = phase
	next.Cursor = cursor
	content, err := json.MarshalIndent(next, "", "  ")
	if err != nil {
		return fmt.Errorf("encode transaction journal: %w", err)
	}
	content = append(content, '\n')
	temp, err := os.CreateTemp(transaction.root, "."+transactionJournalName+".tmp-")
	if err != nil {
		return fmt.Errorf("create temporary transaction journal: %w", err)
	}
	tempPath := temp.Name()
	defer os.Remove(tempPath)
	closeNeeded := true
	defer func() {
		if closeNeeded {
			_ = temp.Close()
		}
	}()
	if err := temp.Chmod(0o600); err != nil {
		return fmt.Errorf("secure temporary transaction journal: %w", err)
	}
	if _, err := temp.Write(content); err != nil {
		return fmt.Errorf("write temporary transaction journal: %w", err)
	}
	if err := temp.Sync(); err != nil {
		return fmt.Errorf("sync temporary transaction journal: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close temporary transaction journal: %w", err)
	}
	closeNeeded = false
	if err := transactionFS.rename(tempPath, transaction.journalPath); err != nil {
		return fmt.Errorf("replace transaction journal: %w", err)
	}
	if err := syncDirectory(transaction.root); err != nil {
		return fmt.Errorf("sync transaction journal directory: %w", err)
	}
	transaction.journal = next
	transactionCheckpointHook(transactionCheckpoint{phase: phase, cursor: cursor})
	return nil
}

func (transaction *artifactTransaction) commit() error {
	for index, relative := range transaction.journal.CreatedDirectories {
		if err := transaction.persist(transactionPhaseDirectoryCreatePending, index); err != nil {
			return err
		}
		target, err := artifactDirectoryPath(transaction.root, relative)
		if err != nil {
			return err
		}
		stagedDirectory, err := transaction.stagedDirectory(index)
		if err != nil {
			return err
		}
		if err := transactionFS.rename(stagedDirectory, target); err != nil {
			return fmt.Errorf("publish transaction-created directory %s: %w", relative, err)
		}
		if err := syncRenameDirectories(stagedDirectory, target); err != nil {
			return fmt.Errorf("sync transaction-created directory %s: %w", relative, err)
		}
		if err := transaction.persist(transactionPhaseDirectoryCreated, index+1); err != nil {
			return err
		}
	}

	for index, entry := range transaction.journal.Entries {
		target, staged, backup, err := transaction.paths(entry)
		if err != nil {
			return err
		}
		if err := transaction.persist(transactionPhaseBackupPending, index); err != nil {
			return err
		}
		if entry.PreviousExists {
			if err := transactionFS.rename(target, backup); err != nil {
				return fmt.Errorf("back up transaction target %s: %w", entry.Target, err)
			}
			if err := syncRenameDirectories(target, backup); err != nil {
				return fmt.Errorf("sync backup rename for %s: %w", entry.Target, err)
			}
		}
		if err := transaction.persist(transactionPhaseBackedUp, index); err != nil {
			return err
		}
		if err := transaction.persist(transactionPhaseReplacePending, index); err != nil {
			return err
		}
		if err := transactionFS.rename(staged, target); err != nil {
			return fmt.Errorf("replace transaction target %s: %w", entry.Target, err)
		}
		if err := syncRenameDirectories(staged, target); err != nil {
			return fmt.Errorf("sync target replacement for %s: %w", entry.Target, err)
		}
		if err := transaction.persist(transactionPhaseReplaced, index+1); err != nil {
			return err
		}
	}
	return transaction.persist(transactionPhaseCommitted, len(transaction.journal.Entries))
}

func (transaction *artifactTransaction) stagedDirectory(index int) (string, error) {
	relative := path.Join(
		transaction.journal.Directory,
		"directories",
		fmt.Sprintf("%06d", index),
	)
	return transactionInternalDirectoryPath(
		transaction.root,
		transaction.journal.Directory,
		relative,
	)
}

func (transaction *artifactTransaction) directoryOwnershipMarker(relative string) (string, error) {
	markerRelative := path.Join(relative, transactionOwnerPrefix+
		strings.TrimPrefix(transaction.journal.Directory, transactionDirectoryPrefix))
	return artifactPath(transaction.root, markerRelative)
}

func (transaction *artifactTransaction) paths(
	entry transactionJournalEntry,
) (target, staged, backup string, err error) {
	target, err = artifactPath(transaction.root, entry.Target)
	if err != nil {
		return "", "", "", err
	}
	staged, err = transactionInternalPath(transaction.root, transaction.journal.Directory, entry.Staged)
	if err != nil {
		return "", "", "", err
	}
	backup, err = transactionInternalPath(transaction.root, transaction.journal.Directory, entry.Backup)
	if err != nil {
		return "", "", "", err
	}
	return target, staged, backup, nil
}

func (transaction *artifactTransaction) rollback() error {
	return transaction.rollbackAndCleanup()
}

func (transaction *artifactTransaction) rollbackAndCleanup() error {
	if err := transaction.rollbackEntries(); err != nil {
		return err
	}
	if err := transaction.rollbackDirectories(); err != nil {
		return err
	}
	return cleanupTransaction(transaction.root, transaction.journalPath, transaction.journal)
}

func (transaction *artifactTransaction) rollbackToRestoredGeneration() error {
	if err := transaction.rollbackEntries(); err != nil {
		return err
	}
	return transaction.persist(
		transactionPhaseRestoredGenerationPending,
		0,
	)
}

func (transaction *artifactTransaction) finalize() error {
	if err := transaction.persist(
		transactionPhaseFinalizing,
		len(transaction.journal.Entries),
	); err != nil {
		return fmt.Errorf("mark transaction for finalization: %w", err)
	}
	return transaction.resumeFinalization()
}

func (transaction *artifactTransaction) resumeFinalization() error {
	if err := transaction.removeDirectoryOwnershipMarkers(true); err != nil {
		return err
	}
	transactionDirectory := filepath.Join(transaction.root, transaction.journal.Directory)
	if err := os.RemoveAll(transactionDirectory); err != nil {
		return fmt.Errorf("remove committed transaction directory: %w", err)
	}
	if err := syncDirectory(transaction.root); err != nil {
		return fmt.Errorf("sync committed transaction cleanup: %w", err)
	}
	if err := os.Remove(transaction.journalPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove committed transaction journal: %w", err)
	}
	if err := syncDirectory(transaction.root); err != nil {
		republishErr := transaction.persist(
			transactionPhaseFinalizing,
			len(transaction.journal.Entries),
		)
		return errors.Join(
			fmt.Errorf("sync committed transaction journal removal: %w", err),
			wrapIfNotNil("republish finalization journal", republishErr),
		)
	}
	return nil
}

func (transaction *artifactTransaction) removeDirectoryOwnershipMarkers(
	allowMissing bool,
) error {
	for _, relative := range transaction.journal.CreatedDirectories {
		marker, err := transaction.directoryOwnershipMarker(relative)
		if err != nil {
			return err
		}
		content, err := os.ReadFile(marker)
		switch {
		case os.IsNotExist(err) && allowMissing:
			continue
		case err != nil:
			return fmt.Errorf("read ownership marker for %s: %w", relative, err)
		case string(content) != transaction.journal.Directory+"\n":
			return fmt.Errorf("ownership marker for %s changed after interruption", relative)
		}
		if err := os.Remove(marker); err != nil {
			return fmt.Errorf("remove ownership marker for %s: %w", relative, err)
		}
		directory, err := artifactDirectoryPath(transaction.root, relative)
		if err != nil {
			return err
		}
		if err := syncDirectory(directory); err != nil {
			return fmt.Errorf("sync ownership marker removal for %s: %w", relative, err)
		}
	}
	return nil
}

func recoverArtifactTransaction(root string) error {
	rootAbsolute, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolve recovery root: %w", err)
	}
	rootResolved, err := filepath.EvalSymlinks(rootAbsolute)
	if err != nil {
		return fmt.Errorf("resolve recovery root: %w", err)
	}
	journalPath, err := artifactPath(rootResolved, transactionJournalName)
	if err != nil {
		return fmt.Errorf("recover transaction: %w", err)
	}
	content, err := os.ReadFile(journalPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read transaction journal: %w", err)
	}
	journal, err := parseTransactionJournal(content)
	if err != nil {
		return fmt.Errorf("parse transaction journal: %w", err)
	}
	if err := validateRecoveryJournal(rootResolved, journal); err != nil {
		return fmt.Errorf("validate transaction journal: %w", err)
	}
	transaction := &artifactTransaction{
		root:        rootResolved,
		journalPath: journalPath,
		journal:     journal,
	}
	switch journal.Phase {
	case transactionPhaseFinalizing:
		if err := transaction.resumeFinalization(); err != nil {
			return fmt.Errorf("resume transaction finalization: %w", err)
		}
		return nil
	case transactionPhaseRestoredGenerationPending:
		if err := validateRestoredJournalState(rootResolved, journal); err != nil {
			return fmt.Errorf("validate restored transaction state: %w", err)
		}
		if err := generateProject(rootResolved); err != nil {
			return fmt.Errorf("regenerate restored sources: %w", err)
		}
		if err := transaction.rollbackDirectories(); err != nil {
			return fmt.Errorf("remove restored transaction directories: %w", err)
		}
		if err := cleanupTransaction(rootResolved, journalPath, journal); err != nil {
			return fmt.Errorf("clean recovered transaction: %w", err)
		}
		return nil
	case transactionPhaseRollbackDirectoriesPending:
		if err := transaction.rollbackDirectories(); err != nil {
			return fmt.Errorf("resume transaction directory rollback: %w", err)
		}
		if err := cleanupTransaction(rootResolved, journalPath, transaction.journal); err != nil {
			return fmt.Errorf("clean recovered transaction: %w", err)
		}
		return nil
	}
	if err := transaction.rollbackEntries(); err != nil {
		return fmt.Errorf("recover transaction: %w", err)
	}
	if transaction.journal.RollbackPhase == transactionPhaseGenerationPending {
		if err := transaction.persist(
			transactionPhaseRestoredGenerationPending,
			0,
		); err != nil {
			return fmt.Errorf("mark restored sources for regeneration: %w", err)
		}
		if err := generateProject(rootResolved); err != nil {
			return fmt.Errorf("regenerate restored sources: %w", err)
		}
	}
	if err := transaction.rollbackDirectories(); err != nil {
		return fmt.Errorf("remove restored transaction directories: %w", err)
	}
	if err := cleanupTransaction(rootResolved, journalPath, transaction.journal); err != nil {
		return fmt.Errorf("clean recovered transaction: %w", err)
	}
	return nil
}

func parseTransactionJournal(content []byte) (transactionJournal, error) {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	var journal transactionJournal
	if err := decoder.Decode(&journal); err != nil {
		return transactionJournal{}, err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return transactionJournal{}, fmt.Errorf("unexpected trailing JSON value")
		}
		return transactionJournal{}, err
	}
	return journal, nil
}

func validateRecoveryJournal(root string, journal transactionJournal) error {
	if journal.Version != transactionJournalVersion {
		return fmt.Errorf("unsupported transaction journal version %d", journal.Version)
	}
	if err := validateTransactionDirectory(journal.Directory); err != nil {
		return err
	}
	switch journal.Phase {
	case transactionPhasePrepared,
		transactionPhaseDirectoryCreatePending,
		transactionPhaseDirectoryCreated,
		transactionPhaseBackupPending,
		transactionPhaseBackedUp,
		transactionPhaseReplacePending,
		transactionPhaseReplaced,
		transactionPhaseCommitted,
		transactionPhaseGenerationPending,
		transactionPhaseRollbackPending,
		transactionPhaseRestoredGenerationPending,
		transactionPhaseRollbackDirectoriesPending,
		transactionPhaseFinalizing:
	default:
		return fmt.Errorf("unknown transaction phase %q", journal.Phase)
	}
	if err := validateTransactionPhaseCursor(journal); err != nil {
		return err
	}
	if journal.Phase == transactionPhaseRollbackPending ||
		journal.Phase == transactionPhaseRestoredGenerationPending ||
		journal.Phase == transactionPhaseRollbackDirectoriesPending {
		if !isRollbackSourcePhase(journal.RollbackPhase) {
			return fmt.Errorf("invalid rollback source phase %q", journal.RollbackPhase)
		}
		source := journal
		source.Phase = journal.RollbackPhase
		source.Cursor = journal.RollbackCursor
		source.RollbackPhase = ""
		source.RollbackCursor = 0
		if err := validateTransactionPhaseCursor(source); err != nil {
			return fmt.Errorf("invalid rollback source state: %w", err)
		}
	} else if journal.RollbackPhase != "" || journal.RollbackCursor != 0 {
		return fmt.Errorf("transaction phase %s has unexpected rollback state", journal.Phase)
	}
	targets := make([]artifact, len(journal.Entries))
	for index, entry := range journal.Entries {
		targets[index] = artifact{RelativePath: entry.Target}
	}
	if err := validateTransactionTargetNamespaces(targets); err != nil {
		return err
	}
	portableTargets := make([][]string, 0, len(journal.Entries))
	for _, entry := range journal.Entries {
		portable := portableArtifactIdentity(entry.Target)
		for previousIndex, previous := range portableTargets {
			switch {
			case portableArtifactIdentitiesEqual(previous, portable):
				return fmt.Errorf(
					"transaction targets %q and %q are portable aliases",
					journal.Entries[previousIndex].Target,
					entry.Target,
				)
			case portableArtifactIdentityIsAncestor(previous, portable):
				return fmt.Errorf(
					"portable transaction target %q is an ancestor of %q",
					journal.Entries[previousIndex].Target,
					entry.Target,
				)
			case portableArtifactIdentityIsAncestor(portable, previous):
				return fmt.Errorf(
					"portable transaction target %q is an ancestor of %q",
					entry.Target,
					journal.Entries[previousIndex].Target,
				)
			}
		}
		portableTargets = append(portableTargets, portable)
	}
	for index, entry := range journal.Entries {
		if _, err := artifactPath(root, entry.Target); err != nil {
			return fmt.Errorf("entry %d target: %w", index, err)
		}
		wantStaged := path.Join(
			journal.Directory,
			"staged",
			fmt.Sprintf("%06d", index),
		)
		if entry.Staged != wantStaged {
			return fmt.Errorf(
				"entry %d staged path %q must be %q",
				index,
				entry.Staged,
				wantStaged,
			)
		}
		if _, err := transactionInternalPath(root, journal.Directory, entry.Staged); err != nil {
			return fmt.Errorf("entry %d staged path: %w", index, err)
		}
		wantBackup := path.Join(
			journal.Directory,
			"backup",
			fmt.Sprintf("%06d", index),
		)
		if entry.Backup != wantBackup {
			return fmt.Errorf(
				"entry %d backup path %q must be %q",
				index,
				entry.Backup,
				wantBackup,
			)
		}
		if _, err := transactionInternalPath(root, journal.Directory, entry.Backup); err != nil {
			return fmt.Errorf("entry %d backup path: %w", index, err)
		}
		if err := validateTransactionHash(entry.IntendedHash); err != nil {
			return fmt.Errorf("entry %d intended hash: %w", index, err)
		}
		if entry.IntendedMode.Perm() != entry.IntendedMode || entry.IntendedMode.Perm() == 0 {
			return fmt.Errorf("entry %d has invalid intended mode %v", index, entry.IntendedMode)
		}
		if entry.PreviousExists {
			if err := validateTransactionHash(entry.PreviousHash); err != nil {
				return fmt.Errorf("entry %d previous hash: %w", index, err)
			}
			if entry.PreviousMode.Perm() != entry.PreviousMode || entry.PreviousMode.Perm() == 0 {
				return fmt.Errorf("entry %d has invalid previous mode %v", index, entry.PreviousMode)
			}
		} else if entry.PreviousHash != "" || entry.PreviousMode != entry.IntendedMode {
			return fmt.Errorf("entry %d records inconsistent state for an absent target", index)
		}
	}
	if err := validateCreatedDirectories(journal); err != nil {
		return err
	}
	if err := validatePublishedDirectoryOwnership(root, journal); err != nil {
		return err
	}
	for _, relative := range journal.CreatedDirectories {
		if _, err := artifactDirectoryPath(root, relative); err != nil {
			return fmt.Errorf("created directory %q: %w", relative, err)
		}
	}
	return nil
}

func validatePublishedDirectoryOwnership(root string, journal transactionJournal) error {
	listed := make([][]string, 0, len(journal.CreatedDirectories))
	for _, relative := range journal.CreatedDirectories {
		listed = append(listed, portableArtifactIdentity(relative))
	}
	markerName := transactionOwnerPrefix +
		strings.TrimPrefix(journal.Directory, transactionDirectoryPrefix)
	checked := make(map[string]bool)
	for _, entry := range journal.Entries {
		parent := path.Dir(entry.Target)
		if parent == "." {
			continue
		}
		segments := strings.Split(parent, "/")
		for index := range segments {
			relative := strings.Join(segments[:index+1], "/")
			identity := strings.Join(portableArtifactIdentity(relative), "\x00")
			if checked[identity] {
				continue
			}
			checked[identity] = true
			directory, err := artifactDirectoryPath(root, relative)
			if err != nil {
				return err
			}
			marker := filepath.Join(directory, markerName)
			info, err := os.Lstat(marker)
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				return fmt.Errorf("inspect ownership marker for %s: %w", relative, err)
			}
			if !info.Mode().IsRegular() {
				return fmt.Errorf("ownership marker for %s is not a regular file", relative)
			}
			content, err := os.ReadFile(marker)
			if err != nil {
				return fmt.Errorf("read ownership marker for %s: %w", relative, err)
			}
			if string(content) != journal.Directory+"\n" {
				return fmt.Errorf("ownership marker for %s changed after interruption", relative)
			}
			found := false
			candidate := portableArtifactIdentity(relative)
			for _, expected := range listed {
				if portableArtifactIdentitiesEqual(candidate, expected) {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf(
					"published ownership marker for %s is omitted from created directories",
					relative,
				)
			}
		}
	}
	return nil
}

func validateTransactionPhaseCursor(journal transactionJournal) error {
	if journal.Cursor < 0 {
		return fmt.Errorf("transaction cursor must not be negative")
	}
	switch journal.Phase {
	case transactionPhasePrepared:
		if journal.Cursor != 0 {
			return fmt.Errorf("prepared transaction cursor must be 0")
		}
	case transactionPhaseDirectoryCreatePending:
		if journal.Cursor >= len(journal.CreatedDirectories) {
			return fmt.Errorf("directory-create-pending cursor must identify a directory")
		}
	case transactionPhaseDirectoryCreated:
		if journal.Cursor == 0 || journal.Cursor > len(journal.CreatedDirectories) {
			return fmt.Errorf("directory-created cursor must identify completed directories")
		}
	case transactionPhaseBackupPending,
		transactionPhaseBackedUp,
		transactionPhaseReplacePending:
		if journal.Cursor >= len(journal.Entries) {
			return fmt.Errorf("%s transaction cursor must identify an entry", journal.Phase)
		}
	case transactionPhaseReplaced:
		if journal.Cursor == 0 || journal.Cursor > len(journal.Entries) {
			return fmt.Errorf("replaced transaction cursor must identify completed entries")
		}
	case transactionPhaseCommitted, transactionPhaseGenerationPending, transactionPhaseFinalizing:
		if journal.Cursor != len(journal.Entries) {
			return fmt.Errorf(
				"%s transaction cursor must equal the entry count",
				journal.Phase,
			)
		}
	case transactionPhaseRollbackPending:
		if journal.Cursor > len(journal.Entries) {
			return fmt.Errorf("%s cursor is outside the entry range", journal.Phase)
		}
	case transactionPhaseRestoredGenerationPending:
		if journal.Cursor != 0 {
			return fmt.Errorf("restored-generation-pending cursor must be 0")
		}
	case transactionPhaseRollbackDirectoriesPending:
		if journal.Cursor > len(journal.CreatedDirectories) {
			return fmt.Errorf("rollback-directories-pending cursor is outside the directory range")
		}
	}
	return nil
}

func isRollbackSourcePhase(phase transactionPhase) bool {
	switch phase {
	case transactionPhasePrepared,
		transactionPhaseDirectoryCreatePending,
		transactionPhaseDirectoryCreated,
		transactionPhaseBackupPending,
		transactionPhaseBackedUp,
		transactionPhaseReplacePending,
		transactionPhaseReplaced,
		transactionPhaseCommitted,
		transactionPhaseGenerationPending:
		return true
	default:
		return false
	}
}

func validateCreatedDirectories(journal transactionJournal) error {
	want := append([]string(nil), journal.CreatedDirectories...)
	slices.SortFunc(want, func(left, right string) int {
		leftDepth := strings.Count(left, "/")
		rightDepth := strings.Count(right, "/")
		if leftDepth != rightDepth {
			return leftDepth - rightDepth
		}
		return strings.Compare(left, right)
	})
	if !slices.Equal(want, journal.CreatedDirectories) {
		return fmt.Errorf("created directories must be in canonical parent-first order")
	}
	portable := make([][]string, 0, len(journal.CreatedDirectories))
	for _, relative := range journal.CreatedDirectories {
		identity := portableArtifactIdentity(relative)
		for _, previous := range portable {
			if portableArtifactIdentitiesEqual(previous, identity) {
				return fmt.Errorf("created directories contain a portable alias")
			}
		}
		portable = append(portable, identity)
		owned := false
		for _, entry := range journal.Entries {
			target := portableArtifactIdentity(entry.Target)
			if portableArtifactIdentityIsAncestor(identity, target) {
				owned = true
				break
			}
		}
		if !owned {
			return fmt.Errorf(
				"created directory %q is not a parent of a transaction target",
				relative,
			)
		}
	}
	return nil
}

func validateTransactionHash(hash string) error {
	if len(hash) != 64 {
		return fmt.Errorf("must be 64 lowercase hexadecimal characters")
	}
	decoded, err := hex.DecodeString(hash)
	if err != nil || hex.EncodeToString(decoded) != hash {
		return fmt.Errorf("must be 64 lowercase hexadecimal characters")
	}
	return nil
}

func validateTransactionDirectory(relative string) error {
	if err := validateManagedPath(relative); err != nil {
		return fmt.Errorf("transaction directory %q: %w", relative, err)
	}
	if path.Dir(relative) != "." || !strings.HasPrefix(relative, transactionDirectoryPrefix) ||
		len(relative) == len(transactionDirectoryPrefix) {
		return fmt.Errorf("invalid transaction directory %q", relative)
	}
	return nil
}

func transactionInternalPath(root, directory, relative string) (string, error) {
	if err := validateManagedPath(relative); err != nil {
		return "", err
	}
	if relative != directory && !strings.HasPrefix(relative, directory+"/") {
		return "", fmt.Errorf("transaction path %q is outside %q", relative, directory)
	}
	return safeProjectPath(root, relative, false)
}

func transactionInternalDirectoryPath(root, directory, relative string) (string, error) {
	if err := validateManagedPath(relative); err != nil {
		return "", err
	}
	if relative != directory && !strings.HasPrefix(relative, directory+"/") {
		return "", fmt.Errorf("transaction path %q is outside %q", relative, directory)
	}
	return safeProjectPath(root, relative, true)
}

type transactionFileState struct {
	exists bool
	hash   string
	mode   fs.FileMode
}

type transactionEntryState struct {
	target transactionFileState
	staged transactionFileState
	backup transactionFileState
}

type expectedTransactionEntryState int

const (
	expectedTransactionEntryOriginal expectedTransactionEntryState = iota
	expectedTransactionEntryBackedUp
	expectedTransactionEntryReplaced
)

func validateRollbackJournalState(root string, journal transactionJournal) error {
	for index, entry := range journal.Entries {
		state, err := inspectTransactionEntryState(root, journal, entry)
		if err != nil {
			return err
		}
		allowed := allowedRollbackEntryStates(journal, index)
		matched := false
		for _, expected := range allowed {
			if transactionEntryStateMatches(entry, state, expected) {
				matched = true
				break
			}
		}
		if !matched {
			return fmt.Errorf(
				"%s transaction target %s changed after interruption",
				journal.Phase,
				entry.Target,
			)
		}
	}
	return nil
}

func allowedRollbackEntryStates(
	journal transactionJournal,
	index int,
) []expectedTransactionEntryState {
	switch journal.Phase {
	case transactionPhasePrepared,
		transactionPhaseDirectoryCreatePending,
		transactionPhaseDirectoryCreated:
		return []expectedTransactionEntryState{expectedTransactionEntryOriginal}
	case transactionPhaseBackupPending:
		switch {
		case index < journal.Cursor:
			return []expectedTransactionEntryState{expectedTransactionEntryReplaced}
		case index == journal.Cursor:
			return []expectedTransactionEntryState{
				expectedTransactionEntryOriginal,
				expectedTransactionEntryBackedUp,
			}
		default:
			return []expectedTransactionEntryState{expectedTransactionEntryOriginal}
		}
	case transactionPhaseBackedUp:
		switch {
		case index < journal.Cursor:
			return []expectedTransactionEntryState{expectedTransactionEntryReplaced}
		case index == journal.Cursor:
			return []expectedTransactionEntryState{expectedTransactionEntryBackedUp}
		default:
			return []expectedTransactionEntryState{expectedTransactionEntryOriginal}
		}
	case transactionPhaseReplacePending:
		switch {
		case index < journal.Cursor:
			return []expectedTransactionEntryState{expectedTransactionEntryReplaced}
		case index == journal.Cursor:
			return []expectedTransactionEntryState{
				expectedTransactionEntryBackedUp,
				expectedTransactionEntryReplaced,
			}
		default:
			return []expectedTransactionEntryState{expectedTransactionEntryOriginal}
		}
	case transactionPhaseReplaced:
		if index < journal.Cursor {
			return []expectedTransactionEntryState{expectedTransactionEntryReplaced}
		}
		return []expectedTransactionEntryState{expectedTransactionEntryOriginal}
	case transactionPhaseCommitted, transactionPhaseGenerationPending:
		return []expectedTransactionEntryState{expectedTransactionEntryReplaced}
	default:
		return nil
	}
}

func inspectTransactionEntryState(
	root string,
	journal transactionJournal,
	entry transactionJournalEntry,
) (transactionEntryState, error) {
	target, err := artifactPath(root, entry.Target)
	if err != nil {
		return transactionEntryState{}, err
	}
	staged, err := transactionInternalPath(root, journal.Directory, entry.Staged)
	if err != nil {
		return transactionEntryState{}, err
	}
	backup, err := transactionInternalPath(root, journal.Directory, entry.Backup)
	if err != nil {
		return transactionEntryState{}, err
	}
	targetState, err := inspectTransactionFile(target)
	if err != nil {
		return transactionEntryState{}, fmt.Errorf("inspect target %s: %w", entry.Target, err)
	}
	stagedState, err := inspectTransactionFile(staged)
	if err != nil {
		return transactionEntryState{}, fmt.Errorf("inspect staged file for %s: %w", entry.Target, err)
	}
	backupState, err := inspectTransactionFile(backup)
	if err != nil {
		return transactionEntryState{}, fmt.Errorf("inspect backup for %s: %w", entry.Target, err)
	}
	return transactionEntryState{
		target: targetState,
		staged: stagedState,
		backup: backupState,
	}, nil
}

func inspectTransactionFile(target string) (transactionFileState, error) {
	info, err := os.Lstat(target)
	if os.IsNotExist(err) {
		return transactionFileState{}, nil
	}
	if err != nil {
		return transactionFileState{}, err
	}
	if !info.Mode().IsRegular() {
		return transactionFileState{}, fmt.Errorf("is not a regular file")
	}
	content, err := os.ReadFile(target)
	if err != nil {
		return transactionFileState{}, err
	}
	return transactionFileState{
		exists: true,
		hash:   contentHash(content),
		mode:   info.Mode().Perm(),
	}, nil
}

func transactionEntryStateMatches(
	entry transactionJournalEntry,
	state transactionEntryState,
	expected expectedTransactionEntryState,
) bool {
	previous := transactionFileState{
		exists: entry.PreviousExists,
		hash:   entry.PreviousHash,
		mode:   entry.PreviousMode.Perm(),
	}
	intended := transactionFileState{
		exists: true,
		hash:   entry.IntendedHash,
		mode:   entry.IntendedMode.Perm(),
	}
	absent := transactionFileState{}
	switch expected {
	case expectedTransactionEntryOriginal:
		return transactionFileStatesEqual(state.target, previous) &&
			transactionFileStatesEqual(state.staged, intended) &&
			transactionFileStatesEqual(state.backup, absent)
	case expectedTransactionEntryBackedUp:
		return transactionFileStatesEqual(state.target, absent) &&
			transactionFileStatesEqual(state.staged, intended) &&
			transactionFileStatesEqual(state.backup, previous)
	case expectedTransactionEntryReplaced:
		return transactionFileStatesEqual(state.target, intended) &&
			transactionFileStatesEqual(state.staged, absent) &&
			transactionFileStatesEqual(state.backup, previous)
	default:
		return false
	}
}

func transactionFileStatesEqual(left, right transactionFileState) bool {
	if left.exists != right.exists {
		return false
	}
	return !left.exists || left.hash == right.hash && left.mode.Perm() == right.mode.Perm()
}

func validateRestoredJournalState(root string, journal transactionJournal) error {
	for _, entry := range journal.Entries {
		state, err := inspectTransactionEntryState(root, journal, entry)
		if err != nil {
			return err
		}
		previous := transactionFileState{
			exists: entry.PreviousExists,
			hash:   entry.PreviousHash,
			mode:   entry.PreviousMode.Perm(),
		}
		intended := transactionFileState{
			exists: true,
			hash:   entry.IntendedHash,
			mode:   entry.IntendedMode.Perm(),
		}
		if !transactionFileStatesEqual(state.target, previous) ||
			state.backup.exists ||
			(state.staged.exists && !transactionFileStatesEqual(state.staged, intended)) {
			return fmt.Errorf(
				"restored transaction target %s changed after interruption",
				entry.Target,
			)
		}
	}
	return nil
}

func (transaction *artifactTransaction) rollbackEntries() error {
	if transaction.journal.Phase != transactionPhaseRollbackPending {
		if err := validateRollbackJournalState(transaction.root, transaction.journal); err != nil {
			return err
		}
		if err := validateTransactionCreatedDirectoryState(
			transaction.root,
			transaction.journal,
		); err != nil {
			return err
		}
		transaction.journal.RollbackPhase = transaction.journal.Phase
		transaction.journal.RollbackCursor = transaction.journal.Cursor
		if err := transaction.persist(
			transactionPhaseRollbackPending,
			len(transaction.journal.Entries),
		); err != nil {
			return fmt.Errorf("begin transaction rollback: %w", err)
		}
	}
	if err := validateRollbackPendingState(transaction.root, transaction.journal); err != nil {
		return err
	}
	for transaction.journal.Cursor > 0 {
		index := transaction.journal.Cursor - 1
		if err := transaction.restoreEntry(index); err != nil {
			return err
		}
		if err := transaction.persist(transactionPhaseRollbackPending, index); err != nil {
			return fmt.Errorf("checkpoint transaction rollback at entry %d: %w", index, err)
		}
	}
	return nil
}

func validateRollbackPendingState(root string, journal transactionJournal) error {
	source := journal
	source.Phase = journal.RollbackPhase
	source.Cursor = journal.RollbackCursor
	source.RollbackPhase = ""
	source.RollbackCursor = 0
	for index, entry := range journal.Entries {
		state, err := inspectTransactionEntryState(root, journal, entry)
		if err != nil {
			return err
		}
		switch {
		case index >= journal.Cursor:
			if !transactionEntryStateIsRestored(entry, state) {
				return fmt.Errorf(
					"rollback-pending target %s changed after interruption",
					entry.Target,
				)
			}
		case index == journal.Cursor-1:
			if !transactionEntryStateIsRestored(entry, state) &&
				!transactionEntryStateIsRollbackTransition(entry, state) &&
				!transactionEntryStateMatchesAny(
					entry,
					state,
					allowedRollbackEntryStates(source, index),
				) {
				return fmt.Errorf(
					"rollback-pending target %s changed after interruption",
					entry.Target,
				)
			}
		default:
			if !transactionEntryStateMatchesAny(
				entry,
				state,
				allowedRollbackEntryStates(source, index),
			) {
				return fmt.Errorf(
					"rollback-pending target %s changed after interruption",
					entry.Target,
				)
			}
		}
	}
	return validateTransactionCreatedDirectoryState(root, source)
}

func transactionEntryStateMatchesAny(
	entry transactionJournalEntry,
	state transactionEntryState,
	allowed []expectedTransactionEntryState,
) bool {
	for _, expected := range allowed {
		if transactionEntryStateMatches(entry, state, expected) {
			return true
		}
	}
	return false
}

func transactionEntryStateIsRestored(
	entry transactionJournalEntry,
	state transactionEntryState,
) bool {
	previous := transactionFileState{
		exists: entry.PreviousExists,
		hash:   entry.PreviousHash,
		mode:   entry.PreviousMode.Perm(),
	}
	intended := transactionFileState{
		exists: true,
		hash:   entry.IntendedHash,
		mode:   entry.IntendedMode.Perm(),
	}
	return transactionFileStatesEqual(state.target, previous) &&
		!state.backup.exists &&
		(!state.staged.exists || transactionFileStatesEqual(state.staged, intended))
}

func transactionEntryStateIsRollbackTransition(
	entry transactionJournalEntry,
	state transactionEntryState,
) bool {
	previous := transactionFileState{
		exists: entry.PreviousExists,
		hash:   entry.PreviousHash,
		mode:   entry.PreviousMode.Perm(),
	}
	intended := transactionFileState{
		exists: true,
		hash:   entry.IntendedHash,
		mode:   entry.IntendedMode.Perm(),
	}
	return !state.target.exists &&
		transactionFileStatesEqual(state.backup, previous) &&
		(!state.staged.exists || transactionFileStatesEqual(state.staged, intended))
}

func (transaction *artifactTransaction) restoreEntry(index int) error {
	entry := transaction.journal.Entries[index]
	state, err := inspectTransactionEntryState(
		transaction.root,
		transaction.journal,
		entry,
	)
	if err != nil {
		return err
	}
	if transactionEntryStateIsRestored(entry, state) {
		return nil
	}
	target, _, backup, err := transaction.paths(entry)
	if err != nil {
		return err
	}
	if state.target.exists {
		intended := transactionFileState{
			exists: true,
			hash:   entry.IntendedHash,
			mode:   entry.IntendedMode.Perm(),
		}
		if !transactionFileStatesEqual(state.target, intended) {
			return fmt.Errorf("transaction target %s changed after interruption", entry.Target)
		}
		if err := os.Remove(target); err != nil {
			return fmt.Errorf("remove replacement for %s: %w", entry.Target, err)
		}
		if err := syncDirectory(filepath.Dir(target)); err != nil {
			return fmt.Errorf("sync replacement removal for %s: %w", entry.Target, err)
		}
		transactionRollbackHook(transactionRollbackPoint{kind: "replacement-removed", index: index})
	}
	if entry.PreviousExists {
		backupState, err := inspectTransactionFile(backup)
		if err != nil {
			return fmt.Errorf("inspect backup for %s: %w", entry.Target, err)
		}
		previous := transactionFileState{
			exists: true,
			hash:   entry.PreviousHash,
			mode:   entry.PreviousMode.Perm(),
		}
		if !transactionFileStatesEqual(backupState, previous) {
			return fmt.Errorf("backup changed after interruption for %s", entry.Target)
		}
		if err := transactionFS.rename(backup, target); err != nil {
			return fmt.Errorf("restore backup for %s: %w", entry.Target, err)
		}
		if err := syncRenameDirectories(backup, target); err != nil {
			return fmt.Errorf("sync backup restoration for %s: %w", entry.Target, err)
		}
		transactionRollbackHook(transactionRollbackPoint{kind: "backup-restored", index: index})
	}
	return nil
}

type transactionDirectoryState int

const (
	transactionDirectoryAbsent transactionDirectoryState = iota
	transactionDirectoryStaged
	transactionDirectoryOwned
	transactionDirectoryRemoved
)

func validateTransactionCreatedDirectoryState(
	root string,
	journal transactionJournal,
) error {
	for index, relative := range journal.CreatedDirectories {
		state, err := inspectTransactionDirectoryState(root, journal, index, false)
		if err != nil {
			return err
		}
		wantStaged := false
		allowStagedOrOwned := false
		switch journal.Phase {
		case transactionPhasePrepared:
			wantStaged = true
		case transactionPhaseDirectoryCreatePending:
			switch {
			case index < journal.Cursor:
			case index == journal.Cursor:
				allowStagedOrOwned = true
			default:
				wantStaged = true
			}
		case transactionPhaseDirectoryCreated:
			wantStaged = index >= journal.Cursor
		default:
		}
		switch {
		case allowStagedOrOwned &&
			(state == transactionDirectoryStaged || state == transactionDirectoryOwned):
		case wantStaged && state == transactionDirectoryStaged:
		case !wantStaged && !allowStagedOrOwned && state == transactionDirectoryOwned:
		default:
			return fmt.Errorf(
				"%s created directory %s changed after interruption",
				journal.Phase,
				relative,
			)
		}
	}
	return nil
}

func inspectTransactionDirectoryState(
	root string,
	journal transactionJournal,
	index int,
	requireMarkerOnly bool,
) (transactionDirectoryState, error) {
	relative := journal.CreatedDirectories[index]
	target, err := artifactDirectoryPath(root, relative)
	if err != nil {
		return 0, err
	}
	targetExists, targetOwned, err := inspectOwnedDirectory(
		target,
		transactionOwnerPrefix+strings.TrimPrefix(
			journal.Directory,
			transactionDirectoryPrefix,
		),
		journal.Directory+"\n",
		requireMarkerOnly,
	)
	if err != nil {
		return 0, fmt.Errorf("inspect transaction-owned directory %s: %w", relative, err)
	}
	stagedRelative := path.Join(journal.Directory, "directories", fmt.Sprintf("%06d", index))
	staged, err := transactionInternalDirectoryPath(root, journal.Directory, stagedRelative)
	if err != nil {
		return 0, err
	}
	stagedExists, stagedOwned, err := inspectOwnedDirectory(
		staged,
		transactionOwnerPrefix+strings.TrimPrefix(
			journal.Directory,
			transactionDirectoryPrefix,
		),
		journal.Directory+"\n",
		true,
	)
	if err != nil {
		return 0, fmt.Errorf("inspect staged transaction directory %s: %w", relative, err)
	}
	removedRelative := path.Join(journal.Directory, "removed", fmt.Sprintf("%06d", index))
	removed, err := transactionInternalDirectoryPath(root, journal.Directory, removedRelative)
	if err != nil {
		return 0, err
	}
	removedExists, removedOwned, err := inspectOwnedDirectory(
		removed,
		transactionOwnerPrefix+strings.TrimPrefix(
			journal.Directory,
			transactionDirectoryPrefix,
		),
		journal.Directory+"\n",
		true,
	)
	if err != nil {
		return 0, fmt.Errorf("inspect removed transaction directory %s: %w", relative, err)
	}
	switch {
	case boolCount(targetExists, stagedExists, removedExists) > 1:
		return 0, fmt.Errorf("transaction-owned directory %s exists at multiple locations", relative)
	case targetExists && !targetOwned:
		return 0, fmt.Errorf("directory %s has no valid transaction ownership marker", relative)
	case stagedExists && !stagedOwned:
		return 0, fmt.Errorf("staged directory %s has no valid transaction ownership marker", relative)
	case removedExists && !removedOwned:
		return 0, fmt.Errorf("removed directory %s has no valid transaction ownership marker", relative)
	case targetExists:
		return transactionDirectoryOwned, nil
	case stagedExists:
		return transactionDirectoryStaged, nil
	case removedExists:
		return transactionDirectoryRemoved, nil
	default:
		return transactionDirectoryAbsent, nil
	}
}

func boolCount(values ...bool) int {
	count := 0
	for _, value := range values {
		if value {
			count++
		}
	}
	return count
}

func inspectOwnedDirectory(
	directory, markerName, markerContent string,
	requireMarkerOnly bool,
) (exists, owned bool, err error) {
	info, err := os.Lstat(directory)
	if os.IsNotExist(err) {
		return false, false, nil
	}
	if err != nil {
		return false, false, err
	}
	if !info.IsDir() {
		return true, false, fmt.Errorf("is not a directory")
	}
	content, err := os.ReadFile(filepath.Join(directory, markerName))
	if err != nil {
		if os.IsNotExist(err) {
			return true, false, nil
		}
		return true, false, err
	}
	if string(content) != markerContent {
		return true, false, nil
	}
	if requireMarkerOnly {
		entries, err := os.ReadDir(directory)
		if err != nil {
			return true, false, err
		}
		if len(entries) != 1 || entries[0].Name() != markerName {
			return true, false, fmt.Errorf("contains files not owned by the transaction")
		}
	}
	return true, true, nil
}

func (transaction *artifactTransaction) rollbackDirectories() error {
	if len(transaction.journal.CreatedDirectories) == 0 {
		return nil
	}
	if transaction.journal.Phase != transactionPhaseRollbackDirectoriesPending {
		if err := validateRestoredJournalState(transaction.root, transaction.journal); err != nil {
			return err
		}
		source := transaction.journal
		source.Phase = transaction.journal.RollbackPhase
		source.Cursor = transaction.journal.RollbackCursor
		source.RollbackPhase = ""
		source.RollbackCursor = 0
		if err := validateTransactionCreatedDirectoryState(transaction.root, source); err != nil {
			return err
		}
		if err := transaction.persist(
			transactionPhaseRollbackDirectoriesPending,
			len(transaction.journal.CreatedDirectories),
		); err != nil {
			return fmt.Errorf("begin transaction directory rollback: %w", err)
		}
	}
	if err := transaction.validateRollbackDirectoriesState(); err != nil {
		return err
	}
	for transaction.journal.Cursor > 0 {
		index := transaction.journal.Cursor - 1
		state, err := inspectTransactionDirectoryState(
			transaction.root,
			transaction.journal,
			index,
			true,
		)
		if err != nil {
			return err
		}
		if state == transactionDirectoryOwned {
			relative := transaction.journal.CreatedDirectories[index]
			target, err := artifactDirectoryPath(transaction.root, relative)
			if err != nil {
				return err
			}
			removedRelative := path.Join(
				transaction.journal.Directory,
				"removed",
				fmt.Sprintf("%06d", index),
			)
			removed, err := transactionInternalDirectoryPath(
				transaction.root,
				transaction.journal.Directory,
				removedRelative,
			)
			if err != nil {
				return err
			}
			if err := transactionFS.rename(target, removed); err != nil {
				return fmt.Errorf("remove transaction-created directory %s: %w", relative, err)
			}
			if err := syncRenameDirectories(target, removed); err != nil {
				return fmt.Errorf("sync removal of transaction-created directory %s: %w", relative, err)
			}
			transactionRollbackHook(transactionRollbackPoint{kind: "directory-removed", index: index})
		}
		if err := transaction.persist(
			transactionPhaseRollbackDirectoriesPending,
			index,
		); err != nil {
			return fmt.Errorf("checkpoint directory rollback at %d: %w", index, err)
		}
	}
	return nil
}

func (transaction *artifactTransaction) validateRollbackDirectoriesState() error {
	source := transaction.journal
	source.Phase = transaction.journal.RollbackPhase
	source.Cursor = transaction.journal.RollbackCursor
	source.RollbackPhase = ""
	source.RollbackCursor = 0
	for index, relative := range transaction.journal.CreatedDirectories {
		state, err := inspectTransactionDirectoryState(
			transaction.root,
			transaction.journal,
			index,
			stateShouldBeOwnedOnly(index, transaction.journal.Cursor),
		)
		if err != nil {
			return err
		}
		mustExist, mayExist := sourceDirectoryOwnership(source, index)
		processed := index >= transaction.journal.Cursor
		current := index == transaction.journal.Cursor-1
		valid := false
		switch {
		case processed && mustExist:
			valid = state == transactionDirectoryRemoved || state == transactionDirectoryAbsent
		case processed:
			valid = state == transactionDirectoryStaged ||
				state == transactionDirectoryAbsent ||
				mayExist && state == transactionDirectoryRemoved
		case current && mustExist:
			valid = state == transactionDirectoryOwned ||
				state == transactionDirectoryRemoved
		case current && mayExist:
			valid = state == transactionDirectoryStaged ||
				state == transactionDirectoryOwned ||
				state == transactionDirectoryRemoved
		case mustExist:
			valid = state == transactionDirectoryOwned
		case mayExist:
			valid = state == transactionDirectoryStaged || state == transactionDirectoryOwned
		default:
			valid = state == transactionDirectoryStaged
		}
		if !valid {
			return fmt.Errorf(
				"rollback-created directory %s changed after interruption",
				relative,
			)
		}
	}
	return nil
}

func sourceDirectoryOwnership(
	journal transactionJournal,
	index int,
) (mustExist, mayExist bool) {
	switch journal.Phase {
	case transactionPhasePrepared:
		return false, false
	case transactionPhaseDirectoryCreatePending:
		switch {
		case index < journal.Cursor:
			return true, true
		case index == journal.Cursor:
			return false, true
		default:
			return false, false
		}
	case transactionPhaseDirectoryCreated:
		if index < journal.Cursor {
			return true, true
		}
		return false, false
	default:
		return true, true
	}
}

func stateShouldBeOwnedOnly(index, cursor int) bool {
	return index == cursor-1
}

func cleanupTransaction(root, journalPath string, journal transactionJournal) error {
	if err := discardRemovedDirectories(root, journal); err != nil {
		return err
	}
	transactionDirectory := filepath.Join(root, journal.Directory)
	if err := os.RemoveAll(transactionDirectory); err != nil {
		return fmt.Errorf("remove transaction directory: %w", err)
	}
	if err := syncDirectory(root); err != nil {
		return fmt.Errorf("sync transaction directory cleanup: %w", err)
	}
	transactionCleanupHook()
	if err := os.Remove(journalPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove transaction journal: %w", err)
	}
	if err := syncDirectory(root); err != nil {
		return fmt.Errorf("sync transaction journal cleanup: %w", err)
	}
	return nil
}

func discardRemovedDirectories(root string, journal transactionJournal) error {
	markerName := transactionOwnerPrefix +
		strings.TrimPrefix(journal.Directory, transactionDirectoryPrefix)
	for index, relative := range journal.CreatedDirectories {
		removedRelative := path.Join(
			journal.Directory,
			"removed",
			fmt.Sprintf("%06d", index),
		)
		removed, err := transactionInternalDirectoryPath(
			root,
			journal.Directory,
			removedRelative,
		)
		if err != nil {
			return err
		}
		exists, owned, err := inspectOwnedDirectory(
			removed,
			markerName,
			journal.Directory+"\n",
			true,
		)
		if err != nil {
			return fmt.Errorf("validate removed directory %s before cleanup: %w", relative, err)
		}
		if !exists {
			continue
		}
		if !owned {
			return fmt.Errorf("removed directory %s changed before cleanup", relative)
		}
		if err := os.Remove(filepath.Join(removed, markerName)); err != nil {
			return fmt.Errorf("remove ownership marker for restored directory %s: %w", relative, err)
		}
		if err := syncDirectory(removed); err != nil {
			return fmt.Errorf("sync restored directory marker removal for %s: %w", relative, err)
		}
		if err := os.Remove(removed); err != nil {
			return fmt.Errorf("remove restored directory %s safely: %w", relative, err)
		}
		if err := syncDirectory(filepath.Dir(removed)); err != nil {
			return fmt.Errorf("sync safe restored directory removal for %s: %w", relative, err)
		}
	}
	return nil
}

func syncRenameDirectories(oldPath, newPath string) error {
	oldDirectory := filepath.Dir(oldPath)
	newDirectory := filepath.Dir(newPath)
	if err := syncDirectory(oldDirectory); err != nil {
		return err
	}
	if newDirectory != oldDirectory {
		if err := syncDirectory(newDirectory); err != nil {
			return err
		}
	}
	return nil
}

func syncDirectory(directory string) error {
	return transactionFS.syncDirectory(directory)
}

func syncOSDirectory(directory string) error {
	file, err := os.Open(directory)
	if err != nil {
		return err
	}
	defer file.Close()
	return file.Sync()
}
