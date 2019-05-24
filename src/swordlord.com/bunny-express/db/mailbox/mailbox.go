package mailbox

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
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"swordlord.com/bunny-express/common"
	"swordlord.com/bunny-express/db"
	"time"
)

type Mailbox struct {
	Mail               string         `db:"mail"`
	Description        sql.NullString `db:"desc"`
	isDescDirty        bool
	Domain             string `db:"domain"`
	isDomainDirty      bool
	Password           string `db:"pwd"`
	isPasswordDirty    bool
	MailDir            string `db:"mail_dir"`
	isMailDirDirty     bool
	LocalPart          string `db:"local_part"`
	isLocalPartDirty   bool
	RelayDomain        sql.NullString `db:"relay_domain"`
	isRelayDomainDirty bool
	Quota              sql.NullString `db:"quota"`
	isQuotaDirty       bool
	IsActive           bool `db:"active"`
	isIsActiveDirty    bool
	// tells us if object is from db or not
	isNew  bool
	CrtDat time.Time `db:"crt_dat"`
	UpdDat time.Time `db:"upd_dat"`
}

func NewMailbox() *Mailbox {

	m := &Mailbox{}
	m.clearDirtyFlags()
	m.isNew = true

	return m
}

func (m *Mailbox) clearDirtyFlags() {
	m.isDescDirty = false
	m.isDomainDirty = false
	m.isPasswordDirty = false
	m.isMailDirDirty = false
	m.isLocalPartDirty = false
	m.isRelayDomainDirty = false
	m.isQuotaDirty = false
	m.isIsActiveDirty = false
}

func (m *Mailbox) GetMail() string                { return m.Mail }
func (m *Mailbox) GetDescription() sql.NullString { return m.Description }
func (m *Mailbox) GetDomain() string              { return m.Domain }
func (m *Mailbox) GetPasssword() string           { return m.Password }
func (m *Mailbox) GetMailDir() string             { return m.MailDir }
func (m *Mailbox) GetLocalPart() string           { return m.LocalPart }
func (m *Mailbox) GetRelayDomain() sql.NullString { return m.RelayDomain }
func (m *Mailbox) GetQuota() sql.NullString       { return m.Quota }
func (m *Mailbox) GetIsActive() bool              { return m.IsActive }

func (m *Mailbox) SetMail(mail string) {
	m.Mail = mail
}

func (m *Mailbox) SetDescription(description sql.NullString) {

	// TODO: add sanity check to all SetXY functions
	if m.Description.String == description.String {
		return
	}

	m.Description = description
	m.isDescDirty = true
}

func (m *Mailbox) SetDomain(domain string) {
	m.Domain = domain
	m.isDomainDirty = true
}

func (m *Mailbox) SetPassword(password string) error {

	// TODO; add check
	pwd, err := generatePassword(password)
	if err != nil {
		return err
	}

	m.Password = pwd
	m.isPasswordDirty = true

	return nil
}

func (m *Mailbox) SetMailDir(mailDir string) {
	m.MailDir = mailDir
	m.isMailDirDirty = true
}

func (m *Mailbox) SetLocalPart(localPart string) {
	m.LocalPart = localPart
	m.isLocalPartDirty = true
}

func (m *Mailbox) SetRelayDomain(relayDomain sql.NullString) {
	m.RelayDomain = relayDomain
	m.isRelayDomainDirty = true
}

func (m *Mailbox) SetQuota(quota sql.NullString) {
	m.Quota = quota
	m.isQuotaDirty = true
}

func (m *Mailbox) SetIsActive(ia bool) {

	if m.IsActive == ia {
		return
	}

	m.IsActive = ia
	m.isIsActiveDirty = true
}

func (m *Mailbox) IsDirty() bool {
	if m.isDescDirty ||
		m.isDomainDirty ||
		m.isPasswordDirty ||
		m.isMailDirDirty ||
		m.isLocalPartDirty ||
		m.isRelayDomainDirty ||
		m.isQuotaDirty ||
		m.isIsActiveDirty {
		return true
	} else {
		return false
	}
}

type MailboxFilter struct {
	Mail        string
	Description string
	Domain      string
	MailDir     string
	LocalPart   string
	RelayDomain string
	Quota       string
	IsActive    sql.NullBool
}

