package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math"
	"net/http"
)

const DESPLAZAMIENTO_ANGULAR = 0.05
const ERROR_PRECISION = 5

type errorResponse struct {
	Message string `json:"Message"`
}

type coordenada struct {
	X float32 `json:"X"`
	Y float32 `json:"Y"`
}

type satelite struct {
	Nombre string `json:"Nombre"`
	Posicion coordenada `json:"Posicion"`
}

type sateliteRequest struct {
	Name string `json:"Name"`
	Distance float32 `json:"Distance"`
	Message []string `json:"Message"`
}

type infoNave struct {
	Satellites []sateliteRequest `json:"Satellites"`
}


type createInfoNave struct {
	Distance float32 `json:"Distance"`
	Message []string `json:"Message"`
}

var infoSatelites []sateliteRequest

type responseInfoNave struct {
	Position coordenada `json:"Position"`
	Message string `json:"Message"`
}

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
		angulo += gradosARadianes(DESPLAZAMIENTO_ANGULAR) // se desplaza 0.05°
	}
	return circunferencia
}

func comparar(coo1 coordenada, coo2 coordenada) bool {
	x := math.Abs(float64(coo1.X - coo2.X))
	y := math.Abs(float64(coo1.Y - coo2.Y))
	valorX := x < ERROR_PRECISION
	valorY := y < ERROR_PRECISION
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

func puntoComunCircunferencias(circunferencias ...[]coordenada) (coord coordenada, err error) {
	var puntosComunes = circunferencias[0]
	for i := 1; i < len(circunferencias); i++ {
		puntosComunes = compararCoordenadas(puntosComunes, circunferencias[i])
	}
	if len(puntosComunes) > 0 {
		coordenadasNave := promedioCoordenadas(puntosComunes)
		return coordenadasNave, nil
	}
	return coordenada{X: 0,	Y: 0}, errors.New("algo salio mal")
}

func validarNombreSatelites(name string) bool {
	for _, s := range satelitesFuncionando {
		if s.Nombre == name {
			return true
		}
	}
	return false
}


func infoNaveHandler(writer http.ResponseWriter, request *http.Request) {
	var infoNave infoNave
	var resInfo responseInfoNave
	var errorsRes []errorResponse
	reqBody, err := ioutil.ReadAll(request.Body)
	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		_, _ = fmt.Fprintf(writer, "Ocurrio un problema")
	}
	err = json.Unmarshal(reqBody, &infoNave)
	if err != nil {
		_, _ = fmt.Fprintf(writer, "Ocurrio un problema")
	}
	if len(infoNave.Satellites) == 3 {
		for _, satellite := range infoNave.Satellites {
			if !validarNombreSatelites(satellite.Name) {
				e := errorResponse{Message: "No se encontro satelite en funcionamiento con nombre:" + satellite.Name}
				errorsRes = append(errorsRes, e)
			}
		}
		if len(errorsRes) != 0 {
			writer.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(writer).Encode(errorsRes)
		} else {
			var distancias []float32
			var mensajes [][]string
			ordenSatelites := []string{"Kenobi", "Skywalker", "Sato"}
			for _, nameSate := range ordenSatelites {
				for _, satellite := range infoNave.Satellites {
					if nameSate == satellite.Name {
						distancias = append(distancias, satellite.Distance)
						mensajes = append(mensajes, satellite.Message)
					}
				}
			}
			msg, erroMensajes := GetMessage(mensajes...)
			x, y, errorLocation := GetLocation(distancias...)
			if erroMensajes != nil || errorLocation != nil {
				writer.WriteHeader(http.StatusNotFound)
				//_ = json.NewEncoder(writer).Encode(error.Error())
			} else {
				resInfo = responseInfoNave{
					Position: coordenada{
						X: x,
						Y: y,
					},
					Message:  msg,
				}
				writer.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(writer).Encode(resInfo)
			}
		}
	} else {
		writer.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(writer).Encode(errorResponse{Message: "Pasar una cantidad de 3 satelites exactamente"})
	}
}

