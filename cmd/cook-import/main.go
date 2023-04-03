package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/solywsh/chatgpt"
	"time"
)

func main() {
		// https://pkg.go.dev/github.com/solywsh/chatgpt
		// The timeout is used to control the situation that the session is in a long and multi session situation.
		// If it is set to 0, there will be no timeout. Note that a single request still has a timeout setting of 30s.
		chat := chatgpt.New("openai_key", "user_id(not required)", 30*time.Second)
		defer chat.Close()
		//
		//select {
		//case <-chat.GetDoneChan():
		//	fmt.Println("time out/finish")
		//}
		question := "你认为2022年世界杯的冠军是谁？"
		fmt.Printf("Q: %s\n", question)
		answer, err := chat.Chat(question)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("A: %s\n", answer)
		log.Info("test")
}
