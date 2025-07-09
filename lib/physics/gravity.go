package physics

func ApplyGravity(z float64, deltaTimeSec float64) float64 {
	const gravityAcc = -9.81
	return z + gravityAcc*deltaTimeSec
}
