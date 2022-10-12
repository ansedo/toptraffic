package models

import "errors"

var (
	ErrPlacementRequestWrongSchema = errors.New("WRONG_SCHEMA")
	ErrPlacementRequestEmptyTiles  = errors.New("EMPTY_TILES")
	ErrPlacementRequestEmptyField  = errors.New("EMPTY_FIELD")
)
