package main

import "fmt"

func logResult(err error) {
	if err != nil {
		problems.WriteString(fmt.Sprintln(err))
		fmt.Print("x")
		return
	}

	fmt.Print(".")
	return
}
