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
	"github.com/sirupsen/logrus"
	"swordlord.com/bunny-express/common"
	"swordlord.com/bunny-express/db"
	"time"
)

type Mailbox struct {
	Mail        string         `db:"mail"`
	Description sql.NullString `db:"desc"`
	Domain      string         `db:"domain"`
	Password    string         `db:"pwd"`
	MailDir     string         `db:"mail_dir"`
	LocalPart   string         `db:"local_part"`
	RelayDomain sql.NullString `db:"relay_domain"`
	Quota       sql.NullString `db:"quota"`
	IsActive    bool           `db:"active"`
	CrtDat      time.Time      `db:"crt_dat"`
	UpdDat      time.Time      `db:"upd_dat"`
}

func GetFieldCaptions() []string {

	captions := []string{"Mail", "Description", "Domain", "Password", "MailDir", "LocalPart",
		"RelayDomain", "Quota", "Active", "Created", "Updated"}

	return captions
}

func GetAllMailboxen() ([]Mailbox, error) {

	db, err := db.OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Preparex(`SELECT * FROM mailbox ORDER BY domain, mail ASC`)
	if err != nil {
		return nil, err
	}

	var m []Mailbox
	err = stmt.Select(&m)

	return m, err
}

func GetMailbox(name string) (Mailbox, error) {

	db, err := db.OpenDB()
	if err != nil {
		return Mailbox{}, err
	}
	defer db.Close()

	stmt, err := db.Preparex(`SELECT * FROM mailbox WHERE mail=?`)
	if err != nil {
		return Mailbox{}, err
	}

	var m Mailbox
	err = stmt.Get(&m, name)
	if err != nil {
		return Mailbox{}, err
	} else {
		return m, nil
	}
}

func AddMailbox(m Mailbox) error {

	db, err := db.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// TODO;
	// Password string `db:"pwd"`
	//	MailDir string `db:"mail_dir"`
	//	LocalPart string `db:"local_part"`
	//	RelayDomain string `db:"relay_domain"`
	//	Quota

	stmt, err := db.Preparex("INSERT INTO mailbox (mail, domain, desc, active) VALUES (?,?,?)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(m.Domain, m.Description, m.IsActive)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fields := logrus.Fields{"mail": m.Mail, "domain": m.Domain, "description": m.Description, "active": m.IsActive}

	if count == 0 {
		common.LogInfo("Nothing added.", fields)
	} else {
		common.LogInfo("Mailbox added.", fields)
	}

	return nil
}

func EditMailbox(m Mailbox) error {

	db, err := db.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Preparex("UPDATE mailbox SET desc = ?, active = ?, upd_dat = ? WHERE mail = ? AND upd_dat <= ?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(m.Description, m.IsActive, time.Now(), m.Mail, m.UpdDat)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fields := logrus.Fields{"mail": m.Mail, "domain": m.Domain, "description": m.Description, "active": m.IsActive}

	if count == 0 {
		common.LogInfo("Nothing changed.", fields)
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
