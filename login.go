/**
 * Copyright 2015 Chungseok Baek
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 **/

/**
 * DECRIPTION:
 * This package receives the httprequest from the client(could be mobile)
 * and returns httpresponse with appropriate information.
 * NOTICE:
 **/

package login

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Sampleinfo was made for logging
type SampleInfo interface {
	ToString() string
}

// format for Request
// it is irritating to have a public fields here but public fields must be used to create json objects.
type Request struct {
	Tag      string // request
	Username string // Username
	Password string // Password
}

func (srq *Request) ToString() string {
	return "Tag: " + srq.Tag + " Username: " + srq.Username + " Password: " + srq.Password
}

// format for responseinfo
type Response struct {
	Tag         string // response
	SuccessCode int    // 1 or 0
}

func (srs *Response) ToString() string {
	return "Tag: " + srs.Tag
}

// initiating go app server
func init() {
	http.HandleFunc("/", handler)
}

// main api handler
func handler(w http.ResponseWriter, r *http.Request) {

	// empty request and response struct address
	var request Request
	var response Response

	// quried values from URL
	queried := r.URL.Query()

	// assign Request parameters from quaried to local values
	request.Tag = "Request"
	request.Username = queried.Get("Username")
	request.Password = queried.Get("Password")

	// building responses based on the db search.
	response.Tag = "Response"
	if (request.Username != "") || (request.Password != "") {
		response.SuccessCode = 0
		jsonObj, err := json.Marshal(response)
		if err != nil {
			fmt.Fprint(w, err)
		}
		fmt.Fprintln(w, string(jsonObj[:]))
	}
}
