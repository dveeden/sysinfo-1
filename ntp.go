// Copyright © 2018 PingCAP Inc.
//
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.
//
// Use ntpq to get basic info of NTPd on the system

package sysinfo

import (
	"bytes"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type TimeStat struct {
	Ver     string  `json:"version,omitempty"`
	Sync    string  `json:"sync,omitempty"`
	Stratum int     `json:"stratum,omitempty"`
	Offset  float64 `json:"offset,omitempty"`
	Jitter  float64 `json:"jitter,omitempty"`
	Status  string  `json:"status,omitempty"`
}

func (si *SysInfo) getNTPInfo() {
	syncd, err := exec.LookPath("ntpq")
	if err != nil {
		si.NTP.Ver = err.Error()
		return
	}

	cmd := exec.Command(syncd, "-c rv")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// set default sync status to false
	si.NTP.Sync = "none"

	output := strings.FieldsFunc(out.String(), multi_split)
	for _, kv := range output {
		_tmp := strings.Split(strings.TrimSpace(kv), "=")
		switch {
		case _tmp[0] == "version":
			si.NTP.Ver = strings.Trim(_tmp[1], "\"")
		case _tmp[0] == "stratum":
			si.NTP.Stratum, _ = strconv.Atoi(_tmp[1])
		case _tmp[0] == "offset":
			si.NTP.Offset, _ = strconv.ParseFloat(_tmp[1], 64)
		case _tmp[0] == "sys_jitter":
			si.NTP.Jitter, _ = strconv.ParseFloat(_tmp[1], 64)
		case strings.Contains(_tmp[0], "sync"):
			si.NTP.Sync = _tmp[0]
		case len(_tmp) > 2 && strings.Contains(_tmp[1], "status"):
			// sample line of _tmp: ["associd", "0 status", "0618 leap_none"]
			si.NTP.Status = strings.Split(_tmp[2], " ")[0]
		default:
			continue
		}
	}
}

func multi_split(r rune) bool {
	switch r {
	case ',':
		return true
	case '\n':
		return true
	default:
		return false
	}
}
