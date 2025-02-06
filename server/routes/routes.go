package routes

import (
	"encoding/json"
	"net/http"
	"timelygator/server/database"
	"timelygator/server/models"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", createTask).Methods("POST")
	r.HandleFunc("/tasks", getTasks).Methods("GET")
}

// CreateTask godoc
//
// @Summary		Create a new task
// @Description	create new task, returns the created task
// @Tags		tasks
// @Accept		json
// @Produce		json
// @Param		task body models.Task true "Task object that needs to be created"
// @Success		200	{object}	models.Task
// @Failure		400	{object}	middleware.HTTPError
// @Failure		404	{object}	middleware.HTTPError
// @Failure		500	{object}	middleware.HTTPError
// @Router		/tasks [post]
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

// GetTasks godoc
// @Summary 	Get all tasks
// @Description Get all tasks
// @Tags 		tasks
// @Produce 	json
// @Success 	200	{object}	[]models.Task
// @Failure 	400	{object}	middleware.HTTPError
// @Failure 	404	{object}	middleware.HTTPError
// @Failure 	500	{object}	middleware.HTTPError
// @Router		/tasks [get]
func getTasks(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	database.DB.Find(&tasks)
	json.NewEncoder(w).Encode(tasks)
}
