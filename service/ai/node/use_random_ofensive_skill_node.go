package node

import (
	"log"
	"math/rand"
	"time"

	"github.com/lunajones/apeiron/lib/combat"
	"github.com/lunajones/apeiron/lib/physics"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/consts"
)

type UseRandomOffensiveSkillNode struct{}

func (n *UseRandomOffensiveSkillNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	log.Printf("[RANDOM SKILL] Tick iniciado para %s (%s)", c.GetHandle().ID, c.PrimaryType)

	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[RANDOM SKILL] [%s (%s)] contexto inválido", c.GetHandle().ID, c.PrimaryType)
		return core.StatusFailure
	}

	grid := svcCtx.GetPathfindingGrid()
	if grid == nil {
		log.Printf("[RANDOM SKILL] [%s (%s)] grid indisponível", c.GetHandle().ID, c.PrimaryType)
		return core.StatusFailure
	}

	nearbyCreatures := svcCtx.GetServiceCreatures(c.GetPosition(), c.DetectionRadius)
	nearbyPlayers := svcCtx.GetServicePlayers(c.GetPosition(), c.DetectionRadius)

	target := c.GetBestTarget(nearbyCreatures, nearbyPlayers)
	if target == nil {
		log.Printf("[RANDOM SKILL] [%s (%s)] nenhum alvo válido encontrado, mudando para SearchFood", c.GetHandle().ID, c.PrimaryType)
		c.ClearTargetHandles()
		c.ChangeAIState(consts.AIStateSearchFood)
		return core.StatusFailure
	}

	if target.IsCreature() {
		if cTarget, ok := target.(*creature.Creature); ok && !cTarget.IsAlive {
			log.Printf("[RANDOM SKILL] [%s (%s)] alvo morto: %s", c.GetHandle().ID, c.PrimaryType, cTarget.GetHandle().ID)
			c.ClearTargetHandles()
			c.ChangeAIState(consts.AIStateSearchFood)
			return core.StatusFailure
		}
	}

	validSkills := []combat.Skill{}
	now := time.Now()

	for _, skillName := range c.Skills {
		skill, ok := combat.SkillRegistry[skillName]
		if !ok {
			continue
		}
		if skill.SkillType != "Physical" && skill.SkillType != "Magic" {
			continue
		}
		if lastUsed, exists := c.SkillCooldowns[skill.Action]; exists {
			if now.Sub(lastUsed).Seconds() < float64(skill.CooldownSec) {
				continue
			}
		}
		validSkills = append(validSkills, skill)
	}

	if len(validSkills) == 0 {
		log.Printf("[RANDOM SKILL] [%s (%s)] sem skills disponíveis", c.GetHandle().ID, c.PrimaryType)
		return core.StatusRunning
	}

	chosenSkill := validSkills[rand.Intn(len(validSkills))]
	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())
	stopAt := c.GetHitboxRadius() + target.GetHitboxRadius() + c.GetDesiredBufferDistance() + 0.2

	if dist > chosenSkill.Range || dist > stopAt {
		if !c.MoveCtrl.IsMoving {
			if physics.IsWalkable(target.GetPosition(), c.GetHitboxRadius()) {
				c.MoveCtrl.SetTarget(target.GetPosition(), c.GetCurrentSpeed(), stopAt)
			} else {
				log.Printf("[RANDOM SKILL] [%s (%s)] destino até alvo bloqueado, abortando aproximação", c.GetHandle().ID, c.PrimaryType)
				return core.StatusFailure
			}
		}
		c.MoveCtrl.Update(c, 0.016, grid)
		c.SetAction(consts.ActionRun)
		log.Printf("[RANDOM SKILL] [%s (%s)] aproximando do alvo (dist %.2f > stopAt %.2f, range %.2f)",
			c.GetHandle().ID, c.PrimaryType, dist, stopAt, chosenSkill.Range)
		return core.StatusRunning
	}

	log.Printf("[RANDOM SKILL] [%s (%s)] usando skill %s em %s",
		c.GetHandle().ID, c.PrimaryType, chosenSkill.Name, target.GetHandle().ID)

	result := combat.UseSkill(c, target, target.GetPosition(), chosenSkill, nearbyCreatures, nearbyPlayers)

	if result.TargetDied && c.IsHungry() {
		log.Printf("[RANDOM SKILL] [%s (%s)] matou %s e está com fome, trocando para SearchFood",
			c.GetHandle().ID, c.PrimaryType, target.GetHandle().ID)
		c.ChangeAIState(consts.AIStateSearchFood)
		return core.StatusSuccess
	}

	if result.Success {
		log.Printf("[RANDOM SKILL] [%s (%s)] skill %s usada com sucesso", c.GetHandle().ID, c.PrimaryType, chosenSkill.Name)
		return core.StatusSuccess
	}

	log.Printf("[RANDOM SKILL] [%s (%s)] skill %s falhou", c.GetHandle().ID, c.PrimaryType, chosenSkill.Name)
	return core.StatusFailure
}

func (n *UseRandomOffensiveSkillNode) Reset() {
	// Nada a resetar
}
