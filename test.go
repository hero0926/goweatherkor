package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// sk플래닛 개발지원센터 주소
// https://developers.skplanetx.com
/* 우선 가입부터 해서 appkey를 발급받으세요!
 */

//간편날씨
const Simpleurl = "http://apis.skplanetx.com/weather/summary?version=1&lat=%s&lon=%s&stnid=%s"

// 시간별 현재날씨
const Nowurl = "http://apis.skplanetx.com/weather/current/hourly?lon=127.02583&village=&county=korea&lat=%s&lon=%s&city=%s&version=1"

// 자외선지수
const UVurl = " http://apis.skplanetx.com/weather/windex/uvindex?lon=%s&lat=%s&version=1"

// 미세먼지
const Pollutionurl = " http://apis.skplanetx.com/weather/dust?lon=%s&lat=%s&version=1"

// ------------------------------------------------------------------------------------------------------------ //

// Get으로 보낼 리퀘스트에 쓰일 값들
type WeatherRequest struct {
	UserId         string    `json:"x-skpop-userId"`
	AcceptLanguage string    `json:"Accept-Language"`
	Date           time.Time `json:"Date"`
	Accept         string    `json:"Accept"`
	AccessToken    string    `json:access_token`
}

// response에서 받아올 값들
// 1. 간편날씨
type JsonMap map[string]interface{}

type JsonTest struct {
	Weather Summary `json:"weather"`
}

type Summary struct {
	Summary []struct {
		Today            Weather `json:"today"`
		Tomorrow         Weather `json:"tomorrow"`
		DayAfterTomorrow Weather `json:"dayAfterTomorrow`
	} `json:"summary"`
}

type Basic struct {
	Today            []Weather `json:"today"`
	Tomorrow         []Weather `json:"tomorrow"`
	DayAfterTomorrow []Weather `json:"dayAfterTomorrow`
}

type Weather struct {
	Temperature Temperature `json:"temperature"`
	Sky         Sky         `json:"sky"`
}

type Temperature struct {
	Tmax string `json:"tmax"`
	Tmin string `json:"tmin"`
}

type Sky struct {
	Name string `json:"name"`
}

// 2. 시간별 날씨
type HourWeather struct {
	Weather Hourly `json:"weather"`
}

type Hourly struct {
	Hourly []struct {
		Grid        Grid        `json:"grid"`
		Temperature Temperature `json:"temperature"`
		Sky         Sky         `json:"sky"`
	} `json:"hourly"`
}

type Grid struct {
	City    string `json:"city"`
	County  string `json:"county"`
	Village string `json:"village"`
}

// ------------------------------------------------------------------------------------------------------------ //

func main() {

	// 1. api를 호출

	jsonreq, err := json.Marshal(WeatherRequest{

		UserId:         "sh8kim@interpark.com",
		AcceptLanguage: "ko_KR",
		Date:           time.Now(),
		Accept:         "application/json",
		AccessToken:    "",
	})

	if err != nil {
		log.Fatalln(err)
		recover()
	}

	// 1. 서초구 간편날씨
	// 위도, 경도, 관측소 주소 얻는 곳(주의 : 느림)
	// http://minwon.kma.go.kr/main/obvStn.do
	seocho := fmt.Sprintf(Simpleurl, "", "", "401")
	req, err := http.NewRequest("GET", seocho, bytes.NewBuffer(jsonreq))
	req.Header.Set("Content-Type", "application/json")
	// 발급받은 appKey를 넣어주세요.
	req.Header.Set("appkey", "0b37b2f4-0d8b-3cac-b53c-20218fe07af8")
	client := &http.Client{}

	// 2. api에서 값을 받아옴

	resp, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var weather JsonTest
	if err := json.Unmarshal(body, &weather); err != nil {
		log.Fatal(err)
	}

	for _, v := range weather.Weather.Summary {
		fmt.Println(v.Today)            // 오늘 날씨
		fmt.Println(v.Tomorrow)         // 내일 날씨
		fmt.Println(v.DayAfterTomorrow) // 모레 날씨
	}

	// 2. 서초구 시간별 날씨

	seocho = fmt.Sprintf(Nowurl, "37.508214", "127.056541", "seoul")
	req, err = http.NewRequest("GET", seocho, bytes.NewBuffer(jsonreq))
	req.Header.Set("Content-Type", "application/json")
	// 발급받은 appKey를 넣어주세요.
	req.Header.Set("appkey", "0b37b2f4-0d8b-3cac-b53c-20218fe07af8")
	client = &http.Client{}

	// 2. api에서 값을 받아옴

	resp, err = client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)

	var hourweather HourWeather
	if err := json.Unmarshal(body, &hourweather); err != nil {
		log.Fatal(err)
	}

	for _, v := range hourweather.Weather.Hourly {
		fmt.Println(v.Grid)        // 위치
		fmt.Println(v.Sky)         // 날씨
		fmt.Println(v.Temperature) // 온도
	}

	// 3. 서초구 자외선
	seocho = fmt.Sprintf(UVurl, "37.508214", "127.056541")

	// 4. 서초구 미세먼지
	seocho = fmt.Sprintf(Pollutionurl, "37.508214", "127.056541")
}
