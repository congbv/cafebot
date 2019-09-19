package chat

import (
	"bytes"
	"strings"
)

const (
	intrSeparator   = '/'
	intrOpSeparator = '?'
)

func newIntrData(
	e intrEndpoint,
	intrData string,
	op orderOp,
	opval string,
) string {
	buf := new(bytes.Buffer)
	buf.WriteString(string(e))
	if intrData != "" {
		buf.WriteRune(intrSeparator)
		buf.WriteString(intrData)
	}
	if op != "" && opval != "" {
		buf.WriteRune(intrOpSeparator)
		buf.WriteString(string(op))
		buf.WriteRune('=')
		buf.WriteString(opval)
	}
	return buf.String()
}

func splitIntrOpData(reqdata string) (string, string, error) {
	splited := strings.Split(reqdata, string(intrOpSeparator))
	opdata := ""
	if len(splited) > 1 {
		opdata = splited[1]
	}
	return splited[0], opdata, nil
}

func splitEndpointIntrData(reqdata string) (intrEndpoint, string) {
	parts := strings.Split(reqdata, string(intrSeparator))
	if len(parts) == 0 {
		return "", ""
	}
	var intrData string
	if len(parts) > 1 {
		intrData = parts[1]
	}
	return intrEndpoint(parts[0]), intrData
}
