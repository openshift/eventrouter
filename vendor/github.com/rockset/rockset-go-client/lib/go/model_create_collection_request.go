/*
 * REST API
 *
 * Rockset's REST API allows for creating and managing all resources in Rockset. Each supported endpoint is documented below.  All requests must be authorized with a Rockset API key, which can be created in the [Rockset console](https://console.rockset.com). The API key must be provided as `ApiKey <api_key>` in the `Authorization` request header. For example: ``` Authorization: ApiKey aB35kDjg93J5nsf4GjwMeErAVd832F7ad4vhsW1S02kfZiab42sTsfW5Sxt25asT ```  All endpoints are only accessible via https.  Build something awesome!
 *
 * API version: v1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package rockset
import (
    "bytes"
    "encoding/json"
    "fmt"
    
)

type CreateCollectionRequest struct {
	// unique identifer for collection, can contain alphanumeric or dash characters
	Name string `json:"name"`
	// text describing the collection
	Description string `json:"description,omitempty"`
	// list of sources from which to ingest data
	Sources []Source `json:"sources,omitempty"`
	// number of seconds after which data is purged, based on event time
	RetentionSecs int64 `json:"retention_secs,omitempty"`
	// configuration for event data
	EventTimeInfo *EventTimeInfo `json:"event_time_info,omitempty"`
	// list of mappings
	FieldMappings []FieldMappingV2 `json:"field_mappings,omitempty"`
}
func (m CreateCollectionRequest) PrintResponse() {
    r, err := json.Marshal(m)
    var out bytes.Buffer
    err = json.Indent(&out, []byte(string(r)), "", "    ")
    if err != nil {
        fmt.Println("error parsing string")
        return
    }

    fmt.Println(out.String())
}

