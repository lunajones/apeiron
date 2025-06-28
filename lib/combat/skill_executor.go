package combat

import (
	"log"
	"time"

	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/creature/consts"
	"github.com/lunajones/apeiron/service/player"
)

func UseSkill(attacker *creature.Creature, target model.Targetable, targetPos position.Position, skill Skill, creatures []*creature.Creature, players []*player.Player) SkillResult {
	result := SkillResult{}

	lastUsed, onCooldown := attacker.SkillCooldowns[skill.Action]
	if onCooldown && time.Since(lastUsed).Seconds() < float64(skill.CooldownSec) {
		log.Printf("[SkillExecutor] [%s] skill %s em cooldown", attacker.GetHandle().ID, skill.Name)
		result.WasOnCooldown = true
		return result
	}

	if !skill.GroundTargeted && target != nil {
		dir := position.Vector2D{
			X: target.GetPosition().FastGlobalX() - attacker.GetPosition().FastGlobalX(),
			Y: target.GetPosition().Z - attacker.GetPosition().Z,
		}
		attacker.FacingDirection = dir.Normalize()
	}

	if skill.GroundTargeted && skill.AOE != nil {
		ApplyAOEDamage(attacker, targetPos, skill, creatures, players)
	} else if skill.Projectile != nil {
		SimulateProjectile(attacker, target, targetPos, skill)
	} else {
		result = ApplyDirectDamage(attacker, target, skill)
	}

	attacker.SkillCooldowns[skill.Action] = time.Now()
	result.Success = true
	return result
}

func ApplyDirectDamage(attacker *creature.Creature, target model.Targetable, skill Skill) SkillResult {
	result := SkillResult{}

	if target == nil || shouldSkipTarget(attacker, target) {
		return result
	}

	// Checa invulnerabilidade se for criatura
	if tgt, ok := target.(*creature.Creature); ok && tgt.Invincibility.IsInvincible {
		log.Printf("[SkillExecutor] [%s] tentou atacar [%s], mas alvo está invulnerável", attacker.GetHandle().ID, tgt.GetHandle().ID)
		return result
	}

	// Checa requisitos de posicionamento
	if skill.Conditions != nil && skill.Conditions.FacingRequirement == "Behind" {
		if !IsBehind(attacker, target) {
			log.Printf("[SkillExecutor] [%s] precisa estar atrás de [%s] para usar %s",
				attacker.GetHandle().ID, target.GetHandle().ID, skill.Name)
			return result
		}
	}

	dir := position.Vector2D{
		X: target.GetPosition().FastGlobalX() - attacker.GetPosition().FastGlobalX(),
		Y: target.GetPosition().Z - attacker.GetPosition().Z,
	}
	attacker.FacingDirection = dir.Normalize()

	// Calcula dano
	damage := calculateDamageGeneric(attacker, target, skill.InitialMultiplier)

	target.TakeDamage(damage)
	result.TargetDied = !target.CheckIsAlive()

	// Postura
	if skill.Impact != nil && skill.Impact.PostureDamageBase > 0 {
		postureDamage := skill.Impact.PostureDamageBase + getPostureScaling(attacker, skill)
		if tgt, ok := target.(*creature.Creature); ok {
			tgt.ApplyPostureDamage(postureDamage)
		}
	}

	// DOT
	if skill.HasDOT && skill.DOT != nil {
		dotPower := damage / (skill.DOT.DurationSec / skill.DOT.TickSec)
		effect := consts.ActiveEffect{
			Type:         skill.DOT.EffectType,
			StartTime:    time.Now(),
			Duration:     time.Duration(skill.DOT.DurationSec) * time.Second,
			TickInterval: time.Duration(skill.DOT.TickSec) * time.Second,
			Power:        dotPower,
			IsDOT:        true,
			IsDebuff:     true,
		}
		target.ApplyEffect(effect)
	}

	log.Printf("[SkillExecutor] [%s (%s)] usou %s em [%s] causando %d de dano",
		attacker.GetHandle().ID, attacker.PrimaryType, skill.Name, target.GetHandle().ID, damage)

	if result.TargetDied && attacker.IsHungry() {
		log.Printf("[AI] [%s] matou [%s], está com fome e vai buscar comida", attacker.GetHandle().ID, target.GetHandle().ID)
		attacker.ChangeAIState(consts.AIStateSearchFood)
	}

	result.Success = true
	return result
}

