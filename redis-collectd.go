// The MIT License (MIT)
//
// Copyright (c) 2014 Jamie Alquiza
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

// Array of stats we want to extract from Redis INFO
var statsFilter = []string{
	"aof_current_rewrite_time_sec$",
	"aof_enabled$",
	"aof_last_bgrewrite_status$",
	"aof_last_rewrite_time_sec$",
	"aof_last_write_status$",
	"aof_rewrite_in_progress$",
	"aof_rewrite_scheduled$",
	"blocked_clients$",
	"client_biggest_input_buf$",
	"client_longest_output_list$",
	"connected_clients$",
	"connected_slaves$",
	"evicted_keys$",
	"expired_keys$",
	"hz$",
	"keyspace_hits$",
	"keyspace_misses$",
	"latest_fork_usec$",
	"loading$",
	"lru_clock$",
	"master_repl_offset$",
	"mem_fragmentation_ratio$",
	"pubsub_channels$",
	"pubsub_patterns$",
	"rdb_bgsave_in_progress$",
	"rdb_changes_since_last_save$",
	"rdb_current_bgsave_time_sec$",
	"rdb_last_bgsave_status$",
	"rdb_last_bgsave_time_sec$",
	"rdb_last_save_time$",
	"rejected_connections$",
	"repl_backlog_active$",
	"repl_backlog_first_byte_offset$",
	"repl_backlog_histlen$",
	"repl_backlog_size$",
	"sync_full$",
	"sync_partial_err$",
	"sync_partial_ok$",
	"total_commands_processed$",
	"total_connections_received$",
	"uptime_in_seconds$",
	"used_cpu_sys$",
	"used_cpu_sys_children$",
	"used_cpu_user$",
	"used_cpu_user_children$",
	"used_memory$",
	"used_memory_lua$",
	"instantaneous_ops_per_sec$",
	"used_memory_peak$",
	"used_memory_rss$",
}

func queryRedis(host, port string) []string {
	// Dial instance
	conn, err := net.DialTimeout("tcp", host+":"+port, time.Duration(3) * time.Second)
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close()

	// Send redis INFO request
	conn.Write([]byte("INFO\n"))

	// Handle response
	respBuf := make([]byte, 2048)
	conn.Read(respBuf)
	resp := string(respBuf)

	// Parse INFO resp
	r := regexp.MustCompile("[a-z_0-9?]*:.*")
	parsed := r.FindAllString(resp, -1)
	return parsed
}

func mapStats(a []string) map[string]string {
	// Map
	stats := make(map[string]string)
	for i := range a {
		line := strings.Split(a[i], ":")
		stats[line[0]] = line[1]
	}
	return stats
}

func arrayToRegex(a []string) string {
	formatted := "\""
	for i := range a {
		if i == len(a)-1 {
			formatted += a[i] + "\""
		} else {
			formatted += a[i] + "|"
		}
	}
	return formatted
}

func filterStats(m *map[string]string, r *regexp.Regexp) {
	// Filter map to match regex
	stats := make(map[string]string)
	for k, v := range *m {
		if r.MatchString(k) {
			stats[k] = v
		}
	}
	*m = stats
}

func formatToCollectd(m map[string]string, id string) map[string]string {
	/*
		Set values to something sane for Collectd/Graphite:
		E.g. convert null values to 0, OKs to 1 (for drawAsInfinite function),
		reset unexpected values to 0, etc.
	*/
	isNan := regexp.MustCompile("ok")
	isUnknown := regexp.MustCompile("[a-zA-Z]")
	for k, v := range m {
		switch {
		case v == "":
			m[k] = "N:0"
		case isNan.MatchString(v):
			m[k] = "N:1"
		case isUnknown.MatchString(v):
			m[k] = "N:0"
		default:
			m[k] = "N:" + v
		}
	}

	// Return formatted map with hostname prefix and other decoration
	var hostname = os.Getenv("COLLECTD_HOSTNAME")
	if hostname == "" {
		hostname, _ = os.Hostname()
	}

	// Default data type to gauge, regex to change metric to counter
	counter := regexp.MustCompile("!") // Currently set to none for simplicity

	formatted := make(map[string]string)
	for k, v := range m {
		switch {
		case counter.MatchString(k):
			formatted["PUTVAL "+hostname+"/redis-"+id+"/counter-"+k] = v
		default:
			formatted["PUTVAL "+hostname+"/redis-"+id+"/gauge-"+k] = v
		}
	}
	return formatted
}

func main() {
	// Get Redis instance, defined by arg format: ./redis-watch 127.0.0.1 6379
	type redisInst struct {
		host string
		port string
	}
	redis := redisInst{}
	args := os.Args[1:]
	if len(args) != 2 {
		os.Exit(1) // Fail silent
	} else {
		redis.host = os.Args[1]
		redis.port = os.Args[2]
	}

	// Convert statsFilter list to regexp string and compile
	statsRegexString := arrayToRegex(statsFilter)
	statsRegex := regexp.MustCompile(statsRegexString)

	// Get INFO
	redisInfo := queryRedis(redis.host, redis.port)
	// Convert INFO resp to stats map
	stats := mapStats(redisInfo)
	// Filter stats that match statsFilter
	filterStats(&stats, statsRegex)
	// Format stats for Collectd
	collectdStats := formatToCollectd(stats, redis.port)

	for k, v := range collectdStats {
		fmt.Println(k, v)
	}
}
