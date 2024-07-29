---
_class: lead
paginate: true
backgroundColor: #fff
marp: true
math: mathjax
---

# Zero-Knowledge Succinct Proofs (Part II)

---

## Recap

### Statement

"I know a number x such that
for a FibSq sequence with $a_0 = 1, a_1 = x$,
$a_{1022} = 2338775057$"

---

- Simple proof
  - not zero-knowledge
  - not succinct
- Proof with EC and pairing
  - encrypted
  - not succinct

---

## Proof with EC

### Proof

$s_1 = \{a_0G_1, a_1G_1,...,a_{1022}G_1\}$
$s_2 = \{a_0G_2, a_1G_2,...,a_{1022}G_2\}$

### Verification

1. $s_{1,0} = G_1$
2. $s_{1,1022} = 2338775057G_1$
3. $e(a_iG_1, a_iG_2) + e(a_{i+1}G_1, a_{i+1}G_2) = e(a_{i+2}G_1, G_2)$
    >$a_i^2 + a_{i+1}^2 = a_{i+2} \cdot 1$
4. $e(a_iG_1, G_2) = e(G_1, a_iG_2)$
    >$a_i \cdot 1 = 1 \cdot a_i$

---

## Problems with Our Proof

- not succinct
- **very specific**, **can't handle degrees higher than quadratic**

---

## Generalizing Our Proof System

Consider this _FibFrankenstein_ formula:

$a_{i+2} = a_{i+1}^2 + a_{i}^{\color{red} 3}$

If we were to keep using our previous constraint form to accomodate this, it would need something that looks like:

$e(a_{i+2}G_1, G_2) = e(a_{i+1}G_1, a_{i+1}G_2) + e(a_iG_1, a_iG_2, {\color{red} a_i???})$

---

## Witness to the Rescue (Again)

Previously, we introduced the witness vectors $s_1$ and $s_2$ to "help" the verifier check our proof.

### The Idea

We can "flatten" the computation of any term of $deg > 2$ by computing it's quadratic components one at a time.

---

## Flattening

$a_{i+2} = a_{i+1}^2 + {\color{red} a_i^3}$
$~~\downarrow$
$v_i = a_i^2$
$a_{i+2} = a_{i+1}^2 + {\color{red} v_i \cdot a_i}$

---

## New Proof

$s_1 = \{a_iG_1\}_{i=0}^{1022}$
$s_2 = \{a_iG_2\}_{i=0}^{1022}$

$\color{red} v = \{v_iG_1\}_{i=0}^{1020}$

## New Verification

$\color{red} e(v_iG_1, G_2) = e(a_iG_1, a_iG_2)$
> $v_i = a_i^2$

$e(a_{i+2}G_1, G_2) = e(a_{i+1}G_1, a_{i+1}G_2) + {\color{red} e(v_iG_1, a_iG_2)}$
> $a_{i+2} = a_{i+1}^2 + v_i \cdot a_i$

---

## The Tools We have

- **Elliptic Curves** Algebraic structure preserving value encryption
- **EC Pairing** Multiplication on encrypted values
- **Witness** Intermediate steps of the computation to help verification
- **Flattening** To help eliminate above quadratic terms

---

## Succinctness

Succinctness = small proof size (verification computation)

- Why is succinctness useful and desired?
- Application in blockchain

---

## Succinctness (Polynomials)

### Schwartz-Zeppel Lemma

Let $d$ be the degree of polynomials $P$ and $Q$, and $S$ be the domain
If we choose $x$ at **random** uniformly and independently from $S$, and if $P(x) = Q(x)$ then
$Pr[P = Q] \geq 1 - \frac{d}{|S|}$

- Probabilistic polynomial identity testing
- If we use an absurdly large $S$ and relatively small $d$, then the probablity that the two polys aren't the same after one successful testing is neglegible.
- Proof of Knowledge vs. Argument of Knowledge

---

## Reducing the Problem to Polynomial Identity Testing

Proving "I know $x$ that passes these checks"
$~~\downarrow$
Proving "I know witness $v$ that satisfies these constraints"
$~~\downarrow$
??? (Probably a ton of steps)
$~~\downarrow$
Proving the equality of two polynomials

> **Warning** Things are going to change drastically, but our understanding of all previous concepts stay valid.

---

Speaking of polynomials...

- $a_{i+2} = a_{i+1}^2 + a_{i}^{3}$
- A way to encode the computation steps (constraints) into polynomials.
- What plays well with polynomials? Matrices!

---

## A Simpler Statement

"I know the solution $x$ to $y = x^3 + 3x - 2 \mod 79$, when $y = 34$"

> $x = 3$

### Flatten

$v = x \cdot x$

$y = v \cdot x + 3x - 2$
$\rightarrow y + 2 = (v + 3) x$

### Witness

$s = [1, y, x, v] = [1, 38, 3, 9]$

---

## Rank 1 Constraint System (R1CS)

### Constraints

$x \cdot x = v$
$(v + 3) x = y + 2$

### Witness

$s = [1, y, x, v] = [1, 38, 3, 9]$

$~~~\downarrow$

${\color{blue}L} \cdot s \odot {\color{blue}R} \cdot s = {\color{blue}O} \cdot s$

