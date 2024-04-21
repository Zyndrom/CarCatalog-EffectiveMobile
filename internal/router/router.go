package router

import (
	"CarsCatalog/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"

	_ "CarsCatalog/docs"

	ginSwagger "github.com/swaggo/gin-swagger"
)

type router struct {
	ginRouter  *gin.Engine
	carService carService
}

type carService interface {
	GetCarsWithOwners(query map[string][]string) ([]model.Car, error)
	DeleteCarById(id int) error
	ContainsCarById(id int) (bool, error)
	UpdateCarById(id int, car model.Car) error
	UpdatePeopleById(id int, people model.People) error
	AddNewCars(regNums []string) error
	ContainsPeopleById(id int) (bool, error)
}

func New(carService carService) router {
	router := router{
		ginRouter:  gin.New(),
		carService: carService,
	}
	router.ginRouter.Use(CORSMiddleware())
	router.ginRouter.GET("/cars", router.getCars())
	router.ginRouter.DELETE("/car", router.deleteCar())
	router.ginRouter.POST("/cars", router.addNewCars())
	router.ginRouter.POST("/car/update", router.updateCar())
	router.ginRouter.POST("/people/update", router.updatePeople())
	router.initSwagger()
	return router
}

func (r *router) initSwagger() {
	r.ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}

func (r *router) StartServer() {
	r.ginRouter.Run(":8080")
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// AddNewCars
//
// @Summary		Add Cars
// @Description	Add cars by regional number
// @Tags			Cars
// @Accept			json
// @Produce		json
// @Param			regNum	body		array 	true	"Cars regional numbers"
// @Success		200	"Cars added"
// @Failure		400 "Bad request"
// @Router			/cars [POST]
func (r *router) addNewCars() func(c *gin.Context) {
	return func(c *gin.Context) {

		var nums struct {
			Nums []string `json:"regNums"`
		}

		if err := c.ShouldBindJSON(&nums); err != nil {
			logrus.Debug(err)
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		logrus.Debug(nums)

		if len(nums.Nums) == 0 {
			logrus.Debug(nums)
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		err := r.carService.AddNewCars(nums.Nums)
		if err != nil {
			logrus.Debug(err)
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		c.JSON(http.StatusOK, "Cars added")
	}
}

// Delete car
//
// @Summary		Delete Car
// @Description	Delete car by Id
// @Tags			Cars
// @Accept			json
// @Produce		json
// @Param			id	body	int 	true	"Car Id"
// @Success		200	"Car deleted"
// @Failure		400 "Bad request"
// @Failure		500 "Internal Server Error"
// @Router		/car [DELETE]
func (r *router) deleteCar() func(c *gin.Context) {
	return func(c *gin.Context) {
		var id struct {
			Id int `json:"id"`
		}
		if err := c.ShouldBindJSON(&id); err != nil {
			logrus.Debug(err)
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		exist, err := r.carService.ContainsCarById(id.Id)
		if err != nil || !exist {
			logrus.Debug(err)

			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		err = r.carService.DeleteCarById(id.Id)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.JSON(http.StatusOK, "Car deleted")
	}
}

// Get Cars
//
// @Summary			Get Cars
// @Description	Return cars. param=constant:value for filter. Equal by default. Constant: <br>lt = lower than <br>ltq = lower than equal <br>gt = Greater Than <br>gtq = Greater Than Equal<br> between = Between
// @Tags				Cars
// @Accept			json
// @Produce			json
// @Param       limit   query     string  false  "Limit on cars in response"
// @Param       offset  query     string  false  "Offset from begin"
// @Param       order   query     string  false  "Order desc or asc"
// @Param       id    	query     string  false  "car id"
// @Param       mark    query     string  false  "Car mark"
// @Param       model   query     string  false  "Car model"
// @Param       year    query     string  false  "year"
// @Success			200	{object} model.Car "Return array of cars"
// @Failure			400 "Bad request"
// @Router			/cars [GET]
func (r *router) getCars() func(c *gin.Context) {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()

		cars, err := r.carService.GetCarsWithOwners(query)
		if err != nil {
			logrus.Debug(err)
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		c.JSON(http.StatusOK, cars)
	}
}

// Update car
//
// @Summary		Update Car
// @Description	Update car by Id
// @Tags			Cars
// @Accept			json
// @Produce		json
// @Param			car	body	model.Car 	true	"Change car date from body by id"
// @Success		200	"Car Updated"
// @Failure		400 "Bad request"
// @Failure		500 "Internal Server Error"
// @Router		/car/update [POST]
func (r *router) updateCar() func(c *gin.Context) {
	return func(c *gin.Context) {
		var car model.Car
		if err := c.ShouldBindJSON(&car); err != nil {
			logrus.Debug(err)
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		exist, err := r.carService.ContainsCarById(car.Id)

		if err != nil || !exist {
			c.JSON(http.StatusBadRequest, "Bad request")
			logrus.Debug(err)
			return
		}

		err = r.carService.UpdateCarById(car.Id, car)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.JSON(http.StatusOK, "Car Updated")
	}
}

// Update people
//
// @Summary		Update people
// @Description	Update people by Id
// @Tags			Peoples
// @Accept			json
// @Produce		json
// @Param			people	body	model.People 	true	"Change people date from body by id"
// @Success		200	"People Updated"
// @Failure		400 "Bad request"
// @Failure		500 "Internal Server Error"
// @Router		/people/update [POST]
func (r *router) updatePeople() func(c *gin.Context) {
	return func(c *gin.Context) {
		var people model.People
		if err := c.ShouldBindJSON(&people); err != nil {
			logrus.Debug(err)
			c.JSON(http.StatusBadRequest, "Bad request")
			return
		}
		exist, err := r.carService.ContainsPeopleById(people.Id)

		if err != nil || !exist {
			c.JSON(http.StatusBadRequest, "Bad request")
			logrus.Debug(err)
			return
		}

		err = r.carService.UpdatePeopleById(people.Id, people)
		if err != nil {
			logrus.Info(err)
			c.JSON(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		c.JSON(http.StatusOK, "People Updated")
	}
}
