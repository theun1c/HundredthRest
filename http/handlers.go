package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

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
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	if err := taskDTO.ValidateOnCreate(); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
	}

	todoTask := todo.NewTask(taskDTO.Title, taskDTO.Text)
	if err := h.todoList.AddTast(todoTask); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}

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

}

// pattern: /tasks?completed=bool
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

}
