package common

import (
	"fmt"
)

type Bet struct {
	Agencia    string
	Nombre     string
	Apellido   string
	Documento  string
	Nacimiento string
	Numero     string
}

func (b Bet) getBetSize() uint32 {
	return uint32(len(b.getBetSerialized()))
}

func (b Bet) getBetSerialized() string {
	log.Infof("bet con %s%s%s%s%s%s", b.Agencia, b.Nombre, b.Apellido, b.Documento, b.Nacimiento, b.Numero)
	return fmt.Sprintf("%s%s%s%s%s%s", b.Agencia, b.Nombre, b.Apellido, b.Documento, b.Nacimiento, b.Numero)
}