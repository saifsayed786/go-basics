package service

import (
	"bufio"
	"bytes"
	"chartdataservice/models"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"runtime"

	// "io"

	"path/filepath"

	"os"
	"strconv"
	"strings"
	"time"

	"github.com/TecXLab/liblogs"
)

var (
	Store             [][ExchSize]models.ChartModel
	ContractStore     [ExchSize]models.ContractMapModel
	Epoch1980         = time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC)
	nDqCount          = 1
	nEnqCount         = 1
	fwriteDqcount     = 1
	fwriteEnqcount    = 1
	fwriteDqcount_FO  = 1
	fwriteEnqcount_FO = 1
	ch1               = make(chan models.NewChartRespModel, 20000)
	chFWriter         = make(chan models.NewChartRespModel, 20000)
	chFWriterFO       = make(chan models.NewChartRespModel, 20000)
	//chNSECD           = make(chan models.NewChartRespModel, 20000)
	time_Index   = 0
	op_Index     = 1
	hp_Index     = 2
	lp_Index     = 3
	cp_Index     = 4
	vl_Index     = 5
	TimeframeMap = make(map[string]models.TimeFrameModel) //Key : Timeframe
	Logger       = liblogs.Logger
)

func Enqueue(ch chan models.NewChartRespModel, val models.NewChartRespModel) {
	ch <- val
}

func EnqueueFile(ch chan models.NewChartRespModel, val models.NewChartRespModel) {
	ch <- val
}

func StoreInit(maxSize int) {
	Store = make([][ExchSize]models.ChartModel, maxSize)

}
func Dequeue(ch chan models.NewChartRespModel) {
	var m models.NewChartRespModel
	for m = range ch {
		nDqCount++
		start := time.Now()

		SaveHistoryDeqNew(m)

		elapsed := time.Since(start)
		str := fmt.Sprintf("%v", elapsed)
		Logger(("$$$ End EQ Store Update: " + " " + strconv.Itoa(len(m.ChartHLOCModel)) + " " + strconv.Itoa(nDqCount) + " " + str), nil, liblogs.Info, liblogs.PRINT|liblogs.ZEROLOG)
		if m.ExchID == 2 {
			Add_CDS_Store_Update_NSEEQ_Count()
			Add_CDS_Dequeue_NSEEQ_Count()
		} else if m.ExchID == 3 {
			Add_CDS_Store_Update_NSECD_Count()
			Add_CDS_Dequeue_NSECD_Count()
		}
		// zerologs.Info().Msg(time.Now().Format("2006-01-02 15:04:05") + " $$$ End EQ Store Update: " + " " + strconv.Itoa(len(m.ChartHLOCModel)) + " " + strconv.Itoa(nDqCount) + " " + str)
		// fmt.Println(time.Now().Format("2006-01-02 15:04:05")+" $$$ End EQ Store Update: ", len(m.ChartHLOCModel), nDqCount, str)
	}
}

func Dequeue_NSECD(ch chan models.NewChartRespModel) {
	var m models.NewChartRespModel
	for m = range ch {
		nDqCount++
		start := time.Now()

		SaveHistoryDeqNew(m)

		elapsed := time.Since(start)
		str := fmt.Sprintf("%v", elapsed)
		Logger(("$$$ End CD Store Update: " + " " + strconv.Itoa(len(m.ChartHLOCModel)) + " " + strconv.Itoa(nDqCount) + " " + str), nil, liblogs.Info, liblogs.PRINT|liblogs.ZEROLOG)
		if m.ExchID == 2 {
			Add_CDS_Store_Update_NSEEQ_Count()
			Add_CDS_Dequeue_NSEEQ_Count()
		} else if m.ExchID == 3 {
			Add_CDS_Store_Update_NSECD_Count()
			Add_CDS_Dequeue_NSECD_Count()
		}
		// zerologs.Info().Msg(time.Now().Format("2006-01-02 15:04:05") + " $$$ End EQ Store Update: " + " " + strconv.Itoa(len(m.ChartHLOCModel)) + " " + strconv.Itoa(nDqCount) + " " + str)
		// fmt.Println(time.Now().Format("2006-01-02 15:04:05")+" $$$ End EQ Store Update: ", len(m.ChartHLOCModel), nDqCount, str)
	}

}

func DequeueFile(ch chan models.NewChartRespModel) {
	var m models.NewChartRespModel
	for m = range ch {
		fwriteDqcount++
		start := time.Now()
		// zerologs.Info().Msg(time.Now().Format("2006-01-02 15:04:05") + " Start File Wrt: " + " " + strconv.Itoa(len(m.ChartHLOCModel)) + " " + strconv.Itoa(fwriteDqcount))
		// fmt.Println(time.Now().Format("2006-01-02 15:04:05")+" Start File Wrt: ", len(m.ChartHLOCModel), fwriteDqcount)
		Logger(time.Now().Format("2006-01-02 15:04:05")+" Start File Wrt: "+" "+strconv.Itoa(len(m.ChartHLOCModel))+" "+strconv.Itoa(fwriteDqcount), nil, liblogs.Info, liblogs.PRINT|liblogs.ZEROLOG)
		WriteRequestFile(m)
		elapsed := time.Since(start)
		str := fmt.Sprintf("%v", elapsed)
		Logger(("$$$ End File Wrt: " + " " + strconv.Itoa(len(m.ChartHLOCModel)) + " " + strconv.Itoa(fwriteDqcount) + " " + str), nil, liblogs.Info, liblogs.PRINT|liblogs.ZEROLOG)
		Add_CDS_Dequeue_RAW_File_Writing_Count()

		// zerologs.Info().Msg(time.Now().Format("2006-01-02 15:04:05") + " $$$ End File Wrt: " + " " + strconv.Itoa(len(m.ChartHLOCModel)) + " " + strconv.Itoa(fwriteDqcount) + " " + str)
		// fmt.Println(time.Now().Format("2006-01-02 15:04:05")+" $$$ End File Wrt: ", len(m.ChartHLOCModel), fwriteDqcount, str)
	}

}

func printAlloc() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%d MB\n", (m.Alloc/1024)/1024)
}

func DequeueFileFO(ch chan models.NewChartRespModel) {
	var m models.NewChartRespModel
	for m = range ch {
		fwriteDqcount_FO++
		start := time.Now()
		// zerologs.Info().Msg(time.Now().Format("2006-01-02 15:04:05") + " Start File Wrt FO: " + " " + strconv.Itoa(len(m.ChartHLOCModel)) + " " + strconv.Itoa(fwriteDqcount_FO))
		// fmt.Println(time.Now().Format("2006-01-02 15:04:05")+" Start File Wrt FO: ", len(m.ChartHLOCModel), fwriteDqcount_FO)
		Logger(time.Now().Format("2006-01-02 15:04:05")+" Start File Wrt FO: "+" "+strconv.Itoa(len(m.ChartHLOCModel))+" "+strconv.Itoa(fwriteDqcount_FO), nil, liblogs.Info, liblogs.PRINT|liblogs.ZEROLOG)
		SaveHistoryDeqNew(m)

		elapsed := time.Since(start)
		str := fmt.Sprintf("%v", elapsed)
		// zerologs.Info().Msg(time.Now().Format("2006-01-02 15:04:05") + " $$$ End FO Store Update: " + " " + strconv.Itoa(len(m.ChartHLOCModel)) + " " + strconv.Itoa(fwriteDqcount_FO) + " " + str)
		// fmt.Println(time.Now().Format("2006-01-02 15:04:05")+" $$$ End FO Store Update: ", len(m.ChartHLOCModel), fwriteDqcount_FO, str)
		Logger((" $$$ End FO Store Update: " + " " + strconv.Itoa(len(m.ChartHLOCModel)) + " " + strconv.Itoa(fwriteDqcount_FO) + " " + str), nil, liblogs.Info, liblogs.PRINT|liblogs.ZEROLOG)

		Add_CDS_Store_Update_NSEFO_Count()
		Add_CDS_Dequeue_NSEFO_Count()
	}
}

func InitChannel() {

	go Dequeue(ch1)
	go DequeueFile(chFWriter)
	go DequeueFileFO(chFWriterFO)
	go DequeueFC(chFileCron)
	go DequeueSendProdFile(chSendProd)
	go DequeueChTrueData(chTruedata)
	//go Dequeue_NSECD(chNSECD)
}

