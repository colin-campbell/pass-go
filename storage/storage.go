/*
 * Copyright (c) 2020. Colin Stewart Campbell <colin.campbell@digitalistgroup.com>
 *  This file is part of Pass-Go.
 *
 *      Pass-Go is free software: you can redistribute it and/or modify
 *      it under the terms of the GNU General Public License as published by
 *      the Free Software Foundation, either version 3 of the License, or
 *      (at your option) any later version.
 *
 *      Pass-Go is distributed in the hope that it will be useful,
 *      but WITHOUT ANY WARRANTY; without even the implied warranty of
 *      MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *      GNU General Public License for more details.
 *
 *      You should have received a copy of the GNU General Public License
 *      along with Pass-Go.  If not, see <https://www.gnu.org/licenses/>.
 */

package storage

import (
	lediscfg "github.com/ledisdb/ledisdb/config"
	"github.com/ledisdb/ledisdb/ledis"
	uuid "github.com/nu7hatch/gouuid"
	"pass-go/config"
	"pass-go/crypto"
	"strconv"
	"strings"
)

type Storage struct {
	db *ledis.DB
	crypto *crypto.Crypto
}

func New(conf config.Config) *Storage {
	s := &Storage{}
	dbCfg := lediscfg.NewConfigDefault()
	dbCfg.DataDir = conf.DB.DataDir
	if l, err := ledis.Open(dbCfg); err != nil {
		panic("Unable to open database")
	} else {
		s.db, _ = l.Select(0)
	}
	s.crypto = crypto.New()
	return s
}

func (s *Storage) SetPassword(password string, ttl string) string {
	u, _ := uuid.NewV4()
	// Compatibility with SnapPass uuid.
	storageKey := strings.Replace(u.String(), "-", "", -1)
	cipherText, encryptionKey := s.crypto.Encrypt(password)
	duration, _  := strconv.Atoi(ttl)
	_ = s.db.SetEX([]byte(storageKey), int64(duration), []byte(cipherText))
	token := strings.Join([]string{storageKey, encryptionKey}, "~")
	return token
}

func (s *Storage) GetPassword(token string) string {
	storageKey, decryptionKey := s.parseToken(token)
	password, _  := s.db.Get([]byte(storageKey))
	_, _ = s.db.Del([]byte(storageKey))
	if password != nil {
		decoded := s.crypto.Decrypt(string(password), decryptionKey, 0)
		return decoded
	}
	return ""
}

func (s *Storage) PasswordExists(token string) bool {
	storageKey, _ := s.parseToken(token)
	exists, _ := s.db.Exists([]byte(storageKey))
	return exists == 1
}

func (s *Storage) parseToken(token string) (string ,string){
		parts := strings.Split(token, "~")
		return parts[0], parts[1]
}