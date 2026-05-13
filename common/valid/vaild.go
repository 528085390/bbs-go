package valid

import "errors"

func IsValidString(str ...string) error {
	for _, s := range str {
		if s == "" {
			return errors.New("string can not be empty")
		}
	}
	return nil
}

type anyInt interface {
	int | int64 | uint | uint64 | int32 | uint32
}

func IsValidInt[T anyInt](Id ...T) error {
	for _, id := range Id {
		if id <= 0 {
			return errors.New("id must be greater than 0")
		}
	}
	return nil
}
