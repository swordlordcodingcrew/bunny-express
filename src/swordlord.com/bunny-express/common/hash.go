package common

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
 ** The MD5-Crypt code is based on source code written by Damian Gryski <damian@gryski.com>
 **
 ** Which is based on the implementation at
 ** http://code.activestate.com/recipes/325204-passwd-file-compatible-1-md5-crypt/
 **
 ** Both are licensed as:
 **   * "THE BEER-WARE LICENSE" (Revision 42):
 **   * <phk@login.dknet.dk> wrote this file.  As long as you retain this notice you
 **   * can do whatever you want with this stuff. If we meet some day, and you think
 **   * this stuff is worth it, you can buy me a beer in return.   Poul-Henning Kamp
 **
-----------------------------------------------------------------------------*/

import (
	"crypto/md5"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

const p64alphabet = "./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

var md5permute [5][3]int

func md5Init() {
	md5permute = [5][3]int{
		[3]int{0, 6, 12},
		[3]int{1, 7, 13},
		[3]int{2, 8, 14},
		[3]int{3, 9, 15},
		[3]int{4, 10, 5},
	}
}

func HashPasswordBCrypt(pwd string) (string, error) {

	password := []byte(pwd)

	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func HashPasswordMD5Crypt(pwd string, salt string) (string, error) {
	md5Init()
	hash, err := md5Crypt([]byte(pwd), []byte(salt))

	return string(hash), err
}

func CheckHashedPasswordBCrypt(hashedPassword string, password string) error {

	pwd := []byte(password)
	hashedPwd := []byte(hashedPassword)

	// Comparing the password with the hash
	err := bcrypt.CompareHashAndPassword(hashedPwd, pwd)

	// nil means it is a match
	return err
}

func CheckHashedPasswordMD5Crypt(hashedPasswordWSalt string, password string) error {
	md5Init()

	return errors.New("")
}

func md5Pass64(b []byte) []byte {
	// not quite base64 encoding
	// 1) bits are encoded in the wrong order
	// 2) the alphabet is different

	pass := make([]byte, 0, (len(b)+1*4)/3)

	for _, v := range md5permute {

		v := int(b[v[0]])<<16 | int(b[v[1]])<<8 | int(b[v[2]])
		for j := 0; j < 4; j++ {
			pass = append(pass, p64alphabet[v&0x3f])
			v >>= 6
		}
	}
	v := b[11]
	pass = append(pass, p64alphabet[v&0x3f])
	v >>= 6
	pass = append(pass, p64alphabet[v&0x3f])

	return pass
}

// Crypt hashes the plaintext password using the salt from the hashed password.
func md5Crypt(plain []byte, salt []byte) ([]byte, error) {

	m := md5.New()
	m.Write(plain)
	m.Write(salt)
	m.Write(plain)
	final := m.Sum(nil)

	m.Reset()
	m.Write(plain)
	m.Write([]byte("$1$"))
	m.Write(salt)

	for idx := len(plain); idx > 0; idx -= 16 {
		if idx > 16 {
			m.Write(final[:16])
		} else {
			m.Write(final[:idx])
		}
	}

	var ctx []byte
	for i := len(plain); i > 0; i >>= 1 {
		if i&1 == 1 {
			ctx = append(ctx, 0)
		} else {
			ctx = append(ctx, plain[0])
		}
	}

	m.Write(ctx)
	final = m.Sum(nil)

	for i := 0; i < 1000; i++ {
		m.Reset()

		if i&1 == 1 {
			m.Write(plain)
		} else {
			m.Write(final[:16])
		}

		if i%3 != 0 {
			m.Write(salt)
		}

		if i%7 != 0 {
			m.Write(plain)
		}

		if i&1 == 1 {
			m.Write(final[:16])
		} else {
			m.Write(plain)
		}

		final = m.Sum(nil)
	}

	var passwd []byte
	passwd = append(passwd, []byte("$1$")...)
	passwd = append(passwd, salt...)
	passwd = append(passwd, '$')
	passwd = append(passwd, md5Pass64(final)...)

	return passwd, nil
}
