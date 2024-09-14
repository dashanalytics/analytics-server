// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package server

import "net/http"

var APIVersions = map[string]func(*ApiEnv) map[string]http.HandlerFunc{
	"v1": V1GetApi,
}
