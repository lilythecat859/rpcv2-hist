#!/usr/bin/env bash
# AGPL-3.0
# Apply production-grade ClickHouse settings for AMD 16-core / 256 GB box

set -euo pipefail

CLICKHOUSE_USER="${CLICKHOUSE_USER:-default}"
CLICKHOUSE_PASSWORD="${CLICKHOUSE_PASSWORD:-}"

echo "=> Applying ClickHouse tuning for rpcv2-hist"

clickhouse-client --query "
ALTER USER '$CLICKHOUSE_USER' SETTINGS
  max_memory_usage = 200000000000, -- 200 GB
  max_memory_usage_for_user = 200000000000,
  max_threads = 32,
  max_parallel_replicas = 2,
  max_distributed_connections = 500,
  max_query_size = 1073741824,
  max_ast_elements = 1000000,
  max_bytes_before_external_group_by = 100000000000,
  max_bytes_before_external_sort = 100000000000;
"

clickhouse-client --query "
SYSTEM DROP QUERY CACHE;
SYSTEM DROP COMPILE EXPRESSION CACHE;
SYSTEM DROP QUERY RESULT CACHE;
"

echo "=> Done"