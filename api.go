package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const defaultPerPage = 30

func GetRepoStargazers(repo string, token string, page int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/stargazers?per_page=%d", repo, defaultPerPage)

	if page != 0 {
		url = fmt.Sprintf("%s&page=%d", url, page)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Add("Authorization", "token "+token)
	}
	req.Header.Add("Accept", "application/vnd.github.v3.star+json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to get stargazers: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func GetRepoStargazersCount(repo string, token string) (int, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s", repo)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	if token != "" {
		req.Header.Add("Authorization", "token "+token)
	}
	req.Header.Add("Accept", "application/vnd.github.v3.star+json")
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return 0, fmt.Errorf("Failed to get stargazers count: %s", res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	stargazersCount, ok := result["stargazers_count"].(float64)
	if !ok {
		return 0, fmt.Errorf("Failed to parse stargazers count")
	}

	return int(stargazersCount), nil
}

// func GetRepoStarRecords(repo string, token string, maxRequestAmount int) ([]map[string]interface{}, error) {
// 	patchRes, err := GetRepoStargazers(repo, token, 1)
// 	if err != nil {
// 		return nil, err
// 	}

// 	headerLink := patchRes.Header.Get("link")
// 	pageCount := 1
// 	re := regexp.MustCompile(`next.*&page=(\d*).*last`)
// 	regResult := re.FindStringSubmatch(headerLink)

// 	if len(regResult) > 0 {
// 		if regResult[1] != "" {
// 			pageCount, err = strconv.Atoi(regResult[1])
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
// 	}

// 	if pageCount == 1 && len(patchRes.Data) == 0 {
// 		return []map[string]interface{}{}, fmt.Errorf("No data found")
// 	}

// 	requestPages := []int{}
// 	if pageCount < maxRequestAmount {
// 		requestPages = utils.Range(1, pageCount)
// 	} else {
// 		for i := 1; i <= maxRequestAmount; i++ {
// 			requestPages = append(requestPages, int(math.Round(float64((i*pageCount)/maxRequestAmount)))-1)
// 		}
// 		if !utils.IntInSlice(1, requestPages) {
// 			requestPages = append([]int{1}, requestPages...)
// 		}
// 	}

// 	var resArray []*http.Response
// 	for _, page := range requestPages {
// 		res, err := GetRepoStargazers(repo, token, page)
// 		if err != nil {
// 			return nil, err
// 		}
// 		resArray = append(resArray, res)
// 	}

// 	starRecordsMap := make(map[string]int)
// 	if len(requestPages) < maxRequestAmount {
// 		starRecordsData := []map[string]interface{}{}
// 		for _, res := range resArray {
// 			starRecordsData = append(starRecordsData, res.Data...)
// 		}
// 		for i := 0; i < len(starRecordsData); {
// 			starRecordsMap[utils.GetDateString(starRecordsData[i]["starred_at"].(string))] = i + 1
// 			i += int(math.Floor(float64(len(starRecordsData))/float64(maxRequestAmount))) || 1
// 		}
// 	} else {
// 		for i, res := range resArray {
// 			if len(res.Data) > 0 {
// 				starRecord := res.Data[0]
// 				starRecordsMap[utils.GetDateString(starRecord["starred_at"].(string))] = DEFAULT_PER_PAGE * (requestPages[i] - 1)
// 			}
// 		}
// 	}

// 	starAmount, err := GetRepoStargazersCount(repo, token)
// 	if err != nil {
// 		return nil, err
// 	}
// 	starRecordsMap[utils.GetDateString(strconv.FormatInt(time.Now().UnixNano(), 10))] = starAmount

// 	starRecords := make([]StarRecord, 0, len(starRecordsMap))
// 	for date, count := range starRecordsMap {
// 		starRecords = append(starRecords, StarRecord{
// 			Date:  date,
// 			Count: count,
// 		})
// 	}

// 	return starRecords, nil
// }

type GithubUser struct {
	AvatarURL string `json:"avatar_url"`
}

func getRepoLogoUrl(repo string, token string) (string, error) {
	owner := strings.Split(repo, "/")[0]
	url := fmt.Sprintf("https://api.github.com/users/%s", owner)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/vnd.github.v3.star+json")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	var user GithubUser
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return "", err
	}

	return user.AvatarURL, nil
}
