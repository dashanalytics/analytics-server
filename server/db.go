// Copyright 2024 Jelly Terra
// This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0
// that can be found in the LICENSE file and https://mozilla.org/MPL/2.0/.

package server

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"strings"
	"time"
)

const (
	kAccessReports_           = "AccessReports:"
	kAccessReports_Timestamps = "AccessReports:Timestamps"
	kAccessCounts_Targets     = "AccessCounts:Targets"
	kAccessAddrs_Target_      = "AccessAddrs:Target:"
	kAccessUuids_Target_      = "AccessUuids:Target:"
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

	err := db.Client.HSet(ctx, kAccessReports_+timestamp, *report).Err()
	if err != nil {
		return err
	}

	err = db.Client.ZAdd(ctx, kAccessReports_Timestamps, redis.Z{Score: score, Member: timestamp}).Err()
	if err != nil {
		return err
	}

	return db.Client.HIncrBy(ctx, kAccessCounts_Targets, report.Target, 1).Err()
}

func (db *Database) AddAccessAddrOfTarget(ctx context.Context, target, addr string) error {
	return db.Client.SAdd(ctx, kAccessAddrs_Target_+target, addr).Err()
}

func (db *Database) AddAccessUuidOfTarget(ctx context.Context, target, uuid string) error {
	return db.Client.SAdd(ctx, kAccessUuids_Target_+target, uuid).Err()
}

func (db *Database) GetAccessReportsTimestamps(ctx context.Context, start, end string) ([]string, error) {
	return db.Client.ZRangeByScore(ctx, kAccessReports_Timestamps, &redis.ZRangeBy{Min: start, Max: end}).Result()
}

func (db *Database) GetAccessReportByTimestamp(ctx context.Context, timestamp string) (report AccessReport, _ error) {
	return report, db.Client.HGetAll(ctx, kAccessReports_+timestamp).Scan(&report)
}

func (db *Database) GetAccessCountOfTarget(ctx context.Context, target string) (string, error) {
	return db.Client.HGet(ctx, kAccessCounts_Targets, target).Result()
}

func (db *Database) GetAccessAddrCountOfTarget(ctx context.Context, target string) (int64, error) {
	return db.Client.SCard(ctx, kAccessAddrs_Target_+target).Result()
}

func (db *Database) GetAccessUuidCountOfTarget(ctx context.Context, target string) (int64, error) {
	return db.Client.SCard(ctx, kAccessUuids_Target_+target).Result()
}

func (db *Database) CountAll(ctx context.Context, userAgentFilter []string) error {
	timestamps, err := db.GetAccessReportsTimestamps(ctx, "0", time.Now().UTC().Format("20060102150405.99"))
	if err != nil {
		return err
	}

	for _, timestamp := range timestamps {
		report, err := db.GetAccessReportByTimestamp(ctx, timestamp)
		if err != nil {
			return err
		}

		isFiltered := false
		for _, substr := range userAgentFilter {
			if strings.Contains(report.UserAgent, substr) {
				isFiltered = true
				break
			}
		}
		if isFiltered {
			continue
		}

		err = db.AddAccessAddrOfTarget(ctx, report.Target, report.SourceIP)
		if err != nil {
			return err
		}

		err = db.AddAccessUuidOfTarget(ctx, report.Target, report.UUID)
		if err != nil {
			return err
		}

		err = db.Client.HIncrBy(ctx, kAccessCounts_Targets, report.Target, 1).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
