package mob

import (
	"time"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/node"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

func BuildChineseSoldierBT(players []*player.Player, creatures []*creature.Creature) core.BehaviorNode {
	return &core.SelectorNode{
		Children: []core.BehaviorNode{
			core.NewCooldownDecorator(&node.FleeIfLowHPNode{}, 1*time.Second),
			core.NewCooldownDecorator(&node.DetectOtherCreatureNode{}, 2*time.Second),
			&node.AttackTargetNode{SkillName: "SoldierSkill1"},
			&node.UseGroundSkillNode{SkillName: "SoldierGroundSlam"},
			&node.AttackIfVulnerableNode{},
			&node.RandomIdleNode{},
		},
	}
}
