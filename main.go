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
    ID          uuid.UUID   `gorm:"type:uuid;primary key"`
    Email       string      `json:"email" gorm:"email;not null"`
    Password    string      `json:"password" gorm:"column:password_hash;not null"`
    CreatedAt   time.Time   `gorm:"created_at"`
    UpdatedAt   time.Time   `gorm:"updated_at"`
}

type Skill struct {
    ID          uuid.UUID   `gorm:"type:uuid;primary key"`
    UserID      uuid.UUID   `json:"user_id" gorm:"type:uuid;not null"`
    Name        string      `json:"name" gorm:"not null"`
    Priority    string      `json:"priority" gorm:"priority"`
    Goal        string      `json:"goal" gorm:"goal"`
    Status      string      `json:"status" gorm:"status"`
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
    router.HandleFunc("/skills", CreateSkillHandler).Methods("POST")
    router.HandleFunc("/login", LoginHandler).Methods("POST")

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
    if result := DB.Where("email = ?", email).First(&user); result != nil {
        http.Error(w,"Error: User not found", http.StatusNotFound)
        return
    }    

    // return user data as json
    json.NewEncoder(w).Encode(&user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    user := User{}
    dbUser := User{}

    json.NewDecoder(r.Body).Decode(&user)

    // get the password from the database for the email provided
    if err := DB.Where("email = ?", user.Email).First(&dbUser); err != nil {
        http.Error(w, "Error: Email not found", http.StatusUnauthorized)
        return
    }

    if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
        http.Error(w, "Error: Invalid email or password", http.StatusUnauthorized)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Login successful!"})
}

func CreateSkillHandler(w http.ResponseWriter, r *http.Request) {
    validPriorities := map[string]bool{
        "High": true,
        "Medium": true,
        "Low": true,
    }

    validStatus := map[string]bool{
        "Not Started": true,
        "In Progress": true,
        "Completed":   true,
    }

    skill := Skill{}
    

    json.NewDecoder(r.Body).Decode(&skill)
    skill.ID = uuid.New()
    
    if !validPriorities[skill.Priority] {
        http.Error(w, "Error: Invalid priority", http.StatusBadRequest)
        return
    }

    if !validStatus[skill.Status] {
        http.Error(w, "Error: Invalid status", http.StatusBadRequest)
        return
    }
    
    result := DB.Create(&skill)
    if result.Error != nil {
        http.Error(w, "Error: Skill not created", http.StatusNotModified)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "New skill created!"})
}