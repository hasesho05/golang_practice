package fetch

import (
	"encoding/json"
	"fmt"
	"go-sqlx-gin/db_client"
	"io"
	"log"
	"net/http"
)

type Hospital struct {
	FacilityId   string `db:"facilityId"`
	FacilityName string `db:"facilityName"`
	PrefName     string `db:"prefName"`
	FacilityAddr string `db:"facilityAddr"`
	FacilityTel  string `db:"facilityTel"`
	SubmitDate   string `db:"submitDate"`
	FacilityType string `db:"facilityType"`
	AnsType      string `db:"ansType"`
}

func GetHospitalinfo() {
	db_client.InitializeDBConnection()
	url := "https://opendata.corona.go.jp/api/covid19DailySurvey"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	var hospitalList []Hospital
	err = json.Unmarshal(body, &hospitalList)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range hospitalList {
		res, err := db_client.DBClient.Exec("INSERT INTO hospital (facilityId, facilityName, prefName, facilityAddr, facilityTel, submitDate, facilityType, ansType) VALUES (?, ?, ?, ?, ?, ?, ?, ?);",
			v.FacilityId,
			v.FacilityName,
			v.PrefName,
			v.FacilityAddr,
			v.FacilityTel,
			v.SubmitDate,
			v.FacilityType,
			v.AnsType,
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(res)
	}

	// results := make[]hos

	// for i, v in
	// 	fmt.Println(i, v)
	// }

	// res, err := db_client.DBClient.Exec("INSERT INTO hospital (username, password, prefecture) VALUES (?, ?, ?);",

}
