package skillinventory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/game-dev-rta-club/gh-skill-linker/internal/skill"
	"github.com/game-dev-rta-club/gh-skill-linker/internal/status"
)

type Scope string

const (
	ScopeProject Scope = "project"
	ScopeUser    Scope = "user"
	ScopeSystem  Scope = "system"
)

type Provider string

const (
	ProviderSkillLinker Provider = "skill-linker"
	ProviderGHSkill     Provider = "gh-skill"
	ProviderCodexPlugin Provider = "codex-plugin"
	ProviderLocal       Provider = "local"
	ProviderCodexSystem Provider = "codex-system"
	ProviderConflict    Provider = "conflict"
)

type Entry struct {
	SkillName    string         `json:"skillName"`
	Path         string         `json:"path"`
	Scope        Scope          `json:"scope"`
	Provider     Provider       `json:"provider"`
	Source       string         `json:"source"`
	Status       string         `json:"status"`
	AbsolutePath string         `json:"-"`
	Managed      *status.Record `json:"-"`
}

type Result struct {
	Entries  []Entry
	Warnings []string
}

type Runner interface {
	Run(ctx context.Context, args ...string) (stdout string, stderr string, err error)
}

type Service struct {
	gh    Runner
	codex Runner
}

func NewService(gh, codex Runner) *Service {
	return &Service{gh: gh, codex: codex}
}

func (s *Service) Inspect(ctx context.Context, projectRoot string) Result {
	result := Result{}
	direct, warning := s.readDirectSkills(ctx)
	result.Entries = append(result.Entries, direct...)
	if warning != "" {
		result.Warnings = append(result.Warnings, warning)
	}
	plugins, warnings := s.readCodexPlugins(ctx, projectRoot)
	result.Entries = append(result.Entries, plugins...)
	result.Warnings = append(result.Warnings, warnings...)
	sort.Slice(result.Entries, func(i, j int) bool {
		if result.Entries[i].AbsolutePath == result.Entries[j].AbsolutePath {
			return result.Entries[i].SkillName < result.Entries[j].SkillName
		}
		return result.Entries[i].AbsolutePath < result.Entries[j].AbsolutePath
	})
	return result
}

type ghSkillRecord struct {
	SkillName string `json:"skillName"`
	SourceURL string `json:"sourceURL"`
	Scope     string `json:"scope"`
	Version   string `json:"version"`
	Path      string `json:"path"`
}

func (s *Service) readDirectSkills(ctx context.Context) ([]Entry, string) {
	stdout, stderr, err := s.gh.Run(
		ctx,
		"skill", "list", "--agent", "codex", "--json",
		"skillName,sourceURL,scope,version,pinned,path,agentHosts",
	)
	if err != nil {
		detail := strings.TrimSpace(stderr)
		if detail == "" {
			detail = err.Error()
		}
		return nil, fmt.Sprintf("gh skill inventory unavailable; update GitHub CLI: %s", detail)
	}
	var records []ghSkillRecord
	if err := json.Unmarshal([]byte(stdout), &records); err != nil {
		return nil, fmt.Sprintf("gh skill inventory returned invalid JSON: %v", err)
	}
	entries := make([]Entry, 0, len(records))
	for _, record := range records {
		entry, ok := directEntry(record)
		if ok {
			entries = append(entries, entry)
		}
	}
	return entries, ""
}

func directEntry(record ghSkillRecord) (Entry, bool) {
	if strings.TrimSpace(record.SkillName) == "" || strings.TrimSpace(record.Path) == "" {
		return Entry{}, false
	}
	absolute, err := filepath.Abs(filepath.Clean(record.Path))
	if err != nil {
		return Entry{}, false
	}
	entry := Entry{
		SkillName: record.SkillName, AbsolutePath: absolute,
		Scope: Scope(record.Scope), Provider: ProviderLocal, Status: "present",
	}
	if record.SourceURL != "" {
		entry.Provider = ProviderGHSkill
		entry.Source = sourceWithVersion(record.SourceURL, record.Version)
	}
	if strings.HasPrefix(record.SkillName, ".system/") {
		entry.SkillName = strings.TrimPrefix(record.SkillName, ".system/")
		entry.Scope = ScopeSystem
		entry.Provider = ProviderCodexSystem
		entry.Source = "Codex"
	}
	if entry.Scope != ScopeProject && entry.Scope != ScopeUser && entry.Scope != ScopeSystem {
		return Entry{}, false
	}
	return entry, true
}

type codexPluginList struct {
	Installed []codexPlugin `json:"installed"`
}

