## 2018-08-23
This is a summary for p0 cmu.

## Go channel
I think it’s most about the go channel. That is, I can finish this much faster if I know better about channel
### Unblocking vs. blocking
1. blocking: make(chan string)
2. Unblocking: make(chan string,1)
It’s like will block the current program or not.
### data race
**go channel feels like parallel.**
So keep data race free.
- go channel itself is data race free.
- In my case,
I have multiple handleConnection to handle clients, which cause data race. 
The solution is simple, use on channel handleMessage to do the clients common message. And use client private channel to send message to outside.

## Producer to Consumer via channel
This code has problem, some data left in buffer when program finished, which would be lost.
``` golang
package main

import (
	"fmt"
)

const maxBufSize = 3
const numToProduce = 1000

var finishedProducing = make(chan bool)
var finishedConsuming = make(chan bool)

var messageBuffer = make(chan int, maxBufSize)

func produce() {
	for i := 0; i < numToProduce; i++ {
		messageBuffer <- i
	}

	finishedProducing <- true
}

func consume() {
	for {
		select {
		case message := <-messageBuffer:
			fmt.Println(message)
		case <-finishedProducing:
			finishedConsuming <- true
			return
		}
	}
}

func main() {
	go produce()
	go consume()
	<-finishedConsuming

	fmt.Println("All go routines ended")
}

```

## reason
Message in buffer haven't been handled yet.


## solution
We would get message from channel via for range, and close messageChannel in producer.
``` golang
package main

import (
	"fmt"
)

const maxBufSize = 3
const numToProduce = 1000

// var finishedProducing = make(chan bool)
var finishedConsuming = make(chan bool)

var messageBuffer = make(chan int, maxBufSize)

func produce() {
	for i := 0; i < numToProduce; i++ {
		messageBuffer <- i
    }

	close(messageBuffer)
}

func consume() {
	for message := range messageBuffer {
		fmt.Println(message)
	}
	finishedConsuming <- true

}

func main() {
	go produce()
	go consume()
	<-finishedConsuming

	fmt.Println("All go routines ended")
}

```

## improvement
we can use `make(chan struct{})` to replace `make(chan bool)`