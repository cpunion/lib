/*
 * Copyright (c) 2026 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sync

import (
	"runtime"

	"github.com/goplus/lib/c"
)

// OnceFunc is a native C callback accepted by Once.Do. It cannot carry a Go
// closure environment; use Once.DoFunc or the Go sync.Once type when the
// callback captures values.
//
//llgo:type C
type OnceFunc func()

//llgo:type C
type onceContextFunc func(c.Pointer)

type onceContext struct {
	fn func()
}

func invokeOnceContext(data c.Pointer) {
	(*onceContext)(data).fn()
}

// DoFunc invokes f once while keeping its Go closure environment above the C
// ABI boundary. Code that does not need a closure should use Do directly.
func (o *Once) DoFunc(f func()) c.Int {
	context := &onceContext{fn: f}
	result := doOnceContext(o, invokeOnceContext, c.Pointer(context))
	runtime.KeepAlive(context)
	return result
}
