package transactions

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
	"github.com/realtobi999/GO_BankDemoApi/src/core/ports"
)

type TransactionService struct {
	TransactionRepository ports.ITransactionRepository
}

func NewTransactionService(transactionRepository ports.ITransactionRepository) *TransactionService {
	return &TransactionService{
		TransactionRepository: transactionRepository,
	}
}

func (ts *TransactionService) Index(accountID uuid.UUID, limit int, offset int) ([]domain.Transaction, error) {
	// Declare variables for transactions, because
	// we can then access them in if/else scope
	var transactions []domain.Transaction
	var err error

	// If the id is provided by the handler
	// fetch the transactions filtered by the account
	if accountID != uuid.Nil {
		transactions, err = ts.TransactionRepository.GetAllTransactionsFromAccount(accountID, limit, offset)
	} else {
		transactions, err = ts.TransactionRepository.GetAllTransactions(limit, offset)
	}

	// Handle error for both options
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.NotFoundError("Transactions not found")
		}
		return nil, domain.InternalFailure(err)
	}

	return transactions, nil	
}

func (ts *TransactionService) Get(transactionID uuid.UUID) (domain.Transaction, error) {
	transaction, err := ts.TransactionRepository.GetTransaction(transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Transaction{}, domain.NotFoundError("Transaction not found")
		}
		return domain.Transaction{}, domain.InternalFailure(err)
	}	

	return transaction, nil
}