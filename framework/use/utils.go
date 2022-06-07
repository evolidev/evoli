package use

func AbortUnless(e interface{}) {
	if e != nil && e != false {
		panic(e)
	}
}