func GetFieldCaptions() []string {

	captions := []string{"Mail", "Description", "Domain", "Password", "MailDir", "LocalPart",
		"RelayDomain", "Quota", "Active", "Created", "Updated"}

	return captions
}

func GetAllMailboxen() ([]Mailbox, error) {

	return GetFilteredMailbox(&MailboxFilter{})
}

func GetFilteredMailbox(mf *MailboxFilter) ([]Mailbox, error) {

	db, err := db.OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sFilter := ""
	params := []string{}

	if len(mf.Domain) > 0 {
		sFilter += "domain LIKE ?"
		params = append(params, mf.Domain)
	}

	if len(mf.MailDir) > 0 {
		if len(sFilter) > 0 {
			sFilter += " AND "
		}
		sFilter += "mail_dir LIKE ?"
		params = append(params, mf.MailDir)
	}

	if len(mf.LocalPart) > 0 {
		if len(sFilter) > 0 {
			sFilter += " AND "
		}
		sFilter += "local_part LIKE ?"
		params = append(params, mf.LocalPart)
	}

	if len(mf.RelayDomain) > 0 {
		if len(sFilter) > 0 {
			sFilter += " AND "
		}
		sFilter += "relay_domain LIKE ?"
		params = append(params, mf.RelayDomain)
	}

	if len(mf.Quota) > 0 {
		if len(sFilter) > 0 {
			sFilter += " AND "
		}
		sFilter += "quota LIKE ?"
		params = append(params, mf.Quota)
	}

	if mf.IsActive.Valid {
		if len(sFilter) > 0 {
			sFilter += " AND "
		}
		sFilter += "active = ?"
		params = append(params, strconv.FormatBool(mf.IsActive.Bool))
	}

	if len(sFilter) > 0 {
		sFilter = "WHERE " + sFilter
	}

	sql := "SELECT * FROM mailbox " + sFilter + " ORDER BY domain, mail ASC"

	stmt, err := db.Preparex(sql)
	if err != nil {
		return nil, err
	}

	// slice of interface != slice of string, which is why we copy the values to a slice of interfaces
	args := make([]interface{}, len(params))
	for i, s := range params {
		args[i] = s
	}

	// select function accepts a slice of interface as variadic, neat
	var m []Mailbox
	err = stmt.Select(&m, args...)

	if err == nil {
		for i := range m {
			m[i].isNew = false
		}
	}

	return m, err
}

func GetMailbox(name string) (*Mailbox, error) {

	db, err := db.OpenDB()
	if err != nil {
		return NewMailbox(), err
	}
	defer db.Close()

	stmt, err := db.Preparex(db.Rebind("SELECT * FROM mailbox WHERE mail=?"))
	if err != nil {
		return NewMailbox(), err
	}

	m := NewMailbox()
	err = stmt.Get(m, name)
	if err != nil {
		return NewMailbox(), err
	} else {
		m.isNew = false
		return m, nil
	}
}

func (m *Mailbox) Persist() error {

	db, err := db.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if !m.IsDirty() {
		common.LogInfo("Mailbox did not change, not persisted.", nil)
		return nil
	}

	if m.isNew {
		err = m.add(db)
	} else {
		err = m.update(db)
	}

	return err
}

func (m *Mailbox) add(db *sqlx.DB) error {

	sFields := ""
	var params []interface{}

	if m.isNew {
		sFields += "mail"
		params = append(params, m.Mail)
	}

	if m.isDomainDirty {
		if len(sFields) > 0 {
			sFields += ", "
		}
		sFields += "domain"
		params = append(params, m.Domain)
	}

	if m.isDescDirty {
		if len(sFields) > 0 {
			sFields += ", "
		}
		sFields += "desc"
		params = append(params, m.Description.String)
	}

	if m.isPasswordDirty {
		if len(sFields) > 0 {
			sFields += ", "
		}
		sFields += "pwd"
		params = append(params, m.Password)
	}

	if m.isMailDirDirty {
		if len(sFields) > 0 {
			sFields += ", "
		}
		sFields += "mail_dir"
		params = append(params, m.MailDir)
	}

	if m.isLocalPartDirty {
		if len(sFields) > 0 {
			sFields += ", "
		}
		sFields += "local_part"
		params = append(params, m.LocalPart)
	}

	if m.isRelayDomainDirty {
		if len(sFields) > 0 {
			sFields += ", "
		}
		sFields += "relay_domain"
		params = append(params, m.RelayDomain.String)
	}

	if m.isQuotaDirty {
		if len(sFields) > 0 {
			sFields += ", "
		}
		sFields += "quota"
		params = append(params, m.Quota.String)
	}

	// generate param string, remove last , from repeater
	sQM := strings.Repeat("?,", len(params))
	sQM = sQM[:len(sQM)-1]

	sQuery := "INSERT INTO mailbox (" + sFields + ") VALUES (" + sQM + ")"

	stmt, err := db.Preparex(sQuery)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(params...)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fields := logrus.Fields{"mail": m.Mail, "domain": m.Domain, "description": m.Description, "active": m.IsActive}

	m.clearDirtyFlags()

	if count == 0 {
		common.LogInfo("Nothing done.", fields)
	} else {
		common.LogInfo("Mailbox added.", fields)
	}

	return nil
}

