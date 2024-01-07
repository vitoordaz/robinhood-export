package utils

type GetIDFunc[T any] func(T) string

func GetIDs[T any](items []T, getIDFunc GetIDFunc[T]) []string {
	var (
		ids []string
		set = map[string]bool{}
	)
	for _, item := range items {
		if id := getIDFunc(item); !set[id] {
			ids = append(ids, id)
			set[id] = true
		}
	}
	return ids
}
