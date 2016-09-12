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
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/hugo/helpers"
	"github.com/spf13/hugo/hugofs"
	"github.com/spf13/hugo/utils"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Shows the bot configuration",
	Long: `Shows the bot current configuration as defined on the config file and the
  environment varialbes.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("config called")
	},
}

func init() {
	configCmd.RunE = config
	RootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func config(cmd *cobra.Command, args []string) error {
	if err := initializeConfig(configCmd); err != nil {
		return err
	}

	allSettings := viper.AllSettings()

	var separator string
	if allSettings["metadataformat"] == "toml" {
		separator = " = "
	} else {
		separator = ": "
	}

	var keys []string
	for k := range allSettings {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		kv := reflect.ValueOf(allSettings[k])
		if kv.Kind() == reflect.String {
			fmt.Printf("%s%s\"%+v\"\n", k, separator, allSettings[k])
		} else {
			fmt.Printf("%s%s%+v\n", k, separator, allSettings[k])
		}
	}

	return nil
}

// InitializeConfig initializes a config file with sensible default configuration flags.
func initializeConfig(subCmdVs ...*cobra.Command) error {
	if err := hugolib.LoadGlobalConfig(source, cfgFile); err != nil {
		return err
	}

	for _, cmdV := range append([]*cobra.Command{hugoCmdV}, subCmdVs...) {

		if flagChanged(cmdV.PersistentFlags(), "verbose") {
			viper.Set("Verbose", verbose)
		}
		if flagChanged(cmdV.PersistentFlags(), "logFile") {
			viper.Set("LogFile", logFile)
		}
		if flagChanged(cmdV.Flags(), "cleanDestinationDir") {
			viper.Set("cleanDestinationDir", cleanDestination)
		}
		if flagChanged(cmdV.Flags(), "buildDrafts") {
			viper.Set("BuildDrafts", draft)
		}
		if flagChanged(cmdV.Flags(), "buildFuture") {
			viper.Set("BuildFuture", future)
		}
		if flagChanged(cmdV.Flags(), "buildExpired") {
			viper.Set("BuildExpired", expired)
		}
		if flagChanged(cmdV.Flags(), "uglyURLs") {
			viper.Set("UglyURLs", uglyURLs)
		}
		if flagChanged(cmdV.Flags(), "canonifyURLs") {
			viper.Set("CanonifyURLs", canonifyURLs)
		}
		if flagChanged(cmdV.Flags(), "disable404") {
			viper.Set("Disable404", disable404)
		}
		if flagChanged(cmdV.Flags(), "disableRSS") {
			viper.Set("DisableRSS", disableRSS)
		}
		if flagChanged(cmdV.Flags(), "disableSitemap") {
			viper.Set("DisableSitemap", disableSitemap)
		}
		if flagChanged(cmdV.Flags(), "enableRobotsTXT") {
			viper.Set("EnableRobotsTXT", enableRobotsTXT)
		}
		if flagChanged(cmdV.Flags(), "pluralizeListTitles") {
			viper.Set("PluralizeListTitles", pluralizeListTitles)
		}
		if flagChanged(cmdV.Flags(), "preserveTaxonomyNames") {
			viper.Set("PreserveTaxonomyNames", preserveTaxonomyNames)
		}
		if flagChanged(cmdV.Flags(), "ignoreCache") {
			viper.Set("IgnoreCache", ignoreCache)
		}
		if flagChanged(cmdV.Flags(), "forceSyncStatic") {
			viper.Set("ForceSyncStatic", forceSync)
		}
		if flagChanged(cmdV.Flags(), "noTimes") {
			viper.Set("NoTimes", noTimes)
		}

	}

	if baseURL != "" {
		if !strings.HasSuffix(baseURL, "/") {
			baseURL = baseURL + "/"
		}
		viper.Set("BaseURL", baseURL)
	}

	if !viper.GetBool("RelativeURLs") && viper.GetString("BaseURL") == "" {
		jww.ERROR.Println("No 'baseurl' set in configuration or as a flag. Features like page menus will not work without one.")
	}

	if theme != "" {
		viper.Set("theme", theme)
	}

	if destination != "" {
		viper.Set("PublishDir", destination)
	}

	if source != "" {
		dir, _ := filepath.Abs(source)
		viper.Set("WorkingDir", dir)
	} else {
		dir, _ := os.Getwd()
		viper.Set("WorkingDir", dir)
	}

	if contentDir != "" {
		viper.Set("ContentDir", contentDir)
	}

	if layoutDir != "" {
		viper.Set("LayoutDir", layoutDir)
	}

	if cacheDir != "" {
		if helpers.FilePathSeparator != cacheDir[len(cacheDir)-1:] {
			cacheDir = cacheDir + helpers.FilePathSeparator
		}
		isDir, err := helpers.DirExists(cacheDir, hugofs.Source())
		utils.CheckErr(err)
		if isDir == false {
			mkdir(cacheDir)
		}
		viper.Set("CacheDir", cacheDir)
	} else {
		viper.Set("CacheDir", helpers.GetTempDir("hugo_cache", hugofs.Source()))
	}

	if verboseLog || logging || (viper.IsSet("LogFile") && viper.GetString("LogFile") != "") {
		if viper.IsSet("LogFile") && viper.GetString("LogFile") != "" {
			jww.SetLogFile(viper.GetString("LogFile"))
		} else {
			jww.UseTempLogFile("hugo")
		}
	} else {
		jww.DiscardLogging()
	}

	if viper.GetBool("verbose") {
		jww.SetStdoutThreshold(jww.LevelInfo)
	}

	if verboseLog {
		jww.SetLogThreshold(jww.LevelInfo)
	}

	jww.INFO.Println("Using config file:", viper.ConfigFileUsed())

	// Init file systems. This may be changed at a later point.
	hugofs.InitDefaultFs()

	themeDir := helpers.GetThemeDir()
	if themeDir != "" {
		if _, err := hugofs.Source().Stat(themeDir); os.IsNotExist(err) {
			return newSystemError("Unable to find theme Directory:", themeDir)
		}
	}

	themeVersionMismatch, minVersion := isThemeVsHugoVersionMismatch()

	if themeVersionMismatch {
		jww.ERROR.Printf("Current theme does not support Hugo version %s. Minimum version required is %s\n",
			helpers.HugoReleaseVersion(), minVersion)
	}

	return nil

}
