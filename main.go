package main

import (
	"example/learn/p2p"
	"fmt"
	"log"
)

func main(){
	tr := p2p.NewTCPTransport(":3000")
	fmt.Println("Listening to New Request!!")
	if err := tr.ListenAndAccept(); err!=nil{
		log.Fatal(err)
	}
	select { }
}