# puter

# Supported Comments Type

```md
/**
 *
 *
 */
// 
<!--  -->
# 
```

# Usage

Immediately after comment begin, put a pipe symbol and type in expressions.

```md
/**
 * 
 * | 1 + 1
 */
// | 2 usd to thb
<!-- | 1 in cm to km -->
#  | (log(10) + 5) kb to gb
```

# Syntax

Any math expressions followed by `unit` to `unit`. 

Incompatible units will be underlined in red.

```md
// 2 usd to centimeters  << wrong!
```

## Supported Units

- All currencies' [ISO 4217](https://en.wikipedia.org/wiki/ISO_4217).
- metrics: cm, m, km
- file size: kb, mb, gb, tb,

## Supported Math Functions
- log
- sqrt
- pow
- sum






