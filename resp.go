package redisprotocol

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const CRLF = "\r\n"

type RedisInt struct {
	Value int
	// RedisValue string
}

type RedisStr struct {
	Value string
	// RedisValue string
}

type RedisError struct {
	Value string
	// RedisValue string
}

type RedisBulkStr struct {
	Value string
	null  bool
	// RedisValue string
}

type RedisArray struct {
	Value []interface{}
	null  bool
	// RedisValue string
}

func EncodeStr(value RedisStr) (string, error) {
	return fmt.Sprintf("+%s%s", value.Value, CRLF), nil
}

func EncodeInt(value RedisInt) (string, error) {
	return fmt.Sprintf(":%d%s", value.Value, CRLF), nil
}

func EncodeErr(value RedisError) (string, error) {
	return fmt.Sprintf("-%s%s", value.Value, CRLF), nil
}

func EncodeBulkStr(value RedisBulkStr) (string, error) {
	if value.null {
		return fmt.Sprintf("$-1%s", CRLF), nil
	}
	return fmt.Sprintf("$%d%s%s%s", len(value.Value), CRLF, value.Value, CRLF), nil
}

func EncodeArray(value RedisArray) (string, error) {
	if value.null {
		return fmt.Sprintf("*-1%s", CRLF), nil
	}
	if len(value.Value) == 0 {
		return fmt.Sprintf("*0%s", CRLF), nil
	}
	res := fmt.Sprintf("*%d%s", len(value.Value), CRLF)
	for _, el := range value.Value {
		elStr, _ := Encode(el)
		res += elStr
	}
	res += CRLF
	return res, nil
}

func Encode(value interface{}) (string, error) {
	switch value.(type) {
	case RedisStr:
		return EncodeStr(value.(RedisStr))
	case RedisInt:
		return EncodeInt(value.(RedisInt))
	case RedisError:
		return EncodeErr(value.(RedisError))
	case RedisBulkStr:
		return EncodeBulkStr(value.(RedisBulkStr))
	case RedisArray:
		return EncodeArray(value.(RedisArray))
	default:
		fmt.Printf("The type of v is unknown\n")
		return "", errors.New("invalid type")
	}
}

func DecodeStr(s string) (RedisStr, error, int) {
	crlfpos := strings.Index(s, CRLF)
	if crlfpos < 0 {
		return RedisStr{}, errors.New(fmt.Sprintf("invalid string %s", s)), -2
	}
	return RedisStr{Value: s[1:crlfpos]}, nil, crlfpos + len(CRLF)
}

func DecodeError(s string) (RedisError, error, int) {
	crlfpos := strings.Index(s, CRLF)
	if crlfpos < 0 {
		return RedisError{}, errors.New("invalid error"), -2
	}
	return RedisError{Value: s[1:crlfpos]}, nil, crlfpos + len(CRLF)
}

func DecodeInt(s string) (RedisInt, error, int) {
	crlfpos := strings.Index(s, CRLF)
	if crlfpos < 0 {
		return RedisInt{}, errors.New("invalid int"), -2
	}
	val, err := strconv.Atoi(s[1:crlfpos])
	if err != nil {
		return RedisInt{}, err, -2
	}
	return RedisInt{Value: val}, nil, crlfpos + len(CRLF)
}

func DecodeBulkStr(s string) (RedisBulkStr, error, int) {
	crlfpos := strings.Index(s, CRLF)
	if crlfpos < 0 {
		return RedisBulkStr{}, errors.New("invalid bulk string"), -2
	}
	bulklen, err := strconv.Atoi(s[1:crlfpos])
	if err != nil {
		return RedisBulkStr{}, err, -2
	}
	if bulklen == -1 {
		// Return null
		return RedisBulkStr{null: true}, nil, crlfpos + len(CRLF)
	}
	nextcrlf := strings.Index(s[crlfpos+len(CRLF):], CRLF)
	bulk := s[crlfpos+len(CRLF) : crlfpos+nextcrlf+len(CRLF)]
	return RedisBulkStr{Value: bulk}, nil, crlfpos + bulklen + len(CRLF)
}

func DecodeArray(s string) (RedisArray, error, int) {
	crlfpos := strings.Index(s, CRLF)
	if crlfpos < 0 {
		return RedisArray{}, errors.New("invalid array"), -2
	}
	arrLen, err := strconv.Atoi(s[1:crlfpos])
	if err != nil {
		return RedisArray{}, err, -2
	}

	if arrLen == -1 {
		return RedisArray{null: true}, nil, 5
	}

	nextobjpos := strings.Index(s, CRLF)
	i := nextobjpos

	if nextobjpos < 0 {
		return RedisArray{}, errors.New("invalid array"), -2
	}

	if arrLen == 0 {
		return RedisArray{Value: []interface{}{}}, nil, 4
	}

	arr := []interface{}{}
	for {
		if i < len(s) && len(arr) < arrLen {
			res, err, iplus := Decode(s[i+len(CRLF) : len(s)])
			if err != nil {
				return RedisArray{}, err, -2
			}
			arr = append(arr, res)
			i += iplus
		} else {
			break
		}
	}
	return RedisArray{Value: arr}, nil, i + len(CRLF)
}

func Decode(s string) (interface{}, error, int) {
	i := 0
	for {
		if i > len(s) {
			break
		}
		switch s[i] {
		case '+':
			res, err, iplus := DecodeStr(s[i : strings.Index(s[i:], CRLF)+len(CRLF)+i])
			if err != nil {
				return nil, err, i
			}
			i += iplus
			return res, nil, i
		case '-':
			res, err, iplus := DecodeError(s[i : strings.Index(s[i:], CRLF)+len(CRLF)+i])
			if err != nil {
				return nil, err, i
			}
			i += iplus
			return res, nil, i
		case '$':
			res, err, iplus := DecodeBulkStr(s[i:len(s)])
			if err != nil {
				return nil, err, i
			}
			i += iplus
			return res, nil, i
		case ':':
			res, err, iplus := DecodeInt(s[i : strings.Index(s[i:], CRLF)+len(CRLF)+i])
			if err != nil {
				return nil, err, i
			}
			i += iplus
			return res, nil, i
		case '*':
			res, err, iplus := DecodeArray(s[i:len(s)])
			if err != nil {
				return nil, err, i
			}
			i += iplus
			return res, nil, i
		default:
			i += 1
		}

	}
	return nil, errors.New("this part shouldn't have been reached. Good Luck Troubleshooting :("), 0
}
