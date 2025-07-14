package neutral

// import (
// 	"log"
// 	"math"
// 	"math/rand"
// 	"time"

// 	constslib "github.com/lunajones/apeiron/lib/consts"
// 	"github.com/lunajones/apeiron/lib/position"
// 	"github.com/lunajones/apeiron/service/ai/core"
// 	"github.com/lunajones/apeiron/service/ai/dynamic_context"
// 	"github.com/lunajones/apeiron/service/creature"
// 	"github.com/lunajones/apeiron/service/helper/finder"
// )

// type FlankOrRepositionNode struct{}

// func (n *FlankOrRepositionNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
// 	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
// 	if !ok {
// 		log.Printf("[FLANK] [%s] contexto inválido", c.Handle.String())
// 		return core.StatusFailure
// 	}

// 	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)
// 	if target == nil {
// 		log.Printf("[FLANK] [%s] sem alvo para flanquear", c.Handle.String())
// 		return core.StatusFailure
// 	}

// 	dirVec3D := target.GetPosition().Sub(c.Position)
// 	dirVec2D := position.Vector2D{X: dirVec3D.X, Z: dirVec3D.Z}.Normalize()

// 	drive := c.GetCombatDrive()

// 	// Escolhe plano tático com base nos impulsos mentais da criatura
// 	roll := rand.Float64()
// 	var plan constslib.CombatState
// 	switch {
// 	case roll < 0.15+drive.Rage*0.5:
// 		plan = constslib.CombatStatePlanAdvance
// 	case roll < 0.4+drive.Value*0.3:
// 		plan = constslib.CombatStatePlanFlank
// 	case roll < 0.6+drive.Caution*0.4+(1.0-drive.Value)*0.2:
// 		plan = constslib.CombatStatePlanRetreat
// 	case roll < 0.85+(1.0-drive.Value)*0.3:
// 		plan = constslib.CombatStatePlanGuard
// 	default:
// 		plan = constslib.CombatStatePlanFake
// 	}

// 	// Define ângulo de movimento com base no plano
// 	var angle float64
// 	switch plan {
// 	case constslib.CombatStatePlanAdvance:
// 		angle = 0.0
// 	case constslib.CombatStatePlanFlank:
// 		if rand.Intn(2) == 0 {
// 			angle = math.Pi / 3 // direita
// 		} else {
// 			angle = -math.Pi / 3 // esquerda
// 		}
// 	case constslib.CombatStatePlanRetreat:
// 		angle = math.Pi
// 	case constslib.CombatStatePlanGuard:
// 		angle = math.Pi / 6
// 	case constslib.CombatStatePlanFake:
// 		angle = -math.Pi
// 	}

// 	dist := 1.2 + rand.Float64()*1.0
// 	rotated := position.RotateVector2D(dirVec2D, angle).Multiply(dist)
// 	dest := c.Position.AddVector3D(position.Vector3D{X: rotated.X, Y: 0, Z: rotated.Z})

// 	if position.CalculateDistance(c.MoveCtrl.CurrentIntentDest, dest) < 0.5 {
// 		return core.StatusRunning
// 	}

// 	if svcCtx.NavMesh.IsWalkable(dest) {
// 		c.MoveCtrl.SetMoveIntent(dest, c.WalkSpeed, 0.0)
// 		c.SetAnimationState(constslib.AnimationWalk)

// 		// Atualiza drive com impulso emocional
// 		drive.Caution += 0.02
// 		drive.Rage += 0.01
// 		drive.Termination = 0
// 		drive.Plan = plan
// 		drive.LastUpdated = time.Now()

// 		c.RecalculateDrive()

// 		log.Printf("[FLANK] [%s] plano: %s → movendo para (%.2f, %.2f)",
// 			c.Handle.String(),
// 			combatPlanName(plan),
// 			dest.X, dest.Z,
// 		)

// 		return core.StatusSuccess
// 	}

// 	log.Printf("[FLANK] [%s] tentativa falhou, destino inalcançável", c.Handle.String())
// 	return core.StatusFailure
// }

// func (n *FlankOrRepositionNode) Reset(c *creature.Creature) {
// 	c.ClearMovementIntent()
// 	c.SetAnimationState(constslib.AnimationIdle)
// }

// func combatPlanName(plan constslib.CombatState) string {
// 	switch plan {
// 	case constslib.CombatStatePlanAdvance:
// 		return "Advance"
// 	case constslib.CombatStatePlanFlank:
// 		return "Flank"
// 	case constslib.CombatStatePlanRetreat:
// 		return "Retreat"
// 	case constslib.CombatStatePlanGuard:
// 		return "Guard"
// 	case constslib.CombatStatePlanFake:
// 		return "Fake"
// 	default:
// 		return "Unknown"
// 	}
// }
