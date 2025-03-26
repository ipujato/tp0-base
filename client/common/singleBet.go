package common

import (
	"fmt"
	"strings"
)

type Bet struct {
	Agencia    string
	Nombre     string
	Apellido   string
	Documento  string
	Nacimiento string
	Numero     string
}

func betFromString(serialized string, agency string) (Bet, error) {
	// var b Bet
	// _, err := fmt.Sscanf(serialized, "%s|%s|%s|%s|%s|%s", &b.Agencia, &b.Nombre, &b.Apellido, &b.Documento, &b.Nacimiento, &b.Numero)
	// if err != nil {
	// 	return Bet{}, err
	// }
	// return b, nil
	splitedString := strings.Split(serialized, ",")
	if len(splitedString) != 5 {
		return Bet{}, fmt.Errorf("invalid bet format")
	}
	bet := Bet {
		Agencia: agency,
		Nombre: splitedString[0],
		Apellido: splitedString[1],
		Documento: splitedString[2],
		Nacimiento: splitedString[3],
		Numero: splitedString[4],
	}
	return bet, nil
}

func (b Bet) getBetSerialized() string {
	return fmt.Sprintf("%s|%s|%s|%s|%s|%s\n", b.Agencia, b.Nombre, b.Apellido, b.Documento, b.Nacimiento, b.Numero)
}