Create a go program first, hook in lsp later.

- [ ] Parser and executor for basic math expression with no additional function. Parser must support auto precedence (PEMDAS) and explict precedence with ().
    - [x] Scanner
    - [ ] Then apply pratt parsing
        - [x] Finish binary operator here (need binary expression): func (b *BinaryOperatorParselet) Parse(parser *Parser, left *ast.Expression, token ast.Token) ast.Expression {
            if left == nil {
                panic("Left is nil. Can't continue! Infix requires left to be present")
            }

            p := b.precedence
            if b.isRight {
                p -= 1
            }
        }
        - [x] Then assignment operator
        - [x] Then map everything in the main parser method. 
        - [x] Boolean keywords (true, false)
        - [ ] Boolean operators (!, ||, &&)
        - [x] `In` keyword for turning something into particular unit (`in ms in hr`)
        - [ ] Once everything is done, try applying iterative pratt parsing
    - [ ] Executor

    Refs:
    - https://github.com/desmosinc/pratt-parser-blog-code/tree/main/src
    - https://github.com/munificent/bantam/blob/master/src/com/stuffwithstuff/bantam/expressions/PrefixExpression.java
    - https://github.com/microsoft/typescript-go/blob/1d138eaa29bc189e6b4f04b87fe278b6afe7e62f/internal/parser/parser.go#L4448
    - https://engineering.desmos.com/articles/pratt-parser/
    - https://github.com/Khongchai/compiler-in-go/blob/dfde56ef6bf93bbe17fff62f9da55633f024866d/src/monkey/parser/parser.go
    - https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/
- [ ] Implement math function calls. Must alo handle nested function calls.
- [ ] More decide later

<!-- end section -->

# 09/12/2025
- [ ] Finish pratt parsing

# 10/12/2025
- [ ] Iterative pratt parsing

<!-- end section -->

# Backlog

- [ ] Support binary value view
    - `| 2 in binary`
- [ ] Support hex value view
    - `| 2 in hex`