func SaveHistory(response models.NewChartRespModel) {
	EnqueueFile(chFWriter, response)
	Add_CDS_Enqueue_RAW_File_Writing_Count()
	if response.ExchID == ExIds["nseeq"] {
		Enqueue(ch1, response)
		Add_CDS_Enqueue_NSEEQ_Count()
	} else if response.ExchID == ExIds["nsefo"] {
		Enqueue(chFWriterFO, response)
		Add_CDS_Enqueue_NSEFO_Count()
	} else if response.ExchID == ExIds["nsecd"] {
		Enqueue(ch1, response)
		Add_CDS_Enqueue_NSECD_Count()
		//TODO add prothmeus counter
	}
}

func SaveHistoryDeqNew(DrainerData models.NewChartRespModel) {
	// sort.Slice(DrainerData.ChartHLOCModel, func(i, j int) bool {
	// 	return DrainerData.ChartHLOCModel[i].DtLogDateTime < DrainerData.ChartHLOCModel[j].DtLogDateTime
	// })

	// zerologs.Info().Msg(time.Now().Format("2006-01-02 15:04:05") + " start onDeq: {0}" + " " + strconv.Itoa(len(DrainerData.ChartHLOCModel)) + " " + strconv.Itoa(nDqCount))
	//fmt.Println(time.Now().Format("2006-01-02 15:04:05")+" start onDeq: {0}", len(DrainerData.ChartHLOCModel), nDqCount)
	Logger(time.Now().Format("2006-01-02 15:04:05")+" start onDeq: {0}"+" "+strconv.Itoa(len(DrainerData.ChartHLOCModel))+" "+strconv.Itoa(nDqCount), nil, liblogs.Info, liblogs.PRINT|liblogs.ZEROLOG)
	for i := range DrainerData.ChartHLOCModel {
		// if DrainerData.ChartHLOCModel[i].NToken != 1001 && DrainerData.ChartHLOCModel[i].NToken != 1002 {
		// 	continue
		// }
		CreateTimeFrameFile(DrainerData.ChartHLOCModel[i], 1, Minute1)
		CreateTimeFrameFile(DrainerData.ChartHLOCModel[i], 2, Minute2)
		CreateTimeFrameFile(DrainerData.ChartHLOCModel[i], 3, Minute3)
		CreateTimeFrameFile(DrainerData.ChartHLOCModel[i], 4, Minute4)
		CreateTimeFrameFile(DrainerData.ChartHLOCModel[i], 5, Minute5)
		CreateTimeFrameFile(DrainerData.ChartHLOCModel[i], 10, Minute10)
		CreateTimeFrameFile(DrainerData.ChartHLOCModel[i], 15, Minute15)
		CreateTimeFrameFile(DrainerData.ChartHLOCModel[i], 30, Minute30)
		CreateTimeFrameFile(DrainerData.ChartHLOCModel[i], 60, Hr1)
		CreateTimeFrameFile(DrainerData.ChartHLOCModel[i], 75, Minute75)
		CreateTimeFrameFile(DrainerData.ChartHLOCModel[i], 120, Hr2)
		CreateTimeFrameFile(DrainerData.ChartHLOCModel[i], 125, Minute125)
		CreateTimeFrameFile(DrainerData.ChartHLOCModel[i], 180, Hr3)
		CreateTimeFrameFile(DrainerData.ChartHLOCModel[i], 720, Day1)
	}
	// zerologs.Info().Msg(time.Now().Format("2006-01-02 15:04:05") + " #End onDeq: {0}}" + " " + strconv.Itoa(len(DrainerData.ChartHLOCModel)) + " " + strconv.Itoa(nDqCount))
	//fmt.Println(time.Now().Format("2006-01-02 15:04:05")+" #End onDeq: {0}}", len(DrainerData.ChartHLOCModel), nDqCount)
	Logger(time.Now().Format("2006-01-02 15:04:05")+" #End onDeq: {0}}"+" "+strconv.Itoa(len(DrainerData.ChartHLOCModel))+" "+strconv.Itoa(nDqCount), nil, liblogs.Info, liblogs.PRINT|liblogs.ZEROLOG)
}

func CreateTimeFrameFile(ObjHLOC models.ChartHLOCModel, minute int16, folderName string) {
	if len(Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap) == 0 {
		Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap = make(map[string]*models.ChartMapModel)
	}

	_, ok := Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName]
	if !ok {
		objChartMapModel := models.ChartMapModel{CSVMap: make(map[int16]*models.CSVModel)}
		Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName] = &objChartMapModel
	}

	timeStamp, err := time.Parse(TimeFormat, ObjHLOC.DtLogDateTime)
	if err != nil {
		Logger("Error in Parsing Time", err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)

		// zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
		// fmt.Println(err)
	}

	//Get Store Index
	var nStoreIndex = GetTFIndex(folderName, ObjHLOC.ExchId, int16(timeStamp.Hour()), int16(timeStamp.Minute()))
	if nStoreIndex >= 0 {
		_, ok := Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName].CSVMap[nStoreIndex]
		if !ok {
			var CSV = GenerateCSV(ObjHLOC)
			var objCSV = new(models.CSVModel)
			objCSV.CSVarr = CSV
			//objCsvModel := models.CSVModel{CSVarr: CSV, DateTime: timeStamp}
			Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName].CSVMap[nStoreIndex] = objCSV
			// objstore, ok := Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName]
			// if ok {
			// 	objstore.CSVMap[nStoreIndex] = objCSV
			// 	objstore.LastReferenceTime = timeStamp
			// 	objstore.TimeDiffernce = minute
			// 	objstore.KeySlice = append(objstore.KeySlice, nStoreIndex)
			// }
			// Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName] = objstore

			// Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName].LastReferenceTime = timeStamp
			// Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName].TimeDiffernce = minute
			Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName].KeySlice = append(Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName].KeySlice, nStoreIndex)
		} else {
			var splitted = strings.Split(Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName].CSVMap[nStoreIndex].CSVarr, ",")

			// hp := StringToFloat64_New(Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName].CSVMap[nStoreIndex].CSVarr[hp_Index])
			var hp = StringToFloat64_New(splitted[hp_Index])
			if hp <= ObjHLOC.NHighPrice {
				// Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName].CSVMap[nStoreIndex].CSVarr[hp_Index] = strconv.FormatFloat(ObjHLOC.NHighPrice, 'f', -1, 64)
				splitted[hp_Index] = strconv.FormatFloat(ObjHLOC.NHighPrice, 'f', -1, 64)
			}

			// lp := StringToFloat64_New(Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName].CSVMap[nStoreIndex].CSVarr[lp_Index])
			var lp = StringToFloat64_New(splitted[lp_Index])
			if lp >= ObjHLOC.NLowPrice {
				// Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName].CSVMap[nStoreIndex].CSVarr[lp_Index] = strconv.FormatFloat(ObjHLOC.NLowPrice, 'f', -1, 64)
				splitted[lp_Index] = strconv.FormatFloat(ObjHLOC.NLowPrice, 'f', -1, 64)
			}
			// Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName].CSVMap[nStoreIndex].CSVarr[cp_Index] = strconv.FormatFloat(ObjHLOC.NClosingPrice, 'f', -1, 64)
			splitted[cp_Index] = strconv.FormatFloat(ObjHLOC.NClosingPrice, 'f', -1, 64)

			// vl := StringToFloat64_New(Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName].CSVMap[nStoreIndex].CSVarr[vl_Index])
			var vl = StringToFloat64_New(splitted[vl_Index])
			TotalVol := math.Abs(vl + ObjHLOC.NVolume)
			// Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName].CSVMap[nStoreIndex].CSVarr[vl_Index] = strconv.FormatFloat(TotalVol, 'f', -1, 64)
			var joinedStr string
			splitted[vl_Index] = strconv.FormatFloat(TotalVol, 'f', -1, 64)
			if ObjHLOC.ExchId == 2 {
				joinedStr = GenerateString(splitted)
			} else if ObjHLOC.ExchId == 1 {
				joinedStr = GenerateStringFO(splitted)
			} else if ObjHLOC.ExchId == 3 {
				joinedStr = GenerateStringCD(splitted)
			}

			Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName].CSVMap[nStoreIndex].CSVarr = joinedStr
		}

		// fmt.Println("###########")
		// for Key := range Store[ObjHLOC.NToken][ObjHLOC.ExchId].StoreMap[folderName].KeySlice {
		// 	fmt.Println(Key)
		// }
	} else {
		//Logger("Store Index is -1 for ExId: "+strconv.Itoa(int(ObjHLOC.ExchId))+" TF:"+folderName+" Time:"+ObjHLOC.DtLogDateTime, nil, liblogs.Info, liblogs.PRINT|liblogs.ZEROLOG)

		// zerologs.Error().Msg("Store Index is -1 for ExId: " + strconv.Itoa(int(ObjHLOC.ExchId)) + " TF:" + folderName + " Time:" + ObjHLOC.DtLogDateTime)

		// fmt.Println("Store Index is -1 for ExId: " + strconv.Itoa(int(ObjHLOC.ExchId)) + " TF:" + folderName + " Time:" + ObjHLOC.DtLogDateTime)
	}
}

