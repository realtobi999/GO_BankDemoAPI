package customer

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
	"github.com/realtobi999/GO_BankDemoApi/src/core/ports"
)

const TOKEN_LENGTH = 64

type CustomerService struct {
	CustomerRepository ports.ICustomerRepository
}

func NewCustomerService(customerRepository ports.ICustomerRepository) *CustomerService {
	return &CustomerService{
		CustomerRepository: customerRepository,
	}
}

func (cs *CustomerService) Index(limit, offset int) ([]domain.Customer, error) {
	customers, err := cs.CustomerRepository.GetAllCustomers(limit, offset)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NotFoundError(errors.New("Customers not found"))
		}
		return nil, domain.InternalFailure(err)
	}

	return customers, nil
}

func (cs *CustomerService) Get(customerID uuid.UUID) (domain.Customer, error) {
	customer, err := cs.CustomerRepository.GetCustomer(customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Customer{}, domain.NotFoundError(errors.New("Customer not found"))
		}
		return domain.Customer{}, domain.InternalFailure(err)
	}

	return customer, nil
}

func (cs *CustomerService) Create(body domain.CreateCustomerRequest) (domain.Customer, error) {
	customer := domain.Customer{
		ID:        uuid.New(),
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Birthday:  body.Birthday,
		Email:     body.Email,
		Phone:     body.Phone,
		State:     body.State,
		Address:   body.Address,
		CreatedAt: time.Now(),
		Token:     GenerateToken(),
	}

	if err := customer.Validate(); err != nil {
		return domain.Customer{}, domain.ValidationError(err)
	}

	_, err := cs.CustomerRepository.CreateCustomer(customer)
	if err != nil {
		return domain.Customer{}, domain.InternalFailure(err)
	}

	return customer, nil
}

func (cs *CustomerService) Update(customerID uuid.UUID, body domain.UpdateCustomerRequest) (int64, error) {
	customer := domain.Customer{
		ID:        customerID,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Birthday:  body.Birthday,
		Email:     body.Email,
		Phone:     body.Phone,
		State:     body.State,
		Address:   body.Address,
	}

	if err := customer.Validate(); err != nil {
		return 0, domain.ValidationError(err)
	}

	affectedRows, err := cs.CustomerRepository.UpdateCustomer(customer)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.NotFoundError(errors.New("Customer not found"))
		}
		return 0, domain.InternalFailure(err)
	}

	if affectedRows == 0 {
		return 0, domain.InternalFailure(errors.New("No rows affected"))
	}

	return affectedRows, nil
}

func (cs *CustomerService) Delete(customerID uuid.UUID) (int64, error) {
	affectedRows, err := cs.CustomerRepository.DeleteCustomer(customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, domain.NotFoundError(errors.New("Account not found"))
		}
		return 0, domain.InternalFailure(err)
	}

	if affectedRows == 0 {
		return 0, domain.InternalFailure(errors.New("No rows affected"))
	}

	return affectedRows, nil
}

func (cs *CustomerService) Auth(customerID uuid.UUID, token string) (bool, error) {
	authenticated, err := cs.CustomerRepository.AuthCustomer(customerID, token)
	if err != nil {
		return false, err
	}

	return authenticated, nil
}
