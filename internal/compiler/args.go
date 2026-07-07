package compiler

type Options struct {
	Compiler  *Compiler
	Includes  []string
	Defines   []string
	ExtraArgs []string
	OutDir    string
	NoCache   bool
	Count     int
}
