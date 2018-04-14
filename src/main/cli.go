package main

import (
	"bufio"
	"fmt"
	"kvraft"
	"os"
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
			if len(rawString) != 3 || rawString[2] == "" {
				fmt.Println("Cannot PUT empty value")
				continue
			} else {
				user.Put(rawString[1], rawString[2])
			}
		} else if rawString[0] == "APPEND" {
			user.Append(rawString[1], rawString[2])
		} else {
			fmt.Println("Not supported method")
		}
	}
}
