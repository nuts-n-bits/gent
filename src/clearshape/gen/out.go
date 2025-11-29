package main

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
)

type Csres01A3enu struct {
    OneofStruct *string
    OneofStruct2 *string
}

type Csres01B1c struct {
    Field0 int64
    Field1 string
}

type C map[string]string

type A struct {
    Session string
    Name []Csres01A4name
    Name2 []string
    Time2 *string
    Time4 *Csres01A5time4
    TechnicalIdentifier map[string]int64
    Another *B
    Str Csres01A3str
    Enu Csres01A3enu
}

type Csres01A5time4 struct {
    Y int64
    M int64
    D int64
}

type Csres01A3str struct {
    Struct string
    Struct2 string
}

type B struct {
    S A
    V Csres01B1v
    C Csres01B1c
    D Csres01B1d
    Map map[string]A
    Bin map[string][][]byte
}

type Csres01B1v struct {
    Field0 string
    Field1 string
}

type Csres01B1d struct {
    Field0 string
}

type D []string

type Csres01A4name struct {
    S string
}

func (a *A) fromJsonCore (b interface{}) _err {
    bMap, err := _parseJsonObject(b)
    if err != nil {
        return err
    }
    var session string
    var name []Csres01A4name
    var name2 []string
    var time2 *string
    var time4 *Csres01A5time4
    var technicalIdentifier map[string]int64
    var another *B
    var str Csres01A3str
    var enu Csres01A3enu
    if v, ok := bMap["s"]; ok {
        // parse v into session here
        parsed, err := _parseString(v)
        if err != nil {
            return _newerr("error when parsing field Session (wire name '"+"s"+"')", err)
        }
        session = parsed
    } else {
        return _newerr("missing required field Session (wire name '"+"s"+"')", nil)
    }
    if v, ok := bMap["n"]; ok {
        // parse v into name here
        t0 := func (a interface{}) ([]Csres01A4name, _err) {
            list, err := _parseJsonList(a)
            if err != nil {
                return []Csres01A4name{}, _newerr("error parsing list", err)
            }
            ret := []Csres01A4name{}
            for i, v := range list {
                parsed := Csres01A4name{}
                err := parsed.fromJsonCore(v)
                if err != nil {
                    iStr := strconv.Itoa(i)
                    return []Csres01A4name{}, _newerr("error parsing item " + iStr + " inside list", err)
                }
                ret = append(ret, parsed)
            }
            return ret, nil
        }
        parsed, err := t0(v)
        if err != nil {
            return _newerr("error when parsing field Name (wire name '"+"n"+"')", err)
        }
        name = parsed
    } else {
        return _newerr("missing required field Name (wire name '"+"n"+"')", nil)
    }
    if v, ok := bMap["n2"]; ok {
        // parse v into name2 here
        t0 := func (a interface{}) ([]string, _err) {
            list, err := _parseJsonList(a)
            if err != nil {
                return []string{}, _newerr("error parsing list", err)
            }
            ret := []string{}
            for i, v := range list {
                parsed, err := _parseString(v)
                if err != nil {
                    iStr := strconv.Itoa(i)
                    return []string{}, _newerr("error parsing item " + iStr + " inside list", err)
                }
                ret = append(ret, parsed)
            }
            return ret, nil
        }
        parsed, err := t0(v)
        if err != nil {
            return _newerr("error when parsing field Name2 (wire name '"+"n2"+"')", err)
        }
        name2 = parsed
    } else {
        return _newerr("missing required field Name2 (wire name '"+"n2"+"')", nil)
    }
    if v, ok := bMap["t2"]; ok {
        // parse v into time2 here
        parsed, err := _parseString(v)
        if err != nil {
            return _newerr("error when parsing field Time2 (wire name '"+"t2"+"')", err)
        }
        time2 = &parsed
    }
    if v, ok := bMap["t4"]; ok {
        // parse v into time4 here
        parsed := Csres01A5time4{}
        err := parsed.fromJsonCore(v)
        if err != nil {
            return _newerr("error when parsing field Time4 (wire name '"+"t4"+"')", err)
        }
        time4 = &parsed
    }
    if v, ok := bMap["ti"]; ok {
        // parse v into technicalIdentifier here
        t0 := func (a interface{}) (map[string]int64, _err) {
            map_, err := _parseJsonObject(a)
            if err != nil {
                return map[string]int64{}, _newerr("error parsing map", err)
            }
            ret := map[string]int64{}
            for k, v := range map_ {
                parsed, err := _parseInt64(v)
                if err != nil {
                    return map[string]int64{}, _newerr("error parsing key '" + k + "' inside map", err)
                }
                ret[k] = parsed
            }
            return ret, nil
        }
        parsed, err := t0(v)
        if err != nil {
            return _newerr("error when parsing field TechnicalIdentifier (wire name '"+"ti"+"')", err)
        }
        technicalIdentifier = parsed
    } else {
        return _newerr("missing required field TechnicalIdentifier (wire name '"+"ti"+"')", nil)
    }
    if v, ok := bMap["a"]; ok {
        // parse v into another here
        parsed := B{}
        err := parsed.fromJsonCore(v)
        if err != nil {
            return _newerr("error when parsing field Another (wire name '"+"a"+"')", err)
        }
        another = &parsed
    }
    if v, ok := bMap["str"]; ok {
        // parse v into str here
        parsed := Csres01A3str{}
        err := parsed.fromJsonCore(v)
        if err != nil {
            return _newerr("error when parsing field Str (wire name '"+"str"+"')", err)
        }
        str = parsed
    } else {
        return _newerr("missing required field Str (wire name '"+"str"+"')", nil)
    }
    if v, ok := bMap["enu"]; ok {
        // parse v into enu here
        parsed := Csres01A3enu{}
        err := parsed.fromJsonCore(v)
        if err != nil {
            return _newerr("error when parsing field Enu (wire name '"+"enu"+"')", err)
        }
        enu = parsed
    } else {
        return _newerr("missing required field Enu (wire name '"+"enu"+"')", nil)
    }
    a.Session = session
    a.Name = name
    a.Name2 = name2
    a.Time2 = time2
    a.Time4 = time4
    a.TechnicalIdentifier = technicalIdentifier
    a.Another = another
    a.Str = str
    a.Enu = enu
    return nil
}

