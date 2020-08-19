package devtools4chains

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// RPCInfo .
type RPCInfo struct {
	Host     string
	User     string
	Password string
	Debug    bool
}

// rpcRequest represent a RCP request
type rpcRequest struct {
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int64       `json:"id"`
	JSONRpc string      `json:"jsonrpc"`
}

type RPCResponse struct {
	ID     int64           `json:"id"`
	Result json.RawMessage `json:"result"`
	Err    struct {
		Code    int16  `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// RPCCallJSON call json rpc
func RPCCallJSON(rpcInfo RPCInfo, method string, params interface{}, result interface{}) ([]byte, error) {
	if !strings.HasPrefix(rpcInfo.Host, "http") {
		rpcInfo.Host = "http://" + rpcInfo.Host
	}
	b, err := json.Marshal(rpcRequest{method, params, time.Now().UnixNano(), "1.0"})
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPost, rpcInfo.Host, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if rpcInfo.User != "" {
		req.SetBasicAuth(rpcInfo.User, rpcInfo.Password)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d, %s", resp.StatusCode, string(body))
	}

	if rpcInfo.Debug {
		log.Println("[dbg] rpc return", string(body))
	}
	var rpcRet RPCResponse
	if err = json.Unmarshal(body, &rpcRet); err != nil {
		return nil, err
	}
	if result != nil {
		if err = json.Unmarshal(rpcRet.Result, &result); err != nil {
			return nil, err
		}
	}
	return rpcRet.Result, nil
}
