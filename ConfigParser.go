package main

import (
	"errors"
	"fmt"
	"strings"
)

func (r *Router) ApplyConfig(config string) error {
	for i, line := range strings.Split(config, "\n") {
		var err error
		toks := strings.Fields(line)

		switch {
		case len(toks) == 0:
			//blank line
			continue
		case strings.HasPrefix(toks[0], "#"):
			//comment
			continue
		case strings.HasPrefix(toks[0], "node"):
			// fmt.Printf("%d node decl\n", i)
			err = r.NodeDecl(toks)

		case strings.HasPrefix(toks[0], "conn"):
			// fmt.Printf("%d conn decl\n", i)
			err = r.ConnDecl(toks)

		default:
			errstr := fmt.Sprintf("Unknown syntax, line %d:\n\t%s", i, line)
			err = errors.New(errstr)

		}
		if err != nil {
			// repackage with line num!
			fmt.Print(err)
			return err
		}
	}

	return nil
}

// syntax is node <name> ::= <nodetype> [args ...]
// returns nil if configuration applied succesfully
func (r *Router) NodeDecl(toks []string) error {
	if toks[2] != "::=" {
		return errors.New("Invalid syntax: missing '::=' operator")
	}

	// name := toks[1]
	// type, ok := NodeTypes[toks[3]]
	if !ok {
		return errors.New("Unknown node type " + toks[3])
	}
	// newnode := type(toks[4:])
	return nil
}

// syntax is conn <name> -> <name> [-> <name> ...]
func (r *Router) ConnDecl(toks []string) error {
	//first pass, check the syntax
	for i, tok := range toks {
		if tok == "->" {
			if (i-1 < 1) || (i+2 > len(toks)) {
				return errors.New("Mismatched '->' operator")
			} else {
				//check that the names are valid in the router
			}
		}
	}
	// second pass, if we got here everything is ok
	// so apply the configuration
	for i, tok := range toks {
		if tok == "->" {
			from := toks[i-1]
			to := toks[i+1]
			fmt.Printf("connecting %s to %s\n", from, to)
		}
	}
	return nil
}
