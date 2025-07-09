package neutral

import (
	"log"
	"math/rand"
	"time"

	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

// EvaluateCombatBehaviorsNode decide se executa ações intermediárias de combate (ex.: circular, ameaçar)
type EvaluateCombatBehaviorsNode struct {
	CircleChance      float64
	ProvokeChance     float64
	FakeAdvanceChance float64
	lastEvalTime      time.Time
}

func (n *EvaluateCombatBehaviorsNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[COMBAT-BEHAVIOR] [%s] contexto inválido", c.Handle.String())
		return core.StatusFailure
	}

	// Evitar avaliação muito frequente (caso não use decorador externo)
	if time.Since(n.lastEvalTime) < 1*time.Second {
		return core.StatusFailure
	}
	n.lastEvalTime = time.Now()

	r := rand.Float64()
	if r < n.CircleChance {
		n.performCircle(c, svcCtx)
		return core.StatusSuccess
	}
	r -= n.CircleChance

	if r < n.ProvokeChance {
		n.performProvoke(c, svcCtx)
		return core.StatusSuccess
	}
	r -= n.ProvokeChance

	if r < n.FakeAdvanceChance {
		n.performFakeAdvance(c, svcCtx)
		return core.StatusSuccess
	}

	// Não escolheu nada
	return core.StatusFailure
}

func (n *EvaluateCombatBehaviorsNode) performCircle(c *creature.Creature, svcCtx *dynamic_context.AIServiceContext) {
	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)

	if target == nil {
		log.Printf("[COMBAT-BEHAVIOR] [%s] sem alvo para circular", c.Handle.String())
		return
	}

	// Direção lateral aleatória (perpendicular ao alvo)
	dirVec3D := target.GetPosition().Sub(c.Position).Normalize()
	dirVec := dirVec3D.ToVector2D().Normalize()

	perp := position.RotateVector2D(dirVec, 1.5708)
	if rand.Float64() < 0.5 {
		perp = perp.Multiply(-1)
	}

	dest := c.Position.AddVector3D(position.Vector3D{X: perp.X, Y: 0, Z: perp.Z}.Multiply(1.0))
	if svcCtx.NavMesh.IsWalkable(dest) {
		c.MoveCtrl.SetTarget(dest, c.WalkSpeed, 0.2)
		log.Printf("[COMBAT-BEHAVIOR] [%s] circulando alvo", c.Handle.String())
	}
}

func (n *EvaluateCombatBehaviorsNode) performProvoke(c *creature.Creature, svcCtx *dynamic_context.AIServiceContext) {
	log.Printf("[COMBAT-BEHAVIOR] [%s] provocação: rosnado, ofendendo ou ameaçando", c.Handle.String())
	// Exemplo: c.SetAnimationState(constslib.AnimationTaunt)
}

func (n *EvaluateCombatBehaviorsNode) performFakeAdvance(c *creature.Creature, svcCtx *dynamic_context.AIServiceContext) {
	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)

	if target == nil {
		log.Printf("[COMBAT-BEHAVIOR] [%s] sem alvo para fake advance", c.Handle.String())
		return
	}

	dirVec := target.GetPosition().Sub(c.Position).Normalize()
	dest := c.Position.AddVector3D(dirVec.Multiply(0.5))
	if svcCtx.NavMesh.IsWalkable(dest) {
		c.MoveCtrl.SetTarget(dest, c.WalkSpeed, 0.2)
		log.Printf("[COMBAT-BEHAVIOR] [%s] avanço falso no alvo", c.Handle.String())
	}
}

func (n *EvaluateCombatBehaviorsNode) Reset() {
	// Nada necessário por enquanto
}
