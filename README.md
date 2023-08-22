# Interior_Visualization_Shop
REST API application on GO with a test front-end view. It is a store with interior visualization services. using SOLID and clean architecture.

Сontains the following functionality:
- CRUD systems for working with basic objects (users, appeal)
- Email feedback system
- JWT-based authentication
- Error Handling
- Logging of the system operation
- The ability to change the configuration

The application uses the following auxiliary and replaceable packages at your discretion:
- Routing: [httprouter](https://github.com/julienschmidt/httprouter)
- Database access: [pgx/v4](https://github.com/jackc/pgx)
- Logging: [logrus](https://github.com/sirupsen/logrus)
- JWT: [jwt-go](https://github.com/dgrijalva/jwt-go)

## Getting Started
The server works at http://localhost:3001. Optionally, you can change the connection settings of both the server and the database in the [config.yml] file. For the feedback system, you need to create an [.env] file and fill it in as follows:
```sh
ENV=dev
USE_HTTP=true
DATABASE_DSN= your postgres patch
MAIL_ADD= your mail
MAIL_PAS= your mail pass
```
## Testing
Tested the application using POSTMAN. Folder with requests [Postman](https://drive.google.com/drive/folders/1Z1P4vn_768SQYPEMikMF6PYoomEck-UB?usp=sharing)
