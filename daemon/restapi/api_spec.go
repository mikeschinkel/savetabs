// Package restapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.1-0.20240325090356-a14414f04fdd DO NOT EDIT.
package restapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RW33PbNgz+Vzhuj6rlNtse/NZuu8a3H8ml6e2hywMtwRZbiuRAyImX0/++I2jJlqXM",
	"6d3yJhHghw8fIECPsnC1dxYsBbl4lF6hqoEA+e09usbf7jz8oWqIByWEArUn7axcyHgq3FpsopugnQeZ",
	"SR0tXlElM2n5mqQOIZMIfzcaoZQLwgYyGYoKahWh+fpCBkJtN7Jt2+gcvLMBmMsl1eZmfxDfC2cJLPFV",
	"eKC8otrElwMiPKjamwj6VzOfXxSl3vIDXIIxLhOXt7//9k0y5b1NZlNMhnnHi8IhA4hgtfdAAoEatFCK",
	"0BQFhLBujNnJNpMfLTx4KAjKXxAdnnBX3htdqIicfw4R/jiF7xDWciG/zQ9VypM15Altgl7TBxSw9+mE",
	"ZiV7Gh6dBySdBC5cOVFjdhZsy+TaYa1ILqS2dPHmIJW2BBvAmG0NIajNk0CdeazycW98kvuAnftdm8ll",
	"OQa9Rl0r3IkvsBO6BEt6rQFPmP74/STTZXm1+gwFjZXQ5TntlyUz3oO6hHOEmWAI6nAe6aq/vYdTiIo7",
	"5waCa7CApxg+I8kGzdTHNUG9C/anpoq/e+5OY67WcvHpv3PoebbZKdFNh3TCIEuW5XPz2HST6FnZ3B3l",
	"E/qEnl+UsRaj6sSo2q4d89HEY+aD2sKtWgXxs4LaWfH2eikzuQUMqVtfz+azeYRyHqzyWi7kBR9lPDCZ",
	"Vl6BMlT9w+IBjVv+pwqKL4IqEMlTBFLUhDiH+/glx59JjoQ8XKLS8j3Q5R7+ZLq+mc/Hoa5+jWR/mF+M",
	"TR8At7oA8dGqrdJGrQykMdPU8YOMQzKxY7ps4gmdr9DdpwG+z27MkGrzLnlNk5yqXO+XD/bEkNINj2jW",
	"jkf32iG/pGhi5cpdJmC2mfGp0YGG2y0IFYSy6XI0H+XFXq/YK3/sFl6bjsO5dPs9G7plu2/YbLCQn/gM",
	"Dy75cGG3dy8n4EAdFqaXVB3/ELBCCKrcnWtpPWpgoYPgq4KcCIBbmGzpm4T+4h0d42gLIRw3NXZjJr/X",
	"VL061Nu7MFHwaxdoajKl7QeB3rly97/9IkxFaoerNv6GtSPpXn8Vha9adGF6lg6LcVuB6JUV94AgCgRF",
	"UEr2XavG0PlWPv31GtbzbVkGUTeGtDdH4TIBqqhELOe+kXnOt23b/hsAAP//LN3/lC4LAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