$$
\begin{bmatrix}
0 & 0 & 1 & 0 \\
3 & 0 & 0 & 1 \\
\end{bmatrix}
\cdot
\begin{bmatrix}
1 \\
y \\
x \\
v
\end{bmatrix}
\odot
\begin{bmatrix}
0 & 0 & 1 & 0 \\
0 & 0 & 1 & 0 \\
\end{bmatrix}
\cdot
\begin{bmatrix}
1 \\
y \\
x \\
v
\end{bmatrix}
=
\begin{bmatrix}
0 & 0 & 0 & 1 \\
2 & 1 & 0 & 0 \\
\end{bmatrix}
\cdot
\begin{bmatrix}
1 \\
y \\
x \\
v
\end{bmatrix}
$$

---

## Interpolation on Columns

$$
{\color{blue}L} = \begin{bmatrix}
0 & 0 & 1 & 0 \\
3 & 0 & 0 & 1 \\
\end{bmatrix}~~
{\color{blue}R} = \begin{bmatrix}
0 & 0 & 1 & 0 \\
0 & 0 & 1 & 0 \\
\end{bmatrix}~~
{\color{blue}O} = \begin{bmatrix}
0 & 0 & 0 & 1 \\
2 & 1 & 0 & 0 \\
\end{bmatrix}
$$

$\texttt{interpolateCols({\color{blue}L})} \rightarrow {\color{red}U} = (3x + 76, 0, 78x + 2, x + 78)$
$\texttt{interpolateCols({\color{blue}R})} \rightarrow {\color{red}V} = (0, 0, 1, 0)$
$\texttt{interpolateCols({\color{blue}O})} \rightarrow {\color{red}W} = (2x + 77, x + 78, 0, 78x + 2)$

$s = [1,y,x,v] = [1,38,3,9]$

We also need to interpolate $s$ to satisfy the ring homomorphism, but it'd look the same:
$\texttt{interpolateCols({\color{blue}s})} \rightarrow {\color{red}a} = (1,38,3,9)$

$Ua \cdot Va = Wa$

---

## Homomorphism

- Vectors under addition and element-wise product forms a ring.
- Polynomials under addition and multiplication forms a ring.

There exists a homomorphism from _vectors under element-wise product_ to _polynomials under multiplication_

${\color{blue}Ls} \rightarrow {\color{red}Ua}$
${\color{blue}Rs} \rightarrow {\color{red}Va}$
${\color{blue}Os} \rightarrow {\color{red}Wa}$
${\color{blue}Ls} \odot {\color{blue}Rs} \rightarrow {\color{red}Ua} \cdot {\color{red}Va}$

---

Wait, how is $Ls$ an element-wise product and how are $L, R, O$ vectors?

${\color{blue}L} = \begin{bmatrix}
0 & 0 & 1 & 0 \\
3 & 0 & 0 & 1 \\
\end{bmatrix}~~~{\color{blue}L}= \begin{bmatrix}0\\3\end{bmatrix},\begin{bmatrix}0\\0\end{bmatrix},\begin{bmatrix}1\\0\end{bmatrix},\begin{bmatrix}0\\1\end{bmatrix}$

${\color{blue}s} = [1,38,3,9]~~~{\color{blue}s} = \begin{bmatrix}1\\1\end{bmatrix},\begin{bmatrix}38\\38\end{bmatrix},\begin{bmatrix}3\\3\end{bmatrix},\begin{bmatrix}9\\9\end{bmatrix}$

And if we take a closer look at ${\color{blue}Ls} \rightarrow {\color{red}Ua}$

${\color{blue}Ls} = \begin{bmatrix}0\\3\end{bmatrix} \odot \begin{bmatrix}1\\1\end{bmatrix},~\begin{bmatrix}0\\0\end{bmatrix} \odot \begin{bmatrix}38\\38\end{bmatrix}, ...$

${\color{red}Ua} = (3x+76)\cdot 1,~ 0 \cdot 38,...$

---

## Reduction

$\text{If } Ua \cdot Va = Wa \text{, then } Ls \odot Rs = Os$

### But There is a Problem

$deg(U_i) + deg(V_i)$ can be greater than $deg(W_i)$

---

## Balancing the Equation

$Ua \cdot Va = Wa + {\color{red}???}$

The idea:

Rewrite the equation for matrices $Ls \odot Rs = Os + {\color{red}\oslash}$

> If we "interpolate" $\oslash$, we get 0. BUT, it can also be a polynomial that vanishes everywhere in the domain.

Let $t(x)$ be the vanishing polynomial $t(x) = (x-1)(x-2)...(x-n)$ of degree $deg(L)$ in our domain,

$Ua \cdot Va = Wa + {\color{red} h(x)t(x)}$

$$
{\color{red}h(x)} = \frac{Ua\cdot Va - Wa}{\color{red}t(x)}
$$

---

### What We Have Now

$Ua \cdot Va = Wa + h(x)t(x)$

- We interpolated the R1CS into a set of polynomial sets $U, V, W$ representing the constraints
- The prover has the data they want to prove (the witness $s$, or $a$ after interpolation)
- We computed $h(x)$ to balance the equation

### Still Missing

- Verifier still can't smaple and test the equality $Ua \cdot Va = Wa + h(x)t(x)$

---

## What's Next?

- Interactive Oracle Proofs (IOP)
- Encrypted polynomial evaluation (using EC)
- Trusted setup (the backbone of everything)
