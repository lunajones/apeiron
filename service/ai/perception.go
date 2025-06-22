package ai

import (
	"log"
	"math"

	"github.com/lunajones/apeiron/service/creature"
)

// Representação simplificada de um player, até você criar o PlayerService real
type Player struct {
	ID       string
	Position creature.Position
}

// Esta função deverá ser substituída no futuro por uma chamada ao PlayerService real
func GetPlayersInWorld() []Player {
	// TODO: No futuro, busque os players reais
	return []Player{}
}

// Verifica se a criatura pode ver algum player
func CanSeePlayer(c *creature.Creature, players []Player) *Player {
	if c.IsBlind || c.VisionRange <= 0 || c.FieldOfViewDegrees <= 0 {
		return nil
	}

	for _, p := range players {
		dx := p.Position.X - c.Position.X
		dz := p.Position.Z - c.Position.Z
		distance := math.Sqrt(dx*dx + dz*dz)

		if distance > c.VisionRange {
			continue
		}

		// TODO: No futuro, use direção real da criatura
		playerAngle := math.Atan2(dz, dx) * (180 / math.Pi)
		creatureFacingAngle := 0.0 // Placeholder: ângulo que a criatura está olhando

		angleDiff := math.Abs(playerAngle - creatureFacingAngle)
		if angleDiff <= c.FieldOfViewDegrees/2 {
			log.Printf("[Perception] Creature %s vê o player %s! Distância: %.2f, Ângulo: %.2f", c.ID, p.ID, distance, angleDiff)
			return &p
		}
	}

	return nil
}

// Verifica se a criatura pode ouvir algum player
func CanHearPlayer(c *creature.Creature, players []Player) *Player {
	if c.IsDeaf || c.HearingRange <= 0 {
		return nil
	}

	for _, p := range players {
		dx := p.Position.X - c.Position.X
		dz := p.Position.Z - c.Position.Z
		distance := math.Sqrt(dx*dx + dz*dz)

		if distance <= c.HearingRange {
			log.Printf("[Perception] Creature %s ouviu o player %s! Distância: %.2f", c.ID, p.ID, distance)
			return &p
		}
	}

	return nil
}

func CanSeeOtherCreatures(c *creature.Creature, creatures []*creature.Creature) *creature.Creature {
	if c.IsBlind || c.VisionRange <= 0 || c.FieldOfViewDegrees <= 0 {
		return nil
	}

	for _, target := range creatures {
		if target.ID == c.ID || !target.IsAlive {
			continue
		}

		dx := target.Position.X - c.Position.X
		dz := target.Position.Z - c.Position.Z
		distance := math.Sqrt(dx*dx + dz*dz)

		if distance > c.VisionRange {
			continue
		}

		// TODO: Usar direção real da criatura no futuro
		targetAngle := math.Atan2(dz, dx) * (180 / math.Pi)
		creatureFacingAngle := 0.0 // Placeholder

		angleDiff := math.Abs(targetAngle - creatureFacingAngle)
		if angleDiff <= c.FieldOfViewDegrees/2 {
			log.Printf("[Perception] Creature %s vê outra criatura %s! Distância: %.2f, Ângulo: %.2f", c.ID, target.ID, distance, angleDiff)
			return target
		}
	}

	return nil
}

// Detecta se há outra criatura dentro do alcance de audição
func CanHearOtherCreatures(c *creature.Creature, creatures []*creature.Creature) *creature.Creature {
	if c.IsDeaf || c.HearingRange <= 0 {
		return nil
	}

	for _, target := range creatures {
		if target.ID == c.ID || !target.IsAlive {
			continue
		}

		dx := target.Position.X - c.Position.X
		dz := target.Position.Z - c.Position.Z
		distance := math.Sqrt(dx*dx + dz*dz)

		if distance <= c.HearingRange {
			log.Printf("[Perception] Creature %s ouviu outra criatura %s! Distância: %.2f", c.ID, target.ID, distance)
			return target
		}
	}

	return nil
}
