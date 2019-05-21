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

package data

import (
	"github.com/pkg/errors"
)

var (
	ErrParamNotFound  = errors.New("param not found")
	ErrWrongParamType = errors.New("wrong param type")
)

/*
	Request
*/
// Request generic request object for Query handling. Request comes with its query
// string and parameters
type Request struct {
	Query  string                 // the query method
	Data   []byte                 // data stored as a serialized byte array
	Params map[string]interface{} // request params as a map of string to interface
}

// NewRequest creates a new request with given Params.
func NewRequest(query string, params map[string]interface{}) (*Request, error) {
	req := &Request{Query: query, Params: params}

	// serialize the params
	d, err := clSzlr.Serialize(params)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new request")
	}

	// assign the data to request
	req.Data = d
	return req, nil
}

// NewRequestFromData creates a new request object from a byte array
func NewRequestFromData(query string, dat []byte) *Request {
	req := &Request{Query: query, Data: dat}
	return req
}

// NewRequestFromObj creates a new request object from an arguments struct passed.
// It serializes the argument struct object for a client channel and sets it against an argname.
// You can check example argument structs in client/request.
func NewRequestFromObj(query string, obj interface{}) (*Request, error) {
	req := &Request{Query: query}

	d, err := clSzlr.Serialize(obj)
	if err != nil {
		return nil, errors.Wrap(err, "error creating new request")
	}

	req.Data = d
	return req, nil
}

// Parse parses the serialized parameters data and saves it to Params
func (r *Request) Parse() {
	p := map[string]interface{}{}

	err := clSzlr.Deserialize(r.Data, &p)
	if err != nil {
		// log
	}

	r.Params = p
}

// GetString retrieves a string parameter saved in a request object. If a string parameter
// is not set or the object type of the parameter is not string an empty string is returned. All calling instances are required to check and handle empty
// string.
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

// GetBytes retrieves a byte array
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

	err := clSzlr.Deserialize(r.Data, obj)
	if err != nil {
		return errors.Wrap(err, "request deserialization error")
	}
	return nil
}

/*
	Response
*/
type Response struct {
	Data     []byte `json:"data"`
	ErrorMsg string `json:"error_msg,omitempty"`
	Success  bool   `json:"success"`
}

func (r *Response) JSON(a interface{}) (err error) {

	r.Data, err = jsonSzlr.Serialize(a)
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

func (r *Response) SetDataObj(obj interface{}) error {
	d, err := clSzlr.Serialize(obj)
	if err != nil {
		return err
	}

	r.SetData(d)
	return nil
}
