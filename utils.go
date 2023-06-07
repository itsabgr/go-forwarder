package main

func must[R any](r R, e error) R {
	if e != nil {
		panic(e)
	}
	return r
}
