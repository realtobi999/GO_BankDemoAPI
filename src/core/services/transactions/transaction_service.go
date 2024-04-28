package transactions

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/realtobi999/GO_BankDemoApi/src/core/domain"
	"github.com/realtobi999/GO_BankDemoApi/src/core/ports"
)

type TransactionService struct {
	TransactionRepository 	ports.ITransactionRepository
	AccountRepository		ports.IAccountRepository
	GeneralRepository		ports.IRepository
}

func NewTransactionService(transactionRepository ports.ITransactionRepository, accountRepository ports.IAccountRepository, generalRepository ports.IRepository) *TransactionService {
	return &TransactionService{
		TransactionRepository: transactionRepository,
		AccountRepository: accountRepository,
		GeneralRepository: generalRepository,
	}
}

func (ts *TransactionService) Index(accountID uuid.UUID, limit int, offset int) ([]domain.Transaction, error) {
	// Declare variables for transactions and error, because
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
			return nil, domain.NotFoundError(errors.New("Transactions not found"))
		}
		return nil, domain.InternalFailure(errors.New("Failed to get transactions: "+err.Error()))
	}

	return transactions, nil	
}

func (ts *TransactionService) Get(transactionID uuid.UUID) (domain.Transaction, error) {
	transaction, err := ts.TransactionRepository.GetTransaction(transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Transaction{}, domain.NotFoundError(errors.New("Transaction not found"))
		}
		return domain.Transaction{}, domain.InternalFailure(errors.New("Failed to get transaction: "+err.Error()))
	}	

	return transaction, nil
}

func (ts *TransactionService) Create(body domain.CreateTransactionRequest) (domain.Transaction, error) {
	transaction := domain.Transaction{
		ID: uuid.New(),
		SenderAccountID: body.SenderAccountID,
		ReceiverAccountID: body.ReceiverAccountID,
		Amount: body.Amount,
		CreatedAt: time.Now(),
	}

	// Get the sender account
	sender, err := ts.AccountRepository.GetAccount(transaction.SenderAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Transaction{}, domain.NotFoundError(errors.New("Account not found"))
		}
		return domain.Transaction{}, domain.InternalFailure(errors.New("Failed to get sender: "+err.Error()))
	}

	// Get the receiver account
	receiver, err := ts.AccountRepository.GetAccount(transaction.ReceiverAccountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Transaction{}, domain.NotFoundError(errors.New("Account not found"))
		}
		return domain.Transaction{}, domain.InternalFailure(errors.New("Failed to get receiver: "+err.Error()))
	}	

	// Create the currency-pair for the transaction
	transaction.CurrencyPair = domain.NewCurrencyPair(sender.Currency, receiver.Currency)

	// Validate the transaction
	if err := transaction.Validate(); err != nil {
		return domain.Transaction{}, domain.ValidationError(err)
	}

	// Validate that the sender can send the money
	if (sender.Balance - transaction.Amount) < 0 {
		return domain.Transaction{}, domain.BadRequestError(errors.New("Sender account doesnt have enough balance"))
	}

	// Calculate the correct amount to add to the receiver account (With the currency conversion)
	receiver.Balance += transaction.CurrencyPair.Calculate(transaction.Amount)
	sender.Balance -= transaction.Amount

	tx, err := ts.GeneralRepository.BeginTransaction()
	if err != nil {
		return domain.Transaction{}, domain.InternalFailure(errors.New("Failed to start a transaction: "+err.Error()))
	}	
	defer tx.Rollback()

	// Update all the accounts
	_, err = ts.AccountRepository.UpdateAccount(sender)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Transaction{}, domain.NotFoundError(errors.New("Account not found"))
		}
		return domain.Transaction{}, domain.InternalFailure(errors.New("Failed to update sender: "+err.Error()))
	}
	
	_, err = ts.AccountRepository.UpdateAccount(receiver)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Transaction{}, domain.NotFoundError(errors.New("Account not found"))
		}
		return domain.Transaction{}, domain.InternalFailure(errors.New("Failed to update receiver: "+err.Error()))
	}

	// Create the transaction
	_, err = ts.TransactionRepository.CreateTransaction(transaction)
	if err != nil {
		return domain.Transaction{}, domain.InternalFailure(errors.New("Failed to create transaction: "+err.Error())) 
	}

	if err := tx.Commit(); err != nil {
		return domain.Transaction{}, domain.InternalFailure(errors.New("Failed to commit: "+err.Error()))
	}

	return transaction, nil
}