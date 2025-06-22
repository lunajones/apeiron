package combat

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/service/creature"
)

func UseSkill(attacker *creature.Creature, target *creature.Creature, targetPos creature.Position, skillName string, creatures []*creature.Creature, players []Player) {
	skill, exists := SkillRegistry[skillName]
	if !exists {
		log.Printf("[SkillExecutor] Skill %s não encontrada", skillName)
		return
	}

	// Cooldown check
	lastUsed, onCooldown := attacker.SkillCooldowns[skill.Action]
	if onCooldown && time.Since(lastUsed).Seconds() < float64(skill.CooldownSec) {
		log.Printf("[SkillExecutor] Skill %s em cooldown para %s", skillName, attacker.ID)
		return
	}

	// Aplicação
	if skill.IsGroundTargeted && skill.AOERadius > 0 {
		ApplyAOEDamage(attacker, targetPos, skill, creatures, players)
	} else if skill.HasProjectile {
		SimulateProjectile(attacker, target, targetPos, skill)
	} else {
		ApplyDirectDamage(attacker, target, skill)
	}

	attacker.SkillCooldowns[skill.Action] = time.Now()
}

func ApplyDirectDamage(attacker *creature.Creature, target *creature.Creature, skill Skill) {
	var damage int
	switch skill.SkillType {
	case "Physical":
		damage = CalculatePhysicalDamage(attacker, target, skill.Multiplier)
	case "Magic":
		damage = CalculateMagicDamage(attacker, target, skill.Multiplier)
	default:
		log.Printf("[SkillExecutor] Tipo %s não implementado", skill.SkillType)
		return
	}

	target.HP -= damage
	log.Printf("[SkillExecutor] %s usou %s contra %s causando %d de dano. HP alvo: %d", attacker.ID, skill.Name, target.ID, damage, target.HP)

	if skill.IsDOT {
		dotPower := damage / (skill.DOTDurationSec / skill.DOTTickSec)
		effect := creature.ActiveEffect{
			Type:         creature.EffectPoison,
			StartTime:    time.Now().Unix(),
			Duration:     int64(skill.DOTDurationSec),
			TickInterval: int64(skill.DOTTickSec),
			Power:        dotPower,
			IsDOT:        true,
			IsDebuff:     true,
		}
		target.ApplyEffect(effect)
	}
}

func ApplyAOEDamage(attacker *creature.Creature, targetPos creature.Position, skill Skill, creatures []*creature.Creature, players []Player) {
	for _, c := range creatures {
		if !c.IsAlive {
			continue
		}
		if Distance(c.Position, targetPos) <= skill.AOERadius {
			ApplyDirectDamage(attacker, c, skill)
		}
	}
	// Players também, se quiser permitir friendly fire ou PvP
	for _, p := range players {
		if Distance(p.Position, targetPos) <= skill.AOERadius {
			// Aqui você pode aplicar o mesmo cálculo para players, criando um ApplyDirectDamageToPlayer()
		}
	}
}

func SimulateProjectile(attacker *creature.Creature, target *creature.Creature, targetPos creature.Position, skill Skill) {
	travelTime := Distance(attacker.Position, targetPos) / skill.ProjectileSpeed
	time.AfterFunc(time.Duration(travelTime*1000)*time.Millisecond, func() {
		ApplyDirectDamage(attacker, target, skill)
		log.Printf("[SkillExecutor] Projetil %s chegou ao alvo %s após %.2f segundos", skill.Name, target.ID, travelTime)
	})
}

func Distance(a, b creature.Position) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
