package pkg

// func Edit(repos []internal.GithubRepository) error {
// 	readme, err := os.Create("../README.md")
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		_ = readme.Close()
// 	}()
// 	editREADME(readme, repos)

// 	return nil
// }

// func editREADME(w io.Writer, repos []internal.GithubRepository) {
// 	fmt.Fprint(w, repos)
// 	for _, repo := range repos {
// 		fmt.Fprintf(w, "| [%s](%s) | %s | %d | %d | %d | %d | %d | %v | %v |\n",
// 			repo.FullName,
// 			repo.URL,
// 			repo.Description,
// 			repo.StargazersCount,
// 			repo.SubscribersCount,
// 			repo.ForksCount,
// 			repo.OpenIssuesCount,
// 			repo.CreatedAt.Format("2006-01-02 15:04:05"),
// 			repo.UpdatedAt.Format("2006-01-02 15:04:05"))
// 	}
// 	// fmt.Fprintf(w, tail, flextime.Now().Format(time.RFC3339))
// }

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
