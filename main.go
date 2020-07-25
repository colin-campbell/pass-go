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

//go:generate pkger
package main

import (
	"crypto/tls"
	"github.com/markbates/pkger"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"pass-go/config"
	"pass-go/router"
	"pass-go/signals"
	"pass-go/storage"
)

func main() {
	// Handle our own signals for when we run as PID 1
	signals.Setup()

	// Tell "pkger" tool  that it must package templates directory.
	// @see go:generate in the file header.
	_ = pkger.Include("/templates")
	conf := config.MustLoad()
	store := storage.New(conf)

	r := router.New(conf, store)

	// If hosts are configured, setup LetsEncrypt and listen on 80 & 443
	if conf.HTTP.Hosts != "" {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(conf.HTTP.Hosts),
			Cache:      autocert.DirCache(conf.CacheDir),
		}
		// Good manners to supply an email.
		if conf.HTTP.Email != "" {
			certManager.Email = conf.HTTP.Email
		}
		server := &http.Server{
			Addr: ":https",
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
				NextProtos:     []string{acme.ALPNProto},
			},
			// Chi router
			Handler: r,
		}
		// With TLS, add auto-redirect 80->443
		go func() {
			h := certManager.HTTPHandler(nil)
			log.Fatal(http.ListenAndServe(":http", h))
		}()

		log.Fatal(server.ListenAndServeTLS("",""))
	} else {
		// Just listen for plain old http (if behind a LB/SSL terminator)
		log.Fatal(http.ListenAndServe(":" + conf.HTTP.Port, r))
	}


}
