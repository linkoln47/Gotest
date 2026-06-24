package main
import ( 
	"fmt"
	"time"
)
func main() { 
	c := make(chan string)
	go count("sheep", c)

	for { //instead checking open, could use range c 
		msg, open := <- c
		fmt.Println(msg)
		if !open {
			break
		}
	}


} 

func count(thing string, c chan string) { 
	for i := 0; i < 5; i++ { 
		c <- thing
		time.Sleep(time.Millisecond * 500) 
	}
	close(c)
}
