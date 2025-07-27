# GoMath

GoMath is a library for parsing and interpreting math expression.

## Usage

Install the library with
```bash
$ go get -u github.com/nyttikord/gomath@latest
```
You can replace `latest` with any valid tags.

You can install the CLI by cloning the repo and building it.
```bash
$ git clone https://github.com/nyttikord/gomath.git
$ cd gomath/cmd 
$ go build -o gomath . 
```
Then, you can put the file gomath where do you want (it could be in `/usr/local/bin` for example).

### Calculate

To parse an expression, use `gomath.Parse(string) (gomath.Result, error)`.
The string is a valid expression, like `1+2` or `2(1/3+4)^5`.

The result will give you everything needed to perform operations, like getting the exact representation of the result, 
or the ability to convert

```go
res, err := gomath.Parse("1/2+1")
// check the error
res.String == "3/2" // true
res.IsExact // true because 1.5 is the exact representation of 3/2
res.LaTeX == `\frac{1}{2} + 1` // true
res.Approx(5) == "1.50000" // true
```

You can also call `gomath.ParseAndCalculate(string, *gomath.Options) (string, error)` to directly get the string 
representation with the given options or `gomath.ParseAndConvertToLatex(string, *gomath.Options) (string, error)` to get
the $\LaTeX$ code.

### Creating a function

You can create a function with `gomath.NewFunction(string) (gomath.Function, int, error)`.
The string is a valid expression representing a function.
It is composed by the arguments and the expression of the function.
You must separate each argument with a coma (`,`) and separate the arguments and the expression with `->`.
```
x, y -> x^y
```
is a valid expression representing a function.
The returned int is the number of arguments.

Then, you can evaluate the function by passing a map containing every argument.
gomath will send you the result (an instance of `gomath.Result`) of the evaluation.

```go
f, _, err := gomath.NewFunction("x, y -> x^y")
if err != nil {
	panic(err) // will not panic because the expression is valid
}
res, err := f(map[string]string{
	"x": "5",
	"y": "2"
})
if err != nil {
	panic(err) // will not panic because the expression is valid and every argument is set
}
res.String == "25" // true
```

### CLI

You can get the help with `gomath help`.

To evaluate an expression, use `gomath eval <expression>`.

To convert an expression to $\LaTeX$, use `gomath latex <expression>`. 

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
