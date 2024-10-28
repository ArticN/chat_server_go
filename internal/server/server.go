package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"
)

type chanClient chan<- string

type User struct {
	chn      chanClient
	nickname string
}

type PrivateMessage struct {
	src string
	dst string
	msg string
}

var (
	entering        = make(chan User)
	leaving         = make(chan User)
	messages        = make(chan string)
	privateMessages = make(chan PrivateMessage)
	activeUsers     = make([]string, 0)
)

func StartServer() {
	fmt.Println("Starting server...")

	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	time.Sleep(time.Second)
	fmt.Println("Server started!")

	go BroadcasterHandler()
	go PrivateHandler()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go HandleConn(conn)
	}
}

func BroadcasterHandler() {
	clients := make(map[chanClient]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}
			log.Println(msg)

		case usr := <-entering:
			clients[usr.chn] = true
			activeUsers = append(activeUsers, usr.nickname)
			log.Printf("User %q connected. Active users: %v\n", usr.nickname, activeUsers)

		case usr := <-leaving:
			delete(clients, usr.chn)
			activeUsers = removeUser(activeUsers, usr.nickname)
			log.Printf("User %q disconnected. Active users: %v\n", usr.nickname, activeUsers)
		}
	}
}

func PrivateHandler() {
	clients := make(map[string]chanClient)

	for {
		select {
		case privateMessage := <-privateMessages:
			if chnSrc, ok := clients[privateMessage.dst[1:]]; ok {
				message := fmt.Sprintf("@%v disse em privado: %v", privateMessage.src, privateMessage.msg)
				log.Print(message)
				chnSrc <- message
			}

		case usr := <-entering:
			clients[usr.nickname] = usr.chn

		case usr := <-leaving:
			delete(clients, usr.nickname)
		}
	}
}

func removeUser(users []string, nickname string) []string {
	for i, user := range users {
		if user == nickname {
			return append(users[:i], users[i+1:]...)
		}
	}
	return users
}

func GetActiveUsers() []string {
	return activeUsers
}

func HandleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 20)
	size, err := conn.Read(buf)
	if err != nil {
		log.Println("Erro ao ler o nickname:", err)
		return
	}

	nickname := strings.TrimSpace(string(buf[:size-1]))
	if nickname == "" {
		log.Println("Nickname não pode ser vazio")
		return
	}
	ch := make(chan string)
	usr := User{ch, nickname}

	go clientWriter(conn, ch)

	ch <- fmt.Sprintf("Conectado como %q", usr.nickname)
	messages <- fmt.Sprintf("Usuário @%v entrou", usr.nickname)
	entering <- usr // Send user to entering channel

	input := bufio.NewScanner(conn)
	for input.Scan() {
		rawTxt := input.Text()
		if handleCommand(rawTxt, usr) {
			leaving <- usr // Send user to leaving channel
			close(ch)
			break
		}
	}
	messages <- fmt.Sprintf("Usuário @%v saiu", usr.nickname)
}

func handleCommand(rawTxt string, usr User) bool {
	re := regexp.MustCompile(`^(\\\w+)(?:\s+(@\w+))?(?:\s+(.*))?$`)
	matches := re.FindStringSubmatch(rawTxt)

	if len(matches) == 0 {
		log.Printf("ERROR:comando invalido %q", rawTxt)
		return false
	}

	cmd, privateUser, msg := matches[1], matches[2], matches[3]

	switch cmd {
	case "\\msg":
		if privateUser != "" {
			privateMessages <- PrivateMessage{usr.nickname, privateUser, msg}
		} else {
			messages <- fmt.Sprintf("@%v disse: %v", usr.nickname, msg)
		}
	case "\\users":
		activeUsers := GetActiveUsers()
		usr.chn <- fmt.Sprintf("Usuarios ativos: %v", activeUsers)
	case "\\exit":
		leaving <- usr
		return true
	case "\\changenickname":
		if msg != "" {
			oldNick := usr.nickname
			usr.nickname = msg
			messages <- fmt.Sprintf("Usuario @%v agora é @%v", oldNick, msg)
		} else {
			log.Println("ERROR: comando invalido")
		}
	default:
		usr.chn <- fmt.Sprintf("ERROR: Invalid command: %v", cmd)
	}
	return false
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		if _, err := fmt.Fprintln(conn, msg); err != nil {
			log.Println("Erro ao enviar mensagem para o cliente:", err)
			return
		}
	}
}
