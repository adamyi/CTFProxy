/*
Copyright 2020 Adam Yi. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

func main() {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		panic(err)
	}
	privPem := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privBytes,
	}
	ioutil.WriteFile("jwt.key", pem.EncodeToMemory(privPem), 0600)

	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		panic(err)
	}
	pubPem := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}
	ioutil.WriteFile("jwt.pub", pem.EncodeToMemory(pubPem), 0644)
}
