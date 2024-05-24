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

package config

type (
	// Config holds the application configuration
	Config struct {
		HTTP struct {
			Hosts string `envconfig:"PASSGO_HTTP_HOSTS"`
			Port  string `envconfig:"PASSGO_HTTP_PORT" default:"5000"`
			Root  string `envconfig:"PASSGO_HTTP_ROOT" default:"/"`
			Email string `envconfig:"PASSGO_LETSENCRYPT_EMAIL"`
		}
		DB struct {
			DataDir string `envconfig:"PASSGO_DB_DATADIR" default:"/var/lib/ledis"`
		}
		CacheDir string `default:"/var/lib/acme"`
		Captcha  bool   `envconfig:"PASSGO_CAPTCHA" default:"false"`
	}
)
