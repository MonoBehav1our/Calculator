package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/Knetic/govaluate"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	}))
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	server := Server{
		Storage: Storage{
			ExpressionResults: make(map[int]ExpressionResult),
		},
	}

	e.GET("/calculations", server.handleGetExpressions)
	e.POST("/calculations", server.handlePostExpression)
	e.PATCH("/calculations/:id", server.handlePatchExpression)
	e.DELETE("/calculations/:id", server.handleDeleteExpression)

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start Server", "error", err)
	}
}

func (s *Server) handleGetExpressions(c *echo.Context) error {
	results := make([]ExpressionResult, 0, len(s.Storage.ExpressionResults))
	for _, expression := range s.Storage.ExpressionResults {
		results = append(results, expression)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Id < results[j].Id
	})

	return c.JSON(http.StatusOK, results)
}

func (s *Server) handlePostExpression(c *echo.Context) error {
	type payload struct {
		Expression string `json:"expression"`
	}
	var p payload
	if err := c.Bind(&p); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	result, err := calculate(p.Expression)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	expResult := ExpressionResult{
		Id:         s.nextId(),
		Expression: p.Expression,
		Result:     result,
	}

	s.Storage.ExpressionResults[expResult.Id] = expResult

	return c.JSON(http.StatusOK, expResult)
}

func (s *Server) handlePatchExpression(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	expResult, ok := s.Storage.ExpressionResults[id]
	if !ok {
		return c.String(http.StatusNotFound, "expression not found")
	}

	type payload struct {
		Expression string `json:"expression"`
	}
	var p payload
	if err := c.Bind(&p); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	newResult, err := calculate(p.Expression)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	expResult.Expression = p.Expression
	expResult.Result = newResult

	s.Storage.ExpressionResults[id] = expResult

	return c.JSON(http.StatusOK, expResult)

}

func (s *Server) handleDeleteExpression(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	_, ok := s.Storage.ExpressionResults[id]
	if !ok {
		return c.String(http.StatusNotFound, "expression not found")
	}

	delete(s.Storage.ExpressionResults, id)
	return c.String(http.StatusOK, "ok")
}

func calculate(exp string) (string, error) {
	expression, err := govaluate.NewEvaluableExpression(exp)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга выражения: %w", err)
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		return "", fmt.Errorf("ошибка математического вычисления: %w", err)
	}

	return fmt.Sprintf("%v", result), nil
}

type ExpressionResult struct {
	Id         int    `json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

type Storage struct {
	ExpressionResults map[int]ExpressionResult
}

type Server struct {
	Storage   Storage
	currentId int
}

func (s *Server) nextId() int {
	s.currentId++
	return s.currentId
}
