package combat

import (
	"log"
	"math"
	"time"

	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/player"
)

func UseSkill(attacker *creature.Creature, target *creature.Creature, targetPos position.Position, skillName string, creatures []*creature.Creature, players []*player.Player) {
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

	// Atualizar FacingDirection para skills single-target
	if !skillData.IsGroundTargeted && target != nil {
		dir := position.Vector2D{
			X: target.Position.X - attacker.Position.X,
			Y: target.Position.Z - attacker.Position.Z,
		}
		attacker.FacingDirection = dir.Normalize()
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

func ApplyDirectDamage(attacker *creature.Creature, target *creature.Creature, skill Skill) {
	if target != nil {
		dir := position.Vector2D{
			X: target.Position.X - attacker.Position.X,
			Y: target.Position.Z - attacker.Position.Z,
		}
		attacker.FacingDirection = dir.Normalize()
	}

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

	if target.HP <= 0 {
		target.IsAlive = false
		target.IsCorpse = true
		target.TimeOfDeath = time.Now()
		log.Printf("[Combat] %s morreu para %s.", target.ID, attacker.ID)
	}

	if skill.IsDOT {
		dotPower := damage / (skill.DOTDurationSec / skill.DOTTickSec)
		effect := creature.ActiveEffect{
			Type:         creature.EffectPoison,
			StartTime:    time.Now(),
			Duration:     time.Duration(skill.DOTDurationSec) * time.Second,
			TickInterval: time.Duration(skill.DOTTickSec) * time.Second,
			Power:        dotPower,
			IsDOT:        true,
			IsDebuff:     true,
		}
		target.ApplyEffect(effect)
	}
}


func ApplyAOEDamage(attacker *creature.Creature, targetPos position.Position, skillData Skill, creatures []*creature.Creature, players []*player.Player) {
	for _, c := range creatures {
		if !c.IsAlive {
			continue
		}
		if Distance(c.Position, targetPos) <= skillData.AOERadius {
			ApplyDirectDamage(attacker, c, skillData)
		}
	}

	// Exemplo: aplicar também em players, se quiser fazer isso depois
	for _, p := range players {
		if Distance(p.Position, targetPos) <= skillData.AOERadius {
			// Exemplo: criar ApplyDirectDamageToPlayer() se quiser dano em player
		}
	}
}

func SimulateProjectile(attacker *creature.Creature, target *creature.Creature, targetPos position.Position, skillData Skill) {
	if target == nil {
		log.Printf("[SkillExecutor] SimulateProjectile recebeu target nil. Skill: %s", skillData.Name)
		return
	}

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
