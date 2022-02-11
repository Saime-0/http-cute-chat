Create .env file in root directory and add following values:
```dotenv
POSTGRES_CONNECTION="host=localhost port=5432 user=postgres password=1234 dbname=chat_db sslmode=disable"

GLOBAL_PASSWORD_SALT=<random string>

MONGODB_URI=mongodb://mongodb:27017

SECRET_SIGNING_KEY=<random string>

SMTP_HOST=smtp.yandex.ru
SMTP_EMAIL_LOGIN=chatix@cute-chat.dot
SMTP_EMAIL_PASSWD=<password>

```
