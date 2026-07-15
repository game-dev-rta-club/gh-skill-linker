package workspace

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveSkillDeletesExactDirectoryAfterCommit(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "sample")
	writeTestFile(t, filepath.Join(target, "SKILL.md"), "original\n")
	expected, err := ReadSkill(target)
	if err != nil {
		t.Fatal(err)
	}
	committed := false

	err = RemoveSkill(target, &expected, func() error {
		committed = true
		return nil
	})

	if err != nil {
		t.Fatalf("RemoveSkill() error = %v", err)
	}
	if !committed {
		t.Fatal("RemoveSkill() did not commit metadata")
	}
	if _, err := os.Lstat(target); !os.IsNotExist(err) {
		t.Fatalf("target still exists: %v", err)
	}
}

func TestRemoveSkillRollsBackWhenWorkspaceChanged(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "sample")
	writeTestFile(t, filepath.Join(target, "SKILL.md"), "original\n")
	expected, err := ReadSkill(target)
	if err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(target, "SKILL.md"), "concurrent\n")

	err = RemoveSkill(target, &expected, nil)

	if !errors.Is(err, ErrWorkspaceChanged) {
		t.Fatalf("RemoveSkill() error = %v, want ErrWorkspaceChanged", err)
	}
	content, readErr := os.ReadFile(filepath.Join(target, "SKILL.md"))
	if readErr != nil || string(content) != "concurrent\n" {
		t.Fatalf("target was not restored: content=%q error=%v", content, readErr)
	}
}

func TestRemoveSkillRollsBackWhenEmptyDirectoryWasAdded(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "sample")
	writeTestFile(t, filepath.Join(target, "SKILL.md"), "original\n")
	expected, err := ReadSkill(target)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(target, "notes"), 0o755); err != nil {
		t.Fatal(err)
	}

	err = RemoveSkill(target, &expected, nil)

	if !errors.Is(err, ErrWorkspaceChanged) {
		t.Fatalf("RemoveSkill() error = %v, want ErrWorkspaceChanged", err)
	}
	if _, statErr := os.Stat(filepath.Join(target, "notes")); statErr != nil {
		t.Fatalf("empty directory was not restored: %v", statErr)
	}
}

func TestRemoveSkillRollsBackWhenCommitFails(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "sample")
	writeTestFile(t, filepath.Join(target, "SKILL.md"), "original\n")

	err := RemoveSkill(target, nil, func() error { return errors.New("manifest failed") })

	if err == nil {
		t.Fatal("RemoveSkill() error = nil, want commit failure")
	}
	content, readErr := os.ReadFile(filepath.Join(target, "SKILL.md"))
	if readErr != nil || string(content) != "original\n" {
		t.Fatalf("target was not restored: content=%q error=%v", content, readErr)
	}
}
