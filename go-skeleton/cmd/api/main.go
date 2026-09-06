package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/klik/fcos-kernel/config"
	deliveryHttp "github.com/klik/fcos-kernel/internal/handler/http"
	"github.com/klik/fcos-kernel/internal/service"
)

func main() {
	fmt.Println("Iniciando Contable Fix by KLIK...")

	// 1. Cargar configuración base
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	// 2. Inicializar repositorios (implementaciones de persistencia)
	// TODO: Inyectar repositorios reales (Postgres, MySQL o In-Memory)

	// 3. Inicializar servicios de casos de uso
	accountSvc := service.NewAccountService(nil)
	journalSvc := service.NewJournalService(nil, nil, nil, nil, nil)

	// 4. Inicializar handlers HTTP
	accountHandler := deliveryHttp.NewAccountHandler(accountSvc)
	journalHandler := deliveryHttp.NewJournalHandler(journalSvc)

	// 5. Configurar router
	router := deliveryHttp.NewRouter(deliveryHttp.RouterConfig{
		AccountHandler: accountHandler,
		JournalHandler: journalHandler,
	})

	// 6. Arrancar servidor HTTP
	port := cfg.Server.Port
	if port == 0 {
		port = 8080
	}
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Servidor Contable Fix escuchando en %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Fallo en el servidor HTTP: %v", err)
	}
}
