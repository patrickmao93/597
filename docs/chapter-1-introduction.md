# Chapter 1: Introduction

## Motivation and Background

In my day job before I decided to take CS597 I mainly dealt with zero-knowledge proof applications in the context of blockchains. That is, writing high-level circuit applications that interacts with smart contracts and other blockchain systems. Because of how novel the field was (and still is), I often find myself reading papers that are rather arcane and obviously not written with a goal of educating regular software engineers. During my own study before taking CS597, I found that I mainly lack knowledge in three subjects: abstract algebra, cryptography in general, and the underlying theory of the proof system that I use which are Groth16 and STARK. The traits of these ZK proof systems, their elegance and the ability to create seemingly "impossible" proof constructs, drawn me to this adventure of learning more about them in an independent study course.

In this course of independent study, I want to acquire a long overdue understanding in the theory backbone of Groth16, and write them down in a tone of teaching to serve both as notes of my own and a piece additional material that others may find helpful while learning about ZKPs.

## Zero-Knowledge Proofs (ZKPs), ZK-SNARKs, and Groth16

If you are an avid blockchain user in the past few years you must have heard of these buzz words already. But what do they mean and how are they conceptually different?

### Zero-Knowledge Proofs

Zero-knowledge proof refers to the concept of proof systems that aim to enable a prover to prove the something (data) is true to a verifier without revealing the data itself. The application of such constructs are currently being seen mostly in privacy and trustless computing, mostly in the blockchain space. A couple simple use cases before we dive into the nitty-gritty:

1. ZKPs can be used to making a blockchain aware of its historical states by proving the correctness of the signatures on previously signed blocks. This enables smart contracts to gain access to much more data than they currently have.
2. ZKPs has also already been used (although in a arguably unethical way) in projects like Tornado Cash to keep information private on a completely transparent blockchain.

To demonstrate proving something without revealing the thing itself is possible, we could use a simple example: imagine that you have to convince a blind person that you have a piece of paper and it has two colors on it (say the top half is red and the bottom is blue). Is it possible? The answer is, of course, yes. And it can be done in the following steps (from the blind person's perspective):

1. Receives a piece of paper from the prover
2. Randomly shuffles the orientation of the paper under the desk
3. Show the paper back to the prover and ask: "did the color switch positions?"
4. Repeat step 1-3 until convinced that there is no way the prover get it right by pure chance every time.

This is also called a _probabilistically checkable proof_ or _PCP_. The more times the above steps are repeated and pass, the more probability that the claim "there are two colors on the paper" is true.

### ZK-SNARK

ZK-SNARK is considered a type of ZKP system and the acronym stands for _Zero-Knowledge Succinct Non-interactive ARgument of Knowledge_. That's a mouthful, let's break down it term by term and see what it actually means.

1. Zero-knowledge, as explained above, means that someone can prove something is true without revealing anything about that something.
2. Succinct means that the proof is short such that whoever wants to verify the truthfulness of the proof can do it in a short amount of time.
3. Non-interactive is an interesting one. In the above "blind person" example, it requires the blind person to "interact" or ask the prover questions in order for them to be confident about the claim's truthfulness. But this process requires both parties (prover and verifier) to be present to execute. The "non-interactive" property allows the prover to prove without the verifier personally asking questions.
4. Argument of Knowledge is a more correct and precise way than saying "proof of knowledge". We can't say "proofs" because proofs need to be absolutely certain. It means no matter what powerful computer an adversarial has, there is no chance that they can find a false positive proof. On the other hand, in our previous "blind person" example we can already see that there is a chance, though neglegible after enough trials, that someone can provide a correct proof by pure luck. It makes argument of knowledge systems unable to handle adversarials that have unbounded computational power (can try an infinite amount of proofs). But do note that in practice, the two terms are often used interchangably.

### Groth16

Groth16 is the code name of a famous paper "On the Size of Pairing-based Non-interactive Arguments" published by Jens Groth in 2016. The paper proposed a ZK proof system that reduces the proof size to the holy grail of complexity of O(1) with respect to the input data size while maintaining a reasonable the proving complexity.

#### Why do we want a small proof size?

The answer is that we actually want a small verification complexity but often times it is assumed that verification complexity is linear to the proof size so these terms are interchangable.

#### Why do we want a low verification complexity?

It is so that we can execute verification logic on resource-limited computers, as of the time of writing, mostly blockchains.

#### What does an O(1) complexity proof size actually mean concretely?

It means no matter the amount of data you are trying to prove, the cost of verification never changes. For example, imagine proving to a tax agency that your taxes are done right while having 10,000 employees under the record. O(1) proof size enables them to check your taxes in as much time as checking the same for one person. And their effort of checking your proof never changes no matter how many more employees you decide to hire.
<script type="text/javascript" src="http://cdn.mathjax.org/mathjax/latest/MathJax.js?config=TeX-AMS-MML_HTMLorMML"></script>
<script type="text/x-mathjax-config"> MathJax.Hub.Config({ tex2jax: {inlineMath: [['$', '$']]}, messageStyle: "none" });</script>
