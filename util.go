// Copyright © 2016 Zlatko Čalušić
//
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.

package sysinfo

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Read one-liner text files, strip newline.
func slurpFile(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}

// Write one-liner text files, add newline, ignore errors (best effort).
func spewFile(path string, data string, perm os.FileMode) {
	_ = ioutil.WriteFile(path, []byte(data+"\n"), perm)
}

func SlurpFile(path string) string {
	return slurpFile(path)
}

func parseMemSize(memInfo string) uint64 {
	for _, line := range strings.Split(memInfo, "\n") {
		if !strings.Contains(line, "MemTotal") {
			continue
		}
		fields := strings.Fields(line)
		size, _ := strconv.ParseUint(fields[1], 10, 64)
		return size
	}
	return 0
}
