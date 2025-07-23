# GoMath

GoMath is a library for parsing and interpreting math expression.

## Usage

Install the library with
```bash
$ go get -u github.com/nyttikord/gomath@latest
```
You can replace `latest` with any valid tags.

You can install a CLI with
```bash
$ go install github.com/nyttikord/gomath/cmd@latest
```
You can replace `latest` with any valid tags.

### Calculate

To parse an expression and calculate it, use `gomath.ParseAndCalculate(string, *gomath.Options) (string, error)`.
The string is a valid expression, like `1+2` or `2(1/3+4)^5`.
It returns the result of the expression in a string according to the given options.

```go
res, err := gomath.ParseAndCalculate("1+2", &gomath.Options{})
err == nil // true
res == "3" // true
```

You can modify the result's type with `gomath.Options`.
Set `Decimal` to `true` if you want to have a decimal approximation.
You can specify the number of digits with `Precision`.

### Convert to LaTeX

To parse an expression and convert it into $\LaTeX$, use `gomath.ParseAndConvertToLatex(string, *gomath.Options) (string, error)`.
The string is a valid expression like `1+2` or `2(1/3+4)^5`.
It returns the result of the expression in a string according to the given options.

```go
res, err := gomath.ParseAndConvertToLatex("(1+2/3)/2", &gomath.Options{})
err == nil // true
res == `\frac{1 + \frac{2}{3}}{2}` // true
```

### CLI

You can get the help with `gomath help`.

To evaluate an expression, use `gomath eval <expression>`.

To convert to $\LaTeX$ an expression, use `gomath latex <expression>`. 

### Special case

The written representation of calculation is definitely not compatible with computers.
We use common signes and common conventions to prevent unwanted behaviors.
For example, the multiplication is represented by `*` and not by `×`.

The computer representation of calculation has some special cases.
How to interpret `5/2(1+2)`? And `-3^2`?
All special cases and how they are interpreted are listed below.

| Computer representation | Interpretation of GoMath |
|-------------------------|--------------------------|
| `-a^b`                  | $-(a^b)$                 |
| `a/b(c+d)`              | $\frac{a}{b}(c+d)$       |
| `a/b*c`                 | $c\frac{a}{b}$           |
| `a/bx`                  | $\frac{a}{bx}$           |
| `a^b^c`                 | $(a^b)^c$                |

These cases are listed on [Wikipedia](https://en.wikipedia.org/wiki/Order_of_operations#Special_cases).

### Supported operation

All common operators (`+`, `-`, `*`, `/`, `^`, `!`) are supported.
Parenthesis (`(`, `)`) are also supported.

We plan to add the support for the modulo (`%`).

### Supported variables

$\pi$ is represented by `pi`.
$e$ is represented by `e`.
$\phi$ is represented by `phi`.

### Supported functions

We plan to add the support for common functions, like $\exp$ or $\sin$ or $\cos$ or $\tan$ or $\ln$ or $\log$.

## Contribution

Before requesting a merge request, be sure that all tests pass.

Do not use any LLMs.

Your contribution will be licensed under MPL 2.0, like the whole project.
