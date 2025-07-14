package model

import (
	"time"
)

type CombatDrive struct {
	Value       float64   // Valor final entre 0.0 e 1.0
	LastUpdated time.Time // Controle de tempo para decaimento

	// Componentes que influenciam o Value
	Rage        float64 // quanto dano recente tomou
	Caution     float64 // quão perigoso percebe o ambiente
	Vengeance   float64 // quanto perdeu aliados
	Termination float64 // tempo sem combate significativo
	Counter     float64 // Diretriz tática atual da criatura (ex: CombatStatePlanFlank, CombatStatePlanRetreat, etc.)

}
