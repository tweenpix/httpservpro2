package server

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// глобальная переменная с подключением к БД
var Db *sqlx.DB

// функция, инициирующая подключение к БД
func InitDB() (err error) {
	//строка, содержащая данные для подключения к БД в следующем формате:
	//login:password@tcp(host:port)/dbname
	var dataSourceName = "root@tcp(localhost:3306)/test"
	//подключаемся к БД, используя нужный драйвер и данные для подключения
	Db, err = sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		return
	}
	err = Db.Ping()
	log.Println(`Обзоначена бд:`, Db)
	return
}
