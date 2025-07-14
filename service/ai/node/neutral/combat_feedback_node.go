package neutral

// import (
// 	"math"

// 	"github.com/fatih/color"
// 	constslib "github.com/lunajones/apeiron/lib/consts"
// 	"github.com/lunajones/apeiron/service/ai/core"
// 	"github.com/lunajones/apeiron/service/creature"
// )

// type CombatFeedbackNode struct{}

// func (n *CombatFeedbackNode) Tick(c *creature.Creature, ctx interface{}) interface{} {
// 	drive := c.GetCombatDrive()
// 	seen := make(map[constslib.CombatAction]bool)

// 	// Aplica efeitos únicos por tipo de ação
// 	for _, action := range c.RecentActions {
// 		if seen[action] {
// 			continue
// 		}
// 		seen[action] = true

// 		switch action {
// 		case constslib.CombatActionBlockSuccess:
// 			drive.Rage += 0.1
// 			drive.Caution -= 0.03
// 			drive.Termination += 0.03
// 			drive.Counter += 0.2

// 		case constslib.CombatActionParrySuccess:
// 			drive.Rage += 0.1
// 			drive.Caution -= 0.06
// 			drive.Termination += 0.45
// 			drive.Counter += 0.45

// 		case constslib.CombatActionDodgeSuccess:
// 			drive.Rage += 0.1
// 			drive.Caution -= 0.1
// 			drive.Termination += 0.03
// 			drive.Counter += 0.2

// 		case constslib.CombatActionMicroRetreat:
// 			drive.Rage -= 0.02
// 			drive.Caution += 0.04
// 			drive.Termination += 0.01
// 			drive.Counter += 0.2

// 		case constslib.CombatActionCircleAround:
// 			drive.Rage += 0.08
// 			drive.Caution -= 0.01
// 			drive.Termination += 0.02
// 			drive.Counter += 0.05

// 		case constslib.CombatActionApproach:
// 			drive.Rage += 0.02
// 			drive.Caution -= 0.02
// 			drive.Termination += 0.01
// 			drive.Counter += 0.05

// 		case constslib.CombatActionChase:
// 			drive.Rage += 0.03
// 			drive.Caution -= 0.02
// 			drive.Termination += 0.03
// 			drive.Counter += 0.05

// 		case constslib.CombatActionAttackPrepared:
// 			drive.Rage -= 0.005
// 			drive.Caution += 0.005
// 			drive.Termination += 0.005
// 			drive.Counter -= 0.001

// 		case constslib.CombatActionAttackSuccess:
// 			drive.Rage += 0.08
// 			drive.Caution -= 0.02
// 			drive.Termination += 0.06
// 			drive.Counter += 0.05

// 		case constslib.CombatActionAttackMissed:
// 			drive.Rage += 0.02
// 			drive.Caution += 0.04
// 			drive.Termination += 0.07
// 			drive.Counter += 0.07

// 		case constslib.CombatActionSkillInterrupted:
// 			drive.Rage -= 0.05
// 			drive.Caution += 0.06
// 			drive.Termination -= 0.01
// 			drive.Counter += 0.12

// 		case constslib.CombatActionCounter:
// 			drive.Rage += 0.01
// 			drive.Caution -= 0.01
// 			drive.Termination -= 0.01
// 			drive.Counter = 0.0

// 		case constslib.CombatActionTookDamage:
// 			drive.Rage += 0.3
// 			drive.Caution -= 0.04
// 			drive.Termination += 0.05
// 			drive.Counter += 0.2
// 		}
// 	}

// 	// Decaimento
// 	drive.Rage *= 0.995
// 	drive.Caution *= 0.995
// 	drive.Termination *= 0.995
// 	drive.Counter *= 0.97

// 	// Clamp entre 0.0 e 1.0
// 	drive.Rage = math.Max(0, math.Min(1, drive.Rage))
// 	drive.Caution = math.Max(0, math.Min(1, drive.Caution))
// 	drive.Termination = math.Max(0, math.Min(1, drive.Termination))
// 	drive.Counter = math.Max(0, math.Min(1, drive.Counter))

// 	// Log visual
// 	color.New(color.FgHiMagenta, color.Bold).Printf(
// 		"[COMBAT-FEEDBACK] [%s] drive atualizado: Rage=%.2f | Caution=%.2f | Termination=%.2f | Counter=%.2f\n",
// 		c.PrimaryType, drive.Rage, drive.Caution, drive.Termination, drive.Counter,
// 	)

// 	// Limpa ações
// 	c.RecentActions = nil
// 	return core.StatusSuccess
// }

// func (n *CombatFeedbackNode) Reset(c *creature.Creature) {}
