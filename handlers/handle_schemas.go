package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"embed"

	"github.com/cmartinc22/devhub-backend-core/controllers"
	"github.com/pedidosya/peya-go/logs"
	"github.com/pedidosya/peya-go/server"
)

var content embed.FS

func HandleGetSchemas(path string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		apiVersion := server.GetStringFromPath(r, "apiVersion", "")
		schema := server.GetStringFromPath(r, "schema", "")

		logs.Debugf("[handlers] handling get schema %s for api version %s", schema, apiVersion)

		if !strings.HasSuffix(schema, ".json") {
			schema = schema + ".json"
		}

		if match, _ := regexp.MatchString(`[a-z0-9]([a-z0-9-]*[a-z0-9])?\.json`, schema); !match {
			server.BadRequest(w, r, "INVALID_PARAMETER", fmt.Sprintf("can't return the requested schema %s", schema))
			return
		}

		schema_content_test, err := content.ReadFile("handle_oas.go")
		logs.Debug(schema_content_test)
		logs.Debug(err)
		schema_content, err := content.ReadFile(fmt.Sprintf("../api/%s/schemas/%s", apiVersion, schema))
		logs.Debug(schema_content)
		logs.Debug(err)
		if err != nil {
			schema_content, err = os.ReadFile(fmt.Sprintf("api/%s/schemas/%s", apiVersion, schema))
		}
		
		if err != nil {
			switch err {
			case os.ErrNotExist:
				server.NotFound(w, r, "NOT_FOUND", "can't return the requested schema")
			default:
				server.BadRequest(w, r, "UNKNOWK_ERROR", "can't return the requested schema")
			}
			return
		}


		// Now let's unmarshall the data into `schema_json
		var schema_json map[string]interface{}
		err = json.Unmarshal(schema_content, &schema_json)
		controllers.UpdateRefAndIdOnSchema(apiVersion, path, schema_json)
		
		if err != nil {
			logs.Debug("[handlers] error processing file content: %s", err)
			server.BadRequest(w, r, "ERROR_PROCESING_FILE", "can't process the requested schema content")
			return
		}

		server.OK(w, r, schema_json)
	}
}