package flagConfig

// FlagConfig stores configuration options
type FlagConfig struct {
	ConfigFile    string
	NoGitignore   bool
	OutputFile    string
	DisplayLineNum bool
	Verbose       bool
}
