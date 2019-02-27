package lite

import "./cli"

func main() {
	//bc := core.NewBlockchain()
	//defer bc.CloseDB()
	//cli := core.NewCLI(bc)
	//cli.Run()
	cli := cli.CLI{}
	cli.Run()
}
