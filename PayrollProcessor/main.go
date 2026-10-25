package main

import (
	"fmt"
	"strconv"
)

type Payable interface {
	fmt.Stringer
	CalculatePay() float64
}

type SalariedEmployee struct {
	Name         string
	AnnualSalary float64
}

func (se SalariedEmployee) CalculatePay() float64 {
	return se.AnnualSalary / 12
}

func (he HourlyEmployee) CalculatePay() float64 {
	return he.HourlyRate * he.HoursWorked
}

func (se SalariedEmployee) String() string {
	return "Salaried Employee: " + se.Name + " Annual Salary: " + strconv.FormatFloat(se.AnnualSalary, 'f', 2, 64)
}

type HourlyEmployee struct {
	Name        string
	HourlyRate  float64
	HoursWorked float64
}

func (he HourlyEmployee) String() string {
	return "Hourly Employee: " + he.Name + " Hourly Rate: " + strconv.FormatFloat(he.HourlyRate, 'f', 2, 64) + " Hours Worked: " + strconv.FormatFloat(he.HoursWorked, 'f', 2, 64)
}

func (ce CommissionEmployee) String() string {
	return "Commissioned Employee: " + ce.Name + " Base Salary: " + strconv.FormatFloat(ce.BaseSalary, 'f', 2, 64) + " Commission Rate: " + strconv.FormatFloat(ce.CommissionRate, 'f', 2, 64) + " Sale Amount: " + strconv.FormatFloat(ce.SaleAmount, 'f', 2, 64)
}

type CommissionEmployee struct {
	Name           string
	BaseSalary     float64
	CommissionRate float64
	SaleAmount     float64
}

func (ce CommissionEmployee) CalculatePay() float64 {
	return ce.BaseSalary + (ce.CommissionRate * ce.SaleAmount)
}

func PrintEmployeeSummary[P fmt.Stringer](e P) {
	fmt.Printf("\n---- Processing Payroll for %s\n", e.String())
}

func ProcessPayroll(employees []Payable) {
	fmt.Println("Processing Payroll for Employees:")
	totalPayroll := 0.0
	for _, employee := range employees {
		PrintEmployeeSummary(employee)
		pay := employee.CalculatePay()
		fmt.Printf("Monthly Pay: $%.2f\n", pay)
		totalPayroll += pay
	}
	fmt.Printf("Total Payroll: $%.2f\n", totalPayroll)
}

func main() {

	fmt.Println("Welcome to the Payroll Processor!")

	salEmp := SalariedEmployee{Name: "Alice Wonderland", AnnualSalary: 72000.00}
	hrEmp := HourlyEmployee{Name: "Bob The Builder", HourlyRate: 25.00, HoursWorked: 160.0}
	comEmp := CommissionEmployee{Name: "Charlie Chaplin", BaseSalary: 2000.00, CommissionRate: 0.10, SaleAmount: 15000.00}

	payrollList := []Payable{
		salEmp,
		hrEmp,
		comEmp,
	}

	ProcessPayroll(payrollList)

}
