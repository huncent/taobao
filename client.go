package taobao

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/smartwalle/going/request"
)

const (
	TAO_BAO_OPEN_API_URL = "http://gw.api.taobao.com/router/rest"
)

var (
	appKey    string
	appSecret string
)

func UpdateKey(appKey, appSecret string) {
	appKey = appKey
	appSecret = appSecret
}

func Request(param ITaoBaoParam) (results interface{}, err error) {
	results, err = RequestWithKey(appKey, appSecret, param)
	return results, err
}

func RequestWithKey(appKey, appSecret string, param ITaoBaoParam) (result ResultResp, err error) {
	var p = url.Values{}
	var keys = make([]string, 6, 6)

	p.Add("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	p.Add("format", "json")
	p.Add("v", "2.0")
	p.Add("sign_method", "md5")
	p.Add("app_key", appKey)
	p.Add("method", param.APIName())

	keys[0] = "timestamp"
	keys[1] = "format"
	keys[2] = "v"
	keys[3] = "sign_method"
	keys[4] = "app_key"
	keys[5] = "method"

	if len(param.ExtJSONParamName()) > 0 {
		p.Add(param.ExtJSONParamName(), param.ExtJSONParamValue())
		keys = append(keys, param.ExtJSONParamName())
	}

	var ps = param.Params()
	if ps != nil {
		for key, value := range ps {
			p.Add(key, value)
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	p.Add("sign", sign(appSecret, keys, p))

	respBodyByte, _ := request.Request("POST", TAO_BAO_OPEN_API_URL, p)
	err = json.Unmarshal(respBodyByte, &result)
	return result, err
	//return results, err
}

func sign(appSecret string, keys []string, param url.Values) (s string) {
	for _, key := range keys {
		s = s + key + param.Get(key)
	}
	s = fmt.Sprintf("%s%s%s", appSecret, s, appSecret)
	var m = md5.New()
	m.Write([]byte(s))
	s = strings.ToUpper(hex.EncodeToString(m.Sum(nil)))
	return s
}
