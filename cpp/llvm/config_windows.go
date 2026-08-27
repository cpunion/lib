//go:build windows

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

package llvm

const (
	LLGoFiles = "$(pkg-config --cflags llvm-19): _wrap/demangle_windows.cpp"
	// The llvm-19 package is also the native Windows contract used by the Go
	// LLVM bindings. It preserves LLVM's static-library order and required
	// system libraries without exposing MSVC-specific C++ names to Go.
	LLGoPackage = "link: $(pkg-config --libs llvm-19); $(llvm-config --libs) $(llvm-config --ldflags) $(llvm-config --system-libs); -lLLVM"
)
