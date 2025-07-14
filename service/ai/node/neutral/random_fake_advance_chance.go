package neutral

import (
	"log"
	"math/rand"
	"time"

	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/service/ai/core"
	"github.com/lunajones/apeiron/service/ai/dynamic_context"
	"github.com/lunajones/apeiron/service/creature"
)

type RandomFakeAdvanceChanceNode struct {
	Chance float64 // Ex: 0.2 para 20%
}

func (n *RandomFakeAdvanceChanceNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
	if rand.Float64() >= n.Chance {
		log.Printf("[FAKE ADVANCE] [%s] decidiu não blefar", c.Handle.String())
		return core.StatusFailure
	}

	svcCtx, ok := ctx.(*dynamic_context.AIServiceContext)
	if !ok {
		log.Printf("[FAKE ADVANCE] [%s] contexto inválido", c.Handle.String())
		return core.StatusFailure
	}

	now := time.Now()
	log.Printf("[FAKE ADVANCE] [%s] decidiu blefar", c.Handle.String())

	targetHandle := c.TargetCreatureHandle
	if !targetHandle.IsValid() {
		log.Printf("[FAKE ADVANCE] [%s] alvo inválido para blefe", c.Handle.String())
		return core.StatusFailure
	}

	// Monta evento FALSO de agressão (sem skill real, mas com Windup simulado)
	event := model.CombatEvent{
		SourceHandle: c.Handle,
		TargetHandle: targetHandle,
		BehaviorType: "FakeAdvanceBroadcast",
		Timestamp:    now,
	}

	// Registra o evento em si mesmo (para rastreamento ou possíveis logísticas internas)
	c.RegisterCombatEvent(event)

	// Registra o evento no alvo (efeito de intimidação ou engano)
	target := svcCtx.FindByHandle(targetHandle)
	if targetCreature, ok := target.(*creature.Creature); ok {
		targetCreature.RegisterCombatEvent(event)
		log.Printf("[FAKE ADVANCE] [%s] alvo [%s] recebeu evento de blefe", c.Handle.String(), targetHandle.String())
	}

	// Atualiza o drive: reduz cautela, zera tédio, reacende iniciativa
	drive := c.GetCombatDrive()
	drive.Caution -= 0.1
	if drive.Caution < 0 {
		drive.Caution = 0
	}
	drive.Termination = 0
	drive.LastUpdated = now
	drive.Value = creature.RecalculateCombatDrive(drive)

	return core.StatusSuccess
}

func (n *RandomFakeAdvanceChanceNode) Reset(c *creature.Creature) {}
