---
_class: lead
paginate: true
backgroundColor: #fff
marp: true
math: mathjax
---

# Zero-Knowledge Succinct Proofs (Part I)

---

## Goal

- What's a knowledge proof and how is it useful?
- What is zero-knowledge and how is it achieved?

---

## Knowledge Proofs

A knowledge proof proves the knowledge of something.

username: my_username
password: p@ssw0rd

> Typical proof protocols consist of a prover and a verifier

---

## FibonacciSq

Just like the familiar Fibonacci sequence but the third element is a sum of the square of the previous two elements

$a_{i+2} = a_{i+1}^2 + a_{i}^2$

- Represented as $a_0, a_1,...,a_n$
- Determined by the first two elements

> Because how fast the sequence grows, we mod the output of each step by a prime number 3221225473

---

## Statement

"I know a number x such that
for a FibSq sequence with $a_0 = 1, a_1 = x$
we have $a_{1022} = 2338775057$"

> 🤫 $x = 3141592$ 🤐

---

## Simple Proof

The prover sends $x = 3141592$ to the verifier
Verifier computes
$$
a_0 = 1, a_1 = 3141592, a_2 = 2986670665, ..., a_{1022} = 2338775057
$$
according to the FibSq formula, and accepts if they arrive at the correct $a_{1022} = 2338775057$

---

## Limitations of the Simple Proof

- **Not zero-knowledge**: The prover gives away their knowledge
- **Not succinct**: The verifier does as much work as the prover

---

## Hiding the Knowledge

We want to encrypt $x$ in a way that its algebraic structure is preseved.

Since the numbers are already on a finite field, we can:

- find a generator $g$ of $F_{p}^{\times}$, which is the multiplicative subgroup of the base field
- encrypt $x$ by exponentiating: $g^x$, this creates a DLP

---

## Hiding the Knowledge (1)

We also need to tweak the FibSq formula to accomedate this change

$g^{a_{i+2}} = g^{a_{i+1}^2} + g^{a_{i}^2}$, where the first and second element in the sequence are $g$ and $g^x$

The prover now sends $g^x$ to the verifier
The verifier computes
$g^{a_2} = g^{1} \cdot g^{x^2}$
$g^{a_3} = g^{x^2} \cdot g^{a_2^2}$
$...$
$g^{a_{1022}} = g^{a_{1010}^2} \cdot g^{a_{1011}^2}$

and accepts if the computed $g^{a_{1022}} = g^{2338775057}$

---

## Hiding the Knowledge (2)

But ... there is a problem: the verifier can't compute $g^{x^2}$ or any $g^{a_i^2}$ because

1. they don't know $x$ or subsequent $a_i$ so they can't compute their square directly
2. this is an exponentiation on the exponent. $g^{x} \cdot g^{x} = g^{2x} \neq g^{x^2}$

Is there a way that allows us to perform multiplication on the encrypted values?

---

## Elliptic Curve

$y^3 = x^2 + ax + b \mod p$

### The set of elliptic curve points forms a group

- Addition P + Q: take a line through P and Q, the sum is the inverse of the third point that the line intersects with the curve
- Identity: 0 (point at infinity)
- Additive Inverse: -(x, y) = (x, -y)

> Scalar multiplication can be achieved by repeated addition $nP = P + P + ... + P \text{(n times)}$.

> Taking $n$ out of $nP$ is weirdly called the "discrete logarithm"

---

## Elliptic Curve Pairing

The EC pairing function $e(P, Q)$ is a bilinear mapping of the curve points $e: H_1 \times H_2 \rightarrow H_T$

### Properties

Let $G_1$, $G_2$ and $G_T$ be the generator points of group 1, 2 and target respectively.

$e(aG_1, bG_2) = e(abG_1, G_2) = e(G_1, abG_2) = abG_T$ (associative)
> $a \cdot b = ab \cdot 1 = 1 \cdot ab$

$e((a + b)G_1, cG_2) = e(aG_1, cG_2) + e(bG_1, cG_2) = (a + b)cG_T$ (distributive)
> $(a + b) \cdot c = ac + bc$

Not all ECs can do pairing. There are some "pairing-friendly" curves such as BN254.

---

## Hiding the Knowledge (3. FibSq with EC)

Rewriting the prove/verify steps with EC:

1. Get a pairing-friendly curve
2. FibSq formula: $a_{i+2}G_T = e(a_{i+1}G_1, a_{i+1}G_2) + e(a_{i}G_1, a_{i}G_2)$

Prover sends $xG_1$ and $xG_2$ to the verifier

Verifier computes
$a_2G_T = 1 \cdot G_T + e(xG_1, xG_2)$
> $a_2 = 1 + x \cdot x$

and ... oops, we can't perform $e$ on $G_T$!

---

## Hiding the Knowledge (4. Witness)

Because the prover knows $x$, they can know the entire sequence of $a_0, a_1, ..., a_n$

The prover encrypts this sequence in both $G_1$ and $G_2$ and sends
$s_1 = \{a_0G_1, a_1G_1,...,a_{1022}G_1\}$
$s_2 = \{a_0G_2, a_1G_2,...,a_{1022}G_2\}$ where $a_1$ is just $x$

and the verifier checks

1. The first element and the 1022nd element in $s_i$ are $1G_1$ and $2338775057G_1$
2. The FibSq fomula is satisfied
   $e(a_iG_1, a_iG_2) + e(a_{i+1}G_1, a_{i+1}G_2) = e(a_{i+2}G_1, G_2)$
    >$a_i^2 + a_{i+1}^2 = a_{i+2} \cdot 1$
3. The $a_i$ encrypted in both $s_1$ and $s_2$ are the same
   $e(a_iG_1, G_2) = e(G_1, a_iG_2)$
    >$a_i \cdot 1 = 1 \cdot a_i$

---

## Recap

### The Statement

"I know a number x such that
for a FibSq sequence with $a_0 = 1, a_1 = x$
we have $a_{1022} = 2338775057$"

### What Both the Prover and the Verifier Know

1. The elliptic curve and the generator $G_1$ and $G_2$
2. The 1022nd number in our FibSq sequence is $2338775057$

---

## Recap (1)

### The Proof

"I can provide vectors $s_1$ and $s_2$ where each element in those vectors represent an intermediate value in the FibSq sequence"

### The Verification

"If the values encrypted in the $s_1$ and $s_2$ you provided are consistent across these two vectors and satisfy the FibSq fomula and if the last element encrypts 2338775057, then indeed you know the answer to the value of $x$ in this setting"
