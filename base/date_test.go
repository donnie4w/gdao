// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package base

import (
	"fmt"
	"testing"
	"time"
)

func Test_date(t *testing.T) {
	testDates := []string{
		"2023-07-08",
		"2023-07-08T15:20:31",
		"2023-07-08 15:20:32",
		"2023-07-08T15:20:33.123",
		"2023 07-08T15:20:34.123456",
		"2023-07-08 15:21:35.1234",
		"2023-07-08 15:21:36.12345",
		"2023-07-08 15:21:37.123456",
		"2023-07-08T15:20:38.123+02:00",
		"2023-07-08T15:20:39.1234567Z-02:00",
		"2023-07-08 15:20:40.123456+04:00",
		"2023-07-08T15:20:41.123Z",
	}

	for _, dateStr := range testDates {
		t, err := convertToDate(dateStr)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Converted date: %s\n", t.Format(time.RFC3339))
		}
	}
}

func BenchmarkConvert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		convertToDate("2023-07-08T15:20:31")
	}
}

func BenchmarkTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Parse(time.DateTime, "2023-07-08T15:20:31")
	}
}
