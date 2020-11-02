package main

import (
	"fmt"
	"math"
)

type coordenada struct {
	X float32 `json:"X"`
	Y float32 `json:"Y"`
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

func calcularCircunferencia(distancia float32) []coordenada {
	dosPi := 2 * math.Pi
	var angulo float64 = 0
	centro := coordenada{
		X: 10,
		Y: -10,
	}
	var coord coordenada
	var circunferencia []coordenada
	for angulo < dosPi {
		coord = calcularCoordenadas(distancia, angulo, centro)
		circunferencia = append(circunferencia, coord)
		angulo += gradosARadianes(0.05) // se desplaza 0.05°
	}
	return circunferencia
}

func main() {
	circulo := calcularCircunferencia(20.5)
	for i, value := range circulo {
		fmt.Println("========= ITERACION:", i, "=========")
		fmt.Println("X:", value.X)
		fmt.Println("Y:", value.Y)
	}
}



//Posición de los satélites actualmente en servicio
//● Kenobi: [-500, -200]
//● Skywalker: [100, -100]
//● Sato: [500, 100]

// input: distancia al emisor tal cual se recibe en cada satélite
// output: las coordenadas ‘x’ e ‘y’ del emisor del mensaje
//func GetLocation(distances ...float32) (x, y float32){
//
//}

// input: el mensaje tal cual es recibido en cada satélite
// output: el mensaje tal cual lo genera el emisor del mensaje
//func GetMessage(messages ...[]string) (msg string) {}
