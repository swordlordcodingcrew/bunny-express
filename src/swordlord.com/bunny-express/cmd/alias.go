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
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"swordlord.com/bunny-express/db/alias"
	"swordlord.com/bunny-express/util"
)

func ListAliases(cmd *cobra.Command, args []string) error {

	aa, err := alias.GetAllAliases()
	if err != nil {
		return fmt.Errorf("command 'list' returns an error %s", err)
	}

	var aliases [][]string

	for _, a := range aa {

		aliases = append(aliases, []string{a.Alias, a.Description.String, a.Domain, a.ForwardAddress, strconv.FormatBool(a.IsActive), a.CrtDat.Format("2006-01-02 15:04:05"), a.UpdDat.Format("2006-01-02 15:04:05")})
	}

	util.WriteTable(alias.GetFieldCaptions(), aliases)

	return nil
}

func AddAlias(cmd *cobra.Command, args []string) error {

	a := alias.Alias{}

	a.Alias = args[0]

	fActive := cmd.Flag("active")
	b, err := strconv.ParseBool(fActive.Value.String())
	if err == nil {
		a.IsActive = b
	}

	fDesc := cmd.Flag("description")
	if fDesc.Changed {

		a.Description.String = fDesc.Value.String()
	}

	return alias.AddAlias(a)
}

func EditAlias(cmd *cobra.Command, args []string) error {

	m, err := alias.GetAlias(args[0])
	if err != nil {
		return fmt.Errorf("command 'edit' returns an error %s", err)
	}

	// todo: dont do it like that, keep dirty flags per field on the struct
	isDirty := false

	fActive := cmd.Flag("active")
	if fActive.Changed {

		b, err := strconv.ParseBool(fActive.Value.String())
		if err == nil {
			m.IsActive = b
			isDirty = true
		}
	}

	fDesc := cmd.Flag("description")
	if fDesc.Changed {

		m.Description.String = fDesc.Value.String()
		isDirty = true
	}

	if isDirty {
		err = alias.EditAlias(m)
	} else {
		err = fmt.Errorf("nothing to change")
	}

	return err

}

func DeleteAlias(cmd *cobra.Command, args []string) error {

	return alias.DeleteAlias(args[0])
}

func init() {

	// calCmd represents the domain command
	var aliasCmd = &cobra.Command{
		Use:   "alias",
		Short: "Add, change and manage aliases.",
		Long:  `Add, change and manage aliases. Requires a subcommand.`,
		RunE:  nil,
	}

	var aliasListCmd = &cobra.Command{
		Use:   "list",
		Short: "List aliases.",
		Long:  `List aliases based on filter given. Name of alias can have wildcard.`,
		RunE:  ListAliases,
	}
	aliasListCmd.Flags().BoolP("active", "a", true, "is alias active")

	var aliasAddCmd = &cobra.Command{
		Use:   "add [alias] [domain] [forward_address]",
		Short: "Add new alias to given domain",
		Long: `Add new alias with parameters given and add it to the given domain.

The forward_address can contain multiple addresses. Please make sure to add an empty character between the addresses`,
		Args: cobra.ExactArgs(3),
		RunE: AddAlias,
	}
	aliasAddCmd.Flags().BoolP("active", "a", true, "is alias active")
	aliasAddCmd.Flags().StringP("description", "d", "", "description for this alias")
	aliasAddCmd.Flags().StringP("maildir", "m", "", "maildir to be used")
	aliasAddCmd.Flags().StringP("localpart", "l", "", "local part, better not change this")
	aliasAddCmd.Flags().StringP("relaydomain", "r", "", "relay domain")
	aliasAddCmd.Flags().StringP("quota", "q", "", "quota for this user")

	var aliasEditCmd = &cobra.Command{
		Use:   "edit [alias]",
		Short: "Edit an existing alias",
		Long:  `Edit.`,
		Args:  cobra.ExactArgs(1),
		RunE:  EditAlias,
	}
	aliasEditCmd.Flags().BoolP("active", "a", true, "is alias active")
	aliasEditCmd.Flags().StringP("description", "d", "", "description for this alias")
	aliasEditCmd.Flags().StringP("password", "p", "", "password in clear")
	aliasEditCmd.Flags().StringP("maildir", "m", "", "maildir to be used")
	aliasEditCmd.Flags().StringP("localpart", "l", "", "local part, better not change this")
	aliasEditCmd.Flags().StringP("relaydomain", "r", "", "relay domain")
	aliasEditCmd.Flags().StringP("quota", "q", "", "quota for this user")

	var aliasDeleteCmd = &cobra.Command{
		Use:   "delete [alias]",
		Short: "Deletes a alias.",
		Long:  `Deletes a alias.`,
		Args:  cobra.ExactArgs(1),
		RunE:  DeleteAlias,
	}

	RootCmd.AddCommand(aliasCmd)

	aliasCmd.AddCommand(aliasListCmd)
	aliasCmd.AddCommand(aliasAddCmd)
	aliasCmd.AddCommand(aliasEditCmd)
	aliasCmd.AddCommand(aliasDeleteCmd)
}