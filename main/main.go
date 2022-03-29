package main

import (
	"flag"
	"log"
)
const(
	MTTU = 1
	EXTERNAL = 2
)

func main() {
	// mongo.InitTunnel()
	var tye int
	var name string
	var version string
	flag.IntVar(&tye,"t", 1, "which type? 1 mttu 2 external 3 print deptree 4 get all version")
	flag.StringVar(&name, "n", "field-descriptions", "depName")
	flag.StringVar(&version,"v","1.0.7","depVersion")
	flag.Parse()
	if tye <= 0 && tye >= 5{
		log.Fatal("please give me true type.")
	}
	if tye == 1{
		NpmMttu(name, version)
	}else if tye ==2{
		NpmExternal(name, version)
	}else if tye == 3{
		NpmAllVersion(name)
	}else{
		NpmSaveTree(name,version)
	}

}
