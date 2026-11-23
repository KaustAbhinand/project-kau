package main

import (
	"fmt"
	"os"
	"path/filepath"

	"kau/services"
)

func runConfigurations(tasks []services.ProjectTask) error {
	for _, task := range tasks {
		if len(task.ConfigTemplates) == 0 {
			continue
		}

		fmt.Printf("Configuring %s in %s...\n", task.Ecosystem, task.RootDir)

		for srcRel, dstRel := range task.ConfigTemplates {
			srcPath := filepath.Join(task.RootDir, srcRel)
			dstPath := filepath.Join(task.RootDir, dstRel)

			if err := generateConfigFile(srcPath, dstPath); err != nil {
				return fmt.Errorf("error configuring %s (%s): %w", task.Ecosystem, dstRel, err)
			}
		}
	}
	return nil
}

func generateConfigFile(srcPath, dstPath string) error {
	// Skip if target config already exists so we never clobber existing local settings
	if _, err := os.Stat(dstPath); err == nil {
		fmt.Printf("  Skipping %s (already exists)\n", filepath.Base(dstPath))
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	data, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", srcPath, err)
	}

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", dstPath, err)
	}

	if err := os.WriteFile(dstPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config %s: %w", dstPath, err)
	}

	fmt.Printf("  Created %s from %s\n", filepath.Base(dstPath), filepath.Base(srcPath))
	return nil
}
