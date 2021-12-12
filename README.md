go run vm-translator

Testando a versão 7:
git checkout e6ca99ab50ae3b62d35f8049a7ea6b488588b2bd

Testando a versão 8:
git checkout 813f9e4a8064dee5a4362721e1540352038ae37e

e comentar o codigo write init:

        func (code *CodeWriter) WriteInit() {
            /*
            code.write("@256")
            code.write("D=A")
            code.write("@SP")
            code.write("M=D")
            code.WriteCall("Sys.init", 0)
            */
        }

// todo: checar se o parser está ok
reComments, err := regexp.Compile("//.\*\n")

## Referências

https://eli.thegreenplace.net/2018/go-and-algebraic-data-types/
