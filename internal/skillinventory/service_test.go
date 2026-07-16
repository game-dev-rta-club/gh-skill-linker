package skillinventory

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/game-dev-rta-club/gh-skill-linker/internal/status"
	"github.com/game-dev-rta-club/gh-skill-linker/internal/syncstate"
)

type fakeRunner struct {
	stdout string
	stderr string
	err    error
	args   []string
}

func (f *fakeRunner) Run(_ context.Context, args ...string) (string, string, error) {
	f.args = append([]string(nil), args...)
	return f.stdout, f.stderr, f.err
}

func TestInspectClassifiesDirectSkillsFromGHInventory(t *testing.T) {
	root := t.TempDir()
	gh := &fakeRunner{stdout: `[
  {"skillName":"review","sourceURL":"","scope":"project","version":"","path":"` + filepath.Join(root, ".agents/skills/review") + `"},
  {"skillName":"doc-master/repo-health","sourceURL":"https://github.com/example/skills","scope":"project","version":"abc123","path":"` + filepath.Join(root, ".agents/skills/repo-health") + `"},
  {"skillName":".system/imagegen","sourceURL":"","scope":"user","version":"","path":"/home/test/.codex/skills/.system/imagegen"}
]`}
	codex := &fakeRunner{stdout: `{"installed":[],"available":[]}`}

	result := NewService(gh, codex).Inspect(context.Background(), root)

	if len(result.Warnings) != 0 {
		t.Fatalf("warnings = %v, want none", result.Warnings)
	}
	if len(result.Entries) != 3 {
		t.Fatalf("entries = %#v, want 3", result.Entries)
	}
	assertEntry(t, result.Entries, "review", ScopeProject, ProviderLocal, "present")
	assertEntry(t, result.Entries, "doc-master/repo-health", ScopeProject, ProviderGHSkill, "present")
	entry := assertEntry(t, result.Entries, "imagegen", ScopeSystem, ProviderCodexSystem, "present")
	if !strings.Contains(entry.Source, "Codex") {
		t.Fatalf("system source = %q, want Codex", entry.Source)
	}
	wantArgs := []string{
		"skill", "list", "--agent", "codex", "--json",
		"skillName,sourceURL,scope,version,pinned,path,agentHosts",
	}
	if strings.Join(gh.args, "\x00") != strings.Join(wantArgs, "\x00") {
		t.Fatalf("gh args = %q, want %q", gh.args, wantArgs)
	}
}

