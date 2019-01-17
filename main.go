// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main // import "github.com/cypherpunkarmory/punch"

import "github.com/cypherpunkarmory/punch/cmd"

func main() {
	cmd.Execute()
}

//this is to setup simple-proxy
//sudo apt-get update
//sudo apt install make(only need to do if compiling code on system)
//sudo apt-get install golang-go
// cd server
//go get
//cd ..
//make
//sudo ./simpleproxy
//sudo apt install docker-compose
//rysnced simple-proxy and holepunch api
//added consul expose port for something not in docker-compose file
