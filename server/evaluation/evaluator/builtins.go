package evaluator

import (
	"math"
	b "puter/evaluation/evaluator/box"
)

type builtinDef struct {
	expectedArgs int
	fn           func(args []b.NumericType) float64
}

var Builtins = map[string]builtinDef{
	"mod":   {2, func(a []b.NumericType) float64 { return math.Mod(a[0].GetNumber(), a[1].GetNumber()) }},
	"log10": {1, func(a []b.NumericType) float64 { return math.Log10(a[0].GetNumber()) }},
	"logE":  {1, func(a []b.NumericType) float64 { return math.Log(a[0].GetNumber()) }},
	"log2":  {1, func(a []b.NumericType) float64 { return math.Log2(a[0].GetNumber()) }},
	"round": {1, func(a []b.NumericType) float64 { return math.Round(a[0].GetNumber()) }},
	"floor": {1, func(a []b.NumericType) float64 { return math.Floor(a[0].GetNumber()) }},
	"ceil":  {1, func(a []b.NumericType) float64 { return math.Ceil(a[0].GetNumber()) }},
	"abs":   {1, func(a []b.NumericType) float64 { return math.Abs(a[0].GetNumber()) }},
	"sin":   {1, func(a []b.NumericType) float64 { return math.Sin(a[0].GetNumber()) }},
	"cos":   {1, func(a []b.NumericType) float64 { return math.Cos(a[0].GetNumber()) }},
	"tan":   {1, func(a []b.NumericType) float64 { return math.Tan(a[0].GetNumber()) }},
	"sqrt":  {1, func(a []b.NumericType) float64 { return math.Sqrt(a[0].GetNumber()) }},
	"lerp": {3, func(a []b.NumericType) float64 {
		v0, v1, t := a[0].GetNumber(), a[1].GetNumber(), a[2].GetNumber()
		return (1-t)*v0 + t*v1
	}},
	"invLerp": {3, func(a []b.NumericType) float64 {
		v0, v1, v := a[0].GetNumber(), a[1].GetNumber(), a[2].GetNumber()
		return (v - v0) / (v1 - v0)
	}},
}
