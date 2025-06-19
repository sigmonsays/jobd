package job

func NewVars(prefix string) *Vars {
	return &Vars{
		Prefix: "",
		Global: make(map[string]Var, 0),
	}
}

type Var string

// variable management for jobs and steps
// top level vars struct
type Vars struct {
	Prefix string
	Global map[string]Var
}

func (me *Vars) SetGlobal(k string, v string) {
	me.Global[me.Prefix+k] = Var(v)
}
