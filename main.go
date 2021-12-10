package main

import (
	"encoding/gob" // pra arquivos binários
	"fmt"
	"log"
	"os"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
)

type waHandler struct {
	c         *whatsapp.Conn
	startTime uint64
}

func (wh *waHandler) HandleTextMessage(message whatsapp.TextMessage) {
	fmt.Println(wh.startTime)
	if message.Info.Timestamp < wh.startTime {
		return
	}

	fmt.Printf("%v %v %v %v\n\t%v\n", message.Info.Timestamp, message.Info.Id, message.Info.RemoteJid, message.ContextInfo.QuotedMessageID, message.Text)
}

func (h *waHandler) HandleError(err error) {

	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := h.c.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	} else {
		log.Printf("error occoured: %v\n", err)
	}
}

func onInit(whatConn *whatsapp.Conn) error {
	var sessionError error = fmt.Errorf("no session")

	// tenta encontrar uma sessão armazenada no sistema
	session, sessionError := readSessionFromFileSystem()
	if sessionError == nil {
		// Tenta recuperar sessão salva
		session, sessionError = whatConn.RestoreWithSession(session)
		if sessionError != nil {
			log.Printf("erro restaurando sessão: %v\n", sessionError)
		}

	} else {
		log.Printf("Nenhuma sessão encontrada nos arquivos do sistema: %v\n", sessionError)
	}

	if sessionError != nil {
		// Faz o login regular
		session, sessionError = login(whatConn)
		if sessionError != nil {
			log.Printf("erro durante o login: %v\n", sessionError)
		}
	}

	// salvando sessão
	sessionError = writeSessionToFileSystem(session)
	if sessionError != nil {
		return fmt.Errorf("erro salvando sessão: %v\n", sessionError)
	}
	return nil
}

func login(whatConn *whatsapp.Conn) (whatsapp.Session, error) {
	qrCode := make(chan string)
	go func() {
		terminal := qrcodeTerminal.New2(qrcodeTerminal.ConsoleColors.BrightBlack, qrcodeTerminal.ConsoleColors.BrightWhite, qrcodeTerminal.QRCodeRecoveryLevels.Low)
		terminal.Get([]byte(<-qrCode)).Print()
	}()

	return whatConn.Login(qrCode)
}

func writeSessionToFileSystem(session whatsapp.Session) error {
	// Cria arquivo temporário no sistema
	file, err := os.Create(os.TempDir() + "waSession.gob")
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}

	return nil
}

func readSessionFromFileSystem() (whatsapp.Session, error) {
	session := whatsapp.Session{}

	// Pega arquivo temporário do sistema
	file, err := os.Open(os.TempDir() + "waSession.gob")
	if err != nil {
		return session, err
	}

	// Executado após o fim da função
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}

	return session, nil
}

func main() {
	whatConn, err := whatsapp.NewConn(10 * time.Second) //10secs of timeout
	if err != nil {
		panic(err)
	}
	// Para o erro no retorno ummarshall
	whatConn.SetClientVersion(3, 2123, 7)

	whatConn.AddHandler(&waHandler{whatConn, uint64(time.Now().Unix())})

	onInit(whatConn)

	/* Envia mensagem */
	/* var numberWhat int

	fmt.Print("Número de whatsapp\nExemplo: 559988525464 \nDigite aqui: ")
	fmt.Scanln(&numberWhat)

	formatNumber := fmt.Sprintf("%d@s.whatsapp.net", numberWhat)

	text := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: formatNumber,
		},
		Text: "Bot v1 em golang",
	}

	whatConn.Send(text)
	*/
}
