package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var (
	tasks     = make(map[int]Task)
	taskMutex sync.Mutex
	nextID    = 1
)

func GetAllTasks(w http.ResponseWriter, r *http.Request) {

	taskMutex.Lock()
	defer taskMutex.Unlock()

	var allTasks []Task
	for _, t := range tasks {
		allTasks = append(allTasks, t)
	}

	json.NewEncoder(w).Encode(allTasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var t Task

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	taskMutex.Lock()
	t.ID = nextID
	nextID++
	tasks[t.ID] = t
	taskMutex.Unlock()
	json.NewEncoder(w).Encode(t)

}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}

	taskMutex.Lock()
	defer taskMutex.Unlock()

	task, ok := tasks[id]
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(task)

}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}

	taskMutex.Lock()
	defer taskMutex.Unlock()
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	delete(tasks, id)
	w.WriteHeader(http.StatusNoContent)

}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id task", http.StatusBadRequest)
		return
	}
	var updated Task
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
	}

	taskMutex.Lock()
	defer taskMutex.Unlock()

	task, ok := tasks[id]
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	task.Title = updated.Title
	task.Done = updated.Done
	tasks[id] = task

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)

}

func GetAllFilter(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()
	filterDone := query.Get("done")
	taskMutex.Lock()
	defer taskMutex.Unlock()

	var results []Task

	for _, task := range tasks {
		if filterDone != "" {
			wantDone, err := strconv.ParseBool(filterDone)
			if err != nil {
				http.Error(w, "invalid 'done' query value", http.StatusBadRequest)
				return
			}
			if task.Done != wantDone {
				continue
			}
		}
		results = append(results, task)
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(results)
}
