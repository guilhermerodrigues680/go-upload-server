package handler

import "errors"

var (
	ErrFileIsNotRegular = errors.New("File is not regular")
	ErrFileIsNotDir     = errors.New("File is not dir")
	ErrConnClosed       = errors.New("Connection Closed")
	ErrUnexpected       = errors.New("Erro inesperado")
)
