package main

import (
	"log"

	"github.com/lunajones/apeiron/lib/model"
	"github.com/lunajones/apeiron/service/grpc"
	"github.com/lunajones/apeiron/service/world"
	"github.com/lunajones/apeiron/service/zone"
)

func main() {
	log.Println("[Main] initializing system...")

	// 🧠 Inicia o servidor gRPC em background
	go grpc.StartGRPCServer("50051")

	// Inicializa habilidades
	model.InitSkills()

	// Inicializa zonas e spawns
	zone.Init()

	// Inicia loop de AI
	world.TickAll()
}
