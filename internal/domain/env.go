package domain

type Env struct {
	Id    int
	Name  string
	Value string
}

func NewEnv(name string, value string) (*Env, error) {
	env := Env{
		Name:  name,
		Value: value,
	}

	return &env, nil
}
