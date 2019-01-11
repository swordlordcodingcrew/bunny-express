package domain

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

type Domain struct {
	Domain      string         `db:"domain"`
	Description sql.NullString `db:"desc"`
	Mailbox     int            `db:"mailbox"`
	Alias       int            `db:"alias"`
	IsActive    bool           `db:"active"`
	CrtDat      time.Time      `db:"crt_dat"`
	UpdDat      time.Time      `db:"upd_dat"`
}

func GetFieldCaptions() []string {

	captions := []string{"Domain", "Description", "Mailbox", "Alias", "Active", "Created", "Updated"}

	return captions
}

func GetAllDomains() ([]Domain, error) {

	db, err := db.OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	q := `SELECT 
			  domain, 
			  desc,  
			  (SELECT count(mail) FROM mailbox WHERE mailbox.domain = domain.domain) as mailbox,
			  (SELECT count(alias) FROM alias WHERE alias.domain = domain.domain) as alias,
			  active,
			  crt_dat,
			  upd_dat
		  FROM 
		  	  domain 
		  ORDER BY domain ASC`

	stmt, err := db.Preparex(q)
	if err != nil {
		return nil, err
	}

	var d []Domain
	err = stmt.Select(&d)

	return d, err
}

func GetDomain(name string) (Domain, error) {

	db, err := db.OpenDB()
	if err != nil {
		return Domain{}, err
	}
	defer db.Close()

	stmt, err := db.Preparex(`SELECT * FROM domain WHERE domain=?`)
	if err != nil {
		return Domain{}, err
	}

	var d Domain
	err = stmt.Get(&d, name)
	if err != nil {
		return Domain{}, err
	} else {
		return d, nil
	}
}

func AddDomain(d Domain) error {

	db, err := db.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Preparex("INSERT INTO domain (domain, desc, active) VALUES (?,?,?)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(d.Domain, d.Description, d.IsActive)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fields := logrus.Fields{"domain": d.Domain, "description": d.Description, "active": d.IsActive}

	if count == 0 {
		common.LogInfo("Nothing added.", fields)
	} else {
		common.LogInfo("Domain added.", fields)
	}

	return nil
}

func EditDomain(d Domain) error {

	db, err := db.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Preparex("UPDATE domain SET desc = ?, active = ?, upd_dat = ? WHERE domain = ? AND upd_dat <= ?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(d.Description, d.IsActive, time.Now(), d.Domain, d.UpdDat)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fields := logrus.Fields{"domain": d.Domain, "description": d.Description, "active": d.IsActive}

	if count == 0 {
		common.LogInfo("Nothing changed.", fields)
	} else {
		common.LogInfo("Domain updated.", fields)
	}

	return nil
}

func DeleteDomain(name string) error {

	db, err := db.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Preparex(`DELETE FROM domain WHERE domain=?`)
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
		common.LogInfo("Nothing deleted. Wrong domain used?", logrus.Fields{"domain": name})
	} else {
		common.LogInfo("Domain deleted.", logrus.Fields{"domain": name, "count": count})
	}

	return nil
}
