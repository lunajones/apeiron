package defensive

import (
	"log"
	"time"

	constslib "github.com/lunajones/apeiron/lib/consts"
	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
	"github.com/lunajones/apeiron/service/helper/finder"
)

type RetreatNode struct{}

func (n *RetreatNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[RETREAT] [%s] contexto inválido", c.Handle.String())
		return core.StatusFailure
	}

	target := finder.FindTargetByHandles(c.Handle, c.TargetCreatureHandle, c.TargetPlayerHandle, svcCtx)
	if target == nil {
		log.Printf("[RETREAT] [%s] sem alvo para recuar", c.Handle.String())
		return core.StatusFailure
	}

	drive := c.GetCombatDrive()

	// Quanto mais cautela, mais longe recua; quanto mais raiva, mais hesita
	recoilFactor := 1.5 + drive.Caution*2.0 - drive.Rage*0.5
	if recoilFactor < 1.0 {
		recoilFactor = 1.0
	}

	dir := position.NewVector3DFromTo(target.GetPosition(), c.GetPosition()).Normalize()
	dest := c.GetPosition().AddVector3D(dir.Scale(recoilFactor))

	c.MoveCtrl.SetMoveIntent(dest, c.WalkSpeed*0.5, 0.0)
	c.SetAnimationState(constslib.AnimationWalk)

	log.Printf("[RETREAT] [%s] recuando para %.2f, %.2f, factor=%.2f", c.Handle.String(), dest.X, dest.Z, recoilFactor)

	now := time.Now()
	event := model.CombatEvent{
		SourceHandle: c.Handle,
		TargetHandle: target.GetHandle(),
		BehaviorType: "RetreatBroadcast",
		Timestamp:    now,
	}

	c.RegisterCombatEvent(event)

	if targetCreature, ok := target.(*creature.Creature); ok {
		targetCreature.RegisterCombatEvent(event)
		log.Printf("[RETREAT] [%s] alvo [%s] registrou hesitação", c.Handle.String(), target.GetHandle().String())
	}

	// Atualiza impulsos mentais da criatura
	drive.Caution += 0.05 // ficou mais cautelosa após o recuo
	drive.Rage += 0.02    // pode ficar irritada por estar sendo forçada a recuar
	drive.Termination = 0
	drive.LastUpdated = now
	drive.Value = creature.RecalculateCombatDrive(drive)

	return core.StatusRunning
}

func (n *RetreatNode) Reset(c *creature.Creature) {
	c.ClearMovementIntent()
	c.SetAnimationState(constslib.AnimationIdle)
}
