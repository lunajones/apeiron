package sensor

import (
	"math"

	"github.com/lunajones/apeiron/lib/position"
	"github.com/lunajones/apeiron/service/creature"
)

// CanSee verifica se o target está dentro do cone de visão do creature
func CanSee(c *creature.Creature, target *creature.Creature) bool {
	dist := position.CalculateDistance2D(c.GetPosition(), target.GetPosition())
	if dist > c.VisionRange {
		return false
	}

	dir := position.NewVector3DFromTo(c.GetPosition(), target.GetPosition()).ToVector2D().Normalize()
	dot := c.GetFacingDirection().Dot(dir)

	angleRad := c.FieldOfViewDegrees * (math.Pi / 180.0)
	fovCos := math.Cos(angleRad / 2)

	return dot >= fovCos
}

// CanHear verifica se o target está dentro do alcance auditivo
func CanHear(c *creature.Creature, target *creature.Creature) bool {
	basePos := c.GetPosition().ToVector2D()
	facing := c.GetFacingDirection()

	// Calcula posições das orelhas (offset lateral)
	offset := 0.5 // você ajusta o quanto a orelha está afastada
	leftEar := basePos.Add(facing.PerpLeft().Scale(offset))
	rightEar := basePos.Add(facing.PerpRight().Scale(offset))

	targetPos := target.GetPosition().ToVector2D()

	// Distâncias para cada orelha
	distLeft := leftEar.Sub(targetPos).Length()
	distRight := rightEar.Sub(targetPos).Length()

	// Se cair em um dos círculos auditivos
	return distLeft <= c.HearingRange || distRight <= c.HearingRange
}

// CanSmell verifica se o target está dentro do alcance olfativo
func CanSmell(c *creature.Creature, target *creature.Creature) bool {
	dist := position.CalculateDistance2D(c.GetPosition(), target.GetPosition())
	return dist <= c.SmellRange
}
