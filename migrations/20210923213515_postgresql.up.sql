create type unit_type as enum ('USER', 'CHAT');

create type message_type as enum ('SYSTEM', 'USER', 'FORMATTED');

create type action_type as enum ('READ', 'WRITE');

create type group_type as enum ('MEMBER', 'CHAR', 'ROLE');

create type fetch_type as enum ('POSITIVE', 'NEUTRAL', 'NEGATIVE');

create type char_type as enum ('ADMIN', 'MODER');

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
        references units,
    app_settings varchar(256),
    password varchar(32) not null,
    email varchar(256) not null
        unique
);

create table if not exists chats
(
    id bigint not null
        primary key
        references units,
    owner_id bigint not null
        references users,
    private boolean not null
);

create table if not exists chat_banlist
(
    chat_id bigint not null
        references chats,
    user_id bigint not null
        references users
);

create table if not exists invites
(
    code varchar(16) default generate_invite_code() not null
        constraint invite_links_pkey
            primary key,
    chat_id bigint not null
        constraint invite_links_chat_id_fkey
            references chats,
    aliens smallint,
    expires_at bigint
);

create table if not exists rooms
(
    id bigserial
        primary key,
    chat_id bigint not null
        references chats,
    parent_id bigint
        references rooms,
    name varchar(32) not null,
    note varchar(64) default NULL::character varying,
    msg_format varchar(8192) default NULL::character varying
);

create table if not exists roles
(
    id bigserial
        primary key,
    chat_id bigint not null
        references chats,
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
        references chats,
    role_id bigint
        references roles,
    char char_type,
    joined_at bigint default unix_utc_now() not null,
    muted boolean default false not null,
    frozen boolean default false not null
);

create table if not exists messages
(
    id bigserial
        primary key,
    reply_to bigint
        references messages,
    author bigint
        references chat_members,
    room_id bigint not null
        references rooms,
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
        references rooms,
    question varchar(64) not null,
    date bigint not null
);

create table if not exists vote_answers
(
    id bigserial
        primary key,
    vote_id bigint not null
        references votes,
    answer varchar(64) not null
);

create table if not exists voters
(
    id bigserial
        primary key,
    answer_id bigint not null
        references vote_answers,
    user_id bigint not null
        references users
);

create table if not exists allows
(
    room_id bigint not null
        references rooms,
    action_type action_type not null,
    group_type group_type not null,
    value varchar(19) not null,
    id bigserial
        primary key
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

