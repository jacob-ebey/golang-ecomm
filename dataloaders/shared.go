package dataloaders

import (
	"fmt"
	"strconv"
	"strings"
)

func HandleErrors(errors []error) error {
	var errs []string
	for _, e := range errors {
		errs = append(errs, e.Error())
	}

	return fmt.Errorf(strings.Join(errs, "\n"))
}

type IntKey int

func (key IntKey) String() string {
	return strconv.Itoa(int(key))
}

func (key IntKey) Raw() interface{} {
	return int(key)
}

type PaginationKey struct {
	Skip  int
	Limit int
}

func (key PaginationKey) String() string {
	return fmt.Sprintf("%d|%d", key.Skip, key.Limit)
}

func (key PaginationKey) Raw() interface{} {
	return key
}
