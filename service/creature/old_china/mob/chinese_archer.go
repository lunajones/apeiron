package mob

// import (
// 	"time"

// 	"github.com/lunajones/apeiron/lib"
// 	"github.com/lunajones/apeiron/lib/handle"
// 	"github.com/lunajones/apeiron/lib/model"
// 	"github.com/lunajones/apeiron/lib/position"
// 	"github.com/lunajones/apeiron/service/creature"
// 	"github.com/lunajones/apeiron/service/creature/consts"
// )

// func NewChineseArcher() *creature.Creature {
// 	id := lib.NewUUID()

// 	return &creature.Creature{
// 		Handle:     handle.NewEntityHandle(id, 1),
// 		Generation: 1,

// 		Creature: model.Creature{
// 			Name:           "Chinese Archer",
// 			MaxHP:          250,
// 			RespawnTimeSec: 60,
// 			SpawnPoint:     position.Position{},
// 			SpawnRadius:    5.0,
// 		},

// 		Actions: []consts.CreatureAction{
// 			consts.ActionIdle,
// 			consts.ActionWalk,
// 			consts.ActionSkill1,
// 			consts.ActionSkill2,
// 			consts.ActionCombo1,
// 			consts.ActionDie,
// 		},
// 		HP:              250,
// 		LastStateChange: time.Now(),
// 		IsAlive:         true,
// 		Position:        position.Position{},
// 		WalkSpeed:       2.5,
// 		RunSpeed:        4.0,
// 	}
// }
