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
	"database/sql"
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

	df := domain.DomainFilter{}

	fIsActive := cmd.Flag("active")
	if fIsActive.Changed {
		bIsActive, err := strconv.ParseBool(fIsActive.Value.String())
		if err == nil {
			df.IsActive.Scan(bIsActive)
		}
	}

	fDomain := cmd.Flag("domain")
	if fDomain.Changed {
		df.Domain = fDomain.Value.String()
	}

	d, err := domain.GetFilteredDomains(&df)
	if err != nil {
		return fmt.Errorf("command 'list' returns an error %s", err)
	}

	var domains [][]string

	for _, domain := range d {

		domains = append(domains, []string{domain.GetDomain(), domain.GetDescription().String,
			strconv.Itoa(domain.GetMailboxCount()),
			strconv.Itoa(domain.GetAliasCount()),
			strconv.FormatBool(domain.GetIsActive()),
			domain.CrtDat.Format("2006-01-02 15:04:05"),
			domain.UpdDat.Format("2006-01-02 15:04:05")})
	}

	util.WriteTable(domain.GetFieldCaptions(), domains)

	return nil
}

func AddDomain(cmd *cobra.Command, args []string) error {

	d := domain.NewDomain()

	d.SetDomain(args[0])

	scanDomainFlagsToObject(cmd, d)

	err := d.Persist()
	if err != nil {
		return err
	}

	fFillDefaultAlias := cmd.Flag("fill")
	bFill, err := strconv.ParseBool(fFillDefaultAlias.Value.String())
	if bFill {

		err = alias.FillDefaultAliasOnDomain(d.Domain)
		if err != nil {
			common.LogInfo("Could not automatically create alias.", logrus.Fields{"domain": d.Domain, "error": err})
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

	scanDomainFlagsToObject(cmd, d)

	return d.Persist()
}

func scanDomainFlagsToObject(cmd *cobra.Command, d *domain.Domain) {

	fActive := cmd.Flag("active")
	if fActive.Changed {

		b, err := strconv.ParseBool(fActive.Value.String())
		if err == nil {
			d.SetIsActive(b)
		}
	}

	fDesc := cmd.Flag("description")
	if fDesc.Changed {

		var s = sql.NullString{}
		err := s.Scan(fDesc.Value.String())
		if err == nil {
			d.SetDescription(s)
		}
	}
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
	domainListCmd.Flags().StringP("domain", "d", "", "mailbox for which domain")

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
