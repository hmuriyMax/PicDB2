package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"strings"
	"time"
)

var server = "localhost"
var port = 3306
var user = "maxim"
var password = "password6989"
var database = "usersPic"

var db *sql.DB

var ExpirationDuration = time.Hour * 24 * 60

func init() {
	var err error
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	db, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("ESTABLISHED CONNECTION AT ", server, ":", port, " db: ", database)
}

func CheckUser(userlog string, pass string) (int, error) {
	log.Println()
	var err error
	log.Println("searching for log/pass in database")
	rows, err := db.Query("SELECT id FROM passwords WHERE login = $1 AND password = $2", userlog, pass)
	if err != nil {
		return -2, err
	}
	err = rows.Close()
	if err != nil {
		return -2, err
	}
	var uid int
	found := false
	for rows.Next() {
		if found {
			log.Println("PAY ATTENTION!!! found more than two users")
			break
		}
		found = true
		err := rows.Scan(&uid)
		if err != nil {
			return -2, err
		}
	}
	if found {
		log.Println("MATCH with user with id = ", uid)
		return uid, nil
	}
	log.Println("user not found")
	return -1, nil
}

func GetRandomString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 15
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

func GetToken(id int) (string, error) {
	log.Println("asked new token for id", id)
	type tok struct {
		token   string
		userid  int
		expires string
	}
	var res tok

	res.token = GetRandomString()
	res.userid = id
	res.expires = time.Now().Add(ExpirationDuration).Format(time.ANSIC)

	row := db.QueryRow("SELECT * FROM Tokens WHERE token = $1", res.token)
	for row != nil {
		log.Println("token \"", res.token, "\" already found, generating another")
		res.token = GetRandomString()
		row = db.QueryRow("SELECT * FROM Tokens WHERE token = $1", res.token)
	}
	log.Println("NEW TOKEN \"", res.token, "\" for id", res.userid, " EXPIRES ", res.expires)

	str, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

func CheckToken(token string) (bool, error) {
	row := db.QueryRow("SELECT * FROM Tokens WHERE token = $1", token)
	if row != nil {
		return true, nil
	}
	return false, nil
}

func InsertUser(userlog string, pass string) (int, error) {
	rows, err := db.Query("SELECT id FROM passwords ORDER BY DESC")
	if err != nil {
		return 0, err
	}
	var newid int
	err = rows.Scan(&newid)
	if err != nil {
		return 0, err
	}
	newid++
	_, err = db.Exec("INSERT INTO password VALUES ($1, $2, $3)", newid, userlog, pass)
	return newid, err
}
