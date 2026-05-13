package compiler

import (
	"fmt"
	"sort"

	"github.com/the-1aw/monkey-business/ast"
	"github.com/the-1aw/monkey-business/code"
	"github.com/the-1aw/monkey-business/object"
)

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

type Compiler struct {
	constants []object.Object

	symbolTable *SymbolTable

	scopes   []CompilationScope
	scopeIdx int
}

type CompilerOption func(*Compiler)

// NOTE: I choose to do this here in order to avoid symbolTable dependence on object package
// but maybe this should be there and be called by default by NewSymbolTable
func defineBuiltins(s *SymbolTable) {
	for idx, builtin := range object.Builtins {
		s.DefineBuiltin(idx, builtin.Name)
	}
}

func New(opts ...CompilerOption) *Compiler {
	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	symbolTable := NewSymbolTable()
	defineBuiltins(symbolTable)

	compiler := &Compiler{
		constants:   []object.Object{},
		symbolTable: symbolTable,
		scopes:      []CompilationScope{mainScope},
		scopeIdx:    0,
	}
	for _, opt := range opts {
		opt(compiler)
	}
	return compiler
}

func WithSymbolTable(s *SymbolTable) CompilerOption {
	return func(c *Compiler) {
		c.symbolTable = s
		defineBuiltins(s)
	}
}

func WithConstants(constants []object.Object) CompilerOption {
	return func(c *Compiler) {
		c.constants = constants
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, stmt := range node.Statements {
			err := c.Compile(stmt)
			if err != nil {
				return err
			}
		}
	case *ast.BlockStatement:
		for _, stmt := range node.Statements {
			err := c.Compile(stmt)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)
	case *ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}
		symbol := c.symbolTable.Define(node.Name.Value)
		if symbol.Scope == GlobalScope {
			c.emit(code.OpSetGlobal, symbol.Index)
		} else {
			c.emit(code.OpSetLocal, symbol.Index)
		}
	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}
		c.loadSymbol(symbol)
	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}
		// NOTE: We intentionally use bogus value here.
		// Value will be back-patched once consequence has been compiled
		consequnceJumpPos := c.emit(code.OpJumpNotTruthy, 9999)
		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}
		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}

		// cf note above about back-patching
		alternativeJumpPos := c.emit(code.OpJump, 9999)
		consequenceEnd := len(c.currentInstructions())
		c.changeOperand(consequnceJumpPos, consequenceEnd)

		if node.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}
			if c.lastInstructionIs(code.OpPop) {
				c.removeLastPop()
			}
		}
		alternativeEnd := len(c.currentInstructions())
		c.changeOperand(alternativeJumpPos, alternativeEnd)

	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		}
	case *ast.InfixExpression:
		// NOTE: "<" is an exception in order to share op code with ">"
		// it showcases compiler's code reordering capabilities
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}
			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			c.emit(code.OpGreaterThan)
			return nil
		}
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case "!=":
			c.emit(code.OpNotEqual)
		case "==":
			c.emit(code.OpEqual)
		case ">":
			c.emit(code.OpGreaterThan)
		default:
			return fmt.Errorf("unkown operator %s", node.Operator)
		}
	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			err := c.Compile(el)
			if err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elements))
	case *ast.HashLiteral:
		keys := []ast.Expression{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}
		// NOTE: Compiler and VM could work just fine without key sort
		// but maintaining consistant key order makes the output way easier to test.
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})
		for _, k := range keys {
			err := c.Compile(k)
			if err != nil {
				return err
			}
			err = c.Compile(node.Pairs[k])
			if err != nil {
				return err
			}
		}
		c.emit(code.OpHash, len(node.Pairs)*2)
	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Index)
		if err != nil {
			return err
		}
		c.emit(code.OpIndex)
	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return err
		}
		c.emit(code.OpReturnValue)
	case *ast.FunctionLiteral:
		c.enterScope()

		for _, p := range node.Parameters {
			c.symbolTable.Define(p.Value)
		}

		err := c.Compile(node.Body)
		if err != nil {
			return err
		}
		// Convert the last statement to an implicit return
		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastPopWithReturn()
		}
		// If the last instruction is not a return it means function body was empty or
		// we were not able to turn the last statement into a return value
		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}
		numLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()
		compiledFn := &object.CompiledFunction{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(node.Parameters),
		}
		c.emit(code.OpConstant, c.addConstant(compiledFn))
	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil {
			return err
		}
		for _, arg := range node.Args {
			err := c.Compile(arg)
			if err != nil {
				return err
			}
		}
		c.emit(code.OpCall, len(node.Args))
	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(str))
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	}
	return nil
}

func (c *Compiler) loadSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(code.OpGetGlobal, s.Index)
	case LocalScope:
		c.emit(code.OpGetLocal, s.Index)
	case BuiltinScope:
		c.emit(code.OpGetBuiltin, s.Index)
	}
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
	c.scopes = append(c.scopes, scope)
	c.scopeIdx++
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() code.Instructions {
	instuctions := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIdx--
	c.symbolTable = c.symbolTable.Outer

	return instuctions
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIdx].instructions
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIdx].lastInstruction.Position
	c.replaceInstruction(lastPos, code.Make(code.OpReturnValue))
	c.scopes[c.scopeIdx].lastInstruction.Opcode = code.OpReturnValue
}

// NOTE: For the sake of simplicity we assume old and new instructions are the same size
// Should this assumption change the stategy needs to change and lastInstruction along with
// previousInstruction must be carefuly dealt with
func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	instructions := c.currentInstructions()
	for i := range len(newInstruction) {
		instructions[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)
	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}
	return c.scopes[c.scopeIdx].lastInstruction.Opcode == op
}

// I have mixed feeling about this as it's not idempotent.
// Sure it works for how we use it, but it throws previousInstruction out of sync.
func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIdx].lastInstruction
	previous := c.scopes[c.scopeIdx].previousInstruction

	oldInstructions := c.currentInstructions()
	newInstructions := oldInstructions[:last.Position]

	c.scopes[c.scopeIdx].instructions = newInstructions
	c.scopes[c.scopeIdx].lastInstruction = previous
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	prev := c.scopes[c.scopeIdx].lastInstruction
	last := EmittedInstruction{op, pos}

	c.scopes[c.scopeIdx].previousInstruction = prev
	c.scopes[c.scopeIdx].lastInstruction = last
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	instruction := code.Make(op, operands...)
	pos := c.addInstruction(instruction)
	c.setLastInstruction(op, pos)
	return pos
}

func (c *Compiler) addInstruction(instruction []byte) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), instruction...)
	c.scopes[c.scopeIdx].instructions = updatedInstructions
	return posNewInstruction
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}
