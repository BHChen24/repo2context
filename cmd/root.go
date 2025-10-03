/*
Copyright © 2025 Baihua Chen <bchen102@myseneca.ca>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BHChen24/repo2context/pkg/core"
	"github.com/BHChen24/repo2context/pkg/flagConfig"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flagCfg flagConfig.FlagConfig

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "r2c [flags] path1 path2 ...",
	Short: "Convert repository context to structured markdown",
	Long: `Repo2Context analyzes repository structure and generates comprehensive
markdown documentation with organized sections for easy sharing and analysis.

Features:
- Scans repository structure and file contents
- Generates organized markdown with proper sections
- Respects .gitignore files by default
- Supports file filtering and exclusion`,
	Version: "v0.0.1",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := core.Run(args, flagCfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

/*
func init() {

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Env Config file support for future use (currently unimplemented yet)
	rootCmd.PersistentFlags().StringVar(&flagCfg.ConfigFile, "config", "", "config file (default is $HOME/.repo2context.yaml)")

	// Gitignore control flag
	rootCmd.Flags().BoolVar(&flagCfg.NoGitignore, "no-gitignore", false, "disable automatic .gitignore filtering")

	// Output file flag
	rootCmd.Flags().StringVarP(&flagCfg.OutputFile, "output", "o", "", "save output to file instead of stdout")

	// Line number display flag
	rootCmd.Flags().BoolVarP(&flagCfg.DisplayLineNum, "line-numbers", "l", false, "display line numbers in file contents")

	// Verbose flag
	rootCmd.Flags().BoolVarP(&flagCfg.Verbose, "verbose", "", false, "display verbose output")

	// Bind flags to Viper for config file support
	viper.BindPFlag("no-gitignore", rootCmd.Flags().Lookup("no-gitignore"))
	viper.BindPFlag("output", rootCmd.Flags().Lookup("output"))
	viper.BindPFlag("line-numbers", rootCmd.Flags().Lookup("line-numbers"))
	viper.BindPFlag("verbose", rootCmd.Flags().Lookup("verbose"))

}
*/
func init() {
	cobra.OnInitialize(initConfig)

	// Persistent config file flag
	rootCmd.PersistentFlags().StringVar(&flagCfg.ConfigFile, "config", "", "config file (default is $HOME/.repo2context.yaml)")

	// Other CLI flags
	rootCmd.Flags().BoolVar(&flagCfg.NoGitignore, "no-gitignore", false, "disable automatic .gitignore filtering")
	rootCmd.Flags().StringVarP(&flagCfg.OutputFile, "output", "o", "", "save output to file instead of stdout")
	rootCmd.Flags().BoolVarP(&flagCfg.DisplayLineNum, "line-numbers", "l", false, "display line numbers in file contents")
	rootCmd.Flags().BoolVarP(&flagCfg.Verbose, "verbose", "", false, "display verbose output")

	// Bind flags to Viper
	viper.BindPFlag("no_gitignore", rootCmd.Flags().Lookup("no-gitignore"))
	viper.BindPFlag("output", rootCmd.Flags().Lookup("output"))
	viper.BindPFlag("display_line_num", rootCmd.Flags().Lookup("line-numbers"))
	viper.BindPFlag("verbose", rootCmd.Flags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
// Env config file support
/*func initConfig() {
	if flagCfg.ConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(flagCfg.ConfigFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".repo2context" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".repo2context")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}*/

/*// Now supports TOML config file in current directory
func initConfig() {
	// Set config to look for .repo2context.toml in current directory
	viper.SetConfigName(".repo2context-config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".") // Current directory

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}*/

// initConfig reads in config file and ENV variables if set.
// Env config file support
/*func initConfig() {
	// If a config file was explicitly provided via --config, use it.
	if flagCfg.ConfigFile != "" {
		viper.SetConfigFile(flagCfg.ConfigFile)
	} else {
		// First, look in current working directory for a dotfile TOML config.
		cwd, err := os.Getwd()
		cobra.CheckErr(err)

		// common candidate names
		candidates := []string{
			filepath.Join(cwd, ".r2c.toml"),
			filepath.Join(cwd, ".r2c-config.toml"),
		}

		found := false
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				viper.SetConfigFile(c)
				viper.SetConfigType("toml")
				found = true
				break
			}
		}

		if !found {
			// look for any file named like ".repo2context*.toml"
			entries, err := os.ReadDir(cwd)
			cobra.CheckErr(err)
			for _, e := range entries {
				name := e.Name()
				if strings.HasPrefix(name, ".r2c") && strings.HasSuffix(name, ".toml") {
					viper.SetConfigFile(filepath.Join(cwd, name))
					viper.SetConfigType("toml")
					found = true
					break
				}
			}
		}

		if !found {
			// fallback to searching home (existing behavior, YAML)
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)
			viper.AddConfigPath(home)
			viper.SetConfigType("yaml")
			viper.SetConfigName(".repo2context")
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// Try to read the config. If there's no file, that's fine.
	if err := viper.ReadInConfig(); err != nil {
		// Distinguish between "file not found" and parsing/other errors.
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// No config file found -> ignore (works as required)
		} else {
			// Config file present but could not be read/parsed -> error out.
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Config file found and read
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())

		// Unmarshal config file into a temporary struct
		var fileCfg flagConfig.FlagConfig
		if err := viper.Unmarshal(&fileCfg); err != nil {
			// parsing/unmarshal error -> exit with a helpful message
			fmt.Fprintf(os.Stderr, "Error parsing config file: %v\n", err)
			os.Exit(1)
		}

		// Merge behavior: only apply config values where the CLI flag was NOT explicitly set.
		// That way CLI flags override config values.
		if !rootCmd.Flags().Changed("output") && fileCfg.OutputFile != "" {
			flagCfg.OutputFile = fileCfg.OutputFile
		}
		if !rootCmd.Flags().Changed("no-gitignore") && fileCfg.NoGitignore {
			flagCfg.NoGitignore = fileCfg.NoGitignore
		}
		if !rootCmd.Flags().Changed("line-numbers") && fileCfg.DisplayLineNum {
			flagCfg.DisplayLineNum = fileCfg.DisplayLineNum
		}
		if !rootCmd.Flags().Changed("verbose") && fileCfg.Verbose {
			flagCfg.Verbose = fileCfg.Verbose
		}
	}
}*/

func initConfig() {
	// Determine config file location
	if flagCfg.ConfigFile != "" {
		viper.SetConfigFile(flagCfg.ConfigFile)
	} else {
		cwd, err := os.Getwd()
		cobra.CheckErr(err)

		// Look for dotfile TOML in cwd
		candidates := []string{
			filepath.Join(cwd, ".r2c.toml"),
			filepath.Join(cwd, ".r2c-config.toml"),
		}

		found := false
		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				viper.SetConfigFile(c)
				viper.SetConfigType("toml")
				found = true
				break
			}
		}

		if !found {
			entries, err := os.ReadDir(cwd)
			cobra.CheckErr(err)
			for _, e := range entries {
				name := e.Name()
				if strings.HasPrefix(name, ".r2c") && strings.HasSuffix(name, ".toml") {
					viper.SetConfigFile(filepath.Join(cwd, name))
					viper.SetConfigType("toml")
					found = true
					break
				}
			}
		}

		if !found {
			// fallback to home YAML config
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)
			viper.AddConfigPath(home)
			viper.SetConfigType("yaml")
			viper.SetConfigName(".repo2context")
		}
	}

	viper.AutomaticEnv() // read env vars

	// Read config
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
			os.Exit(1)
		}
	}

	// Unmarshal into flagCfg (CLI flags already bound; CLI overrides TOML automatically)
	if err := viper.Unmarshal(&flagCfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config file: %v\n", err)
		os.Exit(1)
	}
}
