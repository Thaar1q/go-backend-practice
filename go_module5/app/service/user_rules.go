package service

import (
	"strings"

	"go_module5/app/model"
)

func ValidateCreate(req model.CreateStudentRequest) map[string]string {
	errs := map[string]string{}

	if req.NIM == "" {
		errs["nim"] = "is required"
	}
	if req.Name == "" {
		errs["name"] = "is required"
	}
	if req.Grade < 0 || req.Grade > 4.0 {
		errs["grade"] = "must be between 0.0 and 4.0"
	}

	return errs
}

func ValidateReplace(req model.ReplaceStudentRequest) map[string]string {
	errs := map[string]string{}

	if req.NIM == "" {
		errs["nim"] = "is required for PUT"
	}
	if req.Name == "" {
		errs["name"] = "is required for PUT"
	}
	if req.Grade == nil {
		errs["grade"] = "is required for PUT"
	} else if *req.Grade < 0 || *req.Grade > 4.0 {
		errs["grade"] = "must be between 0.0 and 4.0"
	}

	return errs
}

func ApplyPatch(current model.Student, req model.PatchStudentRequest) (model.Student, map[string]string) {
	errs := map[string]string{}
	if req.NIM != nil {
		nimVal := strings.TrimSpace(*req.NIM)
		if nimVal == "" {
			errs["nim"] = "cannot be empty"
		}
		current.NIM = nimVal
	}
	if req.Name != nil {
		nameVal := strings.TrimSpace(*req.Name)
		if nameVal == "" {
			errs["name"] = "cannot be empty"
		}
		current.Name = nameVal
	}
	if req.Grade != nil {
		if *req.Grade < 0 || *req.Grade > 4.0 {
			errs["grade"] = "must be between 0.0 and 4.0"
		}
		current.Grade = *req.Grade
	}
	if req.IsActive != nil {
		current.IsActive = *req.IsActive
	}

	return current, errs
}

func IsEmptyPatch(req model.PatchStudentRequest) bool {
	return req.Name == nil && req.NIM == nil && req.Grade == nil && req.IsActive == nil
}

func CountTotalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}
	return (total + limit - 1) / limit
}
