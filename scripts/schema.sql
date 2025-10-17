-- AGPL-3.0
-- ClickHouse schema for rpcv2-hist
-- Keeps hot partitions on NVMe, cold on S3 via TTL+storage policy

CREATE DATABASE IF NOT EXISTS solana;

-- blocks table (partitioned by slot)
CREATE TABLE IF NOT EXISTS solana.blocks
(
    slot        UInt64,
    blockhash   FixedString(44),
    parent_slot UInt64,
    block_time  Int64,
    height      UInt64,
    commitment  Enum8('processed' = 1, 'confirmed' = 2, 'finalized' = 3),
    raw         String CODEC(LZ4)
) ENGINE = MergeTree
PARTITION BY intDiv(slot, 864000)     -- ~100k slots per partition
ORDER BY (commitment, slot)
TTL toDateTime(block_time) + INTERVAL 30 DAY TO VOLUME 'cold'
SETTINGS index_granularity = 8192;

-- transactions table
CREATE TABLE IF NOT EXISTS solana.transactions
(
    signature     String,
    slot          UInt64,
    tx_idx        UInt64,
    block_time    Int64,
    signer        String,
    fee           UInt64,
    compute_units UInt64,
    err           Nullable(String),
    commitment    Enum8('processed' = 1, 'confirmed' = 2, 'finalized' = 3),
    raw           String CODEC(LZ4)
) ENGINE = MergeTree
PARTITION BY intDiv(slot, 864000)
ORDER BY (commitment, signature)
TTL toDateTime(block_time) + INTERVAL 30 DAY TO VOLUME 'cold'
SETTINGS index_granularity = 8192;

-- signatures index for getSignaturesForAddress
CREATE TABLE IF NOT EXISTS solana.signatures
(
address     String,
    signature   String,
    slot        UInt64,
    block_time  Int64,
    err         Nullable(String),
    memo        Nullable(String),
    commitment  Enum8('processed' = 1, 'confirmed' = 2, 'finalized' = 3)
) ENGINE = MergeTree
PARTITION BY intDiv(slot, 864000)
ORDER BY (commitment, address, slot, signature)
TTL toDateTime(block_time) + INTERVAL 30 DAY TO VOLUME 'cold'
SETTINGS index_granularity = 8192;

-- materialized view to keep latest confirmed signature per address for fast pagination
CREATE MATERIALIZED VIEW IF NOT EXISTS solana.signatures_latest
ENGINE = ReplacingMergeTree
ORDER BY (address, signature)
AS SELECT
    address,
    signature,
    slot,
    block_time
FROM solana.signatures
WHERE commitment = 'confirmed';