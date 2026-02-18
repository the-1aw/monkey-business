# Monkey Business

Monkey Business is a toolchain for the [Monkey programming language](https://monkeylang.org/) reated by [Thorsten Ball](https://thorstenball.com/).  

## Why Monkey Business ?

These sole purpose of this toolchain is to improve my knowledge of go and how programming languages works.

## Usage

## Sidequests

This part list topics I don't want to dive into straight away, as my main goal is to have a working interpreter and compiler.

- [ ] attach filename, line and column number to token for better error handling
- [ ] refactor lexer test so them don't stop at first failure
- [ ] add support for unicode (currently lexer uses char, we would need to use rune)
- [ ] merge readIdentifier and readNumber into a single function readWord(identityFn fn(ch byte) bool)
- [ ] handle float
- [ ] handle hex numbers
- [ ] add support for <= and >=
- [ ] make the semicolon at end of line optionnal
- [ ] add postfix operators (eg, ++, --)
- [ ] update parser testing functions in order to make it easier to add test cases for each expression types
