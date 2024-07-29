package demo

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"testing"
)

func TestDemo(t *testing.T) {
	// Setup

	// compile the circuit to a constraint system
	constrSys, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &Circuit{})
	check(err)

	// Dummy GROTH16 trusted setup
	// pk and vk has the circuit logic "imprinted" in them
	pk, vk, err := groth16.Setup(constrSys)
	check(err)

	// Proving

	// circuit input assignment
	// u = x ^ y ⊕ z
	// 1 = 1 ^ 0 ⊕ 1
	assignment := Circuit{
		Z: 1, U: 1, X: 1, Y: 1,
	}

	// Compute the witness vector (the array of intermediate variables in the computation)
	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	check(err)
	fmt.Println("witness", witness.Vector())

	publicWitness, err := witness.Public()
	fmt.Println("public witness", publicWitness.Vector())
	check(err)

	// GROTH16: Prove
	proof, err := groth16.Prove(constrSys, pk, witness)
	check(err)

	// GROTH16: Verify
	err = groth16.Verify(proof, vk, publicWitness)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
