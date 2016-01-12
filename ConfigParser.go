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
			// TODO repackage with line num!
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
	kind := toks[3]
	args := toks[4:]

	return r.CreateNode(name, kind, args...)

}

func (r *Router) CreateNode(name, kind string, args ...string) error {
	n, ok := AvailableNodes[kind]
	if !ok {
		return errors.New("Unknown node type " + kind)
	}

	// extract the constructor function as a relfelct.Value
	nodectr := reflect.ValueOf(n)
	nodetype := reflect.TypeOf(n)

	// make sure the arity of the constructor matches
	// the number of args provided
	if nodetype.NumIn() != len(args) {
		return errors.New("Incorrect number of args provided!")
	}

	// now work out the types of the arguments
	var err error
	var val interface{}
	var node Routeable
	rargs := make([]reflect.Value, 0)
	for i := 0; i < nodetype.NumIn(); i++ {
		switch nodetype.In(i).Kind() {
		case reflect.Bool:
			val, err = strconv.ParseBool(args[i])
		case reflect.Int:
			val, err = strconv.ParseInt(args[i], 10, 32)
			val = int(val.(int64))
		case reflect.Int64:
			val, err = strconv.ParseInt(args[i], 10, 64)
		case reflect.Float64:
			val, err = strconv.ParseFloat(args[i], 64)
		case reflect.String:
			val = args[i]
			err = nil
		default:
			return errors.New("Invalid argument in node declaration: " + args[i])
		}
		if err != nil {
			return err
		} else {
			rargs = append(rargs, reflect.ValueOf(val))
		}
	}

	// finally call the constructor to get the new node value
	res := nodectr.Call(rargs)
	node = res[0].Interface().(Routeable)

	if res[1].Interface() != nil {
		return res[1].Interface().(error)
	} else {
		// if all was ok, register it with the router
		return r.RegisterNode(name, node)
	}
}

// syntax is conn <name> -> <name> [-> <name> ...]
func (r *Router) ConnDecl(toks []string) error {
	//first pass, check the syntax
	for i, tok := range toks {
		if tok == "->" {
			if (i-1 < 1) || (i+2 > len(toks)) {
				return errors.New("Mismatched '->' operator")
			}
		}
	}
	// second pass, if we got here everything is ok
	// so apply the configuration
	for i, tok := range toks {
		if tok == "->" {
			src := toks[i-1]
			dst := toks[i+1]
			err := r.Connect(src, dst)
			if err != nil {
				return err
			}
			fmt.Printf("connected %s to %s\n", src, dst)
		}
	}
	return nil
}
