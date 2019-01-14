/*

TODO:

Реализовать интерфейс для Ноды: печать тела (с отступами \n), печать заголовка (аналогично),
печать дерева от этой ноды, печать сортированного дерева от ноды, печать оглавление от ноды,
печать сортированного оглавления от ноды.

Последние два не очень уместны если мы хотим чтобы тулза работала в связке с pandoc, ну да ладно.

Ну и основное — разбор ноды от узла, то есть порубить ноду на кусочки и создать вложенные
ноды...

Остальные функции не являются непосредственно функциями ноды, а обрабатывают пары нод
и списки нод с параметрами, то есть не должны входить в интерфейс собственно самой ноды.

*/

package node

import (
	"bytes"
	"fmt"
	"strings"
)

var TreeList []*Node // список деревьев - корневых узлов - лучше сделать глобальным

type Node struct {
	Title      string
	Body       []string
	NextNode   []*Node
	ParentNode *Node
	Tree       *Node
	Level      int
}

type PairInt struct {
	First, Second int
}

func PrintTree(r *Node) {
	fmt.Println(r.Title)
	for _, line := range r.Body {
		fmt.Println(line)
	}
	for _, n := range r.NextNode {
		PrintTree(n)
	}

}

func HeaderTree(r *Node) {
	if r.Level > 0 {
		fmt.Print("├", strings.Repeat("─", r.Level), " ", Limit(r.Title, 72), "\n")
	} else {
		fmt.Println("\n╒", strings.ToUpper(r.Title))
	}
	for _, n := range r.NextNode {
		HeaderTree(n)
	}

}

func HeaderTreeSorted(r *Node) {
	fmt.Println(r.Title)
	if r.Level > 0 {
		fmt.Print("├", strings.Repeat("─", r.Level), " ", Limit(r.Title, 72), "\n")
	} else {
		fmt.Println("\n╒", strings.ToUpper(r.Title))
	}
	SortNodes(r.NextNode, 0)
	for _, n := range r.NextNode {
		HeaderTreeSorted(n)
	}

}

func PrintTreeSorted(r *Node) {
	fmt.Println(r.Title)
	for _, line := range r.Body {
		fmt.Println(line)
	}
	SortNodes(r.NextNode, 0)
	for _, n := range r.NextNode {
		PrintTreeSorted(n)
	}

}

// this function splits the node to title, body and creates child nodes
// after it does the same to child nodes until Markdown headers are done with
func ParseNode(n *Node) bool {
	head := findHeaders(n.Body, n.Level+1)
	if len(head) > 0 {
		for i, v := range head {
			n.NextNode = append(n.NextNode, new(Node))
			n.NextNode[len(n.NextNode)-1].Level = n.Level + 1
			n.NextNode[len(n.NextNode)-1].Title = n.Body[v]
			if i < len(head)-1 {
				n.NextNode[len(n.NextNode)-1].Body = n.Body[head[i]+1 : head[i+1]]
			} else {
				n.NextNode[len(n.NextNode)-1].Body = n.Body[head[i]+1 : len(n.Body)]
			}
		}
		n.Body = n.Body[:head[0]]
		for _, j := range n.NextNode {
			ParseNode(j)
		}

		return true
	}
	return false
}

// деревья объединяются достаточно просто: сливаем две корневые ноды,
// потом сливаем списки и смотрим, есть ли в полученном списке ноды, которые надо
// объединить. если есть -- выполняем объединение и проводим ту же операцию над
// списком (рекурсия). смысл в том, что если в объединяемых списках нет совпадающих
// позиций, то и проблем нет, это просто разные разделы итогового документа
func MergeNodes(l []*Node, p PairInt) (result *Node) {
	result = new(Node)
	result.Level = l[p.First].Level
	result.Title = l[p.First].Title // + " " + l[p.Second].Title
	result.Body = append(result.Body, l[p.First].Body...)
	result.Body = append(result.Body, l[p.Second].Body...)
	result.NextNode = append(result.NextNode, l[p.First].NextNode...)
	result.NextNode = append(result.NextNode, l[p.Second].NextNode...)

	dupes := findDup(result.NextNode, 0)
	for _, t := range dupes {
		merged := MergeNodes(result.NextNode, t)
		result.NextNode = append(result.NextNode, merged)
		result.NextNode = DeleteNodes(result.NextNode, t)
	}
	return
}

// this function trims a string to a number of symbols
func Limit(s string, lim int) string {
	r := bytes.Runes([]byte(s))
	l := len(r)
	if l > lim {
		return string(r[:lim]) + "..."
	} else {
		return s
	}
}

func swapNodes(l []*Node, n, m int) {
	l[n], l[m] = l[m], l[n]
}

// this function returns a list of int pairs, which represent
// duplicate nodes in a list
func findDup(l []*Node, lb int) (dup []PairInt) {
	dup = []PairInt{}
	if len(l) > 1 {
		start := lb
		for i := lb + 1; i < len(l); i++ {
			if l[i].Title == l[start].Title {
				dup = append(dup, PairInt{First: start, Second: i})
				//fmt.Println(start, l[start].Title, i, l[i].Title)
			}
		}
	}

	return
}

// this function returns an index of a minimum element
// in slice L from left boundary to the right limit
func min(l []*Node, lb int) int {
	if len(l) > 0 {
		min := lb
		for i := lb; i < len(l); i++ {
			if l[i].Title < l[min].Title {
				min = i
			}

		}
		return min
	} else {
		return -1
	}
}

// this function sort a node list from left boundary (LB) limit.
// if length of the list is 0 or 1, there's nothing to sort.
// if not, see if the minimum function to the right of LB returns
// an index number, larger than left boundary. if it does, there's
// something to do: swap this minimum element with LB element, so the
// minimum is on the left of the processed range.
// then the sorting is repeated on the smaller range recursively.
func SortNodes(l []*Node, lb int) {
	if len(l) > 1 {
		swapNodes(l, lb, min(l, lb))
		if lb+1 < len(l) {
			SortNodes(l, lb+1)
		}
	}
}

// this is a new function made to accept a sorting function as
// a parameter (for ASC or DSC order) TODO
func Sort(l []*Node, sorter func(m []*Node, lb int)) {}

func findHeaders(content []string, l int) []int {
	var idx []int
	if l == 0 {
		return idx
	} else {
		var level string = strings.Repeat("#", l) + " "
		for i, _ := range content {
			if strings.HasPrefix(content[i], level) {
				idx = append(idx, i)
			}
		}
		return idx
	}
}

func print(l []*Node) {
	fmt.Print("[")
	for _, n := range l {
		fmt.Print(" ", n.Title)
	}
	fmt.Print(" ]")
}

func DeleteNodes(l []*Node, p PairInt) (result []*Node) {
	swapNodes(l, p.First, len(l)-1)
	swapNodes(l, p.Second, len(l)-2)
	result = l[:len(l)-2]
	//print(result)
	return
}
