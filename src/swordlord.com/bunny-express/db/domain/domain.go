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
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"swordlord.com/bunny-express/common"
	"swordlord.com/bunny-express/db"
	"time"
)

type Domain struct {
	Domain          string         `db:"domain"`
	Description     sql.NullString `db:"desc"`
	isDescDirty     bool
	MailboxCount    int  `db:"mailbox_count"` // dynamically loaded, not stored
	AliasCount      int  `db:"alias_count"`   // dynamically loaded, not stored
	IsActive        bool `db:"active"`
	isIsActiveDirty bool
	// tells us if object is from db or not
	isNew  bool
	CrtDat time.Time `db:"crt_dat"`
	UpdDat time.Time `db:"upd_dat"`
}

func NewDomain() *Domain {

	d := &Domain{}
	d.clearDirtyFlags()
	d.isNew = true
	d.MailboxCount = 0
	d.AliasCount = 0

	return d
}

func (m *Domain) clearDirtyFlags() {
	m.isDescDirty = false
	m.isIsActiveDirty = false
}

func (d *Domain) GetDomain() string              { return d.Domain }
func (d *Domain) GetDescription() sql.NullString { return d.Description }
func (d *Domain) GetMailboxCount() int           { return d.MailboxCount }
func (d *Domain) GetAliasCount() int             { return d.AliasCount }
func (d *Domain) GetIsActive() bool              { return d.IsActive }

func (d *Domain) SetDomain(domain string) {
	d.Domain = domain
}

func (d *Domain) SetDescription(description sql.NullString) {

	// TODO: add sanity check to all SetXY functions
	if d.Description.String == description.String {
		return
	}

	d.Description = description
	d.isDescDirty = true
}

func (d *Domain) SetIsActive(ia bool) {

	if d.IsActive == ia {
		return
	}

	d.IsActive = ia
	d.isIsActiveDirty = true
}

func (d *Domain) IsDirty() bool {
	if d.isDescDirty ||
		d.isIsActiveDirty {
		return true
	} else {
		return false
	}
}

type DomainFilter struct {
	Domain      string
	Description string
	IsActive    sql.NullBool
}

func GetFieldCaptions() []string {

	captions := []string{"Domain", "Description", "MailboxCount", "AliasCount", "Active", "Created", "Updated"}

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
			  (SELECT count(mail) FROM mailbox WHERE mailbox.domain = domain.domain) as mailbox_count,
			  (SELECT count(alias) FROM alias WHERE alias.domain = domain.domain) as alias_count,
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

func GetFilteredDomains(df *DomainFilter) ([]Domain, error) {

	db, err := db.OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sFilter := ""
	params := []string{}

	if len(df.Domain) > 0 {
		sFilter += "domain LIKE ?"
		params = append(params, df.Domain)
	}

	if df.IsActive.Valid {
		if len(sFilter) > 0 {
			sFilter += " AND "
		}
		sFilter += "active = ?"
		params = append(params, strconv.FormatBool(df.IsActive.Bool))
	}

	if len(sFilter) > 0 {
		sFilter = "WHERE " + sFilter
	}

	sql := `SELECT 
			  domain, 
			  desc,  
			  (SELECT count(mail) FROM mailbox WHERE mailbox.domain = domain.domain) as mailbox_count,
			  (SELECT count(alias) FROM alias WHERE alias.domain = domain.domain) as alias_count,
			  active,
			  crt_dat,
			  upd_dat
		  FROM 
		  	  domain ` + sFilter + ` 
		  ORDER BY domain ASC`

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
	var d []Domain
	err = stmt.Select(&d, args...)

	if err == nil {
		for i := range d {
			d[i].isNew = false
		}
	}

	return d, err
}

func GetDomain(domain string) (*Domain, error) {

	db, err := db.OpenDB()
	if err != nil {
		return NewDomain(), err
	}
	defer db.Close()

	stmt, err := db.Preparex(db.Rebind("SELECT * FROM domain WHERE domain=?"))
	if err != nil {
		return NewDomain(), err
	}

	d := NewDomain()
	err = stmt.Get(d, domain)
	if err != nil {
		return NewDomain(), err
	} else {
		d.isNew = false
		return d, nil
	}
}

func (d *Domain) Persist() error {

	db, err := db.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if !d.IsDirty() {
		common.LogInfo("Domain did not change, not persisted.", nil)
		return nil
	}

	if d.isNew {
		err = d.add(db)
	} else {
		err = d.update(db)
	}

	return err
}

// called by d.Persist, never call directly
func (d *Domain) add(db *sqlx.DB) error {

	sFields := ""
	var params []interface{}

	if d.isNew {
		sFields += "domain"
		params = append(params, d.Domain)
	}

	if d.isDescDirty {
		if len(sFields) > 0 {
			sFields += ", "
		}
		sFields += "desc"
		params = append(params, d.Description.String)
	}

	if d.isIsActiveDirty {
		if len(sFields) > 0 {
			sFields += ", "
		}
		sFields += "active"
		params = append(params, d.GetIsActive())
	}

	// generate param string, remove last , from repeater
	sQM := strings.Repeat("?,", len(params))
	sQM = sQM[:len(sQM)-1]

	stmt, err := db.Preparex("INSERT INTO domain (" + sFields + ") VALUES (" + sQM + ")")
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

	fields := logrus.Fields{"domain": d.Domain, "description": d.Description, "active": d.IsActive}

	d.clearDirtyFlags()

	if count == 0 {
		common.LogInfo("Nothing done.", fields)
	} else {
		common.LogInfo("Domain added.", fields)
	}

	return nil
}

// called by a.Persist, never call directly
func (d *Domain) update(db *sqlx.DB) error {

	if !d.IsDirty() {
		return errors.New("trying to update unchanged object")
	}

	sStatement := ""
	var params []interface{}

	// d.Description, d.IsActive, time.Now(), a.UpdDat
	if d.isDescDirty {
		if len(sStatement) > 0 {
			sStatement += ", "
		}
		sStatement += "desc = ?"
		params = append(params, d.Description.String)
	}

	if d.isIsActiveDirty {
		if len(sStatement) > 0 {
			sStatement += ", "
		}
		sStatement += "active = ?"
		params = append(params, d.GetIsActive())
	}

	// update upddat field
	if len(sStatement) > 0 {
		sStatement += ", "
	}
	sStatement += "upd_dat = ?"
	params = append(params, time.Now())

	// append params for where
	params = append(params, d.Domain) // pkey
	params = append(params, d.UpdDat) // optimistic locking

	stmt, err := db.Preparex("UPDATE domain SET " + sStatement + " WHERE domain = ? AND upd_dat <= ?")
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

	fields := logrus.Fields{"domain": d.Domain, "description": d.Description, "active": d.IsActive}

	d.clearDirtyFlags()

	if count == 0 {
		common.LogInfo("Nothing done.", fields)
	} else {
		common.LogInfo("Alias updated.", fields)
	}

	return nil
}

func (d Domain) Delete() error {

	return DeleteDomain(d.Domain)
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
