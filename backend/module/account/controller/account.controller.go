package controller

import (
	"errors"

	"github.com/bricksocoolxd/bengi-investment-system/module/account/dto"
	"github.com/bricksocoolxd/bengi-investment-system/module/account/service"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/common"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/middleware"
	"github.com/bricksocoolxd/bengi-investment-system/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type AccountController struct {
	accountService *service.AccountService
}

func NewAccountController(accountService *service.AccountService) *AccountController {
	return &AccountController{
		accountService: accountService,
	}
}

func (ctrl *AccountController) CreateAccount(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}

	var req dto.CreateAccountRequest
	if validatationErrors := utils.ParseAndValidate(c, &req); validatationErrors != nil {
		return common.ValidationError(c, validatationErrors)
	}

	result, err := ctrl.accountService.CreateAccount(c.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrAccountAlreadyExists) {
			return common.BadRequest(c, "Account with this currency already exists")
		}
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "Account created successfully")
}

func (ctrl *AccountController) GetAccounts(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.NotFound(c, "User not found")
	}

	result, err := ctrl.accountService.GetAccounts(c.Context(), userID)
	if err != nil {
		return common.InternalError(c, err.Error())
	}

	return common.Success(c, result, "Accounts retrieved successfully")
}

func (ctrl *AccountController) GetAccountByID(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}
	accountID := c.Params("id")
	result, err := ctrl.accountService.GetAccountByID(c.Context(), accountID, userID)
	if err != nil {
		if errors.Is(err, service.ErrAccountNotFound) {
			return common.NotFound(c, "Account not found")
		}
		return common.InternalError(c, err.Error())
	}
	return common.Success(c, result, "")
}

func (ctrl *AccountController) Deposit(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}
	accountID := c.Params("id")
	var req dto.DepositRequest
	if validationErrors := utils.ParseAndValidate(c, &req); validationErrors != nil {
		return common.ValidationError(c, validationErrors)
	}
	result, err := ctrl.accountService.Deposit(c.Context(), accountID, userID, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAccountNotFound):
			return common.NotFound(c, "Account not found")
		case errors.Is(err, service.ErrAccountFrozen):
			return common.BadRequest(c, "Account is frozen")
		default:
			return common.InternalError(c, err.Error())
		}
	}
	return common.Success(c, result, "Deposit successful")
}

func (ctrl *AccountController) Withdraw(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}
	accountID := c.Params("id")
	var req dto.WithdrawRequest
	if validationErrors := utils.ParseAndValidate(c, &req); validationErrors != nil {
		return common.ValidationError(c, validationErrors)
	}
	result, err := ctrl.accountService.Withdraw(c.Context(), accountID, userID, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAccountNotFound):
			return common.NotFound(c, "Account not found")
		case errors.Is(err, service.ErrAccountFrozen):
			return common.BadRequest(c, "Account is frozen")
		case errors.Is(err, service.ErrInsufficientBalance):
			return common.BadRequest(c, "Insufficient balance")
		default:
			return common.InternalError(c, err.Error())
		}
	}
	return common.Success(c, result, "Withdrawal successful")
}

func (ctrl *AccountController) GetTransactions(c *fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return common.Unauthorized(c, "User not authenticated")
	}
	accountID := c.Params("id")
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)
	result, err := ctrl.accountService.GetTransactions(c.Context(), accountID, userID, limit, offset)
	if err != nil {
		if errors.Is(err, service.ErrAccountNotFound) {
			return common.NotFound(c, "Account not found")
		}
		return common.InternalError(c, err.Error())
	}
	return common.Success(c, result, "")
}
