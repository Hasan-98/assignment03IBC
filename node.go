package assignment03IBC

import (
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

func ReceiveChain(conn net.Conn) *Block {

	var recvdBlock *Block

	dec := gob.NewDecoder(conn)
	err := dec.Decode(&recvdBlock)

	if err != nil {

	}

	ListBlocks(recvdBlock)
	return recvdBlock

}

func StartListening(portNo string, user string) {

}

func main() {
	//it is better to check for arguments length and throw error
	satoshiAddress := os.Args[1]
	myListeningAddress := os.Args[2]

	if err != nil {
		log.Fatal(err)
	}

	//The function below launches the server, uses different second argument
	//It then starts a routine for each connection request received
	//go StartListening(myListeningAddress, "others")
	//StartListening(myListeningAddress, "others")
	//log.Println("Sending my listening address to Satoshi")
	//Satoshi is there waiting for our address, it stores it somehow
	//WriteString(conn, myListeningAddress)

	//once the satoshi unblocks on Quorum completion it sends peer to connect to
	//log.Println("receiving peer to connect to ... ")
	//receivedString := ReadString(conn)
	//log.Println(receivedString)

	//Then satoshi sends the chain with x+1 blocks
	//log.Println("receiving Chain")
	chainHead := ReceiveChain(conn)
	log.Println(CalculateBalance("Satoshi", chainHead))

	//Each node then connects to the other peer info received from satoshi
	//The topology eventually becomes a ring topology
	//Each node then both writes the hello world to the connected peer
	//and also receives the one from another peer

	//log.Println("connecting to the other peer ... ", receivedString, "hamza", "mohtasim")
	//ad := ":" + receivedString
	//log.Println("connecting to the other peer ... ", ad)

	/*	peerConn, err1 := net.Dial("tcp", receivedString)

		if err1 != nil {
			log.Println(err1)
		} else {
			fmt.Fprintf(peerConn, "Hello From %v to %v\n", myListeningAddress, receivedString)
		}*/
	select {}
}

/////////////////////////////////////////  blockchain code    /////////////////////////////////////////////////////

const miningReward = 100
const rootUser = "Satoshi"

type Block struct {
	Spender     map[string]int
	Receiver    map[string]int
	PrevPointer *Block
	PrevHash    string
	CurrentHash string
}

func CalculateBalance(userName string, chainHead *Block) int {

	var temp *Block
	var in = 0
	var out = 0
	//var net = 0
	var len = 0

	if userName == "Satoshi" {

		for temp = chainHead; temp != nil; {

			for key, value := range temp.Receiver {
				//	fmt.Println("cpmapring=", userName, key)
				if userName == key {

					in = in + value

				}
				len = len + 100

			}

			for key, value := range temp.Spender {

				if userName == key {

					out = out + value
				}
			}

			temp = temp.PrevPointer
		}
		//	fmt.Println("returning valur=", (in-out)+len)
		return (in - out) + len
	}

	if userName != "Satoshi" {

		for temp = chainHead; temp != nil; {

			for key, value := range temp.Receiver {
				if userName == key {

					in = in + value

				}

			}

			for key, value := range temp.Spender {

				if userName == key {

					out = out + value
				}
			}

			temp = temp.PrevPointer
		}
		//	fmt.Println("returning valur=", (in - out))
		return (in - out)
	}
	return 0
}

func InsertBlock(spendingUser string, receivingUser string, miner string, amount int, chainHead *Block) *Block {
	var temp *Block = new(Block)

	if chainHead == nil {

		if spendingUser == "" && receivingUser == "" && amount == 0 && miner == "Satoshi" {

			fmt.Println(" valid conditions for the transactions -> gensis done")
			temp.PrevHash = ""
			temp.PrevPointer = nil
			temp.Spender = make(map[string]int)

			temp.Spender[""] = amount
			temp.Receiver = make(map[string]int)
			temp.Receiver[""] = amount
			temp.CurrentHash = CalculateHash(temp)
			//fmt.Println("hash=", temp.CurrentHash)
			return temp

		}

		fmt.Println("invalid transaction denied")
		temp = nil
		return temp

	}

	if chainHead != nil {
		var temp *Block = new(Block)
		//fmt.Println(spendingUser, CalculateBalance(spendingUser, chainHead))
		if miner == "Satoshi" && CalculateBalance(spendingUser, chainHead) >= amount {
			//	fmt.Println("valid transaction checks444444")
			temp.PrevHash = chainHead.CurrentHash
			temp.PrevPointer = chainHead
			temp.Spender = make(map[string]int)

			temp.Spender[spendingUser] = amount
			temp.Receiver = make(map[string]int)
			temp.Receiver[receivingUser] = amount
			temp.CurrentHash = CalculateHash(temp)
			//fmt.Println("hash=", temp.CurrentHash)
			fmt.Println(" valid conditions for the transactions -> done")
			return temp

		}

		if miner != "Satoshi" || CalculateBalance(spendingUser, chainHead) < amount {
			fmt.Println("not valid conditions for the transactions")
			temp = nil
			return chainHead
		}
		fmt.Println("invalid transaction denied")
		return chainHead

	}
	fmt.Println("invalid transaction denied")
	return chainHead
}

func CalculateHash(inputBlock *Block) string {
	//	a := len(inputBlock.Spender)
	//	b := len(inputBlock.Receiver)
	var temp string

	for key, value := range inputBlock.Receiver {

		temp = temp + key + strconv.Itoa(value)

	}

	for key, value := range inputBlock.Spender {

		temp = temp + key + strconv.Itoa(value)

	}

	//fmt.Println(temp)
	obj := sha256.New()
	obj.Write([]byte(fmt.Sprintf("%x", temp)))

	return fmt.Sprintf("%x", obj.Sum(nil))
}

func ListBlocks(chainHead *Block) {
	var temp *Block
	for temp = chainHead; temp != nil; {
		fmt.Println("transaction=")
		for key, value := range temp.Receiver {
			fmt.Println("receiving=", key, "money:", value)

		}

		for key, value := range temp.Spender {
			fmt.Println("sender=", key, "money:", value)

		}
		fmt.Println("--------")
		temp = temp.PrevPointer
	}
}

func VerifyChain(chainHead *Block) {

	var temp *Block
	for temp = chainHead; temp != nil; {
		temp2 := temp.PrevPointer
		//fmt.Println("transaction=", temp.transactions)
		if temp2 != nil {
			pHash := CalculateHash(temp.PrevPointer)
			if pHash != temp.PrevHash {
				fmt.Println("change detected")
				return
			}
		}
		temp = temp.PrevPointer
	}
}
