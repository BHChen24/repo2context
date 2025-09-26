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

	"github.com/BHChen24/repo2context/pkg/core"
	"github.com/BHChen24/repo2context/pkg/flagConfig"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var flagCfg flagConfig.FlagConfig

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "r2c [flags] path1 path2 ...",
	Short:   "Convert repository context to structured markdown",
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
}

// initConfig reads in config file and ENV variables if set.
// Env config file support 
func initConfig() {
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
}
