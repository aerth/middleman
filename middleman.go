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

// Package middleman makes adding middleware easy
package middleman

import (
	"log"
	"net/http"
)

func wrap(heir, f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			f(w, r)
			heir(w, r)
		})
}

type Boolware func(w http.ResponseWriter, r *http.Request) bool

// return heir(w,r) only if f(w,r) returns true
func wrapbool(heir http.HandlerFunc, f Boolware) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if f(w, r) {
				heir(w, r)
			}
		})
}

func Hello(heir http.HandlerFunc) (middled http.HandlerFunc) {
	hello := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Hello", "World")
	}
	middled = wrap(heir, hello)
	return middled
}

func Log(logger *log.Logger, heir http.HandlerFunc) (middled http.HandlerFunc) {
	loghandler := func(w http.ResponseWriter, r *http.Request) {
		logger.Println(r.Method, r.RemoteAddr, r.URL.Path)
	}

	middled = wrap(heir, loghandler)

	return middled

}

func IfThen(boolfunc Boolware, heir http.HandlerFunc) http.HandlerFunc {
	if boolfunc == nil {
		boolfunc = func(w http.ResponseWriter, r *http.Request) bool {
			http.Error(w, "error", http.StatusMethodNotAllowed)
			return false
		}
	}
	middled := wrapbool(heir, boolfunc)
	return middled
}

func SingleHost(allowedhost string, heir http.HandlerFunc) http.HandlerFunc {
	singlehost := func(w http.ResponseWriter, r *http.Request) bool {
		if r.Host == allowedhost {
			return true
		} else {
			w.WriteHeader(403)
			return false
		}
	}
	middled := wrapbool(heir, singlehost)
	return middled
}
