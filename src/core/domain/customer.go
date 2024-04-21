package domain

import (
	"regexp"
	"time"

	"github.com/google/uuid"
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
	CreatedAt time.Time
	Token     string
}

type CreateCustomerRequest struct {
	FirstName string   
	LastName  string    
	Birthday  time.Time 
	Email     string    
	Phone     string    
	State     string    
	Address   string    
}

type UpdateCustomerRequest struct {
	FirstName string   
	LastName  string    
	Birthday  time.Time 
	Email     string    
	Phone     string    
	State     string    
	Address   string    
}

/* ------------------------------------------------------------ */
func (r Customer) Validate() *ValidationErrors{
	var errors []string

	if r.ID == uuid.Nil {
		errors = append(errors, "id is required")
	}

	if r.FirstName == "" {
		errors = append(errors, "first name is required")
	}

	if r.LastName == "" {
		errors = append(errors, "last name is required")
	}

	if r.Birthday.IsZero() {
		errors = append(errors, "birthday is required")
	} else {
		age := calculateAge(r.Birthday)
		if age < 18 {
			errors = append(errors, "age must be at least 18")
		}
	}

	if r.Email == "" {
		errors = append(errors, "email is required")
	} else {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(r.Email) {
			errors = append(errors, "invalid email format")
		}
	}

	if r.Phone == "" {
		errors = append(errors, "phone is required")
	} else {
		if len(r.Phone) < 12 || len(r.Phone) > 13 {
			errors = append(errors, "phone number must be between 12 and 13 digits")
		}

		// Validate the phone number using a regular expression
		re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
		if !re.MatchString(r.Phone) {
			errors = append(errors, "phone number is not valid")
		}
	}

	if r.State == "" {
		errors = append(errors, "state is required")
	}

	if r.Address == "" {
		errors = append(errors, "address is required")
	}

	if len(errors) > 0 {
        return &ValidationErrors{Errors: errors}
    }
    return nil
}

func calculateAge(birthday time.Time) int {
	today := time.Now()
	age := today.Year() - birthday.Year()
	if today.Month() < birthday.Month() || (today.Month() == birthday.Month() && today.Day() < birthday.Day()) {
		age--
	}
	return age
}