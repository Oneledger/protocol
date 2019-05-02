/*
   ____             _              _                      _____           _                  _
  / __ \           | |            | |                    |  __ \         | |                | |
 | |  | |_ __   ___| |     ___  __| | __ _  ___ _ __     | |__) | __ ___ | |_ ___   ___ ___ | |
 | |  | | '_ \ / _ \ |    / _ \/ _` |/ _` |/ _ \ '__|    |  ___/ '__/ _ \| __/ _ \ / __/ _ \| |
 | |__| | | | |  __/ |___|  __/ (_| | (_| |  __/ |       | |   | | | (_) | || (_) | (_| (_) | |
  \____/|_| |_|\___|______\___|\__,_|\__, |\___|_|       |_|   |_|  \___/ \__\___/ \___\___/|_|
                                      __/ |
                                     |___/


Copyright 2017 - 2019 OneLedger
*/

package app

import "fmt"

func Example_newRouter() {

	r := NewRouter("test")

	err := r.AddHandler("exec", exec)
	if err != nil {
		// r.logger.Fatal("router error", err)
	}

	p := map[string]interface{}{
		"name":   "exec",
		"number": 1,
	}
	req := NewRequest("/exec", p)
	resp := &Response{}
	r.Handle(*req, resp)

	fmt.Println("response:", string(resp.Data))
	// Output
	// name: name
	// not set:
	// number: 1
	// response: function response
}

func exec(req Request, resp *Response) {
	fmt.Println(req.Query)

	fmt.Println("name:", req.GetString("name"))
	fmt.Println("not set:", req.GetString("not_set"))

	fmt.Println("number:", req.GetInt("number"))

	err := resp.JSON("function response")
	if err != nil {
		// log.Error("error marshalling response", err)
	}
}
