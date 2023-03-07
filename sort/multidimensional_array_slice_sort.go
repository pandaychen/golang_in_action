package main

import (
	"fmt"
	"sort"
)

type Matrix [5][3]int

func (m Matrix) Len() int { return len(m) }
func (m Matrix) Less(i, j int) bool {
	for x := range m[i] {
		if m[i][x] == m[j][x] {
			continue
		}
		return m[i][x] < m[j][x]
	}
	return false
}

func (m *Matrix) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

type Matrix2 [][3]int

func (m Matrix2) Len() int { return len(m) }
func (m Matrix2) Less(i, j int) bool {
	for x := range m[i] {
		if m[i][x] == m[j][x] {
			continue
		}
		return m[i][x] < m[j][x]
	}
	return false
}

func (m Matrix2) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

func main() {
	var matrix = [5][3]int{
		{1, 0, 0},
		{1, 1, 1},
		{1, 1, 5},
		{1, 1, 3},
		{1, 1, 4},
	}
	m := Matrix(matrix)
	sort.Sort(&m)
	fmt.Println(m[len(m)-1])

	var matrix2 = [][3]int{
		{1, 0, 0},
		{1, 1, 1},
		{1, 3, 1},
		{1, 1, 5},
		{1, 1, 3},
		{2, 0, 3},
		{2, 1, 3},
		{1, 1, 4},
	}
	m2 := Matrix2(matrix2)
	sort.Sort(&m2)
	fmt.Println(m2[len(m2)-1])
}
