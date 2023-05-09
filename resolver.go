package rushia

import "context"

// ResolveContext
func ResolveContext[T any, M any](c context.Context, data []T, fetch func(T) string, keyer func(T, M) bool, swap func(T, M), solver func(context.Context, []string) ([]M, error)) error {
	//
	if len(data) == 0 {
		return nil
	}
	//
	ids := make([]string, len(data))
	for i, v := range data {
		id := fetch(v)
		if id == "" {
			continue
		}
		ids[i] = id
	}
	//
	solved, err := solver(c, ids)
	if err != nil {
		return err
	}
	//
	for _, v := range data {
		for _, j := range solved {
			if keyer(v, j) {
				swap(v, j)
				break
			}
		}
	}
	return nil
}

// Resolve
func Resolve[T any, M any](data []T, fetch func(T) string, keyer func(T, M) bool, swap func(T, M), solver func([]string) ([]M, error)) error {
	newSovler := func(_ context.Context, ids []string) ([]M, error) {
		return solver(ids)
	}
	return ResolveContext(context.Background(), data, fetch, keyer, swap, newSovler)
}
