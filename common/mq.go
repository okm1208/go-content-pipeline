package common

import (
	"github.com/streadway/amqp"
	"log"
)

type MqConn struct {
	Connection *amqp.Connection
	Channel *amqp.Channel
}

func (mqConn *MqConn) InitConn(rabbitURL string) error {
	log.Println("RabbitMQ GetConnection.")
	//TCP 연결 (물리적)
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return err
	}
	mqConn.Connection = conn
	log.Println("RabbitMQ Connection success.")

	//채널 연결 (논리적)
	chnl , err := mqConn.Connection.Channel()
	if err != nil {
		return err
	}
	mqConn.Channel = chnl
	log.Println("RabbitMQ Channel Open success.")

	return nil
}

func (mqConn *MqConn) CloseConn(){
	log.Printf("RabbitMq close.\n")

	if err :=mqConn.Channel.Cancel("mi-bulk", true); err != nil {
		log.Printf("Consumer_cancel_failed: %s %v\n", mqConn.Channel,err)
	}

	if err :=mqConn.Channel.Close(); err != nil {
		log.Printf("Consumer_Close_failed: %s %v\n", mqConn.Channel,err)
	}

	log.Printf("RabbitMq Channel close.\n")
	if !mqConn.Connection.IsClosed(){
		mqConn.Connection.Close()
	}
	log.Printf("RabbitMq Connection Close.\n")

}