package middleware

import (
	"encoding/base64"
	"strings"

	"github.com/omriz/multigrok/backends"
)

// SEPERATOR is a sepertor indicating a backend.
const SEPERATOR = "__MULTIGROKBACKEND__"

// EncodeBackendAddress is used to encode a given address in base64 with a seperator.
func EncodeBackendAddress(addr string) string {
	return SEPERATOR + base64.StdEncoding.EncodeToString([]byte(addr)) + SEPERATOR
}

// DecodeBackendAddress is used to decode an encoded backend.
func DecodeBackendAddress(part string) (string, error) {
	encoded := strings.TrimSuffix(strings.TrimPrefix(part, SEPERATOR), SEPERATOR)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// CombineResults takes two backends.WebServiceResult and combines them into one.
// The path of the result (not the directory) will be suffixed by an escaped backend address
// in the format: SEPERATORbase64(Address)SEPERATOR
func CombineResults(responses map[string]backends.WebServiceResult) (backends.WebServiceResult, error) {
	qres := make(map[string][]backends.QueryResult)
	total := 0
	for uid, wres := range responses {
		total += wres.Resultcount
		for path, qr := range wres.Results {
			np := "/" + EncodeBackendAddress(uid) + path
			qres[np] = qr
		}
	}
	resp := backends.WebServiceResult{
		Resultcount: total,
		Results:     qres,
	}
	return resp, nil
}
