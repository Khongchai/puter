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
- [ ] Implement math function calls. Must alo handle nested function calls.
- [ ] More decide later

<!-- end section -->

# 09/12/2025
- [x] Finish pratt parsing

# 13/12/2025
- [ ] Evaluator
    - [ ] Takes a text, modify the state of the program.
    - [ ] Evaluator's eval produces a string, ready-to-be-display value. To achieve this with a lot of API call, we need to wrap all computation result with a monad-like struct:
      - During compilation phase, evaluator may encounter an operation that requires an async fetch 
        "x = 2 usd in thb"  // adds this to object pool ({value: 2, unit: "thb", type: "x", line: 0})
        "a = x + 5 in thb"  // A is a computation result (a Promise) 
        Here, line 0 blocks line 2 from displaying any value. These two lines can be parsed and eval by two different go routines. So it's best if they somehow "sync" their values whenever they are ready to display.
        Solution:
            - whenever an api value is needed, the evaluator must submit a request to a `fetcher` instance and wait for returned result.
            Otherwise that lines' value is a Monad that resolves immediately to the line's value.


    - [ ] Error propagation.
- [ ] Executor
    Tree-walking interpreter, goes 
    ```
    executor := NewExecutor() // global scope
    func main() {
        // executor loop
        changes := breakIntoLines(change)
        for _, change range (changes) {
            if (!change.text.startsWith(allAllowedCommentsThenPipe)) {
                continue
            }
            result, errors := executor.exec(change.text) // breaks down change into lines and pass that into the evaluator.
            emitErrorIfAny(errors, change.line) // errors include position within line and such
            if result != nil {
                emitResult(result) // print to lsp, whatever
            }
        }
    }
    ```

<!-- end section -->

# Backlog
- [ ] Apply go worker and channels to optimize wait time when network fetching is required
- [ ] handle non utf8 char with DecodeRuneInString
- [ ] Multiline and parsing initialization in executor layer.
- [ ] Support binary value view
    - `| 2 in binary`
- [ ] Support hex value view
    - `| 2 in hex`
- [ ] Should also allow this `// | x = 2, y = 3` (multiple declaration in one line)