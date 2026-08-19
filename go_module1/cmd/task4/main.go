package main

import "fmt"

type Student struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Grade    float64 `json:"grade"`
	IsActive bool    `json:"is_active"`
}

// Value receiver
func (s Student) GetInfo() string {
	return fmt.Sprintf("%v - %s - %v - Aktif: %v", s.ID, s.Name, s.Grade, s.IsActive)
}

// Pointer Receiver
func (s *Student) UpdateGrade(grade float64) { s.Grade = grade }
func (s *Student) Activate()                 { s.IsActive = true }
func (s *Student) Deactivate()               { s.IsActive = false }

func main() {
	mahasiswa := Student{ID: 13, Name: "MahasiswaAA", Grade: 4.01, IsActive: false}

	fmt.Println("4 - Struct\n---")
	fmt.Println("Kondisi Awal			= ", mahasiswa.GetInfo())
	mahasiswa.UpdateGrade(12.3)
	fmt.Println("Kondisi Update Grade		= ", mahasiswa.GetInfo())
	mahasiswa.Activate()
	fmt.Println("Kondisi Status Aktif		= ", mahasiswa.GetInfo())
	mahasiswa.Deactivate()
	fmt.Println("Kondisi Status Non-aktif	= ", mahasiswa.GetInfo())
}
