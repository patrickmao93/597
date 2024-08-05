# Chapter 7: Arithmetic Circuits

In the previous chapters we have already come across many examples of writing arithmetic circuits (on paper). But most of our examples were cherry-picked simple ones to help introduce other concepts first. The truth is, almost nothing in the real world is as simple as checking if a Fibonacci sequence is correct. Speaking from experience in the blockchain industry, interesting things that have actual applications are topics such as "checking if some data is the preimage of a hash" or "checking if a message is signed by some party". Purposefully designed arithmetic circuits are written to accomplish these types of proofs.

## Arithmetic Circuits

Arithmetic circuits are just a way to express the verification logic in arithmetics, and in the context of Groth16, is something that is transformable to R1CS. Note that in the Groth16 paper it never restricts the type of arithmetization one should use. Anything that can be turned into an R1CS works.

## Arithmetization

Arithmetization is a step where we "translate" verification logic into arithmetics.

### Example

"I have the preimage of this SHA-3 hash"

We work with equations like $y = 3x^3 + 2$, how do we turn SHA-3 into these? We know SHA-3 is just a bunch of AND, XOR, LROT, RROT. If we can represent these operations with arithmetics, then we can arithmetize SHA-3.

## AND

I think the truth table for AND is familiar enough so I will skip describing these. Here is the table:

| x   | y   | z   |
| --- | --- | --- |
| 0   | 0   | 0   |
| 0   | 1   | 0   |
| 1   | 0   | 0   |
| 1   | 1   | 1   |

Since we have been using the word "circuit", let's also make the terms correct for AND. In this AND "gate", we have  $x, y$ as input signals and $z$ as an output signal

In circuit, we can write down $z = xy$ as a constraint.

## XOR

| x   | y   | z   |
| --- | --- | --- |
| 0   | 0   | 0   |
| 0   | 1   | 1   |
| 1   | 0   | 1   |
| 1   | 1   | 0   |

In circuit: $z = x + y - 2xy$

But, there is a catch.

Though this arithmetization seems complete, it still has a "bug".

What if the input signals x and y are 2 and 3?

Substitute the input in, we get $z = 2 + 3 - 2 \cdot 2 \cdot 3 = -7$, which is totally wrong! Our XOR gate isn't doing XOR's job at all! Our input signals must be either 0 or 1 for the XOR gate, and we must also check this attribute.

## Range Check

Range checking itself is an arithmetization step. Checking binary digits is easy:

$(x - 1)x=0$

Generalizing to checking that the range of x is in [m, n] where m < n

$(x - n)(x - n - 1)(x - n - 2)...(x - m)=0$

The LHS is just the vanishing polynomial in range [m, n]

And with that, our AND and XOR gates can now be fully constrained:

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

## Some Big Caveats

Arithmetic circuits can really catch a traditional programmer off-guard as they are so different from ordinary programs. In a way, they are really more like a circuit you build physically with wires than a computer program. Here are some big "traps" that programmers like me fell hard for when learning to write circuits for the first time.

- There is no dynamic array
- There is no loop
- There is no if-else

Dynamic array handling is a semi-advanced topic and is beyond this article's scope. Loops are borderline impossible without building a virtual machine from circuits. So I will skip discussing these.

If-else on the other hand is a much more applicable and useful construct, and we will dive deep into it.

## Signal Selector (if-else)

For the sake of keeping thing "circuity", let's call our if-else a "signal selector". What the selector does is really picking an output signal from two input signals $(x, y)$ depending on the on/off state of a control signal $c$.

$z = select(c, x, y)$

### Arithmetization

$z = c \cdot x + isZero(c) \cdot y$

$(c - 1)c = 0$

where

$isZero: y = 1 - x$

And that's it. This totally works. But I want to make things more interesting. Notice how our $isZero$ gate can only handle 0 and 1? What if we want to achieve C-like if-else where 0 means false and everything else means true? Let's try generalizing $isZero$.

## IsZero (Generalized)

$y = isZero(x)$. Outputs 1 when $x = 0$, and 0 otherwise.

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

Ignore the "if $x \neq 0$" part for a second, substitute $x$ with any numbers and we should see that $isZero$ works the way it's intended. It always outputs 1 when $x = 0$ and 0 when $x \neq 0$.

Now let's address the "if $x \neq 0$" part. What is this "if"?

Since we are using $isZero$ in $select$, wouldn't it cause infinite recursion? Well, the answer is that the "if" here isn't a circuit component, it's just a plain "if statement" you would write using whatever programming language that builds the circuit. It means that whoever wants to prove that their input signals, e.g. $y=0,x=1$, is correct, has to compute $m$ outside their circuit. But wouldn't this mean the prover can provide whatever they want? That is true, and it's a problem. Let's take a moment and see what would happen if we leave this as is and let the prover just **claim** the value of $m$.

A malicious prover have $x = 10$, which is not zero and the correct output should be 0, but because they can just claim $m$, they can supply $m=0$ and making the output 1 instead. We can't let the prover freely control the value of $m$. Therefore, we add the following constraint to **restrict** what the prover can claim.

$$
x \cdot out = 0
$$

This constraint actually partially describes the nature of our $isZero$ gate:

$x$ and $out$ can't both be 1 (fixes the malicious example above)

But what about $x$ and $out$ both being 0? This constraint would be satisfied, but it's not correct. Well, that case is covered by the previous constraint $out = -m \cdot x + 1$. It is easy to see that $x$ and $out$ can never both be zero.

This concept of letting the prover simply "claim" something and then constrain the things they claim has actually appeared in the previous chapters. When the verifier can't compute the Fibonacci sequence on their own, we asked the prover to provide more data (essentially claiming these data), then ask the verifier to check whether these data are correctly claimed using verification logic.