func (a *Csres01A5time4) fromJsonCore (b interface{}) _err {
    bMap, err := _parseJsonObject(b)
    if err != nil {
        return err
    }
    var y int64
    var m int64
    var d int64
    if v, ok := bMap["y"]; ok {
        // parse v into y here
        parsed, err := _parseInt64(v)
        if err != nil {
            return _newerr("error when parsing field Y (wire name '"+"y"+"')", err)
        }
        y = parsed
    } else {
        return _newerr("missing required field Y (wire name '"+"y"+"')", nil)
    }
    if v, ok := bMap["m"]; ok {
        // parse v into m here
        parsed, err := _parseInt64(v)
        if err != nil {
            return _newerr("error when parsing field M (wire name '"+"m"+"')", err)
        }
        m = parsed
    } else {
        return _newerr("missing required field M (wire name '"+"m"+"')", nil)
    }
    if v, ok := bMap["d"]; ok {
        // parse v into d here
        parsed, err := _parseInt64(v)
        if err != nil {
            return _newerr("error when parsing field D (wire name '"+"d"+"')", err)
        }
        d = parsed
    } else {
        return _newerr("missing required field D (wire name '"+"d"+"')", nil)
    }
    a.Y = y
    a.M = m
    a.D = d
    return nil
}

func (a *Csres01A3str) fromJsonCore (b interface{}) _err {
    bMap, err := _parseJsonObject(b)
    if err != nil {
        return err
    }
    var struct_ string
    var struct2 string
    if v, ok := bMap["struct"]; ok {
        // parse v into struct_ here
        parsed, err := _parseString(v)
        if err != nil {
            return _newerr("error when parsing field Struct (wire name '"+"struct"+"')", err)
        }
        struct_ = parsed
    } else {
        return _newerr("missing required field Struct (wire name '"+"struct"+"')", nil)
    }
    if v, ok := bMap["struct2"]; ok {
        // parse v into struct2 here
        parsed, err := _parseString(v)
        if err != nil {
            return _newerr("error when parsing field Struct2 (wire name '"+"struct2"+"')", err)
        }
        struct2 = parsed
    } else {
        return _newerr("missing required field Struct2 (wire name '"+"struct2"+"')", nil)
    }
    a.Struct = struct_
    a.Struct2 = struct2
    return nil
}

