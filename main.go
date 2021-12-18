package main

import (
	"time"

	"github.com/Rhymen/go-whatsapp"
	"github.com/qgx-pagamentos/whatsapp-figBot/connection"
	"github.com/qgx-pagamentos/whatsapp-figBot/msg"
)

func main() {
	whatConn, err := whatsapp.NewConn(10 * time.Second)
	if err != nil {
		panic(err)
	}
	// Para o erro no retorno ummarshall
	whatConn.SetClientVersion(3, 2123, 7)

	connection.NewSession(whatConn)

	msg.AddHandler(whatConn)
	<-time.After(60 * time.Minute)
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

	whatConn.Send(text) */
}
