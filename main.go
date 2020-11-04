package main

import (
	"fmt"
	"math"
	"strconv"
)

type coordenada struct {
	X float32 `json:"X"`
	Y float32 `json:"Y"`
}

type satelite struct {
	Nombre string `json:"Nombre"`
	Posicion coordenada `json:"Posicion"`
}
//Posición de los satélites actualmente en servicio
//● Kenobi: [-500, -200]
//● Skywalker: [100, -100]
//● Sato: [500, 100]

var satelitesFuncionando = []satelite{
	{
		Nombre:   "Kenobi",
		Posicion: coordenada{
			X: -500,
			Y: -200,
		},
	},
	{
		Nombre:   "Skywalker",
		Posicion: coordenada{
			X: 100,
			Y: -100,
		},
	},
	{
		Nombre:   "Sato",
		Posicion: coordenada{
			X: 500,
			Y: 100,
		},
	},
}

func gradosARadianes(deg float64) float64 {
	return deg * (math.Pi / 180.0)
}

func calcularCoordenadas(radio float32, angulo float64, centro coordenada) coordenada {
	x := radio * float32(math.Cos(angulo)) + centro.X
	y := radio * float32(math.Sin(angulo)) + centro.Y
	return coordenada{
		X: x,
		Y: y,
	}
}

func calcularCircunferencia(distancia float32, centroCirculo coordenada) []coordenada {
	dosPi := 2 * math.Pi
	var angulo float64 = 0
	var circunferencia []coordenada
	for angulo < dosPi {
		circunferencia = append(circunferencia, calcularCoordenadas(distancia, angulo, centroCirculo))
		angulo += gradosARadianes(0.1) // se desplaza 0.05°
	}
	return circunferencia
}

func comparar(coo1 coordenada, coo2 coordenada) bool {
	//x2 := math.Pow(float64(coo2.X - coo1.X), 2)
	//y2 := math.Pow(float64(coo2.Y - coo1.Y), 2)
	//d := math.Sqrt(x2 + y2)
	x := math.Abs(float64(coo1.X - coo2.X))
	y := math.Abs(float64(coo1.Y - coo2.Y))
	valorX := x < 5
	valorY := y < 5
	if valorX && valorY {
		return true
	}
	return false
}

func compararCoordenadas(co1 []coordenada, co2 []coordenada) []coordenada {
	var comun []coordenada
	for _, value1 := range co1 {
		for _, value2 := range co2 {
			bandera := comparar(value1, value2)
			if bandera {
				comun = append(comun, value2)
			}
		}
	}
	return comun
}

func promedioCoordenadas(puntos []coordenada) coordenada {
	var x float32
	var y float32
	for _, punto := range puntos {
		x += punto.X
		y += punto.Y
	}
	promedioX := x / float32(len(puntos))
	promedioY := y / float32(len(puntos))
	return coordenada{
		X: promedioX,
		Y: promedioY,
	}
}

func puntoComunCircunferencias(circunferencias ...[]coordenada) (coordenada, string) {
	var puntosComunes = circunferencias[0]
	for i := 1; i < len(circunferencias); i++ {
		puntosComunes = compararCoordenadas(puntosComunes, circunferencias[i])
	}
	for _, comune := range puntosComunes {
		fmt.Println(comune)
	}
	if len(puntosComunes) > 0 {
		coordenadasNave := promedioCoordenadas(puntosComunes)
		return coordenadasNave, ""
	}
	return coordenada{X: 0,	Y: 0}, "Algo salio Mal"
}

func pasarFloat32AString(valor float32) string {
	return fmt.Sprintf("%.1f", valor)
}

func pasarStringAFloat32(valor string) (float32, error) {
	float, err := strconv.ParseFloat(valor, 32)
	if err != nil {
		return 0, err
	}
	return float32(float), err
}

func main() {
	x, y := GetLocation(670.82, 583.1, 761.58)
	fmt.Println(x, y)
}

func foo()  {
	var x float32 = 36.1135843521
	fmt.Println("String:", x)
	fmt.Printf("%T\n", x)
	s := pasarFloat32AString(x)
	fmt.Println("String:", s)
	fmt.Printf("%T\n", s)
	f, _ := pasarStringAFloat32(s)
	fmt.Println("Float32", f)
	fmt.Printf("%T\n", f)
}



//Posición de los satélites actualmente en servicio
//● Kenobi: [-500, -200]
//● Skywalker: [100, -100]
//● Sato: [500, 100]

// input: distancia al emisor tal cual se recibe en cada satélite
// output: las coordenadas ‘x’ e ‘y’ del emisor del mensaje
func GetLocation(distances ...float32) (x, y float32){
	// Se considera que viene la distancia de los satelites en el orden Kenobi, Skywalker y Sato
	// y siempre vienen 3 distancias
	kenobiCirculo := calcularCircunferencia(distances[0], satelitesFuncionando[0].Posicion)
	skywalkerCirculo := calcularCircunferencia(distances[1], satelitesFuncionando[1].Posicion)
	satoCirculo := calcularCircunferencia(distances[2], satelitesFuncionando[2].Posicion)
	coordenadaNave, err := puntoComunCircunferencias(kenobiCirculo, skywalkerCirculo, satoCirculo)
	if err != "" {
		fmt.Println(err)
	}
	return coordenadaNave.X, coordenadaNave.Y
}

// input: el mensaje tal cual es recibido en cada satélite
// output: el mensaje tal cual lo genera el emisor del mensaje
//func GetMessage(messages ...[]string) (msg string) {}
