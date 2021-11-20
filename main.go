package main

import (
	"fmt"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
)

func login(whatConn *whatsapp.Conn) error {
	// Para o erro no retorno ummarshall
	whatConn.SetClientVersion(3, 2123, 7)

	qrCode := make(chan string)
	go func() {
		terminal := qrcodeTerminal.New2(qrcodeTerminal.ConsoleColors.BrightBlack, qrcodeTerminal.ConsoleColors.BrightWhite, qrcodeTerminal.QRCodeRecoveryLevels.Low)
		terminal.Get([]byte(<-qrCode)).Print()
	}()

	_, err := whatConn.Login(qrCode)
	if err != nil {
		return fmt.Errorf("Erro ao efetuar o login:", err)
	}

	return nil
}

func main() {
	waconn, err := whatsapp.NewConn(10 * time.Second) //10secs of timeout
	if err != nil {
		panic(err)
	}

	login(waconn)
}
