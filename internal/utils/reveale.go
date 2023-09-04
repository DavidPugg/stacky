package utils

type ObjectWithTotalCount interface {
	GetTotalCount() int
}

type RevealeObject[T ObjectWithTotalCount] struct {
	Data     T    `json:"data"`
	IsLast   bool `json:"is_last"`
	Page     int  `json:"page"`
	LastPage bool `json:"last_page"`
}

func CreateRevealeObjects[T ObjectWithTotalCount](objects []T, page, limit int) []RevealeObject[T] {
	revealeObjects := []RevealeObject[T]{}
	for i, c := range objects {
		revealeObjects = append(
			revealeObjects,
			RevealeObject[T]{
				Data:     c,
				IsLast:   i == len(objects)-1,
				Page:     page + 1,
				LastPage: c.GetTotalCount() <= limit*page,
			})
	}

	return revealeObjects
}