func (m *Mailbox) update(db *sqlx.DB) error {

	if !m.IsDirty() {
		return errors.New("trying to update unchanged object")
	}

	sStatement := ""
	var params []interface{}

	if m.isDomainDirty {
		sStatement += "domain = ?"
		params = append(params, m.Domain)
	}

	if m.isDescDirty {
		if len(sStatement) > 0 {
			sStatement += ", "
		}
		sStatement += "desc = ?"
		params = append(params, m.Description.String)
	}

	if m.isPasswordDirty {
		if len(sStatement) > 0 {
			sStatement += ", "
		}
		sStatement += "pwd = ?"
		params = append(params, m.Password)
	}

	if m.isMailDirDirty {
		if len(sStatement) > 0 {
			sStatement += ", "
		}
		sStatement += "mail_dir = ?"
		params = append(params, m.MailDir)
	}

	if m.isLocalPartDirty {
		if len(sStatement) > 0 {
			sStatement += ", "
		}
		sStatement += "local_part = ?"
		params = append(params, m.LocalPart)
	}

	if m.isRelayDomainDirty {
		if len(sStatement) > 0 {
			sStatement += ", "
		}
		sStatement += "relay_domain = ?"
		params = append(params, m.RelayDomain.String)
	}

	if m.isQuotaDirty {
		if len(sStatement) > 0 {
			sStatement += ", "
		}
		sStatement += "quota = ?"
		params = append(params, m.Quota.String)
	}

	// update upddat field
	if len(sStatement) > 0 {
		sStatement += ", "
	}
	sStatement += "upd_dat = ?"
	params = append(params, time.Now())

	// append params for where
	params = append(params, m.Mail)   // pkey
	params = append(params, m.UpdDat) // optimistic locking

	stmt, err := db.Preparex("UPDATE mailbox SET " + sStatement + " WHERE mail = ? AND upd_dat <= ?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(params...)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fields := logrus.Fields{"mail": m.Mail, "domain": m.Domain, "description": m.Description, "active": m.IsActive}

	m.clearDirtyFlags()

	if count == 0 {
		common.LogInfo("Nothing done.", fields)
	} else {
		common.LogInfo("Mailbox updated.", fields)
	}

	return nil

}

func DeleteMailbox(name string) error {

	db, err := db.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Preparex(`DELETE FROM mailbox WHERE mail=?`)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(name)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		common.LogInfo("Nothing deleted. Wrong mailbox used?", logrus.Fields{"mailbox": name})
	} else {
		common.LogInfo("Mailbox deleted.", logrus.Fields{"mailbox": name, "count": count})
	}

	return nil
}

func FillDefaultMailboxOnDomain(domain string) error {

	mailboxen := common.GetStringSliceFromConfig("default.mailbox")

	for _, mn := range mailboxen {

		m := NewMailbox()
		m.SetDomain(domain)
		m.SetMail(mn + "@" + domain)
		m.SetPassword(domain)
		m.IsActive = true
		m.Description.String = "filled automatically with default mailbox from config"
		m.Description.Valid = true

		err := m.Persist()
		if err != nil {
			common.LogInfo("AddMailbox returned an error.", logrus.Fields{"mailbox": mn, "error": err})
			return err
		}
	}

	return nil
}

func generatePassword(password string) (string, error) {

	pwd := []byte(password)

	// Hashing the pwd with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil

	// Comparing the pwd with the hash
	// err = bcrypt.CompareHashAndPassword(hashedPassword, pwd)
	// fmt.Println(err) // nil means it is a match
}
