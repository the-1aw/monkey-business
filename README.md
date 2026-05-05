# Monkey Business

Monkey Business is a toolchain for the [Monkey programming language](https://monkeylang.org/) created by [Thorsten Ball](https://thorstenball.com/).

It includes a full front-end (lexer, Pratt parser, AST) shared by two independent execution back-ends: a tree-walking interpreter and a bytecode compiler paired with a stack-based virtual machine. The VM path is roughly 3× faster than the interpreter and serves as a practical study of how real language runtimes work under the hood.

Monkey itself is a dynamically-typed scripting language with first-class functions, closures, integers, booleans, strings, arrays, and hash maps.

## Why Monkey Business?

The sole purpose of this toolchain is to improve my knowledge of Go and how programming language toolchains work.  
**This MUST NOT be considered production-ready under any circumstances.**

## ⬇️ Installation

**Prerequisites:** Go 1.25.5 or later.

```bash
go install github.com/the-1aw/monkey-business@latest
```

## 🚀 Usage

`monkey-business` provides two REPL commands backed by different execution engines:

```bash
# Bytecode compiler + VM (~3x faster)
monkey-business run

# Tree-walking interpreter
monkey-business walk
```

Both commands start an interactive REPL where you can type Monkey expressions:

```
Hello <user>! This is the Monkey programming language!
Feel free to type in commands
>> let add = fn(x, y) { x + y };
>> add(3, 4)
7
```

| Command | Engine | Notes |
|---------|--------|-------|
| `run` | Bytecode compiler + stack-based VM | Preferred for performance |
| `walk` | Recursive AST evaluator | Simpler execution model |

## Sidequests

This section lists topics I don't want to dive into straight away, as my main goal is to have a working interpreter and compiler.

- [ ] Look into project structure (use internal and split in two engine(jit/tree-walk) use eiter 2 cli or one cli with options).
- [ ] Attach filename, line and column number to token for better parser error handling.
- [ ] Add stack trace to the interpreter error.
- [ ] Refactor tests so they don't stop at first failure especially for compiler and vm.
- [ ] Add support for unicode (currently lexer uses char, we would need to use rune).
- [ ] Merge `readIdentifier` and `readNumber` into a single function `readWord(identityFn fn(ch byte) bool)`.
- [ ] Handle floats.
- [ ] Handle hex numbers.
- [ ] Add support for `<=` and `>=`.
- [ ] Add postfix operators (e.g., `++`, `--`).
- [x] Update parser testing functions in order to make it easier to add test cases for each expression type.
- [ ] Implement a language server for Monkey.
- [ ] Implement a debugger adapter for Monkey.
- [ ] Add support for `else if`.
- [ ] Add support for switch statements.
- [ ] Add support for string indexing.
- [ ] Add looping/iteration support.
- [ ] Look into using a more idiomatic error handling approach instead of the `isError` function.
- [ ] Look into the ability to build a wasm module to run the interpreter in the browser.
- [x] Twist the builtin `push` function behavior to allow pushing multiple values at once.
- [ ] Look into "separate chaining" and "open addressing" as a mitigation strategy to avoid `HashKey` collision risk (fnv collision risk).
- [ ] Look into caching `HashKey` results for `Hashable` objects for performance improvement.
- [ ] Look into a register-based VM to see if it might be worth replacing our stack-based one.
- [ ] Consider refactoring ast/lexer/evaluator tests with the same shape as compiler and vm.
- [ ] Use go function option pattern to replace `compiler.NewWithState` and `vm.NewWithGlobalsStore` (overkill for the use case but will gain knowledge)
- [ ] Use go function option pattern to replace `compiler.NewSymbolTable` and `compiler.NewEnclosedSymbolTable` (overkill for the use case but will gain knowledge)
- [ ] Raise an error from the vm on unknown OpCode instead of undefined behavior most likely to panic.
- [ ] Investigate how to make `GetBuiltinByName` O(1) instead of O(n)
