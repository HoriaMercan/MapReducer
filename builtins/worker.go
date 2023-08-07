package builtins

type Student struct {
	Age    int
	Name   string
	School string
}

type Args struct{}
type Reply struct {
	Value int
}

func (S *Student) GetAge() int {
	return S.Age
}

func (S *Student) GetName() string {
	return S.Name
}

func (s *Student) GetSchool() string {
	return s.School
}

func (s *Student) GetAll(args *Args, reply *Reply) error {
	reply.Value = s.Age
	return nil
}
