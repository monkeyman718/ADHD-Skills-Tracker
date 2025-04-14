package main

import (
	"encoding/json"
	//"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	//"github.com/joho/godotenv"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct{
    ID          uuid.UUID   `gorm:"type:uuid;primary key"`
    Username    string      `json:"username" gorm:"username;unique;not null"`
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

/*func init() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
}*/

func main() {
    ConnectDB()

    router := mux.NewRouter()    
    router.HandleFunc("/users", CreateUserHandler).Methods("POST")
    router.HandleFunc("/users", GetUsersHandler).Methods("GET")
    router.HandleFunc("/users/{email}", GetUserByIdHandler).Methods("GET")
    router.HandleFunc("/skills", CreateSkillHandler).Methods("POST")
    router.HandleFunc("/login", LoginHandler).Methods("POST")

    port := os.Getenv("PORT")
    
    if port == "" {
        port = "8080" // fallback for local dev
    }

    http.ListenAndServe(":"+port, enableCORS(router))
}

func enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Or restrict to your frontend's URL
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func ConnectDB() {
    DB, err = gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{Logger: logger.Default.LogMode(logger.Info),})
    if err != nil {
        panic("Error: Could not connect to database.")
    }  
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
    w.Header().Set("Content-Type","application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
	    "message": "User created!",
	    "User": user,
    })
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
    var users []User
	
    result := DB.Find(&users)
    if result.RowsAffected < 1 {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]error{"error": result.Error})
    }

    w.Header().Set("Content-Type","application/json")
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
    w.Header().Set("Content-Type","application/json")
    json.NewEncoder(w).Encode(&user)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    user := User{}
    dbUser := User{}

    json.NewDecoder(r.Body).Decode(&user)

    // get the password from the database for the email provided
    if err := DB.Where("email = ?", user.Email).First(&dbUser); err.Error != nil {
        http.Error(w, "Error: Email not found", http.StatusUnauthorized)
        return
    }

    if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
        http.Error(w, "Error: Invalid email or password", http.StatusUnauthorized)
        return
    }

    tokenStr, err := CreateJWT(w,dbUser.Email)
    if err != nil {
        http.Error(w, "Error creating jwt token", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "message": "LOGIN SUCCESSFUL",
        "jwt": tokenStr,
    })
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

    w.Header().Set("Content-Type","application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "New skill created!"})
}

func CreateJWT(w http.ResponseWriter, email string) (string, error) {
    jwtKey := []byte(os.Getenv("JWT_SECRET"))

    claims := jwt.MapClaims{
        "email": email,
        "exp": time.Now().Add(24 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Error: Creating token string", http.StatusNotFound)
        return "", err
    }

    return tokenString, nil
}