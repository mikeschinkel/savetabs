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

	"H4sIAAAAAAAC/7xYXXPbNhb9K1jsPrQzTOQ47s6s3pymSby164xttQ/ZzA5EXlGoQYAFQNtcj/77zr0g",
	"JZGEbMqj5E0ivs69OPfj4JGnpiiNBu0dnz7yUlhRgAdL/04VWH/hcvydgUutLL00mk/DCCvAOZED84Zl",
	"0pVK1Ew49unm4pwnXOK8vyqwNU+4FgXwKS9czhPu0iUUAjf1dYmfnbdS53yVcOdrhV8Wxhb4n865oVl9",
	"CPiVmQUTL4Pi67IDBXRV8OkXTtvRmoXBCVWagnM84ffCakSZcLDWWP41eR7+O2NuC2FvP0jlwQ5taMfZ",
	"giawhbHsXOpbPA8eSmUy4FNvK4jbMO9YID0ULurV5oOwVtQxmD8LD7mx9S6Y7ThTiI3Na5aGLxJGIk0P",
	"hPQXdH2UkTSyPw3A2v0Y+dGaqkTu7XIW3R/dpGA0mRF/R7kp9wfy069Q3xub7cLYDL+QdrcHAolHnmVD",
	"dGfvCZFfAhGuvblS+OUGBI78V2Y84Rb+qqSFrIW8gfYPCws+5X+fbJLcJIy6yVnGVw2GmVVDELOr8+YS",
	"n8FQWfUkiJ5bVgm/AF2deSgilnsomBJzUHR2AbpiJgxGARTtTvsC8OIJakx+F6oCVoAXL2RI0WGIyDKJ",
	"uwv12ZoSrMe08QRhzPxPSH2MMTdiTiG1C3w7vklUXuQjMecHYvWNyHfDy1+C7DBJYYUccaXRLnj/ky/U",
	"VfMB/6dGe9CedoYHP1n6gsJiq0Q+iKKkHf9THR29TTN5Rz/gEyhlEkqzfwtDk/UYH5RJRNL1Cy5kxtIG",
	"zGlZluCZBV9ZDRlrSvCiUorMmml4KCH1kFHO72EXZalkKnDnyZ8Ot38cmRH+fX35W9gxArFaH8qgmdNe",
	"Cnkzlsc+W1kIW7NbqJnMQHu5kIDlBi9EeD7lUvt/nmw8JLWHHCwaeZZdhijApqwTNTIbmdoG0dTu6To0",
	"enqny/XqPrs27hpATInQ8QpNY10XvD2OuqCp5c+U+gi/tpPhF94c2E7/2iT9lzo24cbKXGqhmrKxtqSy",
	"kkd6Qi+9gmi8VmM2iF0k4v9D+iXlOiK9UpcLPv3yNHwye5X07c7bXQYAaSQwewRl87YzihednhXtPbi1",
	"IeNp2bV/wE08jdp3xBG8z6/FHdyIuWPvBRRGs9PPZzzhd2Bd4NSb10evj+h2S9CilHzK39KnhEouQZos",
	"QSi//B85DXykT15Cekt9S5jJnBe+cqhS1udndP5rTidZylPoYf4R/Kdm+16iPj46Gh51+SuC/eno7XDo",
	"GuydTIHNtLgTUom5gpCtqgLTEebbgI7g0hAl+0mQPruMC5Iv8ADToEBfSlez2VmjwZo2e2iXL9RpI6u2",
	"BeYOum6mTDbyD2k7ZjJqg9XXuAdj69fzJp162PPXzcX5qy3LyZaN36DNglG/hWz1Ur+FDLuv39Yy6Vu7",
	"YlMLgyuou4FtEkWtOm+m7WtXX3yNYEWvZRyxoqeGR52Rj5/clWYjFvQeEUas2OrvD06AK2rLKMlRuxY0",
	"UimwJCIjPOYbzHih0b2XftlIGKEaNSF1zlIr8ZdgUjNqd9kWE1YJL42L0OezcT3+YKUH59+ZrH6iEXx4",
	"dX9/T7R9VVkFGnuCLLz6pCbDCoUL0hCxj4OWvtWa1AN32vVBi71pNbtFdrP5oMpu7T6yKYuUvEF97bZB",
	"iHZ1YCbMykx4YII58HjhQR9ucgGK0+cSAUph/h0Yin/orC68yWOroFdjkF5s5PZ+eWu98LuEY9dYisPJ",
	"vEbmTx7bJwuyt6xiMVZ5usl39cyq9nlkX4PbdbvsPYhg26iToV67aR6P2L1wrCKmUtgcH735zueDBZZa",
	"2AA4GfYIvxn2c4NolfCTWMP3TmTsKiS7MCfS+X0wdi6zDHSYET3Isw+m0lmY8a9ID2v0QsmgM07eHMdU",
	"LaRGhxcd9kFIFew6efNT5BFNu6osjcXCcAGZFKxp5vjJ8XFsemkNan2qIb9oL33NfvgD5u9Pf/+Rk48X",
	"olL++YjpvxF0g+YadMZmV+cMVQK2MmiKN2wOzHmDCXMrcrCCvcrXAmV3YeqrmbHFaT/S9U8Zlej34/xe",
	"TwMurr/i0eAG4XCgGz3NMseKSnlZquaohIFIl6EBoeuj296+2cem8j6Z+8nfzUv1S7IgVutvmQSDqB+f",
	"AMcll+dSx4Gu7SP4TgyKuak8E9RKsB9mV+c/hvuyILL6Od0tByqbScdoKQa3A3sHUd19FXb/5rIbz5Ea",
	"nFsr79Xq/wEAAP//h38RMYEdAAA=",
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
