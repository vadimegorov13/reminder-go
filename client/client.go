package client

import "fmt"

func wrapError(msg string, err error) error {
	return fmt.Errorf("%s : %v\n", msg, err)
}
