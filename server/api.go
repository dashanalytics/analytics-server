// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package server

import (
	"fmt"
	"log"
	"net/http"
)

type ApiEnv struct {
	Database *Database

	AccessToken string

	HeaderKeyForConnectingIP,
	HeaderKeyForIPCountry string
}

type AccessReport struct {
	SourceIP string `redis:"SourceIP" json:"source_ip"`
	Country  string `redis:"Country" json:"country"`

	UUID       string `redis:"UUID" json:"uuid"`
	DeployTime string `redis:"DeployTime" json:"deploy_time"`
	Target     string `redis:"Target" json:"target"`
}

func Wrap(h func(w http.ResponseWriter, r *http.Request) (int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := h(w, r)
		if status == 0 {
			return
		}

		w.WriteHeader(status)
		if err != nil {
			switch status {
			case http.StatusInternalServerError:
				log.Println(err)
			case http.StatusBadRequest:
				_, _ = fmt.Fprintln(w, err)
			}
		}
	}
}

func Register(mux *http.ServeMux, prefix string, table map[string]http.HandlerFunc) {
	for path, handler := range table {
		mux.HandleFunc(prefix+path, handler)
	}
}
