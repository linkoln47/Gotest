package main

import (
	"context"
	"unicode"
)

type cProccessor struct {
	outC chan Cout
	errs chan error
}

func newCProcessor() *cProccessor {
	return &cProccessor{
		outC: make(chan Cout, 1),
		errs: make(chan error, 1),
	}
}

func (p *cProccessor) start(ctx context.Context, inputC cIn) {
	go func() {
		Cout, err := getResultC(ctx, inputC)
		if err != nil {
			p.errs <- err
			return
		}
		p.outC <- Cout
	}()
}

func (p *cProccessor) wait(ctx context.Context) (Cout, error) {
	select {
	case out := <-p.outC:
		return out, nil
	case err := <-p.errs:
		return Cout{}, err
	case <-ctx.Done():
		return Cout{}, ctx.Err()

	}
}

func getResultC(ctx context.Context, c cIn) (Cout, error) {
	if err := ctx.Err(); err != nil {
		return Cout{}, nil
	}

	combined := c.a.Reversed + " | " + c.b.Upper

	freqCount := make(map[rune]int)

	for _, r := range combined {
		if unicode.IsSpace(r) || r == '|' {
			continue
		}

		r = unicode.ToLower(r)
		freqCount[r]++
	}

	return Cout{
		freqCount: freqCount,
	}, nil
}
