package amq

import (
	"github.com/go-stomp/stomp"
	"strconv"
	"net/http"
	"time"
	"fmt"
	"github.com/Fjolnir-Dvorak/manageAMQ/amq/response"
	"encoding/json"
)
type server struct{
	connection *stomp.Conn
	server     string
	username   string
	password   string
	port int
}
var srv server

func Connect(host, username, password string, stompPort int) error {
	connection, err := stomp.Dial("tcp", host + ":" + strconv.Itoa(stompPort),
		stomp.ConnOpt.Login(username, password),
		stomp.ConnOpt.AcceptVersion(stomp.V11),
		stomp.ConnOpt.AcceptVersion(stomp.V12),
		stomp.ConnOpt.Host("localFoobar"),
		stomp.ConnOpt.UseStomp)
	if err != nil {
		return err
	}
	srv = server {
		connection: connection,
		server:     host,
		username:   username,
		password:   password,
		port: stompPort,
	}
	return nil
}

func reconnect(times int, waitBetween time.Duration) error {
	for i := 1; i <= times; i++ {
		fmt.Println("[AMQ] trying to reconnect...")
		connection, err := stomp.Dial("tcp", srv.server + ":" + strconv.Itoa(srv.port),
			stomp.ConnOpt.Login(srv.username, srv.password),
			stomp.ConnOpt.AcceptVersion(stomp.V11),
			stomp.ConnOpt.AcceptVersion(stomp.V12),
			stomp.ConnOpt.Host("localFoobar"),
			stomp.ConnOpt.UseStomp)
		if err != nil {
			if i < times {
				time.Sleep(waitBetween)
			} else {
				return err
			}
		}
		srv.connection = connection
		break
	}
	return nil
}

func SendMessage(destination, message string) (err error) {
	//fmt.Printf("trying to write message into queue %s: %s\n", "/queue/" + destination, message)
	err = srv.connection.Send(destination, "text/plain", []byte(message), stomp.SendOpt.NoContentLength)
	if err != nil {
		fmt.Printf("[AMQ] ERROR happened: %s\n", err)
		err = reconnect(5, 5 * time.Second)
		if err != nil {
			fmt.Printf("[AMQ] could not reconnect: %s\n", err)
			return err
		}
		err = srv.connection.Send(destination, "text/plain", []byte(message), stomp.SendOpt.NoContentLength)
		if err != nil {
			fmt.Printf("[AMQ] ERROR happened: %s\n", err)
			return err
		}
	}
	return nil
}

func Disconnect() error {
	err := srv.connection.Disconnect()
	srv = server{}
	return err
}

func GetEnqueuedCount(queue string, port int, username string, password string) (size int, err error) {
	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	url := fmt.Sprintf("http://%s:%d/api/jolokia/read/org.apache.activemq:type=Broker,brokerName=localhost,destinationType=Queue,destinationName=%s", srv.server, port, queue)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return -1, err
	}
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		return -1, err
	}
	var parsed response.ReadQueue
	err = json.NewDecoder(resp.Body).Decode(&parsed)
	if err != nil {
		return -1, err
	}
	if parsed.Status == 404 {
		return -1, fmt.Errorf("QueueNotFound")
	}
	queueSize := parsed.Value.QueueSize
	return queueSize, nil
}
