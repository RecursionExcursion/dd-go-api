package lib

import (
	"encoding/json"
	"fmt"
	"log"
)

/* Prints data as JSON to the terminal for readability */
func PrettyLog(s any) {
	out, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Fatalf("PrettyLog error: %v", err)
		panic(err)
	}

	fmt.Println(string(out))
}
