package common

/*-----------------------------------------------------------------------------
 ** ______                           _______
 **|   __ \.--.--.-----.-----.--.--.|    ___|.--.--.-----.----.-----.-----.-----.
 **|   __ <|  |  |     |     |  |  ||    ___||_   _|  _  |   _|  -__|__ --|__ --|
 **|______/|_____|__|__|__|__|___  ||_______||__.__|   __|__| |_____|_____|_____|
 **                          |_____|               |__|
 **
 ** CLI-based tool for postfix / dovecot user administration
 **
 ** Copyright 2018-19 by SwordLord - the coding crew - http://www.swordlord.com
 ** and contributing authors
 **
 ** This program is free software; you can redistribute it and/or modify it
 ** under the terms of the GNU Affero General Public License as published by the
 ** Free Software Foundation, either version 3 of the License, or (at your option)
 ** any later version.
 **
 ** This program is distributed in the hope that it will be useful, but WITHOUT
 ** ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 ** FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
 ** for more details.
 **
 ** You should have received a copy of the GNU Affero General Public License
 ** along with this program. If not, see <http://www.gnu.org/licenses/>.
 **
 **-----------------------------------------------------------------------------
 **
 ** Original Authors:
 ** LordEidi@swordlord.com
 **
-----------------------------------------------------------------------------*/

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
)

const toolName = "be"
const versionID = "0.1.0"

func InitConfig() {

	// we look in these dirs for the config file
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/." + toolName)
	viper.AddConfigPath("/etc/" + toolName)

	// the file is expected to be named hrafn.config.json
	viper.SetConfigName(toolName + ".config")
	viper.SetConfigType("json")

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {

		// if not found, write a standard config file and quit...
		writeStandardConfig()

		// quit execution
		LogFatal("Error reading config file. New file dumped.", logrus.Fields{"error": err, "filename": toolName + ".config.json"})
	}
}

// Helper function
func GetStringFromConfig(key string) string {

	return viper.GetString(key)
}

func GetVersion() string {

	return versionID
}

func GetIntFromConfig(key string) int {

	return viper.GetInt(key)
}

// values need to be separated by empty char ( )
func GetStringSliceFromConfig(key string) []string {
	return viper.GetStringSlice(key)
}

func GetBoolFromConfig(key string, def bool) bool {

	if !viper.IsSet(key) {
		return def
	} else {
		return viper.GetBool(key)
	}
}

func GetLogLevel() string {

	loglevel := viper.GetString("log.level")
	if loglevel == "" {

		return "warn"
	} else {

		return loglevel
	}
}

//
func writeStandardConfig() error {

	err := ioutil.WriteFile(toolName+".config.json", defaultConfig, 0700)

	return err
}

//
var defaultConfig = []byte(`
{
  "log": {
    "level": "debug"
  },
  "db": {
    "driver": "sqlite3",
    "file": "be.sqlite",
    "add_demo_data": "false"
  },
  "default.alias": "info abuse"
}
`)
