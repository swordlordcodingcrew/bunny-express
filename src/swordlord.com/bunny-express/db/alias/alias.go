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
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"swordlord.com/bunny-express/common"
	"swordlord.com/bunny-express/db"
	"time"
)

type Alias struct {
	Alias                 string         `db:"alias"`
	Description           sql.NullString `db:"desc"`
	isDescDirty           bool
	Domain                string `db:"domain"`
	isDomainDirty         bool
	ForwardAddress        string `db:"forward_address"`
	isForwardAddressDirty bool
	IsActive              bool `db:"active"`
	isIsActiveDirty       bool
	// tells us if object is from db or not
	isNew  bool
	CrtDat time.Time `db:"crt_dat"`
	UpdDat time.Time `db:"upd_dat"`
}

func NewAlias() *Alias {

	a := &Alias{}
	a.clearDirtyFlags()
	a.isNew = true
	a.CrtDat = time.Now()
	a.UpdDat = time.Now()

	return a
}

func (a *Alias) clearDirtyFlags() {
	a.isDescDirty = false
	a.isDomainDirty = false
	a.isForwardAddressDirty = false
	a.isIsActiveDirty = false
}

func (a *Alias) GetAlias() string               { return a.Alias }
func (a *Alias) GetDescription() sql.NullString { return a.Description }
func (a *Alias) GetDomain() string              { return a.Domain }
func (a *Alias) GetForwardAddress() string      { return a.ForwardAddress }
func (a *Alias) GetIsActive() bool              { return a.IsActive }

func (a *Alias) SetAlias(aliass string) {
	a.Alias = aliass
}

func (a *Alias) SetDescription(description sql.NullString) {

	// TODO: add sanity check to all SetXY functions
	if a.Description.String == description.String {
		return
	}

	a.Description = description
	a.isDescDirty = true
}

func (a *Alias) SetDomain(domain string) {
	a.Domain = domain
	a.isDomainDirty = true
}

func (a *Alias) SetForwardAddress(fa string) {

	if a.ForwardAddress == fa {
		return
	}

	a.ForwardAddress = fa
	a.isForwardAddressDirty = true
}

func (a *Alias) SetIsActive(ia bool) {

	if a.IsActive == ia {
		return
	}

	a.IsActive = ia
	a.isIsActiveDirty = true
}

func (a *Alias) IsDirty() bool {
	if a.isDescDirty || a.isDomainDirty || a.isForwardAddressDirty || a.isIsActiveDirty {
		return true
	} else {
		return false
	}
}

type AliasFilter struct {
	Alias          string
	Description    string
	Domain         string
	ForwardAddress string
	IsActive       sql.NullBool
}

func GetFieldCaptions() []string {

	captions := []string{"Alias", "Description", "Domain", "Forward", "Active", "Created", "Updated"}

	return captions
}

func GetAllAliases() ([]Alias, error) {

	return GetFilteredAliases(&AliasFilter{})
}

func GetFilteredAliases(af *AliasFilter) ([]Alias, error) {

	db, err := db.OpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	sFilter := ""
	params := []string{}

	if len(af.Domain) > 0 {
		sFilter += "domain LIKE ?"
		params = append(params, af.Domain)
	}

	if len(af.ForwardAddress) > 0 {
		if len(sFilter) > 0 {
			sFilter += " AND "
		}
		sFilter += "forward_address LIKE ?"
		params = append(params, af.ForwardAddress)
	}

	if af.IsActive.Valid {
		if len(sFilter) > 0 {
			sFilter += " AND "
		}
		sFilter += "active LIKE ?"

		params = append(params, strconv.FormatBool(af.IsActive.Bool))
	}

	if len(sFilter) > 0 {
		sFilter = "WHERE " + sFilter
	}

	sql := "SELECT * FROM alias " + sFilter + " ORDER BY domain, alias ASC"

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
	var a []Alias
	err = stmt.Select(&a, args...)

	if err == nil {
		for i := range a {
			a[i].isNew = false
		}
	}

	return a, err
}

func GetAlias(name string) (*Alias, error) {

	db, err := db.OpenDB()
	if err != nil {
		return NewAlias(), err
	}
	defer db.Close()

	stmt, err := db.Preparex(db.Rebind("SELECT * FROM alias WHERE alias=?"))
	if err != nil {
		return NewAlias(), err
	}

	a := NewAlias()
	err = stmt.Get(a, name)
	if err != nil {
		return NewAlias(), err
	} else {
		a.isNew = false
		return a, nil
	}
}

func (a *Alias) Persist() error {

	db, err := db.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	if !a.IsDirty() && !a.isNew {
		common.LogInfo("Alias did not change, not persisted.", nil)
		return nil
	}

	if a.isNew {
		err = a.add(db)
	} else {
		err = a.update(db)
	}

	return err
}

