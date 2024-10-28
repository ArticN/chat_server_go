package bot

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func InitializeConnection() net.Conn {
	connection, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Fatal("conexao falhou: ", err)
	}

	log.Println("conexao estabelecida")
	connection.Write([]byte("BOT\n"))
	return connection
}

func HandleMessages(connection net.Conn, signal chan struct{}, responseHandler func(msg string) string) {
	scanner := bufio.NewScanner(connection)

	for scanner.Scan() {
		serverMessage := scanner.Text()
		log.Printf("Mensagem recebida do servidor: %s", serverMessage)
		handleIncomingMessage(connection)
	}

	if err := scanner.Err(); err != nil {
		log.Println("ERROR: ler da conexao do server:", err)
	}

	log.Println("conexao perdida")
	signal <- struct{}{}
}

func handleIncomingMessage(connection net.Conn) {
	scanner := bufio.NewScanner(connection)
	for scanner.Scan() {
		text := scanner.Text()

		if strings.Contains(text, "disse em privado:") {
			fmt.Println("received message:", text)

			privateMsgIndex := strings.Index(text, "disse em privado:") + len("disse em privado:")
			originalMessage := strings.TrimSpace(text[privateMsgIndex:])

			reversedMessage := reverseString(originalMessage)
			connection.Write([]byte("\\msg " + reversedMessage + "\n"))
		}
	}
}

func reverseString(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func StartBot() {
	connection := InitializeConnection()
	defer connection.Close()

	signal := make(chan struct{})

	go HandleMessages(connection, signal, func(msg string) string {
		if strings.Contains(strings.ToLower(msg), "hello") {
			return "Hello there!"
		}
		return ""
	})

	<-signal
	log.Println("Bot finalizado.")
}
