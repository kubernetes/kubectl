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

package selectors

import (
	"github.com/ghodss/yaml"

	p "k8s.io/kubectl/pkg/framework/path/predicates"
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
	u := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(data), &u); err != nil {
		panic(err)
	}

	// The authors of all books in the store. Returns a list of strings.
	Children().Field("book").Children().Field("author").AsString().SelectFrom(u)

	// All authors. Returns a list of strings.
	All().Field("author").AsString().SelectFrom(u)

	// All things in store, which are some books and a red bicycle. Returns a list of maps.
	Field("store").Children().AsMap().SelectFrom(u)

	// The price of everything in the store. Returns a list of "float64".
	Field("store").All().Field("price").AsNumber().SelectFrom(u)

	// The third book. Returns a list of 1 interface{}.
	All().Field("book").At(2).SelectFrom(u)

	// The last book in order. Return a list of 1 interface{}.
	All().Field("book").Last().SelectFrom(u)

	// The first two books. Returns a list of 2 interface{}.
	All().Field("book").AtP(p.NumberLessThan(2)).SelectFrom(u)

	// Filter all books with isbn number. Returns a list of interface{}.
	All().Field("book").Filter(Field("isbn")).SelectFrom(u)

	// Filter all books cheaper than 10. Returns a list of interface{}.
	All().Field("book").Children().Filter(Field("price").AsNumber().Filter(p.NumberLessThan(10))).SelectFrom(u)

	// All elements in structure. Returns a list of interface{}.
	All().SelectFrom(u)
}
