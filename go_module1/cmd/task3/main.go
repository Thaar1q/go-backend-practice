package main

import "fmt"

func swapValue(a, b int) {
	a, b = b, a
}

func swapPointer(a, b *int) {
	*a, *b = *b, *a
}

func updateSliceValue(s []string, newItem string) {
	s = append(s, newItem)
}

func updateSlicePointer(s *[]string, newItem string) {
	*s = append(*s, newItem)
}

func main() {
	angka1, angka2 := 111, 222
	fmt.Println("3A - Swap Integer")
	fmt.Println("---")
	fmt.Println("Kondisi awal		=", angka1, angka2)
	swapValue(angka1, angka2)
	fmt.Println("Kondisi swapValue	=", angka1, angka2)
	swapPointer(&angka1, &angka2)
	fmt.Println("Kondisi swapPointer 	=", angka1, angka2)

	var stringTest []string = []string{"AaAa", "BbBb"}
	fmt.Println("\n3B - Update Slice")
	fmt.Println("---")
	fmt.Println("Kondisi awal		=", stringTest)
	updateSliceValue(stringTest, "CcCc")
	fmt.Println("Kondisi updateValue	=", stringTest)
	updateSlicePointer(&stringTest, "CcCc")
	fmt.Println("Kondisi updatePointer	=", stringTest)
}
