package lib30

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"package30/server"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID      int    `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Age     int    `db:"age" json:"age"`
	Friends []int  `db:"friends" json:"friends"`
}

// var users = make(map[int]*User)

func unique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func Hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintln(w, "hello from 2 server")

}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	_, err2 := server.Db.Exec("INSERT INTO test_table(name, age) VALUES(?, ?)", user.Name, user.Age)
	if err2 != nil {
		fmt.Println(err.Error())
	}

	log.Printf(`User "%v" added to the database!`, user.Name)
}

// 2
func MakeFriends(w http.ResponseWriter, r *http.Request) {
	var friendRequest struct {
		SourceID int `json:"source_id"`
		TargetID int `json:"target_id"`
	}
	var user User
	var listfriends []byte
	var listfriends2 []byte

	err := json.NewDecoder(r.Body).Decode(&friendRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err)
	}

	//проверяем чтобы таргет и сурс не были одним и тем же
	if friendRequest.SourceID == friendRequest.TargetID {
		fmt.Fprint(w, `пользователь source_id такой же, как target_id`)
		return
	}

	// Получение текущего списка друзей
	err1 := server.Db.Get(&listfriends, "SELECT friends FROM test_table WHERE id=?", friendRequest.SourceID)
	if err1 != nil {
		log.Fatal(err1)
	}

	// Получение текущего списка друзей
	err3 := server.Db.Get(&listfriends2, "SELECT friends FROM test_table WHERE id=?", friendRequest.TargetID)

	//пример перехвата ошибок
	if err3 != nil {
		fmt.Fprintf(w, `пользователя с id %d не существует`, friendRequest.TargetID)
		return
	}

	//первый друг
	var friends []int
	err = json.Unmarshal([]byte(listfriends), &friends)
	if err != nil {
		// Handle error
	}
	// Добавление нового друга в список
	user.Friends = unique(append(friends, friendRequest.TargetID))

	//отправляем массив на просеивание, оставляем только уникальные идентификаторы
	uniqueSlice, err := json.Marshal(user.Friends)
	if err != nil {
		// Handle error
	}

	// Обновление списка друзей в базе данных
	_, err2 := server.Db.Exec("UPDATE test_table SET friends=? WHERE id = ?", uniqueSlice, friendRequest.SourceID)
	if err2 != nil {
		log.Fatal(err2)
	}

	//второй друг
	var friends2 []int
	err = json.Unmarshal([]byte(listfriends2), &friends2)
	if err != nil {
		// Handle error
	}
	// Добавление нового друга в список
	user.Friends = unique(append(friends2, friendRequest.SourceID))

	//отправляем массив на просеивание, оставляем только уникальные идентификаторы
	uniqueSlice2, err := json.Marshal(user.Friends)
	if err != nil {
		// Handle error
	}

	// Обновление списка друзей в базе данных
	_, err4 := server.Db.Exec("UPDATE test_table SET friends=? WHERE id = ?", uniqueSlice2, friendRequest.TargetID)
	if err4 != nil {
		log.Fatal(err4)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%s and %s are now friends", friendRequest.TargetID, friendRequest.SourceID)))
}

// 3
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "deleting")

	var deleteUserRequest struct {
		TargetID int `json:"target_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&deleteUserRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получение текущего списка друзей
	var friends []byte
	err2 := server.Db.Get(&friends, "SELECT friends FROM test_table WHERE id=?", deleteUserRequest.TargetID)
	if err2 != nil {
		fmt.Fprint(w, err2)
		return
	}
	fmt.Fprint(w, friends)

	// Обновление списка друзей в базе данных
	_, err3 := server.Db.Exec("UPDATE test_table SET friends = replace(friends, ?, '')", deleteUserRequest.TargetID)
	if err3 != nil {
		return
	}

	// Удаление пользователя из хранилища
	query := "DELETE FROM test_table WHERE id = ?"

	_, err4 := server.Db.Exec(query, deleteUserRequest.TargetID)
	if err2 != nil {
		fmt.Println(err4.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%v has been deleted", deleteUserRequest.TargetID)))
}

// 4
func GetUserFriends(w http.ResponseWriter, r *http.Request) {

	var listfriends []byte
	var listfriendsstring []byte

	// Извлекаем ID пользователя из URL-адреса
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получение текущего списка друзей
	err2 := server.Db.Get(&listfriends, "SELECT friends FROM test_table WHERE id = ?", userID)
	if err2 != nil {
		log.Fatal(err2)
	}

	//первый друг
	var friends []int
	err = json.Unmarshal([]byte(listfriends), &friends)
	if err != nil {
		// Handle error
	}

	for _, f := range friends {
		// Получение имя друзей и выводим
		err3 := server.Db.Get(&listfriendsstring, "SELECT name FROM test_table WHERE id = ?", f)
		if err3 != nil {
			log.Fatal(err3)
		}
		fmt.Fprintln(w, string(listfriendsstring[:]))

	}

}

// 5
func UpdateUserAge(w http.ResponseWriter, r *http.Request) {
	var users []User

	// Извлекаем ID пользователя из URL-адреса
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "SELECT * FROM test_table WHERE id = ?"
	sql2 := server.Db.Select(&users, query, userID)
	if sql2 != nil {
		panic(sql2.Error())
	}

	for _, user := range users {

		fmt.Printf("ID: %d, Name: %s, Age: %d, Friends: %s\n", user.ID, user.Name, user.Age, user.Friends)
	}

	// Декодируем JSON-тело запроса
	var requestBody map[string]string
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Извлекаем новый возраст пользователя из JSON-тела запроса
	newAge, err := strconv.Atoi(requestBody["age"])
	fmt.Printf("Устанавливаем возраст: %d\n", newAge)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обновляем возраст пользователя в мапе

	query_update := "UPDATE test_table SET age = ? WHERE id = ?"

	_, err2 := server.Db.Exec(query_update, newAge, userID)
	if err2 != nil {
		fmt.Println(err.Error())
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "возраст пользователя успешно обновлён")

}
