package types

import (
	"time"
	"fmt"
	"os"
	"github.com/edgexfoundry/edgex-go/support/consul-client"
)

type Endpoint struct {}

func(e Endpoint) Monitor(params EndpointParams, ch chan string) {
	for true {
		data, err := consulclient.GetServiceEndpoint(params.ServiceKey)
		if err != nil {
			fmt.Fprintln(os.Stdout, err.Error())
		}
		url := fmt.Sprintf("http://%s:%v%s", data.Address, data.Port, params.Path)
		ch <- url
		time.Sleep(15000 * time.Millisecond)
	}
}
