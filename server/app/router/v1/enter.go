package v1

type RouterGroup struct {
	AuthRouter
	CurrencyRouter
	UserRouter
	LendingStrategyRouter
	ScheduleOrderRouter
}
