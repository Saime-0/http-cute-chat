Create .env file in root directory and add following values:
```dotenv
POSTGRES_CONNECTION="host=localhost port=5432 user=postgres password=1234 dbname=chat_db sslmode=disable"

GLOBAL_PASSWORD_SALT=<random string>

MONGODB_URI=mongodb://mongodb:27017

SECRET_SIGNING_KEY=<random string>

SMTP_HOST=
SMTP_EMAIL_LOGIN=
SMTP_EMAIL_PASSWD=

```
