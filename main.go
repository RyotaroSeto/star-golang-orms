package main

import (
	"fmt"
	"log"
)

func main() {
	config, err := LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	repo := "gorm/"
	accessToken := config.GithubToken

	res, err := GetRepoStargazersCount(repo, accessToken)
	// res, err := api.GetRepoStargazers(repo, accessToken, 1)
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

// func main() {
// 	repoURL := "https://api.github.com/repos/go-gorm/gorm"

// 	req, err := http.NewRequest("GET", repoURL, nil)
// 	if err != nil {
// 		fmt.Println("Error creating HTTP request:", err)
// 		return
// 	}

// 	client := http.Client{}
// 	res, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error sending HTTP request:", err)
// 		return
// 	}

// 	body, err := io.ReadAll(res.Body)
// 	defer res.Body.Close()
// 	if err != nil {
// 		fmt.Println("Error reading response body:", err)
// 		return
// 	}

// 	var result map[string]interface{}
// 	err = json.Unmarshal(body, &result)
// 	if err != nil {
// 		fmt.Println("Error parsing response body:", err)
// 		return
// 	}

// 	stars := result["stargazers_count"].(float64)
// 	fmt.Println("Stars:", stars)
// }

// func main() {
// 	owner := "owner_name"
// 	repo := "repo_name"
// 	accessToken := "your_access_token"

// 	// create http client
// 	client := &http.Client{}

// 	// create request
// 	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/stats/participation", owner, repo), nil)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	repoURL := `https://api.github.com/repos/beego/stargazers?per_page=1`

// 	req, err := http.NewRequest("GET", repoURL, nil)
// 	if err != nil {
// 		fmt.Println("Error creating HTTP request:", err)
// 		return
// 	}
// 	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))

// 	client := http.Client{}
// 	res, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error sending HTTP request:", err)
// 		return
// 	}

// 	body, err := io.ReadAll(res.Body)
// 	defer res.Body.Close()
// 	if err != nil {
// 		fmt.Println("Error reading response body:", err)
// 		return
// 	}

// 	var result map[string]interface{}
// 	err = json.Unmarshal(body, &result)
// 	if err != nil {
// 		fmt.Println("Error parsing response body:", err)
// 		return
// 	}

// 	log.Println(result)
// 	// starCount := result["all"].([]interface{})[len(result["all"].([]interface{}))-2].(float64)
// 	// fmt.Printf("1 month ago star count: %v\n", starCount)
// }

// func main() {
// 	repoURL := "https://api.github.com/repos/go-gorm/gorm"

// 	req, err := http.NewRequest("GET", repoURL, nil)
// 	if err != nil {
// 		fmt.Println("Error creating HTTP request:", err)
// 		return
// 	}

// 	client := http.Client{}
// 	res, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("Error sending HTTP request:", err)
// 		return
// 	}

// 	body, err := io.ReadAll(res.Body)
// 	defer res.Body.Close()
// 	if err != nil {
// 		fmt.Println("Error reading response body:", err)
// 		return
// 	}

// 	var result map[string]interface{}
// 	err = json.Unmarshal(body, &result)
// 	if err != nil {
// 		fmt.Println("Error parsing response body:", err)
// 		return
// 	}
// 	log.Println(result)
// 	stars := result["stargazers_count"].(float64)
// 	fmt.Println("Stars:", stars)
// }
