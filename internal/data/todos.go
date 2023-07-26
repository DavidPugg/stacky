package data

import "fmt"

type Todo struct {
	ID   int
	Text string
}

var todos = []Todo{
	{ID: 1, Text: "Tidy up the house"},
	{ID: 2, Text: "Buy groceries"},
	{ID: 3, Text: "Walk the dog"},
}

func (d *Data) GetTodos() ([]Todo, error) {
	return todos, nil
}

func (d *Data) AddTodo(text string) (Todo, error) {
	var ID int
	if len(todos) == 0 {
		ID = 1
	} else {
		ID = todos[len(todos)-1].ID + 1
	}

	todo := Todo{
		ID:   ID,
		Text: text,
	}

	todos = append(todos, todo)

	return todo, nil
}

func (d *Data) DeleteTodo(ID int) error {
	for i, todo := range todos {
		if todo.ID == ID {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("Todo with ID %d not found", ID)
}
