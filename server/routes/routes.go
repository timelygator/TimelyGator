package routes

import (
	"encoding/json"
	"net/http"
	"timelygator/server/database"
	"timelygator/server/database/models"
	"timelygator/server/utils/types"

	"github.com/gorilla/mux"
)

func RegisterRoutes(cfg types.Config, datastore *database.Datastore, r *mux.Router) {
	r.HandleFunc("/tasks", createTask).Methods("POST")
	r.HandleFunc("/tasks", getTasks).Methods("GET")
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	database.DB.Create(&task)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	database.DB.Find(&tasks)
	json.NewEncoder(w).Encode(tasks)
}
