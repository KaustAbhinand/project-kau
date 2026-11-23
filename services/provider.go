package services

// RuntimeSpec defines the global toolchain required on the host machine.
type RuntimeSpec struct {
	Name        string   // Human-readable name (e.g., "Node.js", "Java")
	CheckCmd    string   // Binary to verify installation (e.g., "node", "java")
	CheckArgs   []string // Arguments to check version (e.g., ["-version"])
	WinWingetID string   // Optional: Windows winget package ID
}

// ProjectTask represents a discovered subproject ready for installation and configuration.
type ProjectTask struct {
	Ecosystem       string            // e.g., "Node.js", "Java (Maven)", "Java (Gradle)"
	RootDir         string            // Directory where install commands and config templates run
	Runtime         RuntimeSpec       // Required global runtime for this task
	InstallSteps    [][]string        // Ordered commands, e.g., [["npm", "install"]]
	ConfigTemplates map[string]string // Map of template source -> target file (e.g., ".env.example": ".env")
}

// LanguageProvider is the contract every language plugin in services/ must implement.
type LanguageProvider interface {
	// Name returns the display name of the language/ecosystem.
	Name() string

	// IgnoredDirs returns directory names that traverse.go should prune (e.g., "node_modules", "target").
	IgnoredDirs() []string

	// Detect inspects dir. If a matching project manifest is found, it returns a configured ProjectTask; otherwise nil.
	Detect(dir string) *ProjectTask
}

var registry []LanguageProvider

// Register is called automatically inside each language file's init() function.
func Register(p LanguageProvider) {
	registry = append(registry, p)
}

// AllProviders returns every registered language provider.
func AllProviders() []LanguageProvider {
	return registry
}
