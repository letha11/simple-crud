# Simple CRUD & Auth
A simple REST API that demonstrates CRUD operations and user authentication using Go. This project is just a learning project which the goal is to learn how to create a simple REST API using Go and getting familiar with the Go language.

## Getting started
1. Clone the repository
```bash
git clone https://github.com/your-username/simple-crud.git
```
2. Navigate to the project directory
```
cd simple-crud
```
3. Copy .env.example and rename it to .env
```bash
cp .env.example .env
```
4. Don't forget to add your own JWT Secret into the `.env` file
5. Install the required dependencies
```bash
go get -d ./...
```
6. Build the project
```bash
go build
```
7. Run the application
```bash
./simple-crud
```

Now you can access the API at http://localhost:5000/api

Documentation for the API can be found at http://localhost:5000/docs/index.html

