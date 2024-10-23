/*
Copyright © 2024 TaskMasterErnest
Copyrights apply to this source code.
Check LICENSE for details
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pscan",
	Short: "Fast TCP port scanner",
	Long: `pscan - short for Port Scanner - executes TCP port scan on a list of hosts.
	
	pscan allows you to add, list and delete hosts from the list.
	
	pscan executes a port scan on specified TCP ports. 
	You can customize the target ports using a command-line flag`,
	Version: "1.0",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func initConfig() {
	// use config file from file
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// find the home directory
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// search for config in home dir with name .pscan
		viper.AddConfigPath(home)
		viper.SetConfigName(".pscan")
	}

	// read environment variables that match
	viper.AutomaticEnv()

	// if config file is found, read it in
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Using config file: ", viper.ConfigFileUsed())
	}
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pscan.yaml)")

	rootCmd.PersistentFlags().StringP("hosts-file", "f", "pscan.hosts", "pscan hosts file")

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("PSCAN")

	viper.BindPFlag("hosts-file", rootCmd.PersistentFlags().Lookup("hosts-file"))

	versionTemplate := `{{ printf "%s: %s - version %s\n" .Name .Short .Version }}`
	rootCmd.SetVersionTemplate(versionTemplate)
}
