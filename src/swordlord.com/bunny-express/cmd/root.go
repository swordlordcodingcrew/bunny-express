package cmd

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
	"github.com/spf13/cobra"
	"swordlord.com/bunny-express/common"
)

// var cfgFile string // see init() for details

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:           "be",
	SilenceUsage:  true, // only show help if explicitly requested
	SilenceErrors: true, // we show errors on our own...
	Version:       common.GetVersion(),
	Short:         "bunnyexpress, the CLI-based tool for postfix & dovecot mailbox administration.",
	Long: `bunnyexpress, the CLI-based tool for postfix & dovecot mailbox administration.

With bunnyexpress you can manage domains, mailboxes and aliases for your own mail domains and infrastructure. 
Everything is stored within an SQLite3 database. See accompanied ReadMe and help for more details.`,
}

func init() {

	// following lines just for reference.

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gohjasmincli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
