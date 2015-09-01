redis-collectd
==============

Redis plugin for collectd.

Based on https://github.com/jamiealquiza/redis-collectd, but adding the host so we can collect stats per host and port.

### How it works

Since this is used as a collectd plugin, it pulls metrics on a best-effort basis and will fail silent with incorrect command line arguments, failed connections to Redis, and so forth.

The plugin honors the `COLLECTD_HOSTNAME` variable passed by collectd and will failover to detecting the system hostname if no value is passed. The metric namespace uses the `port` argument passed in, allowing several Redis instances per-box to be monitored.

### Example setup

Build / place binary in your plugin path. Example:
<pre>
go build redis-collectd.go
cp redis-collectd /opt/collectd/lib/collectd/
chown root:root /opt/collectd/lib/collectd/redis-collectd
</pre>

In collectd.conf:
<pre>
LoadPlugin exec
<Plugin exec>
  Exec "nobody" "/opt/collectd/lib/collectd/redis-collectd" "127.0.01.1" "6379"
</Plugin>
</pre>

### Example output

<pre>
./redis-collectd 127.0.0.1 6379
PUTVAL some-server/redis-6379/gauge-aof_last_bgrewrite_status N:1
PUTVAL some-server/redis-6379/gauge-used_cpu_sys N:96.59
PUTVAL some-server/redis-6379/gauge-instantaneous_ops_per_sec N:0
PUTVAL some-server/redis-6379/gauge-expired_keys N:0
PUTVAL some-server/redis-6379/gauge-aof_enabled N:0
PUTVAL some-server/redis-6379/gauge-total_connections_received N:287
PUTVAL some-server/redis-6379/gauge-repl_backlog_first_byte_offset N:0
PUTVAL some-server/redis-6379/gauge-sync_full N:0
PUTVAL some-server/redis-6379/gauge-used_cpu_user_children N:0.00
PUTVAL some-server/redis-6379/gauge-keyspace_misses N:0
PUTVAL some-server/redis-6379/gauge-lru_clock N:15260340
PUTVAL some-server/redis-6379/gauge-used_memory_lua N:33792
PUTVAL some-server/redis-6379/gauge-repl_backlog_active N:0
PUTVAL some-server/redis-6379/gauge-rejected_connections N:0
PUTVAL some-server/redis-6379/gauge-aof_last_write_status N:1
PUTVAL some-server/redis-6379/gauge-pubsub_channels N:0
PUTVAL some-server/redis-6379/gauge-used_cpu_sys_children N:0.00
PUTVAL some-server/redis-6379/gauge-rdb_changes_since_last_save N:0
PUTVAL some-server/redis-6379/gauge-sync_partial_err N:0
PUTVAL some-server/redis-6379/gauge-rdb_last_bgsave_status N:1
PUTVAL some-server/redis-6379/gauge-rdb_last_save_time N:1407526200
PUTVAL some-server/redis-6379/gauge-total_commands_processed N:277
PUTVAL some-server/redis-6379/gauge-rdb_current_bgsave_time_sec N:-1
PUTVAL some-server/redis-6379/gauge-used_memory_peak N:810776
PUTVAL some-server/redis-6379/gauge-client_biggest_input_buf N:0
PUTVAL some-server/redis-6379/gauge-connected_slaves N:0
PUTVAL some-server/redis-6379/gauge-repl_backlog_histlen N:0
PUTVAL some-server/redis-6379/gauge-aof_last_rewrite_time_sec N:-1
PUTVAL some-server/redis-6379/gauge-mem_fragmentation_ratio N:2.93
PUTVAL some-server/redis-6379/gauge-rdb_last_bgsave_time_sec N:-1
PUTVAL some-server/redis-6379/gauge-blocked_clients N:0
PUTVAL some-server/redis-6379/gauge-latest_fork_usec N:0
PUTVAL some-server/redis-6379/gauge-master_repl_offset N:0
PUTVAL some-server/redis-6379/gauge-aof_rewrite_in_progress N:0
PUTVAL some-server/redis-6379/gauge-connected_clients N:1
PUTVAL some-server/redis-6379/gauge-aof_rewrite_scheduled N:0
PUTVAL some-server/redis-6379/gauge-keyspace_hits N:0
PUTVAL some-server/redis-6379/gauge-used_memory N:809872
PUTVAL some-server/redis-6379/gauge-uptime_in_seconds N:243068
PUTVAL some-server/redis-6379/gauge-sync_partial_ok N:0
PUTVAL some-server/redis-6379/gauge-client_longest_output_list N:0
PUTVAL some-server/redis-6379/gauge-pubsub_patterns N:0
PUTVAL some-server/redis-6379/gauge-hz N:10
PUTVAL some-server/redis-6379/gauge-rdb_bgsave_in_progress N:0
PUTVAL some-server/redis-6379/gauge-repl_backlog_size N:1048576
PUTVAL some-server/redis-6379/gauge-evicted_keys N:0
PUTVAL some-server/redis-6379/gauge-used_cpu_user N:43.47
PUTVAL some-server/redis-6379/gauge-loading N:0
</pre>
