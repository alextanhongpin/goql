package goql

import (
	"fmt"
	"strconv"
)

func ParseInts(ss []string) ([]int, error) {
	res := make([]int, len(ss))

	for i, s := range ss {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidInt, s)
		}
		res[i] = n
	}

	return res, nil
}

func ParseBools(ss []string) ([]bool, error) {
	res := make([]bool, len(ss))

	for i, s := range ss {
		t, err := strconv.ParseBool(s)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidBool, s)
		}
		res[i] = t
	}

	return res, nil
}

func ParseFloat64s(ss []string) ([]float64, error) {
	res := make([]float64, len(ss))

	for i, s := range ss {
		f, err := strconv.ParseFloat(s, 10)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidFloat, s)
		}
		res[i] = f
	}

	return res, nil
}
