package service

import (
	"github.com/TecXLab/liblogs"
	"github.com/golobby/container/v3"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	CDS_Enqueue_Input_MSG_Received_Count prometheus.Counter
	CDS_Dequeue_Input_MSG_Received_Count prometheus.Counter

	CDS_Enqueue_RAW_File_Writing_Count prometheus.Counter
	CDS_Dequeue_RAW_File_Writing_Count prometheus.Counter

	CDS_Enqueue_NSEEQ_Count      prometheus.Counter
	CDS_Dequeue_NSEEQ_Count      prometheus.Counter
	CDS_Store_Update_NSEEQ_Count prometheus.Counter

	CDS_Enqueue_NSEFO_Count      prometheus.Counter
	CDS_Dequeue_NSEFO_Count      prometheus.Counter
	CDS_Store_Update_NSEFO_Count prometheus.Counter

	CDS_Enqueue_NSECD_Count      prometheus.Counter
	CDS_Dequeue_NSECD_Count      prometheus.Counter
	CDS_Store_Update_NSECD_Count prometheus.Counter

	CDS_NSEEQ_Total_Token_Count        prometheus.Counter
	CDS_NSEFO_Total_Token_Count        prometheus.Counter
	CDS_NSECD_Total_Token_Count        prometheus.Counter
	CDS_INDEX_Master_Total_Token_Count prometheus.Counter
)

