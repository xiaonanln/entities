package main

import (
	"bytes"
	"flag"
	"log"
	"rpc"
	"time"
)

const (
	PROFILE_COUNT = 300000

	TEST_EID    = "testentityid"
	TEST_METHOD = "testmethodname"
)

var (
	TEST_ARGS_SMALL = []interface{}{1, 2.0, "string", false}
	TEST_ARGS_LARGE = []interface{}{}

	profileJson       = flag.Bool("json", false, "profile json")
	profileGob        = flag.Bool("gob", false, "profile gob")
	profileCustom     = flag.Bool("custom", false, "profile custom")
	useLargeArguments = flag.Bool("large", false, "use large arguments")
)

func init() {
	for i := 0; i < 1000; i++ {
		for _, j := range TEST_ARGS_SMALL {
			TEST_ARGS_LARGE = append(TEST_ARGS_LARGE, j)
		}
	}
}

func main() {
	flag.Parse()
	if *profileJson {
		jsonBuffer := bytes.NewBuffer([]byte{})
		jsonRpcEncoder := rpc.NewJsonRPCEncoder(jsonBuffer)
		jsonRpcDecoder := rpc.NewJsonRPCDecoder(jsonBuffer)
		profile("json", jsonRpcEncoder, jsonRpcDecoder)
	}

	if *profileGob {
		gobBuffer := bytes.NewBuffer([]byte{})
		gobRpcEncoder := rpc.NewGobRPCEncoder(gobBuffer)
		gobRpcDecoder := rpc.NewGobRPCDecoder(gobBuffer)
		profile("gob", gobRpcEncoder, gobRpcDecoder)
	}

	if *profileCustom {
		customBuffer := bytes.NewBuffer([]byte{})
		customRpcEncoder := rpc.NewCustomRPCEncoder(customBuffer)
		customRpcDecoder := rpc.NewCustomRPCDecoder(customBuffer)
		profile("custom", customRpcEncoder, customRpcDecoder)
	}
}

func profile(id string, encoder rpc.RPCEncoder, decoder rpc.RPCDecoder) {
	args := TEST_ARGS_SMALL
	profileCount := PROFILE_COUNT
	if *useLargeArguments {
		log.Printf("%s: using large arguments!!!\n", id)
		args = TEST_ARGS_LARGE
		profileCount = profileCount / 100
	}

	startTime := time.Now()

	for i := 0; i < profileCount; i++ {
		err := encoder.Encode(TEST_EID, TEST_METHOD, args)
		if err != nil {
			log.Panic(err)
		}
		var eid, method string
		var args []interface{}
		err = decoder.Decode(&eid, &method, &args)
		if err != nil {
			log.Panic(err)
		}
		if eid != TEST_EID || method != TEST_METHOD {
			log.Panicln("Encode & Decode mismatch")
		}
	}
	elapsedTime := time.Since(startTime)
	log.Printf("profile %s takes %f seconds\n", id, elapsedTime.Seconds())
}
