package main

import (
	"context"
	"strings"
)

type abProcessor struct {
	outA chan Aout
	outB chan Bout
	errs chan error
}

func newABProcessor() *abProcessor {
	return &abProcessor{
		outA: make(chan Aout, 1),
		outB: make(chan Bout, 1),
		errs: make(chan error, 2),
	}
}

type Aout struct {
	Original string
	Lower    string
	Reversed string
}

type Bout struct {
	Original string
	Upper    string
	Compact  string
}

type cIn struct {
	a Aout
	b Bout
}

func (p *abProcessor) start(ctx context.Context, data Input) {
	go func() {
		Aout, err := getResultA(ctx, data.A)
		if err != nil {
			p.errs <- err
			return
		}
		p.outA <- Aout
	}()

	go func() {
		Bout, err := getResultB(ctx, data.B)
		if err != nil {
			p.errs <- err
			return
		}
		p.outB <- Bout
	}()
}

func (p *abProcessor) wait(ctx context.Context) (cIn, error) {
	var cData cIn
	for count := 0; count < 2; count++ {
		select {
		case a := <-p.outA:
			cData.a = a
		case b := <-p.outB:
			cData.b = b
		case err := <-p.errs:
			return cIn{}, err
		case <-ctx.Done():
			return cIn{}, ctx.Err()
		}
	}
	return cData, nil
}

func getResultA(ctx context.Context, in string) (Aout, error) {
	if err := ctx.Err(); err != nil {
		return Aout{}, err
	}

	n := strings.ToLower(strings.TrimSpace(in))
	reversed := reverseString(n)

	return Aout{
		Original: in,
		Lower:    n,
		Reversed: reversed,
	}, nil
}

func getResultB(ctx context.Context, in string) (Bout, error) {
	if err := ctx.Err(); err != nil {
		return Bout{}, nil
	}
	words := strings.Fields(in)
	compact := strings.Join(words, " ")
	upper := strings.ToUpper(compact)

	return Bout{
		Original: in,
		Upper:    upper,
		Compact:  compact,
	}, nil
}

func reverseString(in string) string {
	runes := []rune(in)
	for l, r := 0, len(runes)-1; l < r; l, r = l+1, r-1 {
		runes[l], runes[r] = runes[r], runes[l]
	}
	return string(runes)
}
