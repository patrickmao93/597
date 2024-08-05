# Chapter 5: GROTH16

## Path to Succinctness

Succinctness was briefly mentioned in the previous chapters, let's reiterate on it.

Succinctness = small proof size (verification computation)

It's really that simple, but what does "small" mean? As with a lot of the concept in this domain, there is no absolute and strict definition for it. Most of the time, small means that the complexity of the proof is low with respect to the data being proven. O(log n) can be considered small, O(polylog n) can be considered small. But what we are after is the holy grail of small, which is O(1). It means no matter the amount of data being proven, the proof is always the same size. Also note that the terms proof size and verification complexity are interchangeable because it mostly assumed that verification complexity grows linears with proof size.

There are many vastly different proof systems designed over the years but most of them rely on the same property of polynomials to achieve succinctness.

### Schwartz-Zeppel Lemma

Let $d$ be the degree of polynomials $P$ and $Q$, and $S$ be the domain
If we choose $x$ at **random** uniformly and independently from $S$, and if $P(x) = Q(x)$ then
$Pr[P = Q] \geq 1 - \frac{d}{|S|}$

Note: the original Schwartz-Zeppel lemma is multivariate. But in our use case it only needs to be monovariate, so I took the freedom simplify it here.

- Probabilistic polynomial identity testing
- If we use an absurdly large $S$ and relatively small $d$, then the probability that the two polys aren't the same after one successful testing is negligible.

This lemma itself is dependent on a basic fact about polynomials: two polynomials with degree less than or equal to $d$ can have **at most** n intersections. Moreover, since our evaluation domain $S$ is discrete (because we use a finite field), it means that we have a limited number of "$x$" values to evaluate the polynomials at. When we take an evaluation of the polynomials at a **random** $x$, then the term $\frac{d}{|S|}$ reflects the probability that the $x$ we take happen to be the intersection of the two polynomials.

Note that if we use Schwartz-Zeppel Lemma in a proof system, then it can no longer be called a "proof" because a proof needs to be certain, or in other words, has no soundness error.

## Groth16

Groth16, like many other proof systems, leverages the Schwartz-Zeppel lemma to do polynomial identity tests, and in turn, achieve succinctness. The big question is: how can we reduce the verifying steps (checking whether a bunch of constraint equations hold) into polynomial identity testing?

It might seem tempting to ask: can't we just somehow use our constraints directly? We had constraints like $a_{i+2} = a_{i+1}^2 + a_{i}^2$. the LHS and RHS are both polynomials, can't we just do something with them? The answer is Yes and No. While we definitely have to incorporate these constraints into our final polynomials, they aren't the only things needed. These constraints alone does not reflect the data being proven, and we want our polynomials to capture and reflect BOTH the constraints and the data being proven.

Here is a rough outline of Groth16's reduction steps for reference. This probably doesn't make any sense now, but we will visit each part in the subsequent writing.

Arithmetic Circuit (Arithmetization)

$\downarrow$

Rank 1 Constraint System (R1CS)

$\downarrow$

QAP (polynomials)

$\downarrow$

Encrypted polynomial evaluation (with EC)

$\downarrow$

Polynomial identity test

## Arithmetic Circuit

Arithmetic circuits are the "programs" we write to express a set of constraints that are essentially "how the data should be checked". We have actually already introduced this part, we just didn't specifically give it a name. When we wrote down the constraints such as in the Fibonacci examaple, $a_{i+2} = a_{i+1}^2 + a_{i}^2$, we are defining an arithmetic circuit.

## A Simpler Statement

In order to demonstrate how rank-1 constraint system (R1CS) works, we will have to come up with a simpler example because complex examples will be a nightmare to demonstrate when working with matrices.

A simpler statement that we want to prove:

"I know the solution $x$ to $y = x^3 + 3x - 2 \mod 79$, when $y = 34$"

> The "secret" answer is $x = 3$

### Flatten

