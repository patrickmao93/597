# Chapter 3: Zero-Knowledge Proofs

The goal of this chapter is to demonstrate the fundamental concepts used in achieving zero-knowledge in a proof. We will use a toy proof system built by myself that solves a semi-realistic problem.

## Knowledge Proofs

While there isn't a strict definition of knowledge proofs (or more precisely, arguments of knowledge), when we say such term we usually refer to a cryptographic construct that demonstrates the validity of a statement or claim while providing certain security guarantees. When we say argument of knowledge, we usually mean a proof system or protocol. Typical proof protocols consist of a prover and a verifier. To help develop an instinct, we are going to take the username/password combo that we use every day as an analogy. When the user submit their username and password to a server, they are essentially providing a knowledge proof to the server, proving their identity because they can provide a password to a username. In this case, the user provides the proof of identify, therefore, the user is a prover; the server checks the correctness of the user's password given the username, therefore, the server is a verifier. When a prover send a proof to a verifier, the result is a boolean value: the verifier either accepts the proof or rejects the proof.

## What is Zero-Knowledge

The academic way of describing zero-knowledge involves showing that the proof provided by the prover passes in a simulation where a polynomial time verifier cannot distinguish between a real proof and a simulated one. But this is not the focus of this chapter. I want to focus on developing an intuition instead through a semi-familiar example.

## FibonacciSq

FibonacciSq is a slight twist on the familiar Fibonacci sequence. The only difference is that the third element in the recurring relation is the sum of the sqaures of the previous two elements.

$a_{i+2} = a_{i+1}^2 + a_{i}^2$

- Represented as $a_0, a_1,...,a_n$
- Determined by the first two elements

> Because of how fast the sequence grows, we mod the output of each step by a prime number 3221225473

## Statement

Every proof must come with a statement which is the thing we want to prove. For our FibonacciSq example, the statement is:

"I know a number x such that
for a FibSq sequence with $a_0 = 1, a_1 = x$
we have $a_{1022} = 2338775057$"

> The secret answer is $x = 3141592$. Let's pretend only the prover knows about it.

## Simple Proof

The prover sends $x = 3141592$ to the verifier

Verifier computes

$a_0 = 1, a_1 = 3141592, a_2 = 2986670665, ..., a_{1022} = 2338775057$

according to the FibSq formula, and accepts if they arrive at the correct $a_{1022} = 2338775057$

While this proof system works, it has two obvious limitations:

- **Not zero-knowledge**: The prover gives away their knowledge of $x$. Going back to the username/password analogy, this is like giving away the password. While in real world, users do give away the passwords to server, but it is conceivable that if we could not give anyone any passwords while maintaining the ability to log in, it'd be much safer.
- **Not succinct**: Succinctness is defined by the complexity of size of the proof in terms of the data being proven. In order to verify the proof, the verifier needs to compute the entire sequence of FibSq which means that the verifier does as much work as the prover, or $O(n)$ complexity. Succinctness with be discussed much later in the Groth16 chapter.

In our toy proof system, we will attempt to hide the knowledge. We will not try to make it succinct as it is much harder to shorten the proof than hiding the data.

## Hiding the Knowledge

The knowledge we want to hide in the FibSq example is $x$, the second element in the sequence. The prover wants to prove that they know it but without revealing it. There are many choices in terms of encryption algorithms but what we want is a specific kind that can encryt a message while preserving its algebraic structure, namely, the field of integers mod the prime number 3221225473.

One path of thinking could be:

- find a generator $g$ of $F_{p}^{\times}$, which is the multiplicative subgroup of the base field
- encrypt $x$ by exponentiation: $g^x$, this creates a Discrete Logarithm Problem (discussed in the previous chapter).

We also need to tweak the FibSq formula to accommodate this change

$g^{a_{i+2}} = g^{a_{i+1}^2} + g^{a_{i}^2}$, where the first and second element in the sequence are $g$ and $g^x$

