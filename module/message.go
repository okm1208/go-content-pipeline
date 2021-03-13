package module

import (
	"content-pipline-example/common"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

func MsgDelivery( inputStream <-chan amqp.Delivery ,
	outputStream chan common.JobContext,
){

	log.Println("Msg delivery Start ===")
	go func(){
		defer func() {
			log.Println("Msg delivery End ====")
			close(outputStream)
		}()
		for {
			select {
			case delivery,more := <- inputStream:
				if more {
					msg := common.QueueMessage{}
					err := json.Unmarshal(delivery.Body, &msg)
					delivery.Ack(false)
					if err != nil {
						log.Printf("Msg unmarshal err : %v\n",err)
						continue
					}
					outputStream <- common.JobContext{TrackId: msg.TrackId}
				}else {
					return
				}
			}
		}
	}()
}
