package main

import (
    "html/template"
    "net/http"
    "strconv"
    "sync"
)

var (
    todos     []Todo
    idCounter int
    mu        sync.Mutex
)

func main() {
    http.HandleFunc("/", listTasks)
    http.HandleFunc("/add", addTask)
    http.HandleFunc("/done", markTaskAsDone)
    http.HandleFunc("/delete", deleteTask)
    http.ListenAndServe(":8080", nil)
}

func listTasks(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, todos)
}

func addTask(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    task := r.FormValue("task")
    if task == "" {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    mu.Lock()
    idCounter++
    todo := Todo{ID: idCounter, Task: task, Done: false}
    todos = append(todos, todo)
    mu.Unlock()

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func markTaskAsDone(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    mu.Lock()
    for i, todo := range todos {
        if todo.ID == id {
            todos[i].Done = true
            break
        }
    }
    mu.Unlock()

    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    mu.Lock()
    for i, todo := range todos {
        if todo.ID == id {
            todos = append(todos[:i], todos[i+1:]...)
            break
        }
    }
    mu.Unlock()

    http.Redirect(w, r, "/", http.StatusSeeOther)
}
