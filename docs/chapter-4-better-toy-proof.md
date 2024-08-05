# Chapter 4: Improving Our Toy Proof

## Recap

In the last chapter, we introduced an example involving proving the statement:

"I know a number x such that
for a FibSq sequence with $a_0 = 1, a_1 = x$,
$a_{1022} = 2338775057$"

We came up with the simplest proof which is just giving away the value $x$. We then introduced encryption using elliptic curves to hide $x$ while maintaining the ability for the prover to prove that they know $x$.

Our proof system works (even though it's not succinct yet), but it still has a big problem: it is extremely limited as it cannot handle degrees higher than 2. Let's see how a simple tweak can break our entire system.

Consider this new Fibonacci formula:

$a_{i+2} = a_{i+1}^2 + a_{i}^{\color{red} 3}$

If we were to keep using our previous constraint form to accommodate this, it would need something that looks like:

$e(a_{i+2}G_1, G_2) = e(a_{i+1}G_1, a_{i+1}G_2) + e(a_iG_1, a_iG_2, {\color{red} a_iG_?})$

This is asking for a trillinear pairing function, but we don't have that over elliptic curve points. And even if we do, it's hard to have both the bilinear pairing and trilinear one map to the same target group. Also, what if the degree goes to 4? We are really asking for a variadic pairing function and it's infeasible.

The good news is, instead of finding a solution in the ground level mathematical constructs, we can fix that in our protocol.

## Witness to the Rescue (Again)

The idea is that we can "flatten" the computation of any term of $deg > 2$ by computing it's quadratic components one at a time. Taking our tweaked Fib formula as an example:

$a_{i+2} = a_{i+1}^2 + {\color{red} a_i^3}$

$~~\downarrow$

$v_i = a_i^2$

$a_{i+2} = a_{i+1}^2 + {\color{red} v_i \cdot a_i}$

Notice that the $a_i^3$ term can be factored to $a_i^2 \cdot a_i$. And if we ask the prover to compute $v_i = a_i^2$ and encrypt $v_i$ into curve points, then the verifier can compute the $a_i^3$ themselves which is just $v_i \cdot a_i$. Let's apply this in the example.

Another way to seeing this is: through the introduction of intermediate variable $v_i$, we are essentially "rewriting" the Fib formula in an equivalent way with a "redundant" step.

### New Proof

$s_1 = \{a_iG_1\}_{i=0}^{1022}$

$s_2 = \{a_iG_2\}_{i=0}^{1022}$

$\color{red} v = \{v_iG_1\}_{i=0}^{1020}$

### New Verification Steps

Verifier checks:

The computation of the intermediate variables $v_i$ is correct
  
$\color{red} e(v_iG_1, G_2) = e(a_iG_1, a_iG_2)$
> analogous to $v_i = a_i^2$

The recurring relations are correct
  
$e(a_{i+2}G_1, G_2) = e(a_{i+1}G_1, a_{i+1}G_2) + {\color{red} e(v_iG_1, a_iG_2)}$
> analogous to $a_{i+2} = a_{i+1}^2 + v_i \cdot a_i$

The $a_i$ encrypted in both $s_1$ and $s_2$ are the same

$e(a_iG_1, G_2) = e(G_1, a_iG_2)$
> analogous to $a_i \cdot 1 = 1 \cdot a_i$

The verifier accepts if all checks pass.

Now we have a working proof system that not only hides the values of $x$, but also can handle arbitrary degree equations (we didn't show this but I hope it's easy to see how it's easy through flattening). Now, the only trait left is succinctness. In fact, through our "improvements" over the simplest proof, even though the verifying complexity is still O(n), we actually worsened it by requiring the verifier to do pairings over and over through the entire sequence. In the next chapter, we will go over the proof system proposed in the Groth16 paper and see how it achieves succinctness.
