package vm

import (
	"fmt"

	"github.com/the-1aw/monkey-business/code"
	"github.com/the-1aw/monkey-business/compiler"
	"github.com/the-1aw/monkey-business/object"
)

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	// Always point to next value, stack top is stack[stackPointer - 1]
	stackPointer int
}

const StackSize = 2048

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constants:    bytecode.Constants,
		instructions: bytecode.Instructions,
		stack:        make([]object.Object, StackSize),
		stackPointer: 0,
	}
}

func (vm *VM) Run() error {
	for instructionPointer := 0; instructionPointer < len(vm.instructions); instructionPointer++ {
		op := code.Opcode(vm.instructions[instructionPointer])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[instructionPointer+1:])
			instructionPointer += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()
			rightValue := right.(*object.Integer).Value
			leftValue := left.(*object.Integer).Value
			result := leftValue + rightValue
			vm.push(&object.Integer{Value: result})
		case code.OpPop:
			vm.pop()
		}
	}
	return nil
}

func (vm *VM) pop() object.Object {
	obj := vm.stack[vm.stackPointer-1]
	vm.stackPointer--
	return obj
}

func (vm *VM) push(constant object.Object) error {
	if vm.stackPointer >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.stackPointer] = constant
	vm.stackPointer++
	return nil
}

// NOTE: This method is meant for test puposes only.
// It relies on the fact we don't set free stack space to nil when we pop.
func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.stackPointer]
}
