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
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/markbates/pkger"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"pass-go/config"
	"pass-go/storage"
	"path/filepath"
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
	"604800", // Week
	"86400", // Day
	"3600":
		return nil
	}
	return  errors.New("invalid ttl")
}
func parseFiles(t *template.Template, filenames ...string) (*template.Template, error) {
	if len(filenames) == 0 {
		// Not really a problem, but be consistent.
		return nil, fmt.Errorf("template: no files named in call to ParseFiles")
	}
	for _, filename := range filenames {
		f , err := pkger.Open(filename)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}
		s := string(b)
		name := filepath.Base(filename)
		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}


func New(conf config.Config, storage *storage.Storage) chi.Router {
	// TODO: Parse all templates at program start.
	r := chi.NewRouter()
	r.Route(conf.HTTP.Root, func(root chi.Router) {
		root.Get("/", func(w http.ResponseWriter, r *http.Request) {
			if tpl, err := parseFiles(nil,
				"/templates/set_password.html",
				"/templates/base.html"); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				if err = tpl.ExecuteTemplate(w, "set_password.html", nil); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
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

			tpl, _ := parseFiles(nil,
				"/templates/confirm.html",
				"/templates/base.html")
			_ = tpl.ExecuteTemplate(w, "confirm.html", u.String())
		})

		root.Get("/{password_key}", func(w http.ResponseWriter, r *http.Request) {
			passwordKey := chi.URLParam(r, "password_key")
			passwordKey, _ = url.QueryUnescape(passwordKey)
			if !storage.PasswordExists(passwordKey) {
				http.NotFound(w,r)
				return
			}
			tpl, _ := parseFiles(nil,
				"/templates/preview.html",
				"/templates/base.html")
			_ = tpl.ExecuteTemplate(w, "preview.html", nil)
		})
		root.Post("/{password_key}", func(w http.ResponseWriter, r *http.Request) {
			passwordKey, _ := url.QueryUnescape(chi.URLParam(r, "password_key"))
			password := storage.GetPassword(passwordKey)
			if password == "" {
				http.NotFound(w,r)
				return
			}
			tpl, _ := parseFiles(nil,
				"/templates/password.html",
				"/templates/base.html")
			_ = tpl.ExecuteTemplate(w, "password.html", password)
		} )

		fs := http.FileServer(pkger.Dir("/static/"))
		root.Handle("/static/*", http.StripPrefix("/static/", fs))
	})
	return r
}
