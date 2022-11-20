package main

import (
	"fmt"
	"github.com/gotranspile/cxgo/runtime/libc"
)

func main() {
	fmt.Println("Hello dbat!")

	InitCmdInfo()
	InitObjCommands()
	InitWldCmd()
	config_info.CONFFILE = libc.CString(CONFIG_FILE)
	load_config()
	boot_db()
	fmt.Println("happening 3")

}