The prover now sends $g^x$ to the verifier

The verifier computes

- $g^{a_2} = g^{1} \cdot g^{x^2}$
- $g^{a_3} = g^{x^2} \cdot g^{a_2^2}$
- $...$
- $g^{a_{1022}} = g^{a_{1010}^2} \cdot g^{a_{1011}^2}$

and accepts if the computed $g^{a_{1022}} = g^{2338775057}$

But ... there is a problem: the verifier can't compute $g^{x^2}$ or any $g^{a_i^2}$ because

1. they don't know $x$ or subsequent $a_i$ so they can't compute their square directly
2. this is an exponentiation on the exponent. $g^{x} \cdot g^{x} = g^{2x} \neq g^{x^2}$

Is there a way that allows us to perform multiplication on the encrypted values? Enter elliptic curves!

## Elliptic Curve

The goal of using curves is to replace the $g$ in $g^{a_{i+2}} = g^{a_{i+1}^2} + g^{a_{i}^2}$ with an elliptic curve point so that performing multiplicaon between the encrypted values $a_i$ is possible. As discussed in chapter 2, an elliptic curve over a finite field is defined by the function:

$y^3 = x^2 + ax + b \mod p$

### Addition Capability

In our FibSq formula, we have additions (subtraction is just adding the inverse of an element to another). Elliptic curves can accommodate this.

The set of elliptic curve points forms a group:

- Addition P + Q: take a line through P and Q, the sum is the inverse of the third point that the line intersects with the curve
- Identity: 0 (point at infinity)
- Additive Inverse: -(x, y) = (x, -y)

Scalar multiplication can be achieved by repeated addition $nP = P + P + ... + P \text{(n times)}$.

> Taking $n$ out of $nP$ is weirdly called the "discrete logarithm" instead of somehting like "discrete factoring". But if you really think about it, $nP$ is just a made up notation. The binary operator of "addition" is also completely defined by us (or by the mathematicians invented it). There is nothing stopping us from using the notation $P^n$. If we could define its written form this way, then it makes sense to call it a "discrete logarithm". This term is used mostly because it is widely used in cryptography to express that it's hard to separate the encrypted value from the result value.

### Multiplication Capability

Elliptic curves by themselves isn't useful enough to us. This is because we also have multiplication in the form of a square in the FibSq formula $a_{i+2} = a_{i+1}^2 + a_{i}^2$ (this is why we chose to have a slight twist of the original Fibonacci).

It turns out a specific kind of elliptic curves can "mimic" multiplication through a bilinear mapping of points (a.k.a pairing).

## Elliptic Curve Pairing

The EC pairing function $e(P, Q)$ is a bilinear mapping of the curve points $e: H_1 \times H_2 \rightarrow H_T$

### Properties

Digging up the theory behind pairings is beyond the scope of this article. But luckily we don't need to know what's behind pairings to use them. Let's figure out some properties of curve pairings:

Let $G_1$, $G_2$ and $G_T$ be the generator points of group 1, 2 and target respectively.

The pairing function over $H_1$ and $H_2$ curve points demonstrates associativity

$e(aG_1, bG_2) = e(abG_1, G_2) = e(G_1, abG_2) = abG_T$
> analogous to $a \cdot b = ab \cdot 1 = 1 \cdot ab$

They also demonstrates distributivity similar to multiplication

$e((a + b)G_1, cG_2) = e(aG_1, cG_2) + e(bG_1, cG_2) = (a + b)cG_T$ (distributive)
> analogous to $(a + b) \cdot c = ac + bc$

Not all elliptic curves have pairing capabilities. Most curves' embedding degree is too high for realistic pairing usage. If a curve can realistically do pairing, we call it a "pairing-friendly" curve. BN254 for example, is a pairing friendly curve used on the Ethereum virtual machine.

## Hiding the Knowledge (Cont.)

