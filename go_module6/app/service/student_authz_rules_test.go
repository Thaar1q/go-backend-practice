package service_test

import (
	"testing"

	"go_module6/app/model"
	"go_module6/app/service"
	"go_module6/helper"
)

func TestCanAccessStudent(t *testing.T) {
	perms := helper.NewPermissionSet(map[string][]string{
		"admin":   {"student:read:any", "student:update:any"},
		"student": {},
	})

	callerUser := model.AuthUser{UserID: 1, Role: "student"}
	callerAdmin := model.AuthUser{UserID: 99, Role: "admin"}

	// Own record access
	if !service.CanAccessStudent(callerUser, 1, perms, "student:read:any") {
		t.Errorf("expected true for record owner")
	}

	// Other record without permission (fail closed)
	if service.CanAccessStudent(callerUser, 2, perms, "student:read:any") {
		t.Errorf("expected false for non-owner student")
	}

	// Admin access with :any bypass
	if !service.CanAccessStudent(callerAdmin, 2, perms, "student:read:any") {
		t.Errorf("expected true for admin with read:any permission")
	}
}

func TestValidateAssignRole(t *testing.T) {
	perms := helper.NewPermissionSet(map[string][]string{
		"admin": {"role:assign"},
		"staff": {},
	})

	admin := model.AuthUser{UserID: 10, Role: "admin"}

	// Error: self-demotion
	errs := service.ValidateAssignRole(admin, 10, model.AssignRoleRequest{Role: "staff"}, perms)
	if errs["role"] != "cannot change own role" {
		t.Errorf("expected cannot change own role error, got: %v", errs)
	}

	// Error: unknown role
	errs = service.ValidateAssignRole(admin, 11, model.AssignRoleRequest{Role: "wizard"}, perms)
	if errs["role"] == "" {
		t.Errorf("expected unknown role error")
	}

	// Success
	errs = service.ValidateAssignRole(admin, 11, model.AssignRoleRequest{Role: "staff"}, perms)
	if len(errs) > 0 {
		t.Errorf("expected no errors, got: %v", errs)
	}
}
