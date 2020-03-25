package devtools4chains

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// RPCInfo .
type RPCInfo struct {
	Host     string
	User     string
	Password string
}

// RPCCallJSON call json rpc
func RPCCallJSON(rpcInfo RPCInfo, method string, params interface{}, result interface{}) ([]byte, error) {
	if !strings.HasPrefix(rpcInfo.Host, "http") {
		rpcInfo.Host = "http://" + rpcInfo.Host
	}
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	bodyReader := bytes.NewReader(b)
	req, err := http.NewRequest(http.MethodPost, rpcInfo.Host, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

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
	if result != nil {
		if err = json.Unmarshal(body, result); err != nil {
			return nil, err
		}
	}
	return body, nil
}
