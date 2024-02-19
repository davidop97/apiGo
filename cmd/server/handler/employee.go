package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/davidop97/apiGo/internal/domain"
	"github.com/davidop97/apiGo/internal/employee"
	"github.com/gin-gonic/gin"
)

type Employee struct {
	employeeService employee.Service
}

func NewEmployee(employeeService employee.Service) *Employee {
	return &Employee{
		employeeService: employeeService,
	}
}

// Body request parameters
var (
	ID           = "id"
	CardNumberID = "card_number_id"
	FirstName    = "first_name"
	LastName     = "last_name"
	WarehouseID  = "warehouse_id"
)

// EmployeeRequest is a struct that represents a body request for a employee
type EmployeeRequest struct {
	ID           int    `json:"id"`
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  int    `json:"warehouse_id"`
}

// @Summary Get a employee by id or an error if that id not exists.
// @Tags domain.Employee
// @Produce json
// @Success 200
// @Failure 404
// @Param id path int true "id from the employee"
// @Router /employees/{id} [get]
func (e *Employee) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		}
		currentEmployee, err := e.employeeService.GetEmployeeByID(c, id)
		if err != nil {
			switch {
			case errors.Is(err, employee.ErrNotFound):
				c.JSON(http.StatusNotFound, gin.H{"message": "employee not found"})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"data": currentEmployee})
	}
}

// @Summary Get all the employees available or an error if the list is empty.
// @Tags domain.Employee
// @Produce json
// @Success 200
// @Failure 404
// @Router /employees [get]
func (e *Employee) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		employees, err := e.employeeService.GetAllEmployees(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": employees})
	}
}

// @Summary Create a new ∂ employee or return an error if the new employee has invalid
// @Tags domain.Employee
// @Produce json
// @Accept json
// @Success 201
// @Failure 422
// @Failure 409
// @Param body body Employee true "Struct of Employee domain"
// @Router /employees [post]
func (e *Employee) Create() gin.HandlerFunc {
	return func(c *gin.Context) {

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			return
		}

		// - check if empty fields
		var bodyJson map[string]interface{}
		err = json.Unmarshal(body, &bodyJson)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			return
		}
		err = checkEmptyValues(bodyJson, CardNumberID, FirstName, LastName, WarehouseID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}

		// - Create request struct
		req := EmployeeRequest{}
		err = json.Unmarshal(body, &req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			return
		}

		// - Check if negative IDs
		err = checkNegative(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// - Create a new section in the database
		newEmployee := requestToEmployee(req)
		id, err := e.employeeService.SaveEmployee(c, newEmployee)
		if err != nil {
			switch {
			case errors.Is(err, employee.ErrEmployeeAlreadyExists):
				c.JSON(http.StatusConflict, gin.H{"message": "duplicate employee number"})
				return
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			}
		}

		// Response
		newEmployee.ID = id
		c.JSON(http.StatusCreated, gin.H{"data": newEmployee})
	}
}

// @Summary Update employee or return an error if that employee not exist, or has invalid format.
// @Tags domain.Employee
// @Produce json
// @Accept json
// @Success 200
// @Failure 404
// @Param id path int true "id from the employee"
// @Router /employees/{id} [patch]
func (e *Employee) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Request
		// - get ID
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "bad id"})
			return
		}

		// - get original employee
		res, err := e.employeeService.GetEmployeeByID(c, id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "employee not found"})
			return
		}

		// - Create request struct
		req := employeeToRequest(res)
		// - apply changes
		err = c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
			return
		}

		req.ID = id
		// - check if negative IDs
		err = checkNegative(req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// Process
		// - save changes
		employeeUpdated := requestToEmployee(req)
		err = e.employeeService.UpdateEmployee(c, employeeUpdated)
		if err != nil {
			switch {
			case errors.Is(err, employee.ErrEmployeeAlreadyExists):
				c.JSON(http.StatusConflict, gin.H{"message": "duplicate section number"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			}
			return
		}

		// Response
		c.JSON(http.StatusOK, gin.H{"data": employeeUpdated})

		// const (
		// 	statusBadRequest = http.StatusBadRequest
		// 	statusNotFound   = http.StatusNotFound
		// 	statusCreated    = http.StatusCreated
		// 	statusConflict   = http.StatusConflict
		// 	statusOK         = http.StatusOK
		// )

		// idParam := c.Param("id")
		// id, err := strconv.Atoi(idParam)

		// if err != nil {
		// 	c.JSON(statusBadRequest, gin.H{"error": "Invalid employee ID"})
		// 	return
		// }
		// var partialEmployee employee.PartialEmployee
		// if err := c.ShouldBindJSON(&partialEmployee); err != nil {
		// 	c.JSON(statusNotFound, gin.H{"error": "Invalid employee"})
		// }

		// if err := e.employeeService.UpdateEmployee(ctx, id, partialEmployee); err != nil {
		// 	c.JSON(statusConflict, gin.H{"error": err.Error()})
		// 	return
		// }
		// c.Status(statusOK)
	}
}

// @Summary Delete a employee using its id or return an error if that employee not exist.
// @Tags domain.Employee
// @Success 204
// @Failure 404
// @Param id path int true "id from the employee"
// @Router /employees/{id} [delete]
func (e *Employee) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"message": "bad id"})
			return
		}

		// Process
		// - Delete item that matches given id
		err = e.employeeService.DeleteEmployee(c, id)
		// Response
		if err != nil {
			switch {
			case errors.Is(err, employee.ErrNotFound):
				c.JSON(http.StatusNotFound, gin.H{"message": "employee not found"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"message": "internal error"})
			}
			return
		}

		// Response
		c.JSON(http.StatusNoContent, nil)
	}
}

// Checks wether a map contain given keys
func checkEmptyValues(m map[string]interface{}, params ...string) error {
	for _, param := range params {
		_, ok := m[param]
		if !ok {
			err := fmt.Sprintf("field %s is empty", param)
			return errors.New(err)
		}
	}
	return nil
}

// Checks wether negative IDS on struct
func checkNegative(req EmployeeRequest) error {
	if req.WarehouseID <= 0 {
		return fmt.Errorf("negative warehouse_id")
	}
	return nil
}

// requestToEmployee creates a Section struct from a Request struct
func requestToEmployee(request EmployeeRequest) (employee domain.Employee) {
	employee.ID = request.ID
	employee.CardNumberID = request.CardNumberID
	employee.FirstName = request.FirstName
	employee.LastName = request.LastName
	employee.WarehouseID = request.WarehouseID
	return
}

// employeeToRequest creates a Request struct from a Section struct
func employeeToRequest(employee domain.Employee) (request EmployeeRequest) {
	request.ID = employee.ID
	request.CardNumberID = employee.CardNumberID
	request.FirstName = employee.FirstName
	request.LastName = employee.LastName
	request.WarehouseID = employee.WarehouseID
	return
}
