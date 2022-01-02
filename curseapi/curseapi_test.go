package curseapi

import (
	"encoding/json"
	"testing"

	"fmt"
)

func TestSearchmod(t *testing.T) {
	l, err := Searchmod("", 0, 6)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.Marshal(l)
	fmt.Println(string(b))
}
