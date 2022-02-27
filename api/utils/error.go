package utils

import (
	"errors"
)

var (
	ErrFileExtension    = errors.New("invalid format of input file")
	ErrFilePathNotExist = errors.New("file path doesn't exists")
	ErrMissingArguments = errors.New("missing parameter, provide file path")
	ErrToOpenFile       = errors.New("failed to open the file, maybe doesn't exists or wrong name")
	ErrWriteToFile      = errors.New("failed to write on the file")
)
