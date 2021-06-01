package tests

import (
	"fmt"
	"net/rpc"
	"strconv"

	"github.com/tofa-project/client-daemon/glob"
)

func RPC() {
	client, err := rpc.DialHTTP("tcp", "localhost:"+strconv.Itoa(glob.V_RPC_PORT))
	if err != nil {
		panic(err)
	}

	var reply interface{}
	err = client.Call("Methods.SomeMeth", []interface{}{0}, &reply)
	if err != nil {
		panic(err)
	}

	fmt.Println(reply)
}
