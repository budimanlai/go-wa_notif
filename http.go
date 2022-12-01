package main

import (
	"time"

	"github.com/valyala/fasthttp"
)

func Get(url string, query string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	if len(query) > 0 {
		req.SetRequestURI(url + "?" + query)
	} else {
		req.SetRequestURI(url)
	}

	req.Header.DisableNormalizing()

	req.Header.SetContentType(`application/json`)
	req.Header.SetMethod(fasthttp.MethodGet)

	respClone := &fasthttp.Response{}
	e := fasthttp.DoTimeout(req, resp, 60*time.Second)
	resp.CopyTo(respClone)

	return respClone, e
}
