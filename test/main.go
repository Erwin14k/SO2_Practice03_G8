package main

import (
	"fmt"
	"time"
)

func main() {
	var memory []byte

	for {
		memory = append(memory, make([]byte, 1024*1024)...) // Añadir 1MB de memoria en cada iteración
		time.Sleep(time.Millisecond * 100)                  // Esperar 10 milisegundos en cada iteración
		fmt.Println("Memoria asignada:", len(memory)/1024/1024, "MB")
	}
}
