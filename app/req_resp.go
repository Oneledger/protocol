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
	"github.com/pkg/errors"
)

var (
	ErrParamNotFound = errors.New("param not found")
	ErrWrongParamType = errors.New("wrong param type")
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

func NewRequestFromObj(query, argName string, obj interface{}) (*Request, error) {
	req := &Request{Query: query}

	d, err := clSzlr.Serialize(obj)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new request")
	}

	req.Params = map[string]interface{}{argName: d}
	return req, nil
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

func (r *Request) GetBytes(key string) []byte {
	s, ok := r.Params[key]
	if !ok {
		return []byte{}
	}

	b, ok := s.([]byte)
	if !ok {
		return []byte{}
	}

	return b
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

func (r *Request) GetFloat64(key string) float64 {
	s, ok := r.Params[key]
	if !ok {
		return 0.0
	}

	i, ok := s.(float64)
	if !ok {
		return 0.0
	}

	return i
}

func (r *Request) GetBool(key string) (bool, error) {
	s, ok := r.Params[key]
	if !ok {
		return false, ErrParamNotFound
	}

	i, ok := s.(bool)
	if !ok {
		return false, ErrWrongParamType
	}

	return i, nil
}


func (r *Request) ClientDeserialize(name string, obj interface{}) error {
	d := r.GetBytes(name)
	if len(d) == 0 {
		return ErrParamNotFound
	}

	err := clSzlr.Deserialize(d, obj)
	if err != nil {
		return errors.Wrap(err, "request deserialization error")
	}
	return nil
}


/*
		Response
 */
type Response struct {
	Data []byte			`json:"data"`
	ErrorMsg string		`json:"error_msg,omitempty"`
	Success bool		`json:"success"`
}


func (r *Response) JSON(a interface{}) (err error) {

	r.Data, err = jsonSerializer.Serialize(a)
	r.Success = true
	return
}

func (r *Response) Error(msg string) {
	r.ErrorMsg = msg
	r.Success = false
}

func (r *Response) SetData(dat []byte) {
	r.Data = dat
	r.Success = true
}