Now we have shown how elliptic curve pairing work, it should fix our problem of "not being able to multiply". Let's see how we can use them in our toy proof system.

Rewriting the prove/verify steps with EC:

1. Get a pairing-friendly curve
2. FibSq formula: $a_{i+2}G_T = e(a_{i+1}G_1, a_{i+1}G_2) + e(a_{i}G_1, a_{i}G_2)$

Prover sends $xG_1$ and $xG_2$ to the verifier

Verifier computes

$a_2G_T = e(1G_1, 1G_2) + e(xG_1, xG_2)$

all the way to

$a_{1022}G_T = e(a_{1010}G_1, a_{1010}G_2) + e(a_{1011}G_1, a_{1011}G_2)$

But there is a problem:

After performing $a_2G_T = e(1G_1, 1G_2) + e(xG_1, xG_2)$, we get $a_2$ encrypted in a $G_T$ point. But our subsequent computation needs to be done over $G_1$ and $G_2$ points because we want to use the pairing $e$ function. We also can't convert $G_T$ points to $G_1$ or $G_2$ points, and there is no pairing function for $G_T$ points. What do we do?

The situation now is that the prover is willing to prove, but the verifier can't verify it. Like in court, the judge is probably going to require the prosecutor to provide more evidence! The natural next step is: can we ask the prover to provide more than the thing they want to prove to support their claim?

Yes, and let's see how we could design that into our protocol.

Because the prover knows $x$, they can know the entire sequence of $a_0, a_1, ..., a_n$

The prover encrypts this sequence in both $G_1$ and $G_2$ and sends

$s_1 = \{a_0G_1, a_1G_1,...,a_{1022}G_1\}$

$s_2 = \{a_0G_2, a_1G_2,...,a_{1022}G_2\}$ where $a_1$ is just $x$

and the verifier checks

1. The first element and the 1022nd element in $s_i$ are $1G_1$ and $2338775057G_1$
2. The FibSq formula is satisfied

   $e(a_iG_1, a_iG_2) + e(a_{i+1}G_1, a_{i+1}G_2) = e(a_{i+2}G_1, G_2)$
    > analogous to $a_i^2 + a_{i+1}^2 = a_{i+2} \cdot 1$
3. The $a_i$ encrypted in both $s_1$ and $s_2$ are the same

   $e(a_iG_1, G_2) = e(G_1, a_iG_2)$
    > analogous to $a_i \cdot 1 = 1 \cdot a_i$

The verifier accepts the proof if all checks pass.

Notice that we now also have to check that the $a_i$ elements encrypted in the two vectors are the same. We let the prover send two vectors of $G_1$ and $G_2$ because the verifier needs to perform pairing for multiplication (the square operation), and that requires the same value but encrypted in both groups.

And we are done. The verifier can now verify the proof without an issue. In this process, the verifier has no idea of what value $x$ is because it's encrypted in a $G_1$ point where DLP exists. All they know is that the 1022nd element in the sequence is 2338775057.

### Witness

Also note that the prover only wanted to prove they know $x$ but supplied additional data in the proof. It turns out this additional data is so commonly used and important it has a special name, called "witness", for witnessing the computation of the sequence.

## Recap

### The Statement

"I know a number x such that
for a FibSq sequence with $a_0 = 1, a_1 = x$
we have $a_{1022} = 2338775057$"

### What Both the Prover and the Verifier Know

1. The elliptic curve and the generator $G_1$ and $G_2$
2. The 1022nd number in our FibSq sequence is $2338775057$

### The Proof

"I can provide vectors $s_1$ and $s_2$ where each element in those vectors represent an intermediate value in the FibSq sequence"

### The Verification

"If the values encrypted in the $s_1$ and $s_2$ you provided are consistent across these two vectors and satisfy the FibSq fomula and if the last element encrypts 2338775057, then indeed you know the answer to the value of $x$ in this setting"
