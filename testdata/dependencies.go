package testdata

type TestOne struct {
	Name string
}

func NewTestOne() TestOne {
	return TestOne{"testone"}
}

type TestTwo struct {
	Name string
}

func NewTestTwo() TestTwo {
	return TestTwo{"test two"}
}

type TestThree struct {
	One TestOne
	Two TestTwo
}

func NewTestThree(t1 TestOne, t2 TestTwo) TestThree {
	return TestThree{t1, t2}
}
