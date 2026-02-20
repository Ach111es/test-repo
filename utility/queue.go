package utility

import (
	"fmt"
	"strconv"

	"github.com/go-stomp/stomp"
)

func SendMessageToQueue(configuration Configuration, strMessage string) error {
	strServer := configuration.Queue.Host + ":" + strconv.Itoa(configuration.Queue.Port)

	connOpts := []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(configuration.Queue.Username, configuration.Queue.Password),
	}

	conn, err := stomp.Dial("tcp", strServer, connOpts...)
	if err != nil {
		PrintConsole(fmt.Sprintf("%v", err.Error()), "error")
		return err
	}

	defer conn.Disconnect()

	err = conn.Send(configuration.Queue.QueueName, "text/plain", []byte(strMessage), nil)
	if err != nil {
		PrintConsole(fmt.Sprintf("%v", err.Error()), "error")
		return err
	}
	PrintConsole("Queue : Message successfully sent to queue", "info")
	return nil
}
