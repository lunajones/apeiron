package neutral

import (
	"math"
	"math/rand"
	"time"

	"github.com/fatih/color"
	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/movement"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type ApproachUntilInRangeNode struct{}

func (n *ApproachUntilInRangeNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	if c.MoveCtrl.MovementPlan != nil && c.MoveCtrl.MovementPlan.Type == constslib.MovementPlanApproach {
		if time.Now().Before(c.MoveCtrl.MovementPlan.ExpiresAt) {
			color.Yellow("[APPROACH-IN-RANGE] [%s] plano ativo, ignorando novo", c.GetPrimaryType())
			return core.StatusRunning
		} else {
			color.HiRed("[APPROACH-IN-RANGE] [%s] plano expirado", c.GetPrimaryType())
			c.MoveCtrl.MovementPlan = nil
			// FSM vai cuidar da limpeza
		}
	}

	if c.IsDodging() || c.CombatState != constslib.CombatStateMoving {
		return core.StatusFailure
	}

	drive := c.GetCombatDrive()
	if drive.Caution < 0.2 {
		color.Yellow("[APPROACH-IN-RANGE] [%s] Caution %.2f fora do intervalo", c.GetPrimaryType(), drive.Caution)
		return core.StatusFailure
	}

	if c.NextSkillToUse == nil {
		color.HiMagenta("[APPROACH-IN-RANGE] [%s] Sem skill planejada", c.GetPrimaryType())
		return core.StatusFailure
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, c.GetContext())
	if target == nil {
		color.HiRed("[APPROACH-IN-RANGE] [%s] alvo não encontrado", c.GetPrimaryType())
		return core.StatusFailure
	}

	c.AddRecentAction(constslib.CombatActionApproach)

	dist := position.CalculateDistance(c.GetPosition(), target.GetPosition())
	rangeNeeded := c.NextSkillToUse.Range

	if dist <= rangeNeeded {
		color.Green("[APPROACH-IN-RANGE] [%s] no range para %s (dist=%.2f)", c.GetPrimaryType(), c.NextSkillToUse.Name, dist)
		return core.StatusSuccess
	}

	// Tentativa de offset lateral com fallback
	offsetDir := rand.Float64()*2 - 1
	forwardVec := position.NewVector2DFromTo(c.GetPosition(), target.GetPosition()).Normalize()
	lateralVec := forwardVec.Perpendicular().Scale(offsetDir * 1.0)

	attempts := 0
	var destination position.Position
	for attempts < 4 {
		destination = target.GetPosition().AddVector2D(lateralVec)
		if c.GetContext().ClaimPosition(destination, c.GetHandle()) {
			break
		}
		angle := math.Pi / 4
		rotation := (rand.Float64()*2 - 1) * angle
		sin := math.Sin(rotation)
		cos := math.Cos(rotation)
		lateralVec = position.Vector2D{
			X: lateralVec.X*cos - lateralVec.Z*sin,
			Z: lateralVec.X*sin + lateralVec.Z*cos,
		}
		attempts++
	}

	if attempts >= 4 {
		color.Red("[APPROACH-IN-RANGE] [%s] falha ao encontrar célula livre próxima ao alvo", c.GetPrimaryType())
		return core.StatusFailure
	}

	// FSM já controla se está se movendo ou não, só planejar aqui
	c.MoveCtrl.SetMoveTarget(destination, c.WalkSpeed*0.5, rangeNeeded)

	c.MoveCtrl.MovementPlan = movement.NewMovementPlan(
		constslib.MovementPlanApproach,
		target.GetHandle(),
		rangeNeeded,
		5*time.Second,
		target.GetPosition(),
	)

	color.HiCyan("[APPROACH-IN-RANGE] [%s] avançando com cautela (dist=%.2f)", c.GetPrimaryType(), dist)
	return core.StatusRunning
}

func (n *ApproachUntilInRangeNode) Reset(c *creature.Creature) {
	color.Yellow("[APPROACH-IN-RANGE] [RESET] [%s] limpando plano", c.GetPrimaryType())
	c.MoveCtrl.MovementPlan = nil
	// FSM limpará o resto
}
