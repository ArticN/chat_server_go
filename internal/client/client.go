package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

func ConnectToServer() (net.Conn, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Nickname: ")
	apelido, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("ERROR:read-nickname: %w", err)
	}

	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		return nil, fmt.Errorf("ERROR:connect-to-server: %w", err)
	}
	log.Println("Conectado ao servidor!")

	_, err = conn.Write([]byte(strings.TrimSpace(apelido) + "\n"))
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("ERROR:send-nickname: %w", err)
	}

	return conn, nil
}

func KeepAlive(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("Conexão com o servidor encerrada.")
			} else {
				log.Println("Erro ao manter a conexão ativa:", err)
			}
			break
		}
		log.Println("Mensagem do servidor:", strings.TrimSpace(message))
	}
}

func HandleInput(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	reader := bufio.NewReader(os.Stdin)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Erro ao ler entrada:", err)
			break
		}
		_, err = conn.Write([]byte(strings.TrimSpace(input) + "\n"))
		if err != nil {
			log.Println("Erro ao enviar mensagem:", err)
			break
		}
		if strings.TrimSpace(input) == "\\exit" {
			log.Println("Saindo do chat...")
			break
		}
	}
}

func StartClient() {
	conn, err := ConnectToServer()
	if err != nil {
		log.Fatal("Erro ao conectar ao servidor:", err)
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go KeepAlive(conn, &wg)
	go HandleInput(conn, &wg)

	wg.Wait()
	log.Println("Cliente encerrado.")
}
