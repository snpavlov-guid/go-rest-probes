package util

type Numeric interface {
    int | int8 | int16 | int32 | int64 | float32 | float64 
}

// Универсальная функция для любого типа
func Ptr[T any](v T) *T {
    return &v
}

func PrtOrNil[T, U any](v* T, f func(T) *U) *U {
    if v == nil {
        return  nil
    }
    return f(*v)
}

func Filter[T any](ts []T, f func(T) bool) []T {
    var items []T
    for _, item := range ts {
        if f(item) {
            items = append(items, item)
        }
    }
    return items
}

// Универсальная функция Map для срезов
func Map[T, U any](ts []T, f func(T) U) []U {
    us := make([]U, len(ts))
    for i := range ts {
        us[i] = f(ts[i])
    }
    return us
}

// Универсальная функция Map для срезов
func Map2[T, U any](ts []T, f func(T) (U, error)) ([]U, error) {
    us := make([]U, len(ts))
    for i := range ts {
        itm, err := f(ts[i])
        if (err != nil) {
            return nil, err
        }
        us[i] = itm
    }
    return us, nil
}

// Универсальная функция MapMany для срезов
func MapMany[T, C any, U any](ts []T, s func(T) []C, f func(C) U) []U {
    us := make([]U, 0)

    for i := range ts {
        iar := s(ts[i])
        for j := range iar {
            us = append(us, f(iar[j]))
        }
    }
   
    return us
}

// Универсальная функция MapMany для map
func MapMany2[K comparable, V any, C any, U any](ts map[K]V, s func(V) []C, f func(C) U) []U {
    us := make([]U, 0)

    for _, vi := range ts {
        for _, iar  := range s(vi) {
            us = append(us, f(iar))
        }
    }
   
    return us
}

func SliceToMap[T any, K comparable](slice []T, keyFunc func(T) K) map[K][]T {
    result := make(map[K][]T)
    for _, item := range slice {
		key := keyFunc(item)
        result[key] = append(result[key], item)
    }
    return result
}

func Sum[T any, R Numeric](ts []T, f func(T) R) R {
	var sum R = 0
    for i := range ts {
        sum += f(ts[i])
    }
    return sum	 
}

func Concat[T any](slices ...[]T) []T {
    var result []T
    for _, slice := range slices {
        result = append(result, slice...)
    }
    return result
}

func Distict[T any, K comparable](slice []T, keyFunc func(T) K) []T {
    var result []T
    cmap := make(map[K]struct{}) 

    for _, item := range slice {
        key := keyFunc(item)
        _, exists := cmap[key]
        if !exists {
            cmap[key] = struct{}{}
            result = append(result, item)
        }
     
    }
    return result
}

func DistictSet[T any, K comparable](slice []T, keyFunc func(T) K) map[K]struct{} {
    cmap := make(map[K]struct{}) 

    for _, item := range slice {
        key := keyFunc(item)
        _, exists := cmap[key]
        if !exists {
            cmap[key] = struct{}{}
        }
     
    }
    return cmap
}
