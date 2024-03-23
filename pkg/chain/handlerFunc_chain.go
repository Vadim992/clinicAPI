package chain

import (
	"fmt"
	"net/http"
)

type MW func(handlerFunc http.HandlerFunc) http.HandlerFunc

type HandleFuncChain struct {
	mwSlice []MW
}

func New(funcs ...MW) *HandleFuncChain {
	fmt.Println(len(funcs))
	mwSlice := make([]MW, len(funcs))
	mwSlice = append(mwSlice, funcs...)

	return &HandleFuncChain{mwSlice: mwSlice}
}

func (c *HandleFuncChain) Then(f http.HandlerFunc) http.HandlerFunc {
	if f == nil {
		return nil
	}
	fmt.Println(len(c.mwSlice) - 1)
	for i := len(c.mwSlice) - 1; i > 0; i-- {
		fmt.Println("here")
		handler := c.mwSlice[i]
		f = handler(f)
	}

	return f
}

func (c *HandleFuncChain) Append(mw MW) {
	c.mwSlice = append(c.mwSlice, mw)
}
