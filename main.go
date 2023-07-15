package main

import (
	"github.com/SCP-FoundationHQ/SCPFOUNDATION/src"
)

func main() {
	src.Gin.Run(":"+src.Envars.Port)
}
