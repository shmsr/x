// Reference: https://eli.thegreenplace.net/2018/the-expression-problem-in-go/
package main

import (
	"fmt"
	"strconv"
)

type expr interface{}

type eval interface {
	Eval() float64
}

type toString interface {
	ToString() string
}

type constant struct {
	value float64
}

type binaryPlus struct {
	left  expr
	right expr
}

type binaryMinus struct {
	left  expr
	right expr
}

type binaryMultiply struct {
	left  expr
	right expr
}

type binaryDivide struct {
	left  expr
	right expr
}

func (c *constant) Eval() float64 {
	return c.value
}

func (c *constant) ToString() string {
	return strconv.FormatFloat(c.value, 'f', -1, 64)
}

func (bp *binaryPlus) Eval() float64 {
	return bp.left.(eval).Eval() + bp.right.(eval).Eval()
}

func (bp *binaryPlus) ToString() string {
	l := bp.left.(toString)
	r := bp.right.(toString)
	return fmt.Sprintf("(%s + %s)", l.ToString(), r.ToString())
}

func (bp *binaryMinus) Eval() float64 {
	return bp.left.(eval).Eval() - bp.right.(eval).Eval()
}

func (bp *binaryMinus) ToString() string {
	l := bp.left.(toString)
	r := bp.right.(toString)
	return fmt.Sprintf("(%s - %s)", l.ToString(), r.ToString())
}

func (bp *binaryMultiply) Eval() float64 {
	return bp.left.(eval).Eval() * bp.right.(eval).Eval()
}

func (bp *binaryMultiply) ToString() string {
	l := bp.left.(toString)
	r := bp.right.(toString)
	return fmt.Sprintf("(%s * %s)", l.ToString(), r.ToString())
}

func (bp *binaryDivide) Eval() float64 {
	return bp.left.(eval).Eval() / bp.right.(eval).Eval()
}

func (bp *binaryDivide) ToString() string {
	l := bp.left.(toString)
	r := bp.right.(toString)
	return fmt.Sprintf("(%s / %s)", l.ToString(), r.ToString())
}

func createNewExpr() expr {
	a := constant{value: 1.0}
	b := constant{value: 2.0}
	c := constant{value: 3.0}
	d := constant{value: 4.0}

	// (a + (b - (c * (d / a)))) = (pow(a, 2) + ab - cd)/a
	bp := binaryPlus{
		&a, &binaryMinus{
			&b, &binaryMultiply{
				&c, &binaryDivide{&d, &a},
			},
		},
	}
	return &bp
}

func main() {
	expr := createNewExpr()
	fmt.Println("Eval: ", expr.(eval).Eval())           // Eval:  -9
	fmt.Println("String: ", expr.(toString).ToString()) // String:  (1 + (2 - (3 * (4 / 1))))
}
