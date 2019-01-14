/*

TODO:
	Обработка всех возможных ошибок!

	сделать функцию сортировки, которая будет принимать другую функцию в качестве параметра (max, min)

	в далекой перспективе научить программу автоматически обрабатывать файлы разных языков, типа chap01.en.md, chap01.de.md, chap01.ru.md и по параметру командной строки "en-ru" генерить на выход таблицу в markdown, чтобы pandoc мог такую таблицу обработать и выгрузить двухколоночную таблицу контракта, к примеру... двуязычные контракты делать

	Прикрутить обработку параметров командной строки. В командной строке можно реализовать следующие команды:
		-r		-- обратная сортировка заголовков
		*.ext	-- расширения файлов, которые обрабатываются
		-l ru-en, ru-de	-- обработка языковых пар, документы должны именоваться file.en.md, file.ru.md
*/
package main

import (
	"bufio"
	"fmt"
	"github.com/alexflint/go-arg"
	"markdownshuffle/internal/node"
	"os"
	"path/filepath"
	"strings"
)

var args struct {
	Reverse   bool
	Extension string `arg:"positional"`
}

func main() {
	args.Extension = "*.md"
	arg.MustParse(&args)
	files, err := filepath.Glob(args.Extension)
	if err != nil {
		panic(err)
	}
	if len(files) < 2 {
		fmt.Printf("Markdown Shuffle -- shuffle all the %s files in current directory to stdout.\nVersion: 0.1 (14.01.2019)\nError: There's no files to merge (less then 2), exiting.", args.Extension)
		fmt.Printf("\nArguments: %t %s", args.Reverse, args.Extension)
		os.Exit(0)
	}

	for i, file := range files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		node.TreeList = append(node.TreeList, new(node.Node))
		node.TreeList[i].Title = strings.ToUpper(f.Name())
		node.TreeList[i].Body = loadFile(f)
		//fmt.Println(">", f.Name(), "[", len(node.TreeList[i].Body), "lines ]")
	}

	//fmt.Println("Length of treelist:", len(node.TreeList))

	for _, n := range node.TreeList {
		node.ParseNode(n) // the parsing routine can be run concurrently
	}

	for len(node.TreeList) > 1 {
		// присоединяем к списку деревьев новый элемент, являющийся объединением пары (0, 1)
		node.TreeList = append(node.TreeList, node.MergeNodes(node.TreeList, node.PairInt{0, 1}))
		// удаляем из списка пару, которую только что объединили
		node.TreeList = node.DeleteNodes(node.TreeList, node.PairInt{0, 1})
		// смысл в том, что мы последовательно сливаем первые два элемента списка, пока в конце не останется только один
	}

	node.PrintTreeSorted(node.TreeList[len(node.TreeList)-1])
}

func loadFile(f *os.File) []string {
	var content []string
	input := bufio.NewScanner(f)
	for input.Scan() {
		content = append(content, input.Text())
	}
	return content
}

func Split(r rune) bool {
	return r == '\n'
}
