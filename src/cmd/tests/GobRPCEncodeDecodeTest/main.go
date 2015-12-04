package main

import (
	"fmt"
	"log"
	"os"
	"rpc"
)

func main() {
	f, err := os.Create("gob.txt")
	if err != nil {
		log.Panic(err)
	}
	encoder := rpc.NewGobRPCEncoder(f)
	encoder.Encode("eid1", "method1", []interface{}{1, 2, "3\n3"})
	encoder.Encode("eid2", "method2", []interface{}{1, false, 3.333})

	_ = f.Close()

	f, _ = os.Open("json.txt")
	decoder := rpc.NewGobRPCDecoder(f)
	var eid, method string
	var args []interface{}
	_ = decoder.Decode(&eid, &method, &args)
	fmt.Printf("Decode: %v.%v(%v)\n", eid, method, args)

	_ = decoder.Decode(&eid, &method, &args)
	fmt.Printf("Decode: %v.%v(%v)\n", eid, method, args)

	f.Close()
}
