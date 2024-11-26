package analitics

import "fmt"

type Statistic struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Statistic   string `json:"statistic"`
	Prefix      string `json:"prefix"`
	Suffix      string `json:"suffix"`
}

func NewTotalMonthlyInteracions(ammount string) *Statistic {
	return &Statistic{
		Title:       "Total Monthly Interactions",
		Description: "Total number of interactions in the last month",
		Statistic:   ammount,
		Prefix:      "",
		Suffix:      "",
	}

}

func NewMonthlyMostConsultedData(file string) *Statistic {
	return &Statistic{
		Title:       "Monthly Most Consulted Data",
		Description: "Most consulted data in the last month",
		Statistic:   file,
		Prefix:      "",
		Suffix:      "",
	}
}

func NewTotalUsers(ammount string) *Statistic {
	return &Statistic{
		Title:       "Total Users",
		Description: "Total number of users",
		Statistic:   ammount,
		Prefix:      "",
		Suffix:      "Users",
	}
}

func NewInteractionsBetweenDates(amount string, startDate string, endDate string) *Statistic {
	return &Statistic{
		Title:       "Interactions in Date Range",
		Description: fmt.Sprintf("Total interactions between %s and %s", startDate, endDate),
		Statistic:   amount,
		Prefix:      "",
		Suffix:      "interactions",
	}
}

func NewUsersBetweenDates(amount string, startDate string, endDate string) *Statistic {
	return &Statistic{
		Title:       "Users in Date Range",
		Description: fmt.Sprintf("Total new users between %s and %s", startDate, endDate),
		Statistic:   amount,
		Prefix:      "",
		Suffix:      "users",
	}
}

func NewMostConsultedDataBetweenDates(file string, startDate string, endDate string) *Statistic {
	return &Statistic{
		Title:       "Most Consulted Data in Date Range",
		Description: fmt.Sprintf("Most consulted data between %s and %s", startDate, endDate),
		Statistic:   file,
		Prefix:      "",
		Suffix:      "",
	}
}
