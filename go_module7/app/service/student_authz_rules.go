package service

import (
	"strings"

	"go_module7/app/model"
	"go_module7/helper"
)

func CanAccessUser(
	current model.AuthUser,
	targetID int,
	perms *helper.PermissionSet,
	anyPermission string,
) bool {
	if current.UserID == targetID {
		return true
	}
	return perms.Can(current.Role, anyPermission)
}

func CanAccessStudent(
	current model.AuthUser,
	targetOwnerID int,
	perms *helper.PermissionSet,
	anyPermission string,
) bool {
	if current.UserID == targetOwnerID {
		return true
	}
	return perms.Can(current.Role, anyPermission)
}

func ValidateAssignRole(
	current model.AuthUser,
	targetID int,
	req model.AssignRoleRequest,
	perms *helper.PermissionSet,
) map[string]string {
	errs := map[string]string{}

	role := strings.TrimSpace(req.Role)
	if role == "" {
		errs["role"] = "must be filled"
		return errs
	}

	if !perms.IsKnownRole(role) {
		errs["role"] = "role unknown, pick one from: " +
			strings.Join(perms.KnownRoles(), ", ")
	}

	if current.UserID == targetID {
		errs["role"] = "cannot change own role"
	}

	return errs
}
