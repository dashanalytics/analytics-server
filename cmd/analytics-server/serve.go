// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"errors"
	"github.com/dashanalytics/analytics-server/server"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
)

func Serve(ctx context.Context, config *Config) error {

	options, err := redis.ParseURL(config.Database)
	if err != nil {
		return err
	}

	var (
		client = redis.NewClient(options)

		db = &server.Database{
			Client: client,
		}

		mux = http.NewServeMux()

		apiEnv = &server.ApiEnv{
			Database: db,

			AccessToken: config.AccessToken,

			HeaderKeyForConnectingIP: config.Header.Key.ConnectingIP,
		}
	)

	for version, get := range server.APIVersions {
		server.Register(mux, "/api/"+version, get(apiEnv))
	}

	log.Println("Test connection to the Redis database.")
	err = client.Ping(ctx).Err()
	if err != nil {
		return err
	}

	serverConfig := &http.Server{
		Addr:           config.Listen,
		Handler:        mux,
		MaxHeaderBytes: 4096,
	}

	go func() {
		<-ctx.Done()
		log.Println("Shutting down.")

		err = serverConfig.Shutdown(context.Background())
		if err != nil {
			log.Println("Occurred an error while shutting down:", err)
		}
	}()

	switch {
	case config.Key == "" || config.Cert == "":
		log.Println("Listen on http://" + config.Listen)

		err = serverConfig.ListenAndServe()
	default:
		log.Println("Listen on https://" + config.Listen)

		err = serverConfig.ListenAndServeTLS(config.Cert, config.Key)
	}
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
