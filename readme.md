-How to use

#u1
go run main.go //with default value
go run main.go -file hindi.xml -genre drama //with custom value


#u2, u3 and u4
- go to movie folder
- go to server folder
- go to .env file and set the json file path
- go run main.go

- go to client folder
- go run main.go

- import the postman suite
- Load - http://localhost:8080/movie-library/load (post)
- Fetch all - http://localhost:8080/movie-library/movie/ (get)
- Fetch by filter - http://localhost:8080/movie-library/movie/01-10-2023 (get)
- Update - http://localhost:8080/movie-library/movie/2 (post)