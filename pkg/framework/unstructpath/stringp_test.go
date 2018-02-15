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
	"regexp"
	"testing"

	. "k8s.io/kubectl/pkg/framework/unstructpath"
)

func TestStringEqual(t *testing.T) {
	if !StringEqual("some string").Match("some string") {
		t.Fatal(`StringEqual("some string") should match "some string"`)
	}
	if StringEqual("some string").Match("another string") {
		t.Fatal(`StringEqual("some string") should not match "another string"`)
	}
}

func TestStringNot(t *testing.T) {
	if StringNot(StringEqual("some string")).Match("some string") {
		t.Fatal(`StringNot(StringEqual("some string")) should not match "some string"`)
	}
	if !StringNot(StringEqual("some string")).Match("another string") {
		t.Fatal(`StringNot(StringEqual("some string")) should match "another string"`)
	}
}

func TestStringAnd(t *testing.T) {
	if !StringAnd().Match("some string") {
		t.Fatal(`StringAnd() should match "some string"`)
	}
	if !StringAnd(StringEqual("some string")).Match("some string") {
		t.Fatal(`StringAnd(StringEqual("some string")) should match "some string"`)
	}
	if StringAnd(StringEqual("some string")).Match("another string") {
		t.Fatal(`StringAnd(StringEqual("some string")) should not match "another string"`)
	}
	if StringAnd(StringEqual("some string"), StringEqual("another string")).Match("some string") {
		t.Fatal(`StringAnd(StringEqual("some string"), StringEqual("another string")) should not match "some string"`)
	}
}

func TestStringOr(t *testing.T) {
	if StringOr().Match("some string") {
		t.Fatal(`StringOr() should not match "some string"`)
	}
	if !StringOr(StringEqual("some string")).Match("some string") {
		t.Fatal(`StringOr(StringEqual("some string")) should match "some string"`)
	}
	if StringOr(StringEqual("some string")).Match("another string") {
		t.Fatal(`StringOr(StringEqual("some string")) should not match "another string"`)
	}
	if !StringOr(StringEqual("some string"), StringEqual("another string")).Match("some string") {
		t.Fatal(`StringOr(StringEqual("some string"), StringEqual("another string")) should match "some string"`)
	}
}

func TestStringLength(t *testing.T) {
	if !StringLength(NumberEqual(11)).Match("some string") {
		t.Fatal(`StringLength(NumberEqual(11)) should match "some string"`)
	}

	if StringLength(NumberLessThan(6)).Match("some string") {
		t.Fatal(`StringLength(NumberLessThan(6)) should not match "some string"`)
	}

	if !StringLength(NumberLessThan(16)).Match("some string") {
		t.Fatal(`StringLength(NumberLessThan(16)) should match "some string"`)
	}

}

func TestStringHasPrefix(t *testing.T) {
	if !StringHasPrefix("some ").Match("some string") {
		t.Fatal(`StringHasPrefix("some ") should match "some string"`)
	}
	if StringHasPrefix("some ").Match("another string") {
		t.Fatal(`StringHasPrefix("some ") should not match "some string"`)
	}
}

func TestStringHasSuffix(t *testing.T) {
	if !StringHasSuffix("string").Match("some string") {
		t.Fatal(`StringHasSuffix("string") should match "some string"`)
	}
	if StringHasSuffix("integer").Match("some string") {
		t.Fatal(`StringHasSuffix("integer") should not match "some string"`)
	}
}

func TestStringRegexp(t *testing.T) {
	if !StringRegexp(regexp.MustCompile(".*")).Match("") {
		t.Fatal(`StringRegexp(regexp.MustCompile(".*")) should match ""`)
	}
	if !StringRegexp(regexp.MustCompile(".*")).Match("Anything") {
		t.Fatal(`StringRegexp(regexp.MustCompile(".*")) should match "Anything"`)
	}
	if !StringRegexp(regexp.MustCompile("word")).Match("word") {
		t.Fatal(`StringRegexp(regexp.MustCompile("word")) should match "word"`)
	}
}
