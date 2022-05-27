package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

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
	file,_:=os.Create("file.log")
	log.SetOutput(file)
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", index)
	r.HandleFunc("/login",login)
	r.HandleFunc("/signup", signup)
	r.HandleFunc("/logout", logout)


	r.HandleFunc("/create", create) //inserting values into users_data table
	 r.HandleFunc("/read/{l}/{o}", read)        // reading all values from users_data table
	 r.HandleFunc("/logout", logout)   //running read template

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

	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	fmt.Println("means logged in")


	var u user
	if req.Method == http.MethodPost {

	// tpl.ExecuteTemplate(w, "bar.gohtml", u)
	// _ = json.NewDecoder(req.Body).Decode(&u)


	u.Email = req.FormValue("Email")
	u.First_name= req.FormValue("password")
	u.Last_name = req.FormValue("f_name")
	u.Password= req.FormValue("l_name")
	u.Dob= req.FormValue("dob")

	fmt.Println("email is ",u.Email)
	e:=req.FormValue("Email")
	fmt.Println(e)
	f:=req.FormValue("f_name")
	fmt.Println(f)
	fmt.Println("reached after email")
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
		log.Println("email already in use")
		fmt.Println("email already in use")
		return
	}
	l:=len(u.Email)
	if l>20{
		log.Println("length greater than 20")
		return
	}
	u.Archived="fal"

	lf:=len(u.First_name)  //length of first name
	ll:=len(u.Last_name)   //length of last name

	if lf+ll>30{
		log.Println("name should be less than 30")
		return

	}
	lp:=len(u.Password)
	if lp<8 || lp>20 {
		log.Println("password not in range")
		return
	}
	fmt.Println("reached before query")
		
	query := "Insert into users_data(first_name,last_name,email,passw,dob,archived) values(?,?,?,?,?,?)"
	_, err = db.Exec(query,u.First_name, u.Last_name, u.Email,u.Password,u.Dob,u.Archived) //cascade injection

	
	if err!= nil {
		log.Println("error found in create\n", err)
		return
	}

	http.Redirect(w, req, "/r", http.StatusSeeOther)
}
	tpl.ExecuteTemplate(w, "entry.gohtml", u)

	// uj, _ := json.Marshal(u)
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// fmt.Fprintf(w, "%s", uj)
}



// func r(w http.ResponseWriter, req *http.Request) {
// 	var u user
// 	tpl.ExecuteTemplate(w, "limit.gohtml", u)
	
// }



func read(w http.ResponseWriter, req *http.Request) {
	var u user
	m := mux.Vars(req)
	fmt.Println(m["l"])
	fmt.Println(m["o"])
	l, _ := strconv.Atoi(m["l"])
	o, _ := strconv.Atoi(m["o"])



	_= json.NewDecoder(req.Body).Decode(&u)
	// fmt.Println(u)


	// _= getUser(w, req)

	if !alreadyLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}
	fmt.Println("rached in read")

	// if req.Method == http.MethodPost {
	// 	l:=req.FormValue("limit")    //limit value
	// 	o:=req.FormValue("offset")    //offset value
	// 	fmt.Print("limit is ",l)

	rows, b := db.Query("select * from users_data limit ? offset ?",l,o)

	for rows.Next() {
		var u user
		if err := rows.Scan(&u.Id, &u.First_name,&u.Last_name,&u.Email,&u.Password,&u.Dob,&u.Created_at,&u.Last_access,&u.Updated_at,&u.Archived); err != nil {
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

	// http.Redirect(w, req, "/r", http.StatusSeeOther)

	// }	
	// tpl.ExecuteTemplate(w, "limit.gohtml", u)
}






func delet(w http.ResponseWriter, req *http.Request) {
	var u user
	m := mux.Vars(req)
	fmt.Println(m["mid"])
	_= json.NewDecoder(req.Body).Decode(&u)
	// fmt.Println(u)

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
// 	db, err := sql.Open("mysql", "admin:qwerty123@tcp(localhost:3306)/users_data?charset=utf8")
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
// 	db, err := sql.Open("mysql", "admin:qwerty123@tcp(localhost:3306)/users_data?charset=utf8")
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
