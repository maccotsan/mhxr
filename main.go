package main

import (
	"github.com/maccotsan/mhxr/schedule"
	"fmt"
	"os"
)
func main() {
	html, err := schedule.CreateHTML()
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	file, err := os.OpenFile("./pages/schedule.html", os.O_CREATE | os.O_WRONLY, 0666)
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	file.WriteString(html)
	defer file.Close()

	fmt.Println("Generate scadule.html success.")
}
