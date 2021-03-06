/*
* The MIT License (MIT)
*
* Copyright (c) 2017  aerth <aerth@riseup.net>
*
* Permission is hereby granted, free of charge, to any person obtaining a copy
* of this software and associated documentation files (the "Software"), to deal
* in the Software without restriction, including without limitation the rights
* to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
* copies of the Software, and to permit persons to whom the Software is
* furnished to do so, subject to the following conditions:
*
* The above copyright notice and this permission notice shall be included in all
* copies or substantial portions of the Software.
*
* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
* IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
* FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
* LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
* OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
* SOFTWARE.
 */

// Package middleman makes adding http middleware handlers easy
package middleman

import (
	"net/http"
)

// Middleware
type Middleware struct {
	f http.Handler // before h
	h http.Handler // after f
}

func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.f.ServeHTTP(w, r)
	m.h.ServeHTTP(w, r)
}

func Wrap(heir, f http.Handler) http.Handler {
	var m Middleware
	m.f = f
	m.h = heir
	return m
}

func WrapFunc(heir, f http.HandlerFunc) http.HandlerFunc {
	var m Middleware
	m.f = http.HandlerFunc(f)
	m.h = http.HandlerFunc(heir)
	return m.ServeHTTP
}

// Boolware returns false if should not continue
type Boolware func(w http.ResponseWriter, r *http.Request) bool

// WrapBoolware returns heir(w,r) only if f(w,r) returns true
func WrapBoolware(heir http.Handler, f Boolware) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if f(w, r) {
				heir.ServeHTTP(w, r)
			}
		})
}

func IfThen(boolfunc Boolware, heir http.Handler) http.Handler {
	if boolfunc == nil {
		boolfunc = func(w http.ResponseWriter, r *http.Request) bool {
			http.Error(w, "error", http.StatusMethodNotAllowed)
			return false
		}
	}
	middled := WrapBoolware(heir, boolfunc)
	return middled
}
