package phpserialize

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

func Encode(value interface{}) (result string, err error) {
	buf := new(bytes.Buffer)
	err = encodeValue(buf, value)
	if err == nil {
		result = buf.String()
	}
	return
}

func encodeValue(buf *bytes.Buffer, value interface{}) (err error) {

	switch t := value.(type) {
	default:
		err = fmt.Errorf("Unexpected type %T", t)
	case bool:
		buf.WriteString("b")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		if t {
			buf.WriteString("1")
		} else {
			buf.WriteString("0")
		}
		buf.WriteRune(VALUES_SEPARATOR)
	case nil:
		buf.WriteString("N")
		buf.WriteRune(VALUES_SEPARATOR)
	case int, int64, int32, int16, int8:
		buf.WriteString("i")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		strValue := fmt.Sprintf("%v", t)
		buf.WriteString(strValue)
		buf.WriteRune(VALUES_SEPARATOR)
	case float32:
		buf.WriteString("d")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		strValue := strconv.FormatFloat(float64(t), 'f', -1, 64)
		buf.WriteString(strValue)
		buf.WriteRune(VALUES_SEPARATOR)
	case float64:
		buf.WriteString("d")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		strValue := strconv.FormatFloat(float64(t), 'f', -1, 64)
		buf.WriteString(strValue)
		buf.WriteRune(VALUES_SEPARATOR)
	case string:
		buf.WriteString("s")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		encodeString(buf, t)
		buf.WriteRune(VALUES_SEPARATOR)
	case map[interface{}]interface{}:
		buf.WriteString("a")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		var int64Key bool
		var maxNumber int64 = math.MinInt64
		for k, _ := range t {
			if reflect.TypeOf(k).String() == "int64" {
				if k.(int64) > maxNumber {
					maxNumber = k.(int64)
				}
				int64Key = true
			} else {
				int64Key = false
				break
			}
		}
		if int64Key {
			err = encodeArrayIntKey(buf, t, maxNumber)
		} else {
			err = encodeArrayCore(buf, t)
		}

	case *PhpObject:
		buf.WriteString("O")
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		encodeString(buf, t.GetClassName())
		buf.WriteRune(TYPE_VALUE_SEPARATOR)
		err = encodeArrayCore(buf, t.GetMembers())
	}
	return
}

func encodeString(buf *bytes.Buffer, strValue string) {
	valLen := strconv.Itoa(len(strValue))
	buf.WriteString(valLen)
	buf.WriteRune(TYPE_VALUE_SEPARATOR)
	buf.WriteRune('"')
	buf.WriteString(strValue)
	buf.WriteRune('"')
}

func encodeArrayCore(buf *bytes.Buffer, arrValue map[interface{}]interface{}) (err error) {
	valLen := strconv.Itoa(len(arrValue))
	buf.WriteString(valLen)
	buf.WriteRune(TYPE_VALUE_SEPARATOR)

	buf.WriteRune('{')
	for k, v := range arrValue {
		if intKey, _err := strconv.Atoi(fmt.Sprintf("%v", k)); _err == nil {
			if err = encodeValue(buf, intKey); err != nil {
				break
			}
		} else {
			if err = encodeValue(buf, k); err != nil {
				break
			}
		}
		if err = encodeValue(buf, v); err != nil {
			break
		}
	}
	buf.WriteRune('}')
	return err
}

func encodeArrayIntKey(buf *bytes.Buffer, arrValue map[interface{}]interface{}, maxKey int64) (err error) {
	valLen := strconv.Itoa(len(arrValue))
	buf.WriteString(valLen)
	buf.WriteRune(TYPE_VALUE_SEPARATOR)
	buf.WriteRune('{')

	var mapInt64 map[int64]interface{} = make(map[int64]interface{})
	for k, v := range arrValue {
		mapInt64[k.(int64)] = v
	}

	for k := int64(0); k <= maxKey; k++ {
		if v, ok := mapInt64[k]; ok {
			if intKey, _err := strconv.Atoi(fmt.Sprintf("%v", k)); _err == nil {
				if err = encodeValue(buf, intKey); err != nil {
					break
				}
			} else {
				if err = encodeValue(buf, k); err != nil {
					break
				}
			}
			if err = encodeValue(buf, v); err != nil {
				break
			}
		}
	}

	buf.WriteRune('}')
	return err
}
