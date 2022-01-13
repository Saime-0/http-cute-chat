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
    password varchar(32) not null,
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
        references users,
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
    code varchar(16) default generate_invite_code() not null
        constraint invite_links_pkey
            primary key,
    chat_id bigint not null
        constraint invite_links_chat_id_fkey
            references chats
            on delete cascade,
    aliens smallint,
    expires_at bigint
);

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

create trigger on_delete_role
    after delete
    on roles
    for each row
execute procedure delete_allow();

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
    muted boolean default false not null,
    frozen boolean default false not null
);

create trigger on_delete_member
    after delete
    on chat_members
    for each row
execute procedure delete_allow();

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

create or replace function validate_allows(chatid bigint, arr findallow[]) returns boolean
    language plpgsql
as $$
DECLARE
    act varchar;
    gr varchar;
    val varchar;
    exist bool = false;
BEGIN
    Foreach act, gr, val IN ARRAY arr
        loop
        if gr = 'MEMBER' then
            select exists(
                           select 1
                           from chats c
                                    join chat_members m
                                         on  c.id = m.chat_id
                           where c.id = chatID
                             AND m.id = val::BIGINT
                       ) into exist;
        else if gr = 'ROLE' then
            select exists(
                           select 1
                           from chats c
                                    join roles r
                                         on c.id = r.chat_id
                           where c.id = chatID
                             AND r.id = val::BIGINT
                       ) into exist;
        else if gr = 'CHAR' then
            exist = true;
        end if;
        end if;
        end if;
        if exist is false THEN
            return false;
        end if;
        end loop;
    return true;
END;
$$;

create or replace function allows_exists(expect boolean, roomid bigint, arr findallow[]) returns boolean
    language plpgsql
as $$
DECLARE
    act varchar;
    gr varchar;
    val varchar;
BEGIN
    Foreach act, gr, val IN ARRAY arr
        loop
        SELECT exists(
                       select 1
                       from allows a
                       where a.room_id = roomID
                         AND a.value = val
                         AND a.group_type = gr::group_type
                         AND a.action_type = act::action_type
                   ) into FOUND;

        if FOUND <> expect THEN
            return false;
        end if;

        end loop;

    return true;

END;
$$;