func createInfoNaveHandler(writer http.ResponseWriter, request *http.Request) {
	var newInfoNave createInfoNave
	var errorsResponse []errorResponse
	reqBody, err := ioutil.ReadAll(request.Body)
	vars := mux.Vars(request)
	name, ok := vars["satellite_name"]
	if !ok {
		_, _ = fmt.Fprintf(writer, "Ocurrio un problema")
	}
	writer.Header().Set("Content-Type", "application/json")

	if err != nil {
		_, _ = fmt.Fprintf(writer, "Ocurrio un problema")
	}
	err = json.Unmarshal(reqBody, &newInfoNave)
	if err != nil {
		_, _ = fmt.Fprintf(writer, "Ocurrio un problema")
	}

	if !validarNombreSatelites(name) {
		e := errorResponse{Message: "No se encontro satelite en funcionamiento con nombre:" + name}
		errorsResponse = append(errorsResponse, e)
	}
	for _, sate := range infoSatelites {
		if sate.Name == name {
			e := errorResponse{Message: "Satelite " + name + " ya existe"}
			errorsResponse = append(errorsResponse, e)
		}
	}
	if len(infoSatelites) == 3 {
		e := errorResponse{Message: "Solo se pueden crear 3 satelites"}
		errorsResponse = append(errorsResponse, e)
	}
	if len(errorsResponse) != 0 {
		writer.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(writer).Encode(errorsResponse)
	} else {
		info := sateliteRequest{
			Name:     name,
			Distance: newInfoNave.Distance,
			Message:  newInfoNave.Message,
		}
		infoSatelites = append(infoSatelites, info)

		writer.WriteHeader(http.StatusCreated)

	}
}

func readInfoNaveHandler(writer http.ResponseWriter, request *http.Request) {
	var resInfo responseInfoNave
	writer.Header().Set("Content-Type", "application/json")

	if len(infoSatelites) == 3 {
		var distancias []float32
		var mensajes [][]string
		ordenSatelites := []string{"Kenobi", "Skywalker", "Sato"}
		for _, nameSate := range ordenSatelites {
			for _, satellite := range infoSatelites {
				if nameSate == satellite.Name {
					distancias = append(distancias, satellite.Distance)
					mensajes = append(mensajes, satellite.Message)
				}
			}
		}
		msg, erroMensajes := GetMessage(mensajes...)
		x, y, errorLocation := GetLocation(distancias...)
		if erroMensajes != nil || errorLocation != nil {
			writer.WriteHeader(http.StatusNotFound)
			//_ = json.NewEncoder(writer).Encode(error.Error())
		} else {
			resInfo = responseInfoNave{
				Position: coordenada{
					X: x,
					Y: y,
				},
				Message:  msg,
			}
			writer.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(writer).Encode(resInfo)
		}
	} else {
		writer.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	//x, y := GetLocation(670.82, 583.1, 761.58)
	//fmt.Println(x, y)

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/topsecret", infoNaveHandler).Methods("POST")
	router.HandleFunc("/topsecret_split/{satellite_name}", createInfoNaveHandler).Methods("POST")
	router.HandleFunc("/topsecret_split", readInfoNaveHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8090", router))
}

//Posición de los satélites actualmente en servicio
//● Kenobi: [-500, -200]
//● Skywalker: [100, -100]
//● Sato: [500, 100]

// input: distancia al emisor tal cual se recibe en cada satélite
// output: las coordenadas ‘x’ e ‘y’ del emisor del mensaje
func GetLocation(distances ...float32) (x, y float32, err error){
	// Se considera que viene la distancia de los satelites en el orden Kenobi, Skywalker y Sato
	// y siempre vienen 3 distancias
	kenobiCirculo := calcularCircunferencia(distances[0], satelitesFuncionando[0].Posicion)
	skywalkerCirculo := calcularCircunferencia(distances[1], satelitesFuncionando[1].Posicion)
	satoCirculo := calcularCircunferencia(distances[2], satelitesFuncionando[2].Posicion)
	coordenadaNave, err := puntoComunCircunferencias(kenobiCirculo, skywalkerCirculo, satoCirculo)
	if err != nil {
		return 0, 0, err
	}
	return coordenadaNave.X, coordenadaNave.Y, nil
}


func descifrar(palabras ...string) string {
	indecifrable := false
	var palabraRes = ""
	for _, p := range palabras {
		if p != "" {
			// p no viene vacia
			if palabraRes != ""{
				// palabraRes ya tiene algo
				if palabraRes != p{
					// si las palabras no coinciden
					indecifrable = true
					palabraRes = ""
				}
			} else {
				// palabraRes no tiene nada
				if !indecifrable {
					// Aun no hubo inconsistencia
					palabraRes = p
				}
			}
		}
	}
	return palabraRes
}

// input: el mensaje tal cual es recibido en cada satélite
// output: el mensaje tal cual lo genera el emisor del mensaje
func GetMessage(messages ...[]string) (msg string, err error) {
	mismoTam := len(messages[0]) == len(messages[1]) && len(messages[1]) == len(messages[2])
	tam := len(messages[2])
	if mismoTam {
		for i := 0; i < tam; i++ {
			if (i + 1) == tam {
				msg += descifrar(messages[0][i], messages[1][i], messages[2][i])
			} else {
				msg += descifrar(messages[0][i], messages[1][i], messages[2][i]) + " "
			}
		}
		return msg, nil
	} else {
		return "", errors.New("no se puede descifrar el mensaje")
	}
}


