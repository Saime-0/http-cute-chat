CREATE TABLE units (
	id bigserial primary key,
	domain varchar(32) not null unique,
	name varchar(32) not null
);
CREATE TABLE users (
	id bigint primary key references units (id),
	app_settings varchar(256),
	password varchar(32) not null,
	email varchar(256) not null
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
create function generate_invite_code() returns text language sql as $$
  SELECT string_agg (substr('abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', ceil (random() * 62)::integer, 1), '')
  FROM generate_series(1, 16)
$$;
CREATE TABLE invite_links (
	code varchar(16) primary key default generate_invite_code(),
	chat_id bigint references chats (id) not null,
	aliens smallint,
	exp bigint not null
);
CREATE TABLE rooms (
	id bigserial primary key,
	chat_id bigint references chats (id) not null,
	parent_id bigint references rooms (id),
	name varchar(32) not null,
	note varchar(64),
	msg_format varchar(8192),
	private boolean not null
);
CREATE TABLE roles (
	id bigserial primary key,
	chat_id bigint references chats (id) not null,
	role_name varchar(32) not null,
	color varchar(7) not null,
	visible boolean not null,
	manage_rooms boolean not null,
	room_id bigint references rooms (id),
	manage_chat boolean not null,
	manage_roles boolean not null,
	manage_members boolean not null
);
CREATE TABLE chat_members (
	id bigserial primary key,
	user_id bigint references users (id) not null,
	chat_id bigint references chats (id) not null,
	role_id bigint references roles (id) not null,
	joined_at bigint not null
);
CREATE TABLE room_whitelist (
	room_id bigint references rooms (id) not null,
	role_id bigint references roles (id) not null
);
CREATE TABLE dialogs (
	id bigserial primary key,
	user1 bigint references users (id) not null,
	user2 bigint references users (id) not null
);
CREATE TABLE messages (
	id bigserial primary key,
	reply_to bigint references messages (id),
	author bigint references units (id) not null,
	body varchar(8192) not null,
	type smallint not null,
	created_at bigint not null 
);
CREATE TABLE dialog_msg_pool (
	dialog_id bigint references dialogs (id) not null,
	message_id bigint references messages (id) not null
);
CREATE TABLE room_msg_pool (
	room_id bigint references rooms (id) not null,
	message_id bigint references messages (id) not null
);
CREATE TABLE refresh_sessions (
    id bigserial primary key,
    user_id bigint references users (id) not null,
    refresh_token varchar(16) not null,
    user_agent varchar(64) not null,
    exp bigint not null,
    created_at bigint not null
);
