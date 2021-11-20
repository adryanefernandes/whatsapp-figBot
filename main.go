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

	onInit(whatConn)
}
