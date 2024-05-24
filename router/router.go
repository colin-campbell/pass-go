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

package router

import (
	"embed"
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"pass-go/config"
	"pass-go/storage"
	"strings"

	"github.com/go-chi/chi"
	"github.com/leonelquinteros/gotext"
	"golang.org/x/text/language"
)

var (
	StaticFS embed.FS
	Temps    template.Template
)

func checkInput(form map[string][]string) error {
	if len(form["password"][0]) == 0 {
		return errors.New("password is required")
	}
	if len(form["ttl"][0]) == 0 {
		return errors.New("ttl is required")
	}

	switch form["ttl"][0] {
	case
		"1209600", // Two weeks
		"604800",  // Week
		"86400",   // Day
		"3600":
		return nil
	}
	return errors.New("invalid ttl")
}

var (
	matcher = language.NewMatcher([]language.Tag{
		language.English,
		language.Swedish,
		language.Spanish,
		language.Bengali,
		language.German,
	})
)

func setLanguage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, _, _ := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
		tag, _, _ := matcher.Match(t...)
		gotext.Configure("locales", tag.String(), "messages")
		next.ServeHTTP(w, r)
	})
}

// blockUserAgents stops Slackbot parsing the secret.
func blockUserAgents(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userAgent := r.Header.Get("User-Agent")
		if strings.HasPrefix(userAgent, "Slack") {
			http.Error(w, "Come in Harold, your tea is ready!", http.StatusTeapot)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// New creates a new router with handlers
func New(conf config.Config, storage *storage.Storage) chi.Router {
	// TODO: Parse all templates at program start.
	r := chi.NewRouter()

	r.Use(blockUserAgents)
	r.Use(setLanguage)

	r.Route(conf.HTTP.Root, func(root chi.Router) {
		root.Get("/", func(w http.ResponseWriter, r *http.Request) {

			if err := Temps.ExecuteTemplate(w, "set_password.html", nil); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		})

		root.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
			const robot = `
User-agent: *
Disallow: /
`
			w.Write([]byte(robot))
		})

		root.Post("/", func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid form data", http.StatusBadRequest)
				return
			}
			if err := checkInput(r.Form); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			scheme := "http://"
			if r.TLS != nil {
				scheme = "https://"
			}
			token := storage.SetPassword(r.Form["password"][0], r.Form["ttl"][0])
			full := scheme + r.Host + "/" + token
			u, _ := url.Parse(full)

			err := Temps.ExecuteTemplate(w, "confirm.html", u.String())
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

		})

		root.Get("/{password_key}", func(w http.ResponseWriter, r *http.Request) {
			passwordKey := chi.URLParam(r, "password_key")
			passwordKey, _ = url.QueryUnescape(passwordKey)
			if !storage.PasswordExists(passwordKey) {
				http.NotFound(w, r)
				return
			}

			_ = Temps.ExecuteTemplate(w, "preview.html", nil)
		})
		root.Post("/{password_key}", func(w http.ResponseWriter, r *http.Request) {
			passwordKey, _ := url.QueryUnescape(chi.URLParam(r, "password_key"))
			password := storage.GetPassword(passwordKey)
			if password == "" {
				http.NotFound(w, r)
				return
			}
			_ = Temps.ExecuteTemplate(w, "password.html", password)
		})

		fs := http.FileServer(http.FS(StaticFS))
		root.Handle("/static/*", fs)
	})
	return r
}
