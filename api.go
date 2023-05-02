package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultPerPage = 30
const defaultPage = 1

func getRepoStargazers(repo string, token string, page int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/stargazers?per_page=%d", repo, defaultPerPage)
	if page != 0 {
		url = fmt.Sprintf("%s&page=%d", url, page)
	}

	client := NewHttpClient(url, http.MethodGet, token)
	body, err := client.Execute()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func getRepoStargazersCount(repo string, token string) (int, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s", repo)

	client := NewHttpClient(url, http.MethodGet, token)
	body, err := client.Execute()
	if err != nil {
		return 0, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return 0, err
	}

	stargazersCount, ok := result["stargazers_count"].(float64)
	if !ok {
		return 0, fmt.Errorf("failed to parse stargazers count")
	}

	log.Println(result["subscribers_count"])
	log.Println(result["forks_count"])
	log.Println(result["open_issues_count"])
	log.Println(result["updated_at"])
	return int(stargazersCount), nil
}

func getStarsForMonthAgo(repo, token string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/stargazers?per_page=%d&page=%d", repo, defaultPerPage, defaultPage)

	client := NewHttpClient(url, http.MethodGet, token)
	body, err := client.Execute()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println(result[0])
	log.Println(result[1])

	return nil, nil
	// return result, nil
}

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

// type StarRecord struct {
// 	Date  string `json:"date"`
// 	Count int    `json:"count"`
// }

// func GetRepoStarRecords(repo string, token string, maxRequestAmount int) ([]StarRecord, error) {
// 	patchRes, err := getRepoStargazers(repo, token, 0)
// 	if err != nil {
// 		return nil, err
// 	}

// 	headerLink := patchRes.Header.Get("link")
// 	pageCount := 1
// 	if headerLink != "" {
// 		var nextPage, lastPage string
// 		fmt.Sscanf(headerLink, "<%s>; rel=\"next\", <%s>; rel=\"last\"", &nextPage, &lastPage)
// 		fmt.Sscanf(lastPage, "&page=%d", &pageCount)
// 	}

// 	if pageCount == 1 && len(patchRes.Data) == 0 {
// 		return nil, fmt.Errorf("no stargazers found")
// 	}

// 	var requestPages []int
// 	if pageCount < maxRequestAmount {
// 		requestPages = make([]int, pageCount)
// 		for i := range requestPages {
// 			requestPages[i] = i + 1
// 		}
// 	} else {
// 		requestPages = make([]int, maxRequestAmount)
// 		for i := range requestPages {
// 			requestPages[i] = int((float64(i) * float64(pageCount)) / float64(maxRequestAmount))
// 		}
// 		if requestPages[0] != 1 {
// 			requestPages = append([]int{1}, requestPages...)
// 		}
// 	}

// 	resArray := make([]repoStargazersResponse, len(requestPages))
// 	for i, page := range requestPages {
// 		res, err := getRepoStargazers(repo, token, page)
// 		if err != nil {
// 			return nil, err
// 		}
// 		resArray[i] = *res
// 	}

// 	starRecordsMap := make(map[string]int)
// 	if len(requestPages) < maxRequestAmount {
// 		var starRecordsData []struct {
// 			StarredAt string `json:"starred_at"`
// 		}
// 		for _, res := range resArray {
// 			starRecordsData = append(starRecordsData, res.Data...)
// 		}
// 		for i := 0; i < len(starRecordsData); {
// 			starRecordsMap[GetDateString(starRecordsData[i].StarredAt)] = i + 1
// 			i += len(starRecordsData) / maxRequestAmount
// 			if i == len(starRecordsData) {
// 				i--
// 			}
// 		}
// 	} else {
// 		for i, res := range resArray {
// 			if len(res.Data) > 0 {
// 				starRecord := res.Data[0]
// 				starRecordsMap[GetDateString(starRecord.StarredAt)] = defaultPerPage * (requestPages[i] - 1)
// 			}
// 		}
// 	}

// 	starAmount, err := getRepoStargazersCount(repo, token)
// 	if err != nil {
// 		return nil, err
// 	}
// 	starRecordsMap[GetDateString(time.Now().Unix())] = starAmount

// 	starRecords := make([]StarRecord, 0, len(starRecordsMap))
// 	for date, count := range starRecordsMap {
// 		starRecords = append(starRecords, StarRecord{
// 			Date:  date,
// 			Count: count,
// 		})
// 	}

// 	return starRecords, nil
// }

// func GetDateString(t interface{}, format string) string {
// 	var ts int64
// 	switch v := t.(type) {
// 	case int64:
// 		ts = int64(v)
// 	case int:
// 		ts = int64(v)
// 	case string:
// 		parsed, err := strconv.Atoi(v)
// 		if err != nil {
// 			panic(fmt.Sprintf("unable to parse timestamp: %s", v))
// 		}
// 		ts = int64(parsed)
// 	case time.Time:
// 		ts = v.Unix()
// 	default:
// 		panic("unsupported input type")
// 	}

// 	d := time.Unix(ts, 0)
// 	year, month, date := d.Date()
// 	hours, minutes, seconds := d.Clock()

// 	formattedString := format
// 	formattedString = regexp.MustCompile("yyyy").ReplaceAllString(formattedString, strconv.Itoa(year))
// 	formattedString = regexp.MustCompile("MM").ReplaceAllString(formattedString, strconv.Itoa(int(month)))
// 	formattedString = regexp.MustCompile("dd").ReplaceAllString(formattedString, strconv.Itoa(date))
// 	formattedString = regexp.MustCompile("hh").ReplaceAllString(formattedString, strconv.Itoa(hours))
// 	formattedString = regexp.MustCompile("mm").ReplaceAllString(formattedString, strconv.Itoa(minutes))
// 	formattedString = regexp.MustCompile("ss").ReplaceAllString(formattedString, strconv.Itoa(seconds))

// 	return formattedString
// }
