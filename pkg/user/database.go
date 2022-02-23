package user

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"strings"
	"time"
)

// TODO: скачать модуль viper... и добавить файлы конфига
var server = "localhost"
var port = 5432
var user = "maxim"
var password = "cringe2001"
var database = "usersPic"

var db *sql.DB

var ExpirationDuration = time.Hour * 24 * 60

type LToken struct {
	Token   string
	Userid  int
	Expires string
}

func init() {
	var err error
	connString := fmt.Sprintf("user=%s password=%s port=%d database=%s",
		user, password, port, database)
	log.Printf("RECEIVED RESPONCE to start server with: \n%s", connString)
	db, err = sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Printf("ESTABLISHED CONNECTION AT %s:%d db: %s\n\n", server, port, database)
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

func DBCheckUser(login string, pass string) (int, error) {
	var err error
	log.Printf("RECEIVED RESPONCE to check user %s/%s in database", login, GetMd5(pass))
	rows, err := db.Query("SELECT id FROM passwords WHERE (username = $1 OR email = $1) AND password = $2", login, GetMd5(pass))
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
	log.Println("FA" +
		"ILED to find user not found")
	return -1, nil
}

func DBGetToken(id int) (LToken, error) {
	log.Printf("RECEIVED RESPONCE to create token for id=%d ", id)

	var res LToken

	res.Userid = id
	res.Expires = time.Now().Add(ExpirationDuration).Format(time.ANSIC)
	res.Token = GetMd5(string(rune(id))) + GetMd5(res.Expires)

	var strr string
	row := db.QueryRow("SELECT Token FROM tokens WHERE Token = $1", res.Token)
	for row.Scan(&strr); strr == res.Token; {
		log.Println("Token \"", res.Token, "\" already found, generating another")
		res.Token = GetRandomString()
		row = db.QueryRow("SELECT * FROM tokens WHERE Token = $1", res.Token)
	}
	_, err := db.Exec("INSERT INTO tokens VALUES ($1, $2, $3)", res.Token, res.Userid, res.Expires)
	if err != nil {
		return LToken{}, err
	}
	log.Printf("NEW TOKEN \"%s\" for id=%d EXPIRES ON %s", res.Token, res.Userid, res.Expires)

	//str, err := json.Marshal(res)
	//if err != nil {
	//	return nil, err
	//}
	return res, nil
}

func DBCheckToken(token string) (bool, error) {
	log.Printf("RECEIVED RESPONCE to check token: %s ", token)
	var gotToken LToken
	row := db.QueryRow("SELECT token, id, expires FROM tokens WHERE token = $1", token)
	err := row.Scan(&gotToken.Token, &gotToken.Userid, &gotToken.Expires)
	if err != nil {
		log.Printf("FAILED token not found")
		return false, nil
	}
	tm, err := time.Parse(time.RFC3339, gotToken.Expires)
	if err != nil {
		return false, err
	}
	if gotToken.Token == token && tm.After(time.Now()) {
		log.Printf("SUCCESS: token %s for id=%d expires on %s", gotToken.Token, gotToken.Userid, gotToken.Expires)
		return true, nil
	}
	log.Printf("FAILED token expired!")

	return false, nil
}

func DBInsertUser(login string, pass string) (int, error) {
	log.Printf("RECEIVED RESPONCE to create new user w login: %s passhash: %s", login, GetMd5(pass))
	rows, err := db.Query("SELECT id FROM passwords WHERE username = $1", login)
	if err != nil {
		return 0, err
	}
	if rows.Next() {
		log.Printf("FAILED to create new user: %s aready exists", login)
		return -1, nil
	}
	rows, err = db.Query("SELECT id FROM passwords ORDER BY id DESC")
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
	_, err = db.Exec("INSERT INTO passwords VALUES ($1, $2, $3)", newid, login, GetMd5(pass))
	log.Printf("SUCCESS: new user %s with id=%d", login, newid)
	return newid, err
}

func DBUpdateUserData(userdata *Data) error {
	log.Printf("RECEIVED RESPONCE to update user data with: \n"+
		"id:            %d \n"+
		"name:          %s \n"+
		"email:         %s \n"+
		"birthday on:   %s \n"+
		"gender:        %s \n"+
		"profilePicURL: %s \n"+
		"unique key:    %s", userdata.userId, userdata.name, userdata.email, userdata.birthday, userdata.gender,
		userdata.profilePicURL, userdata.uniqueKey)
	rows, err := db.Query("SELECT id FROM user_info WHERE id = $1", userdata.userId)
	if err != nil {
		return err
	}
	if rows.Next() {
		command := fmt.Sprintf("UPDATE user_info SET name = '%s'", userdata.name)
		if userdata.birthday != "" {
			command += fmt.Sprintf(", birthday = '%s'", userdata.birthday)
		}
		if userdata.gender != "" {
			command += fmt.Sprintf(", gender = '%s'", userdata.gender)
		}
		if userdata.profilePicURL != "" {
			command += fmt.Sprintf(", userpic_url = '%s'", userdata.profilePicURL)
		}
		if userdata.uniqueKey != "" {
			command += fmt.Sprintf(", unique_key = '%s'", userdata.uniqueKey)
		}
		command += fmt.Sprintf(" WHERE id = %d", userdata.userId)
		_, err := db.Exec(command)
		if err != nil {
			return err
		}
	} else {
		command := fmt.Sprintf("INSERT INTO user_info VALUES (%d", userdata.userId)
		if userdata.gender != "" {
			command += fmt.Sprintf(", '%s'", userdata.gender)
		} else {
			command += fmt.Sprintf(", null")
		}
		if userdata.birthday != "" {
			command += fmt.Sprintf(", '%s'", userdata.birthday)
		} else {
			command += fmt.Sprintf(", null")
		}
		if userdata.uniqueKey != "" {
			command += fmt.Sprintf(", '%s'", userdata.uniqueKey)
		} else {
			command += fmt.Sprintf(", null")
		}
		if userdata.profilePicURL != "" {
			command += fmt.Sprintf(", '%s'", userdata.profilePicURL)
		} else {
			command += fmt.Sprintf(", DEFAULT")
		}
		command += fmt.Sprintf(", '%s')", userdata.name)
		_, err := db.Exec(command)
		if err != nil {
			return err
		}
	}
	log.Printf("SUCCESS!")
	_, err = db.Exec("UPDATE passwords SET email = $1 WHERE id = $2",
		userdata.email, userdata.userId)
	if err != nil {
		return err
	}
	return nil
}

func DBDeleteUser(id int32) error {
	log.Printf("RECEIVED RESPONCE to delete user w id=%d", id)
	_, err := db.Exec("DELETE from passwords WHERE id = $1", id)
	return err
}

func DBGetFullUserData(id int32) (*Data, error) {
	log.Printf("RECEIVED RESPONCE to get FULL userdata w id=%d", id)
	row := db.QueryRow("SELECT user_info.id, gender, birthday, unique_key, userpic_url, name, username, email FROM user_info RIGHT JOIN passwords ON user_info.id = passwords.id WHERE user_info.id = $1", id)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var ret Data
	err := row.Scan(&ret.userId, &ret.gender, &ret.birthday, &ret.uniqueKey,
		&ret.profilePicURL, &ret.name, &ret.username, &ret.email)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func DBGetShortUserData(id int32) (*Data, error) {
	log.Printf("RECEIVED RESPONCE to get SHORT userdata w id=%d", id)
	row := db.QueryRow("SELECT user_info.id, name, username, userpic_url FROM user_info RIGHT JOIN passwords ON (user_info.id = passwords.id) WHERE user_info.id = $1", id)
	if row.Err() != nil {
		return nil, row.Err()
	}
	var ret Data
	err := row.Scan(&ret.userId, &ret.name, &ret.username, &ret.profilePicURL)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
