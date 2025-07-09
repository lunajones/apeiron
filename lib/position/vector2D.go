package position

import "math"

// Vector2D representa um vetor 2D no plano X-Z (Y = altura separada)
type Vector2D struct {
	X float64 // Coordenada X no plano
	Z float64 // Coordenada Z no plano
}

func (v Vector2D) Add(other Vector2D) Vector2D {
	return Vector2D{
		X: v.X + other.X,
		Z: v.Z + other.Z,
	}
}

// Magnitude retorna o comprimento (módulo) do vetor.
// É calculado como a raiz quadrada da soma dos quadrados dos componentes.
func (v Vector2D) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Z*v.Z)
}

// Normalize retorna um vetor unitário (mesma direção, magnitude 1).
// Se o vetor original tem magnitude zero, retorna o vetor (0,0).
func (v Vector2D) Normalize() Vector2D {
	mag := v.Magnitude()
	if mag == 0 {
		return Vector2D{0, 0}
	}
	return Vector2D{v.X / mag, v.Z / mag}
}

// Dot calcula o produto escalar (dot product) entre dois vetores 2D.
// O resultado é um número escalar que representa o quanto os dois vetores estão alinhados:
// - 1 significa que apontam na mesma direção.
// - 0 significa que são perpendiculares (90 graus entre si).
// - -1 significa que apontam em direções opostas.
//
// Isso é útil para calcular ângulos, determinar se um alvo está de frente ou de costas, etc.
func (v Vector2D) Dot(other Vector2D) float64 {
	return v.X*other.X + v.Z*other.Z
}

// Sub retorna o vetor resultante da subtração entre este vetor e outro.
// Útil para calcular o vetor direção entre dois pontos.
func (v Vector2D) Sub(other Vector2D) Vector2D {
	return Vector2D{
		X: v.X - other.X,
		Z: v.Z - other.Z,
	}
}

// RotateVector2D retorna o vetor original rotacionado por um ângulo em radianos.
// A rotação ocorre no plano X-Z no sentido anti-horário.
func RotateVector2D(v Vector2D, angleRad float64) Vector2D {
	cosA := math.Cos(angleRad)
	sinA := math.Sin(angleRad)

	return Vector2D{
		X: v.X*cosA - v.Z*sinA,
		Z: v.X*sinA + v.Z*cosA,
	}
}

func (v Vector2D) Multiply(scalar float64) Vector2D {
	return Vector2D{
		X: v.X * scalar,
		Z: v.Z * scalar,
	}
}

// PerpLeft retorna o vetor perpendicular à esquerda (90 graus CCW)
func (v Vector2D) PerpLeft() Vector2D {
	return Vector2D{
		X: -v.Z,
		Z: v.X,
	}
}

// PerpRight retorna o vetor perpendicular à direita (90 graus CW)
func (v Vector2D) PerpRight() Vector2D {
	return Vector2D{
		X: v.Z,
		Z: -v.X,
	}
}

func (v Vector2D) Scale(scalar float64) Vector2D {
	return v.Multiply(scalar)
}

func (v Vector2D) Length() float64 {
	return v.Magnitude()
}

func NewVector2DFromTo(from, to Position) Vector2D {
	return Vector2D{
		X: to.X - from.X,
		Z: to.Z - from.Z,
	}.Normalize()
}
