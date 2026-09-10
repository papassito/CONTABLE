package domain

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
)

var (
	ErrOverflowCents = errors.New("FCOS_ERR_MONEY: El monto monetario excede los límites seguros de int64")
	ErrInvalidFormat = errors.New("FCOS_ERR_MONEY: El formato de cadena monetaria no es canónico ni válido")
)

// Cents representa dinero exacto en centavos de MXN para prevenir errores de coma flotante.
// Se utiliza un int64 subyacente. $1.00 se almacena como 100 centavos.
type Cents int64

// Add realiza una suma de valores monetarios con protección contra desbordamientos.
func (c Cents) Add(o Cents) (Cents, error) {
	if o > 0 && c > math.MaxInt64-o {
		return 0, ErrOverflowCents
	}
	if o < 0 && c < math.MinInt64-o {
		return 0, ErrOverflowCents
	}
	return c + o, nil
}

// Sub realiza una resta con protección contra desbordamientos.
func (c Cents) Sub(o Cents) (Cents, error) {
	if o < 0 && c > math.MaxInt64+o {
		return 0, ErrOverflowCents
	}
	if o > 0 && c < math.MinInt64+o {
		return 0, ErrOverflowCents
	}
	return c - o, nil
}

// Neg invierte el signo aritmético de los centavos con protección contra desbordamiento de extremos.
func (c Cents) Neg() (Cents, error) {
	if c == math.MinInt64 {
		return 0, ErrOverflowCents
	}
	return -c, nil
}

// MarshalJSON implementa la serialización determinista.
// Todos los montos se transmiten cruzando el IPC/API como cadenas decimales exactas.
func (c Cents) MarshalJSON() ([]byte, error) {
	strVal := strconv.FormatInt(int64(c), 10)
	return []byte("\"" + strVal + "\""), nil
}

// UnmarshalJSON deserializa cadenas de texto decimales a enteros de centavos de manera robusta.
func (c *Cents) UnmarshalJSON(data []byte) error {
	var rawString string
	if err := json.Unmarshal(data, &rawString); err != nil {
		return ErrInvalidFormat
	}

	if len(rawString) == 0 {
		return ErrInvalidFormat
	}

	// Restricción de ARCHITECTURE.md Sección 4: Se rechaza -0 de forma explícita
	if rawString == "-0" {
		return ErrInvalidFormat
	}

	// Restricción de REQUIREMENTS.md (NUM-02): Rechazar ceros iniciales (e.g. "05", "-05")
	if len(rawString) > 1 {
		if rawString[0] == '0' {
			return ErrInvalidFormat
		}
		if rawString[0] == '-' && rawString[1] == '0' {
			return ErrInvalidFormat
		}
	}

	// Asegurar que contiene únicamente dígitos y un signo menos opcional al inicio (sin puntos, e ni espacios intermedios)
	start := 0
	if rawString[0] == '-' {
		start = 1
		if len(rawString) == 1 {
			return ErrInvalidFormat
		}
	}
	for i := start; i < len(rawString); i++ {
		if rawString[i] < '0' || rawString[i] > '9' {
			return ErrInvalidFormat
		}
	}

	val, err := strconv.ParseInt(rawString, 10, 64)
	if err != nil {
		return ErrInvalidFormat
	}

	*c = Cents(val)
	return nil
}
