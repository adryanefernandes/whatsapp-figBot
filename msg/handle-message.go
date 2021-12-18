package msg

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Rhymen/go-whatsapp"
)

type waHandler struct {
	c         *whatsapp.Conn
	startTime uint64
}

func (wh *waHandler) HandleTextMessage(message whatsapp.TextMessage) {
	if message.Info.Timestamp < wh.startTime {
		return
	}

	fmt.Printf("%v %v %v %v\n\t%v\n", message.Info.Timestamp, message.Info.Id, message.Info.RemoteJid, message.ContextInfo.QuotedMessageID, message.Text)
}

func (wh *waHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	data, err := message.Download()

	// TODO: pegar isso do whatsapp
	errMediaDownloadFailedWith410 := errors.New("download failed with status code 410")
	errMediaDownloadFailedWith404 := errors.New("download failed with status code 404")

	if err != nil {
		if err != errMediaDownloadFailedWith410 && err != errMediaDownloadFailedWith404 {
			return
		}
		if _, err = wh.c.LoadMediaInfo(message.Info.RemoteJid, message.Info.Id, strconv.FormatBool(message.Info.FromMe)); err == nil {
			data, err = message.Download()
			if err != nil {
				return
			}
		}
	}

	filename := fmt.Sprintf("%v/%v.%v", os.TempDir(), message.Info.Id, strings.Split(message.Type, "/")[1])
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		return
	}
	file.Write(data)
	if err != nil {
		return
	}
	log.Printf("%v %v\n\timage received, saved at:%v\n", message.Info.Timestamp, message.Info.RemoteJid, filename)
}

func AddHandler(whatConn *whatsapp.Conn) {
	whatConn.AddHandler(&waHandler{whatConn, uint64(time.Now().Unix())})
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
