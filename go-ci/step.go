package main

import "os/exec"

type step struct {
	name    string
	exe     string
	args    []string
	message string
	project string
}

// a constructor function to initialize the step struct
func newStep(name, exe, message, project string, args []string) step {
	return step{
		name:    name,
		exe:     exe,
		args:    args,
		message: message,
		project: project,
	}
}

func (s step) Execute() (string, error) {
	cmd := exec.Command(s.exe, s.args...)
	cmd.Dir = s.project

	if err := cmd.Run(); err != nil {
		return "", &stepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err,
		}
	}

	return s.message, nil
}