func CreateDirectory(symbol string, ExchId int16, folderName string) string {
	var newpath string
	if ExchId == ExIds["nseeq"] {
		newpath = Env.VOLUMNE_PATH + "/" + BasePath + "/" + EqPath + "/" + symbol + "/" + folderName + "/"
		if _, err := os.Stat(newpath); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(newpath, os.ModePerm)
			if err != nil {
				Logger("Error Creating directory"+newpath, err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)

				// zerologs.Error().Err(err).Msg("Error Creating directory" + newpath)
				// fmt.Println(err, "Error token directory")
			}
		}
	} else if ExchId == ExIds["nsefo"] {
		newpath = Env.VOLUMNE_PATH + "/" + BasePath + "/" + FoPath + "/" + symbol + "/" + folderName + "/"
		if _, err := os.Stat(newpath); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(newpath, os.ModePerm)
			if err != nil {
				Logger("Error Creating directory"+newpath, err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)

				// zerologs.Error().Err(err).Msg("Error Creating directory" + newpath)
				// fmt.Println(err, "Error token directory")
			}
		}
	} else if ExchId == ExIds["nsecd"] {
		newpath = Env.VOLUMNE_PATH + "/" + BasePath + "/" + CdPath + "/" + symbol + "/" + folderName + "/"
		if _, err := os.Stat(newpath); errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(newpath, os.ModePerm)
			if err != nil {
				Logger("Error Creating directory"+newpath, err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)

				// zerologs.Error().Err(err).Msg("Error Creating directory" + newpath)
				// fmt.Println(err, "Error token directory")
			}
		}
	}
	return newpath
}

func GenerateCSV(objHOLC models.ChartHLOCModel) string {
	//sLineToInsert := objHOLC.DtLogDateTime + "+0530" + "," + fmt.Sprint(objHOLC.NOpenPrice) + "," + fmt.Sprint(objHOLC.NHighPrice) + "," + fmt.Sprint(objHOLC.NLowPrice) + "," + fmt.Sprint(objHOLC.NClosingPrice) + "," + fmt.Sprint(objHOLC.NVolume) //+ "\n"

	var sLineToInsert = objHOLC.DtLogDateTime + "+0530" + "," + strconv.FormatFloat(objHOLC.NOpenPrice, 'f', -1, 64) + "," + strconv.FormatFloat(objHOLC.NHighPrice, 'f', -1, 64) + "," + strconv.FormatFloat(objHOLC.NLowPrice, 'f', -1, 64) + "," + strconv.FormatFloat(objHOLC.NClosingPrice, 'f', -1, 64) + "," + strconv.FormatFloat(objHOLC.NVolume, 'f', -1, 64)

	return sLineToInsert
}

var buffer bytes.Buffer
var bufferFO bytes.Buffer
var bufferCD bytes.Buffer

func GenerateString(objHOLC []string) string {
	//sLineToInsert := objHOLC.DtLogDateTime + "+0530" + "," + fmt.Sprint(objHOLC.NOpenPrice) + "," + fmt.Sprint(objHOLC.NHighPrice) + "," + fmt.Sprint(objHOLC.NLowPrice) + "," + fmt.Sprint(objHOLC.NClosingPrice) + "," + fmt.Sprint(objHOLC.NVolume) + "\n"
	// sLineToInsert := []string{objHOLC.DtLogDateTime + "+0530", fmt.Sprint(objHOLC.NOpenPrice), fmt.Sprint(objHOLC.NHighPrice), fmt.Sprint(objHOLC.NLowPrice), fmt.Sprint(objHOLC.NClosingPrice), fmt.Sprint(objHOLC.NVolume)}
	buffer.Reset()
	for n := 0; n < len(objHOLC); n++ {
		if n == len(objHOLC)-1 {
			buffer.WriteString(objHOLC[n])
		} else {
			buffer.WriteString(objHOLC[n] + ",")
		}
	}
	return buffer.String()
}
func GenerateStringFO(objHOLC []string) string {
	//sLineToInsert := objHOLC.DtLogDateTime + "+0530" + "," + fmt.Sprint(objHOLC.NOpenPrice) + "," + fmt.Sprint(objHOLC.NHighPrice) + "," + fmt.Sprint(objHOLC.NLowPrice) + "," + fmt.Sprint(objHOLC.NClosingPrice) + "," + fmt.Sprint(objHOLC.NVolume) + "\n"
	// sLineToInsert := []string{objHOLC.DtLogDateTime + "+0530", fmt.Sprint(objHOLC.NOpenPrice), fmt.Sprint(objHOLC.NHighPrice), fmt.Sprint(objHOLC.NLowPrice), fmt.Sprint(objHOLC.NClosingPrice), fmt.Sprint(objHOLC.NVolume)}
	bufferFO.Reset()
	for n := 0; n < len(objHOLC); n++ {
		if n == len(objHOLC)-1 {
			bufferFO.WriteString(objHOLC[n])
		} else {
			bufferFO.WriteString(objHOLC[n] + ",")
		}
	}
	return bufferFO.String()
}

func GenerateStringCD(objHOLC []string) string {
	//sLineToInsert := objHOLC.DtLogDateTime + "+0530" + "," + fmt.Sprint(objHOLC.NOpenPrice) + "," + fmt.Sprint(objHOLC.NHighPrice) + "," + fmt.Sprint(objHOLC.NLowPrice) + "," + fmt.Sprint(objHOLC.NClosingPrice) + "," + fmt.Sprint(objHOLC.NVolume) + "\n"
	// sLineToInsert := []string{objHOLC.DtLogDateTime + "+0530", fmt.Sprint(objHOLC.NOpenPrice), fmt.Sprint(objHOLC.NHighPrice), fmt.Sprint(objHOLC.NLowPrice), fmt.Sprint(objHOLC.NClosingPrice), fmt.Sprint(objHOLC.NVolume)}
	bufferCD.Reset()
	for n := 0; n < len(objHOLC); n++ {
		if n == len(objHOLC)-1 {
			bufferCD.WriteString(objHOLC[n])
		} else {
			bufferCD.WriteString(objHOLC[n] + ",")
		}
	}
	return bufferCD.String()
}

func GetHistoryData(Exch int, FileName, MinuteFolder, FromDate, ToDate string) interface{} {
	var newpath string
	switch {
	case Exch == 1:
		newpath = Env.VOLUMNE_PATH + "/" + BasePath + "/" + FoPath + "/" + FileName + "/" + MinuteFolder + "/"
	case Exch == 2:
		newpath = Env.VOLUMNE_PATH + "/" + BasePath + "/" + EqPath + "/" + FileName + "/" + MinuteFolder + "/"
	case Exch == 3:
		newpath = Env.VOLUMNE_PATH + "/" + BasePath + "/" + CdPath + "/" + FileName + "/" + MinuteFolder + "/"
	default:
		return nil
	}
	if FileName == "" || MinuteFolder == "" {
		return nil
	}

	intFromDate, _ := strconv.Atoi(FromDate)
	intToDate, _ := strconv.Atoi(ToDate)
	if intFromDate > intToDate {
		ToDate, FromDate = FromDate, ToDate
	}

	DateFiles := GetdateRange(FromDate, ToDate)
	// fmt.Println("Daterange ", DateFiles)
	var content [][]string
	current_date := strings.Split(time.Now().Format(FileTimeFormat), "T")[0]
	for i := len(DateFiles) - 1; i >= 0; i-- {
		if DateFiles[i] == current_date {
			// fmt.Println("Checking data in Todays Store for file " + DateFiles[i])
			data := GetTodaysStore(Exch, FileName, MinuteFolder, ToDate)
			if len(data) == 0 {
				// fmt.Println("No data found from store for todays " + FileName)
				// return nil
			}
			if len(data) > 0 {
				content = append(content, data...)
			} else {
				// fmt.Println("file found in todays store " + DateFiles[i])
				readContent, err := CsvReadFile(newpath, DateFiles[i])
				if err != nil {
					fmt.Println(err)
				}
				content = append(content, readContent...)
			}
		} else {
			readContent, err := CsvReadFile(newpath, DateFiles[i])
			if err != nil {
				//fmt.Println(err)
			}
			if readContent != nil {
				content = append(content, readContent...)
			}
		}
	}
	return content
}

