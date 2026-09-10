package domain

import (
	"encoding/json"
	"errors"
	"math"
	"testing"
)

func TestCents_AritmeticaSegura(t *testing.T) {
	a := Cents(1000) // $10.00
	b := Cents(550)  // $5.50

	// Suma normal
	sum, err := a.Add(b)
	if err != nil {
		t.Fatalf("No se esperaba error en suma: %v", err)
	}
	if sum != 1550 {
		t.Errorf("Se esperaba 1550, obtenido %d", sum)
	}

	// Resta normal
	sub, err := a.Sub(b)
	if err != nil {
		t.Fatalf("No se esperaba error en resta: %v", err)
	}
	if sub != 450 {
		t.Errorf("Se esperaba 450, obtenido %d", sub)
	}
}

func TestCents_Desbordamiento(t *testing.T) {
	max := Cents(math.MaxInt64)
	min := Cents(math.MinInt64)

	_, err := max.Add(1)
	if err == nil {
		t.Error("Se esperaba error por desbordamiento superior en suma")
	}

	_, err = min.Sub(1)
	if err == nil {
		t.Error("Se esperaba error por desbordamiento inferior en resta")
	}

	_, err = min.Neg()
	if err == nil {
		t.Error("Se esperaba error por desbordamiento al negar MinInt64")
	}
}

func TestCents_UnmarshalJSON_FormatosCanonicosValidos(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Cents
	}{
		{
			name:     "cero",
			input:    `"0"`,
			expected: 0,
		},
		{
			name:     "positivo",
			input:    `"12345"`,
			expected: 12345,
		},
		{
			name:     "negativo",
			input:    `"-12345"`,
			expected: -12345,
		},
		{
			name:     "max int64",
			input:    `"9223372036854775807"`,
			expected: Cents(math.MaxInt64),
		},
		{
			name:     "min int64",
			input:    `"-9223372036854775808"`,
			expected: Cents(math.MinInt64),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Cents

			if err := json.Unmarshal([]byte(tt.input), &got); err != nil {
				t.Fatalf("entrada válida %s rechazada: %v", tt.input, err)
			}

			if got != tt.expected {
				t.Fatalf("esperado %d, obtenido %d", tt.expected, got)
			}
		})
	}
}

func TestCents_UnmarshalJSON_RechazaFormatosNoCanonicos(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"numero JSON crudo", `12345`},
		{"numero JSON negativo crudo", `-12345`},
		{"cero negativo", `"-0"`},
		{"cero inicial", `"05"`},
		{"cero inicial negativo", `"-05"`},
		{"signo positivo", `"+5"`},
		{"fraccion", `"12.34"`},
		{"exponente minuscula", `"1e3"`},
		{"exponente mayuscula", `"1E3"`},
		{"espacio inicial", `" 12345"`},
		{"espacio final", `"12345 "`},
		{"espacios ambos", `" 12345 "`},
		{"tabulador", `"\t12345"`},
		{"salto de linea", `"12345\n"`},
		{"vacio", `""`},
		{"solo menos", `"-"`},
		{"letras", `"abc"`},
		{"separador coma", `"1,000"`},
		{"separador underscore", `"1_000"`},
		{"boolean", `true`},
		{"null", `null`},
		{"array", `["123"]`},
		{"objeto", `{"value":"123"}`},

		// Fuera de int64
		{"overflow positivo", `"9223372036854775808"`},
		{"overflow negativo", `"-9223372036854775809"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Cents

			err := json.Unmarshal([]byte(tt.input), &got)

			if !errors.Is(err, ErrInvalidFormat) {
				t.Fatalf(
					"entrada %s: esperado ErrInvalidFormat, obtenido %v",
					tt.input,
					err,
				)
			}
		})
	}
}

func TestCents_MarshalJSON_FormatoCanonico(t *testing.T) {
	tests := []struct {
		value    Cents
		expected string
	}{
		{0, `"0"`},
		{1, `"1"`},
		{-1, `"-1"`},
		{12345, `"12345"`},
		{-12345, `"-12345"`},
		{Cents(math.MaxInt64), `"9223372036854775807"`},
		{Cents(math.MinInt64), `"-9223372036854775808"`},
	}

	for _, tt := range tests {
		data, err := json.Marshal(tt.value)
		if err != nil {
			t.Fatalf("Marshal(%d): %v", tt.value, err)
		}

		if string(data) != tt.expected {
			t.Fatalf(
				"Marshal(%d): esperado %s, obtenido %s",
				tt.value,
				tt.expected,
				string(data),
			)
		}
	}
}

func TestCents_JSON_RoundTripExacto(t *testing.T) {
	values := []Cents{
		0,
		1,
		-1,
		100,
		-100,
		123456789,
		Cents(math.MaxInt64),
		Cents(math.MinInt64),
	}

	for _, original := range values {
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Marshal(%d): %v", original, err)
		}

		var recovered Cents

		if err := json.Unmarshal(data, &recovered); err != nil {
			t.Fatalf("Unmarshal(%d): %v", original, err)
		}

		if recovered != original {
			t.Fatalf(
				"round-trip alterado: original=%d recovered=%d",
				original,
				recovered,
			)
		}
	}
}
