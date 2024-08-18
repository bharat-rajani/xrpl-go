package main

import (
	"context"
	"fmt"
	"github.com/xrpscan/xrpl-go/client"
	"github.com/xrpscan/xrpl-go/pkg/ripplexcodec"
	"log"
)

func main() {

	x := "XVrHo9Rzj1LCrK2oVJ5H6skgAZUKHtD7SoQ8j7PFN7Nm7AP"
	address, err := ripplexcodec.DecodeXAddress(x)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(address)

	ctx := context.Background()
	websocketCM := client.NewWebsocketConnectionManager()
	connection, err := websocketCM.NewWebsocketConnection(ctx, "wss://s.devnet.rippletest.net:51233")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err := websocketCM.Close()
		fmt.Println(err)
		return
	}()

	c := client.NewClient(ctx, connection, client.ClientConfig{})
	go func() {
		err := c.ReadMessages()
		if err != nil {
			log.Println(err)
		}
	}()
	_, err = c.GetXRPBalance(context.Background(), "rf1BiGeXwwQoi8Z2ueFYTEXSwuJYfV2Jpn")
	fmt.Println(err)
}
