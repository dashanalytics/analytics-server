// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

type Config struct {
	Key  string `yaml:"key"`
	Cert string `yaml:"cert"`

	Listen   string `yaml:"listen"`
	Database string `yaml:"db"`

	AccessToken string `yaml:"access_token"`

	Header struct {
		Key struct {
			ConnectingIP string `yaml:"connecting_ip"`
		} `yaml:"key"`
	} `yaml:"header"`
}
