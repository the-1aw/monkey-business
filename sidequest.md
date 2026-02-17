# Sidequest (possible future improvements)
- [ ] attach filename, line and column number to token for better error handling
- [ ] refactor lexer test so them don't stop at first failure
- [ ] add support for unicode (currently lexer uses char, we would need to use rune)
- [ ] merge readIdentifier and readNumber into a single function readWord(identityFn fn(ch byte) bool)
- [ ] handle float
- [ ] handle hex numbers
- [ ] add support for <= and >=
- [ ] make the semicolon at end of line optionnal
- [ ] add postfix operators (eg, ++, --)
