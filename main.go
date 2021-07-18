package main

import (
	"bufio"
	"fmt"
	"os"

	us "github.com/imhotep-nb/units/quantity"
)

// main is just simple conversion program.
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Type 'quit' to exit the loop.")
	for {
		fmt.Println()
		fmt.Print("Value:    ")
		scanner.Scan()
		s := scanner.Text()
		if s == "quit" {
			break
		}
		qu1, err := us.Parse(s)
		if err != nil {
			fmt.Println("Cannot parse")
			continue
		}
		fmt.Print("New unit: ")
		var symbol string
		scanner.Scan()
		symbol = scanner.Text()
		qu2, ok := qu1.ConvertTo(symbol)
		if ok {
			fmt.Println("New value:", qu2)
		} else {
			fmt.Println("Cannot convert")
		}
	}
}
