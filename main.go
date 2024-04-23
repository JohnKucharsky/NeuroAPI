package main

import (
	"fmt"
	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Patient struct {
	Fullname string `json:"fullname"`
	Birthday string `json:"birthday"`
	Gender   int    `json:"gender"`
	Guid     string `json:"guid"`
}

type PatientInput struct {
	Fullname string `json:"fullname" binding:"required"`
	Birthday string `json:"birthday" binding:"required"`
	Gender   int    `json:"gender" binding:"required"`
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		patients, err := patientJSONtoStruct()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"patients": patients,
		})
	})

	r.POST("/", func(c *gin.Context) {
		var input PatientInput
		if err := c.Bind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		g := guid.NewString()
		patient := Patient{
			Fullname: input.Fullname,
			Birthday: input.Birthday,
			Gender:   input.Gender,
			Guid:     g,
		}

		patients, err := patientJSONtoStruct()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = appendToJSON(patient, *patients)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"patient": patient})
	})

	r.PUT("/:id", func(c *gin.Context) {
		var input PatientInput
		if err := c.Bind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id := c.Param("id")
		if ok := guid.IsGuid(id); !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error getting id: " + id})
			return
		}

		g := guid.NewString()
		patient := Patient{
			Fullname: input.Fullname,
			Birthday: input.Birthday,
			Gender:   input.Gender,
			Guid:     g,
		}

		patients, err := patientJSONtoStruct()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = updateJSON(id, patient, *patients)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"patient": patient})
	})

	r.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if ok := guid.IsGuid(id); !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error getting id: " + id})
			return
		}

		patients, err := patientJSONtoStruct()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = removeFromJSON(id, *patients)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"patient": id})
	})

	err := r.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}
