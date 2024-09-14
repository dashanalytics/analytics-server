// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package server

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type Database struct {
	Client *redis.Client
}

func (db *Database) AddAccessReport(ctx context.Context, report *AccessReport) error {

	var (
		now       = time.Now().UTC()
		timestamp = now.Format("2006:01:02:15:04:05.99")
		score, _  = strconv.ParseFloat(now.Format("20060102150405.99"), 64)
	)

	err := db.Client.HSet(ctx, "AccessReports:"+timestamp, *report).Err()
	if err != nil {
		return err
	}

	err = db.Client.ZAdd(ctx, "AccessReports:Timestamps", redis.Z{Score: score, Member: timestamp}).Err()
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetAccessReportsTimestamps(ctx context.Context, start, end string) ([]string, error) {
	return db.Client.ZRangeByScore(ctx, "AccessReports:Timestamps", &redis.ZRangeBy{Min: start, Max: end}).Result()
}

func (db *Database) GetAccessReportByTimestamp(ctx context.Context, timestamp string) (report AccessReport, _ error) {
	return report, db.Client.HGetAll(ctx, "AccessReports:"+timestamp).Scan(&report)
}
