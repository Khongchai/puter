# Change Tracking

A change per evaluation context is tracked with:

```go
type EvaluationState struct {
	evaluator *Evaluator;
	// inclusive 
	lineFrom int
	// inclusive
	lineTo int
}
```

# Partial Line Change Detection


Given 

```go
type Change struct {
	lineStart int
	lineEnd int
	content string
}
```

1.

We must first break it down into individual lines:

```go
type LineChangeDetail struct {
	lineIndex int
	text string
}
func diffChange(change Change) {
	// ...
}
var lineChangeDetails []LineChangeDetail = diffChange(change)
```

2.

Then for each line sorted by line index descending, loop through each evaluator in this document:
if line was part of evaluator but is no longer: remove that line and reevaluate
if line was part of evaluator but still is: no-op
if line was not part of evaluator but is: add to evaluator and re-evaluate
if line was not part of evaluator and is not: no-op
if line was not part of evaluator but 

Finding whether a line is valid for evaluation context is:
	line starts with `# |` or `// |` 



