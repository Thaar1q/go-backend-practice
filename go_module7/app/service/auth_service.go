package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"go_module7/app/model"
	"go_module7/app/repository"
	"go_module7/helper"
)

const refreshTokenBytes = 32

type AuthService struct {
	students    repository.StudentRepository
	tokens      repository.TokenRepository
	jwt         *helper.JWTManager
	permissions *helper.PermissionSet
	refreshTTL  time.Duration
}

func NewAuthService(
	students repository.StudentRepository,
	tokens repository.TokenRepository,
	jwtManager *helper.JWTManager,
	permissions *helper.PermissionSet,
	refreshTTL time.Duration,
) *AuthService {
	return &AuthService{
		students:    students,
		tokens:      tokens,
		jwt:         jwtManager,
		permissions: permissions,
		refreshTTL:  refreshTTL,
	}
}

func (s *AuthService) Register(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()

	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body must be valid JSON")
	}

	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.Validation(errs)
	}

	hashed, err := helper.HashPassword(req.Password)
	if err != nil {
		return helper.Internal(err)
	}

	created, err := s.students.Create(ctx, model.Student{
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		Password: hashed,
		Role:     "student",
		IsActive: true,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Conflict("NIM is already registered")
		}
		return helper.Internal(err)
	}

	return helper.Created(c, "registration successful", created,
		"/api/v1/students/"+strconv.Itoa(created.ID))
}

func (s *AuthService) Login(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()

	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body must be valid JSON")
	}

	// Support both .NIM or .Username
	identifier := strings.TrimSpace(req.NIM)
	if identifier == "" {
		identifier = strings.TrimSpace(req.Username)
	}

	if identifier == "" || req.Password == "" {
		return helper.Validation(map[string]string{
			"nim":      "is required",
			"password": "is required",
		})
	}

	student, err := s.students.FindByUsername(ctx, identifier)
	if err != nil {
		helper.VerifyDummyPassword(req.Password)
		return helper.Unauthorized("invalid username or password")
	}

	if !helper.VerifyPassword(student.Password, req.Password) {
		return helper.Unauthorized("invalid username or password")
	}

	if !student.IsActive {
		return helper.Forbidden("account is deactivated")
	}

	pair, err := s.issueTokenPair(ctx, student)
	if err != nil {
		return helper.Internal(err)
	}

	return helper.Success(c, "login successful", pair)
}

func (s *AuthService) Refresh(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()

	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body must be valid JSON")
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		return helper.BadRequest("refresh_token is required")
	}

	hash := helper.SHA256Hex(req.RefreshToken)

	stored, err := s.tokens.FindActive(ctx, hash)
	if err != nil {
		return helper.Unauthorized(
			"refresh token is invalid or has expired")
	}

	student, err := s.students.FindByID(ctx, stored.UserID)
	if err != nil || !student.IsActive {
		return helper.Unauthorized("account cannot be accessed")
	}

	if err := s.tokens.Revoke(ctx, hash); err != nil {
		return helper.Internal(err)
	}

	pair, err := s.issueTokenPair(ctx, student)
	if err != nil {
		return helper.Internal(err)
	}

	return helper.Success(c, "token successfully refreshed", pair)
}

func (s *AuthService) Logout(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()

	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.BadRequest("body must be valid JSON")
	}

	if strings.TrimSpace(req.RefreshToken) != "" {
		_ = s.tokens.Revoke(ctx, helper.SHA256Hex(req.RefreshToken))
	}

	return helper.Success(c, "logout successful", nil)
}

func (s *AuthService) Me(c *fiber.Ctx) error {
	ctx, cancel := helper.ReqCtx(c)
	defer cancel()

	authUser, ok := helper.CurrentUser(c)
	if !ok {
		return helper.Unauthorized("not authenticated")
	}

	student, err := s.students.FindByID(ctx, authUser.UserID)
	if err != nil {
		return helper.NotFound("student not found")
	}

	return helper.Success(c, "profile retrieved successfully", fiber.Map{
		"student":     student,
		"permissions": s.permissions.PermissionsOf(student.Role),
	})
}

func (s *AuthService) issueTokenPair(ctx context.Context, student model.Student) (model.TokenPair, error) {
	accessToken, err := s.jwt.GenerateAccess(student)
	if err != nil {
		return model.TokenPair{}, err
	}

	refreshToken, err := helper.RandomToken(refreshTokenBytes)
	if err != nil {
		return model.TokenPair{}, err
	}

	err = s.tokens.Save(ctx, model.RefreshToken{
		UserID:    student.ID,
		TokenHash: helper.SHA256Hex(refreshToken),
		ExpiresAt: time.Now().Add(s.refreshTTL),
	})
	if err != nil {
		return model.TokenPair{}, err
	}

	return model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(s.jwt.AccessTTL().Seconds()),
	}, nil
}
