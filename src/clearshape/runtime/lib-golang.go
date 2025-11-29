package WHATPACKAGENAME
import (
	"encoding/base64"
	"encoding/json"
	"strconv"
)

// start rest of runtime lib

type null struct {}

type _errstr struct {
	ErrorMessage string
	Cause_ *_err
}

type _err interface {
	error
	Cause() _err
}

func(e *_errstr) Error() string {
	return e.ErrorMessage
}

func(e *_errstr) Cause() _err {
	return *e.Cause_
}

func _newerr(s string, e _err) _err {
	return &_errstr{
		ErrorMessage: s,
		Cause_: &e,
	}
}

func _parseString(a interface{}) (string, _err) {
	if str, ok := a.(string); ok {
		return str, nil
	} else {
		return "", _newerr("expected string", nil)
	}
}

func _parseBoolean(a interface{}) (bool, _err) {
	if b, ok := a.(bool); ok {
		return b, nil
	} else {
		return false, _newerr("expected string", nil)
	}
}

func _parseUint64(a interface{}) (uint64, _err) {
	if num, ok := a.(json.Number); ok {
		u64, err := strconv.ParseUint(num.String(), 10, 64)
		if err != nil {
			return 0, _newerr("error when parsing json number into uint64", _newerr(err.Error(), nil))
		}
		return u64, nil
	} else {
		return 0, _newerr("expecting json number", nil)
	}
}

func _parseInt64(a interface{}) (int64, _err) {
	if num, ok := a.(json.Number); ok {
		i64, err := num.Int64()
		if err != nil {
			return 0, _newerr("error when parsing json number into int64", _newerr(err.Error(), nil))
		}
		return i64, nil
	} else {
		return 0, _newerr("expecting json number", nil)
	}
}

func _parseFloat64(a interface{}) (float64, _err) {
	if num, ok := a.(json.Number); ok {
		f64, err := num.Float64()
		if err != nil {
			return 0, _newerr("error when parsing json number into uint64", _newerr(err.Error(), nil))
		}
		return f64, nil
	} else {
		return 0, _newerr("expecting json number", nil)
	}
}

func _parseNull(a interface{}) (any, _err) {
	switch a.(type) {
	case nil: 
		return nil, nil
	default:
		return nil, _newerr("expecting null", nil)
	}
}

func _parseBinary(a interface{}) ([]byte, _err) {
	if str, ok := a.(string); ok {
		strictDecoder := base64.StdEncoding.Strict()
		decoded, err := strictDecoder.DecodeString(str)
		if err != nil {
			return []byte{}, _newerr("base64 decode error", _newerr(err.Error(), nil))
		}
		return decoded, nil
	} else {
		return []byte{}, _newerr("expecting string while parsing binary data", nil)
	}
}

func _parseJsonObject(a interface{}) (map[string]interface{}, _err) {
	if obj, ok := a.(map[string]interface{}); ok {
		return obj, nil
	} else {
		return nil, _newerr("expected json object", nil)
	}
}

func _parseJsonList(a interface{}) ([]interface{}, _err) {
	if obj, ok := a.([]interface{}); ok {
		return obj, nil
	} else {
		return nil, _newerr("expected json list", nil)
	}
}