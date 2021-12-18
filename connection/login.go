package connection

import (
	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
)

func Login(whatConn *whatsapp.Conn) (whatsapp.Session, error) {
	qrCode := make(chan string)
	go func() {
		terminal := qrcodeTerminal.New2(qrcodeTerminal.ConsoleColors.BrightBlack, qrcodeTerminal.ConsoleColors.BrightWhite, qrcodeTerminal.QRCodeRecoveryLevels.Low)
		terminal.Get([]byte(<-qrCode)).Print()
	}()

	return whatConn.Login(qrCode)
}
