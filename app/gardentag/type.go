package gardentag

import (
	"garden/hashtree"
	"time"
)

type GardenTag struct {
	Signature string
	Message   string
	Timestamp time.Time
	tree      *hashtree.HashTree
}

func NewGardenTag() *GardenTag {
	return &GardenTag{}
}
