package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "go-sqlite3"
	"jwt-go"
	"net/http"
	"time"
)

var mySigningKey = []byte("captainjacksparrowsaysh123123i")

/*
type User struct {
	Id       int
	Name     string
	Surname  string
	Username string
	Password string
}
*/

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

/*
func getUser(db *sql.DB, id2 int) User {
	rows, err := db.Query("select * from users")
	checkErr(err)
	for rows.Next() {
		var tempUser User
		err = rows.Scan(&tempUser.Username, &tempUser.Password, &tempUser.Id)
		checkErr(err)
		if tempUser.Id == id2 {
			return tempUser
		}

	}
	return User{}
}
*/
//jsonHandler returns http respone in JSON format.
/*
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user := User{Id: 1,
		Name:  "John Doe",
		Email: "johndoe@gmail.com",
		Phone: "000099999"}
	json.NewEncoder(w).Encode(user)
}
*/
//templateHandler renders a template and returns as http response
/*.
func templateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("template.html")
	if err != nil {
		fmt.Fprintf(w, "Unable to load template")
	}

	user := User{Id: 1,
		Name:  "John Doe",
		Email: "johndoe@gmail.com",
		Phone: "000099999"}
	t.Execute(w, user)
}
*/
/*
type Kullanici struct {
	id       int
	name     string
	surname  string
	username string
	password string
	// 213123
}
*/
/*
func getKullanici(db *sql.DB) {

	rows, err := db.Query("SELECT * FROM users")

	checkErr(err)

	for rows.Next() {
		var temp Kullanici
		err = rows.Scan(&temp.id, &temp.name, &temp.surname, &temp.username, &temp.password)
		checkErr(err)
		//	fmt.Println(temp.name)
	}

}
*/

/*
func GenerateJWT() (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	tokenString, err := token.SignedString(mySigningKey)

	return tokenString, err
}

*/

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//fmt.Println(claims)
		return claims, err
	} else {
		return nil, err
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(string("letwitt service")))

}

type user struct {
	id       int
	name     string
	username string
	password string
}

func createToken(userClass user) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["id"] = userClass.id
	claims["username"] = userClass.username
	claims["password"] = userClass.password
	token.Claims = claims
	// Sign and get the complete encoded token as a stringg
	tokenString, err := token.SignedString(mySigningKey)
	return tokenString, err
}

func login(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")

	db, err := sql.Open("sqlite3", "./paylasim.db")

	checkErr(err)

	keys := r.URL.Query()
	username := keys.Get("username")
	password := keys.Get("password")

	rows, err := db.Query("SELECT * FROM users where username='" + username + "' and password='" + password + "'")

	checkErr(err)

	//	w.Header().Set("Content-Type", "application/json")
	//	w.WriteHeader(http.StatusOK)

	var tempUser user
	var count int = 0

	for rows.Next() {
		err = rows.Scan(&tempUser.id, &tempUser.name, &tempUser.username, &tempUser.password)
		checkErr(err)
		count++
		//w.Write([]byte(string(name)))
	}
	jsonMap := make(map[string]string)
	if count > 0 {
		jwt, err := createToken(tempUser)
		checkErr(err)
		jsonMap["message"] = "successful"
		jsonMap["jwt"] = jwt
		jsonMap["usernamesurname"] = tempUser.name
		//	json = "{\"message\":\"succesful\",\"jwt\":\"" + jwt + "\"}"
	} else {
		jsonMap["message"] = "unsuccessful"
	}
	json, err := json.Marshal(jsonMap)
	checkErr(err)
	w.Write([]byte(json))
}

func checkLogin(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Content-type", "multipart/form-data")

	db, err := sql.Open("sqlite3", "./paylasim.db")
	checkErr(err)

	keys := r.URL.Query()
	userToken := keys.Get("usertoken")

	userMap, err := ParseJWT(userToken)
	if err != nil {
		w.Write([]byte("{\"message\":\"false\"}"))
		return
	}

	username := fmt.Sprintf("%v", userMap["username"])
	password := fmt.Sprintf("%v", userMap["password"])

	rows, err := db.Query("SELECT * FROM users where username='" + username + "' and password='" + password + "'")
	checkErr(err)

	var count int = 0

	for rows.Next() {
		count++
	}

	if count > 0 {
		w.Write([]byte("{\"message\":\"true\"}"))
	} else {
		w.Write([]byte("{\"message\":\"false\"}"))
	}

}

type Post struct {
	Id     string
	Post   string
	Userid string
	User   string
	Date   string
}

func getposts(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")

	db, err := sql.Open("sqlite3", "./paylasim.db")

	checkErr(err)
	/*	mapD := map[string]int{"apple": 5, "lettuce": 7}
		mapB, _ := json.Marshal(mapD)
		fmt.Println(string(mapB))

		w.Write([]byte(string(mapB)))
	*/

	rows, err := db.Query("SELECT posts.*,users.name FROM posts inner join users on users.id=posts.userid order by posts.id desc")

	checkErr(err)
	posts := []Post{}

	for rows.Next() {
		var temp Post
		err = rows.Scan(&temp.Id, &temp.Post, &temp.Userid, &temp.Date, &temp.User)
		checkErr(err)
		posts = append(posts, temp)
	}

	/*	for i := range posts {
		//fmt.Println(posts[i].post)
		w.Write([]byte(string(posts[i].post)))
	}*/

	//	fmt.Println(posts[0].post)

	pagesJson, err := json.Marshal(posts)
	checkErr(err)

	w.Write([]byte(pagesJson))

}

func addPost(w http.ResponseWriter, r *http.Request) {

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")

	keys := r.URL.Query()
	userToken := keys.Get("usertoken")
	postText := keys.Get("post")

	userMap, err := ParseJWT(userToken)
	if err != nil {
		w.Write([]byte("{\"message\":\"unsuccessful\"}"))
		return
	}

	userid := fmt.Sprintf("%v", userMap["id"])
	date := time.Now()
	date2 := date.Format("2006-01-02")
	db, err := sql.Open("sqlite3", "./paylasim.db")

	checkErr(err)

	tx, _ := db.Begin()
	stmt, _ := tx.Prepare("insert into posts (post,userid,date) values (?,?,?)")
	_, err = stmt.Exec(postText, userid, date2)
	checkErr(err)
	tx.Commit()

	w.Write([]byte("{\"message\":\"successful\"}"))
}

func main() {

	http.HandleFunc("/", handler)
	http.HandleFunc("/posts", getposts)
	http.HandleFunc("/login", login)
	http.HandleFunc("/checkLogin", checkLogin)
	http.HandleFunc("/addPost", addPost)
	http.ListenAndServe(":8081", nil)
}
