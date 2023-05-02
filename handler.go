package main

import (
	"fmt"
	"log"
)

func handlerAPI(config Config) {
	repo := "beego/beego"
	// repo := "go-gorm/gorm"
	accessToken := config.GithubToken

	res, err := GetRepoStargazers(repo, accessToken, 1)
	// res, err := GetRepoStargazersCount(repo, accessToken)
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Println(res)

	// client := &http.Client{}

	// req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s", repo), nil)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// req.Header.Set("Accept", "application/vnd.github.v3.star+json")
	// req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))

	// resp, err := client.Do(req)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer resp.Body.Close()

	// log.Println("22222222222222")
	// log.Println(resp.Body)
	// log.Println("22222222222222")
	// var body []byte
	// _, err = resp.Body.Read(body)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// log.Println("11111111111111111")
	// log.Println(body)
	// log.Println("11111111111111111")
	// var data map[string]interface{}
	// err = json.Unmarshal(body, &data)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// log.Println("3333333333333")
	// log.Println(data)
	// log.Println("3333333333333")
	// starCount := data["all"].([]interface{})[len(data["all"].([]interface{}))-2].(float64)

	// fmt.Printf("1 month ago star count: %v\n", starCount)
}
