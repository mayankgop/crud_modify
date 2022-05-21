package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Id          int    `json:"id"`
	First_name  string `json:"fname"`
	Last_name   string `json:"lname"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Dob         string `json:"dob"`
	Created_at  string `json:"createdat"`
	Last_access string `json:"last_access"`
	Updated_at  string `json:"updated"`
	Archived    string `json:"archived"`
}

var db *sql.DB
var err error

func init() {

	db, err = sql.Open("mysql", "admin:qwerty123@tcp(localhost:3306)/bookstore?charset=utf8")
	if err != nil {
		fmt.Println("error in opening database")
	}
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/login",login).Methods("GET")
	r.HandleFunc("/signup", signup)


	r.HandleFunc("/create", create) //inserting values into users_data table
	 r.HandleFunc("/read", read).Methods("GET")          // reading all values from users_data table
	 r.HandleFunc("/del/{mid}", delet).Methods("GET") // deleting values by id from books table
	// r.HandleFunc("/book/{mid}", getbyid).Methods("GET")  // reading values by id from books table
	// r.HandleFunc("/book/{mid}", upd).Methods("PUT")      // updating values by id from books table
	r.Handle("/favicon.ico", http.NotFoundHandler())


	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
	}
}

func create(w http.ResponseWriter, req *http.Request) {
	d:= getUser(w, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "bar.gohtml", d)
	var u user
	_ = json.NewDecoder(req.Body).Decode(&u)


	bs, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	u.Password=string(bs)
	rows,err:=db.Query("select count(*) from users_data where ?=users_data.email",u.Email)
	if err!=nil{
		log.Fatal(err)
	}
	var count int
	defer rows.Close()
	for rows.Next(){
		if err:=rows.Scan(&count);err!=nil{
			log.Fatal(err)
		}
	}
	// fmt.Printf("total count is %d",count)
	if count>0{
		fmt.Println("email already in use")
		return
	}else{

		// l:=len(u.Email)
		
	}



	query := "Insert into users_data(user_id,first_name,last_name,email,passw,dob,archived) values(?,?,?,?,?,?,?)"
	_, b := db.Exec(query,u.Id,u.First_name, u.Last_name, u.Email,u.Password,u.Dob,u.Archived) //cascade injection

	if b != nil {
		fmt.Println("error found in create\n", b)
		return
	}

	uj, _ := json.Marshal(u)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", uj)
}

func read(w http.ResponseWriter, req *http.Request) {
	var u user
	err := json.NewDecoder(req.Body).Decode(&u)
	
	d:= getUser(w, req)
	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "bar.gohtml", d)


	db, err := sql.Open("mysql", "admin:qwerty123@tcp(localhost:3306)/bookstore?charset=utf8")
	if err != nil {
		fmt.Println(err)
	}

	x, b := db.Query("select * from users_data")

	for x.Next() {
		var u user
		if err := x.Scan(&u.Id, &u.First_name,&u.Last_name,&u.Email,&u.Password,&u.Dob,&u.Created_at,&u.Last_access,&u.Updated_at,&u.Archived); err != nil {
			log.Fatal(err)
		}
		uj, _ := json.Marshal(u)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "%s", uj)
	}

	if b != nil {
		fmt.Println("error found in read")
	}
	
}
func delet(w http.ResponseWriter, req *http.Request) {
	var u user
	m := mux.Vars(req)
	fmt.Println(m["mid"])
	err := json.NewDecoder(req.Body).Decode(&u)
	// fmt.Println(u)
	db, err := sql.Open("mysql", "admin:qwerty123@tcp(localhost:3306)/bookstore?charset=utf8")
	if err != nil {
		fmt.Println(err)
	}
	x, _ := strconv.Atoi(m["mid"])
	query := fmt.Sprintf("DELETE FROM users_data WHERE user_id=%d", x)
	_, b := db.Exec(query)

	if b != nil {
		fmt.Println("error found in pbook")
	}

	uj, _ := json.Marshal(&u)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", uj)
}

// func getbyid(w http.ResponseWriter, req *http.Request) {
// 	var u user
// 	m := mux.Vars(req)
// 	// fmt.Println(m["mid"])
// 	err := json.NewDecoder(req.Body).Decode(&u)
// 	// fmt.Println(u)
// 	db, err := sql.Open("mysql", "admin:qwerty123@tcp(localhost:3306)/bookstore?charset=utf8")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	x, _ := strconv.Atoi(m["mid"])
// 	query := fmt.Sprintf("select * from books where id=%d", x)
// 	row, b := db.Query(query)

// 	for row.Next() {
// 		row.Scan(&u.Id, &u.Name)
// 	}

// 	if b != nil {
// 		fmt.Println("error found in get by id", b)
// 	}

// 	uj, _ := json.Marshal(&u)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	fmt.Fprintf(w, "%s", uj)
// }


// func upd(w http.ResponseWriter, req *http.Request) {
// 	var u user
// 	m := mux.Vars(req)
// 	// fmt.Println(m["mid"])
// 	err := json.NewDecoder(req.Body).Decode(&u)
// 	// fmt.Println(u)
// 	db, err := sql.Open("mysql", "admin:qwerty123@tcp(localhost:3306)/bookstore?charset=utf8")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	x, _ := strconv.Atoi(m["mid"])
// 	query := fmt.Sprintf("update books set bookname='%s' where id=%d", u.Name, x)
// 	_, b := db.Exec(query)

// 	if b != nil {
// 		fmt.Println("error found in pbook")
// 	}

// 	uj, _ := json.Marshal(&u)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	fmt.Fprintf(w, "%s", uj)
// }