type codexPlugin struct {
	PluginID  string `json:"pluginId"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Installed bool   `json:"installed"`
	Enabled   bool   `json:"enabled"`
	Source    struct {
		Path string `json:"path"`
	} `json:"source"`
}

type pluginManifest struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Skills  string `json:"skills"`
}

func (s *Service) readCodexPlugins(ctx context.Context, projectRoot string) ([]Entry, []string) {
	stdout, stderr, err := s.codex.Run(ctx, "plugin", "list", "--json")
	if err != nil {
		if commandMissing(err) {
			return nil, nil
		}
		detail := strings.TrimSpace(stderr)
		if detail == "" {
			detail = err.Error()
		}
		return nil, []string{"Codex plugin inventory unavailable: " + detail}
	}
	var list codexPluginList
	if err := json.Unmarshal([]byte(stdout), &list); err != nil {
		return nil, []string{fmt.Sprintf("Codex plugin inventory returned invalid JSON: %v", err)}
	}
	entries := []Entry{}
	warnings := []string{}
	for _, plugin := range list.Installed {
		if !plugin.Installed || !plugin.Enabled || plugin.Source.Path == "" {
			continue
		}
		pluginEntries, pluginWarnings, err := readPluginSkills(projectRoot, plugin)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("skip Codex plugin %s: %v", plugin.PluginID, err))
			continue
		}
		entries = append(entries, pluginEntries...)
		warnings = append(warnings, pluginWarnings...)
	}
	return entries, warnings
}

func readPluginSkills(projectRoot string, plugin codexPlugin) ([]Entry, []string, error) {
	root, err := filepath.Abs(filepath.Clean(plugin.Source.Path))
	if err != nil {
		return nil, nil, fmt.Errorf("resolve plugin path: %w", err)
	}
	content, err := os.ReadFile(filepath.Join(root, ".codex-plugin", "plugin.json"))
	if err != nil {
		return nil, nil, fmt.Errorf("read plugin manifest: %w", err)
	}
	var manifest pluginManifest
	if err := json.Unmarshal(content, &manifest); err != nil {
		return nil, nil, fmt.Errorf("parse plugin manifest: %w", err)
	}
	if manifest.Skills == "" {
		return nil, nil, nil
	}
	skillsRoot, err := containedPluginPath(root, manifest.Skills)
	if err != nil {
		return nil, nil, err
	}
	pluginName := manifest.Name
	if pluginName == "" {
		pluginName = plugin.Name
	}
	if pluginName == "" {
		return nil, nil, fmt.Errorf("plugin name is required")
	}
	version := plugin.Version
	if version == "" {
		version = manifest.Version
	}
	source := plugin.PluginID
	if version != "" {
		source += " (" + version + ")"
	}
	scope := ScopeUser
	if pathContained(projectRoot, root) {
		scope = ScopeProject
	}
	entries := []Entry{}
	warnings := []string{}
	err = filepath.WalkDir(skillsRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&os.ModeSymlink != 0 {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.IsDir() || entry.Name() != "SKILL.md" {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		name, err := skill.ParseDeclaredName(content)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("skip Codex plugin skill %s: %v", path, err))
			return nil
		}
		entries = append(entries, Entry{
			SkillName: pluginName + ":" + name, AbsolutePath: filepath.Dir(path),
			Scope: scope, Provider: ProviderCodexPlugin, Source: source, Status: "enabled",
		})
		return nil
	})
	if err != nil {
		return nil, warnings, fmt.Errorf("scan plugin skills: %w", err)
	}
	return entries, warnings, nil
}

func containedPluginPath(root, relative string) (string, error) {
	clean := filepath.Clean(relative)
	if filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("plugin skills path escapes plugin root")
	}
	target := filepath.Join(root, clean)
	if !pathContained(root, target) {
		return "", fmt.Errorf("plugin skills path escapes plugin root")
	}
	return target, nil
}

func pathContained(root, target string) bool {
	root, err := filepath.Abs(filepath.Clean(root))
	if err != nil {
		return false
	}
	target, err = filepath.Abs(filepath.Clean(target))
	if err != nil {
		return false
	}
	relative, err := filepath.Rel(root, target)
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) && !filepath.IsAbs(relative)
}

func commandMissing(err error) bool {
	return errors.Is(err, exec.ErrNotFound)
}

func sourceWithVersion(sourceURL, version string) string {
	sourceURL = strings.TrimSuffix(sourceURL, ".git")
	sourceURL = strings.TrimPrefix(sourceURL, "https://github.com/")
	if version == "" {
		return sourceURL
	}
	return sourceURL + "@" + version
}
