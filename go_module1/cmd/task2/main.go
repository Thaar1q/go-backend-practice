package main

import "fmt"

func deklarasiVariabel() {
	fmt.Println("2A - Deklarasi Variabel")
	fmt.Println("---")

	var stringTest string = "Testing"
	fmt.Println("String		=", stringTest)

	var intTest int = 1234
	fmt.Println("Int		=", intTest)

	var floatTest float64 = 1.234
	fmt.Println("Float		=", floatTest)

	var boolTest bool
	fmt.Println("Bool		=", boolTest)

	var sliceTest []int
	sliceTest = append(sliceTest, 1, 3, 5, 7, 9)
	fmt.Println("Slice		=", sliceTest)
}

func mapMahasiswa() {
	fmt.Println("\n2B - Deklarasi Map Nilai Mahasiswa")
	fmt.Println("---")
	nilaiMahasiswa := map[string]int{
		"MahasiswaAA": 79,
		"MahasiswaBB": 86,
		"MahasiswaCC": 69,
	}

	// Menelusuri isi setelah perubahan
	fmt.Println("Kondisi awal:")
	for key, value := range nilaiMahasiswa {
		fmt.Printf("%s: %v\n", key, value)
	}

	// Tambah entry baru
	nilaiMahasiswa["MahasiswaDD"] = 77

	// Cek nilai
	nilai, exists := nilaiMahasiswa["MahasiswaAA"]
	if exists {
		fmt.Println("\nCek Nilai MahasiswaAA	=", nilai)
	}

	// Delete nilai
	delete(nilaiMahasiswa, "MahasiswaAA")

	// Menelusuri isi setelah perubahan
	fmt.Println("\nKondisi akhir:")
	for key, value := range nilaiMahasiswa {
		fmt.Printf("%s: %v\n", key, value)
	}
}

func main() {
	deklarasiVariabel()
	mapMahasiswa()
}