func UpDownServiceNew(sPath string, nExID int16) error {
	Logger("UpDownService Start", nil, liblogs.Info, liblogs.PRINT|liblogs.ZEROLOG)

	// zerologs.Info().Msg("UpDownService Start")
	// fmt.Println("UpDownService Start")
	start := time.Now()
	// var (
	// 	// EmptycontractDetails models.ContractDetails
	// 	current_date = GetCurrentDate()
	// 	folderMap    = map[string]string{
	// 		"nseforoot": Env.VOLUMNE_PATH + "/chartdataservice/rawdata/nse1/" + current_date,
	// 		"nseeqroot": Env.VOLUMNE_PATH + "/chartdataservice/rawdata/nsefo1/" + current_date,
	// 		// "nsecdroot": Env.VOLUMNE_PATH + "/chartdataservice/rawdata" + current_date + "/NSECD/",
	// 	}
	// )
	//for _, folder := range folderMap {
	if _, err := os.Stat(sPath); errors.Is(err, os.ErrNotExist) {
		Logger("Message: ", err, liblogs.Error, liblogs.PRINT)

		// fmt.Println(err)
		return err
	}
	fmt.Println(time.Now().Format("2006-01-02 15:04:05")+" Start Reading Folder: ", sPath)
	files, err := os.ReadDir(sPath)
	if err != nil {
		Logger("Message: ", err, liblogs.Error, liblogs.PRINT)

		// fmt.Println(err)
	}
	for i := range files {
		//current_date := "20220827"
		//fmt.Println(files[i].Name())
		fileInfo, err := os.ReadFile(sPath + "/" + files[i].Name())
		if err != nil {
			Logger("Message: ", err, liblogs.Error, liblogs.PRINT)

			// fmt.Println(err)
		}
		var request models.NewChartRespModel
		err = json.Unmarshal(fileInfo, &request)
		if err != nil {
			Logger("Message: ", err, liblogs.Error, liblogs.PRINT)

			// fmt.Println(err)
			return err
		}
		//SaveHistoryDeqNew(request)
		if nExID == ExIds["nseeq"] {
			Enqueue(ch1, request)
		} else if nExID == ExIds["nsefo"] {
			Enqueue(chFWriterFO, request)
		} else if nExID == ExIds["nsecd"] {
			Enqueue(ch1, request)
		}
	}
	//}
	Logger(fmt.Sprintf("UpDownService, execution time %s\n", time.Since(start)), nil, liblogs.Info, liblogs.PRINT|liblogs.ZEROLOG)

	// zerologs.Info().Msg("UpDownService End")
	// fmt.Println("UpDownService End")
	// log.Printf("UpDownService, execution time %s\n", time.Since(start))
	return nil
}

// func GenerateFileModel(contractDetails models.ContractDetails, folderName string, objHLOC models.ChartHLOCModel) interface{} {
// 	fileModel := models.FileWriterModel{
// 		// CsvMap:         Store[objHLOC.NToken][objHLOC.ExchId].StoreMap[folderName].CSVMap,
// 		ExchId:         objHLOC.ExchId,
// 		Ntoken:         objHLOC.NToken,
// 		DtLogTime:      objHLOC.DtLogDateTime,
// 		ChartHLOCModel: Store[objHLOC.NToken][objHLOC.ExchId].StoreMap[folderName].ChartHLOCModel,
// 	}
// 	return fileModel
// }

// func writer(objData []models.FileWriterModel) {
// 	for i := range objData {
// 		objFileM := objData[i]
// 		FileWriting2(Store[objFileM.Ntoken][objFileM.ExchId].ContractDetails, "1Min", true, objFileM.CsvMap, objFileM.ExchId, objFileM.DtLogTime, objFileM.ChartHLOCModel)
// 	}
// }

