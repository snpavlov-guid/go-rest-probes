package util

type Numeric interface {
    int | int8 | int16 | int32 | int64 | float32 | float64 
}

// Универсальная функция для любого типа
func Ptr[T any](v T) *T {
    return &v
}

// Универсальная функция Map для срезов
func Map[T, U any](ts []T, f func(T) U) []U {
    us := make([]U, len(ts))
    for i := range ts {
        us[i] = f(ts[i])
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
