package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/go-chi/chi"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var (
	tasks      = make(map[int]Task)
	tasksMutex sync.Mutex
	nextID     = 1
)

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	var allTasks []Task
	for _, t := range tasks {
		allTasks = append(allTasks, t)
	}

	json.NewEncoder(w).Encode(allTasks)
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSpace(chi.URLParam(r, "id"))
	fmt.Println("Current tasks:", tasks)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}
	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	task, ok := tasks[id]
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var t Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tasksMutex.Lock()
	t.ID = nextID

	nextID++
	tasks[t.ID] = t
	tasksMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(t)

}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}
	tasksMutex.Lock()
	defer tasksMutex.Unlock()

	if _, ok := tasks[id]; !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return

	}
	delete(tasks, id)
	w.WriteHeader(http.StatusNoContent)
}
