package main

import (
	"chartdataservice/controller"
	"chartdataservice/service"
	"fmt"
	"log"
	"strconv"
	"strings"

	"runtime"
	"time"

	cache "github.com/TecXLab/libcache"
	"github.com/TecXLab/libcache/persistence"
	"github.com/TecXLab/libdb"
	"github.com/TecXLab/libdrainer"
	"github.com/TecXLab/libenv"
	"github.com/TecXLab/liblogs"
	"github.com/TecXLab/libprometheus"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/golobby/container/v3"
	"github.com/rs/zerolog"
)

var (
	// router   = gin.Default()
	router        = gin.New()
	Env           libenv.Env
	zerologs      zerolog.Logger
	Logger        = liblogs.Logger
	EnableCaching = true
)

func init() {
	// os.Setenv("VOLUME_PATH", "D:/GitHub/KubeVolumeData")
	// os.Setenv("MSQL_PASSWORD", "Mksl.1234") //Tecxlabs@123
	// os.Setenv("MSQL_USERNAME", "admin")
	// os.Setenv("MSQL_SERVICENAME", "nuuulightdb.c3vdaxham06l.ap-south-1.rds.amazonaws.com") //Tecxlabs@123
	// os.Setenv("MSQL_PORT", "3306")
	// os.Setenv("MSQL_contractmaster_DBNAME", "contractmaster_production")

	err := Env.Init()
	if err != nil {
		fmt.Println("Env Lib ", err)
		panic("Env Lib Not Initialize" + err.Error())
	}
	fmt.Println(Env.VOLUMNE_PATH)
	//use libredis
	err = liblogs.Init()
	if err != nil {
		fmt.Println("liblogs: ", err)
		panic("Logs Lib Not Initialize" + err.Error())
	}

	container.NamedResolve(&zerologs, Env.DI_DB_ZEROLOGS)

	var q libdrainer.Q
	err = q.Init("userdrainer")
	if err != nil {
		panic("Drainer Lib Not Initialize" + err.Error())
	}
	var csdb libdb.ContractMasterDB
	err = csdb.Connect()
	if err != nil {
		zerologs.Error().Err(err).Msg("Error connecting Database")
		fmt.Println(err)
		panic("Contract Database Not Initialize" + err.Error())
	}
	if Env.STAGE == "local" {
		Env.PROMETHEUS_COUNTER = true
		Env.PROMETHEUS_HOOK = true
		Env.APPLABEL = "chartdataservice"
	}

	fmt.Println(Env.PROMETHEUS_COUNTER)
	fmt.Println(Env.PROMETHEUS_HOOK)
	fmt.Println(Env.APPLABEL)
	err = libprometheus.InitCounter([]libprometheus.CounterModel{
		{Name: "CDS_Enqueue_NSEEQ_Counter_Request", Info: "Total CDS_Enqueue_NSEEQ_Counter_Request"},
		{Name: "CDS_Dequeue_NSEEQ_Counter_Request", Info: "Total CDS_Dequeue_NSEEQ_Counter_Request"},
		{Name: "CDS_Store_Update_NSEEQ_Counter_Request", Info: "Total CDS_Store_Update_NSEEQ_Counter_Request"},

		{Name: "CDS_Enqueue_NSEFO_Counter_Request", Info: "Total CDS_Enqueue_NSEFO_Counter_Request"},
		{Name: "CDS_Dequeue_NSEFO_Counter_Request", Info: "Total CDS_Dequeue_NSEFO_Counter_Request"},
		{Name: "CDS_Store_Update_NSEFO_Counter_Request", Info: "Total CDS_Store_Update_NSEFO_Counter_Request"},

		{Name: "CDS_Enqueue_NSECD_Counter_Request", Info: "Total CDS_Enqueue_NSECD_Counter_Request"},
		{Name: "CDS_Dequeue_NSECD_Counter_Request", Info: "Total CDS_Dequeue_NSECD_Counter_Request"},
		{Name: "CDS_Store_Update_NSECD_Counter_Request", Info: "Total CDS_Store_Update_NSECD_Counter_Request"},

		{Name: "CDS_Enqueue_Input_MSG_Received_Count_Request", Info: "Total CDS_Enqueue_Input_MSG_Received_Count_Request"},
		{Name: "CDS_Dequeue_Input_MSG_Received_Count_Request", Info: "Total CDS_Dequeue_Input_MSG_Received_Count_Request"},

		{Name: "CDS_Enqueue_RAW_File_Writing_Count_Request", Info: "Total CDS_Enqueue_RAW_File_Writing_Count_Request"},
		{Name: "CDS_Dequeue_RAW_File_Writing_Count_Request", Info: "CDS_Dequeue_RAW_File_Writing_Count_Request"},

		{Name: "CDS_NSEEQ_Total_Token_Count_Request", Info: "Total CDS_NSEEQ_Total_Token_Count_Request"},
		{Name: "CDS_NSEFO_Total_Token_Count_Request", Info: "Total CDS_NSEFO_Total_Token_Count_Request"},
		{Name: "CDS_NSECD_Total_Token_Count_Request", Info: "Total CDS_NSECD_Total_Token_Count_Request"},
		{Name: "CDS_INDEX_Master_Total_Token_Count_Request", Info: "CDS_INDEX_Master_Total_Token_Count_Request"},
	})
	if err != nil {
		fmt.Println("Error while Init Promethus ")
	}

	err = service.InitDI()
	if err != nil {
		panic("Error while DI of All  Lib " + err.Error())
	}

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		zerologs.Error().Err(err).Msg("Error While setting local time zone.")
		fmt.Println("Error While setting local time zone. ", err)
	} else {
		time.Local = loc // -> this is setting the global timezone
	}
	service.InitVariable()
	service.AddTimeframeMap()
	service.InitPromethues()
	service.InitAll()
	service.InitChannel()
	//Run GC on startup
	runtime.GC()
	//Temp: Start
	// var size = flag.Int("size", 8, "file size in GiB")
	// flag.Parse()
	// fSize := int64(*size) * (1024 * 1024 * 1024)
	// err = service.WriteFile1(fSize)
	// if err != nil {
	// 	fmt.Fprintln(os.Stderr, fSize, err)
	// }
	//Temp: End
}
func main() {

	defer service.LogPanic()
	var zerologs zerolog.Logger
	container.NamedResolve(&zerologs, "zerologs")
	store := persistence.NewInMemoryStore(time.Hour * 24)

	router.Use(gzip.Gzip(gzip.DefaultCompression))
	if Env.PROMETHEUS_HOOK {
		router.Use(libprometheus.PrometheusMiddleware)
	}
	pprof.Register(router)

	api := router.Group("/chartdataservice")
	{
		api.POST("/historicalChartData", controller.SaveChartHistory)
		//
		api.GET("/gethistory/:exch/:symbol/:timefilter/:from/:to", cache.CachePage(store, time.Hour*24, controller.GetChartHistory, ChartHistoryCacheCheck))
		api.OPTIONS("/gethistory/:exch/:symbol/:timefilter/:from/:to", controller.GetChartHistory)
		api.GET("/gethistoryz/:exch/:symbol/:timefilter/:from/:to", cache.CachePage(store, time.Hour*24, controller.GetChartHistoryZ, ChartHistoryCacheCheck))
		api.OPTIONS("/gethistoryz/:exch/:symbol/:timefilter/:from/:to", controller.GetChartHistoryZ)

		api.GET("/getstroredata/:token/:exch", controller.GetDataFromStoreController)
		api.GET("/getrawdata/:token/:exch/:date", controller.GetRawdataService)

		api.GET("rawfilewriting", controller.RawFileWritingController)
		api.GET("InitAllCron", controller.InitAllCron)
		api.POST("/generatetffile", controller.GenerateTFFileController)
		api.POST("sendfiletoprod", controller.SendFileToProd)
		api.POST("getFileFromlocal", controller.GetFileFromLocal)
		api.POST("historicalChartDataOld", controller.SaveHistory_Old)
		api.GET("/filecrongeneratetf/:exch", controller.Filecrongeneratetf)
		api.POST("/filecron", controller.FileCronController)
		api.POST("DeleteFiles", controller.DeleteFilesController)
		api.POST("/marketstatus", controller.MarketStatusController)
		api.POST("/changefolderpermission", controller.ChangeFolderPermission)
		api.POST("/copyfolders", controller.CopyFolders)
		// api.GET("/getSize", service.GetStoreSize)
		api.GET("/changefoldername", controller.ChangeFolderName)

		api.POST("/clearCache", func(c *gin.Context) {
			cache.ClearStoreCache(store)
			c.JSON(200, "Cleared Cache")
		})
		api.POST("/toggleCaching", ToggleCaching)

	}
	controller.ReportsApi(router)
	// router.Run(":8080")
	sPort := Env.EXPOSEPORT
	if sPort == "" {
		fmt.Println("Started on Port 8080")
		zerologs.Fatal().Err(router.Run(":8080"))
		log.Fatal(router.Run(":8080"))
	} else {
		zerologs.Fatal().Err(router.Run(":" + sPort))
		log.Fatal(router.Run(":" + sPort))
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service.CoreHeader(c)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func ChartHistoryCacheCheck(c *gin.Context) bool {
	if !EnableCaching {
		return false
	}
	var FromDate string
	var ToDate string
	if c.Params[3].Key == "from" {
		FromDate = c.Params[3].Value
	}
	if c.Params[4].Key == "to" {
		ToDate = c.Params[4].Value
	}
	fromDate, _ := GetDateFromFileName(FromDate, "20060102")
	toDate, _ := GetDateFromFileName(ToDate, "20060102")

	if DateEqual(toDate, time.Now()) || DateEqual(fromDate, time.Now()) {
		return false
	}

	if fromDate.Before(time.Now()) && (toDate.After(time.Now())) {
		return false
	}
	return true
}

func GetDateFromFileName(timeString string, format string) (time.Time, error) {
	timeString = strings.Split(timeString, ".")[0]
	return time.Parse(format, timeString)
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func ToggleCaching(c *gin.Context) {
	temp := c.Query("status")
	b1, _ := strconv.ParseBool(temp)
	EnableCaching = b1
	if EnableCaching {
		c.JSON(200, "Caching is now Enabled")
	} else {
		c.JSON(200, "Caching is now Disabled")

	}

}
