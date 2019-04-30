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

import (
	"github.com/Oneledger/protocol/serialize"
)


/*
		Request
 */
type Request struct {
	Query string
	Params map[string]interface{}
	Ctx context
}

func NewRequest(query string, params map[string]interface{}) *Request {
	req := &Request{Query:query, Params:params}

	return req
}


func (r *Request) GetString(key string) string {
	s, ok := r.Params[key]
	if !ok {
		return ""
	}

	str, ok := s.(string)
	if !ok {
		return ""
	}

	return str
}

func (r *Request) GetInt(key string) int {
	s, ok := r.Params[key]
	if !ok {
		return 0
	}

	i, ok := s.(int)
	if !ok {
		return 0
	}

	return i
}


/*
		Response
 */


type Response struct {
	Data []byte			`json:"data"`
	ErrorMsg string		`json:"error_msg,omitempty"`
	Success bool		`json:"success"`
}

var jsonSerializer serialize.Serializer

func init() {
	jsonSerializer = serialize.GetSerializer(serialize.JSON)
}

func (r *Response) JSON(a interface{}) (err error) {

	r.Data, err = jsonSerializer.Serialize(a)
	return
}