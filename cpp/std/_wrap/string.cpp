#if defined(_WIN32) && defined(__clang__) && __clang_major__ < 20
// The current Visual Studio STL rejects older Clang releases by default.
// LLGo intentionally supports LLVM 19, whose MSVC ABI is sufficient for this
// wrapper; use the STL's documented escape hatch instead of changing the
// compiler or C++ ABI selected by the caller.
#define _ALLOW_COMPILER_AND_STL_VERSION_MISMATCH
#endif

#include <string>
#include <new>

extern "C" {

// -----------------------------------------------------------------------------

void stdStringInitEmpty(std::string* s) {
    new(s) std::string();
}

void stdStringInitFrom(std::string* s, std::string* v) {
    new(s) std::string(*v);
}

void stdStringInitFromCStr(std::string* s, const char* cstr) {
    new(s) std::string(cstr);
}

void stdStringInitFromCStrLen(std::string* s, const char* cstr, size_t len) {
    new(s) std::string(cstr, len);
}

void stdStringDispose(std::string* s) {
    s->~basic_string();
}

// -----------------------------------------------------------------------------

const char* stdStringCStr(const std::string* s) {
    return s->c_str();
}

const char* stdStringData(const std::string* s) {
    return s->data();
}

size_t stdStringSize(const std::string* s) {
    return s->size();
}

// -----------------------------------------------------------------------------

} // extern "C"
