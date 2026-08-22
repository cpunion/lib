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

package setjmp

import (
	_ "unsafe"

	"github.com/goplus/lib/c"
)

const LLGoPackage = "decl"

// MSVC exposes jmp_buf as an opaque, pointer-aligned array. Reserve more than
// the current AMD64 and ARM64 layouts so the public Go type remains usable
// across both Windows architectures without copying private CRT definitions.
type JmpBuf [32]uintptr
type SigjmpBuf = JmpBuf

//go:linkname Setjmp llgo.setjmp
func Setjmp(env *JmpBuf) c.Int

//go:linkname Longjmp llgo.longjmp
func Longjmp(env *JmpBuf, val c.Int)

// Windows has no signal mask to preserve. LLGo lowers these intrinsics to the
// architecture-specific Windows context save/restore pair used by its runtime.
//
//go:linkname Sigsetjmp llgo.sigsetjmp
func Sigsetjmp(env *SigjmpBuf, savemask c.Int) c.Int

//go:linkname Siglongjmp llgo.siglongjmp
func Siglongjmp(env *SigjmpBuf, val c.Int)
