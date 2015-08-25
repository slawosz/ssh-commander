package communication

type CommIn interface {
	StartInput()
}

type CommOut interface {
	StartOutput()
}
