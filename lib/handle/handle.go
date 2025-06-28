package handle

import "fmt"

// EntityHandle identifica de forma segura uma entidade viva no mundo.
type EntityHandle struct {
	ID         string // ID único da entidade (UUID ou ID de sistema)
	Generation int    // Incrementado a cada respawn da entidade
}

// NewEntityHandle cria um novo handle.
func NewEntityHandle(id string, generation int) EntityHandle {
	return EntityHandle{ID: id, Generation: generation}
}

// Equals verifica se dois handles apontam para a mesma entidade e geração.
func (h EntityHandle) Equals(other EntityHandle) bool {
	return h.ID == other.ID && h.Generation == other.Generation
}

// IsValid confirma se o handle é válido.
func (h EntityHandle) IsValid() bool {
	return h.ID != "" && h.Generation > 0
}

// String formata o handle para logs/debug.
func (h EntityHandle) String() string {
	return fmt.Sprintf("Handle[ID=%s Gen=%d]", h.ID, h.Generation)
}

// Checa se o handle está vazio
func (h EntityHandle) IsEmpty() bool {
	return h.ID == "" && h.Generation == 0
}
