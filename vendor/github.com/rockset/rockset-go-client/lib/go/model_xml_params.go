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

type XmlParams struct {
	// tag until which xml is ignored
	RootTag string `json:"root_tag,omitempty"`
	// encoding in which data source is encoded
	Encoding string `json:"encoding,omitempty"`
	// tags with which documents are identified
	DocTag string `json:"doc_tag,omitempty"`
	// tag used for the value when there are attributes in the element having no child
	ValueTag string `json:"value_tag,omitempty"`
	// tag to differentiate between attributes and elements
	AttributePrefix string `json:"attribute_prefix,omitempty"`
}
func (m XmlParams) PrintResponse() {
    r, err := json.Marshal(m)
    var out bytes.Buffer
    err = json.Indent(&out, []byte(string(r)), "", "    ")
    if err != nil {
        fmt.Println("error parsing string")
        return
    }

    fmt.Println(out.String())
}

