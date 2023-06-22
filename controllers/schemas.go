package controllers

import (
	"os"
	"strings"
)

func UpdateRefAndIdOnSchema(apiVersion string, path string, mp map[string]interface{}) {
	modifyValueFromKeyWithPrefix(mp, "$ref", getBaseSchemaUrl(apiVersion, path), true)
	modifyValueFromKeyWithPrefix(mp, "$id", getBaseSchemaUrl(apiVersion, path), false)
}

func modifyValueFromKeyWithPrefix(mp map[string]interface{}, key_name string, prefix string, recursive bool) {
	for k, v := range mp {
		if recursive {
			switch a := v.(type) {
			case map[string]interface{}:
				modifyValueFromKeyWithPrefix(a, key_name, prefix, true)
			case []interface{}:
				for _, i := range a {
					switch ii := i.(type) {
					case map[string]interface{}:
						modifyValueFromKeyWithPrefix(ii, key_name, prefix, true)
					}
				}
			}
		}
		if k == key_name {
			value := v.(string)
			if strings.HasPrefix(value, "schemas/") {
				value = prefix + value
			} else {
				value = prefix + "schemas/" + value
			}
			mp[k] = value
		}
	}
}

func getBaseSchemaUrl(apiVersion string, path string) (baseUrl string) {
	if os.Getenv("SERVICE_BASE_URL") != "" {
		baseUrl = os.Getenv("SERVICE_BASE_URL")
	} else {
		baseUrl = "http://localhost"
	}
	if len(path) == 0 {
		baseUrl = baseUrl + "/" + apiVersion + "/"
	} else {
		baseUrl = baseUrl + "/" + path + "/" + apiVersion + "/"
	}
	return
}
