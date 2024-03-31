# gofiber-http-server

HTTP Server built with gofiber framework

GoFiber is http framework is golang.
Syntax is equivalent of Express JS.

Spin up HTTP server quickly in just few lines.


### To Run the project:

1) Setting up Mysql
   
  a) Get a Mysql Database server with a database
   
  Or
  
  b) Create a Mysql docker container and create a database inside it

Run `docker run --name mysql-gofiber -d -p 3307:3306 -e MYSQL_ROOT_PASSWORD=change-me mysql:8`

Run `docker ps` to find container id

SSH into container - Run `docker exec -it <container-id> sh`

Access mysql shell - Run `mysql -u root -p` (Provide root password if prompted)

Create Database - Run `CREATE DATABASE gofibersampledb;`


2) Update Database URI in database.go accordingly


3) Build whole project and generate executable output file

Run `go build -o gofiber-server .`


4) Run executable file directly

Run `./gofiber-server`


### Examples:

1) Get Users API without Authorization token - Works fine
<img width="1280" alt="Screenshot 2024-03-18 at 8 14 41 PM" src="https://github.com/ArnavRupde/gofiber-http-server/assets/34592221/897ae91e-2ce0-4ad5-aeac-5b7f327cd303">

2) POST Users API without Authorization token - Throws Unauthenticated error
   
<img width="1280" alt="Screenshot 2024-03-18 at 8 15 37 PM" src="https://github.com/ArnavRupde/gofiber-http-server/assets/34592221/d4af4820-80cc-42a2-bcf0-52ae6929477d">

3) POST Login API to get jwt
   
<img width="1280" alt="Screenshot 2024-03-18 at 8 15 49 PM" src="https://github.com/ArnavRupde/gofiber-http-server/assets/34592221/6853a553-2103-4c29-a88b-26ee3d491b41">

4) POST Users API with Authorization token - Works fine

<img width="1280" alt="Screenshot 2024-03-18 at 8 16 14 PM" src="https://github.com/ArnavRupde/gofiber-http-server/assets/34592221/d6749de6-1722-4180-9e98-1e7da6177947">

5) Connect over Websocket and Send message
<img width="1720" alt="Screenshot 2024-03-31 at 9 03 18 PM" src="https://github.com/ArnavRupde/gofiber-http-server/assets/34592221/4f864487-741c-4000-8647-84a50ad85b6a">
