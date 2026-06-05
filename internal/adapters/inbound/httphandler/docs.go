package httphandler

import _ "embed"

//go:embed files/openapi.yaml
var openAPISpec []byte

//go:embed files/swagger_ui.html
var swaggerUIHTML []byte
