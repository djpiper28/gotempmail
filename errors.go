package gotempmail

import (
	"fmt"
)

func BodyReadErr(err error) error {
	return fmt.Errorf("CANNOT READ BODY %s", err)
}

func StatusCodeErr(code int) error {
	return fmt.Errorf("UNEXPECTED RETURN CODE (%d)", code)
}

func JsonParseErr(err error) error {
	return fmt.Errorf("CANNOT PARSE JSON %s", err)
}
