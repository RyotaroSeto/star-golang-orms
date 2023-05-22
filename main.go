package main

import (
	"star-golang-orms/cmd"
	"time"
)

var timeTime time.Time = time.Date(2022, 4, 1, 0, 0, 0, 0, time.UTC) //2022-04-01 00:00:00 +0000 UTC

func main() {
	// fmt.Println(timeTime)
	// fmt.Println(time.Now().UTC())
	// fmt.Println(time.Now().UTC().AddDate(0, -1, 0))
	// fmt.Println(time.Now().UTC().AddDate(0, -34, 0))
	cmd.Execute()

}

// type CheckMouth struct {
// 	StarCount12MouthAgo int
// 	StarCount9MouthAgo  int
// 	StarCount6MouthAgo  int
// 	StarCount3MouthAgo  int
// 	StarCountNow        int
// }

// func main() {
// 	checkMouth := CheckMouth{
// 		StarCount12MouthAgo: 100,
// 		StarCount9MouthAgo:  200,
// 		StarCount6MouthAgo:  300,
// 		StarCount3MouthAgo:  400,
// 		StarCountNow:        500,
// 	}

// 	now := time.Now()
// 	y := fmt.Sprintf("%04d", now.Year())
// 	m := fmt.Sprintf("%02d", now.Month())
// 	d := fmt.Sprintf("%02d", now.Day())
// 	fmt.Println(y)
// 	fmt.Println(m)
// 	fmt.Println(d)

// 	nowStr := y + m + d

// 	layout := "20060102"
// 	date, err := time.Parse("20060102", nowStr)
// 	if err != nil {
// 		fmt.Println("無効な日付形式です。")
// 		return
// 	}
// 	tmp3MonthsAgo := date.AddDate(0, -3, 0)
// 	date3MonthsAgo := tmp3MonthsAgo.Format(layout)
// 	tmp6MonthsAgo := date.AddDate(0, -6, 0)
// 	date6MonthsAgo := tmp6MonthsAgo.Format(layout)
// 	tmp9MonthsAgo := date.AddDate(0, -9, 0)
// 	date9MonthsAgo := tmp9MonthsAgo.Format(layout)
// 	tmp12MonthsAgo := date.AddDate(0, -12, 0)
// 	date12MonthsAgo := tmp12MonthsAgo.Format(layout)

// 	header := []string{
// 		date12MonthsAgo,
// 		date9MonthsAgo,
// 		date6MonthsAgo,
// 		date3MonthsAgo,
// 		nowStr,
// 	}

// 	separator := strings.Repeat("-", 8)

// 	fmt.Println(formatRow(header))
// 	fmt.Println(separator)

// 	row := []string{
// 		fmt.Sprintf("%d", checkMouth.StarCount12MouthAgo),
// 		fmt.Sprintf("%d", checkMouth.StarCount9MouthAgo),
// 		fmt.Sprintf("%d", checkMouth.StarCount6MouthAgo),
// 		fmt.Sprintf("%d", checkMouth.StarCount3MouthAgo),
// 		fmt.Sprintf("%d", checkMouth.StarCountNow),
// 	}

// 	fmt.Println(formatRow(row))
// }

// func formatRow(row []string) string {
// 	var formattedRow strings.Builder

// 	for _, cell := range row {
// 		formattedRow.WriteString("| ")
// 		formattedRow.WriteString(cell)
// 		formattedRow.WriteString(" ")
// 	}

// 	formattedRow.WriteString("|")

// 	return formattedRow.String()
// }
