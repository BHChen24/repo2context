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

func initConfig() {
	// Determine config file location
	if flagCfg.ConfigFile != "" {
		viper.SetConfigFile(flagCfg.ConfigFile)
	} else {
		cwd, err := os.Getwd()
		cobra.CheckErr(err)

		// Look for specific dotfile TOML in cwd
		configFile := filepath.Join(cwd, ".r2c-config.toml")
		if _, err := os.Stat(configFile); err == nil {
			viper.SetConfigFile(configFile)
			viper.SetConfigType("toml")
		} else {
			// fallback to home YAML config
			home, err := os.UserHomeDir()
			cobra.CheckErr(err)
			viper.AddConfigPath(home)
			viper.SetConfigType("yaml")
			viper.SetConfigName(".repo2context")
		}
	}

	viper.AutomaticEnv() // read env vars

	// Read config with enhanced error handling
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok && flagCfg.ConfigFile == "" {
			// No config file found, but none was explicitly specified - this is OK
		} else {
			// Config file was specified or another error occurred
			fmt.Fprintf(os.Stderr, "Error reading config file: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	// Unmarshal into flagCfg (CLI flags already bound; CLI overrides TOML automatically)
	if err := viper.Unmarshal(&flagCfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config file: %v\n", err)
		os.Exit(1)
	}
}
