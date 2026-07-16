package skillinventory

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/game-dev-rta-club/gh-skill-linker/internal/status"
)

func Merge(projectRoot string, managed []status.Record, discovered []Entry) Result {
	result := Result{Entries: make([]Entry, 0, len(managed)+len(discovered))}
	byPath := make(map[string]int, len(managed)+len(discovered))
	for index := range managed {
		record := &managed[index]
		absolute := record.Path
		if !filepath.IsAbs(absolute) {
			absolute = filepath.Join(projectRoot, filepath.FromSlash(record.Path))
		}
		absolute, _ = filepath.Abs(filepath.Clean(absolute))
		entry := Entry{
			SkillName: record.SkillName,
			Path:      displayPath(projectRoot, absolute), AbsolutePath: absolute,
			Scope: ScopeProject, Provider: ProviderSkillLinker,
			Source: managedSource(record.SourceURL, record.SourceRef),
			Status: "unknown", Managed: record,
		}
		if record.State != nil {
			entry.Status = string(*record.State)
		}
		byPath[absolute] = len(result.Entries)
		result.Entries = append(result.Entries, entry)
	}
	for _, candidate := range discovered {
		absolute, err := filepath.Abs(filepath.Clean(candidate.AbsolutePath))
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("skip invalid skill path %q", candidate.AbsolutePath))
			continue
		}
		candidate.AbsolutePath = absolute
		candidate.Path = displayPath(projectRoot, absolute)
		if existingIndex, ok := byPath[absolute]; ok {
			existing := &result.Entries[existingIndex]
			if candidate.Provider == ProviderLocal || candidate.Provider == existing.Provider {
				continue
			}
			existing.Provider = ProviderConflict
			existing.Status = "provider-conflict"
			existing.Source = joinSources(existing.Source, candidate.Source)
			result.Warnings = append(
				result.Warnings,
				fmt.Sprintf("%s has multiple providers", existing.Path),
			)
			continue
		}
		byPath[absolute] = len(result.Entries)
		result.Entries = append(result.Entries, candidate)
	}
	sort.Slice(result.Entries, func(i, j int) bool {
		if result.Entries[i].Path == result.Entries[j].Path {
			return result.Entries[i].SkillName < result.Entries[j].SkillName
		}
		return result.Entries[i].Path < result.Entries[j].Path
	})
	sort.Strings(result.Warnings)
	return result
}

func displayPath(projectRoot, absolute string) string {
	if pathContained(projectRoot, absolute) {
		relative, err := filepath.Rel(projectRoot, absolute)
		if err == nil {
			return filepath.ToSlash(relative)
		}
	}
	return filepath.Clean(absolute)
}

func managedSource(repository, ref *string) string {
	if repository == nil {
		return ""
	}
	value := strings.TrimSuffix(*repository, ".git")
	value = strings.TrimPrefix(value, "https://github.com/")
	if ref == nil || *ref == "" {
		return value
	}
	shortRef := strings.TrimPrefix(*ref, "refs/heads/")
	shortRef = strings.TrimPrefix(shortRef, "refs/tags/")
	return value + "@" + shortRef
}

func joinSources(first, second string) string {
	if first == "" {
		return second
	}
	if second == "" || first == second {
		return first
	}
	return first + " | " + second
}
