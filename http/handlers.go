package http

import (
	"net/http"

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

// pattern: /tasks/
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
