package data

type Todo struct {
	ID   int
	Text string
}

func (d *Data) GetTodos() ([]Todo, error) {
	todos := []Todo{
		{ID: 1, Text: "Tidy up the house"},
		{ID: 2, Text: "Buy groceries"},
		{ID: 3, Text: "Walk the dog"},
	}

	return todos, nil
}
