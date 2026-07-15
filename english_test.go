package main

import (
	"os"
	"path/filepath"
	"testing"
	"unicode"
	"unicode/utf8"
)

func TestProjectUsesEnglishScripts(t *testing.T) {
	err := filepath.WalkDir(".", func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if entry.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		if containsJapaneseText(path) {
			t.Errorf("path contains a Japanese character: %s", path)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if utf8.Valid(content) && containsJapaneseText(string(content)) {
			t.Errorf("file contains a Japanese character: %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func containsJapaneseText(value string) bool {
	for _, character := range value {
		if unicode.In(character, unicode.Han, unicode.Hiragana, unicode.Katakana) ||
			character >= 0x3000 && character <= 0x303f ||
			character >= 0xff00 && character <= 0xffef {
			return true
		}
	}
	return false
}
