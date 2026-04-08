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
		}
	}
	return nil
}

func (vm *VM) push(constant object.Object) error {
	if vm.stackPointer >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	vm.stack[vm.stackPointer] = constant
	vm.stackPointer++
	return nil
}

func (vm *VM) StackTop() object.Object {
	if vm.stackPointer == 0 {
		return nil
	}
	return vm.stack[vm.stackPointer-1]
}
