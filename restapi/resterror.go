// Punch CLI used for interacting with holepunch.io
// Copyright (C) 2018-2019  Orb.House, LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package restapi

import (
	"errors"
	"fmt"
)

var errorUnownedSubdomain = errors.New("you do not own this subdomain")
var errorCantConnectRestCall = errors.New("problem contacting the server")
var errorUnableToParse = errors.New("can't parse the json response")
var errorUnownedTunnel = errors.New("you do not own this subdomain")
var errorUnableToDelete = errors.New("failed to delete")

//ResponseError JSONapi response error
type ResponseError struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			Title  string `json:"title"`
			Status string `json:"status"`
			Detail string `json:"detail"`
		} `json:"attributes"`
		ID string `json:"id"`
	} `json:"data"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf(e.Data.Attributes.Detail)
}
