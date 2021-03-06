package main

import (
	"io/ioutil"
	"log"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// f, _ := os.Create("cpuprof2")
	// defer f.Close()
	// pprof.StartCPUProfile(f)
	// pprof.StopCPUProfile()
	// c := make(chan TsPacket, 10)
	// n, _ := serv.Router.GetNodeByName("rmp")
	// n.RegisterListener(c)
	// for _ = range c {
	// }

	log.Print("starting")
	var err error
	cfgdat, _ := ioutil.ReadFile("basic.chq")

	serv := NewServer()
	err = serv.Router.ApplyConfig(string(cfgdat))
	Check(err)
	serv.Router.ToggleAll()

	err = serv.Start()
	Check(err)

	/////////

	// fname := flag.String("if", "", "input file")

	// addr := flag.String("a", "", "input multicast/unicast address")
	// port := flag.Int("p", 0, "input UDP port")
	// flag.Parse()

	// rtr := NewRouter()
	// cfgdat, _ := ioutil.ReadFile("basic.chq")
	// err := rtr.ApplyConfig(string(cfgdat))
	// Check(err)

	// out := make(chan TsPacket)

	// n := rtr.Nodes["t"]
	// n.RegisterListener(out)

	// for _ = range out {
	// }
	// for nn, np := range rtr.Nodes {
	// 	j, err := np.MarshalJSON()
	// 	log.Printf("node: %s\n json: %s\n err: %v\n", nn, j, err)
	// }
}
