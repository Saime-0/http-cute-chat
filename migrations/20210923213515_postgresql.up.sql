CREATE TYPE unit_type AS ENUM ('USER', 'CHAT');
CREATE TYPE char_type AS ENUM ('ADMIN', 'MODER', 'NONE');
CREATE TYPE message_type AS ENUM ('SYSTEM', 'USER', 'FORMATTED');
CREATE TYPE action_type AS ENUM ('READ', 'WRITE');
CREATE TYPE group_type AS ENUM ('USERS', 'CHARS', 'ROLES');
CREATE TYPE fetch_type AS ENUM ('POSITIVE', 'NEUTRAL', 'NEGATIVE');

create table schema_migrations (
    version bigint  not null primary key,
    dirty   boolean not null
);

alter table schema_migrations owner to postgres;

create function generate_invite_code() returns text language sql as $$
  SELECT string_agg (substr('abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', ceil (random() * 62)::integer, 1), '')
  FROM generate_series(1, 16)
$$;

create function unix_utc_now(bigint = 0) returns bigint language sql as $$
    SELECT (date_part('epoch'::text, now()))::bigint + $1
$$;

CREATE TABLE units (
	id bigserial primary key,
	domain varchar(32) not null unique,
	name varchar(32) not null,
    type unit_type not null
);
CREATE TABLE users (
	id bigint primary key references units (id),
	app_settings varchar(256),
	password varchar(32) not null,
	email varchar(256) not null unique
);

CREATE TABLE chats (
	id bigint primary key references units (id) not null,
	owner_id bigint references users (id) not null,
	private boolean not null
);
CREATE TABLE chat_banlist (
	chat_id bigint references chats (id) not null,
	user_id bigint references users (id) not null
);
CREATE TABLE invites (
	code varchar(16) primary key default generate_invite_code() not null,
	chat_id bigint references chats (id) not null,
	aliens smallint,
	expires_at bigint
);
CREATE TABLE rooms (
	id bigserial primary key,
	chat_id bigint references chats (id) not null,
	parent_id bigint references rooms (id),
	name varchar(32) not null,
	note varchar(64),
	msg_format varchar(8192)
);
CREATE TABLE roles (
	id bigserial primary key,
	chat_id bigint references chats (id) not null,
	name varchar(32) not null,
	color varchar(7) not null
);
CREATE TABLE chat_members (
	id bigserial primary key,
	user_id bigint references users (id) not null,
	chat_id bigint references chats (id) not null,
	role_id bigint references roles (id),
	char char_type not null,
	joined_at bigint default unix_utc_now() not null,
	muted boolean default false not null,
	frozen boolean default false not null
);
CREATE TABLE messages (
	id bigserial primary key,
	reply_to bigint references messages (id),
	author bigint references chat_members (id),
	room_id bigint references rooms (id) not null,
	body varchar(8192) not null,
	type message_type not null,
	created_at bigint default unix_utc_now() not null
);
CREATE TABLE refresh_sessions (
    id bigserial primary key,
    user_id bigint references users (id) not null,
    refresh_token varchar(16) not null,
    user_agent varchar(1024) not null,
    expires_at bigint not null,
    created_at bigint default unix_utc_now() not null
);
CREATE TABLE votes (
    id bigserial primary key,
    room_id bigint references rooms (id) not null,
    question varchar(64) not null ,
    date bigint not null
);
CREATE TABLE vote_answers (
    id bigserial primary key,
    vote_id bigint references votes (id) not null,
    answer varchar(64) not null
);
CREATE TABLE voters (
    id bigserial primary key,
    answer_id bigint references vote_answers (id) not null,
    user_id bigint references users (id) not null
);
CREATE TABLE allows (
    room_id bigint references rooms (id) not null,
    action_type action_type not null,
    group_type group_type not null,
    value varchar(19) not null
);

