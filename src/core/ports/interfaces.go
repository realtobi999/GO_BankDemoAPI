package ports

import "github.com/realtobi999/GO_BankDemoApi/src/core/domain"

type ISerializable interface {
	ToDTO() domain.DTO
}
