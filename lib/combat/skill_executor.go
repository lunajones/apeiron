package combat

import (
	"log"
	"math"
	"time"

	"github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/model" // adicionado para ApplySkillMovement
	"github.com/lunajones/apeiron/lib/navmesh"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/helper/finder"
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

	// 🔁 Direciona para o alvo se necessário
	if !skill.GroundTargeted && target != nil {
		dir := position.Vector2D{
			X: target.GetPosition().X - attacker.GetPosition().X,
			Z: target.GetPosition().Z - attacker.GetPosition().Z,
		}
		attacker.SetFacingDirection(dir.Normalize())
	}

	// ⚡ Movimento (Leap, Dash, etc)
	if skill.Movement != nil {
		if target != nil {
			attacker.SetSkillMovementState(ApplySkillMovement(attacker, target, skill))
		} else {
			log.Printf("[SkillExecutor] [%s] Skill %s possui movimento, mas target é nil", attacker.GetHandle().ID, skill.Name)
		}
		result.Success = true
		return result
	}

	// 🌀 Área no chão (Ground AOE)
	if skill.GroundTargeted && skill.AOE != nil {
		ApplyAOEDamage(attacker, targetPos, skill, creatures, players, svcCtx)
		result.Success = true
		return result
	}

	// 🏹 Projétil
	if skill.Projectile != nil {
		SimulateProjectile(attacker, target, targetPos, skill, svcCtx)
		result.Success = true
		return result
	}

	// 🗡️ Dano direto
	result = ApplyDirectDamage(attacker, target, skill, svcCtx)
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

	moveDist := state.Speed * deltaTime
	if moveDist > distanceToTargetPos {
		moveDist = distanceToTargetPos
	}
	moveVec := state.Direction.Scale(moveDist)
	newPos := currentPos.AddVector3D(moveVec)
	mov.SetPosition(newPos)

	// 💥 Ativa bloqueio frontal se configurado
	if state.Config.BlockDuringMovement {
		mov.SetBlocking(true)
	}

	// 🔥 HITBOX DURANTE MOVIMENTO
	if state.Skill.Hitbox != nil {
		nearby := finder.FindNearbyTargets(svcCtx, mov, state.Config.MaxDistance+1.0)
		for _, other := range nearby {
			if other == nil || other.GetHandle().Equals(mov.GetHandle()) {
				continue
			}
			if shouldSkipTarget(mov, other) {
				continue
			}
			if state.HasAlreadyHit(other) {
				continue
			}
			if !isTargetInsideHitbox(mov, other, state.Skill) {
				continue
			}

			log.Printf("[MOVE-HIT] [%s] atingiu [%s] durante movimento %s",
				mov.GetHandle().ID,
				other.GetHandle().ID,
				state.Skill.Name,
			)

			ApplyDirectDamage(mov, other, state.Skill, svcCtx)
			state.MarkAsHit(other)

			// 💥 Push contínuo: empurrar alvo junto até o fim
			if state.Config.PushType == consts.MoveToImpact && state.EngagedTarget == nil {
				state.EngagedTarget = other
			}

			// 🛑 Stop imediato após primeiro impacto, se configurado
			if state.Config.StopOnFirstHit && state.EngagedTarget == nil {
				state.EngagedTarget = other
				state.StartTime = time.Now().Add(-time.Duration(state.Config.DurationSec * 0.99 * float64(time.Second)))
				break
			}

			// ⚙️ Push configurável alternativo
			if state.Config.PushTargetDuringMovement && state.EngagedTarget == nil {
				state.EngagedTarget = other
			}
		}
	}

	// 💨 Empurra o alvo engatado continuamente durante o movimento
	if state.Config.PushTargetDuringMovement && state.EngagedTarget != nil && state.EngagedTarget.IsAlive() {
		// Ativa empurrão inicial
		state.EngagedTarget.ApplyImpulseFrom(mov.GetPosition(), 300*time.Millisecond)

		// Corrige sobreposição contínua enquanto se move
		dist := position.CalculateDistance2D(mov.GetPosition(), state.EngagedTarget.GetPosition())
		if dist < 0.5 {
			dir := position.NewVector2DFromTo(mov.GetPosition(), state.EngagedTarget.GetPosition()).Normalize()
			push := dir.Scale(0.5 - dist)
			state.EngagedTarget.SetPosition(state.EngagedTarget.GetPosition().AddVector2D(push))
		}
	}

	// ✅ Fim do movimento + separação
	if state.IsComplete(time.Now(), newPos) {
		if state.Config.SeparationRadius > 0 && state.Config.SeparationForce > 0 {
			nearby := finder.FindNearbyTargets(svcCtx, mov, state.Config.SeparationRadius)
			for _, other := range nearby {
				if other.GetHandle().Equals(mov.GetHandle()) || !other.IsAlive() {
					continue
				}
				dir := position.NewVector2DFromTo(mov.GetPosition(), other.GetPosition()).Normalize()
				separation := dir.Scale(state.Config.SeparationForce)
				pushedPos := other.GetPosition().AddVector2D(separation)
				other.SetPosition(pushedPos)

				log.Printf("[SEPARATION] [%s] empurrou [%s] após fim do movimento %s",
					mov.GetHandle().ID,
					other.GetHandle().ID,
					state.Skill.Name,
				)
			}
		}

		// 🛑 Desativa bloqueio e limpa estado
		if state.Config.BlockDuringMovement {
			mov.SetBlocking(false)
		}
		state.EngagedTarget = nil
		return true
	}

	return false
}

