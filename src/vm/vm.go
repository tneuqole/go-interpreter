package vm

import (
	"code"
	"compiler"
	"fmt"
	"object"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions code.Instructions
	stack        []object.Object

	// always points to next free space
	// top of stack is stack[sp-1]
	sp int
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd:
			obj1 := vm.pop()
			obj2 := vm.pop()

			int1 := obj1.(*object.Integer).Value
			int2 := obj2.(*object.Integer).Value

			vm.push(&object.Integer{Value: int1 + int2})
		case code.OpPop:
			vm.pop()
		}
	}

	return nil
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("STACK OVERFLOW")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	obj := vm.stack[vm.sp-1]
	vm.sp--

	return obj
}