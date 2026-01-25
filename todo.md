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
        - [x] Boolean operators (!, ||, &&)
        - [x] `In` keyword for turning something into particular unit (`in ms in hr`)

    Refs:
    - https://github.com/desmosinc/pratt-parser-blog-code/tree/main/src
    - https://github.com/munificent/bantam/blob/master/src/com/stuffwithstuff/bantam/expressions/PrefixExpression.java
    - https://github.com/microsoft/typescript-go/blob/1d138eaa29bc189e6b4f04b87fe278b6afe7e62f/internal/parser/parser.go#L4448
    - https://engineering.desmos.com/articles/pratt-parser/
    - https://github.com/Khongchai/compiler-in-go/blob/dfde56ef6bf93bbe17fff62f9da55633f024866d/src/monkey/parser/parser.go
    - https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/
- [x] Implement math function calls. Must alo handle nested function calls.
- [x] Finish pratt parsing
- [x] Executor
- [ ] Working on evaluator loop
    - [x] Can now evaluate simple add expression
    - [x] rest of tests
    - [ ] percent unit
    - [x] Remaining operators 
        - [x] continue implementing comparision and equality operators in eval
        - [x] write tests for boolean comparsion expression
    - [x] Add builtins support (continue from here (e *Evaluator) evalCallExpression(functionName ast.Expression)
        - [x] implement
        - [x] tests
    - [x] add % postfix 
    ~~- [-] Deal with n^2 expression problem in `evalBinaryNumberExpression`~~
        - Ok there is actually a plus side to this -- using switch statement like that is going to allow for very clear error branches.
    - [x] Add error propagation and handling
        - [x] parser/scanner
        - [x] clear all panic from evaluator
    - [ ] Wrap with lsp
        - [x] Finish engine.go Do it like microsoft to learn about queue handling in go 
        - [x] Start up server and see it being called by vscode!
        - [ ] Test that 
        ```
        // | 1 + 2 
        const something = 2;
        ```
        Emits text decoration of 3
- [ ]  Handle decoration -- this case not yet passing
- func (e *Engine) handleTextDocumentDidChange(ctx context.Context, params *lsproto.DidChangeTextDocumentParams) error {
	// uri := params.TextDocument.Uri

	return nil
}
- [ ] Handle diagnostics -- not implemented on the client side yet.

<!-- end section -->

# Backlog
- [ ] mod keyword (builtin invocation)
- [ ] pi and e constant (allow overriding)
- [ ] Summing everything above current line (can be done in a later stage that detects special keyword)
- [ ] Add more unit supports (metrics, gbs)
- [ ] Apply go worker and channels to optimize wait time when network fetching is required
- [ ] handle non utf8 char with DecodeRuneInString
- [ ] Multiline and parsing initialization in executor layer.
- [ ] Support binary value view
    - `| 2 in binary`
- [ ] Support hex value view
    - `| 2 in hex`
- [ ] Should also allow this `// | x = 2, y = 3` (multiple declaration in one line)
- [ ] Apply promise to make the program's rate conversion faster.