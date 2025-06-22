package combat

import (
	"log"
	"math"
	"time"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/player"
	"github.com/lunajones/apeiron/service/creature"
)

func UseSkill(attacker *creature.Creature, target *creature.Creature, targetPos position.Position, skillName string, creatures []*creature.Creature, players []player.Player) {
	skillData, exists := SkillRegistry[skillName]
	if !exists {
		log.Printf("[SkillExecutor] Skill %s não encontrada", skillName)
		return
	}

	// Cooldown check
	lastUsed, onCooldown := attacker.SkillCooldowns[skillData.Action]
	if onCooldown && time.Since(lastUsed).Seconds() < float64(skillData.CooldownSec) {
		log.Printf("[SkillExecutor] Skill %s em cooldown para %s", skillName, attacker.ID)
		return
	}

	// Aplicação
	if skillData.IsGroundTargeted && skillData.AOERadius > 0 {
		ApplyAOEDamage(attacker, targetPos, skillData, creatures, players)
	} else if skillData.HasProjectile {
		SimulateProjectile(attacker, target, targetPos, skillData)
	} else {
		ApplyDirectDamage(attacker, target, skillData)
	}

	attacker.SkillCooldowns[skillData.Action] = time.Now()
}

func ApplyDirectDamage(attacker *creature.Creature, target *creature.Creature, skillData Skill) {
	var damage int
	switch skillData.SkillType {
	case "Physical":
		damage = CalculatePhysicalDamage(attacker, target, skillData.Multiplier)
	case "Magic":
		damage = CalculateMagicDamage(attacker, target, skillData.Multiplier)
	default:
		log.Printf("[SkillExecutor] Tipo %s não implementado", skillData.SkillType)
		return
	}

	target.HP -= damage
	log.Printf("[SkillExecutor] %s usou %s contra %s causando %d de dano. HP alvo: %d", attacker.ID, skillData.Name, target.ID, damage, target.HP)

	if skillData.IsDOT {
		dotPower := damage / (skillData.DOTDurationSec / skillData.DOTTickSec)
		effect := creature.ActiveEffect{
			Type:         creature.EffectPoison,
			StartTime:    time.Now().Unix(),
			Duration:     int64(skillData.DOTDurationSec),
			TickInterval: int64(skillData.DOTTickSec),
			Power:        dotPower,
			IsDOT:        true,
			IsDebuff:     true,
		}
		target.ApplyEffect(effect)
	}
}

func ApplyAOEDamage(attacker *creature.Creature, targetPos position.Position, skillData Skill, creatures []*creature.Creature, players []player.Player) {
	for _, c := range creatures {
		if !c.IsAlive {
			continue
		}
		if Distance(c.Position, targetPos) <= skillData.AOERadius {
			ApplyDirectDamage(attacker, c, skillData)
		}
	}
	// Players também, se quiser permitir friendly fire ou PvP
	for _, p := range players {
		if Distance(p.Position, targetPos) <= skillData.AOERadius {
			// Aqui você pode aplicar o mesmo cálculo para players, criando um ApplyDirectDamageToPlayer()
		}
	}
}

func SimulateProjectile(attacker *creature.Creature, target *creature.Creature, targetPos position.Position, skillData Skill) {
	travelTime := Distance(attacker.Position, targetPos) / skillData.ProjectileSpeed
	time.AfterFunc(time.Duration(travelTime*1000)*time.Millisecond, func() {
		ApplyDirectDamage(attacker, target, skillData)
		log.Printf("[SkillExecutor] Projetil %s chegou ao alvo %s após %.2f segundos", skillData.Name, target.ID, travelTime)
	})
}

func Distance(a, b position.Position) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
