package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"biblioteca_go-whatsapp/lib/whatsmeow"

	"lib/whatsmeow/store/sqlstore"
	"lib/whatsmeow/types/events"
	waLog "lib/whatsmeow/util/log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
)

func main() {
	dbLog := waLog.Stdout("Client", "DATABASE: ", true)
	container, err := sqlstore.New("sqlite3", "file:whatsApp_api.db?_foreign_keys=on", dbLog)
	if err != nil {
		log.Fatalf("Erro ao criar armazenamento: %v", err)
	}

	clientLog := waLog.Stdout("Client", "CLIENT: ", true)
	device, err := container.GetFirstDevice()
	if err != nil {
		log.Fatalf("Erro ao obter dispositivo: %v", err)
	}

	client := whatsmeow.NewClient(device, clientLog)

	// Adiciona manipulador para o QR code
	client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			fmt.Printf("Nova mensagem recebida: %+v\n", v.Message.GetConversation())
		case *events.QR:
			// Exibe o QR code no terminal
			fmt.Println("QR Code recebido:")
			if len(v.Codes) > 0 {
				qrCode := v.Codes[0]
				qrterminal.GenerateHalfBlock(qrCode, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Nenhum código QR disponível.")
			}
		default:
			fmt.Printf("Evento desconhecido: %+v\n", v)
		}
	})

	// Conecta ao WhatsApp
	err = client.Connect()
	if err != nil {
		log.Fatalf("Erro ao conectar: %v", err)
	}

	fmt.Println("Aguardando autenticação...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
	fmt.Println("Cliente desconectado. Programa encerrado.")
}
