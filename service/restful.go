package service

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"

	"github.com/Oneledger/protocol/data/keys"
	"github.com/Oneledger/protocol/serialize"
)

type RestfulRouter map[string]http.HandlerFunc

//basic RestfulService entry point
type RestfulService struct {
	ctx    *Context
	router RestfulRouter
}

func NewRestfulService(ctx *Context) RestfulService {
	svc := RestfulService{
		ctx:    ctx,
		router: make(RestfulRouter),
	}

	svc.router["/"] = svc.restfulRoot()
	svc.router["/health"] = svc.health()
	svc.router["/token"] = svc.GetToken()

	return svc
}

func (rs RestfulService) Router() RestfulRouter {
	return rs.router
}

func (rs RestfulService) restfulRoot() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		_, err := fmt.Fprintln(w, "Available endpoints: ")
		if err != nil {
			rs.ctx.Logger.Errorf("failed to display available endpoints info")
		}
		for path := range rs.router {
			_, err = fmt.Fprintln(w, r.Host+path)
			if err != nil {
				rs.ctx.Logger.Errorf("failed to display available endpoints info")
			}
		}
	}
}

func (rs RestfulService) health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "health check for SDK port %v : OK\n", rs.ctx.Cfg.Network.SDKAddress)
		if err != nil {
			rs.ctx.Logger.Errorf("failed to display SDK port health check info")
		}
	}
}

type clientTokenReq struct {
	ClientID string `json:"clientId"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type clientTokenResp struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

func (rs RestfulService) GetToken() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		request := make([]byte, 0, 10)
		serializer := serialize.GetSerializer(serialize.NETWORK)
		clientReq := &clientTokenReq{}
		clientResp := &clientTokenResp{}

		defer func() {
			bytes, _ := serializer.Serialize(clientResp)
			_, _ = w.Write(bytes)
			_ = r.Body.Close()
		}()

		h := sha256.New()

		rs.ctx.Logger.Error("Deserializing request into Object")
		request, err := ioutil.ReadAll(r.Body)
		err = serializer.Deserialize(request, clientReq)
		if err != nil {
			clientResp.Message = errors.Wrap(err, "failed to get request body").Error()
			return
		}

		rs.ctx.Logger.Error("Validating Username and Password")
		//Validate Username and Password from request
		if len(rs.ctx.Cfg.Node.Auth.OwnerCredentials) > 0 {
			userList := rs.ctx.Cfg.Node.Auth.OwnerCredentials

			for index, value := range userList {
				//parse value into a pair
				pair := strings.Split(value, ":")
				if pair[0] == clientReq.Username && pair[1] == clientReq.Password {
					break
				}

				if index == len(userList)-1 {
					clientResp.Message = "Invalid Username and/or Password"
					return
				}
			}
		}

		if rs.ctx.Cfg.Node.Auth.RPCPrivateKey != "" {

			//Hash clientID given in the request and create an encoded token based on the hash
			h.Write([]byte(clientReq.ClientID))
			data := h.Sum(nil)[:20]

			//Get Private Key for signature.
			keyData, err := base64.StdEncoding.DecodeString(rs.ctx.Cfg.Node.Auth.RPCPrivateKey)
			if err != nil {
				rs.ctx.Logger.Error("Error decoding rpc private key", err)
				clientResp.Message = err.Error()
				return
			}

			privateKey, err := keys.GetPrivateKeyFromBytes(keyData, keys.ED25519)
			if err != nil {
				rs.ctx.Logger.Error("Error Getting private key from Bytes", err)
				clientResp.Message = err.Error()
				return
			}

			//Get Key handler
			privateKeyHandler, err := privateKey.GetHandler()
			if err != nil {
				rs.ctx.Logger.Error("Error Getting private key Handler", err)
				clientResp.Message = err.Error()
				return
			}

			//Generate signature using private key
			signature, err := privateKeyHandler.Sign(data)
			tokenBytes := append(data, signature...)

			//Populate response fields.
			clientResp.Token = base58.Encode(tokenBytes)
			clientResp.Message = "Token generated successfully"

		} else {
			clientResp.Message = "No Private Key configured for this NODE"
		}
	}
}
