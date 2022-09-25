package action

import "testing"

type leoAction struct {
}

func (la leoAction) Execute(json []byte) (string, error) {
	return "leo", nil
}

func TestActionRepository(t *testing.T) {
	ar := NewActionRepository()

	ar.Register("leo", &leoAction{})
	la := ar.Retrieve("leo")
	s, e := la.Execute([]byte{})

	if e != nil {
		t.Fatalf("Execute produced an error: %+v", e)
	}
	if s != "leo" {
		t.Fatalf("Execute returned something else instead of `leo`: %s", s)
	}
}
