package translategooglefree

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	Proxy   string
	TimeOut int
}

type TranslateGoogle struct {
	Config *Config
	Client *http.Client
}

func NewClient(cfg *Config) *TranslateGoogle {
	transport := &http.Transport{}
	proxy := strings.TrimSpace(cfg.Proxy)
	if strings.HasPrefix(proxy, "http") {
		proxyUrl, _ := url.Parse(proxy)
		transport.Proxy = http.ProxyURL(proxyUrl)                         // set proxy
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} // skip verify
	}

	timeOut := cfg.TimeOut
	if timeOut == 0 {
		timeOut = 20
	}

	return &TranslateGoogle{
		Config: cfg,
		Client: &http.Client{
			Timeout:   time.Duration(timeOut) * time.Second,
			Transport: transport,
		},
	}
}

func (tsl *TranslateGoogle) Translate(source, sourceLang, targetLang string) (string, string, error) {
	var text []string
	var result []interface{}

	/*encodedSource, err := encodeURI(source)
	if err != nil {
		return "err", err
	}*/

	encodedSource := url.QueryEscape(source)

	uri := "https://translate.googleapis.com/translate_a/single?client=gtx&sl=" +
		sourceLang + "&tl=" + targetLang + "&dt=t&q=" + encodedSource

	//r, err := http.Get(uri)

	r, err := tsl.Client.Get(uri)
	if err != nil {
		return "err", "", errors.New("Error getting translate.googleapis.com")
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "err", "", errors.New("Error reading response body")
	}

	bReq := strings.Contains(string(body), `<title>Error 400 (Bad Request)`)
	if bReq {
		return "err", "", errors.New("Error 400 (Bad Request)")
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "err", "", errors.New("Error unmarshaling data")
	}

	if len(result) > 0 {
		inner := result[0]
		if inner == nil {
			return "err", "", errors.New("Error result data")
		}
		for _, slice := range inner.([]interface{}) {
			for _, translatedText := range slice.([]interface{}) {
				text = append(text, fmt.Sprintf("%v", translatedText))
				break
			}
		}
		cText := strings.Join(text, "")

		lang := "auto"
		if len(result) > 3 {
			lang = result[2].(string)
		}
		return cText, lang, nil
	} else {
		return "err", "", errors.New("No translated data in responce")
	}
}
