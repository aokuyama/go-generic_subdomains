package errs

import (
	"errors"
	"fmt"
)

type Errors []error

func New() *Errors {
	errs := Errors{}
	return &errs
}

func (errs *Errors) Append(err error) {
	*errs = append(*errs, err)
}

func (errs *Errors) Err() error {
	var e string
	for _, err := range *errs {
		if err == nil {
			continue
		}
		e += fmt.Sprintf("%s\n", err.Error())
	}
	if len(e) == 0 {
		return nil
	}
	return errors.New(e)
}
