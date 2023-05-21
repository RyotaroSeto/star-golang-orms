package main

func main() {
	Execute()
	// fp, err := os.Open("./bun.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// defer fp.Close()

	// data, err := io.ReadAll(fp)
	// if err != nil {
	// 	panic(err)
	// }

	// starTimes := []string{}
	// tmpStarTimes := strings.Split(string(data), "}")
	// for k, v := range tmpStarTimes {
	// 	if k == 0 || k+1 == len(tmpStarTimes) {
	// 		continue
	// 	}
	// 	v = strings.Replace(v, "{", "", -1)
	// 	result := v[1:]
	// 	starTimes = append(starTimes, result)
	// }

	// var starCount202111 int
	// var starCount202204 int
	// var starCount202301 int
	// var starCount202302 int
	// var starCount202303 int
	// for _, val := range starTimes {
	// 	if val < "2021-11-11 00:00:00 +0000 UTC" {
	// 		starCount202111++
	// 	}
	// 	if val < "2022-04-01 00:00:00 +0000 UTC" {
	// 		starCount202204++
	// 	}
	// 	if val < "2023-01-01 00:00:00 +0000 UTC" {
	// 		starCount202301++
	// 	}
	// 	if val < "2023-02-01 00:00:00 +0000 UTC" {
	// 		starCount202302++
	// 	}
	// 	if val < "2023-03-01 00:00:00 +0000 UTC" {
	// 		starCount202303++
	// 	}
	// }
	// fmt.Println(starCount202111)
	// fmt.Println(starCount202204)
	// fmt.Println(starCount202301)
	// fmt.Println(starCount202302)
	// fmt.Println(starCount202303)
	// stargazer := Stargazer{StarredAt: stringToTime("2017-07-07 02:50:15 +0000 UTC")}
}
