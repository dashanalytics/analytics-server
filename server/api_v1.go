// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package server

import (
	"context"
	"encoding/json"
	"github.com/jellyterra/go-httpform"
	"net/http"
)

func V1GetApi(e *ApiEnv) map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/reportAccess": Wrap(func(w http.ResponseWriter, r *http.Request) (int, error) {
			wrap, err := httpform.WrapFromRequest(r)
			if err != nil {
				return http.StatusBadRequest, err
			}

			var (
				uuid       = wrap.String("uuid", "")
				deployTime = wrap.StringRequired("deploy_time")
				target     = wrap.StringRequired("target")
			)
			err = wrap.Parse()
			if err != nil {
				return http.StatusBadRequest, err
			}

			var (
				country  = r.Header.Get(e.HeaderKeyForIPCountry)
				sourceIP = r.Header.Get(e.HeaderKeyForConnectingIP)
			)
			err = e.Database.AddAccessReport(context.Background(), &AccessReport{
				SourceIP:   sourceIP,
				Country:    country,
				UUID:       *uuid,
				DeployTime: *deployTime,
				Target:     *target,
			})
			if err != nil {
				return http.StatusInternalServerError, err
			}

			return 0, nil
		}),

		"/getAccessReportTimestamps": Wrap(func(w http.ResponseWriter, r *http.Request) (int, error) {
			ctx := r.Context()

			wrap, err := httpform.WrapFromRequest(r)
			if err != nil {
				return http.StatusBadRequest, err
			}

			var (
				token = wrap.StringRequired("token")
				start = wrap.StringRequired("start")
				end   = wrap.StringRequired("end")
			)

			err = wrap.Parse()
			if err != nil {
				return http.StatusBadRequest, err
			}

			if *token != e.AccessToken {
				return http.StatusForbidden, nil
			}

			timestamps, err := e.Database.GetAccessReportsTimestamps(ctx, *start, *end)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			b, err := json.Marshal(timestamps)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			_, _ = w.Write(b)

			return 0, nil
		}),

		"/getAccessReportByTimestamp": Wrap(func(w http.ResponseWriter, r *http.Request) (int, error) {
			ctx := r.Context()

			wrap, err := httpform.WrapFromRequest(r)
			if err != nil {
				return http.StatusBadRequest, err
			}

			var (
				token     = wrap.StringRequired("token")
				timestamp = wrap.StringRequired("timestamp")
			)

			err = wrap.Parse()
			if err != nil {
				return http.StatusBadRequest, err
			}

			if *token != e.AccessToken {
				return http.StatusForbidden, nil
			}

			report, err := e.Database.GetAccessReportByTimestamp(ctx, *timestamp)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			b, err := json.Marshal(&report)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			_, _ = w.Write(b)

			return 0, nil
		}),

		"/getAccessReportsByRange": Wrap(func(w http.ResponseWriter, r *http.Request) (int, error) {
			ctx := r.Context()

			wrap, err := httpform.WrapFromRequest(r)
			if err != nil {
				return http.StatusBadRequest, err
			}

			var (
				token = wrap.StringRequired("token")
				start = wrap.StringRequired("start")
				end   = wrap.StringRequired("end")
			)

			err = wrap.Parse()
			if err != nil {
				return http.StatusBadRequest, err
			}

			if *token != e.AccessToken {
				return http.StatusForbidden, nil
			}

			timestamps, err := e.Database.GetAccessReportsTimestamps(ctx, *start, *end)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			reportsMap := map[string]AccessReport{}

			for _, timestamp := range timestamps {
				report, err := e.Database.GetAccessReportByTimestamp(ctx, timestamp)
				if err != nil {
					return http.StatusInternalServerError, err
				}

				reportsMap[timestamp] = report
			}

			b, err := json.Marshal(reportsMap)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			_, _ = w.Write(b)

			return 0, nil
		}),
	}
}
