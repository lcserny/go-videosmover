package videosmover

type SearchAction struct {
}

func (a *SearchAction) Execute(jsonFile string) (string, error) {
	return "Arg sent = " + jsonFile + ", Something back", nil
}
