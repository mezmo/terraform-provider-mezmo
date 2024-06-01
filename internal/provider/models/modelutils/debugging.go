package modelutils

import (
	"encoding/json"
	"fmt"
	"log"
)

func Json(label string, obj any) string {
	json, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return fmt.Sprintf("%s %s", label, string(json))
}
