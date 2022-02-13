#!/bin/sh

set -e

pathToMigrations="$1"
shift
dbInitScript="$1"
shift
host="$1"
shift
cmd="$@"

# Regular Colors
Black='\033[0;30m'        # Black
Red='\033[0;31m'          # Red
Green='\033[0;32m'        # Green
Yellow='\033[0;33m'       # Yellow
Blue='\033[0;34m'         # Blue
Purple='\033[0;35m'       # Purple
Cyan='\033[0;36m'         # Cyan
White='\033[0;37m'        # White

# 1 ===
echo -e "\n${Cyan}[STAGE 1] Wait For Postgres${White}"

until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$host" -U "$POSTGRES_USER" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 2
done

>&2 echo "Postgres is up"

# 2 ===
echo -e "\n${Cyan}[STAGE 2] Run InitScript If DB Not Exists${White}"

if [ "$(PGPASSWORD=$POSTGRES_PASSWORD psql -h "$host" -U "$POSTGRES_USER" -tAc "SELECT 1 FROM pg_database WHERE datname='chat_db'" )" = '1' ]
then
    echo "Database 'chat_db' already exists"
else
    echo "Database 'chat_db' does not exist"
    echo "Run the init script"
#    PGPASSWORD=$POSTGRES_PASSWORD
#  PGPASSWORD=$POSTGRES_PASSWORD psql -U "$POSTGRES_USER -c "CREATE DATABASE chat_db"
   PGPASSWORD=$POSTGRES_PASSWORD psql -h "$host" -U "$POSTGRES_USER" -a -f "$dbInitScript"
fi

# 3 ===
echo -e "\n${Cyan}[STAGE 3] Apply Migrations${White}"

mkdir -p ./temp/migrator \
&& wget https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz -O ./temp/migrator/migrator.tar.gz \
&& tar -xzf ./temp/migrator/migrator.tar.gz -C ./temp/migrator \
&& ./temp/migrator/migrate -path "$pathToMigrations" -database "$POSTGRES_CONNECTION" up

exec $cmd
