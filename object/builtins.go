package object

import "fmt"

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		"len",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				switch arg := args[0].(type) {
				case *String:
					return &Integer{Value: int64(len(arg.Value))}
				case *Array:
					return &Integer{Value: int64(len(arg.Elements))}
				default:
					return newError("argument to `len` not supported, got %s", args[0].Type())
				}
			},
		},
	},
	{
		"puts",
		&Builtin{
			Fn: func(args ...Object) Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return nil
			},
		},
	}, {
		"first",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to `first` must be %s, got %s", ARRAY_OBJ, args[0].Type())
				}
				arg := args[0].(*Array)
				if len(arg.Elements) > 0 {
					return arg.Elements[0]
				}
				return nil
			},
		},
	}, {
		"last",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to `last` must be %s, got %s", ARRAY_OBJ, args[0].Type())
				}
				array := args[0].(*Array)
				arrayLen := len(array.Elements)
				if arrayLen > 0 {
					return array.Elements[arrayLen-1]
				}
				return nil
			},
		},
	}, {
		"rest",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				if args[0].Type() != ARRAY_OBJ {
					return newError("argument to `rest` must be %s, got %s", ARRAY_OBJ, args[0].Type())
				}
				array := args[0].(*Array)
				arrayLen := len(array.Elements)
				if arrayLen == 0 {
					return nil
				}
				elementRest := make([]Object, arrayLen-1)
				copy(elementRest, array.Elements[1:])
				return &Array{
					Elements: elementRest,
				}
			},
		},
	},
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}

func newError(format string, a ...any) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
