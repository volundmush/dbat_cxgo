package main

import (
	"fmt"
	"os"
)

func main() {

	fmt.Println("Loading config...")
	config_info.CONFFILE = "config"
	fmt.Println(os.Getwd())
	load_config()

	port = config_info.Operation.DFLT_PORT

	setup_log(config_info.Operation.LOGNAME, 2)

	fmt.Println("It compiles!", port)

	init_game(port)

	fmt.Println("We got here!")

}
