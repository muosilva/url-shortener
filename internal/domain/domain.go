package domain

import (
	"errors"
	"strings"
)

type CreateURLRequest struct {
	Url string `json:"url"`
}

type JSONResponse struct {
	Msg        string
	StatusCode int
	Code       string
}

//https://redis.io/docs/latest/develop/clients/go/

func (cr *CreateURLRequest) Validation(url string) (err error) {
	if len(url) == 0 {
		err = errors.New("URL cannot be empty!")
		return err
	}

	if !strings.Contains(url, "https") {
		err = errors.New("URL needs to be complete (https)!")
		return err
	}

	return nil
}
