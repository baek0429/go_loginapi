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
 * This package receives the httprequest from the client(mobile)
 * and returns response with appropriate information.
 * NOTICE:
 **/

package login

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
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
	UniqueId    string // uuid
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

	// setting up db
	db, err := sql.Open("mysql", "root@cloudsql(loginapi-1125:instance-name)/dbname")
	if err != nil {
		fmt.Fprint(w, err)
	}
	defer db.Close()

	// empty request and response struct address
	var request Request
	var response Response

	// quried values from URL
	queried := r.URL.Query()

	// assign Request parameters from quaried to local values
	request.Tag = queried.Get("Tag")
	request.Username = queried.Get("Username")
	request.Password = queried.Get("Password")

	// querying db based on the request info
	response.Tag = "Response"

	// if tag : login, register, ?
	if request.Tag == "Login" {
		if (request.Username != "") || (request.Password != "") {
			// sql query for selecting user
			sqlStmt := "SELECT unique_id, username, encrypted_password, created_at, updated_at FROM users" +
				" WHERE username ='" + request.Username + "'"

			// rows with the same request.Username
			rows, err := db.Query(sqlStmt)
			if err != nil {
				fmt.Fprint(w, err)
				return
			}
			defer rows.Close()

			// []byte for storing the scanned values
			var col1, col2, col3, col4, col5 []byte
			for rows.Next() {
				// scan the value to []byte
				err := rows.Scan(&col1, &col2, &col3, &col4, &col5)
				if err != nil {
					fmt.Fprint(w, err)
					return
				}

				// if password is right successcode= 1 else 0
				if bcrypt.CompareHashAndPassword(col3, []byte(request.Password)) == nil {
					response.UniqueId = string(col1)
					if response.UniqueId != "" {
						response.SuccessCode = 1
					}
				} else {
					response.SuccessCode = 0
				}
			}
		}
	} else if request.Tag == "Register" {
		// statment frame
		sqlStmt := "INSERT INTO users (unique_id, username, encrypted_password) VALUES( ? , ? , ?, ?)"

		// prepare statment for inserting data
		insStmt, err := db.Prepare(sqlStmt)
		if err != nil {
			fmt.Fprint(w, err)
			return

		}
		defer insStmt.Close()
		var sqlValues []string

		// TODO: Hashing password, assiging unique_id
		sqlValues = GetStoreUserQuery(request.Username, request.Password)

		// executes query
		_, err = insStmt.Exec(sqlValues[0], sqlValues[1], sqlValues[2], sqlValues[3])
		if err != nil {
			response.SuccessCode = 0
		} else {
			response.SuccessCode = 1
		}

	} else {
		response.SuccessCode = 0
	}

	// save info as json object
	fmt.Fprint(w, response.ToString())
	jsonObj, err := json.Marshal(response)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	// print json obj
	fmt.Fprintln(w, string(jsonObj[:]))
}

// takes string var name and password
// creates uuid, hashed password
// returns string[] with uuid, username, hashedpassword
func GetStoreUserQuery(name, password string) []string {
	var result []string

	u4, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error: ", err)
		return result
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	result = []string{u4.String(), name, string(hashed)}
	return result
}

// for test only
func ConfirmPassword(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err == nil {
		return true
	} else {
		return false
	}
}
