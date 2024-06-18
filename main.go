package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v2"
)

var (
	settings Settings
)

type Rule struct {
	Path    string `yaml:"path"`
	Pattern string `yaml:"pattern"`
	Days    int    `yaml:"days"`
}

type Settings struct {
	Rules   []Rule   `yaml:"rules"`
	Exclude []string `yaml:"exclude"`
	Root    string   `yaml:"root"`
	Trash   string   `yaml:"trash"`
}

func main() {
	getSettings()
	runRules()
}

func getSettings() {
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalf("error getting executable path: %v", err)
	}
	executableDir := filepath.Dir(executablePath)
	settingsPath := filepath.Join(executableDir, "settings.yaml")

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = yaml.Unmarshal(data, &settings)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func runRules() {
	for _, rule := range settings.Rules {
		r, err := regexp.Compile(rule.Pattern)
		if err != nil {
			log.Fatalf("error compiling regex: %v", err)
		}

		var folder string = filepath.Join(settings.Root, rule.Path)

		items, err := os.ReadDir(folder)

		if err != nil {
			log.Fatalf("error reading folder: %v", err)
		}

		for _, item := range items {
			var matchesRuleRegex bool = r.MatchString(item.Name())
			var notExcluded bool = !isExcluded(folder, item.Name())
			var isOldEnough bool = oldEnough(item, rule.Days)

			if matchesRuleRegex && notExcluded && isOldEnough {
				deleteFile(folder, item.Name())
			}
		}
	}
}

func isExcluded(folder string, name string) bool {
	var path string = filepath.Join(folder, name)

	for _, exclude := range settings.Exclude {
		r, err := regexp.Compile(exclude)

		if err != nil {
			log.Fatalf("error compiling regex: %v", err)
		}

		if r.MatchString(path) {
			return true
		}
	}

	return false
}

func oldEnough(item fs.DirEntry, days int) bool {
	info, err := item.Info()

	if err != nil {
		log.Fatalf("error getting file info: %v", err)
	}

	modTime := info.ModTime()
	limit := modTime.AddDate(0, 0, -days)

	return modTime.Before(limit)
}

func deleteFile(folder string, name string) {
	var oldPath string = filepath.Join(folder, name)
	var newPath string = filepath.Join(settings.Trash, name)

	fmt.Println(oldPath, " -> ", newPath)

	os.Rename(oldPath, newPath)
}
