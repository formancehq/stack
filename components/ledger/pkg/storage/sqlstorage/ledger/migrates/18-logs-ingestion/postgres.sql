--statement
create table if not exists "VAR_LEDGER_NAME".logs_ingestion (
    log_id bigint primary key
);
--statement
create table "VAR_LEDGER_NAME".transactions2 AS TABLE "VAR_LEDGER_NAME".transactions;
--statement
create table "VAR_LEDGER_NAME".postings2 AS TABLE "VAR_LEDGER_NAME".postings;
--statement
CREATE TABLE IF NOT EXISTS "VAR_LEDGER_NAME".accounts2
(
    "address"  varchar NOT NULL,
    "metadata" jsonb   DEFAULT '{}',

    UNIQUE ("address")
);
--statement
create table if not exists "VAR_LEDGER_NAME".volumes2
(
    "account" varchar,
    "asset"   varchar,
    "input"   bigint,
    "output"  bigint,

    UNIQUE ("account", "asset")
);
