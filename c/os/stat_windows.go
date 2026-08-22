//go:build windows

/*
 * Copyright (c) 2026 The GoPlus Authors (goplus.org). All rights reserved.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package os

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

//go:linkname Stat C._stat64
func Stat(path *c.Char, buf *StatT) c.Int

// The Universal CRT has no lstat variant. Its stat operation exposes the
// portable Windows file metadata available through StatT.
//
//go:linkname Lstat C._stat64
func Lstat(path *c.Char, buf *StatT) c.Int

//go:linkname Fstat C._fstat64
func Fstat(fd c.Int, buf *StatT) c.Int
