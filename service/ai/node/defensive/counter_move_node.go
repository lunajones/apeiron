package defensive

import (
	"log"
	"math"
	"time"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/movement"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type CounterMoveNode struct {
	actionRegistered bool
}

func (n *CounterMoveNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[COUNTER-MOVE] [%s] contexto inválido", c.Handle.String())
		return core.StatusFailure
	}

	// ⚠️ Verifica se o valor de Counter está alto o suficiente
	if c.GetCombatDrive().Counter < 0.5 && c.GetCombatDrive().Counter >= 0.8 {
		return core.StatusFailure
	}

	// ✅ Reseta o counter (ou decai) após consumir
	c.GetCombatDrive().Counter = 0

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)
	if target == nil {
		log.Printf("[COUNTER-MOVE] [%s] sem alvo válido", c.Handle.String())
		return core.StatusFailure
	}

	if !n.actionRegistered {
		c.RecentActions = append(c.RecentActions, constslib.CombatActionCounter)
		n.actionRegistered = true
	}

	// Direção entre criatura e alvo
	dirVec := target.GetPosition().Sub2D(c.Position).Normalize()
	perp := position.RotateVector2D(dirVec, math.Pi/2).Normalize()

	// Tenta mover para o lado direito ou esquerdo
	offsets := []position.Vector2D{
		perp.Multiply(2.0),
		perp.Multiply(-2.0),
	}

	for _, offset := range offsets {
		dest := c.Position.AddVector3D(position.Vector3D{X: offset.X, Y: 0, Z: offset.Z})
		if !svcCtx.NavMesh.IsWalkable(dest) {
			continue
		}

		c.MoveCtrl.ImpulseState = &movement.ImpulseMovementState{
			Active:   true,
			StartPos: c.Position,
			EndPos:   dest,
			Duration: 350 * time.Millisecond,
			Start:    time.Now(),
		}

		log.Printf("[COUNTER-MOVE] [%s] ativou impulse lateral para (%.2f, %.2f)", c.Handle.String(), dest.X, dest.Z)
		return core.StatusSuccess
	}

	log.Printf("[COUNTER-MOVE] [%s] nenhum destino lateral disponível", c.Handle.String())
	return core.StatusFailure
}

func (n *CounterMoveNode) Reset(c *creature.Creature) {}