We will perform the same kind of flattening over the constraints to reduce the degree of the $x^3$ term.

$v = x \cdot x$

$y = v \cdot x + 3x - 2$

$\rightarrow y + 2 = (v + 3) x$ (rearrange to prepare for R1CS)

### Witness

For this proof, we still need some witness vectors like before. But notice that now we aren't encrypting the witness values into EC points yet.

$s = [1, y, x, v] = [1, 38, 3, 9]$

Also, what is this strange element of $1$ in front of everything? This is called the constant term and it's function will become immediately apparent when we use them in our R1CS.

## Rank-1 Constraint System (R1CS)

R1CS is a constraint system where each row represents an equality constraint:

$L \cdot s \odot R \cdot s = O \cdot s$

(Note: If this equality is true, then the witness vector $s$ passes the verification check)

where $L, R, O$ are the "operation" matrices on the witness vector s.

### R1CS with Our Simple Example

Now, we know that we have a set of constraints

$x \cdot x = v$

$(v + 3) x = y + 2$

and a witness vector

$s = [1, y, x, v] = [1, 38, 3, 9]$

Let's see how we transform these into R1CS and how the new constraint system captures both the constraints and the data.

${\color{blue}L} \cdot s \odot {\color{blue}R} \cdot s = {\color{blue}O} \cdot s$

$\downarrow$

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

The "$\cdot$" operator is just regular matrix dot product and the "$\odot$" operator is matrix element-wise product.

In the first row of the $L$ matrix, we can see that $l_{1,4}$ is 1. If we attempt a matrix dot product with the witness $s$, the 1 is essentially "flipping on" the $x$ element in the witness. The same goes for the $R$ matrix. The $x$ is flipped on again. Then we move to the RHS. Because $o_{1,4}$ is 1, that means the $v$ element in $s$ is flipped on. Rememebr that the "$\odot$" symbol is element-wise product, that means the entire first line of this constraint system reads "$x \cdot x = v$", which is exactly our first constraint!

Let's look at the second row now. $l_{2,1} = 3$, $l_{2,4} = 1$. After dot product with the witness, the second row becomes $[3 + 0 + 0 + v] = [3 + v]$. For the $R$ matrix, $r_{2,3} = 1$, that means we flip on $x$. For the $O$ matrix, $o_{2,1}=2$ and $o_{2,2} = 1$, so the RHS matrix, after the dot product, becomes $[2 + y + 0 + 0] = [2 + y]$. And now, the entire second row after the dot product reads $(3 + v) \cdot x = 2 + y$, which is again, exactly our second constraint!

I hope now it's clear why we needed to put a 1 in front of the witness vector. In case it isn't, this is because we have constant terms such as the 3 in $(v+3)x$ and the 2 in $y + 2$.

Something that's not so obvious but worth mentioning now is that the reason why we had to "flatten" our circuit (reducing degree to at most 2) is because we have exactly one "$\odot$" in R1CS! Therefore, we can't represent $x \cdot x \cdot x$. And before we reveal it, the reader is encouraged to think about why we have to limit ourselves inside this seemingly "restrictive" system (the reason was mentioned already in the previous chapters).

## From Matrices to Polynomials

We mentioned that the key to succinctness is polynomials, but we have been working with matrices. It is because matrices turn out to work very well with polynomials.

### Polynomial Interpolation

The way we turn matrices into polynomials is through interpolating the columns of the matrices.

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

### What We Have Now

$Ua \cdot Va = Wa + h(x)t(x)$

- We interpolated the R1CS into a set of polynomial sets $U, V, W$ representing the constraints
- The prover has the data they want to prove (the witness $s$, or $a$ after interpolation)
- We computed $h(x)$ to balance the equation

What's still missing is that the verifier can't sample and test the equality $Ua \cdot Va = Wa + h(x)t(x)$

In the next chapter, we will fully dive into the Groth16 backend to see how we can leverage the polynomials we get from this chapter to craft efficiently verified proofs.
