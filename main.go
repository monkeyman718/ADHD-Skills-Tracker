package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct{
    ID          uuid.UUID   `gorm:"id"`
    Email       string      `json:"email" gorm:"email"`
    Password    string      `json:"password" gorm:"column:password_hash"`
    CreatedAt   time.Time   `gorm:"created_at"`
    UpdatedAt   time.Time   `gorm:"updated_at"`
}

var DB *gorm.DB
var err error

func init() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
}

func main() {
    ConnectDB()

    router := mux.NewRouter()    
    router.HandleFunc("/users", CreateUserHandler).Methods("POST")
    router.HandleFunc("/users", GetUsersHandler).Methods("GET")
    router.HandleFunc("/users/{email}", GetUserByIdHandler).Methods("GET")

    fmt.Println("Listening on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}

func ConnectDB() {
    DB, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{Logger: logger.Default.LogMode(logger.Info),})
    if err != nil {
        panic("Error: Could not connect to database.")
    }   

    fmt.Println("Connected to database!")
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    user := User{}

    // Take the json body and put the info into a user variable
    json.NewDecoder(r.Body).Decode(&user)
    user.ID = uuid.New()

    // check that user info is valid data
    if user.Email == "" || string(user.Password) == "" {
        http.Error(w, "Error: email or password not valid entries.", http.StatusBadRequest) 
        return
    }

    // hash user password
    hashedPw, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    user.Password = string(hashedPw)
    // create user in database user table
    DB.Create(&user)

    // send info to the user saying the new user was created
    fmt.Fprintf(w, "User created!\n %v\n", user)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
    var users []User

    // store request in []users array

    // get users from database
    result := DB.Find(&users)
    if result.RowsAffected < 1 {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]error{"error": result.Error})
    }

    // return users as json response
    json.NewEncoder(w).Encode(&users)
}

func GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    email := vars["email"]
    user := User{}

    // search for user info with that id
    DB.Where("email = ?", email).First(&user)

    // return user data as json
    json.NewEncoder(w).Encode(&user)
}