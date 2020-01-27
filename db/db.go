package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Use PostgreSQL in gorm
	"golang.org/x/crypto/bcrypt"
)

var (
	db  *gorm.DB
	err error
)

type User struct {
	gorm.Model
	Name     string
	Email    string `gorm:"type:varchar(100);unique_index"`
	Password []byte
}

type History struct {
	gorm.Model
	Word   string
	UserID int
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

//connectInit はOpenの引数を作る
func connectInit() string {
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "taku")
	dbname := getEnv("DB_NAME", "postgres")
	password := getEnv("DB_PASS", "password")
	sslmode := getEnv("DB_SSLMODE", "disable")

	dbinfo := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		user,
		password,
		host,
		port,
		dbname,
		sslmode,
	)
	return dbinfo
}

// Init is initialize db from main function
func Init() {
	dbinfo := connectInit()
	db, err = gorm.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&User{}, &History{})
}

func CreateHistory(history History) {
	dbinfo := connectInit()
	db, err = gorm.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Create(&history)
	if db.NewRecord(history) == false {
		log.Printf("History %d Recorded\n", history.ID)
	} else {
		log.Printf("History  not %d Recorded\n", history.ID) //エラー内容を表示したい http://gorm.io/ja_JP/docs/error_handling.html
	}
}

func CreateUser(user User) {
	dbinfo := connectInit()
	db, err = gorm.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err := db.Create(&user).Error
	if err != nil {
		log.Printf("User %s not created: error: %s",
			user.Name,
			err,
		)
	}
}

func FindUser(email string, pass string) (User, error) {
	dbinfo := connectInit()
	db, err = gorm.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var user User
	db.Where("email = ?", email).First(&user)
	passError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	return user, passError
}

func GetHistory(id int) []string {
	dbinfo := connectInit()
	db, err = gorm.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var history []string
	db.Where("user_id = ?", id).Find(&history)
	return history
}
