// Copyright © 2016 Pablo de la Concepcion <pconcepcion@gmail.com>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
//    may be used to endorse or promote products derived from this software
//    without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	// Verbose defines if the output log should be verbose
	Verbose bool

	// Storage holds the backend storage access string
	Storage string

	// Program version
	version = "0.1.0"

	// BuildDate to show in the version command, should be set at build time (see Makefile)
	BuildDate = "unknown"

	// CommitHash to show in the version command, should be set at build time (see Makefile)
	CommitHash = "unknown"
)

// Flag names
const (
	FlagConfig  = "config"
	FlagVerbose = "verbose"
	FlagStorage = "storage"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "telegram_dice_bot",
	Short: "A Telegram bot to roll dices",
	Long: `The bot can roll dice expressios like 3d6 or 5d6k3 to roll 5 dices of six sides and keep the highest 3
  More information and samples on dice gramar on http://pconcepcion.github.io/dice/`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	Version: fmt.Sprintf("version: %s\ncommmit hash: %s\nBuild date: %s\n", version, CommitHash, BuildDate),
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, FlagConfig, "", "config file (default is $HOME/.telegram_dice_bot.yaml)")
	RootCmd.PersistentFlags().BoolVarP(&Verbose, FlagVerbose, "v", false, "verbose output")
	viper.BindPFlag(FlagStorage, RootCmd.PersistentFlags().Lookup(FlagVerbose))

	// storage
	RootCmd.PersistentFlags().StringVarP(&Storage, FlagStorage, "s", "", "Storage backend")
	viper.BindPFlag(FlagStorage, RootCmd.PersistentFlags().Lookup(FlagStorage))
	viper.SetDefault(FlagStorage, "sqlite://telegram_dice_bot.sqlite")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		fmt.Println("Setting config file", cfgFile)
		viper.SetConfigName(cfgFile) // name of config file (without extension)
	} else {
		viper.SetConfigName(".telegram_dice_bot") // name of config file (without extension)
	}
	viper.SetConfigType("yaml")                     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/telegram_dice_bot/")  // path to look for the config file in
	viper.AddConfigPath("$HOME")                    // adding home directory as first search path
	viper.AddConfigPath("$HOME/.telegram_dice_bot") // call multiple times to add many search paths
	viper.AddConfigPath(".")                        // optionally look for config in the working directory
	viper.SetEnvPrefix("tdb")                       // Env variables will start with TDB_
	viper.AutomaticEnv()                            // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		if Verbose {
			fmt.Println(err)
		} else {
			fmt.Println("No config file found")
		}
	}
}
