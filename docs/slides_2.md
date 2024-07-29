---
_class: lead
paginate: true
backgroundColor: #fff
marp: true
math: mathjax
---

# Zero-Knowledge Succinct Proofs (Part III)

Arithmetic Circuits

---

## Arithmetic Circuits

A way to express the verification logic in arithmetics, and is something that is transformable to R1CS.

---

## Arithmetic Circuits & GROTH16

Arithmetic Circuit (Arithmetization)
$\downarrow$
Rank 1 Constraint System (R1CS)
$\downarrow$
QAP (polynomials)
$\downarrow$
Trusted setup $\rightarrow$ proving key, verifying key
$\downarrow$
Encrypted polynomial evaluation (with EC)
$\downarrow$
Polynomial identity test

---

## Arithmetization

... is a step where we "translate" logic into arithmetics.
... is just the verification logic.

### Example

"I have the preimage of this SHA-3 hash"

We work with equations like $y = 3x^3 + 2$, how do we turn SHA-3 into these?

SHA-3 = a bunch of AND, XOR, LROT, RROT

---

## AND

| x   | y   | z   |
| --- | --- | --- |
| 0   | 0   | 0   |
| 0   | 1   | 0   |
| 1   | 0   | 0   |
| 1   | 1   | 1   |

$x, y$: input signals
$z$: output signal

In circuit: $z = xy$

---

## XOR

| x   | y   | z   |
| --- | --- | --- |
| 0   | 0   | 0   |
| 0   | 1   | 1   |
| 1   | 0   | 1   |
| 1   | 1   | 0   |

In circuit: $z = x + y - 2xy$

But, there is a catch.

---

## XOR (Cont.)

The previous arithmetization is almost complete, but it still has a "bug".

What if the input signals x and y are 2 and 3?

$z = 2 + 3 - 2 \cdot 2 \cdot 3 = -7$

Our input signals must be either 0 or 1 for the XOR gate!

---

## Range Check

Range checking itself is an arithmetization step. Checking binary digits is easy:

$(x - 1)x=0$

Generalizing to checking that the range of x is in [m, n] where m < n

$(x - n)(x - n - 1)(x - n - 2)...(x - m)=0$

The LHS is just the vanishing polynomial in range [m, n]

---

And with that, our AND and XOR gates are now fully constrained:

XOR
$z = x + y - 2xy$
$\color{red}(x-1)x = 0$
$\color{red}(y-1)y = 0$
$\color{red}(z-1)z = 0$

AND
$z = xy$
$\color{red}(x-1)x = 0$
$\color{red}(y-1)y = 0$
$\color{red}(z-1)z = 0$

---

## Code

Composing the AND and XOR Gates in Circuit.

For demo purposes, our circuit will be extremely simple:

$u_{Pub} = x \land y \bigoplus z_{Pub}$

The prover must prove they have the correct private values $x, y$ and public value $z, u$ that satisfies this circuit.

We are going to use Gnark. It is a Go library for developing arithmetic circuits.

---

## Some Big Caveats

- There is no dynamic array
- There is no loop
- There is no if-else

---

## Signal Selector (if-else)

$z = select(c, x, y)$

Arithmetization
$z = c \cdot x + isZero(c) \cdot y$
$(c - 1)c = 0$

where
$isZero: y = 1 - x$

---

## IsZero (Generalized)

$isZero(x)$. Outputs 1 when $x = 0$, and 0 otherwise.

Arithmetization
$$
m =
\begin{cases}
\frac{1}{x}, & \text{if } x \neq 0 \\
0 & \text{if } x = 0 \\
\end{cases}
$$

$$
out = -m \cdot x + 1
$$

---

## IsZero (Cont.)

Recall that we have shown how the prover can "help" verification through supplying more data than the data they want to prove...

> $s_1 = \{a_0G_1, a_1G_1,...,a_{1022}G_1\}$
> $s_2 = \{a_0G_2, a_1G_2,...,a_{1022}G_2\}$ where $a_1$ is just $x$

$m$ is exactly like that, something computed by the prover to help the verification.

$m$ is computed outside of the circuit by the prover
$$
m =
\begin{cases}
\frac{1}{x}, & \text{if } x \neq 0 \\
0 & \text{if } x = 0 \\
\end{cases}
$$

$$
out = -m \cdot x + 1
$$

$$
\color{red}x \cdot out = 0
$$
