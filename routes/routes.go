package routes

import (
    "github.com/gorilla/mux"
    "ADHD-Skills-Tracker/handlers"
)

func CreateRoutes() *mux.Router {
    router := mux.NewRouter()

    router.HandleFunc("/users", handlers.CreateUserHandler).Methods("POST")
    //router.HandleFunc("/skills", handlers.GetSkillsHandler).Methods("GET")
    //router.HandleFunc("/skills", handlers.CreateSkillHandler).Methods("POST")
    //router.HandleFunc("/skills/{id}", handlers.UpdateSkillHandler).Methods("PUT")
    //router.HandleFunc("/skills/{id}", handlers.DeleteSkillHandler).Methods("DELETE")

    //router.HandleFunc("/users", handlers.GetUsersHandler).Methods("GET")
    //router.HandleFunc("/users/{id}", handlers.GetUserHandler).Methods("GET")
    //router.HandleFunc("/users", handlers.CreateUserHandler).Methods("POST")
    

    return router
}
