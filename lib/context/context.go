package context

type CreatureContext interface {
	GetCreatures() []interface{}
	GetPlayers() []interface{}
}