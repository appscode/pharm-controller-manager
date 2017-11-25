package lightsail

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestZone(t *testing.T) {
	//.Println(fetchRegion(""))
}

func TestToken(t *testing.T) {
	tokenSource := &tokenSource{
		Token: "",
	}
	d, _ := json.Marshal(tokenSource)
	fmt.Println(string(d))

}
