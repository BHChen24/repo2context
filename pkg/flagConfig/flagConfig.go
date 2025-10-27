package flagConfig

// FlagConfig stores configuration options
type FlagConfig struct {
	ConfigFile     string `mapstructure:"config"`
	NoGitignore    bool   `mapstructure:"no_gitignore"`
	OutputFile      string `mapstructure:"output"`
	DisplayLineNum  bool   `mapstructure:"display_line_num"`
	Verbose         bool   `mapstructure:"verbose"`
	CountTokens     bool   `mapstructure:"count_tokens"`
}