func (a *B) fromJsonCore (b interface{}) _err {
    bMap, err := _parseJsonObject(b)
    if err != nil {
        return err
    }
    var s A
    var v Csres01B1v
    var c Csres01B1c
    var d Csres01B1d
    var map_ map[string]A
    var bin map[string][][]byte
    if v, ok := bMap["s"]; ok {
        // parse v into s here
        parsed := A{}
        err := parsed.fromJsonCore(v)
        if err != nil {
            return _newerr("error when parsing field S (wire name '"+"s"+"')", err)
        }
        s = parsed
    } else {
        return _newerr("missing required field S (wire name '"+"s"+"')", nil)
    }
    if v, ok := bMap["v"]; ok {
        // parse v into v here
        parsed := Csres01B1v{}
        err := parsed.fromJsonCore(v)
        if err != nil {
            return _newerr("error when parsing field V (wire name '"+"v"+"')", err)
        }
        v = parsed
    } else {
        return _newerr("missing required field V (wire name '"+"v"+"')", nil)
    }
    if v, ok := bMap["c"]; ok {
        // parse v into c here
        parsed := Csres01B1c{}
        err := parsed.fromJsonCore(v)
        if err != nil {
            return _newerr("error when parsing field C (wire name '"+"c"+"')", err)
        }
        c = parsed
    } else {
        return _newerr("missing required field C (wire name '"+"c"+"')", nil)
    }
    if v, ok := bMap["d"]; ok {
        // parse v into d here
        parsed := Csres01B1d{}
        err := parsed.fromJsonCore(v)
        if err != nil {
            return _newerr("error when parsing field D (wire name '"+"d"+"')", err)
        }
        d = parsed
    } else {
        return _newerr("missing required field D (wire name '"+"d"+"')", nil)
    }
    if v, ok := bMap["map"]; ok {
        // parse v into map_ here
        t0 := func (a interface{}) (map[string]A, _err) {
            map_, err := _parseJsonObject(a)
            if err != nil {
                return map[string]A{}, _newerr("error parsing map", err)
            }
            ret := map[string]A{}
            for k, v := range map_ {
                parsed := A{}
                err := parsed.fromJsonCore(v)
                if err != nil {
                    return map[string]A{}, _newerr("error parsing key '" + k + "' inside map", err)
                }
                ret[k] = parsed
            }
            return ret, nil
        }
        parsed, err := t0(v)
        if err != nil {
            return _newerr("error when parsing field Map (wire name '"+"map"+"')", err)
        }
        map_ = parsed
    } else {
        return _newerr("missing required field Map (wire name '"+"map"+"')", nil)
    }
    if v, ok := bMap["bin"]; ok {
        // parse v into bin here
        t0 := func (a interface{}) (map[string][][]byte, _err) {
            map_, err := _parseJsonObject(a)
            if err != nil {
                return map[string][][]byte{}, _newerr("error parsing map", err)
            }
            ret := map[string][][]byte{}
            for k, v := range map_ {
                t0 := func (a interface{}) ([][]byte, _err) {
                    list, err := _parseJsonList(a)
                    if err != nil {
                        return [][]byte{}, _newerr("error parsing list", err)
                    }
                    ret := [][]byte{}
                    for i, v := range list {
                        parsed, err := _parseBinary(v)
                        if err != nil {
                            iStr := strconv.Itoa(i)
                            return [][]byte{}, _newerr("error parsing item " + iStr + " inside list", err)
                        }
                        ret = append(ret, parsed)
                    }
                    return ret, nil
                }
                parsed, err := t0(v)
                if err != nil {
                    return map[string][][]byte{}, _newerr("error parsing key '" + k + "' inside map", err)
                }
                ret[k] = parsed
            }
            return ret, nil
        }
        parsed, err := t0(v)
        if err != nil {
            return _newerr("error when parsing field Bin (wire name '"+"bin"+"')", err)
        }
        bin = parsed
    } else {
        return _newerr("missing required field Bin (wire name '"+"bin"+"')", nil)
    }
    a.S = s
    a.V = v
    a.C = c
    a.D = d
    a.Map = map_
    a.Bin = bin
    return nil
}

func (a *Csres01B1v) fromJsonCore (b interface{}) _err {

}

func (a *Csres01B1d) fromJsonCore (b interface{}) _err {
    UNIMPLEMENTED
}

func (a *D) fromJsonCore (b interface{}) _err {
    listParser := func (a interface{}) ([]string, _err) {
        list, err := _parseJsonList(a)
        if err != nil {
            return []string{}, _newerr("error parsing list", err)
        }
        ret := []string{}
        for i, v := range list {
            parsed, err := _parseString(v)
            if err != nil {
                iStr := strconv.Itoa(i)
                return []string{}, _newerr("error parsing item " + iStr + " inside list", err)
            }
            ret = append(ret, parsed)
        }
        return ret, nil
    }
    parsed, err := listParser(b)
    if err != nil {
        return _newerr("error while parsing list", err)
    }
    *a = parsed
    return nil
}

func (a *Csres01A4name) fromJsonCore (b interface{}) _err {
    bMap, err := _parseJsonObject(b)
    if err != nil {
        return err
    }
    var s string
    if v, ok := bMap["s"]; ok {
        // parse v into s here
        parsed, err := _parseString(v)
        if err != nil {
            return _newerr("error when parsing field S (wire name '"+"s"+"')", err)
        }
        s = parsed
    } else {
        return _newerr("missing required field S (wire name '"+"s"+"')", nil)
    }
    a.S = s
    return nil
}

func (a *Csres01A3enu) fromJsonCore (b interface{}) _err {
    UNIMPLEMENTED
}

func (a *Csres01B1c) fromJsonCore (b interface{}) _err {
    UNIMPLEMENTED
}

func (a *C) fromJsonCore (b interface{}) _err {
    mapParser := func (a interface{}) (map[string]string, _err) {
        map_, err := _parseJsonObject(a)
        if err != nil {
            return map[string]string{}, _newerr("error parsing map", err)
        }
        ret := map[string]string{}
        for k, v := range map_ {
            parsed, err := _parseString(v)
            if err != nil {
                return map[string]string{}, _newerr("error parsing key '" + k + "' inside map", err)
            }
            ret[k] = parsed
        }
        return ret, nil
    }
    parsed, err := mapParser(b)
    if err != nil {
        return _newerr("error while parsing map", err)
    }
    *a = parsed
    return nil
}

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
