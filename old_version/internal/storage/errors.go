package storage

import "errors"

var (
	ErrConnectionBucketNotFound = errors.New("Cannot find connections bucket")
	ErrConnectionNotFound       = errors.New("Cannot find connection in bucket")
	ErrConnectionAlreadyExists  = errors.New("Connection with this label already exists in bucket")
	ErrConnectionToBytes        = errors.New("Failed to convert connection to bytes")
	ErrFailedDelete             = errors.New("Cannot delete connection from bucket")
	ErrUpdateConnection         = errors.New("Failed update connection")
)
