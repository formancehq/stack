drop trigger "insert_log" on "{{.Bucket}}".logs;

drop index transactions_reference;
create unique index transaction_reference on "{{.Bucket}}".transactions (ledger, reference);

alter table "{{.Bucket}}".transactions
add column inserted_at timestamp without time zone
default now();

alter table "{{.Bucket}}".transactions
alter column timestamp
set default now();

alter table "{{.Bucket}}".transactions
alter column id
type bigint;

alter table "{{.Bucket}}".moves
alter column post_commit_volumes
drop not null,
alter column post_commit_effective_volumes
drop not null;

alter table "{{.Bucket}}".logs
alter column hash
drop not null;

-- Change from jsonb to json to keep keys order and ensure consistent hashing
alter table "{{.Bucket}}".logs
alter column data
type json;

create table "{{.Bucket}}".balances (
    ledger varchar,
    account varchar,
    asset varchar,
    balance numeric,

    primary key (ledger, account, asset)
);

insert into "{{.Bucket}}".balances
select distinct on (ledger, account_address, asset)
	ledger,
	account_address as account,
	asset,
	(moves.post_commit_volumes).inputs - (moves.post_commit_volumes).outputs as balance
from (
	select *
	from moves
	order by seq desc
) moves;

drop index moves_post_commit_volumes;
drop index moves_effective_post_commit_volumes;

drop trigger "insert_account"  on "{{.Bucket}}".accounts;
drop trigger "update_account"  on "{{.Bucket}}".accounts;
drop trigger "insert_transaction"  on "{{.Bucket}}".transactions;
drop trigger "update_transaction"  on "{{.Bucket}}".transactions;

drop index moves_account_address_array;
drop index moves_account_address_array_length;
drop index transactions_sources_arrays;
drop index transactions_destinations_arrays;
drop index accounts_address_array;
drop index accounts_address_array_length;

drop index transactions_sources;
drop index transactions_destinations;

drop aggregate "{{.Bucket}}".aggregate_objects(jsonb);
drop aggregate "{{.Bucket}}".first(anyelement);

drop function "{{.Bucket}}".array_distinct(anyarray);
drop function "{{.Bucket}}".insert_posting(_transaction_seq bigint, _ledger character varying, _insertion_date timestamp without time zone, _effective_date timestamp without time zone, posting jsonb, _account_metadata jsonb);
drop function "{{.Bucket}}".upsert_account(_ledger character varying, _address character varying, _metadata jsonb, _date timestamp without time zone, _first_usage timestamp without time zone);
drop function "{{.Bucket}}".get_latest_move_for_account_and_asset(_ledger character varying, _account_address character varying, _asset character varying, _before timestamp without time zone);
drop function "{{.Bucket}}".update_transaction_metadata(_ledger character varying, _id numeric, _metadata jsonb, _date timestamp without time zone);
drop function "{{.Bucket}}".delete_account_metadata(_ledger character varying, _address character varying, _key character varying, _date timestamp without time zone);
drop function "{{.Bucket}}".delete_transaction_metadata(_ledger character varying, _id numeric, _key character varying, _date timestamp without time zone);
drop function "{{.Bucket}}".balance_from_volumes(v "{{.Bucket}}".volumes);
drop function "{{.Bucket}}".get_all_account_volumes(_ledger character varying, _account character varying, _before timestamp without time zone);
drop function "{{.Bucket}}".first_agg(anyelement, anyelement);
drop function "{{.Bucket}}".volumes_to_jsonb(v "{{.Bucket}}".volumes_with_asset);
drop function "{{.Bucket}}".get_account_aggregated_effective_volumes(_ledger character varying, _account_address character varying, _before timestamp without time zone);
drop function "{{.Bucket}}".handle_log();
drop function "{{.Bucket}}".get_account_aggregated_volumes(_ledger character varying, _account_address character varying, _before timestamp without time zone);
drop function "{{.Bucket}}".get_aggregated_volumes_for_transaction(_ledger character varying, tx numeric);
drop function "{{.Bucket}}".insert_move(_transactions_seq bigint, _ledger character varying, _insertion_date timestamp without time zone, _effective_date timestamp without time zone, _account_address character varying, _asset character varying, _amount numeric, _is_source boolean, _account_exists boolean);
drop function "{{.Bucket}}".get_all_assets(_ledger character varying);
drop function "{{.Bucket}}".insert_transaction(_ledger character varying, data jsonb, _date timestamp without time zone, _account_metadata jsonb);
drop function "{{.Bucket}}".get_all_account_effective_volumes(_ledger character varying, _account character varying, _before timestamp without time zone);
drop function "{{.Bucket}}".get_account_balance(_ledger character varying, _account character varying, _asset character varying, _before timestamp without time zone);
drop function "{{.Bucket}}".get_aggregated_effective_volumes_for_transaction(_ledger character varying, tx numeric);
drop function "{{.Bucket}}".aggregate_ledger_volumes(_ledger character varying, _before timestamp without time zone, _accounts character varying[], _assets character varying[] );
drop function "{{.Bucket}}".get_transaction(_ledger character varying, _id numeric, _before timestamp without time zone);
drop function "{{.Bucket}}".explode_address(_address character varying);
drop function "{{.Bucket}}".revert_transaction(_ledger character varying, _id numeric, _date timestamp without time zone);

