package actions

type OutputAction struct {
}

func (a *OutputAction) Execute(jsonPayload []byte) (string, error) {
	return "", nil
}
