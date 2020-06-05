package additionally

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"sort"
	"strconv"

	"BizCoinWebSocket/config"
	"github.com/fatih/structs"
)

func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsValid is vk sign checker
func IsValid(query config.Data, secret string) bool {
	workQuery := structs.New(query)
	vkSubset := map[string]interface{}{}
	var keyQuery []string

	// get data keys
	for _, field := range workQuery.Fields() {
		if field.IsZero() {
			return false
		}

		if field.Tag("json")[0:3] == "vk_" {
			vkSubset[field.Tag("json")] = field.Value()
			keyQuery = append(keyQuery, field.Tag("json"))
		}
	}

	sort.Slice(keyQuery, func(i, j int) bool {
		return keyQuery[int(i)] < keyQuery[int(j)]
	})

	urlParams := url.Values{}
	for key := range keyQuery {
		urlParams.Add(strconv.Itoa(key), vkSubset["key"].(string))
	}
	urlEncode := urlParams.Encode()

	decodedHashCode := []byte(reverse(computeHmac256(urlEncode, secret)))
	for i := 0; i < len(decodedHashCode); i++ {
		if decodedHashCode[i] == '+' {
			decodedHashCode[i] = '-'
		}
		if decodedHashCode[i] == '/' {
			decodedHashCode[i] = '_'
		}
	}

	// fmt.Println(string(decodedHashCode))

	return query.Sign == string(decodedHashCode)
}