func FileWriting2(contractDetails models.ContractDetails, folderName string, IsUpdate bool, csvMap map[string]string, ExchId int16, DtLogDateTime string, arrChartM []models.ChartHLOCModel) {
	var filePath string
	if len(contractDetails.SSymbol) == 0 {
		return
	}
	if ExchId == ExIds["nseeq"] {
		filePath = CreateDirectory(contractDetails.SSymbol, ExchId, folderName)
	} else if ExchId == ExIds["nsefo"] {
		if contractDetails.NStrikePrice > 0 {
			expiryDate := ConvertDateTime1980(contractDetails.NExpiryDate)
			strikePrice := strconv.FormatInt((contractDetails.NStrikePrice / 100), 10)
			//TCS22AUG3300CE
			// TCS29SEP223200CE
			var b bytes.Buffer
			b.WriteString(contractDetails.SSymbol)
			b.WriteString(expiryDate)
			b.WriteString(strikePrice)
			b.WriteString(contractDetails.SOptionType)
			filePath = CreateDirectory(b.String(), ExchId, folderName)

		} else if contractDetails.NStrikePrice == -1 {
			//TCS22AUGFUT
			expiryDate := ConvertDateTime1980(contractDetails.NExpiryDate)
			var b bytes.Buffer
			b.WriteString(contractDetails.SSymbol)
			b.WriteString(expiryDate)
			b.WriteString("FUT")
			filePath = CreateDirectory(b.String(), ExchId, folderName)
		}
	}

	fileTimeStamp, err := time.Parse(TimeFormat, DtLogDateTime)
	if err != nil {
		Logger("Error parsing Time fileTimeStamp", err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)

		// zerologs.Error().Err(err).Msg("Error parsing Time fileTimeStamp" + err.Error())
		// fmt.Println(err)
	}
	fileName := strings.Split(fileTimeStamp.Format(FileTimeFormat), "T")[0]

	path := filePath + fileName

	file, err := os.OpenFile(path+".csv", os.O_RDWR, 0777)
	if err != nil {
		file, _ = os.OpenFile(path+".csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	}
	defer file.Close()
	// w := csv.NewWriter(file)
	// defer w.Flush()
	var buffer bytes.Buffer
	// iterate over the slice
	for _, data := range csvMap {
		buffer.WriteString(data)
	}
	// for _, data := range arrChartM {
	// 	var sCSV = GenerateCSV(data)
	// 	buffer.WriteString(sCSV)
	// }
	_, err = file.Write(buffer.Bytes())
	if err != nil {
		Logger("Error parsing Time fileTimeStamp", err, liblogs.Error, liblogs.PRINT)

		// fmt.Println(err)
		return
	}
}

var sCSVBuilder string

func FileWriting3(contractDetails models.ContractDetails, folderName string, IsUpdate bool, objChartMapModel models.ChartMapModel, ExchId int16, DtLogDateTime string, arrChartM []models.ChartHLOCModel, date string) {

	if len(contractDetails.SFullName) == 0 {
		return
	}

	var filePath = CreateDirectory(contractDetails.SFullName, ExchId, folderName)

	// fileTimeStamp, err := time.Parse(TimeFormat, DtLogDateTime)
	// if err != nil {
	// 	zerologs.Error().Err(err).Msg("Error parsing Time fileTimeStamp" + err.Error())
	// 	fmt.Println(err)
	// }
	// fileName := strings.Split(fileTimeStamp.Format(FileTimeFormat), "T")[0]

	path := filePath + date

	file, err := os.OpenFile(path+".csv", os.O_RDWR, 0777)
	if err != nil {
		file, _ = os.OpenFile(path+".csv", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	}
	defer file.Close()

	wr := bufio.NewWriter(file)

	// for _, line := range arrCSVModel {

	// 	var sCSV = GenerateString(line.CSVarr)
	// 	wr.WriteString(sCSV + "\n")
	// }

	if contractDetails.NToken > 0 {
		sCSVBuilder = ""
		for i := len(objChartMapModel.KeySlice) - 1; i >= 0; i-- {
			// var sCSV = GenerateString(objChartMapModel.CSVMap[objChartMapModel.KeySlice[i]].CSVarr)
			// wr.WriteString(sCSV + "\n")
			sCSVBuilder += objChartMapModel.CSVMap[objChartMapModel.KeySlice[i]].CSVarr + "\n"
		}
		wr.WriteString(sCSVBuilder)
	}

	// for index := len(arrChartM) - 1; index >= 0; index-- {
	// 	var sCSV = GenerateCSVnew(arrChartM[index])
	// 	// buffer.WriteString(sCSV)
	// 	wr.WriteString(sCSV + "\n")
	// }
	err = wr.Flush()
	// if err != nil {
	// 	return err
	// }
	err = file.Sync()
	// if err != nil {
	// 	return err
	// }
}

func CreateFileCron(ExchID int16, inputDate string, sExName string) string {
	start := time.Now()
	var SlackMsg string
	// fmt.Println(sExName + " Filecron Start")
	Logger(sExName+" Filecron Start", nil, liblogs.Info, liblogs.PRINT)

	SlackMsg = "EOD CRON\n" + sExName + " Filecron Start\n"
	for token := range Store {
		// if i != 22 {
		// 	continue
		// }
		for j := range Store[token] {
			if j != int(ExchID) {
				continue
			}
			if Store[token][j].ContractDetails.NToken <= 0 {
				continue
			}

			if Store[token][j].StoreMap == nil {
				continue
			}

			contractDetails := Store[token][j].ContractDetails
			if contractDetails.NToken%100 == 0 {
				fmt.Println(contractDetails)
			}
			for TF := range Store[token][j].StoreMap {

				// if k != "1Min" {
				// 	continue
				// }

				if len(Store[token][j].StoreMap[TF].CSVMap) == 0 {
					continue
				}

				if len(Store[token][j].StoreMap[TF].KeySlice) == 0 {
					continue
				}

				// fmt.Println("TF: " + k + " CSVMap Count :" + strconv.Itoa(len(Store[i][j].StoreMap[k].CSVModelarr)))
				//time1 := time.Now()
				FileWriting3(contractDetails, TF, true, *Store[token][j].StoreMap[TF], int16(j), inputDate, nil, inputDate)
				//log.Printf("file writing time %s\n", time.Since(time1))
			}
		}
	}

	msg := fmt.Sprintf(sExName+" Filecron, execution time %s\n", time.Since(start))
	Logger(msg, nil, liblogs.Info, liblogs.PRINT)
	SlackMsg += sExName + " Filecron, execution time " + time.Since(start).String() + "\n"
	SlackMsg += sExName + " Filecron End\n"

	return SlackMsg
}

func WriteFile1(fSize int64) error {
	// fName := `/home/peter/diskio` // test file
	fName := Env.VOLUMNE_PATH + "ABC.txt" // test file
	//defer os.Remove(fName)
	f, err := os.Create(fName)
	if err != nil {
		return err
	}
	const defaultBufSize = 4096
	buf := make([]byte, defaultBufSize)
	buf[len(buf)-1] = '\n'
	buf[len(buf)-2] = 'P'
	w := bufio.NewWriterSize(f, len(buf))

	start := time.Now()
	written := int64(0)
	for i := int64(0); i < fSize; i += int64(len(buf)) {
		nn, err := w.Write(buf)
		written += int64(nn)
		if err != nil {
			return err
		}
		if i == 12288 {
			break
		}
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}
	since := time.Since(start)

	err = f.Close()
	if err != nil {
		return err
	}

	Logger(fmt.Sprintf("written: %dB %dns %.2fGB %.2fs %.2fMB/s\n",
		written, since,
		float64(written)/1000000000, float64(since)/float64(time.Second),
		(float64(written)/1000000)/(float64(since)/float64(time.Second)),
	), nil, liblogs.Info, liblogs.PRINT)

	return nil
}

func CreateAllDir(sFullName string, folderName string, ExchId int16, sTodaysDate string) {
	//CreateDirectory(sFullName, ExchId, folderName)
	//var sDirPath = CreateDirectory(sFullName, ExchId, folderName)
	//sCreateAllFiles(sDirPath, sTodaysDate)
}

func CreateAllFiles(sDirPath string, sTodaysDate string) *os.File {
	sFilepath := sDirPath
	var file *os.File
	if sTodaysDate != "" {
		sFilepath += sTodaysDate + ".csv"
	}
	if _, err := os.Stat(sFilepath); errors.Is(err, os.ErrNotExist) {
		// path/to/whatever does not exist
		file, err = os.Create(sFilepath)
		if err != nil {
			Logger("Error Creating file "+sFilepath, err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)
			// zerologs.Error().Err(err).Msg("Error Creating file " + sFilepath)
			// fmt.Println(err, "Error Creating file "+sFilepath)
		}
	}
	return file
}

func RawFileWriting(objdata models.NewChartRespModel) {
	WriteRequestFile(objdata)
}

func InitAll() string {
	// for i := range Store {
	// 	for j := range Store[i] {
	// 		for keys := range Store[i][j].StoreMap {
	// 			delete(Store[i][j].StoreMap, keys)
	// 		}
	// 		Store[i][j].ContractDetails = models.ContractDetails{}
	// 	}
	// 	Store[i] = [ExchSize]models.ChartModel{}
	// }
	//Store = nil
	runtime.GC()
	// Store = [][4]models.ChartModel{}
	// ContractStore = [4]models.ContractMapModel{}
	ClearMemory()
	runtime.GC()
	msg := StoreContract(0, nil)
	Logger(msg, nil, liblogs.Info, liblogs.SLACK)

	// err := SentToContractCronJob(msg)
	// if err != nil {
	// 	Logger(fmt.Sprintf("Error while sending to slack, msg : "+msg), err, liblogs.Error, liblogs.PRINT)

	// }

	current_date := GetCurrentDate()
	UpDownServiceNew(Env.VOLUMNE_PATH+"/"+BasePath+"/"+RawPath+"/"+FoPath+"/"+current_date, 1)
	UpDownServiceNew(Env.VOLUMNE_PATH+"/"+BasePath+"/"+RawPath+"/"+EqPath+"/"+current_date, 2)
	UpDownServiceNew(Env.VOLUMNE_PATH+"/"+BasePath+"/"+RawPath+"/"+CdPath+"/"+current_date, 3)
	
	return msg
}

func AddTimeframeMap() {
	if len(TimeframeMap) == 0 {
		TimeframeMap = make(map[string]models.TimeFrameModel)
	}

	NSEEQ_Start_Time, err := time.Parse(NSETimeFormat, NSEEQ_Start_Time)
	if err != nil {
		Logger("Error in Parsing Time", err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)

		// zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
		// fmt.Println(err)
	}

	NSEEQ_End_Time, err := time.Parse(NSETimeFormat, NSEEQ_End_Time)
	if err != nil {
		Logger("Error in Parsing Time", err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)

		// zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
		// fmt.Println(err)
	}

	NSEFO_Start_Time, err := time.Parse(NSETimeFormat, NSEFO_Start_Time)
	if err != nil {
		Logger("Error in Parsing Time", err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)

		// zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
		// fmt.Println(err)
	}

	NSEFO_End_Time, err := time.Parse(NSETimeFormat, NSEFO_End_Time)
	if err != nil {
		Logger("Error in Parsing Time", err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)

		// zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
		// fmt.Println(err)
	}

	NSECD_Start_Time, err := time.Parse(NSETimeFormat, NSECD_Start_Time)
	if err != nil {
		Logger("Error in Parsing Time", err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)

		// zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
		// fmt.Println(err)
	}

	NSECD_End_Time, err := time.Parse(NSETimeFormat, NSECD_End_Time)
	if err != nil {
		Logger("Error in Parsing Time", err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)

		// zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
		// fmt.Println(err)
	}

	// fmt.Println(strconv.Itoa(NSEEQ_Start_Time.Hour()) + " Minute:" + strconv.Itoa(NSEEQ_Start_Time.Minute()))
	// fmt.Println(strconv.Itoa(NSEEQ_End_Time.Hour()) + " Minute:" + strconv.Itoa(NSEEQ_End_Time.Minute()))

	for Key, val := range GetMinuteMap() {
		// if Key != "4Min" {
		// 	continue
		// }
		// fmt.Println("###########")
		// fmt.Println(Key + " " + strconv.Itoa(int(val)))
		_, ok := TimeframeMap[Key]
		if !ok {
			ExTFarr := make([][24][60]int16, ExchSize)
			objTimeFrameModel := models.TimeFrameModel{ExTFarr: ExTFarr}
			TimeframeMap[Key] = objTimeFrameModel
			for exId := range ExTFarr {
				switch {
				case exId == (int)(ExIds["nseeq"]):
					var nStoreIndex int16 = 0
					JumpTime := NSEEQ_Start_Time
					for i := 0; i < len(TimeframeMap[Key].ExTFarr[exId]); i++ {
						if i >= NSEEQ_Start_Time.Hour() && i <= NSEEQ_End_Time.Hour() {
							for j := 0; j < len(TimeframeMap[Key].ExTFarr[exId][i]); j++ {
								var sHour = "0"
								var sMinute = "0"
								if i <= 9 {
									sHour += strconv.Itoa(i)
								} else {
									sHour = strconv.Itoa(i)
								}

								if j <= 9 {
									sMinute += strconv.Itoa(j)
								} else {
									sMinute = strconv.Itoa(j)
								}
								t := sHour + ":" + sMinute

								objTime, err := time.Parse(NSETimeFormat, t)
								if err != nil {
									//zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
									Logger("Error:", err, liblogs.Error, liblogs.PRINT)

									// fmt.Println(err)
								}
								if objTime.Equal(NSEEQ_Start_Time) || objTime.Equal(NSEEQ_End_Time) || (objTime.After(NSEEQ_Start_Time) && objTime.Before(NSEEQ_End_Time)) {
									sNumTime := sHour + sMinute
									intVar, err := strconv.Atoi(sNumTime)
									if err != nil {
										//zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
										Logger("Error:", err, liblogs.Error, liblogs.PRINT)

										// fmt.Println(err)
									}
									if objTime.Equal(NSEEQ_Start_Time) {
										nStoreIndex = 0
										JumpTime = objTime
										//fmt.Println(objTime)
										JumpTime = JumpTime.Add(time.Minute * time.Duration(int64(val)))
										//fmt.Println(JumpTime)
										TimeframeMap[Key].ExTFarr[exId][i][j] = nStoreIndex
										//fmt.Println(sHour + ":" + sMinute + " = " + strconv.Itoa(int(nStoreIndex)))
									} else if objTime.Before(JumpTime) {
										TimeframeMap[Key].ExTFarr[exId][i][j] = nStoreIndex
									} else if objTime.Equal(JumpTime) || objTime.After(JumpTime) {
										nIndex := GetIndex(int16(intVar), exId)
										nStoreIndex = (int16(nIndex))
										TimeframeMap[Key].ExTFarr[exId][i][j] = nStoreIndex
										JumpTime = objTime
										JumpTime = JumpTime.Add(time.Minute * time.Duration(int64(val)))
									}
								}
							}
						}
					}

				case exId == (int)(ExIds["nsefo"]):
					var nStoreIndex int16 = 0
					JumpTime := NSEFO_Start_Time
					for i := 0; i < len(TimeframeMap[Key].ExTFarr[exId]); i++ {
						if i >= NSEFO_Start_Time.Hour() && i <= NSEFO_End_Time.Hour() {
							for j := 0; j < len(TimeframeMap[Key].ExTFarr[exId][i]); j++ {
								var sHour = "0"
								var sMinute = "0"
								if i <= 9 {
									sHour += strconv.Itoa(i)
								} else {
									sHour = strconv.Itoa(i)
								}

								if j <= 9 {
									sMinute += strconv.Itoa(j)
								} else {
									sMinute = strconv.Itoa(j)
								}
								t := sHour + ":" + sMinute

								objTime, err := time.Parse(NSETimeFormat, t)
								if err != nil {
									//zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
									// fmt.Println(err)
									Logger("Error:", err, liblogs.Error, liblogs.PRINT)

								}
								if objTime.Equal(NSEFO_Start_Time) || objTime.Equal(NSEFO_End_Time) || (objTime.After(NSEFO_Start_Time) && objTime.Before(NSEFO_End_Time)) {
									sNumTime := sHour + sMinute
									intVar, err := strconv.Atoi(sNumTime)
									if err != nil {
										//zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
										// fmt.Println(err)
										Logger("Error:", err, liblogs.Error, liblogs.PRINT)

									}
									if objTime.Equal(NSEFO_Start_Time) {
										nStoreIndex = 0
										JumpTime = objTime
										//fmt.Println(objTime)
										JumpTime = JumpTime.Add(time.Minute * time.Duration(int64(val)))
										//fmt.Println(JumpTime)
										TimeframeMap[Key].ExTFarr[exId][i][j] = nStoreIndex
										//fmt.Println(sHour + ":" + sMinute + " = " + strconv.Itoa(int(nStoreIndex)))
									} else if objTime.Before(JumpTime) {
										TimeframeMap[Key].ExTFarr[exId][i][j] = nStoreIndex
									} else if objTime.Equal(JumpTime) || objTime.After(JumpTime) {
										nIndex := GetIndex(int16(intVar), exId)
										nStoreIndex = (int16(nIndex))
										TimeframeMap[Key].ExTFarr[exId][i][j] = nStoreIndex
										JumpTime = objTime
										JumpTime = JumpTime.Add(time.Minute * time.Duration(int64(val)))
									}
								}
							}
						}
					}

				case exId == (int)(ExIds["nsecd"]):
					var nStoreIndex int16 = 0
					JumpTime := NSECD_Start_Time
					for i := 0; i < len(TimeframeMap[Key].ExTFarr[exId]); i++ {
						if i >= NSECD_Start_Time.Hour() && i <= NSECD_End_Time.Hour() {
							for j := 0; j < len(TimeframeMap[Key].ExTFarr[exId][i]); j++ {
								var sHour = "0"
								var sMinute = "0"
								if i <= 9 {
									sHour += strconv.Itoa(i)
								} else {
									sHour = strconv.Itoa(i)
								}

								if j <= 9 {
									sMinute += strconv.Itoa(j)
								} else {
									sMinute = strconv.Itoa(j)
								}
								t := sHour + ":" + sMinute

								objTime, err := time.Parse(NSETimeFormat, t)
								if err != nil {
									//zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
									// fmt.Println(err)
									Logger("Error:", err, liblogs.Error, liblogs.PRINT)

								}
								if objTime.Equal(NSECD_Start_Time) || objTime.Equal(NSECD_End_Time) || (objTime.After(NSECD_Start_Time) && objTime.Before(NSECD_End_Time)) {
									sNumTime := sHour + sMinute
									intVar, err := strconv.Atoi(sNumTime)
									if err != nil {
										//zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
										// fmt.Println(err)
										Logger("Error:", err, liblogs.Error, liblogs.PRINT)

									}
									if objTime.Equal(NSECD_Start_Time) {
										nStoreIndex = 0
										JumpTime = objTime
										//fmt.Println(objTime)
										JumpTime = JumpTime.Add(time.Minute * time.Duration(int64(val)))
										//fmt.Println(JumpTime)
										TimeframeMap[Key].ExTFarr[exId][i][j] = nStoreIndex
										//fmt.Println(sHour + ":" + sMinute + " = " + strconv.Itoa(int(nStoreIndex)))
									} else if objTime.Before(JumpTime) {
										TimeframeMap[Key].ExTFarr[exId][i][j] = nStoreIndex
									} else if objTime.Equal(JumpTime) || objTime.After(JumpTime) {
										nIndex := GetIndex(int16(intVar), exId)
										nStoreIndex = (int16(nIndex))
										TimeframeMap[Key].ExTFarr[exId][i][j] = nStoreIndex
										JumpTime = objTime
										JumpTime = JumpTime.Add(time.Minute * time.Duration(int64(val)))
									}
								}
							}
						}
					}
				}
			}
		}
	}
	// fmt.Println("Doneeee")
	Logger("Doneeee", nil, liblogs.Info, liblogs.PRINT)

}

func GetTFIndex(sTF string, nExId int16, nHours int16, nMinute int16) int16 {
	if len(TimeframeMap) == 0 {
		Logger("Error: Zero entry TimeFrame Map.", nil, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)

		// zerologs.Error().Msg("Error: Zero entry TimeFrame Map.")
		// fmt.Println("Error: Zero entry TimeFrame Map.")
		return -1
	}
	var MktStartTime, MktEndTime time.Time

	switch nExId {
	case ExIds["nseeq"]:
		NSEEQ_Start_Time, err := time.Parse(NSETimeFormat, NSEEQ_Start_Time)
		if err != nil {
			Logger("Error in Parsing Time"+err.Error(), err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)
			// zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
			// fmt.Println(err)
		}

		NSEEQ_End_Time, err := time.Parse(NSETimeFormat, NSEEQ_End_Time)
		if err != nil {
			Logger("Error in Parsing Time"+err.Error(), err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)
			// zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
			// fmt.Println(err)
		}
		MktStartTime, MktEndTime = NSEEQ_Start_Time, NSEEQ_End_Time
	case ExIds["nsefo"]:
		NSEFO_Start_Time, err := time.Parse(NSETimeFormat, NSEFO_Start_Time)
		if err != nil {
			Logger("Error in Parsing Time"+err.Error(), err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)
			// zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
			// fmt.Println(err)
		}

		NSEFO_End_Time, err := time.Parse(NSETimeFormat, NSEFO_End_Time)
		if err != nil {
			Logger("Error in Parsing Time"+err.Error(), err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)
			// zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
			// fmt.Println(err)
		}
		MktStartTime, MktEndTime = NSEFO_Start_Time, NSEFO_End_Time
	case ExIds["nsecd"]:
		NSECD_Start_Time, err := time.Parse(NSETimeFormat, NSECD_Start_Time)
		if err != nil {
			Logger("Error in Parsing Time"+err.Error(), err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)
			// zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
			// fmt.Println(err)
		}

		NSECD_End_Time, err := time.Parse(NSETimeFormat, NSECD_End_Time)
		if err != nil {
			Logger("Error in Parsing Time"+err.Error(), err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)
			// zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
			// fmt.Println(err)
		}
		MktStartTime, MktEndTime = NSECD_Start_Time, NSECD_End_Time
	}

	var sHour = "0"
	var sMinute = "0"
	if nHours <= 9 {
		sHour += strconv.Itoa(int(nHours))
	} else {
		sHour = strconv.Itoa(int(nHours))
	}
	if nMinute <= 9 {
		sMinute += strconv.Itoa(int(nMinute))
	} else {
		sMinute = strconv.Itoa(int(nMinute))
	}
	t := sHour + ":" + sMinute
	objTime, err := time.Parse(NSETimeFormat, t)
	if err != nil {
		//zerologs.Error().Err(err).Msg("Error in Parsing Time" + err.Error())
		Logger("Error in Parsing Time"+err.Error(), err, liblogs.Error, liblogs.PRINT|liblogs.ZEROLOG)
		// fmt.Println(err)
	}
	if objTime.Equal(MktStartTime) || objTime.Equal(MktEndTime) || (objTime.After(MktStartTime) && objTime.Before(MktEndTime)) {
		objTFModel, ok := TimeframeMap[sTF]
		if ok {
			return objTFModel.ExTFarr[nExId][nHours][nMinute]
		}
	}
	return -1
}

func GetIndex(nTime int16, exId int) int {
	//fmt.Println(strconv.Itoa(int(nTime)) + " = " + strconv.Itoa(int((nTime - NSEEQ_Min_Time))))
	var min_time = 0
	switch exId {
	case (int)(ExIds["nseeq"]):
		min_time, _ = strconv.Atoi(NSEEQ_Min_Time)
	case (int)(ExIds["nsefo"]):
		min_time, _ = strconv.Atoi(NSEFO_Min_Time)
	case (int)(ExIds["nsecd"]):
		min_time, _ = strconv.Atoi(NSECD_Min_Time)
	}

	return (int)(nTime - int16(min_time))
}

var chFileCron = make(chan interface{}, 500)

func EnqueueFilecronData(val interface{}) {
	EnqueueFC(chFileCron, val)
}

func EnqueueFC(ch chan interface{}, val interface{}) {
	ch <- val
}

func DequeueFC(ch chan interface{}) {
	for objData := range ch {
		switch objData.(type) {
		case models.FileCronInput:
			GenerateTFFiles(objData)

		case models.GenerateTFFileModel:
			GenerateStoreDataForDate(objData)
		case string:
			// ChangeFolderPermission(objData)
			CopyFolders(objData)
		}
	}
}

var chSendProd = make(chan interface{}, 500000)

func EnqueuechSendProd(val interface{}) {
	EnqueuechSendProdFile(chSendProd, val)
}

func EnqueuechSendProdFile(ch chan interface{}, val interface{}) {
	ch <- val
}

func DequeueSendProdFile(ch chan interface{}) {
	for objData := range ch {
		switch objData.(type) {
		case models.BulkUploadFileModel:
			GetFileFromLocalHelper(objData)
		default:
			fmt.Println()
		}
	}
}
func GetFileFromLocalHelper(objdata interface{}) {
	var path string
	UploadFileModel := objdata.(models.BulkUploadFileModel)
	for i := range UploadFileModel.Bulkdata {
		if UploadFileModel.Bulkdata[i].ExchId == ExIds["nseeq"] {
			path = Env.VOLUMNE_PATH + "/" + "chartdataservice1" + "/" + EqPath + "/" + UploadFileModel.Bulkdata[i].Path
		} else if UploadFileModel.Bulkdata[i].ExchId == ExIds["nsefo1"] {
			path = Env.VOLUMNE_PATH + "/" + "chartdataservice1" + "/" + FoPath + "/" + UploadFileModel.Bulkdata[i].Path
		}
		splittedPath := strings.FieldsFunc(path, Split)
		folders := strings.Join(splittedPath[:len(splittedPath)-1], "/")
		if _, err := os.Stat(folders); os.IsNotExist(err) {
			err = os.MkdirAll(folders, 0777)
			if err != nil {
				Logger("Error in creating prodPath", err, liblogs.Error, liblogs.ZEROLOG)

				// zerologs.Error().Msg("Error in creating prodPath" + err.Error())
				return
			}
		}
		file := CreateAllFiles(path, "")
		defer file.Close()

		wr := bufio.NewWriter(file)
		wr.Write(UploadFileModel.Bulkdata[i].Data)
		wr.Flush()
		file.Sync()
		// err := os.WriteFile(path, UploadFileModel.Data, 0777)
		// if err != nil {
		// 	fmt.Println(err)
		// }
	}
}
func GenerateTFFiles(objData interface{}) {
	InputData := objData.(models.FileCronInput)
	Logger("File Cron Deque: Start", nil, liblogs.Info, liblogs.PRINT|liblogs.SLACK)
	// SentToContractCronJob(time.Now().String() + " File Cron Deque: Start")
	// fmt.Println(time.Now().String() + " File Cron Deque: Start")
	var current_date = GetCurrentDate()
	if InputData.Crondate != "" {
		current_date = InputData.Crondate
	}

	var slackMsg string
	if InputData.ExchId == "" {
		for ExName, val := range ExIds {
			runtime.GC()
			Logger(time.Now().String()+" IF Inside Loop: "+ExName+" TF File writing Started.", nil, 1, liblogs.SLACK|liblogs.PRINT)
			slackMsg = CreateFileCron(val, current_date, ExName)
			Logger(slackMsg, nil, liblogs.Info, liblogs.SLACK)
			runtime.GC()
			// SentToContractCronJob(slackMsg)
			// if err != nil {
			// 	fmt.Println("Error while sending to slack, msg : " + slackMsg)
			// }
		}
	} else {
		runtime.GC()
		Logger(time.Now().String()+" Inside Else: "+InputData.ExchId+" TF File writing Started.", nil, liblogs.Info, liblogs.PRINT|liblogs.ZEROLOG|liblogs.SLACK)

		nexid, _ := strconv.Atoi(InputData.ExchId)
		slackMsg = CreateFileCron(int16(nexid), current_date, GetExchngeName(int16(nexid)))
		Logger(slackMsg, nil, 1, liblogs.SLACK)

		Logger(" Inside Else: "+InputData.ExchId+" TF File writing Started.", nil, liblogs.Info, liblogs.PRINT|liblogs.SLACK)
		runtime.GC()
		// SentToContractCronJob(time.Now().String() + " Inside Else: " + InputData.ExchId + " TF File writing Started.")
		// fmt.Println(time.Now().String() + " Inside Else: " + InputData.ExchId + " TF File writing Started.")

		// if err != nil {
		// 	fmt.Println("Error while sending to slack, msg : " + slackMsg)
		// }
	}
	Logger(" File Cron Deque: End", nil, liblogs.Info, liblogs.PRINT|liblogs.SLACK)

	// SentToContractCronJob(time.Now().String() + " File Cron Deque: End")
	// fmt.Println(time.Now().String() + " File Cron Deque: End")
}

func GenerateStoreDataForDate(obj interface{}) {
	Logger(" GenerateStoreDataForDate Deque: Start", nil, liblogs.Info, liblogs.PRINT|liblogs.SLACK)

	// SentToContractCronJob(time.Now().String() + " GenerateStoreDataForDate Deque: Start")
	// fmt.Println(time.Now().String() + " GenerateStoreDataForDate Deque: Start")

	objData := obj.(models.GenerateTFFileModel)
	var sDate = GetCurrentDate()
	if objData.Date != "" {
		sDate = objData.Date
	}

	if objData.ExId == "" {
		for ExName, val := range ExIds {
			if ExName == "nseeq" {
				UpDownServiceNew(Env.VOLUMNE_PATH+"/"+BasePath+"/"+RawPath+"/"+EqPath+"/"+sDate, val)
			} else if ExName == "nsefo" {
				UpDownServiceNew(Env.VOLUMNE_PATH+"/"+BasePath+"/"+RawPath+"/"+FoPath+"/"+sDate, val)
			}
		}
	} else {
		nexid, _ := strconv.Atoi(objData.ExId)
		if GetExchngeName(int16(nexid)) == "nseeq" {
			UpDownServiceNew(Env.VOLUMNE_PATH+"/"+BasePath+"/"+RawPath+"/"+EqPath+"/"+sDate, int16(nexid))
		} else if GetExchngeName(int16(nexid)) == "nsefo" {
			UpDownServiceNew(Env.VOLUMNE_PATH+"/"+BasePath+"/"+RawPath+"/"+FoPath+"/"+sDate, int16(nexid))
		}
	}
	Logger(" GenerateStoreDataForDate Deque: End", nil, liblogs.Info, liblogs.PRINT|liblogs.SLACK)

	// SentToContractCronJob(time.Now().String() + " GenerateStoreDataForDate Deque: End")
	// fmt.Println(time.Now().String() + " GenerateStoreDataForDate Deque: End")
}

func DeleteFiles(Exch int16, symbol, Tf, Date string) error {
	var newpath string
	switch {
	case Exch == 1:
		newpath = Env.VOLUMNE_PATH + "/" + BasePath + "/" + FoPath
	case Exch == 2:
		newpath = Env.VOLUMNE_PATH + "/" + BasePath + "/" + EqPath
	case Exch == 3:
		// newpath = Env.VOLUMNE_PATH + "/chartdataservice" + "/nsecd1/" + FileName + "/" + MinuteFolder + "/"
	}
	_, err := os.Stat(newpath)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Println(err)
	}
	files, err := filepath.Glob(filepath.Join(newpath, "*"))
	if err != nil {
		return err
	}
	start := time.Now()
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
		fmt.Println("Deleting File", file)
	}
	fmt.Println("DeleteFiles, execution time " + time.Since(start).String())
	return nil
}
func GetDataFromStoreService(Exch int16, token int64) (map[string]*models.ChartMapModel, error) {
	data := Store[token][Exch].StoreMap
	if len(Store[token][Exch].StoreMap) == 0 {
		tk := strconv.Itoa(int(token))
		zerologs.Error().Msgf("No Data Found For Token = " + tk)
		return nil, errors.New("noData Found For Token = " + tk)
	}
	return data, nil
}
func GetRawdata(date string, Exch int16, token int64) (interface{}, error) {
	// UpDownServiceNew(Env.VOLUMNE_PATH+"/"+BasePath+"/"+RawPath+"/"+FoPath+"/"+current_date, 1)
	var newpath string
	var NewChartRespModel models.NewChartRespModel
	var response []models.ChartHLOCModel

	if Exch == 1 {
		newpath = Env.VOLUMNE_PATH + "/" + BasePath + "/" + RawPath + "/" + FoPath + "/" + date
	} else {
		newpath = Env.VOLUMNE_PATH + "/" + BasePath + "/" + RawPath + "/" + EqPath + "/" + date
	}
	folder, err := os.ReadDir(newpath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for file := range folder {
		data, err := os.ReadFile(newpath + "//" + folder[file].Name())
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		err = json.Unmarshal(data, &NewChartRespModel)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		for i := range NewChartRespModel.ChartHLOCModel {
			if token == int64(NewChartRespModel.ChartHLOCModel[i].NToken) {
				response = append(response, NewChartRespModel.ChartHLOCModel[i])
			}
		}
	}

	return response, nil
}
func ChangeFolderName() error {
	path := Env.VOLUMNE_PATH + "/" + BasePath + "/" + FoPath
	fomap := createSfullNameMap()
	fopath, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	contractMap := fomap.(map[string]models.ContractDetails)
	fmt.Println("folder name conversion Start" + time.Now().String())
	for folderName := range fopath {
		contract, ok := contractMap[fopath[folderName].Name()]
		if ok {
			datetime, _ := NewConvertDateTime1980(int64(contract.NExpiryDate))
			if datetime.Day() < 10 {
				err := os.Rename(path+"//"+fopath[folderName].Name(), path+"//"+contract.SFullName)
				if err != nil {
					fmt.Println(err)
					//return err
				}
			}
		}
	}
	fmt.Println("Succesfully Renamed folders" + time.Now().String())
	return nil
}

func ChangeFolderPermission(root interface{}) error {
	fmt.Println("ChangePermission Start")
	folder, err := os.ReadDir(root.(string))
	if err != nil {
		fmt.Println("Path Not Fount", root)
		return err
	}
	for i := range folder {
		if folder[i].IsDir() {
			info, _ := os.Stat(root.(string) + "//" + folder[i].Name())
			if info.Mode().Perm() != 0777 {
				err := os.Chmod(root.(string)+"//"+folder[i].Name(), 0777)
				if err != nil {
					Logger("ChangeFolderPermission() : Error While Changing Permission"+err.Error(), err, 2, liblogs.ZEROLOG|liblogs.PRINT)
					return err
				}
			}
		}
		subfolder, err := os.ReadDir(root.(string) + "//" + folder[i].Name())
		if err != nil {
			fmt.Println("Path Not Fount", root)
			return err
		}
		for i := range subfolder {
			if subfolder[i].IsDir() {
				info, _ := os.Stat(root.(string) + "//" + folder[i].Name() + "//" + subfolder[i].Name())
				if info.Mode().Perm() != 0777 {
					err := os.Chmod(root.(string)+"//"+folder[i].Name()+"//"+subfolder[i].Name(), 0777)
					if err != nil {
						Logger("ChangeFolderPermission() : Error While Changing Permission"+err.Error(), err, 2, liblogs.ZEROLOG|liblogs.PRINT)
						return err
					}
				}
			}
		}
	}
	fmt.Println("ChangePermission End")
	return nil
}

func CopyFolders(root interface{}) error {
	path := root.(string)
	fomap := createSfullNameMap()
	fopath, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	contractMap := fomap.(map[string]models.ContractDetails)
	fmt.Println("folder Copy  Start" + time.Now().String())
	for folderName := range fopath {
		contract, ok := contractMap[fopath[folderName].Name()]
		if ok {
			datetime, _ := NewConvertDateTime1980(int64(contract.NExpiryDate))
			if datetime.Day() < 10 {
				files, _ := Readdir(path + "//" + fopath[folderName].Name())
				if Exists(path + "//" + contract.SFullName) {
					for i := range files {
						tf := strings.Split(files[i], "//")
						tffoldername := tf[len(tf)-2]
						filename := tf[len(tf)-1]
						destinationfolder := path + "//" + contract.SFullName + "//" + tffoldername
						if _, err := os.Stat(destinationfolder); errors.Is(err, os.ErrNotExist) {
							err := os.MkdirAll(destinationfolder, 0777)
							if err != nil {
								Logger("Message: ", err, liblogs.Error, liblogs.PRINT)

							}
						}
						dstfile := destinationfolder + "//" + filename
						copyFile(files[i], dstfile)
					}
				} else {
					err := os.Rename(path+"//"+fopath[folderName].Name(), path+"//"+contract.SFullName)
					if err != nil {
						fmt.Println(err)
						return err
					}
				}
			}
		}
	}
	fmt.Println("folder Copy End" + time.Now().String())
	return nil
}

func ClearMemory() {
	fmt.Println("Alloc before M1 is made")
	printAlloc()
	for i := range Store {
		for j := range Store[i] {
			Store[i][j].ContractDetails = models.ContractDetails{}
			for keys := range Store[i][j].StoreMap {
				Store[i][j].StoreMap[keys].KeySlice = []int16{}
				for m := range Store[i][j].StoreMap[keys].CSVMap {
					delete(Store[i][j].StoreMap[keys].CSVMap, m)
				}
				// Store[i][j].StoreMap[keys].CSVMap = make(map[int16]*models.CSVModel)
			}
			// Store[i][j].StoreMap = make(map[string]*models.ChartMapModel)
		}
		Store[i] = [ExchSize]models.ChartModel{}
	}
	Store = nil

	for i := range ContractStore {
		for keys2 := range ContractStore[i].Contracts {
			delete(ContractStore[i].Contracts, keys2)
		}
		// ContractStore[i].Contracts = make(map[string]*models.ContractDetails)
	}

	runtime.GC()
	fmt.Println("Alloc before M1 is made")
	printAlloc()
}
