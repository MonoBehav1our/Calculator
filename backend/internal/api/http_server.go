package api

import (
	"Calculator/internal/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type echoServer struct {
	e          *echo.Echo
	calculator CalculationService
}

type CalculationService interface {
	GetAllCalculations() ([]model.ExpressionResult, error)
	CreateCalculation(expression string) (model.ExpressionResult, error)
	UpdateCalculation(id int, expression string) (model.ExpressionResult, error)
	DeleteCalculation(id int) error
}

func NewEchoServer(calculator CalculationService) *echoServer {
	server := echoServer{
		e:          echo.New(),
		calculator: calculator,
	}

	return &server
}

func (s *echoServer) StartListening() error {
	s.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))
	s.e.Use(middleware.RequestLogger())
	s.e.Use(middleware.Recover())

	s.e.GET("/calculations", s.handleGetExpressions)
	s.e.POST("/calculations", s.handlePostExpression)
	s.e.PATCH("/calculations/:id", s.handlePatchExpression)
	s.e.DELETE("/calculations/:id", s.handleDeleteExpression)

	if err := s.e.Start(":8080"); err != nil {
		return err
	}

	return nil
}

func (s *echoServer) handleGetExpressions(c *echo.Context) error {
	expResults, err := s.calculator.GetAllCalculations()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, expResults)
}

func (s *echoServer) handlePostExpression(c *echo.Context) error {
	type payload struct {
		Expression string `json:"expression"`
	}
	var p payload
	if err := c.Bind(&p); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	expResult, err := s.calculator.CreateCalculation(p.Expression)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, expResult)
}

func (s *echoServer) handlePatchExpression(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	type payload struct {
		Expression string `json:"expression"`
	}
	var p payload
	if err := c.Bind(&p); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	expResult, err := s.calculator.UpdateCalculation(id, p.Expression)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, expResult)

}

func (s *echoServer) handleDeleteExpression(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = s.calculator.DeleteCalculation(id)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "ok")
}
