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

package predicates

import (
	"regexp"
	"strings"
)

// String is a "string predicate". It's a type that decides if a
// string matches or not.
type String interface {
	Match(string) bool
}

// StringNot will inverse the result of the predicate.
func StringNot(predicate String) String {
	return stringNot{sp: predicate}
}

type stringNot struct {
	sp String
}

func (p stringNot) Match(str string) bool {
	return !p.sp.Match(str)
}

// StringAnd returns true if all the sub-predicates are true. If there are
// no sub-predicates, always returns true.
func StringAnd(predicates ...String) String {
	return stringAnd{sps: predicates}
}

type stringAnd struct {
	sps []String
}

func (p stringAnd) Match(str string) bool {
	for _, sp := range p.sps {
		if !sp.Match(str) {
			return false
		}
	}
	return true
}

// StringOr returns true if any sub-predicate is true. If there are no
// sub-predicates, always returns false.
func StringOr(predicates ...String) String {
	sps := []String{}

	// Implements "De Morgan's law"
	for _, sp := range predicates {
		sps = append(sps, StringNot(sp))
	}
	return StringNot(StringAnd(sps...))
}

// StringEqual returns a predicate that matches only the exact string.
func StringEqual(str string) String {
	return stringEqual{str: str}
}

type stringEqual struct {
	str string
}

func (p stringEqual) Match(str string) bool {
	return p.str == str
}

// StringLength matches if the length of the string matches the given
// integer predicate.
func StringLength(predicate Number) String {
	return stringLength{ip: predicate}
}

type stringLength struct {
	ip Number
}

func (p stringLength) Match(str string) bool {
	return p.ip.Match(float64(len(str)))
}

// StringHasPrefix matches if the string starts with the given prefix.
func StringHasPrefix(prefix string) String {
	return stringHasPrefix{prefix: prefix}
}

type stringHasPrefix struct {
	prefix string
}

func (p stringHasPrefix) Match(str string) bool {
	return strings.HasPrefix(str, p.prefix)
}

// StringHasSuffix matches if the string ends with the given suffix.
func StringHasSuffix(suffix string) String {
	return stringHasSuffix{suffix: suffix}
}

type stringHasSuffix struct {
	suffix string
}

func (p stringHasSuffix) Match(str string) bool {
	return strings.HasSuffix(str, p.suffix)
}

// StringRegexp matches if the string matches with the given regexp.
func StringRegexp(regex *regexp.Regexp) String {
	return stringRegexp{regex: regex}
}

type stringRegexp struct {
	regex *regexp.Regexp
}

func (p stringRegexp) Match(str string) bool {
	return p.regex.MatchString(str)
}
