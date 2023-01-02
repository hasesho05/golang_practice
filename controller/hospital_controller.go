package controller

import (
	"go-sqlx-gin/db_client"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
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

func GetHospitals(c *gin.Context) {
	var hospitals []Hospital

	var resBody struct {
		PrefName string `json:"prefName"`
	}
	if err := c.ShouldBindJSON(&resBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   true,
			"message": "Invalid request body",
			"content": err.Error(),
		})
		return
	}

	if IsNilValue(resBody) {
		rows, err := db_client.DBClient.Query("SELECT * FROM hospital WHERE prefName = ?;", resBody.PrefName)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": "Invalid request body",
			})
			return
		}
	} else {
		rows, err := db_client.DBClient.Query("SELECT * FROM hospital;")
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": "Invalid request body",
			})
			return
		}
	}

	for rows.Next() {
		var singleHospital Hospital
		if err := rows.Scan(
			&singleHospital.FacilityId,
			&singleHospital.FacilityName,
			&singleHospital.PrefName,
			&singleHospital.FacilityAddr,
			&singleHospital.FacilityTel,
			&singleHospital.SubmitDate,
			&singleHospital.FacilityType,
			&singleHospital.AnsType,
		); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error":   true,
				"message": "Invalid request body",
			})
			return
		}
		hospitals = append(hospitals, singleHospital)
	}

	c.JSON(http.StatusOK, hospitals)
}

func IsNilValue(value interface{}) bool {
	if (value == nil) || reflect.ValueOf(value).IsNil() {
		return true
	} else {
		return false
	}
}
