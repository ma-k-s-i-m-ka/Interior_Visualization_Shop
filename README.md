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
The server is running at http://localhost:3001. If desired, you can change the connection settings in the [config.yml] file. To launch the application and further connect the database and the feedback system, you need to create the [.env] file and fill it as follows:
```sh
ENV=dev
USE_HTTP=true
DATABASE_DSN= your postgres patch
MAIL_ADD= your mail
MAIL_PAS= your mail pass
```

After the [.env] file has been created, and you have entered the project directory, you must run the command
```sh
cd ...your path...\Interior_Visualization_Shop
```
```sh
.\main.exe  
```
If your environment is not Windows use docker
```sh
 docker build -t main-docker .
```
```sh
 docker run -it main-docker   
```
But if you run it through docker, before building, remove the line
```sh
	/// открытие веб-страницы в браузере \\\
	err := browser.OpenURL("http://" + s.srv.Addr + "/")
	if err != nil {
		return err
	}
```
in file
```sh
app/internal/server/server.go
```

The site will run at the specified address in the config.yml file. The origin is http://localhost:3001
## Testing
To quickly test the application, I used POSTMAN. Folder with requests [Postman](https://drive.google.com/drive/folders/1Z1P4vn_768SQYPEMikMF6PYoomEck-UB?usp=sharing)
