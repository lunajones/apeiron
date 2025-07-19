package combat

import (
	"log"
	"math"
	"time"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/model" // adicionado para ApplySkillMovement
	"github.com/lunajones/apeiron/lib/navmesh"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
)

func UseSkill(
	attacker model.Attacker,
	target model.Targetable,
	targetPos position.Position,
	skill *model.Skill,
	creatures []model.Targetable,
	players []model.Targetable,
	navMesh *navmesh.NavMesh,
	svcCtx *dynamic_context.AIServiceContext,
) model.SkillResult {
	result := model.SkillResult{}

	if !skill.GroundTargeted && target != nil {
		dir := position.Vector2D{
			X: target.GetPosition().X - attacker.GetPosition().X,
			Z: target.GetPosition().Z - attacker.GetPosition().Z,
		}
		attacker.SetFacingDirection(dir.Normalize())
	}

	if skill.Movement != nil {
		if target != nil {
			attacker.SetSkillMovementState(ApplySkillMovement(attacker, target, skill))
		} else {
			log.Printf("[SkillExecutor] [%s] Skill %s requer target para movimento, mas target é nil", attacker.GetHandle().ID, skill.Name)
		}
	} else if skill.GroundTargeted && skill.AOE != nil {
		ApplyAOEDamage(attacker, targetPos, skill, creatures, players, svcCtx)
		result.Success = true
	} else if skill.Projectile != nil {
		SimulateProjectile(attacker, target, targetPos, skill, svcCtx)
		result.Success = true
	} else {
		result = ApplyDirectDamage(attacker, target, skill, svcCtx)
	}

	return result
}

func ApplySkillMovement(
	attacker model.Attacker,
	target model.Targetable,
	skill *model.Skill,
) *model.SkillMovementState {
	now := time.Now()

	// Calcula direção até o alvo
	dirVec := position.NewVector3DFromTo(attacker.GetPosition(), target.GetPosition()).Normalize()

	// Calcula distância final desejada
	var targetPosAdjusted position.Position
	if skill.Movement.ExtraDistance != 0 {
		targetPosAdjusted = attacker.GetPosition().AddVector3D(dirVec.Scale(skill.Movement.ExtraDistance))
	} else {
		creatureHitbox := attacker.GetHitboxRadius()
		targetHitbox := target.GetHitboxRadius()
		buffer := target.GetDesiredBufferDistance()

		distance := position.CalculateDistance2D(attacker.GetPosition(), target.GetPosition())
		desiredDistance := distance + creatureHitbox + targetHitbox + buffer

		if desiredDistance > skill.Movement.MaxDistance {
			desiredDistance = skill.Movement.MaxDistance
		}

		targetPosAdjusted = attacker.GetPosition().AddVector3D(dirVec.Scale(desiredDistance))
	}

	finalDir := position.NewVector3DFromTo(attacker.GetPosition(), targetPosAdjusted).Normalize()

	log.Printf(
		"[LEAP] [%s] Direção=(%.2f,%.2f,%.2f) AlvoFinal=%v",
		attacker.GetHandle().ID,
		finalDir.X, finalDir.Y, finalDir.Z,
		targetPosAdjusted,
	)

	state := &model.SkillMovementState{
		Active:    true,
		StartTime: now,
		Duration:  time.Duration(skill.Movement.DurationSec * float64(time.Second)),
		Speed:     skill.Movement.Speed,
		Direction: finalDir,
		TargetPos: targetPosAdjusted,
		Config:    skill.Movement,
		Skill:     skill, // Caso precise do skill para o Update
	}

	attacker.SetSkillMovementState(state)

	return state
}

