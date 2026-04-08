# Monkey Business

Monkey Business is a toolchain for the [Monkey programming language](https://monkeylang.org/) created by [Thorsten Ball](https://thorstenball.com/).  

## Why Monkey Business ?

The sole purpose of this toolchain is to improve my knowledge of go and how programming languages toolchains' work.  
**This MUST NOT be considered production-ready under any circumstances.**

## Usage

## Sidequests

This part list topics I don't want to dive into straight away, as my main goal is to have a working interpreter and compiler.

- [ ] look into project structure (use internal and split in two engine(jit/tree-walk) use eiter 2 cli or one cli with options)
- [ ] attach filename, line and column number to token for better parser error handling
- [ ] add stack trace to the interpreter error
- [ ] refactor lexer test so them don't stop at first failure
- [ ] add support for unicode (currently lexer uses char, we would need to use rune)
- [ ] merge readIdentifier and readNumber into a single function readWord(identityFn fn(ch byte) bool)
- [ ] handle float
- [ ] handle hex numbers
- [ ] add support for <= and >=
- [ ] add postfix operators (eg, ++, --)
- [x] update parser testing functions in order to make it easier to add test cases for each expression types
- [ ] implement a language server for monkey
- [ ] implement a debbuger adapter for monkey
- [ ] add support for else if
- [ ] add support for switch statement
- [ ] add looping/iteration support 
- [ ] look into the ability to use a more idiomatic error handling instead of isError function
- [ ] look into the ability to build a module in wasm to run the interpreter in the browser
- [x] twist builtin `push` function behavior to allow pushing multiple values at once
- [ ] look into "separate chaining" and "open addressing" as a mitigation stratagy to avoid `HashKey` collision risk(fnv collision risk).
- [ ] one could argue we should take a look into caching `HashKey` result for `Hashable` objects for performance improvement.
- [ ] look into register based VM to see if it might be interesting to replace our stack based one.
- [ ] consider refactoring ast/lexer/evaluator test with the same shpae as compiler and vm.