// called by a.Persist, never call directly
func (a *Alias) add(db *sqlx.DB) error {

	sFields := ""
	var params []interface{}

	if a.isNew {
		sFields += "alias"
		params = append(params, a.Alias)
	}

	// a.ForwardAddress, a.Description, a.IsActive, time.Now(), a.Alias, a.UpdDat
	if a.isDomainDirty {
		if len(sFields) > 0 {
			sFields += ", "
		}
		sFields += "domain"
		params = append(params, a.Domain)
	}

	if a.isDescDirty {
		if len(sFields) > 0 {
			sFields += ", "
		}
		sFields += "desc"
		params = append(params, a.Description.String)
	}

	if a.isForwardAddressDirty {
		if len(sFields) > 0 {
			sFields += ", "
		}
		sFields += "forward_address"
		params = append(params, a.ForwardAddress)
	}

	if a.isIsActiveDirty {
		if len(sFields) > 0 {
			sFields += ", "
		}
		sFields += "active"
		params = append(params, a.GetIsActive())
	}

	if len(sFields) > 0 {
		sFields += ", "
	}
	sFields += "crt_dat, upd_dat"
	params = append(params, a.CrtDat)
	params = append(params, a.UpdDat)

	// generate param string, remove last , from repeater
	sQM := strings.Repeat("?,", len(params))
	sQM = sQM[:len(sQM)-1]

	stmt, err := db.Preparex("INSERT INTO alias (" + sFields + ") VALUES (" + sQM + ")")
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

	fields := logrus.Fields{"alias": a.Alias, "domain": a.Domain, "forward": a.ForwardAddress, "description": a.Description, "active": a.IsActive}

	a.clearDirtyFlags()

	if count == 0 {
		common.LogInfo("Nothing done.", fields)
	} else {
		common.LogInfo("Alias added.", fields)
	}

	return nil
}

// called by a.Persist, never call directly
func (a *Alias) update(db *sqlx.DB) error {

	if !a.IsDirty() {
		return errors.New("trying to update unchanged object")
	}

	sStatement := ""
	var params []interface{}

	// a.ForwardAddress, a.Description, a.IsActive, time.Now(), a.Alias, a.UpdDat
	if a.isDomainDirty {
		sStatement += "domain = ?"
		params = append(params, a.Domain)
	}

	if a.isDescDirty {
		if len(sStatement) > 0 {
			sStatement += ", "
		}
		sStatement += "desc = ?"
		params = append(params, a.Description.String)
	}

	if a.isForwardAddressDirty {
		if len(sStatement) > 0 {
			sStatement += ", "
		}
		sStatement += "forward_address = ?"
		params = append(params, a.ForwardAddress)
	}

	if a.isIsActiveDirty {
		if len(sStatement) > 0 {
			sStatement += ", "
		}
		sStatement += "active = ?"
		params = append(params, a.GetIsActive())
	}

	// update upddat field
	if len(sStatement) > 0 {
		sStatement += ", "
	}
	sStatement += "upd_dat = ?"
	params = append(params, time.Now())

	// append params for where
	params = append(params, a.Alias)  // pkey
	params = append(params, a.UpdDat) // optimistic locking

	stmt, err := db.Preparex("UPDATE alias SET " + sStatement + " WHERE alias = ? AND upd_dat <= ?")
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

	fields := logrus.Fields{"alias": a.Alias, "domain": a.Domain, "forward": a.ForwardAddress, "description": a.Description, "active": a.IsActive}

	a.clearDirtyFlags()

	if count == 0 {
		common.LogInfo("Nothing done.", fields)
	} else {
		common.LogInfo("Alias updated.", fields)
	}

	return nil
}

func (a Alias) Delete() error {

	return DeleteAlias(a.Alias)
}

func DeleteAlias(alias string) error {

	db, err := db.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Preparex(`DELETE FROM alias WHERE alias=?`)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(alias)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		common.LogInfo("Nothing deleted. Wrong Alias used?", logrus.Fields{"alias": alias})
	} else {
		common.LogInfo("Alias deleted.", logrus.Fields{"alias": alias, "count": count})
	}

	return nil
}

func FillDefaultAliasOnDomain(domain string) error {

	aliases := common.GetStringSliceFromConfig("default.alias")

	for _, an := range aliases {

		alias := NewAlias()
		alias.SetDomain(domain)
		alias.SetAlias(an + "@" + domain)
		alias.SetForwardAddress("root@" + domain)
		alias.SetIsActive(true)

		var desc sql.NullString
		desc.Scan("filled automatically with default alias from config")
		alias.SetDescription(desc)

		err := alias.Persist()
		if err != nil {
			common.LogInfo("AddAlias returned an error.", logrus.Fields{"alias": an, "error": err})
			return err
		}
	}

	return nil
}
