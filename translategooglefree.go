package translategooglefree

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/robertkrimen/otto"
)

// javascript "encodeURI()"
// so we embed js to our golang programm
func encodeURI(s string) (string, error) {
	eUri := `eUri = encodeURI(sourceText);`
	vm := otto.New()
	err := vm.Set("sourceText", s)
	if err != nil {
		return "err", errors.New("Error setting js variable")
	}
	_, err = vm.Run(eUri)
	if err != nil {
		return "err", errors.New("Error executing jscript")
	}
	val, err := vm.Get("eUri")
	if err != nil {
		return "err", errors.New("Error getting variable value from js")
	}
	v, err := val.ToString()
	if err != nil {
		return "err", errors.New("Error converting js var to string")
	}
	return v, nil
}

func Translate(source, sourceLang, targetLang string) (string, string, error) {
	var text []string
	var result []interface{}

	/*encodedSource, err := encodeURI(source)
	if err != nil {
		return "err", err
	}*/

	encodedSource := url.QueryEscape(source)

	uri := "https://translate.googleapis.com/translate_a/single?client=gtx&sl=" +
		sourceLang + "&tl=" + targetLang + "&dt=t&q=" + encodedSource

	r, err := http.Get(uri)
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
