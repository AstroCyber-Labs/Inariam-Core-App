package gcp

import "github.com/labstack/echo/v4"

type IHandler interface {
	GetUser(echo.Context) error
	ListUsers(echo.Context) error
	UpdateUser(echo.Context) error
	DeleteUser(echo.Context) error

	GetGroup(echo.Context) error
	ListGroups(echo.Context) error
	UpdateGroup(echo.Context) error
	DeleteGroup(echo.Context) error

	ListPolicies(echo.Context) error
	DeletePolicies(echo.Context) error
	UpdatePolicies(echo.Context) error
	GetPolicies(echo.Context) error

	ListServiceAccounts(echo.Context) error
	GetServiceAccount(echo.Context) error
	UpdateServiceAccount(echo.Context) error
	DeleteServiceAccount(echo.Context) error
}
