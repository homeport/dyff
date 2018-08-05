// Copyright © 2018 Matthias Diester
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package dyff_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/HeavyWombat/dyff/pkg/v1/dyff"
	. "github.com/gorilla/mux"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Input test cases", func() {
	Context("Input data from local sources", func() {
		It("should load multiple JSON documents from one stream", func() {
			doc0, doc1 := `{ "key": "value" }`, `[ { "foo": "bar" } ]`

			documents, err := LoadDocuments([]byte(doc0 + "\n" + doc1))
			Expect(err).To(BeNil())
			Expect(len(documents)).To(BeEquivalentTo(2))
			Expect(documents[0]).To(BeEquivalentTo(yml(doc0)))
			Expect(documents[1]).To(BeEquivalentTo(list(doc1)))
		})
	})

	Context("Input data from remote locations", func() {
		var server *httptest.Server

		BeforeEach(func() {
			r := NewRouter()
			r.HandleFunc("/v1/assets/{directory}/{filename}", func(w http.ResponseWriter, r *http.Request) {
				vars := Vars(r)
				directory := vars["directory"]
				filename := vars["filename"]

				location := "../../../assets/" + directory + "/" + filename
				if _, err := os.Stat(location); os.IsNotExist(err) {
					w.WriteHeader(404)
					fmt.Fprintf(w, "File not found: %s/%s", directory, filename)
					return
				}

				data, err := ioutil.ReadFile(location)
				if err != nil {
					Fail(err.Error())
				}

				w.WriteHeader(200)
				w.Write(data)
			})

			server = httptest.NewServer(r)
		})

		AfterEach(func() {
			if server != nil {
				server.Close()
			}
		})

		It("should load a YAML via a HTTP request", func() {
			inputfile, err := LoadFile(server.URL + "/v1/assets/examples/from.yml")
			Expect(err).To(BeNil())
			Expect(inputfile).ToNot(BeNil())
		})

		It("should fail if the HTTP request fails", func() {
			_, err := LoadFile(server.URL + "/v1/assets/examples/does-not-exist.yml")
			Expect(err.Error()).To(BeEquivalentTo("Unable to load data from " + server.URL + "/v1/assets/examples/does-not-exist.yml: failed to load from location: File not found: examples/does-not-exist.yml"))
		})
	})
})
