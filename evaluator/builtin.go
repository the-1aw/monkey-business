package evaluator

import (
	"fmt"

	"github.com/the-1aw/monkey-business/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be %s, got %s", object.ARRAY_OBJ, args[0].Type())
			}
			arg := args[0].(*object.Array)
			if len(arg.Elements) > 0 {
				return arg.Elements[0]
			}
			return NULL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` must be %s, got %s", object.ARRAY_OBJ, args[0].Type())
			}
			array := args[0].(*object.Array)
			arrayLen := len(array.Elements)
			if arrayLen > 0 {
				return array.Elements[arrayLen-1]
			}
			return NULL
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` must be %s, got %s", object.ARRAY_OBJ, args[0].Type())
			}
			array := args[0].(*object.Array)
			arrayLen := len(array.Elements)
			if arrayLen == 0 {
				return NULL
			}
			elementRest := make([]object.Object, arrayLen-1)
			copy(elementRest, array.Elements[1:])
			return &object.Array{
				Elements: elementRest,
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 2 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `push` must be %s, got %s", object.ARRAY_OBJ, args[0].Type())
			}
			array := args[0].(*object.Array)
			arrayLen := len(array.Elements)
			newElements := make([]object.Object, arrayLen)
			copy(newElements, array.Elements)
			newElements = append(newElements, args[1:]...)
			return &object.Array{
				Elements: newElements,
			}
		},
	},
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