func ApplyDirectDamage(attacker model.Attacker, target model.Targetable, skill *model.Skill, svcCtx *dynamic_context.AIServiceContext) model.SkillResult {
	log.Printf("[SKILL] [%s] tentando aplicar %s em [%s]",
		attacker.GetPrimaryType(),
		skill.Name,
		target.GetHandle().ID,
	)
	result := model.SkillResult{}

	if target == nil || shouldSkipTarget(attacker, target) {
		log.Printf("[SKILL] [%s] alvo inválido ou deve ser ignorado", target.GetHandle().ID)
		attacker.SetLastMissedSkillAt(time.Now())
		return result
	}

	if target.IsInvulnerableNow() {
		log.Printf("[SKILL] [%s] invulnerável — ataque de [%s] anulado", target.GetHandle().ID, attacker.GetPrimaryType())
		attacker.SetLastMissedSkillAt(time.Now())
		return result
	}

	if skill.Hitbox != nil && !isTargetInsideHitbox(attacker, target, skill) {
		log.Printf("[HITBOX] [%s] fora da hitbox de [%s] para skill %s",
			target.GetHandle().ID,
			attacker.GetPrimaryType(),
			skill.Name,
		)
		attacker.SetLastMissedSkillAt(time.Now())
		return result
	}

	dir := position.NewVector2DFromTo(attacker.GetPosition(), target.GetPosition())
	attacker.SetFacingDirection(dir)

	damage := calculateDamageGeneric(attacker, target, skill.InitialMultiplier)

	// 🔄 BLOQUEIO E PARRY
	if target.IsBlocking() {
		blockDir := target.GetFacingDirection()
		attackDir := position.NewVector2DFromTo(target.GetPosition(), attacker.GetPosition())
		dot := blockDir.Dot(attackDir)

		if dot > 0.5 {
			if target.IsInParryWindow() {
				log.Printf("[PARRY] [%s] parry bem-sucedido contra [%s]", target.GetHandle().ID, attacker.GetPrimaryType())
				attacker.SetLastMissedSkillAt(time.Now())
				return result
			}

			log.Printf("[BLOCK] [%s] bloqueou ataque de [%s]", target.GetHandle().ID, attacker.GetPrimaryType())
			if skill.Impact != nil && skill.Impact.PostureDamage > 0 {
				postureDamage := skill.Impact.PostureDamage + getPostureScaling(attacker, skill)
				target.ApplyPostureDamage(postureDamage * 2)
				log.Printf("[BLOCK] [%s] sofreu %.1f de dano de postura (dobrado)", target.GetHandle().ID, postureDamage*2)
			}
			attacker.SetLastMissedSkillAt(time.Now())
			return result
		} else {
			log.Printf("[BLOCK] [%s] bloqueou em direção errada — ataque de [%s] passou",
				target.GetHandle().ID, attacker.GetPrimaryType())
		}
	}

	target.TakeDamage(damage)

	result.TargetDied = !target.IsAlive()

	if skill.Impact != nil && skill.Impact.PostureDamage > 0 {
		postureDamage := skill.Impact.PostureDamage + getPostureScaling(attacker, skill)
		target.ApplyPostureDamage(postureDamage)
		log.Printf("[POSTURE] [%s] sofreu %.1f de dano de postura de [%s]",
			target.GetHandle().ID,
			postureDamage,
			attacker.GetPrimaryType(),
		)
	}

	if skill.HasDOT && skill.DOT != nil {
		dotPower := damage / (skill.DOT.DurationSec / skill.DOT.TickSec)
		effect := consts.ActiveEffect{
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
		log.Printf("[DOT] [%s] aplicou efeito %s em [%s] por %ds",
			attacker.GetPrimaryType(),
			skill.DOT.EffectType,
			target.GetHandle().ID,
			skill.DOT.DurationSec,
		)
	}

	result.Success = true
	return result
}

// Verifica se o alvo está dentro da hitbox da skill
func isTargetInsideHitbox(attacker model.Attacker, target model.Targetable, skill *model.Skill) bool {
	if skill.Hitbox == nil {
		return true
	}

	shape := skill.Hitbox.Shape
	attackerPos := attacker.GetPosition().ToVector2D()
	targetPos := target.GetPosition().ToVector2D()

	switch shape {
	case model.HitboxBox:
		return isInsideBox(attackerPos, targetPos, attacker.GetFacingDirection(), skill.Hitbox.Length, skill.Hitbox.Width)

	case model.HitboxCone:
		dist := position.CalculateDistance(attacker.GetPosition(), target.GetPosition())
		if dist < skill.Hitbox.MinRadius || dist > skill.Hitbox.MaxRadius {
			return false
		}
		dirToTarget := position.NewVector2DFromTo(attacker.GetPosition(), target.GetPosition())
		angle := math.Acos(attacker.GetFacingDirection().Normalize().Dot(dirToTarget.Normalize())) * 180 / math.Pi
		return angle <= (skill.Hitbox.Angle / 2)

	case model.HitboxCircle:
		dist := position.CalculateDistance(attacker.GetPosition(), target.GetPosition())
		return dist <= skill.Hitbox.MaxRadius

	case model.HitboxLine:
		return isInsideLine(attackerPos, targetPos, attacker.GetFacingDirection(), skill.Hitbox.Length, skill.Hitbox.Width)

	default:
		return false
	}
}

// Função auxiliar para box hitbox
func isInsideBox(origin, point position.Vector2D, facing position.Vector2D, length, width float64) bool {
	dir := facing.Normalize()
	rel := point.Sub(origin)

	forward := dir.Dot(rel)
	if forward < 0 || forward > length {
		return false
	}

	side := dir.Perpendicular().Dot(rel)
	return math.Abs(side) <= width/2
}

// Função auxiliar para linha tipo Leap
func isInsideLine(start, point position.Vector2D, facing position.Vector2D, length, width float64) bool {
	return isInsideBox(start, point, facing, length, width)
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

func applySeparationOnSkillEnd(attacker model.Attacker, radius float64, force float64, svcCtx *dynamic_context.AIServiceContext) {
	nearby := finder.FindNearbyTargets(svcCtx, attacker, radius)

	for _, t := range nearby {
		if t.GetHandle().Equals(attacker.GetHandle()) || !t.IsAlive() {
			continue
		}

		dir := position.NewVector2DFromTo(attacker.GetPosition(), t.GetPosition()).Normalize()
		separation := dir.Scale(force)
		pushedPos := t.GetPosition().AddVector2D(separation)
		t.SetPosition(pushedPos)

		log.Printf("[SEPARATION] [%s] empurrou [%s] após movimento", attacker.GetHandle().ID, t.GetHandle().ID)
	}
}