func InitPromethues() {
	if Env.PROMETHEUS_COUNTER {
		var CDS_Enqueue_NSEEQ_Counter_Request *prometheus.Counter
		err := container.NamedResolve(&CDS_Enqueue_NSEEQ_Counter_Request, "CDS_Enqueue_NSEEQ_Counter_Request")
		if err != nil {
			Logger("CDS_Enqueue_NSEEQ_Counter_Request Lib Not Initialize", err, liblogs.Error, liblogs.ZEROLOG)
		}
		if CDS_Enqueue_NSEEQ_Counter_Request == nil {
			Logger("CDS_Enqueue_NSEEQ_Counter_Request counterConnection is NULL", err, liblogs.Error, liblogs.ZEROLOG)
		}
		CDS_Enqueue_NSEEQ_Count = *CDS_Enqueue_NSEEQ_Counter_Request

		var CDS_Dequeue_NSEEQ_Counter_Request *prometheus.Counter
		err = container.NamedResolve(&CDS_Dequeue_NSEEQ_Counter_Request, "CDS_Dequeue_NSEEQ_Counter_Request")
		if err != nil {
			Logger("CDS_Dequeue_NSEEQ_Counter_Request Lib Not Initialize", err, liblogs.Error, liblogs.ZEROLOG)
		}
		if CDS_Dequeue_NSEEQ_Counter_Request == nil {
			Logger("CDS_Dequeue_NSEEQ_Counter_Request counterConnection is NULL", err, liblogs.Error, liblogs.ZEROLOG)
		}
		CDS_Dequeue_NSEEQ_Count = *CDS_Dequeue_NSEEQ_Counter_Request

		var CDS_Store_Update_NSEEQ_Counter_Request *prometheus.Counter
		err = container.NamedResolve(&CDS_Store_Update_NSEEQ_Counter_Request, "CDS_Store_Update_NSEEQ_Counter_Request")
		if err != nil {
			Logger("CDS_Store_Update_NSEEQ_Count_Request Lib Not Initialize", err, liblogs.Error, liblogs.ZEROLOG)
		}
		if CDS_Store_Update_NSEEQ_Counter_Request == nil {
			Logger("CDS_Store_Update_NSEEQ_Count_Request counterConnection is NULL", err, liblogs.Error, liblogs.ZEROLOG)
		}
		CDS_Store_Update_NSEEQ_Count = *CDS_Store_Update_NSEEQ_Counter_Request

		var CDS_Enqueue_NSEFO_Counter_Request *prometheus.Counter
		err = container.NamedResolve(&CDS_Enqueue_NSEFO_Counter_Request, "CDS_Enqueue_NSEFO_Counter_Request")
		if err != nil {
			Logger("CDS_Enqueue_NSEFO_Counter_Request Lib Not Initialize", err, liblogs.Error, liblogs.ZEROLOG)
		}
		if CDS_Enqueue_NSEFO_Counter_Request == nil {
			Logger("CDS_Enqueue_NSEFO_Counter_Request counterConnection is NULL", err, liblogs.Error, liblogs.ZEROLOG)
		}
		CDS_Enqueue_NSEFO_Count = *CDS_Enqueue_NSEFO_Counter_Request

		var CDS_Dequeue_NSEFO_Counter_Request *prometheus.Counter
		err = container.NamedResolve(&CDS_Dequeue_NSEFO_Counter_Request, "CDS_Dequeue_NSEFO_Counter_Request")
		if err != nil {
			Logger("CDS_Dequeue_NSEFO_Counter_Request Lib Not Initialize", err, liblogs.Error, liblogs.ZEROLOG)
		}
		if CDS_Dequeue_NSEFO_Counter_Request == nil {
			Logger("CDS_Dequeue_NSEFO_Counter_Request counterConnection is NULL", err, liblogs.Error, liblogs.ZEROLOG)
		}
		CDS_Dequeue_NSEFO_Count = *CDS_Dequeue_NSEFO_Counter_Request

		var CDS_Store_Update_NSEFO_Counter_Request *prometheus.Counter
		err = container.NamedResolve(&CDS_Store_Update_NSEFO_Counter_Request, "CDS_Store_Update_NSEFO_Counter_Request")
		if err != nil {
			Logger("CDS_Store_Update_NSEFO_Counter_Request Lib Not Initialize", err, liblogs.Error, liblogs.ZEROLOG)
		}
		if CDS_Store_Update_NSEFO_Counter_Request == nil {
			Logger("CDS_Store_Update_NSEFO_Count_Request counterConnection is NULL", err, liblogs.Error, liblogs.ZEROLOG)
		}
		CDS_Store_Update_NSEFO_Count = *CDS_Store_Update_NSEFO_Counter_Request

		var CDS_Enqueue_Input_MSG_Received_Count_Request *prometheus.Counter
		err = container.NamedResolve(&CDS_Enqueue_Input_MSG_Received_Count_Request, "CDS_Enqueue_Input_MSG_Received_Count_Request")
		if err != nil {
			Logger("CDS_Enqueue_Input_MSG_Received_Count_Request Lib Not Initialize", err, liblogs.Error, liblogs.ZEROLOG)
		}
		if CDS_Enqueue_Input_MSG_Received_Count_Request == nil {
			Logger("CDS_Enqueue_Input_MSG_Received_Count_Request counterConnection is NULL", err, liblogs.Error, liblogs.ZEROLOG)
		}
		CDS_Enqueue_Input_MSG_Received_Count = *CDS_Enqueue_Input_MSG_Received_Count_Request

		var CDS_Dequeue_Input_MSG_Received_Count_Request *prometheus.Counter
		err = container.NamedResolve(&CDS_Dequeue_Input_MSG_Received_Count_Request, "CDS_Dequeue_Input_MSG_Received_Count_Request")
		if err != nil {
			Logger("CDS_Dequeue_Input_MSG_Received_Count_Request Lib Not Initialize", err, liblogs.Error, liblogs.ZEROLOG)
		}
		if CDS_Dequeue_Input_MSG_Received_Count_Request == nil {
			Logger("CDS_Dequeue_Input_MSG_Received_Count_Request counterConnection is NULL", err, liblogs.Error, liblogs.ZEROLOG)
		}
		CDS_Dequeue_Input_MSG_Received_Count = *CDS_Dequeue_Input_MSG_Received_Count_Request

		var CDS_Enqueue_RAW_File_Writing_Count_Request *prometheus.Counter
		err = container.NamedResolve(&CDS_Enqueue_RAW_File_Writing_Count_Request, "CDS_Enqueue_RAW_File_Writing_Count_Request")
		if err != nil {
			Logger("CDS_Enqueue_RAW_File_Writing_Count_Request Lib Not Initialize", err, liblogs.Error, liblogs.ZEROLOG)
		}
		if CDS_Enqueue_RAW_File_Writing_Count_Request == nil {
			Logger("CDS_Enqueue_RAW_File_Writing_Count_Request counterConnection is NULL", err, liblogs.Error, liblogs.ZEROLOG)
		}
		CDS_Enqueue_RAW_File_Writing_Count = *CDS_Enqueue_RAW_File_Writing_Count_Request

		var CDS_Dequeue_RAW_File_Writing_Count_Request *prometheus.Counter
		err = container.NamedResolve(&CDS_Dequeue_RAW_File_Writing_Count_Request, "CDS_Dequeue_RAW_File_Writing_Count_Request")
		if err != nil {
			Logger("CDS_Dequeue_RAW_File_Writing_Count_Request Lib Not Initialize", err, liblogs.Error, liblogs.ZEROLOG)
		}
		if CDS_Dequeue_RAW_File_Writing_Count_Request == nil {
			Logger("CDS_Dequeue_RAW_File_Writing_Count_Request counterConnection is NULL", err, liblogs.Error, liblogs.ZEROLOG)
		}
		CDS_Dequeue_RAW_File_Writing_Count = *CDS_Dequeue_RAW_File_Writing_Count_Request

		var CDS_NSEEQ_Total_Token_Count_Request *prometheus.Counter
		err = container.NamedResolve(&CDS_NSEEQ_Total_Token_Count_Request, "CDS_NSEEQ_Total_Token_Count_Request")
		if err != nil {
			Logger("CDS_NSEEQ_Total_Token_Count_Request Lib Not Initialize", err, liblogs.Error, liblogs.ZEROLOG)
		}
		if CDS_NSEEQ_Total_Token_Count_Request == nil {
			Logger("CDS_NSEEQ_Total_Token_Count_Request counterConnection is NULL", err, liblogs.Error, liblogs.ZEROLOG)
		}
		CDS_NSEEQ_Total_Token_Count = *CDS_NSEEQ_Total_Token_Count_Request

		var CDS_NSEFO_Total_Token_Count_Request *prometheus.Counter
		err = container.NamedResolve(&CDS_NSEFO_Total_Token_Count_Request, "CDS_NSEFO_Total_Token_Count_Request")
		if err != nil {
			Logger("CDS_NSEFO_Total_Token_Count_Request Lib Not Initialize", err, liblogs.Error, liblogs.ZEROLOG)
		}
		if CDS_NSEFO_Total_Token_Count_Request == nil {
			Logger("CDS_NSEFO_Total_Token_Count_Request counterConnection is NULL", err, liblogs.Error, liblogs.ZEROLOG)
		}
		CDS_NSEFO_Total_Token_Count = *CDS_NSEFO_Total_Token_Count_Request

		var CDS_INDEX_Master_Total_Token_Count_Request *prometheus.Counter
		err = container.NamedResolve(&CDS_INDEX_Master_Total_Token_Count_Request, "CDS_INDEX_Master_Total_Token_Count_Request")
		if err != nil {
			Logger("CDS_INDEX_Master_Total_Token_Count_Request Lib Not Initialize", err, liblogs.Error, liblogs.ZEROLOG)
		}
		if CDS_INDEX_Master_Total_Token_Count_Request == nil {
			Logger("CDS_INDEX_Master_Total_Token_Count_Request counterConnection is NULL", err, liblogs.Error, liblogs.ZEROLOG)
		}
		CDS_INDEX_Master_Total_Token_Count = *CDS_INDEX_Master_Total_Token_Count_Request

		var CDS_Enqueue_NSECD_Counter_Request *prometheus.Counter
		err = container.NamedResolve(&CDS_Enqueue_NSECD_Counter_Request, "CDS_Enqueue_NSECD_Counter_Request")
		if err != nil {
			Logger("CDS_Enqueue_NSECD_Counter_Request Lib Not Initialize", err, liblogs.Error, liblogs.ZEROLOG)
		}
		if CDS_Enqueue_NSECD_Counter_Request == nil {
			Logger("CDS_Enqueue_NSECD_Counter_Request counterConnection is NULL", err, liblogs.Error, liblogs.ZEROLOG)
		}
		CDS_Enqueue_NSECD_Count = *CDS_Enqueue_NSECD_Counter_Request

		var CDS_Dequeue_NSECD_Counter_Request *prometheus.Counter
		err = container.NamedResolve(&CDS_Dequeue_NSECD_Counter_Request, "CDS_Dequeue_NSECD_Counter_Request")
		if err != nil {
			Logger("CDS_Dequeue_NSECD_Counter_Request Lib Not Initialize", err, liblogs.Error, liblogs.ZEROLOG)
		}
		if CDS_Dequeue_NSECD_Counter_Request == nil {
			Logger("CDS_Dequeue_NSECD_Counter_Request counterConnection is NULL", err, liblogs.Error, liblogs.ZEROLOG)
		}
		CDS_Dequeue_NSECD_Count = *CDS_Dequeue_NSECD_Counter_Request

		var CDS_Store_Update_NSECD_Counter_Request *prometheus.Counter
		err = container.NamedResolve(&CDS_Store_Update_NSECD_Counter_Request, "CDS_Store_Update_NSECD_Counter_Request")
		if err != nil {
			Logger("CDS_Store_Update_NSECD_Counter_Request Lib Not Initialize", err, liblogs.Error, liblogs.ZEROLOG)
		}
		if CDS_Store_Update_NSECD_Counter_Request == nil {
			Logger("CDS_Store_Update_NSECD_Count_Request counterConnection is NULL", err, liblogs.Error, liblogs.ZEROLOG)
		}
		CDS_Store_Update_NSECD_Count = *CDS_Store_Update_NSECD_Counter_Request

		var CDS_NSECD_Total_Token_Count_Request *prometheus.Counter
		err = container.NamedResolve(&CDS_NSECD_Total_Token_Count_Request, "CDS_NSECD_Total_Token_Count_Request")
		if err != nil {
			Logger("CDS_NSECD_Total_Token_Count_Request Lib Not Initialize", err, liblogs.Error, liblogs.ZEROLOG)
		}
		if CDS_NSECD_Total_Token_Count_Request == nil {
			Logger("CDS_NSECD_Total_Token_Count_Request counterConnection is NULL", err, liblogs.Error, liblogs.ZEROLOG)
		}
		CDS_NSECD_Total_Token_Count = *CDS_NSECD_Total_Token_Count_Request
	}
}