func UpdateSkillMovement(
	mov model.Attacker,
	state *model.SkillMovementState,
	target model.Targetable,
	navMesh *navmesh.NavMesh,
	svcCtx *dynamic_context.AIServiceContext,
	deltaTime float64,
) bool {
	currentPos := mov.GetPosition()
	distanceToTargetPos := position.CalculateDistance2D(currentPos, state.TargetPos)

	moveDist := state.Speed * deltaTime // ~60fps tick
	if moveDist > distanceToTargetPos {
		moveDist = distanceToTargetPos
	}
	moveVec := state.Direction.Scale(moveDist)
	newPos := currentPos.AddVector3D(moveVec)

	mov.SetPosition(newPos)

	realDist := position.CalculateDistance2D(newPos, target.GetPosition())
	log.Printf("[LEAP-REALDIST] Distância real após avanço: %.2f", realDist)

	if !state.DamageApplied && realDist <= state.Config.MaxDistance {

		ApplyDirectDamage(mov, target, state.Skill, svcCtx)
		state.DamageApplied = true
	}

	return state.IsComplete(time.Now(), newPos)
}

func ApplyDirectDamage(attacker model.Attacker, target model.Targetable, skill *model.Skill, svcCtx *dynamic_context.AIServiceContext) model.SkillResult {
	result := model.SkillResult{}

	if target == nil || shouldSkipTarget(attacker, target) {
		attacker.SetLastMissedSkillAt(time.Now())
		return result
	}

	if target.IsInvulnerableNow() {
		log.Printf("[SkillExecutor] [%s] invulnerável no momento, dano evitado de [%s]", target.GetHandle().ID, attacker.GetHandle().ID)
		attacker.SetLastMissedSkillAt(time.Now())

		return result
	}

	// Atualiza direção do atacante
	dir := position.NewVector2DFromTo(attacker.GetPosition(), target.GetPosition())
	attacker.SetFacingDirection(dir)

	damage := calculateDamageGeneric(attacker, target, skill.InitialMultiplier)

	// 🔄 BLOQUEIO E PARRY
	if target.IsBlocking() {
		blockDir := target.GetFacingDirection()
		attackDir := position.NewVector2DFromTo(target.GetPosition(), attacker.GetPosition())
		dot := blockDir.Dot(attackDir)

		if dot > 0.5 {
			// PARRY
			if target.IsInParryWindow() {
				log.Printf("[PARRY] [%s] executou parry em [%s]", target.GetHandle().ID, attacker.GetHandle().ID)
				attacker.SetLastMissedSkillAt(time.Now())
				return result // Parry bem-sucedido cancela ataque
			}

			// BLOQUEIO bem-sucedido
			log.Printf("[BLOCK] [%s] bloqueou ataque de [%s]", target.GetHandle().ID, attacker.GetHandle().ID)

			if skill.Impact != nil && skill.Impact.PostureDamage > 0 {
				postureDamage := skill.Impact.PostureDamage + getPostureScaling(attacker, skill)
				target.ApplyPostureDamage(postureDamage * 2)
				log.Printf("[BLOCK] [%s] aplicou %.1f de posture damage (dobrado)", target.GetHandle().ID, postureDamage*2)
			}

			attacker.SetLastMissedSkillAt(time.Now())

			return result // Bloqueio nega dano
		} else {
			log.Printf("[BLOCK-FAILED] [%s] bloqueou em direção errada, ataque passou", target.GetHandle().ID)
		}
	}

	// Aplica dano
	target.TakeDamage(damage)
	result.TargetDied = !target.IsAlive()

	if skill.Impact != nil && skill.Impact.PostureDamage > 0 {
		postureDamage := skill.Impact.PostureDamage + getPostureScaling(attacker, skill)
		target.ApplyPostureDamage(postureDamage)
	}

	if skill.HasDOT && skill.DOT != nil {
		dotPower := damage / (skill.DOT.DurationSec / skill.DOT.TickSec)
		effect := constslib.ActiveEffect{
			Type:            skill.DOT.EffectType,
			StartTime:       time.Now(),
			Duration:        time.Duration(skill.DOT.DurationSec) * time.Second,
			TickInterval:    time.Duration(skill.DOT.TickSec) * time.Second,
			Power:           dotPower,
			IsDOT:           true,
			IsDebuff:        true,
			Elapsed:         0,
			LastTickElapsed: 0,
		}
		target.ApplyEffect(effect)
	}

	result.Success = true
	return result
}

