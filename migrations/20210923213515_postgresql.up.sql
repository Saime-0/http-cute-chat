create type unit_type as enum ('USER', 'CHAT');

create type message_type as enum ('SYSTEM', 'USER', 'FORMATTED');

create type action_type as enum ('READ', 'WRITE');

create type group_type as enum ('MEMBER', 'CHAR', 'ROLE');

create type fetch_type as enum ('POSITIVE', 'NEUTRAL', 'NEGATIVE');

create type char_type as enum ('ADMIN', 'MODER');

create type findallow as
(
    act varchar,
    gr varchar,
    val varchar
);

create table if not exists schema_migrations
(
    version bigint not null
        primary key,
    dirty boolean not null
);

create table if not exists units
(
    id bigserial
        primary key,
    domain varchar(32) not null
        unique,
    name varchar(32) not null,
    type unit_type not null
);

create table if not exists users
(
    id bigint not null
        primary key
        references units
            on delete cascade,
    hashed_password varchar(128) not null,
    email varchar(256) not null
        unique
);

create table if not exists chats
(
    id bigint not null
        primary key
        references units
            on delete cascade,
    owner_id bigint not null
        references users
            on delete cascade,
    private boolean not null
);

create table if not exists chat_banlist
(
    chat_id bigint not null
        references chats
            on delete cascade,
    user_id bigint not null
        references users
);

create table if not exists invites
(
    code varchar(16) default generate_secret((16)::bigint) not null,
    chat_id bigint not null
        constraint invite_links_chat_id_fkey
            references chats
            on delete cascade,
    aliens smallint,
    expires_at bigint,
    id bigserial
        constraint invites_pk
            primary key
);

create unique index if not exists invites_code_uindex
    on invites (code);

create table if not exists rooms
(
    id bigserial
        primary key,
    chat_id bigint not null
        references chats
            on delete cascade,
    parent_id bigint
                   references rooms
                       on delete set null,
    name varchar(32) not null,
    note varchar(64) default NULL::character varying,
    msg_format varchar(8192) default NULL::character varying
);

create table if not exists roles
(
    id bigserial
        primary key,
    chat_id bigint not null
        references chats
            on delete cascade,
    name varchar(32) not null,
    color varchar(7) not null
);

create table if not exists chat_members
(
    id bigserial
        primary key,
    user_id bigint not null
        references users,
    chat_id bigint not null
        references chats
            on delete cascade,
    role_id bigint
                   references roles
                       on delete set null,
    char char_type,
    joined_at bigint default unix_utc_now() not null,
    muted boolean default false not null
);

create table if not exists messages
(
    id bigserial
        primary key,
    reply_to bigint
                                             references messages
                                                 on delete set null,
    user_id bigint
        constraint messages_users_id_fk
            references users,
    room_id bigint not null
        references rooms
            on delete cascade,
    body varchar(8192) not null,
    type message_type not null,
    created_at bigint default unix_utc_now() not null
);

create table if not exists refresh_sessions
(
    id bigserial
        primary key,
    user_id bigint not null
        references users,
    refresh_token varchar(32) not null,
    user_agent varchar(1024) not null,
    expires_at bigint not null,
    created_at bigint default unix_utc_now() not null
);

create table if not exists votes
(
    id bigserial
        primary key,
    room_id bigint not null
        references rooms
            on delete cascade,
    question varchar(64) not null,
    date bigint not null
);

create table if not exists vote_answers
(
    id bigserial
        primary key,
    vote_id bigint not null
        references votes
            on delete cascade,
    answer varchar(64) not null
);

create table if not exists voters
(
    id bigserial
        primary key,
    answer_id bigint not null
        references vote_answers
            on delete cascade,
    user_id bigint not null
        references users
);

create table if not exists allows
(
    room_id bigint not null
        references rooms
            on delete cascade,
    action_type action_type not null,
    group_type group_type not null,
    value varchar(19) not null,
    id bigserial
        primary key
);

create table if not exists count_members
(
    id bigserial
        primary key,
    chat_id bigint not null
        references chats,
    count_value integer default 0 not null
);

create table if not exists registration_session
(
    id bigserial,
    domain varchar(32) not null
        unique,
    name varchar(32) not null,
    email varchar(256) not null
        unique,
    hashed_password varchar(128) not null,
    verify_code varchar(6) default generate_num_secret((6)::bigint) not null,
    expires_at bigint
);

create or replace function generate_invite_code() returns text
    language sql
as $$
SELECT string_agg (substr('abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', ceil (random() * 62)::integer, 1), '')
FROM generate_series(1, 16)
$$;

create or replace function unix_utc_now(bigint DEFAULT 0) returns bigint
    language sql
as $$
SELECT (date_part('epoch'::text, now()))::bigint + $1
$$;

create or replace function delete_allow() returns trigger
    language plpgsql
as $$
BEGIN
    IF tg_table_name = 'chat_members' THEN
        DELETE FROM allows
        WHERE group_type = 'MEMBER' AND value = old.id::VARCHAR;
    ELSE
        DELETE FROM allows
        WHERE group_type = 'ROLE' AND value = OLD.ID::VARCHAR;
    end if;
    raise notice '%', old.id;
    RETURN NULL;
end;
$$;

create trigger on_delete_role
    after delete
    on roles
    for each row
execute procedure delete_allow();

create trigger on_delete_member
    after delete
    on chat_members
    for each row
execute procedure delete_allow();

create or replace function change_count_members() returns trigger
    language plpgsql
as $$
BEGIN
    if tg_op = 'INSERT' then
        update count_members
        set count_value = count_value + 1
        where count_members.chat_id = new.chat_id;
        return new;
    else if tg_op = 'DELETE' then
        update count_members
        set count_value = count_value - 1
        where count_members.chat_id = new.chat_id;
        return new;
    end if;
    end if;
    raise exception 'operation could not be detected';
end;
$$;

create trigger on_change_members_table
    after insert or delete
    on chat_members
    for each row
execute procedure change_count_members();

create or replace function create_or_delete_count_members_row() returns trigger
    language plpgsql
as $$
begin
    if tg_op = 'INSERT' then
        insert into count_members (chat_id) values (new.id);
        return new;
    else if tg_op = 'DELETE' then
        delete from count_members WHERE chat_id = old.id;
        return old;
    end if;
    end if;
    raise exception 'operation could not be detected';
end;
$$;

create trigger on_create_chat
    after insert or delete
    on chats
    for each row
execute procedure create_or_delete_count_members_row();

create or replace function generate_secret(len bigint) returns text
    language sql
as $$
SELECT string_agg (substr('abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', ceil (random() * 62)::integer, 1), '')
FROM generate_series(1, len)
$$;

create or replace function generate_num_secret(len bigint) returns text
    language sql
as $$
SELECT string_agg (substr('0123456789', ceil (random() * 10)::integer, 1), '')
FROM generate_series(1, len)
$$;

