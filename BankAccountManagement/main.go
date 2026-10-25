package main

import (
	"errors"
	"fmt"
)

type Account struct {
	AccountNumber string
	Balance       float64
	OwnerName     string
}

type SavingsAccount struct {
	Account
	InterestRate float64
}

type OverdraftAccount struct {
	Account
	OverdraftLimit float64
}

func (acc *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	acc.Balance += amount
	fmt.Println("Deposit successful")
	return nil
}

func (acc *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	if acc.Balance < amount {
		return errors.New("insufficient funds")
	}
	acc.Balance -= amount
	fmt.Println("Withdrawal successful")
	return nil
}

func (acc *Account) GetBalance() float64 {
	return acc.Balance
}

func (acc *Account) String() string {
	return fmt.Sprintf("Account Number: %s, Owner Name: %s, Balance: %.2f", acc.AccountNumber, acc.OwnerName, acc.Balance)
}

func (sa *SavingsAccount) AddInterest() {
	interest := sa.Balance * sa.InterestRate
	fmt.Println("Adding Interest $ ", interest, " at rate ", sa.InterestRate, "% to Savings Account")
	err := sa.Deposit(interest)
	if err != nil {
		fmt.Println("Error depositing interest: ", err)
	} else {
		fmt.Println("Interest deposited successfully")
	}
}

func (oa *OverdraftAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("withdrawal amount must be positive")
	}
	if (oa.Balance + oa.OverdraftLimit) < amount {
		return fmt.Errorf("withdrawal of $%.2f exceeds overdraft limit for %s. Available including overdraft: $%.2f",
			amount, oa.AccountNumber, oa.Balance+oa.OverdraftLimit)
	}
	oa.Balance -= amount
	fmt.Printf("Withdrew $%.2f from overdraft account %s. New Balance: $%.2f\n", amount, oa.AccountNumber, oa.Balance)
	return nil
}

func main() {

	fmt.Println("--- Bank Account System ---")

	savAcc := SavingsAccount{
		Account: Account{ // Initialize the embedded Account
			AccountNumber: "SAV001",
			Balance:       1000.00,
			OwnerName:     "Alice Saver",
		},
		InterestRate: 0.02, // 2%
	}
	fmt.Println("\n--- Savings Account Operations ---")
	fmt.Println(savAcc.Account.String())

	err := savAcc.Deposit(200.00)
	if err != nil {
		fmt.Printf("Error depositing $%.2f to savings account. %v\n", 200.00, err)
	}
	savAcc.AddInterest()
	err = savAcc.Withdraw(50.00)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println("Final Savings Details:", savAcc.Account.String())

	ovdAcc := OverdraftAccount{
		Account: Account{
			AccountNumber: "OVD002",
			Balance:       100.00,
			OwnerName:     "Bob Spender",
		},
		OverdraftLimit: 200.00,
	}

	fmt.Println("\n--- Overdraft Account Operations ---")
	fmt.Println(ovdAcc.Account.String())

	err = ovdAcc.Deposit(50.00)
	if err != nil {
		fmt.Println("Error:", err)
	}
	err = ovdAcc.Withdraw(200.00)
	if err != nil {
		fmt.Println("Error:", err)
	}

	err = ovdAcc.Withdraw(100.00)
	if err != nil {
		fmt.Println("Error:", err)
	}
	err = ovdAcc.Withdraw(100.00)
	if err != nil {
		fmt.Println("Error (expected for overdraft limit):", err)
	}

	fmt.Println("Final Overdraft Details:", ovdAcc.Account.String())
}