func ApplyAOEDamage(attacker model.Attacker, targetPos position.Position, skill *model.Skill, creatures []model.Targetable, players []model.Targetable, svcCtx *dynamic_context.AIServiceContext) {
	for _, t := range creatures {
		if t.GetHandle().Equals(attacker.GetHandle()) {
			continue
		}
		if shouldSkipTarget(attacker, t) {
			continue
		}
		if position.CalculateDistance(t.GetPosition(), targetPos) <= skill.AOE.Radius {
			ApplyDirectDamage(attacker, t, skill, svcCtx)
		}
	}

	for _, t := range players {
		if t.GetHandle().Equals(attacker.GetHandle()) {
			continue
		}
		if shouldSkipTarget(attacker, t) {
			continue
		}
		if position.CalculateDistance(t.GetPosition(), targetPos) <= skill.AOE.Radius {
			ApplyDirectDamage(attacker, t, skill, svcCtx)
		}
	}
}

func SimulateProjectile(attacker model.Attacker, target model.Targetable, targetPos position.Position, skill *model.Skill, svcCtx *dynamic_context.AIServiceContext) {
	if target == nil || skill.Projectile == nil {
		log.Printf("[SkillExecutor] SimulateProjectile inválido. Skill: %s", skill.Name)
		return
	}

	travelTime := position.CalculateDistance(attacker.GetPosition(), targetPos) / skill.Projectile.Speed
	time.AfterFunc(time.Duration(travelTime*1000)*time.Millisecond, func() {
		ApplyDirectDamage(attacker, target, skill, svcCtx)
		log.Printf("[SkillExecutor] Projetil %s chegou ao alvo %s após %.2f segundos",
			skill.Name, target.GetHandle().ID, travelTime)
	})
}

func IsBehind(attacker model.Attacker, target model.Targetable) bool {
	dirToAttacker := position.Vector2D{
		X: attacker.GetPosition().X - target.GetPosition().X,
		Z: attacker.GetPosition().Z - target.GetPosition().Z,
	}.Normalize()

	targetFacing := target.GetFacingDirection().Normalize()
	dot := dirToAttacker.Dot(targetFacing)

	return dot > 0.5
}

func shouldSkipTarget(attacker model.Attacker, target model.Targetable) bool {
	if attacker == nil || target == nil {
		return true
	}

	if attacker.GetHandle().Equals(target.GetHandle()) {
		return true // Não atacar a si mesmo
	}

	if !target.IsAlive() {
		return true
	}

	if target.IsHostile() {
		return false // Pode atacar se é hostil
	}

	// Caso adicional: se atacante precisa de hostilidade ou fome (ex: criatura faminta)
	// Você pode adaptar essa regra conforme o contexto do seu AI
	if attacker.IsHungry() {
		return false // Pode atacar mesmo não sendo hostil se está com fome
	}

	// PvP check (caso alvo seja player)
	if !target.IsPvPEnabled() {
		return true
	}

	return false
}

func getPostureScaling(attacker model.Attacker, skill *model.Skill) float64 {
	if skill.Impact == nil {
		return 0
	}
	switch skill.Impact.ScalingStat {
	case "Strength":
		return float64(attacker.GetStrength()) * skill.Impact.ScalingMultiplier
	case "Dexterity":
		return float64(attacker.GetDexterity()) * skill.Impact.ScalingMultiplier
	case "Intelligence":
		return float64(attacker.GetIntelligence()) * skill.Impact.ScalingMultiplier
	case "Focus":
		return float64(attacker.GetFocus()) * skill.Impact.ScalingMultiplier
	}
	return 0
}

