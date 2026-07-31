package dto

type TransactionStatistics struct {
	AccountNo            string `json:"account_no"`
	TodayTransaction     int64  `json:"today_transaction"`
	ThisWeekTransaction  int64  `json:"this_week_transaction"`
	ThisMonthTransaction int64  `json:"this_month_transaction"`
	TotalTransaction     int64  `json:"total_transaction"`
}
