package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/theun1c/HundredthRest/todo"
)

type HTTPHandlers struct {
	todoList *todo.List
}

func NewHTTPHandlers(todoList *todo.List) *HTTPHandlers {
	return &HTTPHandlers{
		todoList: todoList,
	}
}

// pattern: /tasks
// method: POST
// info: JSON from request body
//
// suc:
//
//   - status code: 201 Created
//
//   - response body: JSON represent created task
//
//     fail:
//
//   - status code: 400, 500 ...
//
//   - response body: JSON is error + time
func (h *HTTPHandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		errDTO := NewErrorDTO(err.Error())
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := taskDTO.ValidateOnCreate(); err != nil {
		errDTO := NewErrorDTO(err.Error())
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	todoTask := todo.NewTask(taskDTO.Title, taskDTO.Text)
	if err := h.todoList.AddTast(todoTask); err != nil {
		errDTO := NewErrorDTO(err.Error())

		if errors.Is(err, todo.ErrTaskAlreadyExist) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(todoTask, "", "    ")

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response")
	}

}

// pattern: /tasks/{title}
// method: GET
// info: title from pattern
//
// suc:
//   - status code: 200
//   - response body: JSON is suc
//
// fail:
//   - status code: 400, 500 ...
//   - response body: JSON is error + time
func (h *HTTPHandlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	// mux.Vars(r) <- возвращает мапу ключ-значение
	// где ключом является название в {Title} в /tasks/{title}
	// а значением является то, что стоит на этом месте

	title := mux.Vars(r)["title"]

	task, err := h.todoList.GetTask(title)
	if err != nil {
		errDTO := NewErrorDTO(err.Error())
		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}

	b, err := json.MarshalIndent(task, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response")
		return
	}

}

// pattern: /tasks
// method: GET
// info: -
//
// suc:
//   - status code: 200
//   - response body: JSON is suc
//
// fail:
//   - status code: 400, 500 ...
//   - response body: JSON is error + time
func (h *HTTPHandlers) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.todoList.GetTasks()
	b, err := json.MarshalIndent(tasks, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response")
		return
	}
}

// pattern: /tasks?completed=true
// method: GET
// info: query params
//
// suc:
//   - status code 200
//   - response body: JSON is suc
//
// fail:
//   - status code: 400, 500 ...
//   - response body: JSON is error + time
func (h *HTTPHandlers) HandleGetUncompletedTasks(w http.ResponseWriter, r *http.Request) {
	uncompletedTasks := h.todoList.ListNotCompletedTasks()

	b, err := json.MarshalIndent(uncompletedTasks, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write HTTP response")
		return
	}

}

// pattern: /tasks/{title}
// method: PATCH
// info: from pattern and JSON with completed field
//
// suc:
//   - status code: 200
//   - response body: JSON is suc
//
// fail:
//   - status code: 400, 500 ...
//   - response body: JSON is error + time
func (h *HTTPHandlers) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	var completeDTO CompleteTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&completeDTO); err != nil {
		errDTO := NewErrorDTO(err.Error())

		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if completeDTO.Complete {
		if err := h.todoList.CompleteTask(title); err != nil {
			errDTO := NewErrorDTO(err.Error())
			if errors.Is(err, todo.ErrTaskNotFound) {
				http.Error(w, errDTO.ToString(), http.StatusNotFound)
			} else {
				http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			}

			return
		}
	} else {
		if err := h.todoList.UncompleteTask(title); err != nil {
			errDTO := NewErrorDTO(err.Error())
			if errors.Is(err, todo.ErrTaskNotFound) {
				http.Error(w, errDTO.ToString(), http.StatusNotFound)
			} else {
				http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			}

			return
		}
	}
}

// pattern: /tasks/{title}
// method: DELETE
// info: in pattern
//
// suc:
//   - status code: 200
//   - responde body: -
//
// fail:
//   - status code: 400, 500, ...
//   - response body: JSON error + time
func (h *HTTPHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	if err := h.todoList.DeleteTask(title); err != nil {
		errDTO := NewErrorDTO(err.Error())
		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}

		return
	}
}