func calculateDamageGeneric(attacker model.Attacker, target model.Targetable, mult float64) int {
	baseDamage := 0
	strength := float64(attacker.GetStrength())

	switch t := target.(type) {
	case model.Attacker:
		// Se o alvo também for um Attacker, podemos acessar defesas via interface (se tiver)
		// Vamos assumir que temos GetPhysicalDefense() na interface ou tratamos como 0
		var defense float64
		if pd, ok := t.(interface{ GetPhysicalDefense() float64 }); ok {
			defense = pd.GetPhysicalDefense()
		}
		baseDamage = int(strength*mult - (defense * strength))

	default:
		baseDamage = int(strength * mult)
	}

	if baseDamage <= 0 {
		baseDamage = 1
	}

	return baseDamage
}

func CalculatePhysicalDamage(attacker model.Attacker, target model.Targetable, skillMultiplier float64) int {
	baseAttack := float64(attacker.GetStrength())
	dexBonus := float64(attacker.GetDexterity()) * 0.1
	rawDamage := (baseAttack + dexBonus) * skillMultiplier

	var defense float64
	if d, ok := target.(interface{ GetPhysicalDefense() float64 }); ok {
		defense = d.GetPhysicalDefense()
	}

	finalDamage := rawDamage * (1 - defense)
	if finalDamage < 1 {
		finalDamage = 1
	}

	return int(math.Round(finalDamage))
}

func CalculateMagicDamage(attacker model.Attacker, target model.Targetable, skillMultiplier float64) int {
	baseMagic := float64(attacker.GetIntelligence())
	focusBonus := float64(attacker.GetFocus()) * 0.05
	rawDamage := (baseMagic + focusBonus) * skillMultiplier

	var defense float64
	if d, ok := target.(interface{ GetMagicDefense() float64 }); ok {
		defense = d.GetMagicDefense()
	}

	finalDamage := rawDamage * (1 - defense)
	if finalDamage < 1 {
		finalDamage = 1
	}

	return int(math.Round(finalDamage))
}

func CalculatePoisonDamage(attacker model.Attacker, target model.Targetable) int {
	base := 5.0
	strBonus := float64(attacker.GetStrength()) * 0.2
	intBonus := float64(attacker.GetIntelligence()) * 0.1
	rawDOT := base + strBonus + intBonus

	var resist float64
	if r, ok := target.(interface{ GetStatusResistance() float64 }); ok {
		resist = r.GetStatusResistance()
	}

	finalDOT := rawDOT * (1 - resist)
	if finalDOT < 1 {
		finalDOT = 1
	}

	return int(math.Round(finalDOT))
}

func CalculateBurnDamage(attacker model.Attacker, target model.Targetable) int {
	base := 7.0
	intBonus := float64(attacker.GetIntelligence()) * 0.3
	rawDOT := base + intBonus

	var resist float64
	if r, ok := target.(interface{ GetStatusResistance() float64 }); ok {
		resist = r.GetStatusResistance()
	}

	finalDOT := rawDOT * (1 - resist)
	if finalDOT < 1 {
		finalDOT = 1
	}

	return int(math.Round(finalDOT))
}

func CalculateHealing(attacker model.Attacker, skillMultiplier float64) int {
	baseHeal := float64(attacker.GetFocus()*2 + attacker.GetIntelligence())
	finalHeal := baseHeal * skillMultiplier

	if finalHeal < 1 {
		finalHeal = 1
	}

	return int(math.Round(finalHeal))
}

func CalculateEffectiveCCDuration(baseDuration float64, target model.Targetable) float64 {
	var resist float64
	if r, ok := target.(interface{ GetControlResistance() float64 }); ok {
		resist = r.GetControlResistance()
	}

	reduction := baseDuration * resist
	finalDuration := baseDuration - reduction

	if finalDuration < 0.1 {
		finalDuration = 0.1 // Mínimo de 0.1s pra evitar CC zero
	}

	return finalDuration
}
