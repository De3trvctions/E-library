# E-library
A simple RESTful API to manage loan of e-book in an electronic library.


1. Use `docker compose up` command to build the db.
2. If you have your own db, you can change the configuration on `init/db` 
3. To run the project, generate the routers file with `bee generate routers`
    PS: You should be downloading all the needed packages by `go mod tidy`. Be sure that you have `Bee` installed.
        If you dont have the bee installed, change the `routers/commentsRouter.txt` to  `routers/commentsRouter.go` to generate the route path
4. Then run the peohect with `go run main.go` OR `bee run`
5. You should be able to browse the API on localhost:3000/{API name}
6. To run the unit test, just do run code `go test -v ./tests`