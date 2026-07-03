package service

import (
	"Calculator/internal/model"
	"fmt"

	"github.com/Knetic/govaluate"
)

type Storage interface {
	GetAllCalculations() ([]model.ExpressionResult, error)
	GetCalculation(id int) (model.ExpressionResult, error)
	CreateCalculation(expression string, result string) (int, error)
	UpdateCalculation(model.ExpressionResult) error
	DeleteCalculation(id int) error
}

type Calculator struct {
	storage Storage
}

func NewCalculator(storage Storage) *Calculator {
	return &Calculator{
		storage: storage,
	}
}

func (c *Calculator) GetAllCalculations() ([]model.ExpressionResult, error) {
	results, err := c.storage.GetAllCalculations()
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (c *Calculator) CreateCalculation(expression string) (model.ExpressionResult, error) {
	result, err := calculate(expression)
	if err != nil {
		return model.ExpressionResult{}, err
	}
	createdId, err := c.storage.CreateCalculation(expression, result)
	if err != nil {
		return model.ExpressionResult{}, err
	}

	expResult := model.ExpressionResult{
		Id:         createdId,
		Expression: expression,
		Result:     result,
	}

	return expResult, nil
}

func (c *Calculator) UpdateCalculation(id int, expression string) (model.ExpressionResult, error) {
	expResult, err := c.storage.GetCalculation(id)
	if err != nil {
		return model.ExpressionResult{}, err
	}

	newResult, err := calculate(expression)
	if err != nil {
		return model.ExpressionResult{}, err
	}

	newExpResult := model.ExpressionResult{
		Id:         expResult.Id,
		Expression: expression,
		Result:     newResult,
	}

	err = c.storage.UpdateCalculation(newExpResult)
	if err != nil {
		return model.ExpressionResult{}, err
	}

	return newExpResult, nil
}

func (c *Calculator) DeleteCalculation(id int) error {
	return c.storage.DeleteCalculation(id)
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
