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
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strconv"
	"swordlord.com/bunny-express/common"
	"swordlord.com/bunny-express/db/alias"
	"swordlord.com/bunny-express/db/domain"
	"swordlord.com/bunny-express/util"
)

func ListDomain(cmd *cobra.Command, args []string) error {

	d, err := domain.GetAllDomains()
	if err != nil {
		return fmt.Errorf("command 'list' returns an error %s", err)
	}

	var domains [][]string

	for _, domain := range d {

		domains = append(domains, []string{domain.Domain, domain.Description.String,
			strconv.Itoa(domain.Mailbox),
			strconv.Itoa(domain.Alias),
			strconv.FormatBool(domain.IsActive),
			domain.CrtDat.Format("2006-01-02 15:04:05"),
			domain.UpdDat.Format("2006-01-02 15:04:05")})
	}

	util.WriteTable(domain.GetFieldCaptions(), domains)

	return nil
}

func AddDomain(cmd *cobra.Command, args []string) error {

	d := domain.Domain{}

	d.Domain = args[0]

	fActive := cmd.Flag("active")
	b, err := strconv.ParseBool(fActive.Value.String())
	if err == nil {
		d.IsActive = b
	}

	fDesc := cmd.Flag("description")
	if fDesc.Changed {

		d.Description.String = fDesc.Value.String()
	}

	err = domain.AddDomain(d)
	if err != nil {
		return err
	}

	fFillDefaultAlias := cmd.Flag("fill")
	b, err = strconv.ParseBool(fFillDefaultAlias.Value.String())
	if b {

		err = alias.FillDefaultAliasOnDomain(d.Domain)
		if err != nil {
			common.LogInfo("Could not automatically create Alias.", logrus.Fields{"domain": d.Domain, "error": err})
			return err
		}
	}

	return nil
}

func EditDomain(cmd *cobra.Command, args []string) error {

	d, err := domain.GetDomain(args[0])
	if err != nil {
		return fmt.Errorf("command 'edit' returns an error %s", err)
	}

	// todo: dont do it like that, keep dirty flags per field on the struct
	isDirty := false

	fActive := cmd.Flag("active")
	if fActive.Changed {

		b, err := strconv.ParseBool(fActive.Value.String())
		if err == nil {
			d.IsActive = b
			isDirty = true
		}
	}

	fDesc := cmd.Flag("description")
	if fDesc.Changed {

		d.Description.String = fDesc.Value.String()
		isDirty = true
	}

	if isDirty {
		err = domain.EditDomain(d)
	} else {
		err = fmt.Errorf("nothing to change")
	}

	return err
}

func DeleteDomain(cmd *cobra.Command, args []string) error {

	return domain.DeleteDomain(args[0])
}

func init() {

	// calCmd represents the domain command
	var domainCmd = &cobra.Command{
		Use:   "domain",
		Short: "Add, change and manage domains.",
		Long:  `Add, change and manage domains. Requires a subcommand.`,
		RunE:  nil,
	}

	var domainListCmd = &cobra.Command{
		Use:   "list",
		Short: "GetAllDomains domains.",
		Long:  `GetAllDomains domains based on filter given. Domain can have wildcard.`,
		RunE:  ListDomain,
	}
	domainListCmd.Flags().BoolP("active", "a", true, "is domain active")

	var domainAddCmd = &cobra.Command{
		Use:   "add [domain]",
		Short: "Add new domain",
		Long:  `Add new domain with parameters given.`,
		Args:  cobra.ExactArgs(1),
		RunE:  AddDomain,
	}
	domainAddCmd.Flags().BoolP("active", "a", true, "is domain active")
	domainAddCmd.Flags().StringP("description", "d", "", "description for this domain")
	domainAddCmd.Flags().BoolP("fill", "f", false, "add default aliases to the new domain")

	var domainEditCmd = &cobra.Command{
		Use:   "edit [domain]",
		Short: "Edit an existing domain",
		Long:  `Edit an existing domain. Only parameters given in flags are changed`,
		Args:  cobra.ExactArgs(1),
		RunE:  EditDomain,
	}
	domainEditCmd.Flags().BoolP("active", "a", true, "is domain active")
	domainEditCmd.Flags().StringP("description", "d", "", "description for this domain")

	var domainDeleteCmd = &cobra.Command{
		Use:   "delete [domain]",
		Short: "Deletes a domain.",
		Long:  `Deletes a domain.`,
		Args:  cobra.ExactArgs(1),
		RunE:  DeleteDomain,
	}

	flag.Parse()

	RootCmd.AddCommand(domainCmd)

	domainCmd.AddCommand(domainListCmd)
	domainCmd.AddCommand(domainAddCmd)
	domainCmd.AddCommand(domainEditCmd)
	domainCmd.AddCommand(domainDeleteCmd)
}
