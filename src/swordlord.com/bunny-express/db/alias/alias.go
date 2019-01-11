package alias

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

type Alias struct {
	Alias          string         `db:"alias"`
	Description    sql.NullString `db:"desc"`
	Domain         string         `db:"domain"`
	ForwardAddress string         `db:"forward_address"`
	IsActive       bool           `db:"active"`
	CrtDat         time.Time      `db:"crt_dat"`
	UpdDat         time.Time      `db:"upd_dat"`
}

func GetFieldCaptions() []string {

	captions := []string{"Alias", "Description", "Domain", "Forward", "Active", "Created", "Updated"}

	return captions
}

func GetAllAliases() ([]Alias, error) {

	db, err := db.OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Preparex(`SELECT * FROM alias ORDER BY domain, alias ASC`)
	if err != nil {
		return nil, err
	}

	var a []Alias
	err = stmt.Select(&a)

	return a, err
}

func GetAlias(name string) (Alias, error) {

	db, err := db.OpenDB()
	if err != nil {
		return Alias{}, err
	}
	defer db.Close()

	stmt, err := db.Preparex(`SELECT * FROM alias WHERE alias=?`)
	if err != nil {
		return Alias{}, err
	}

	var a Alias
	err = stmt.Get(&a, name)
	if err != nil {
		return Alias{}, err
	} else {
		return a, nil
	}
}

func AddAlias(a Alias) error {

	db, err := db.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	// TODO
	// Other fields...

	stmt, err := db.Preparex("INSERT INTO alias (alias, domain, desc, forward_address, active) VALUES (?, ?, ?,?,?)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(a.Alias, a.Domain, a.Description, a.ForwardAddress, a.IsActive)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fields := logrus.Fields{"alias": a.Alias, "domain": a.Domain, "forward": a.ForwardAddress, "description": a.Description, "active": a.IsActive}

	if count == 0 {
		common.LogInfo("Nothing done.", fields)
	} else {
		common.LogInfo("Alias added.", fields)
	}

	return nil
}

func EditAlias(a Alias) error {

	db, err := db.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Preparex("UPDATE alias SET forward_address = ? desc = ?, active = ?, upd_dat = ? WHERE alias = ? AND upd_dat <= ?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(a.ForwardAddress, a.Description, a.IsActive, time.Now(), a.Alias, a.UpdDat)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fields := logrus.Fields{"alias": a.Alias, "domain": a.Domain, "forward": a.ForwardAddress, "description": a.Description, "active": a.IsActive}

	if count == 0 {
		common.LogInfo("Nothing done.", fields)
	} else {
		common.LogInfo("Alias updated.", fields)
	}

	return nil
}

func DeleteAlias(name string) error {

	db, err := db.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Preparex(`DELETE FROM alias WHERE alias=?`)
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
		common.LogInfo("Nothing deleted. Wrong Alias used?", logrus.Fields{"alias": name})
	} else {
		common.LogInfo("Alias deleted.", logrus.Fields{"alias": name, "count": count})
	}

	return nil
}

func FillDefaultAliasOnDomain(domain string) error {

	aliases := common.GetStringSliceFromConfig("default.alias")

	for _, an := range aliases {

		alias := Alias{}
		alias.Domain = domain
		alias.Alias = an + "@" + domain
		alias.ForwardAddress = "root@" + domain
		alias.IsActive = true
		alias.Description.String = "filled automatically with default alias from config"
		alias.Description.Valid = true

		err := AddAlias(alias)
		if err != nil {
			common.LogInfo("AddAlias returned an error.", logrus.Fields{"alias": an, "error": err})
			return err
		}
	}

	return nil
}
