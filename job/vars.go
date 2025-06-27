package job

func NewVars(prefix string) *Vars {
	return &Vars{
		Prefix: "",
		Global: NewVarSet(),
		Step:   make(map[string]*VarSet, 0),
	}
}

type Var string

func (me Var) String() string {
	return string(me)
}

// variable management for jobs and steps
// top level vars struct
type Vars struct {
	Prefix string
	Global *VarSet
	Step   map[string]*VarSet
}

func (me *Vars) GetStep(step string) *VarSet {
	vars, found := me.Step[step]
	if !found {
		vars = NewVarSet()
		me.Step[step] = vars
	}
	return vars
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