create function "{{.Bucket}}".set_volumes()
    returns trigger
    security definer
    language plpgsql
as
$$
begin
    new.post_commit_volumes = coalesce((
        select (
            (post_commit_volumes).inputs + case when new.is_source then 0 else new.amount end,
            (post_commit_volumes).outputs + case when new.is_source then new.amount else 0 end
        )
        from "{{.Bucket}}".moves
        where accounts_seq = new.accounts_seq
            and asset = new.asset
            and ledger = new.ledger
        order by seq desc
        limit 1
    ), (
        case when new.is_source then 0 else new.amount end,
        case when new.is_source then new.amount else 0 end
    ));

    return new;
end;
$$;

create function "{{.Bucket}}".set_effective_volumes()
    returns trigger
    security definer
    language plpgsql
as
$$
begin
    new.post_commit_effective_volumes = coalesce((
        select (
            (post_commit_effective_volumes).inputs + case when new.is_source then 0 else new.amount end,
            (post_commit_effective_volumes).outputs + case when new.is_source then new.amount else 0 end
        )
        from "{{.Bucket}}".moves
        where accounts_seq = new.accounts_seq
            and asset = new.asset
            and ledger = new.ledger
            and (effective_date < new.effective_date or (effective_date = new.effective_date and seq < new.seq))
        order by effective_date desc, seq desc
        limit 1
    ), (
        case when new.is_source then 0 else new.amount end,
        case when new.is_source then new.amount else 0 end
    ));

    return new;
end;
$$;

create function "{{.Bucket}}".update_effective_volumes()
    returns trigger
    security definer
    language plpgsql
as
$$
begin
    update "{{.Bucket}}".moves
    set post_commit_effective_volumes =
            (
             (post_commit_effective_volumes).inputs + case when new.is_source then 0 else new.amount end,
             (post_commit_effective_volumes).outputs + case when new.is_source then new.amount else 0 end
                )
    where accounts_seq = new.accounts_seq
        and asset = new.asset
        and effective_date > new.effective_date
        and ledger = new.ledger;

    return new;
end;
$$;

create function "{{.Bucket}}".set_log_hash()
	returns trigger
	security definer
	language plpgsql
as
$$
declare
	previousHash bytea;
	marshalledAsJSON varchar;
begin
	select hash into previousHash
	from "{{.Bucket}}".logs
	where ledger = new.ledger
	order by seq desc
	limit 1;

	-- select only fields participating in the hash on the backend and format json representation the same way
	select public.json_compact(json_build_object(
		'type', new.type,
		'data', new.data,
		'date', to_json(new.date::timestamp)#>>'{}' || 'Z',
		'idempotencyKey', coalesce(new.idempotency_key, ''),
		'id', 0,
		'hash', null
	)) into marshalledAsJSON;

	new.hash = (
		select public.digest(
			case
				when previousHash is null
					then marshalledAsJSON::bytea
					else '"' || encode(previousHash::bytea, 'base64')::bytea || E'"\n' || convert_to(marshalledAsJSON, 'LATIN1')::bytea
			end || E'\n', 'sha256'::text
        )
    );

	return new;
end;
$$;