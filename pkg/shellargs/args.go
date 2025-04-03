package shellargs

type Args struct {
	args [][]string
}

func New() *Args {
	return &Args{
		args: make([][]string, 0),
	}
}

func (a *Args) Add(value string) *Args {
	a.args = append(a.args, []string{value})
	return a
}

func (a *Args) AddKeyed(key, value string) *Args {
	a.args = append(a.args, []string{key, value})
	return a
}

func (a *Args) AddKeyedSingleQuoted(key, value string) *Args {
	a.args = append(a.args, []string{key, "'" + value + "'"})
	return a
}

func (a *Args) String() string {
<<<<<<< HEAD
	// fmt.Print(len(a.args))
=======
>>>>>>> main
	output := ""
	for _, values := range a.args {
		if len(values) == 2 {
			output = output + " " + values[0] + " " + values[1]
		} else {
			output = output + " " + values[0]
		}
	}
	return output
}
