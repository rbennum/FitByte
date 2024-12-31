package helper

import (
	"net/http"
)

func GetErrorStatusCode(err error) int {
	if err.Error() == ErrorNotFound.Error() {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}
