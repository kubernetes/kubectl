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

package predicates_test

import (
	"fmt"
	"testing"

	. "k8s.io/kubectl/pkg/framework/predicates"
)

// This example shows you how you can create a IntP, and how it's use to
// compare with the actual value.
//
// XXX: This could definitely be improved to add a better example.
func ExampleNumber() {
	fmt.Println(NumberEqual(5).Match(5))
	// Output: true
}

func TestNumberEqual(t *testing.T) {
	if !NumberEqual(5).Match(5) {
		t.Fatal("NumberEqual(5) should match 5")
	}
	if NumberEqual(5).Match(4) {
		t.Fatal("NumberEqual(5) should not match 4")
	}
}

func TestNumberNot(t *testing.T) {
	if NumberNot(NumberEqual(5)).Match(5) {
		t.Fatal("NumberNot(NumberEqual(5)) should not match 5")
	}
	if !NumberNot(NumberEqual(5)).Match(4) {
		t.Fatal("NumberNot(NumberEqual(5)) should match 4")
	}
}

func TestNumberAnd(t *testing.T) {
	if !NumberAnd().Match(5) {
		t.Fatal("NumberAnd() should match 5")
	}
	if !NumberAnd(NumberEqual(5)).Match(5) {
		t.Fatal("NumberAnd(NumberEqual(5)) should match 5")
	}
	if NumberAnd(NumberEqual(5)).Match(4) {
		t.Fatal("NumberAnd(NumberEqual(5)) should not match 4")
	}
	if NumberAnd(NumberEqual(5), NumberEqual(4)).Match(5) {
		t.Fatal("NumberAnd(NumberEqual(5), NumberEqual(4)) should not match 5")
	}
}

func TestNumberOr(t *testing.T) {
	if NumberOr().Match(5) {
		t.Fatal("NumberOr() should not match 5")
	}
	if !NumberOr(NumberEqual(5)).Match(5) {
		t.Fatal("NumberOr(NumberEqual(5)) should match 5")
	}
	if NumberOr(NumberEqual(5)).Match(4) {
		t.Fatal("NumberOr(NumberEqual(5)) should not match 4")
	}
	if !NumberOr(NumberEqual(5), NumberEqual(4)).Match(5) {
		t.Fatal("NumberOr(NumberEqual(5), NumberEqual(4)) should match 5")
	}
}

func TestNumberLessThan(t *testing.T) {
	if NumberLessThan(3).Match(5) {
		t.Fatal("NumberLessThan(3) should not match 5")
	}
	if NumberLessThan(3).Match(3) {
		t.Fatal("NumberLessThan(3) should not match 3")
	}
	if !NumberLessThan(3).Match(1) {
		t.Fatal("NumberLessThan(3) should match 1")
	}
}

func TestNumberEqualOrLessThan(t *testing.T) {
	if NumberEqualOrLessThan(3).Match(5) {
		t.Fatal("NumberEqualOrLessThan(3) should not match 5")
	}
	if !NumberEqualOrLessThan(3).Match(3) {
		t.Fatal("NumberEqualOrLessThan(3) should match 3")
	}
	if !NumberEqualOrLessThan(3).Match(1) {
		t.Fatal("NumberEqualOrLessThan(3) should match 1")
	}
}

func TestNumberGreaterThan(t *testing.T) {
	if !NumberGreaterThan(3).Match(5) {
		t.Fatal("NumberGreaterThan(3) should match 5")
	}
	if NumberGreaterThan(3).Match(3) {
		t.Fatal("NumberGreaterThan(3) should not match 3")
	}
	if NumberGreaterThan(3).Match(1) {
		t.Fatal("NumberGreaterThan(3) should not match 1")
	}
}

func TestNumberEqualOrGreaterThan(t *testing.T) {
	if !NumberGreaterThan(3).Match(5) {
		t.Fatal("NumberGreaterThan(3) should match 5")
	}
	if NumberGreaterThan(3).Match(3) {
		t.Fatal("NumberGreaterThan(3) should not match 3")
	}
	if NumberGreaterThan(3).Match(1) {
		t.Fatal("NumberGreaterThan(3) should not match 1")
	}
}
