package types

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/utils"
)

type Customer struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Birthday  time.Time
	Email     string
	Phone     string
	State     string
	Address   string
	Accounts  []Account
}

type CustomerDTO struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Birthday  time.Time `json:"birthday"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	State     string    `json:"state"`
	Address   string    `json:"address"`
}

type CreateCustomerRequest struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Birthday  time.Time `json:"birthday"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	State     string    `json:"state"`
	Address   string    `json:"address"`
}

func (c Customer) ToDTO() DTO {
	return CustomerDTO{
		ID:        c.ID.String(),
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Birthday:  c.Birthday,
		Email:     c.Email,
		Phone:     c.Phone,
		State:     c.State,
		Address:   c.Address,
	}
}

func (r CreateCustomerRequest) ToCustomer() Customer {
	return Customer{
		ID:        uuid.New(),
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Birthday:  r.Birthday,
		Email:     r.Email,
		Phone:     r.Phone,
		State:     r.State,
		Address:   r.Address,
		Accounts:  []Account{},
	}
}

func (r CreateCustomerRequest) Validate() []error {
	var validationErrors []error

	if r.FirstName == "" {
		validationErrors = append(validationErrors, errors.New("first name is required"))
	}

	if r.LastName == "" {
		validationErrors = append(validationErrors, errors.New("last name is required"))
	}

	if r.Birthday.IsZero() {
		validationErrors = append(validationErrors, errors.New("birthday is required"))
	} else {
		age := utils.CalculateAge(r.Birthday)
		if age < 18 {
			validationErrors = append(validationErrors, errors.New("age must be at least 18"))
		}
	}

	if r.Email == "" {
		validationErrors = append(validationErrors, errors.New("email is required"))
	} else {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(r.Email) {
			validationErrors = append(validationErrors, errors.New("invalid email format"))
		}
	}

	if r.Phone == "" {
		validationErrors = append(validationErrors, errors.New("phone is required"))
	} else {
		if len(r.Phone) < 12 || len(r.Phone) > 13 {
			validationErrors = append(validationErrors, errors.New("phone number must be between 12 and 13 digits"))
		}

		// Validate the phone number using a regular expression
		re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
		if !re.MatchString(r.Phone) {
			validationErrors = append(validationErrors, errors.New("phone number is not valid"))
		}
	}

	if r.State == "" {
		validationErrors = append(validationErrors, errors.New("state is required"))
	}

	if r.Address == "" {
		validationErrors = append(validationErrors, errors.New("address is required"))
	}

	return validationErrors
}
