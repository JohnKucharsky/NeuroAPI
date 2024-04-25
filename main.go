package main

import (
	"fmt"
	"github.com/beevik/guid"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
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
	Gender   *int   `json:"gender" binding:"required"`
}

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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
		fmt.Println(input)
		if err := c.Bind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		g := guid.NewString()
		patient := Patient{
			Fullname: input.Fullname,
			Birthday: input.Birthday,
			Gender:   *input.Gender,
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

		patient := Patient{
			Fullname: input.Fullname,
			Birthday: input.Birthday,
			Gender:   *input.Gender,
			Guid:     id,
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

		_, ok := lo.Find(*patients, func(p Patient) bool {
			return p.Guid == id
		})
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can't find patient with id: " + id})
			return
		}

		err = removeFromJSON(id, *patients)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"patientID": id})
	})

	err := r.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
}
