package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/the-1aw/monkey-business/lexer"
	"github.com/the-1aw/monkey-business/parser"
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

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
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

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}
