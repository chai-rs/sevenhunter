package errx

import (
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func Mongo(err error) *Error {
	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return E(http.StatusNotFound, err, "resource not found")
	case mongo.IsDuplicateKeyError(err):
		return E(http.StatusConflict, err, "resource already exists")
	}

	return InternalServerError
}
