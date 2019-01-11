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
	"swordlord.com/bunny-express/db/mailbox"
	"swordlord.com/bunny-express/util"
)

func ListMailbox(cmd *cobra.Command, args []string) error {

	m, err := mailbox.GetAllMailboxen()
	if err != nil {
		return fmt.Errorf("command 'list' returns an error %s", err)
	}

	var mailboxen [][]string

	for _, mb := range m {

		mailboxen = append(mailboxen, []string{mb.Mail, mb.Description.String, mb.Domain, mb.Password, mb.MailDir,
			mb.LocalPart, mb.RelayDomain.String, mb.Quota.String,
			strconv.FormatBool(mb.IsActive),
			mb.CrtDat.Format("2006-01-02 15:04:05"),
			mb.UpdDat.Format("2006-01-02 15:04:05")})
	}

	util.WriteTable(mailbox.GetFieldCaptions(), mailboxen)

	return nil
}

func AddMailbox(cmd *cobra.Command, args []string) error {

	m := mailbox.Mailbox{}

	m.Domain = args[0]

	fActive := cmd.Flag("active")
	b, err := strconv.ParseBool(fActive.Value.String())
	if err == nil {
		m.IsActive = b
	}

	fDesc := cmd.Flag("description")
	if fDesc.Changed {

		m.Description.String = fDesc.Value.String()
	}

	return mailbox.AddMailbox(m)
}

func EditMailbox(cmd *cobra.Command, args []string) error {

	m, err := mailbox.GetMailbox(args[0])
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
		err = mailbox.EditMailbox(m)
	} else {
		err = fmt.Errorf("nothing to change")
	}

	return err

}

func DeleteMailbox(cmd *cobra.Command, args []string) error {

	return mailbox.DeleteMailbox(args[0])
}

func init() {

	// calCmd represents the domain command
	var mailboxCmd = &cobra.Command{
		Use:   "mailbox",
		Short: "Add, change and manage mailboxes.",
		Long:  `Add, change and manage mailboxes. Requires a subcommand.`,
		RunE:  nil,
	}

	var mailboxListCmd = &cobra.Command{
		Use:   "list",
		Short: "List mailboxes.",
		Long:  `List mailboxes based on filter given. Name of mailbox can have wildcard.`,
		RunE:  ListMailbox,
	}
	mailboxListCmd.Flags().BoolP("active", "a", true, "is mailbox active")

	var mailboxAddCmd = &cobra.Command{
		Use:   "add [mailbox] [password] [domain]",
		Short: "Add new mailbox to given domain",
		Long:  `Add new mailbox with parameters given and add it to the given domain.`,
		Args:  cobra.ExactArgs(3),
		RunE:  AddMailbox,
	}
	mailboxAddCmd.Flags().BoolP("active", "a", true, "is mailbox active")
	mailboxAddCmd.Flags().StringP("description", "d", "", "description for this mailbox")
	mailboxAddCmd.Flags().StringP("maildir", "m", "", "maildir to be used")
	mailboxAddCmd.Flags().StringP("localpart", "l", "", "local part, better not change this")
	mailboxAddCmd.Flags().StringP("relaydomain", "r", "", "relay domain")
	mailboxAddCmd.Flags().StringP("quota", "q", "", "quota for this user")

	var mailboxEditCmd = &cobra.Command{
		Use:   "edit [mailbox]",
		Short: "Edit an existing mailbox",
		Long:  `Edit.`,
		Args:  cobra.ExactArgs(1),
		RunE:  EditMailbox,
	}
	mailboxEditCmd.Flags().BoolP("active", "a", true, "is mailbox active")
	mailboxEditCmd.Flags().StringP("description", "d", "", "description for this mailbox")
	mailboxEditCmd.Flags().StringP("password", "p", "", "password in clear")
	mailboxEditCmd.Flags().StringP("maildir", "m", "", "maildir to be used")
	mailboxEditCmd.Flags().StringP("localpart", "l", "", "local part, better not change this")
	mailboxEditCmd.Flags().StringP("relaydomain", "r", "", "relay domain")
	mailboxEditCmd.Flags().StringP("quota", "q", "", "quota for this user")

	var mailboxDeleteCmd = &cobra.Command{
		Use:   "delete [mailbox]",
		Short: "Deletes a mailbox.",
		Long:  `Deletes a mailbox.`,
		Args:  cobra.ExactArgs(1),
		RunE:  DeleteMailbox,
	}

	RootCmd.AddCommand(mailboxCmd)

	mailboxCmd.AddCommand(mailboxListCmd)
	mailboxCmd.AddCommand(mailboxAddCmd)
	mailboxCmd.AddCommand(mailboxEditCmd)
	mailboxCmd.AddCommand(mailboxDeleteCmd)
}
