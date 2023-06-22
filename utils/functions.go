package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	PROTOCOL    string = "http"
	API_ADRESS  string = "localhost:3000"
	DATE_LAYOUT string = time.RFC3339Nano
)

const (
	OPERATION_CREATED  string = "created"
	OPERATION_UPDATED  string = "updated"
	OPERATION_DELETED  string = "deleted"
	OPERATION_ARCHIVED string = "archived"
)

func GetSha256Hash(text string) string {
	h := sha256.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func GetSha256HashFromMap(obj map[string]interface{}) string {
	s, err := json.Marshal(obj)
	if err != nil {
		// TODO: it's okk to return empty string for hash?
		return ""
	}
	return GetSha256Hash(string(s))
}

func SliceHasItem[T comparable](slice_obj []T, item T) bool {
	for _, i := range slice_obj {
		if i == item {
			return true
		}
	}
	return false
}

func SliceHasItemGeneric[T any](slice_obj []T, item T) bool {
	switch item := any(item).(type) {
	case map[string]interface{}, []interface{}:
		for _, i := range slice_obj {
			if eq := reflect.DeepEqual(i, item); eq {
				return true
			}
		}
		return false
	}
	return false
}

func FindDuplicatesOnSlice[T any](sliceList []T, keys ...string) []T {
	allKeys := make(map[string]bool, 0)
	list := []T{}
	item_key := ""

	for _, item := range sliceList {
		switch item := any(item).(type) {
		case map[string]interface{}:
			item_key = makeKeyForMapItem(item, keys...)
		default:
			item_key = fmt.Sprintf("%v", item)
		}

		if _, ok := allKeys[item_key]; !ok {
			allKeys[item_key] = true
		} else {
			list = append(list, item)
		}
	}
	return list
}

func RemoveDuplicatesOnSlice[T any ](sliceList []T, keys ...string) []T {
	allKeys := make(map[string]bool, 0)
	list := []T{}
	item_key := ""

	for _, item := range sliceList {
		switch item := any(item).(type) {
		case map[string]interface{}:
			item_key = makeKeyForMapItem(item, keys...)
		default:
			item_key = fmt.Sprintf("%v", item)
		}
		if _, ok := allKeys[item_key]; !ok {
			allKeys[item_key] = true
			list = append(list, item)
		}
	}
	return list
}

func makeKeyForMapItem(item map[string]interface{}, keys ...string) string {
	key_format := strings.Repeat("%s", len(keys))
	ket_args := []interface{}{}
	for _, k := range keys {
		ket_args = append(ket_args, item[k])
	}
	return fmt.Sprintf(key_format, ket_args...)
}

func OverrideMap(original map[string]interface{}, new_data map[string]interface{}) {
	for k, v := range new_data {
		switch v := any(v).(type) {
		case map[string]interface{}:
			OverrideMap(original[k].(map[string]interface{}), v)
		// case []interface{}:
		// TODO override all Array to news, or append if can be added? (problem with links for example, to avoid duplicates)
		default:
			original[k] = v
		}
	}
}

// TODO if needed
func MergeSlice(original []interface{}, new_data []interface{}) {
	for k, v := range new_data {
		switch v := any(v).(type) {
		case map[string]interface{}:
			OverrideMap(original[k].(map[string]interface{}), v)
		case []interface{}:
			// TODO override Array items if there is object
		default:
			original[k] = v
		}
	}
}

func ConcatStringSlices(slices ...[]string) []string {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]string, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}