func Add_CDS_Enqueue_NSEEQ_Count() {
	CDS_Enqueue_NSEEQ_Count.Inc()
}

func Add_CDS_Dequeue_NSEEQ_Count() {
	CDS_Dequeue_NSEEQ_Count.Inc()
}

func Add_CDS_Store_Update_NSEEQ_Count() {
	CDS_Enqueue_NSEEQ_Count.Inc()
}

func Add_CDS_Enqueue_NSEFO_Count() {
	CDS_Dequeue_NSEFO_Count.Inc()
}

func Add_CDS_Dequeue_NSEFO_Count() {
	CDS_Enqueue_NSEFO_Count.Inc()
}

func Add_CDS_Store_Update_NSEFO_Count() {
	CDS_Dequeue_NSEFO_Count.Inc()
}
func Add_CDS_Enqueue_NSECD_Count() {
	CDS_Dequeue_NSECD_Count.Inc()
}

func Add_CDS_Dequeue_NSECD_Count() {
	CDS_Enqueue_NSECD_Count.Inc()
}

func Add_CDS_Store_Update_NSECD_Count() {
	CDS_Dequeue_NSECD_Count.Inc()
}

func Add_CDS_Enqueue_Input_MSG_Received_Count() {
	CDS_Enqueue_Input_MSG_Received_Count.Inc()
}

func Add_CDS_Dequeue_Input_MSG_Received_Count() {
	CDS_Dequeue_Input_MSG_Received_Count.Inc()
}

func Add_CDS_Enqueue_RAW_File_Writing_Count() {
	CDS_Enqueue_RAW_File_Writing_Count.Inc()
}

func Add_CDS_Dequeue_RAW_File_Writing_Count() {
	CDS_Dequeue_RAW_File_Writing_Count.Inc()
}

func Add_CDS_NSEEQ_Total_Token_Count() {
	CDS_NSEEQ_Total_Token_Count.Inc()
}

func Add_CDS_NSECD_Total_Token_Count() {
	CDS_NSECD_Total_Token_Count.Inc()
}
func Add_CDS_NSEFO_Total_Token_Count() {
	CDS_NSEFO_Total_Token_Count.Inc()
}
func Add_CDS_INDEX_Master_Total_Token_Count() {
	CDS_INDEX_Master_Total_Token_Count.Inc()
}