func calculateDamageGeneric(attacker *creature.Creature, target model.Targetable, mult float64) int {
	baseDamage := 0
	switch skillTarget := target.(type) {
	case *creature.Creature:
		baseDamage = int(float64(attacker.Strength)*mult - (skillTarget.PhysicalDefense * float64(attacker.Strength)))
	case *player.Player:
		baseDamage = int(float64(attacker.Strength) * mult) // Ajuste para player
	default:
		baseDamage = int(float64(attacker.Strength) * mult)
	}
	if baseDamage <= 0 {
		baseDamage = 1
	}
	return baseDamage
}

func ApplyAOEDamage(attacker *creature.Creature, targetPos position.Position, skill Skill, creatures []*creature.Creature, players []*player.Player) {
	for _, c := range creatures {
		if c.GetHandle().Equals(attacker.GetHandle()) {
			continue
		}
		if shouldSkipTarget(attacker, c) {
			continue
		}
		if position.CalculateDistance(c.GetPosition(), targetPos) <= skill.AOE.Radius {
			ApplyDirectDamage(attacker, c, skill)
		}
	}

	for _, p := range players {
		if position.CalculateDistance(p.Position, targetPos) <= skill.AOE.Radius {
			// No futuro: ApplyDirectDamageToPlayer
		}
	}
}

func SimulateProjectile(attacker *creature.Creature, target model.Targetable, targetPos position.Position, skill Skill) {
	if target == nil || skill.Projectile == nil {
		log.Printf("[SkillExecutor] SimulateProjectile inválido. Skill: %s", skill.Name)
		return
	}

	travelTime := position.CalculateDistance(attacker.GetPosition(), targetPos) / skill.Projectile.Speed
	time.AfterFunc(time.Duration(travelTime*1000)*time.Millisecond, func() {
		ApplyDirectDamage(attacker, target, skill)
		log.Printf("[SkillExecutor] Projetil %s chegou ao alvo %s após %.2f segundos", skill.Name, target.GetHandle().ID, travelTime)
	})
}

func IsBehind(attacker *creature.Creature, target model.Targetable) bool {
	dirToAttacker := position.Vector2D{
		X: attacker.GetPosition().FastGlobalX() - target.GetPosition().FastGlobalX(),
		Y: attacker.GetPosition().Z - target.GetPosition().Z,
	}.Normalize()
	targetFacing := target.GetFacingDirection().Normalize()
	dot := dirToAttacker.X*targetFacing.X + dirToAttacker.Y*targetFacing.Y
	return dot > 0.5
}

func shouldSkipTarget(attacker, target model.Targetable) bool {
	if attacker == nil || target == nil {
		return true
	}
	if attacker.GetHandle().Equals(target.GetHandle()) {
		return true
	}

	switch tgt := target.(type) {
	case *creature.Creature:
		creatureAttacker, ok := attacker.(*creature.Creature)
		if !ok {
			return !tgt.IsAlive || !tgt.IsHostile
		}
		return !tgt.IsAlive || (!tgt.IsHostile && !creatureAttacker.IsHungry())
	case *player.Player:
		return !tgt.IsAlive || !tgt.IsPvPEnabled
	}

	return false
}

func getPostureScaling(attacker *creature.Creature, skill Skill) float64 {
	if skill.Impact == nil {
		return 0
	}
	switch skill.Impact.ScalingStat {
	case "Strength":
		return float64(attacker.Strength) * skill.Impact.ScalingMultiplier
	case "Dexterity":
		return float64(attacker.Dexterity) * skill.Impact.ScalingMultiplier
	case "Intelligence":
		return float64(attacker.Intelligence) * skill.Impact.ScalingMultiplier
	case "Focus":
		return float64(attacker.Focus) * skill.Impact.ScalingMultiplier
	}
	return 0
}
