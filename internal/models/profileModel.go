package models

type Todo struct {
	Id      int64
	UserId  int64
	Title   string
	Content string
	IsDone  int
}

type TodoQuery struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

func Page(q *TodoQuery) int {
	var page int
	if q.Limit > 0 {
		page = (q.Offset / q.Limit) + 1
	}
	return page
}

type CreateTodo struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	IsDone  int    `json:"is_done"`
}

type UpdateTodo struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	IsDone  int    `json:"is_done"`
}
