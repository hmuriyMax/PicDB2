package auth

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"strings"
	"time"
)

var server = "localhost"
var port = 5432
var user = "maxim"
var password = "xamburger6989"
var database = "usersPic"

var db *sql.DB

var ExpirationDuration = time.Hour * 24 * 60

type tok struct {
	Token   string
	Userid  int
	Expires string
}

func init() {
	var err error
	connString := fmt.Sprintf("user=%s password=%s port=%d database=%s",
		user, password, port, database)
	fmt.Printf("%s\n", connString)
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
	var err error
	log.Printf("searching for %s/%s in database", userlog, pass)
	rows, err := db.Query("SELECT id FROM passwords WHERE login = $1 AND password = $2", userlog, GetMd5(pass))
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
		log.Printf("MATCH with user with id=%d", uid)
		return uid, nil
	}
	err = rows.Close()
	if err != nil {
		return -2, err
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

func GetMd5(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func GetToken(id int) (string, error) {
	//log.Println("asked new Token for id", id)

	var res tok

	res.Userid = id
	res.Expires = time.Now().Add(ExpirationDuration).Format(time.ANSIC)
	res.Token = GetMd5(string(rune(id))) + GetRandomString()

	var strr string
	row := db.QueryRow("SELECT Token FROM tokens WHERE Token = $1", res.Token)
	for row.Scan(&strr); strr == res.Token; {
		log.Println("Token \"", res.Token, "\" already found, generating another")
		res.Token = GetRandomString()
		row = db.QueryRow("SELECT * FROM tokens WHERE Token = $1", res.Token)
	}
	_, err := db.Exec("INSERT INTO tokens VALUES ($1, $2, $3)", res.Token, res.Userid, res.Expires)
	if err != nil {
		return "", err
	}
	log.Printf("NEW TOKEN \"%s\" for id=%d EXPIRES ON %s", res.Token, res.Userid, res.Expires)

	str, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(str), nil
}

func CheckToken(token string) (bool, error) {
	row := db.QueryRow("SELECT * FROM Tokens WHERE Token = $1", token)
	if row != nil {
		return true, nil
	}
	return false, nil
}

func InsertUser(userlog string, pass string) (int, error) {
	rows, err := db.Query("SELECT id FROM passwords ORDER BY id DESC")
	if err != nil {
		return 0, err
	}
	newid := 0
	if rows.Next() {
		err = rows.Scan(&newid)
		if err != nil {
			return 0, err
		}
	}
	newid++
	_, err = db.Exec("INSERT INTO passwords VALUES ($1, $2, $3)", newid, userlog, GetMd5(pass))
	return newid, err
}
