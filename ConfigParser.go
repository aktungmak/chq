package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
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

	name := toks[1]

	_, ok := AvailableNodes[toks[3]]
	if !ok {
		return errors.New("Unknown node type " + toks[3])
	}

	// extract the constructor function as a relfelct.Value
	nodectr := reflect.ValueOf(AvailableNodes[toks[3]])
	nodetype := reflect.TypeOf(AvailableNodes[toks[3]])

	// now work out the types of the arguments
	args := make([]Value, 0)
	for i := 0; i < nodetype.NumIn(); i++ {

		switch nodetype.In(i) {
		case reflect.Bool:
			val := strconv.ParseBool(toks[3+i])
		case reflect.Int:
			val := strconv.ParseInt(toks[3+i])
		case reflect.Float64:
			val := strconv.ParseFloat(toks[3+i])
		case reflect.String:
			val := toks[3]
		default:
			return errors.New("Invalid argument in node declaration: " + "TODO value")
		}
		args = append(args, reflect.TypeOf(val))
	}

	// finally call the constructor to ge the new node value
	nodectr.Call(args)

	name = name
	return nil
}

// syntax is conn <name> -> <name> [-> <name> ...]
func (r *Router) ConnDecl(toks []string) error {
	//first pass, check the syntax
	for i, tok := range toks {
		if tok == "->" {
			if (i-1 < 1) || (i+2 > len(toks)) {
				return errors.New("Mismatched '->' operator")
				{
				}
			}
			// second pass, if we got here everything is ok
			// so apply the configuration
			for i, tok := range toks {
				if tok == "->" {
					src := toks[i-1]
					dst := toks[i+1]
					fmt.Printf("connecting %s to %s\n", src, dst)
					return r.Connect(src, dst)
				}
			}
		}
	}
	return nil
}
