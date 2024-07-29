package demo

import "github.com/consensys/gnark/frontend"

type Circuit struct {
	Z frontend.Variable `gnark:",public"`
	U frontend.Variable `gnark:",public"`
	X frontend.Variable
	Y frontend.Variable
}

func (c *Circuit) Define(api frontend.API) error {
	// asserts all inputs are binary
	assertBinary(api, c.X)
	assertBinary(api, c.Y)
	assertBinary(api, c.Z)
	assertBinary(api, c.U)

	// Circuit: u = x ^ y ⊕ z

	// perform AND: a * b
	and := api.Mul(c.X, c.Y)

	// perform XOR: a + b - 2ab
	rhs := api.Sub(
		api.Add(and, c.Z),    // a + b
		api.Mul(2, and, c.Z), // 2ab
	)

	// asserts the relationship u = x ^ y ⊕ z
	api.AssertIsEqual(c.U, rhs)

	return nil
}

func assertBinary(api frontend.API, x frontend.Variable) {
	// (x - 1)x = 0
	res := api.Mul(api.Sub(x, 1), x)
	api.AssertIsEqual(res, 0)
}
