package domain

import (
	"encoding/json"
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

func TestCents_SerializacionJSON(t *testing.T) {
	monto := Cents(123456789) // $123,456.789 céntimos

	// Serialización determinista a string
	data, err := json.Marshal(monto)
	if err != nil {
		t.Fatalf("Error en serialización: %v", err)
	}

	expectedJSON := `"123456789"`
	if string(data) != expectedJSON {
		t.Errorf("Se esperaba %s, obtenido %s", expectedJSON, string(data))
	}

	// Deserialización exitosa desde string
	var resultado Cents
	err = json.Unmarshal(data, &resultado)
	if err != nil {
		t.Fatalf("Error en deserialización: %v", err)
	}

	if resultado != monto {
		t.Errorf("Se esperaba recuperar %d, obtenido %d", monto, resultado)
	}
}

func TestCents_DeserializacionFormatosNoCanonicos(t *testing.T) {
	// Soporte robusto de deserialización ante entradas tipo número nativo en JSON.
	rawNumJSON := []byte(`9999`)
	var res Cents
	if err := json.Unmarshal(rawNumJSON, &res); err != nil {
		t.Fatalf("Error deserializando número crudo: %v", err)
	}
	if res != 9999 {
		t.Errorf("Se esperaba 9999, obtenido %d", res)
	}

	badJSON := []byte(`"texto_invalido"`)
	if err := json.Unmarshal(badJSON, &res); err == nil {
		t.Error("Se esperaba fallo por cadena no convertible a número")
	}
}

func TestCents_FormatosProhibidosNUM02(t *testing.T) {
	// Casos de prueba que violan explícitamente la canonicidad del baseline normativo
	prohibidos := []string{
		`"05"`,    // Ceros iniciales positivos
		`"-05"`,   // Ceros iniciales negativos
		`"-0"`,    // Negativo nulo explícitamente prohibido por arquitectura
		`"12.34"`, // Fracciones/Decimales en cadena
		`"1e3"`,   // Exponentes
		`" 100"`,  // Espacios internos iniciales
		`"100 "`,  // Espacios internos finales
		`""`,      // Vacío
	}
	for _, raw := range prohibidos {
		var c Cents
		if err := json.Unmarshal([]byte(raw), &c); err == nil {
			t.Errorf("Se esperaba fallo de validación canónica NUM-02 para la entrada: %s", raw)
		}
	}
}
