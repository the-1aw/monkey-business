package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/the-1aw/monkey-business/compiler"
	"github.com/the-1aw/monkey-business/evaluator"
	"github.com/the-1aw/monkey-business/lexer"
	"github.com/the-1aw/monkey-business/object"
	"github.com/the-1aw/monkey-business/parser"
	"github.com/the-1aw/monkey-business/vm"
)

const PROMPT = ">> "
const MONKEY_FACE = `	    __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParserErrors(out io.Writer, error []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Whoops! We ran into some mokey business here!\n")
	io.WriteString(out, "Parser errors:\n")
	for _, msg := range error {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func StartInterpreter(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func StartCompiler(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		comp := compiler.New(
			compiler.WithSymbolTable(symbolTable),
			compiler.WithConstants(constants),
		)
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}

		bytecode := comp.Bytecode()
		// NOTE: This is required because compiler can use append on c.constants during compilation
		// which would prevent us from from maintaining a constant context.
		constants = bytecode.Constants

		machine := vm.New(bytecode, vm.WithGlobalsStore(globals))
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executin bytecode failed:\n %s\n", err)
			continue
		}

		lastPopped := machine.LastPoppedStackElem()
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
	}
}
