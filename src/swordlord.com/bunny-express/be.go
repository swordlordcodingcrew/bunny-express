package main

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
	"fmt"
	"os"
	"swordlord.com/bunny-express/cmd"
	"swordlord.com/bunny-express/common"
	"swordlord.com/bunny-express/db"
)

func main() {

	common.InitConfig()
	common.InitLog()

	db.CheckDatabase()

	fmt.Println(` ______                           _______                                        `)
	fmt.Println(`|   __ \.--.--.-----.-----.--.--.|    ___|.--.--.-----.----.-----.-----.-----.    (\(\`)
	fmt.Println(`|   __ <|  |  |     |     |  |  ||    ___||_   _|  _  |   _|  -__|__ --|__ --|   ( =':')`)
	fmt.Println(`|______/|_____|__|__|__|__|___  ||_______||__.__|   __|__| |_____|_____|_____|   (..(")(")`)
	fmt.Println(`                          |_____|               |__|                             `)

	fmt.Println("")
	fmt.Println("CLI based mailbox configuration for Postfix and Dovecot")
	fmt.Println("(c) 2018-19 by SwordLord - the coding crew")
	fmt.Println("")

	// initialise the command structure
	if err := cmd.RootCmd.Execute(); err != nil {

		fmt.Println("Your command returned an error. You might want to run with --help.")
		fmt.Println(err)

		os.Exit(1)
	}
}

// TODO: alias -> forward -> multi line, one line per address
// TODO: @ in alias == catchall? -> fix in sql syntax
// TODO: remove logrus, doesnt make sense in a cli tool...
