package modelutils

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrintJSON(obj any) {
	json, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(string(json))
}
