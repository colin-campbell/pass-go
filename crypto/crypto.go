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

package crypto

import (
	"github.com/fernet/fernet-go"
	"time"
)
type Crypto struct {

}

func (c *Crypto) Encrypt(plainText string) (string, string) {
	k := &fernet.Key{}
	_ = k.Generate()
	cipherText, _ :=  fernet.EncryptAndSign([]byte(plainText), k)
	return string(cipherText), k.Encode()
}

func (c *Crypto) Decrypt(cipherText string, key string, ttl time.Duration) string {
	k := fernet.MustDecodeKeys(key)
	plainText := fernet.VerifyAndDecrypt([]byte(cipherText), ttl, k )
	return string(plainText)
}

func New() *Crypto {
	c := &Crypto{}
	return c
}