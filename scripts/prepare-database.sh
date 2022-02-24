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


exec $cmd
