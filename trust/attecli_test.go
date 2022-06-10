package trust

import (
	"fmt"
	"testing"
)

func TestNewAttestationHttpClient(t *testing.T) {
	cururl := "http://192.168.1.221:8999"
	cururl = "http://106.3.133.179:8999"
	//cururl = "http://164.52.51.10:8999"
	cli := NewAttestationHttpClient(cururl)
	//127.0.0.1:8999
	creditNodes, err := cli.GetNodeRank()
	if err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Println(creditNodes)
}
