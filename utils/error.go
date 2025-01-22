package utils

import (
	"errors"
	"strconv"
)

func GenErrorLine(i int) error {
	return errors.New("line " + strconv.Itoa(i+1) + " has an error")
}
