package peerProc

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Holds the Snip Information
type Snip struct {
	message    string
	senderAddr string
	timeStamp  int
}

// This function formats the list into a string and returns a list of snips received as a string
func PreparelistSnipsToString() string {
	var numSnips string = strconv.Itoa(len(listSnips))
	var snipList string = numSnips + "\n"

	for _, snip := range listSnips {
		snipList += strconv.Itoa(snip.timeStamp) + " " + snip.message + " " + snip.senderAddr + "\n"
	}
	return snipList
}

// Handles what happends when you get a snip
func SnipHandler(sourceAddress string, conn *net.UDPConn, ctx context.Context) {
	ch := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			sendSnip(msg, sourceAddress, conn)
		}
	}
}

//Send a snip to the peer
func sendSnip(msg string, sourceAddress string, conn *net.UDPConn) {

	currentTime++
	snipCurrentTime := strconv.Itoa(currentTime)
	msg = "snip" + snipCurrentTime + " " + msg
	mutex.Lock()
	// Send the message to all peers
	// fmt.Println("Sending messages")
	for _, peer := range listPeers {
		if CheckForValidAddress(peer.peerAddress) {
			go sendMessage(peer.peerAddress, msg, conn)
		} else {

		}
	}
	mutex.Unlock()
}

// After receiving a snip store it for report
func storeSnips(command string, senderAddr string) {
	msg := strings.Split(command, " ")
	timestamp, err := strconv.Atoi(msg[0])
	if err != nil {
		fmt.Println("Timestamp is not a valid number")
		return
	}
	if len(msg) < 2 {
		fmt.Printf("Invalid snip command: \n message: %s%s\n", command, msg)
		return
	}
	// Store the snip to list
	// join the rest of the message
	snipContent := strings.Join(msg[1:], " ")

	// check which time is the latest
	if senderAddr != mainUdpAddress {
		currentTime = getMAxValue(currentTime, timestamp)
	}
	mutex.Lock()
	listSnips = append(listSnips, Snip{snipContent, senderAddr, currentTime})
	mutex.Unlock()

	// update last seen
	for i := 0; i < len(listPeers); i++ {
		if listPeers[i].peerAddress == senderAddr {
			listPeers[i].lastSeen = time.Now()
			listPeers[i].isAlive = true
		}
	}

}
