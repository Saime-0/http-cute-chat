set -e

export folder="/usr/local/migrator"
mkdir -p $folder

export tmp="./temp/migrator"
mkdir -p $tmp

wget https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz -O $tmp/migrator.tar.gz
tar -xzf $tmp/migrator.tar.gz -C $folder