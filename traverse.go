package main

import (
	"io/fs"
	"path/filepath"

	"kau/services"
)

func traverseRepo(root string) ([]services.ProjectTask, error) {
	var tasks []services.ProjectTask
	providers := services.AllProviders()

	// 1. Collect ignored directories from all registered language providers
	ignored := map[string]bool{
		".git": true,
		".svn": true,
		".hg":  true,
	}
	for _, p := range providers {
		for _, dir := range p.IgnoredDirs() {
			ignored[dir] = true
		}
	}

	// 2. Walk the repository tree
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}

		// Skip ignored directories (e.g., node_modules, target, .git)
		if ignored[d.Name()] {
			return filepath.SkipDir
		}

		// 3. Hand the directory to every registered language specialist
		for _, p := range providers {
			if task := p.Detect(path); task != nil {
				tasks = append(tasks, *task)
			}
		}

		return nil
	})

	return tasks, err
}
