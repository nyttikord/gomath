# GoMath

GoMath is a new lang aiming to be between a formal calculator and a complete computer language.

## Syntax

You can evaluate any expressions by just writing it
```gomath
1 + (2 + 3)*4^2/3
```

You can create a variable with `let` statement.
```gomath
let x = 5
```
Every variable is a fraction (here it's `5/1`).
Every calculus is free of double float (float on 64 bits) approximation while these remains in $\mathbb{Q}$.

You can create a function with this statement:
```gomath
for x in R, f{x} = 5*x+1
```
`for x in R` defines the variable in your function and their space (here it's $\mathbb{R}$: all commons sets are 
integrated).
Then, `f{x} =` defines the name of the function (`f`).
Finally, `5*x+1` defines the relation between $x$ and $f(x)$.
:warning: It is `f{x}` and not `f(x)`!


## Todo

- [x] Simple calcul
- [x] Variable
- [x] Using predefined variable (like $\pi$)
- [x] Functions on specific ensemble (e.g., $\mathbb{R}$, $\mathbb{Z}$)
- [ ] Simplification of equations
- [ ] Functions on all space (e.g., $\{1, 2\}$, $[0, 3]$, $[| -4, 8 |]$, $[-8, -4]\cup [|2, 23|]$)
- [ ] Derivation
- [ ] Equations solver

## Technos

Made with Go.

Licensed under MPL 2.0.