package modelutils

import (
	"encoding/json"
	"fmt"
	"log"
)

func PrintJSON(label string, obj any) {
	json, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(label, string(json))
}
