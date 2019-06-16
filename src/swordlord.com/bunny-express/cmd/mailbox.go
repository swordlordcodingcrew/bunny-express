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
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"swordlord.com/bunny-express/common"
	"swordlord.com/bunny-express/db/mailbox"
	"swordlord.com/bunny-express/util"
)

func ListMailbox(cmd *cobra.Command, args []string) error {

	mbf := mailbox.MailboxFilter{}

	fIsActive := cmd.Flag("active")
	if fIsActive.Changed {
		bIsActive, err := strconv.ParseBool(fIsActive.Value.String())
		if err == nil {
			mbf.IsActive.Scan(bIsActive)
		}
	}

	fDomain := cmd.Flag("domain")
	if fDomain.Changed {
		mbf.Domain = fDomain.Value.String()
	}

	ms, err := mailbox.GetFilteredMailbox(&mbf)
	if err != nil {
		return fmt.Errorf("command 'list' returns an error %s", err)
	}

	var mailboxen [][]string

	for _, mb := range ms {

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

	pwdScheme := checkSchemeFlag(cmd)

	m := mailbox.NewMailbox()

	m.SetMail(args[0])
	m.SetPassword(args[1], pwdScheme)
	m.SetDomain(args[2])

	m.SetMailDir("")
	m.SetLocalPart("")

	m.SetQuota(0)

	scanMailboxFlagsToObject(cmd, m)

	return m.Persist()
}

func EditMailbox(cmd *cobra.Command, args []string) error {

	m, err := mailbox.GetMailbox(args[0])
	if err != nil {
		return fmt.Errorf("command 'edit' returns an error %s", err)
	}

	scanMailboxFlagsToObject(cmd, m)

	return m.Persist()
}

func scanMailboxFlagsToObject(cmd *cobra.Command, m *mailbox.Mailbox) {

	fActive := cmd.Flag("active")
	if fActive.Changed {

		b, err := strconv.ParseBool(fActive.Value.String())
		if err == nil {
			m.SetIsActive(b)
		}
	}

	fDesc := cmd.Flag("description")
	if fDesc.Changed {

		var s = sql.NullString{}
		err := s.Scan(fDesc.Value.String())
		if err == nil {
			m.SetDescription(s)
		}
	}

	// check for nil since this flag is not used in all commands
	fPassword := cmd.Flag("password")
	if fPassword != nil && fPassword.Changed {

		pwdScheme := checkSchemeFlag(cmd)

		m.SetPassword(fPassword.Value.String(), pwdScheme)
	}

	fMaildir := cmd.Flag("maildir")
	if fMaildir.Changed {

		m.SetMailDir(fMaildir.Value.String())
	}

	fLocalPart := cmd.Flag("localpart")
	if fLocalPart.Changed {

		m.SetLocalPart(fLocalPart.Value.String())
	}

	fRelayDomain := cmd.Flag("relaydomain")
	if fRelayDomain.Changed {

		var s = sql.NullString{}
		err := s.Scan(fRelayDomain.Value.String())
		if err == nil {
			m.SetRelayDomain(s)
		}
	}

	fQuota := cmd.Flag("quota")
	if fQuota.Changed {

		var s = sql.NullString{}
		err := s.Scan(fQuota.Value.String())
		if err == nil {
			m.SetQuotaAsNullString(s)
		}
	}
}

func checkSchemeFlag(cmd *cobra.Command) string {

	pwdScheme := common.GetStringFromConfig("default.scheme")
	if pwdScheme == "" {

		pwdScheme = "MD5-CRYPT"
	}

	fPwdScheme := cmd.Flag("pwdscheme")
	if fPwdScheme != nil && fPwdScheme.Changed {

		scheme := fPwdScheme.Value.String()
		switch scheme {
		case "md5crypt":
			pwdScheme = "MD5-CRYPT"
		case "bcrypt":
			pwdScheme = "BLF-CRYPT"
		}
	}

	return pwdScheme
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
	mailboxListCmd.Flags().StringP("domain", "d", "", "mailbox for which domain")

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
	mailboxAddCmd.Flags().StringP("pwdscheme", "s", "", "password hashing scheme to be used")

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
	mailboxEditCmd.Flags().StringP("pwdscheme", "s", "", "password hashing scheme to be used")

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
