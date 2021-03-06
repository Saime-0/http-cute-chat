drop type if exists fetch_type cascade;

drop type if exists findallow cascade;

drop table if exists schema_migrations cascade;

drop table if exists chat_banlist cascade;

drop table if exists invites cascade;

drop table if exists chat_members cascade;

drop type if exists char_type cascade;

drop table if exists roles cascade;

drop table if exists messages cascade;

drop type if exists message_type cascade;

drop table if exists refresh_sessions cascade;

drop table if exists voters cascade;

drop table if exists vote_answers cascade;

drop table if exists votes cascade;

drop table if exists allows cascade;

drop type if exists action_type cascade;

drop type if exists group_type cascade;

drop table if exists rooms cascade;

drop table if exists count_members cascade;

drop table if exists chats cascade;

drop table if exists users cascade;

drop table if exists units cascade;

drop type if exists unit_type cascade;

drop table if exists registration_session cascade;

drop function if exists generate_invite_code() cascade;

drop function if exists unix_utc_now(bigint) cascade;

drop function if exists delete_allow() cascade;

drop function if exists change_count_members() cascade;

drop function if exists create_or_delete_count_members_row() cascade;

drop function if exists generate_secret(bigint) cascade;

drop function if exists generate_num_secret(bigint) cascade;

