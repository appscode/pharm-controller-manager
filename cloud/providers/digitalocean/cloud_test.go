package digitalocean

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestRegion(t *testing.T) {
	resp, err := http.Get(instanceInfoURL + "/region")
	fmt.Println(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

}
