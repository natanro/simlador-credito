package model

type Simulation struct {
	ID              int     `gorm:"id"`
	RequestedAmount float64 `gorm:"requested_amount"`
	Installments    int     `gorm:"installments"`
	Status          string  `gorm:"status"`
	Age             int     `gorm:"age"`
	AnnualRate      float64 `gorm:"annual_rate"`
	MonthlyRate     float64 `gorm:"monthly_rate"`
	MonthlyPayment  float64 `gorm:"monthly_payment"`
	TotalAmount     float64 `gorm:"total_amount"`
}

func (s *Simulation) TableName() string {
	return "simulacoes"
}