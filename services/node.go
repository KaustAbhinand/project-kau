package services

import (
	"os"
	"path/filepath"
)

type NodeProvider struct{}

// init auto-registers the Node provider with the global registry on startup
func init() {
	Register(&NodeProvider{})
}

func (n *NodeProvider) Name() string {
	return "Node.js"
}

func (n *NodeProvider) IgnoredDirs() []string {
	return []string{"node_modules"}
}

func (n *NodeProvider) Detect(dir string) *ProjectTask {
	pkgPath := filepath.Join(dir, "package.json")
	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		return nil
	}

	// Determine package manager and install command based on lockfile presence
	installCmd := []string{"npm", "install"}
	if _, err := os.Stat(filepath.Join(dir, "pnpm-lock.yaml")); err == nil {
		installCmd = []string{"pnpm", "install"}
	} else if _, err := os.Stat(filepath.Join(dir, "yarn.lock")); err == nil {
		installCmd = []string{"yarn", "install"}
	}

	// Check for environment configuration templates
	configTemplates := make(map[string]string)
	if _, err := os.Stat(filepath.Join(dir, ".env.example")); err == nil {
		configTemplates[".env.example"] = ".env"
	} else if _, err := os.Stat(filepath.Join(dir, ".env.template")); err == nil {
		configTemplates[".env.template"] = ".env"
	}

	return &ProjectTask{
		Ecosystem: "Node.js",
		RootDir:   dir,
		Runtime: RuntimeSpec{
			Name:        "Node.js",
			CheckCmd:    "node",
			CheckArgs:   []string{"-v"},
			WinWingetID: "OpenJS.NodeJS.LTS",
		},
		InstallSteps:    [][]string{installCmd},
		ConfigTemplates: configTemplates,
	}
}
