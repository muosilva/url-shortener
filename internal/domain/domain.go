package domain

import (
	"errors"
	"strings"
)

type CreateURLRequest struct {
	Url string `json:"url"`
}

type JSONResponse struct {
	Msg        string `json:"message"`
	StatusCode int    `json:"status_code"`
	Url        string `json:"url,omitempty"`
}

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
