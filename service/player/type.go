package player

type PlayerRole string

const (
	RoleNone     PlayerRole = "none"
	RoleMerchant PlayerRole = "merchant"
	RoleHunter   PlayerRole = "hunter"
	RoleGuard    PlayerRole = "guard"
	// Adicione mais conforme sua lógica
)