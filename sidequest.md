# Sidequest (possible future improvements)
- [ ] use iota for token values instead of string
- [ ] attach filename, line and column number to token for better error handling
- [ ] refactor lexer test so them don't stop at first failure
- [ ] add support for unicode
- [ ] merge readIdentifier and readNumber into a single function readWord(identityFn fn(ch byte) bool)
- [ ] handle float
- [ ] handle hex numbers
- [ ] add support for <= and >=
