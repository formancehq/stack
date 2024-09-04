
-- create a sequence for transactions by ledger instead of a sequence of the table as we want to have contiguous ids
-- notes: we can still have "holes" on id since a sql transaction can be reverted after a usage of the sequence
create sequence "{{.Bucket}}"."{{.Name}}_transaction_id" owned by "{{.Bucket}}".transactions.id;
select setval('"{{.Bucket}}"."{{.Name}}_transaction_id"', coalesce((
    select max(id) + 1
    from "{{.Bucket}}".transactions
    where ledger = '{{ .Name }}'
), 1)::bigint, false);

-- create a sequence for logs by ledger instead of a sequence of the table as we want to have contiguous ids
-- notes: we can still have "holes" on id since a sql transaction can be reverted after a usage of the sequence
create sequence "{{.Bucket}}"."{{.Name}}_log_id" owned by "{{.Bucket}}".logs.id;
select setval('"{{.Bucket}}"."{{.Name}}_log_id"', coalesce((
    select max(id) + 1
    from "{{.Bucket}}".logs
    where ledger = '{{ .Name }}'
), 1)::bigint, false);


-- enable post commit volumes synchronously
{{ if .HasFeature "POST_COMMIT_VOLUMES" "SYNC" }}
create index "{{.Name}}_pcv" on "{{.Bucket}}".moves (accounts_seq, asset, seq) where ledger = '{{.Name}}';

create trigger "{{.Name}}_set_volumes"
before insert
on "{{.Bucket}}"."moves"
for each row
when (
    new.ledger = '{{.Name}}'
)
execute procedure "{{.Bucket}}".set_volumes();
{{ end }}

-- enable post commit effective volumes synchronously

{{ if .HasFeature "POST_COMMIT_EFFECTIVE_VOLUMES" "SYNC" }}
create index "{{.Name}}_pcev" on "{{.Bucket}}".moves (accounts_seq, asset, effective_date desc) where ledger = '{{.Name}}';

create trigger "{{.Name}}_set_effective_volumes"
before insert
on "{{.Bucket}}"."moves"
for each row
when (
    new.ledger = '{{.Name}}'
)
execute procedure "{{.Bucket}}".set_effective_volumes();

create trigger "{{.Name}}_update_effective_volumes"
after insert
on "{{.Bucket}}"."moves"
for each row
when (
    new.ledger = '{{.Name}}'
)
execute procedure "{{.Bucket}}".update_effective_volumes();
{{ end }}

-- logs hash

{{ if .HasFeature "HASH_LOGS" "SYNC" }}
create trigger "{{.Name}}_set_log_hash"
before insert
on "{{.Bucket}}"."logs"
for each row
when (
    new.ledger = '{{.Name}}'
)
execute procedure "{{.Bucket}}".set_log_hash();
{{ end }}

{{ if .HasFeature "ACCOUNT_METADATA_HISTORIES" "SYNC" }}
create trigger "{{.Name}}_update_account_metadata_history"
after update
on "{{.Bucket}}"."accounts"
for each row
when (
    new.ledger = '{{.Name}}'
)
execute procedure "{{.Bucket}}".update_account_metadata_history();

create trigger "{{.Name}}_insert_account_metadata_history"
after insert
on "{{.Bucket}}"."accounts"
for each row
when (
    new.ledger = '{{.Name}}'
)
execute procedure "{{.Bucket}}".insert_account_metadata_history();
{{ end }}

{{ if .HasFeature "TRANSACTION_METADATA_HISTORIES" "SYNC" }}
create trigger "{{.Name}}_update_transaction_metadata_history"
after update
on "{{.Bucket}}"."transactions"
for each row
when (
    new.ledger = '{{.Name}}'
)
execute procedure "{{.Bucket}}".update_transaction_metadata_history();

create trigger "{{.Name}}_insert_transaction_metadata_history"
after insert
on "{{.Bucket}}"."transactions"
for each row
when (
    new.ledger = '{{.Name}}'
)
execute procedure "{{.Bucket}}".insert_transaction_metadata_history();
{{ end }}

{{ if .HasFeature "INDEX_ADDRESS_SEGMENTS" "ON" }}
--todo: names became too long, are stripped
create index "{{.Name}}_moves_account_address_array" on "{{.Bucket}}".moves using gin (account_address_array jsonb_ops) where ledger = '{{.Name}}';
create index "{{.Name}}_moves_account_address_array_length" on "{{.Bucket}}".moves (jsonb_array_length(account_address_array)) where ledger = '{{.Name}}';

create index "{{.Name}}_accounts_address_array" on "{{.Bucket}}".accounts using gin (address_array jsonb_ops) where ledger = '{{.Name}}';
create index "{{.Name}}_accounts_address_array_length" on "{{.Bucket}}".accounts (jsonb_array_length(address_array)) where ledger = '{{.Name}}';

{{ if .HasFeature "INDEX_TRANSACTION_ACCOUNTS" "ON" }}
create index "{{.Name}}_transactions_sources_arrays" on "{{.Bucket}}".transactions using gin (sources_arrays jsonb_path_ops) where ledger = '{{.Name}}';
create index "{{.Name}}_transactions_destinations_arrays" on "{{.Bucket}}".transactions using gin (destinations_arrays jsonb_path_ops) where ledger = '{{.Name}}';
{{ end }}
{{ end }}

{{ if .HasFeature "INDEX_TRANSACTION_ACCOUNTS" "ON" }}
create index "{{.Name}}_transactions_sources" on "{{.Bucket}}".transactions using gin (sources jsonb_path_ops) where ledger = '{{.Name}}';
create index "{{.Name}}_transactions_destinations" on "{{.Bucket}}".transactions using gin (destinations jsonb_path_ops) where ledger = '{{.Name}}';
{{ end }}