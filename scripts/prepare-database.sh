#!/bin/sh
echo "echo \$PATH:"
echo $PATH
echo "echo \$GOPATH:"
echo $GOPATH
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
   PGPASSWORD=$POSTGRES_PASSWORD psql -h "$host" -U "$POSTGRES_USER" -a -f "$dbInitScript"
fi

# 3 ===
echo -e "\n${Cyan}[STAGE 3] Apply Migrations${White}"
migrate -path "$pathToMigrations" -database "$POSTGRES_CONNECTION" up

# 4 ===
#echo -e "\n${Cyan}[STAGE 4] Go Generate${White}"
##gqlgen
#echo "skip..."

# 5 ===
#echo -e "\n${Cyan}[STAGE 5] Build project${White}"
#go build -v ./server.go

exec $cmd
