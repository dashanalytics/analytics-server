// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"os/signal"
)

var licenseText = `
Copyright 2024 Jelly Terra
This program and its Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0 that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.
`

func main() {
	err := _main()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func _main() error {

	flag.BoolFunc("license", "Print license information.", func(_ string) error {
		fmt.Println(licenseText)
		return nil
	})

	var (
		configPath = flag.String("c", "analytics-server.yaml", "Path to config file.")
	)

	flag.Parse()

	body, err := os.ReadFile(*configPath)
	if err != nil {
		return err
	}

	var config Config

	err = yaml.Unmarshal(body, &config)
	if err != nil {
		return err
	}

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	err = Serve(ctx, &config)
	if err != nil {
		return err
	}

	return nil
}
