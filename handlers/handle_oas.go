package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/deliveryhero/devhub-backend-core/controllers"
	"github.com/ghodss/yaml"
	"github.com/pedidosya/peya-go/logs"
	"github.com/pedidosya/peya-go/server"
)

func HandleGetOAS(path string) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			apiVersion := server.GetStringFromPath(r, "apiVersion", "")

			logs.Debugf("[handlers] handling get OAS for api version %s", apiVersion)
	
			oas_content, err := os.ReadFile(fmt.Sprintf("api/%s/oas.yaml", apiVersion))
	
			if err != nil {
				switch err {
				case os.ErrNotExist:
					server.NotFound(w, r, "NOT_FOUND", "can't return OAS spec")
				default:
					server.BadRequest(w, r, "UNKNOWK_ERROR", "can't return OAS spec")
				}
				return
			}
	
			oas_json_content, err2 := yaml.YAMLToJSON(oas_content)
			if err2 != nil {
				logs.Debug("[handlers] error processing file content (yaml to json): %s", err2)
				server.BadRequest(w, r, "ERROR_PROCESING_FILE", "can't process the requested OAS content")
				return
			}
	
			var oas_json map[string]interface{}
			err = json.Unmarshal(oas_json_content, &oas_json)
				
			if err != nil {
				logs.Debug("[handlers] error processing file content (json unmarshal): %s", err)
				server.BadRequest(w, r, "ERROR_PROCESING_FILE", "can't process the requested OAS content")
				return
			}
	
			controllers.UpdateRefAndIdOnSchema(apiVersion, path, oas_json)
	
			patched_json_bytes, err2 := json.Marshal(oas_json)
	
			if err2 != nil {
				logs.Debug("[handlers] error processing file content (json marshal): %s", err2)
				server.BadRequest(w, r, "ERROR_PROCESING_FILE", "can't process the requested OAS content")
				return
			}
	
			oas_yaml_bytes, err2 := yaml.JSONToYAML(patched_json_bytes)
	
			if err2 != nil {
				logs.Debug("[handlers] error processing file content (json to yaml): %s", err2)
				server.BadRequest(w, r, "ERROR_PROCESING_FILE", "can't process the requested OAS content")
				return
			}
			
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/yaml")
			w.Write(oas_yaml_bytes)
		}
}
