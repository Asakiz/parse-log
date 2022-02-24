package utils

import (
	"errors"
)

var (
	ErrFileExtension    = errors.New("invalid format of input file")
	ErrFilePathNotExist = errors.New("file path doesn't exists")
	ErrMissingArguments = errors.New("missing parameter, provide file path")
)
