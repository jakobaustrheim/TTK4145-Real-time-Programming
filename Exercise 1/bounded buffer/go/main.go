
package main

import "fmt"
import "time"


func producer(ch chan<- int){

    for i := 0; i < 10; i++ {
        time.Sleep(100 * time.Millisecond)
        fmt.Printf("[producer]: pushing %d\n", i)
        // TODO: push real value to buffer
        ch <- i //Legger den nye verdien til i inn i bufferen

    }

}

func consumer(ch <-chan int){

    time.Sleep(1 * time.Second)
    for {
        i := 0 //TODO: get real value from buffer
        i = <-ch //Henter verdien til i fra bufferen
        fmt.Printf("[consumer]: %d\n", i)
        time.Sleep(50 * time.Millisecond)
    }
    
}


func main(){
    
    // TODO: make a bounded buffer
    buffer := make(chan int, 5) //Bruker make for å lage en buffer med plass til 5 elementer
    
    go consumer(buffer)
    go producer(buffer)
    
    select {}
}