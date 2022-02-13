Create docker.env file in root directory and add following values:
```dotenv
POSTGRES_CONNECTION=postgres://<user>:<password>@db:5432/chat_db?sslmode=disable
POSTGRES_PASSWORD=<password>
POSTGRES_USER=<user>
GLOBAL_PASSWORD_SALT=<random string>
MONGODB_URI=mongodb+srv://<user>:<password>@<cluster>/test?tlsInsecure=true
SECRET_SIGNING_KEY=<random string>
SMTP_HOST=smtp.yandex.ru
SMTP_EMAIL_LOGIN=example@exampe.ru
SMTP_EMAIL_PASSWD=<password>

```
