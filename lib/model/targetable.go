package model

import "github.com/lunajones/apeiron/lib/position"

type Targetable interface {
	GetPosition() position.Position
	GetID() string
}