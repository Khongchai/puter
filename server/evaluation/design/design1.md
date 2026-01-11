# Overall

The language consists for 4 main components

```
scan -> parse -> eval -> (EvaluationResult.View() -> string value for the lsp)
```

## EvaluationResult

The EvaluationResult from evaluation step is just an object -- a `job` that tells the executor whether the result is ready to be viewed or needs further processing.

For example, consider this:

```ts
/**
 * | rate = 0.06 percent
 * | gross = 2 in usd
 * | net = gross * rate in thb
 */
function calculateIncome() {
    // ...
}
```

On the last line, gross needs to be converted to thb from usd first. This requires a network fetch. The executor must perform this.

The evaluation output and side effects of each line is as follows:

**1**: `rate = 0.06 percent`

1. Stores `rate` in heap as `{"name": "rate", "unit": "percent", "value": Future(2)}`
2. Outputs `{value: "0.06", unit: percent}` // we can then print this is 0.06% later down the line.

**1**: `gross = 2 in usd`

1. Stores `gross` in heap as `{"name": "gross", "unit": "usd", "value": Future(2)}`
2. Outputs `{value: 2, unit: "usd"}`

**2**: `net = gross * rate in thb`

1. Stores `net` in heap as `{"name": "net", "unit": "thb", "value": Future(go func() { fetch data and return })}`
2. Outputs `{"result": {"type": "variable", name: "gross"}}`




# Components

# Evaluator 

The evaluator is meant to be initialized once at the first pipe character in a contiguous comment block:

```
// | << here
//
//
```
