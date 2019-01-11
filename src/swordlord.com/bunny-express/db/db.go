package db

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
	"github.com/jmoiron/sqlx"
	"log"
	"swordlord.com/bunny-express/common"
)

var createDomainTbl = `
CREATE TABLE domain (
  domain varchar(255) PRIMARY KEY,
  desc varchar(2000),
  active bool DEFAULT true,
  crt_dat timestamp DEFAULT CURRENT_TIMESTAMP,
  upd_dat timestamp DEFAULT CURRENT_TIMESTAMP
);`

var createMailboxTbl = `
CREATE TABLE mailbox (
  mail varchar(255) PRIMARY KEY,
  desc varchar(2000),
  domain varchar(255) NOT NULL,
  pwd varchar(2000) NOT NULL,
  mail_dir varchar(2000) NOT NULL,
  local_part varchar(500) NOT NULL,
  relay_domain varchar(500),
  quota varchar(100),
  active bool DEFAULT true,
  crt_dat timestamp DEFAULT CURRENT_TIMESTAMP,
  upd_dat timestamp DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT mailbox_domain_fk FOREIGN KEY (domain) REFERENCES domain (domain)
);`

var createAliasTbl = `
CREATE TABLE alias (
  alias varchar(255) PRIMARY KEY,
  desc varchar(2000),
  domain varchar(255) NOT NULL,
  forward_address varchar(255) NOT NULL,
  active bool DEFAULT true,
  crt_dat timestamp DEFAULT CURRENT_TIMESTAMP,
  upd_dat timestamp DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT alias_domain_fk FOREIGN KEY (domain) REFERENCES domain (domain)
);`

var checkTblExists = `SELECT COUNT(name) FROM sqlite_master WHERE type='table' AND tbl_name=?;`

func CheckDatabase() {

	checkTables()

	checkDemoData()
}

func OpenDB() (*sqlx.DB, error) {

	return sqlx.Open(getDatabaseDriver(), getDatabaseName())
}

func getDatabaseDriver() string {
	return "sqlite3"
}

func getDatabaseName() string {

	dbName := common.GetStringFromConfig("db.file")
	if dbName == "" {

		return "be.sqlite"
	} else {

		return dbName
	}
}

func checkDemoData() {

	insertDemoData := common.GetBoolFromConfig("db.add_demo_data", false)
	if !insertDemoData {
		return
	}

	db, err := OpenDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	runStatement(db, "INSERT INTO domain (domain, active) VALUES ('demo1.com', true)")
	runStatement(db, "INSERT INTO domain (domain, active) VALUES ('demo2.com', true)")
	runStatement(db, "INSERT INTO domain (domain, active) VALUES ('demo3.com', true)")
	runStatement(db, "INSERT INTO mailbox (mail, domain, pwd, mail_dir, local_part) VALUES ('a@demo1.com', 'demo1.com', 'pwd', '/var/mail/', 'demo1.com')")
	runStatement(db, "INSERT INTO mailbox (mail, domain, pwd, mail_dir, local_part) VALUES ('a@demo2.com', 'demo2.com', 'pwd', '/var/mail/', 'demo2.com')")
	runStatement(db, "INSERT INTO alias (alias, forward_address, domain) VALUES ('alias@demo2.com', 'a@demo2.com', 'demo2.com')")
}

func runStatement(db *sqlx.DB, s string) {

	_, err := db.Exec(s)
	if err != nil {
		common.LogErrorFmt("Could not add Demo Data for statement '%s' with error '%s'", s, err)
	}
}

func checkTables() {

	db, err := OpenDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	err = checkTable(db, "domain", createDomainTbl)
	if err != nil {
		log.Fatalln(err)
	}

	err = checkTable(db, "mailbox", createMailboxTbl)
	if err != nil {
		log.Fatalln(err)
	}

	err = checkTable(db, "alias", createAliasTbl)
	if err != nil {
		log.Fatalln(err)
	}

}

func checkTable(db *sqlx.DB, name string, sqlCrt string) error {

	var exists bool

	err := db.QueryRow(checkTblExists, name).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		_, err := db.Exec(sqlCrt)
		if err != nil {
			return err
		}
	}

	return nil
}
