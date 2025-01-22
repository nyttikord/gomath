# GoMath

GoMath is a new lang aiming to be between a formal calculator and a complete computer language.

## Syntax

```gomath
2^3 + 4.5*4.17
```
will give you the result of this addition.

```gomath
let a = 5
```
is defining a new variable.

```gomath
for x in R, f(x) = 5x + 1
```
is defining a new function.

```gomath
for x in ]-\pi / 2, \pi / 2[, f(x) = tan{x}
f{\pi / 2}
```
is defining an alias for $\tan$ on $\left]-\frac{\pi}{2}, \frac{\pi}{2}\right[$. 
You are evaluating it on $\frac{\pi}{2}$ (which is not possible): it will throw an error.

## Todo

- [x] Simple calcul
- [ ] Variable
- [ ] Using predefined variable (like $\pi$)
- [ ] Functions on specific ensemble (e.g., $\mathbb{R}$, $\mathbb{Z}$)
- [ ] Simplification of equations
- [ ] Functions on all space (e.g., $\{1, 2\}$, $[0, 3]$, $[| -4, 8 |]$, $[-8, -4]\cup [|2, 23|]$)
- [ ] Derivation
- [ ] Equations solver

## Technos

Made with Go.

Licensed under MPL 2.0.