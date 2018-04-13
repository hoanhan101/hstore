package main

import (
	"kvraft"
    "bufio"
    "os"
    "fmt"
    "strings"
)

func main() {
	user := raftkv.Command{}

    fmt.Printf("Enter number of servers: ")
    var nservers int
    fmt.Scan(&nservers)

    user.Setup(nservers)
	user.Put("foo", "bar")
	user.Get("foo")

    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        rawString := strings.Split(scanner.Text(), " ")
        
        if rawString[0] == "GET" {
            user.Get(rawString[1])
        } else if rawString[0] == "PUT" {
            user.Put(rawString[1], rawString[2])
        } else {
            fmt.Println("Not supported method")
        }
    }
}
