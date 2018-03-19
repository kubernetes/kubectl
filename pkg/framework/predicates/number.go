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

// Number is a "number predicate". It's a type that decides if a
// number matches or not.
type Number interface {
	Match(float64) bool
}

// NumberNot inverts the result of the sub-predicate.
func NumberNot(predicate Number) Number {
	return numberNot{ip: predicate}
}

type numberNot struct {
	ip Number
}

func (p numberNot) Match(i float64) bool {
	return !p.ip.Match(i)
}

// NumberAnd returns true if all the sub-predicates are true. If there are
// no sub-predicates, always returns true.
func NumberAnd(predicates ...Number) Number {
	return numberAnd{ips: predicates}
}

type numberAnd struct {
	ips []Number
}

func (p numberAnd) Match(i float64) bool {
	for _, ip := range p.ips {
		if !ip.Match(i) {
			return false
		}
	}
	return true
}

// NumberOr returns true if any sub-predicate is true. If there are no
// sub-predicates, always returns false.
func NumberOr(predicates ...Number) Number {
	ips := []Number{}

	// Implements "De Morgan's law"
	for _, ip := range predicates {
		ips = append(ips, NumberNot(ip))
	}
	return NumberNot(NumberAnd(ips...))
}

// NumberEqual returns true if the value is exactly i.
func NumberEqual(i float64) Number {
	return numberEqual{i: i}
}

type numberEqual struct {
	i float64
}

func (p numberEqual) Match(i float64) bool {
	return i == p.i
}

// NumberGreaterThan returns true if the value is strictly greater than i.
func NumberGreaterThan(i float64) Number {
	return numberGreaterThan{i: i}
}

type numberGreaterThan struct {
	i float64
}

func (p numberGreaterThan) Match(i float64) bool {
	return i > p.i
}

// NumberEqualOrGreaterThan returns true if the value is equal or greater
// than i.
func NumberEqualOrGreaterThan(i float64) Number {
	return NumberOr(NumberEqual(i), NumberGreaterThan(i))
}

// NumberLessThan returns true if the value is strictly less than i.
func NumberLessThan(i float64) Number {
	// It's not equal, and it's not greater than i.
	return NumberAnd(NumberNot(NumberEqual(i)), NumberNot(NumberGreaterThan(i)))
}

// NumberEqualOrLessThan returns true if the value is equal or less than i.
func NumberEqualOrLessThan(i float64) Number {
	return NumberOr(NumberEqual(i), NumberLessThan(i))
}
