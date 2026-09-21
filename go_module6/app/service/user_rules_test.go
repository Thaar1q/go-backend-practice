package service

import (
	"testing"

	"go_module6/app/model"
)

func TestValidateCreate(t *testing.T) {
	req := model.CreateStudentRequest{NIM: "", Name: "mahasiswaAA", Grade: 4.5}
	errs := ValidateCreate(req)

	if _, ok := errs["nim"]; !ok {
		t.Errorf("expected error on nim, got nil")
	}
	if _, ok := errs["grade"]; !ok {
		t.Errorf("expected error on grade out of range, got nil")
	}
}

func TestValidateReplace_MissingGrade(t *testing.T) {
	req := model.ReplaceStudentRequest{NIM: "000000001", Name: "mahasiswaAA", Grade: nil}
	errs := ValidateReplace(req)

	if _, ok := errs["grade"]; !ok {
		t.Errorf("expected error for missing mandatory grade on PUT, got nil")
	}
}

func TestApplyPatch(t *testing.T) {
	current := model.Student{ID: 1, NIM: "000000001", Name: "Lama", Grade: 3.0, IsActive: true}
	newGrade := 3.8
	newStatus := false
	req := model.PatchStudentRequest{Grade: &newGrade, IsActive: &newStatus}

	updated, errs := ApplyPatch(current, req)
	if len(errs) > 0 {
		t.Fatalf("unexpected validation errors: %v", errs)
	}
	if updated.Grade != 3.8 || updated.IsActive != false {
		t.Errorf("patch fields were not applied correctly")
	}
	if updated.Name != "Lama" {
		t.Errorf("unmodified field should retain original value")
	}
}
