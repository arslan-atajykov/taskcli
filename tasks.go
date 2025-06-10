package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT id, title, done FROM tasks")
	if err != nil {
		http.Error(w, "failed to fetch tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Done)
		if err != nil {
			http.Error(w, "scan error", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, t)
	}
	json.NewEncoder(w).Encode(&tasks)

}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var t Task

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	id, err := InserTask(t.Title, t.Done)
	if err != nil {
		http.Error(w, "failed to insert task", http.StatusInternalServerError)
		return
	}
	t.ID = id
	json.NewEncoder(w).Encode(t)

}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}

	var t Task
	err = db.QueryRow("SELECT id, title, done FROM tasks WHERE id = ?", id).Scan(&t.ID, &t.Title, &t.Done)
	if err != sql.ErrNoRows {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(t)

}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var t Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE tasks SET title = ?, done = ? WHERE id = ?", t.Title, t.Done, id)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	t.ID = id
	json.NewEncoder(w).Encode(t)
}

func GetAllFilter(w http.ResponseWriter, r *http.Request) {
	filterDone := r.URL.Query().Get("done")

	query := "SELECT id, title, done FROM tasks"
	var args []interface{}

	if filterDone != "" {
		query += " WHERE done = ?"
		val, err := strconv.ParseBool(filterDone)
		if err != nil {
			http.Error(w, "invalid 'done' value", http.StatusBadRequest)
			return
		}
		args = append(args, val)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Done); err != nil {
			http.Error(w, "scan error", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, t)
	}

	json.NewEncoder(w).Encode(tasks)
}