func TestInspectExpandsEnabledCodexPluginSkills(t *testing.T) {
	root := t.TempDir()
	pluginRoot := filepath.Join(t.TempDir(), "sample-plugin")
	skillRoot := filepath.Join(pluginRoot, "skills", "hello")
	invalidRoot := filepath.Join(pluginRoot, "skills", "invalid")
	if err := os.MkdirAll(filepath.Join(pluginRoot, ".codex-plugin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(skillRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(invalidRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := `{"name":"sample","version":"1.2.3","skills":"./skills/"}`
	if err := os.WriteFile(filepath.Join(pluginRoot, ".codex-plugin", "plugin.json"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	document := "---\nname: Hello\ndescription: Say hello.\n---\n\n# Hello\n"
	if err := os.WriteFile(filepath.Join(skillRoot, "SKILL.md"), []byte(document), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(invalidRoot, "SKILL.md"), []byte("invalid"), 0o644); err != nil {
		t.Fatal(err)
	}
	gh := &fakeRunner{stdout: `[]`}
	codex := &fakeRunner{stdout: `{"installed":[{
  "pluginId":"sample@personal","name":"sample","version":"1.2.3",
  "installed":true,"enabled":true,"source":{"source":"local","path":"` + pluginRoot + `"}
}],"available":[]}`}

	result := NewService(gh, codex).Inspect(context.Background(), root)

	if len(result.Warnings) != 1 || !strings.Contains(result.Warnings[0], "invalid/SKILL.md") {
		t.Fatalf("warnings = %v, want invalid skill warning", result.Warnings)
	}
	entry := assertEntry(t, result.Entries, "sample:Hello", ScopeUser, ProviderCodexPlugin, "enabled")
	if entry.AbsolutePath != skillRoot {
		t.Fatalf("plugin path = %q, want %q", entry.AbsolutePath, skillRoot)
	}
	if entry.Source != "sample@personal (1.2.3)" {
		t.Fatalf("plugin source = %q", entry.Source)
	}
}

func TestInspectIgnoresMissingCodexCommand(t *testing.T) {
	result := NewService(
		&fakeRunner{stdout: `[]`},
		&fakeRunner{err: &exec.Error{Name: "codex", Err: exec.ErrNotFound}},
	).Inspect(context.Background(), t.TempDir())

	if len(result.Warnings) != 0 {
		t.Fatalf("warnings = %v, want none", result.Warnings)
	}
}

func TestInspectWarnsWhenCodexCommandCannotExecute(t *testing.T) {
	result := NewService(
		&fakeRunner{stdout: `[]`},
		&fakeRunner{stderr: "permission denied", err: &exec.Error{Name: "codex", Err: os.ErrPermission}},
	).Inspect(context.Background(), t.TempDir())

	if len(result.Warnings) != 1 || !strings.Contains(result.Warnings[0], "permission denied") {
		t.Fatalf("warnings = %v", result.Warnings)
	}
}

func TestInspectWarnsWhenGHInventoryIsUnavailable(t *testing.T) {
	result := NewService(
		&fakeRunner{stderr: "unknown command skill", err: errors.New("exit status 1")},
		&fakeRunner{stdout: `{"installed":[],"available":[]}`},
	).Inspect(context.Background(), t.TempDir())

	if len(result.Warnings) != 1 || !strings.Contains(result.Warnings[0], "update GitHub CLI") {
		t.Fatalf("warnings = %v", result.Warnings)
	}
}

func TestMergeKeepsManagedHealthAndDeduplicatesLocalDiscovery(t *testing.T) {
	root := t.TempDir()
	state := syncstate.Push
	managed := []status.Record{{
		SkillName: "review", Path: ".agents/skills/review",
		SourceURL: stringPointer("https://github.com/example/skills.git"),
		SourceRef: stringPointer("refs/heads/main"), State: &state,
		PullEligibility: status.Eligible, PushEligibility: status.Eligible,
	}}
	discovered := []Entry{{
		SkillName: "review", AbsolutePath: filepath.Join(root, ".agents/skills/review"),
		Scope: ScopeProject, Provider: ProviderLocal, Status: "present",
	}}

	result := Merge(root, managed, discovered)

	if len(result.Entries) != 1 {
		t.Fatalf("entries = %#v, want one", result.Entries)
	}
	entry := result.Entries[0]
	if entry.Provider != ProviderSkillLinker || entry.Status != "push" || entry.Managed == nil {
		t.Fatalf("entry = %#v", entry)
	}
	if len(result.Warnings) != 0 {
		t.Fatalf("warnings = %v, want none", result.Warnings)
	}
}

func TestMergeReportsConflictingProviderClaimsForSamePath(t *testing.T) {
	root := t.TempDir()
	state := syncstate.Clean
	managed := []status.Record{{
		SkillName: "review", Path: ".agents/skills/review",
		SourceURL: stringPointer("https://github.com/example/skills.git"),
		SourceRef: stringPointer("refs/heads/main"), State: &state,
	}}
	discovered := []Entry{{
		SkillName: "review", AbsolutePath: filepath.Join(root, ".agents/skills/review"),
		Scope: ScopeProject, Provider: ProviderGHSkill, Source: "other/repo", Status: "present",
	}}

	result := Merge(root, managed, discovered)

	if len(result.Entries) != 1 || result.Entries[0].Provider != ProviderConflict {
		t.Fatalf("entries = %#v", result.Entries)
	}
	if len(result.Warnings) != 1 || !strings.Contains(result.Warnings[0], "multiple providers") {
		t.Fatalf("warnings = %v", result.Warnings)
	}
}

func assertEntry(
	t *testing.T, entries []Entry, name string, scope Scope, provider Provider, status string,
) Entry {
	t.Helper()
	for _, entry := range entries {
		if entry.SkillName == name {
			if entry.Scope != scope || entry.Provider != provider || entry.Status != status {
				t.Fatalf("entry %q = %#v", name, entry)
			}
			return entry
		}
	}
	t.Fatalf("entry %q not found in %#v", name, entries)
	return Entry{}
}

func stringPointer(value string) *string { return &value }
