package position

import "math"

// Vector3D representa um vetor 3D genérico para velocidade, aceleração etc.
type Vector3D struct {
	X float64 // plano X
	Y float64 // altura Y
	Z float64 // plano Z
}

// Add soma dois vetores
func (v Vector3D) Add(other Vector3D) Vector3D {
	return Vector3D{
		X: v.X + other.X,
		Y: v.Y + other.Y,
		Z: v.Z + other.Z,
	}
}

// Sub subtrai dois vetores
func (v Vector3D) Sub(other Vector3D) Vector3D {
	return Vector3D{
		X: v.X - other.X,
		Y: v.Y - other.Y,
		Z: v.Z - other.Z,
	}
}

// Scale multiplica o vetor por um escalar
func (v Vector3D) Scale(scalar float64) Vector3D {
	return Vector3D{
		X: v.X * scalar,
		Y: v.Y * scalar,
		Z: v.Z * scalar,
	}
}

// Magnitude calcula o comprimento do vetor
func (v Vector3D) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Normalize retorna um vetor unitário na mesma direção
func (v Vector3D) Normalize() Vector3D {
	mag := v.Magnitude()
	if mag == 0 {
		return Vector3D{}
	}
	return v.Scale(1 / mag)
}

// Zero zera o vetor
func (v *Vector3D) Zero() {
	v.X = 0
	v.Y = 0
	v.Z = 0
}

// NewVector3DFromTo cria um vetor do ponto A até o ponto B
func NewVector3DFromTo(a, b Position) Vector3D {
	return Vector3D{
		X: b.X - a.X,
		Y: b.Y - a.Y,
		Z: b.Z - a.Z,
	}
}

func (v Vector3D) Multiply(scalar float64) Vector3D {
	return Vector3D{
		X: v.X * scalar,
		Y: v.Y * scalar,
		Z: v.Z * scalar,
	}
}

// ToVector2D converte um vetor 3D em 2D, descartando a componente Y
func (v Vector3D) ToVector2D() Vector2D {
	return Vector2D{
		X: v.X,
		Z: v.Z, // Descartando a componente Y
	}
}

func (v Vector3D) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}
