package main

import (
	"fmt"
	"log"
	"net/http"
	"package30/lib30"
	"package30/server"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println(`hello server 2`)
	server.InitDB()
	fmt.Println(`db init`)
	r := chi.NewRouter()
	r.Get("/", lib30.Hello)
	r.Post("/create", lib30.CreateUser)
	r.Post("/make_friends", lib30.MakeFriends)
	r.Delete("/user/{id}", lib30.DeleteUser)
	r.Get("/friends/{id}", lib30.GetUserFriends)
	r.Put("/{id}", lib30.UpdateUserAge)

	http.ListenAndServe(":9001", r)

	//инициализируем подключение к базе данных
	err := server.InitDB()
	fmt.Println(err)
	if err != nil {
		log.Fatal(err)
	}

}
