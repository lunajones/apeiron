package consts

import (
	"time"
)

type EffectType string

const (
    // DOTs
    EffectBleed      EffectType = "Bleed"
    EffectPoison     EffectType = "Poison"
    EffectBurn       EffectType = "Burn"
    EffectFreeze     EffectType = "Freeze"

    // Control (CC)
    EffectSlow       EffectType = "Slow"
    EffectStun       EffectType = "Stun"
    EffectInterrupt  EffectType = "Interrupt"
    EffectStagger    EffectType = "Stagger"
    EffectKnockback  EffectType = "Knockback"
    EffectFear       EffectType = "Fear"

    // Buffs
    EffectBerserk    EffectType = "Berserk"    // +Atk, +Speed, +Regen, -Cooldown, -ConjurationTime, -Defense
    EffectFocus      EffectType = "Focus"      // +Accuracy, +Critical
    EffectInsanity   EffectType = "Insanity"   // Modifica comportamento, +Speed, -Cooldown
    EffectShield     EffectType = "Shield"     // Extra HP temporário
    EffectRegen      EffectType = "Regen"      // Heal over Time
)

type ActiveEffect struct {
	Type          EffectType
	StartTime     time.Time
	Duration      time.Duration
	TickInterval  time.Duration
	LastTickTime  time.Time
	Power         int   // Intensidade (ex: quanto de dano ou quanto de slow)
	IsDOT         bool  // Se é Damage Over Time
	IsDebuff      bool  // Se é um efeito negativo
	IsCC          bool  // Se é controle de grupo (stun, etc)
}


// ---- Categorização por comportamento de jogo ----

func (e EffectType) IsAggressiveBuff() bool {
    return e == EffectBerserk || e == EffectInsanity
}

func (e EffectType) IsDefensiveBuff() bool {
    return e == EffectShield || e == EffectRegen
}

func (e EffectType) IsDebuff() bool {
    return e == EffectStun || e == EffectSlow || e == EffectInterrupt || e == EffectStagger || e == EffectKnockback || e == EffectFear || e == EffectFreeze
}

func (e EffectType) IsDOT() bool {
    return e == EffectBleed || e == EffectPoison || e == EffectBurn
}

// ---- Mapa centralizado de efeito visual ----

var visualEffectKeys = map[EffectType]string{
    EffectBleed:     "bleed_overlay",
    EffectPoison:    "poison_green_cloud",
    EffectBurn:      "burning_flames",
    EffectFreeze:    "freeze_ice_shards",
    EffectSlow:      "slow_blue_glow",
    EffectStun:      "stun_stars_above_head",
    EffectInterrupt: "interrupt_flash",
    EffectStagger:   "stagger_shake",
    EffectKnockback: "knockback_dust",
    EffectFear:      "fear_dark_aura",
    EffectBerserk:   "red_aura_pulse",
    EffectFocus:     "focus_glow",
    EffectInsanity:  "insanity_distortion",
    EffectShield:    "shield_barrier_effect",
    EffectRegen:     "regen_green_particles",
}

// ---- Método de acesso ao visual ----

func (e EffectType) VisualEffectKey() string {
    if key, exists := visualEffectKeys[e]; exists {
        return key
    }
    return ""
}