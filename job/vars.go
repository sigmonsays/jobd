package job

func NewVars(prefix string) *Vars {
	return &Vars{
		Prefix: "",
		Global: NewVarSet(),
		Step:   make(map[string]*VarSet, 0),
	}
}

type Var string

// variable management for jobs and steps
// top level vars struct
type Vars struct {
	Prefix string
	Global *VarSet
	Step   map[string]*VarSet
}

func (me *Vars) SetStep(step string, k string, v string) {
	vars, found := me.Step[step]
	if !found {
		vars = NewVarSet()
	}

	vars.SetVar(k, v)
}
func NewVarSet() *VarSet {
	return &VarSet{
		Vars: make(map[string]Var, 0),
	}
}

type VarSet struct {
	Vars map[string]Var
}

func (me *VarSet) SetVar(k, v string) {
	me.Vars[k] = Var(v)
}
