package nginx_parser

import (
	"strings"
)

func ParseRequest(log string, visitor *Visitor) *Visitor {
	parts := strings.Split(log, " ")

	if len(parts) < 3 {
		return visitor
	}

	// remove " from the first and last part
	parts[0] = strings.Trim(parts[0], "\"")
	parts[len(parts)-1] = strings.Trim(parts[len(parts)-1], "\"")

	method := parts[0]
	uri := parts[1]
	protocol := parts[2]

	visitor.RequestMethod = method
	visitor.RequestURI = uri
	visitor.ResponseStatus = protocol
	return visitor
}
