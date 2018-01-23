/*
Copyright 2018 The Kubernetes Authors.

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

package unstructpath_test

import (
	"github.com/ghodss/yaml"

	"k8s.io/kubectl/pkg/framework/unstruct"
	. "k8s.io/kubectl/pkg/framework/unstruct/unstructpath"
)

// This example is inspired from http://goessner.net/articles/JsonPath/#e3.
func Example() {
	data := `{ "store": {
	    "book": [
	      { "category": "reference",
	        "author": "Nigel Rees",
	        "title": "Sayings of the Century",
	        "price": 8.95
	      },
	      { "category": "fiction",
	        "author": "Evelyn Waugh",
	        "title": "Sword of Honour",
	        "price": 12.99
	      },
	      { "category": "fiction",
	        "author": "Herman Melville",
	        "title": "Moby Dick",
	        "isbn": "0-553-21311-3",
	        "price": 8.99
	      },
	      { "category": "fiction",
	        "author": "J. R. R. Tolkien",
	        "title": "The Lord of the Rings",
	        "isbn": "0-395-19395-8",
	        "price": 22.99
	      }
	    ],
	    "bicycle": {
	      "color": "red",
	      "price": 19.95
	    }
	  }
	}`
	y := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(data), &y); err != nil {
		panic(err)
	}
	u := unstruct.New(y)

	// The authors of all books in the store. Returns a list of strings.
	Children().Map().Field("book").Children().Map().Field("author").String().Select(u)

	// All authors. Returns a list of strings.
	All().Map().Field("author").String().Select(u)

	// All things in store, which are some books and a red bicycle. Returns a list of Values.
	Map().Field("store").Children().Select(u)

	// The price of everything in the store. Returns a list of "float64".
	Map().Field("store").All().Map().Field("price").Number().Select(u)

	// The third book. Returns a list of 1 Value.
	All().Map().Field("book").Slice().At(2).Select(u)

	// The last book in order. Return a list of 1 Value.
	All().Map().Field("book").Slice().Last().Select(u)

	// The first two books. Returns a list of 2 Values.
	All().Map().Field("book").Slice().AtP(NumberLessThan(2)).Select(u)

	// Filter all books with isbn number. Returns a list of Values.
	All().Map().Field("book").Filter(Map().Field("isbn")).Select(u)

	// Filter all books cheaper than 10. Returns a list of "float64".
	All().Map().Field("book").Children().Filter(Map().Field("price").Number().Filter(NumberLessThan(10))).Select(u)

	// All elements in structure. Returns a list of Values.
	All().Select(u)
}
