-- Download extensions
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE units (
	id bigserial primary key,
	domain varchar(32) not null unique,
	name varchar(32) not null
);
CREATE TABLE users (
	id bigint primary key references units (id),
	app_settings varchar(512)
);
CREATE TABLE chats (
	id bigint primary key references units (id) not null,
	owner_id bigint references users (id) not null
);
CREATE TABLE chat_members (
	id bigserial primary key,
	user_id bigint references users (id) not null,
	chat_id bigint references chats (id) not null
);
CREATE TABLE rooms (
	id bigserial primary key,
	chat_id bigint references chats (id) not null,
	parent_id bigint references rooms (id),
	name varchar(32) not null,
	note varchar(64)
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
	type smallint not null
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
    created_at bigint not null,
	counter smallint not null default 0
);