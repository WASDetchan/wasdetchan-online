package receipt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/WASDetchan/wasdetchan-online/auth"
	"github.com/WASDetchan/wasdetchan-online/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type InvalidData struct{}

func (d InvalidData) Error() string {
	return "the qr-code or receipt string contained invalid data"
}

type OtherError struct {
	status_code int
	code        int
}

func (d OtherError) Error() string {
	return fmt.Sprintf("An error has occured (code %v, status %v). Try again later.", d.code, d.status_code)
}

type LocalTime struct {
	time.Time
}

func (ct *LocalTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}
	*ct = LocalTime{t}
	return nil
}

type ReceiptItem struct {
	Name     string  `json:"name"`
	Nds      int64   `json:"nds"`
	Price    int64   `json:"price"`
	Quantity float64 `json:"quantity"`
	Sum      int64   `json:"sum"`
}

type Receipt struct {
	User               string        `json:"user"`
	RetailPlaceAddress string        `json:"retailPlaceAddress"`
	RetailPlace        string        `json:"retailPlace"`
	UserInn            string        `json:"userInn"`
	Time               LocalTime     `json:"dateTime"`
	RequestNumber      int64         `json:"requestNumber"`
	ShiftNumber        int64         `json:"shiftNumber"`
	Operator           string        `json:"operator"`
	OperationType      int8          `json:"operationType"`
	Items              []ReceiptItem `json:"items"`
	Nds18              int64         `json:"nds18"`
	Nds                int64         `json:"nds"`
	Nds0               int64         `json:"nds0"`
	NdsNo              int64         `json:"ndsNo"`
	TotalSum           int64         `json:"totalSum"`
	CashTotalSum       int64         `json:"cashTotalSum"`
	EcashTotalSum      int64         `json:"ecashTotalSum"`
	TaxationType       int8          `json:"taxationType"`
	FiscalSign         int64         `json:"fiscalSign"`
}

type responseBase struct {
	Code int `json:"code"`
}

type responseData struct {
	Code  int  `json:"code"`
	First int8 `json:"first"`
	Data  *struct {
		Json *Receipt `json:"json"`
		Html string   `json:"html"`
	} `json:"data"`
}

type requestData struct {
	Qrraw string `json:"qrraw"`
	Token string `json:"token"`
}

const API_URL = "https://proverkacheka.com/api/v1/check/get"

func GetReceiptFromString(qrraw string) (*Receipt, string, error) {
	var buffer bytes.Buffer
	err := json.NewEncoder(&buffer).Encode(&requestData{qrraw, os.Getenv("PROVERKACHEKA_TOKEN")})
	if err != nil {
		return nil, "", fmt.Errorf("error encoding request: %v", err)
	}

	resp, err := http.Post(API_URL, "application/json", bytes.NewReader(buffer.Bytes()))
	if err != nil {
		return nil, "", fmt.Errorf("error requesting api: %v", err)
	}

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("error requesting api: %v", err)
	}

	var base responseBase
	err = json.Unmarshal(body, &base)
	if err != nil {
		log.Printf("API response: %s", body)
		return nil, "", fmt.Errorf("error decoding response: %v", err)
	}

	if base.Code == 0 {
		return nil, "", InvalidData{}
	}

	if base.Code > 1 {
		return nil, "", OtherError{resp.StatusCode, base.Code}
	}

	var data responseData
	// err = json.NewDecoder(resp.Body).Decode(&data)
	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Printf("API response: %s", body)
		return nil, "", fmt.Errorf("error decoding response: %v", err)
	}

	return data.Data.Json, data.Data.Html, nil
}

// 0 - OK
// 1..=6 - proverkacheka API error
// 7 - Missing data
// 8 - Duplicate
// API errors:
// 1 - чек некорректен,
// 2 - данные чека пока не получены,
// 3 - превышено кол-во запросов,
// 4 - ожидание перед повторным запросом,
// 5 - прочее (данные не получены)
type receiptPostResponse struct {
	Code    int      `json:"code"`
	Status  string   `json:"status"`
	Receipt *Receipt `json:"receipt"`
}

func HandlePostReceipt(c *gin.Context) {
	qrraw := c.PostForm("qrraw")

	log.Printf("[INFO] Processing receipt \"%v\"", qrraw)

	if qrraw == "" {
		c.JSON(http.StatusBadRequest, receiptPostResponse{
			Code:   7,
			Status: "The form is missing the \"qrraw\" field",
		})
		return
	}

	data, _, err := GetReceiptFromString(qrraw)
	if err != nil {
		switch err := err.(type) {
		case OtherError:
			c.JSON(http.StatusInternalServerError, receiptPostResponse{
				Code:   err.code,
				Status: err.Error(),
			})
		case InvalidData:
			c.JSON(http.StatusBadRequest, receiptPostResponse{
				Code:   1,
				Status: err.Error(),
			})
		default:
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	userId := auth.AssertAuth(c).ID
	q := repository.GetQueries(c)

	_, err = q.CreateReceipt(context.Background(), repository.CreateReceiptParams{
		UserID: userId,
		Fpd:    data.FiscalSign,
		Total:  data.TotalSum,
		Time: pgtype.Timestamp{
			Time:  data.Time.Time,
			Valid: true,
		},
		Place: data.RetailPlace,
	})

	if err != nil {
		if err.(*pgconn.PgError).Code == "23505" {
			c.JSON(http.StatusOK, receiptPostResponse{
				Code:    8,
				Status:  "The receipt was already registered",
				Receipt: data,
			})
			return
		}
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, receiptPostResponse{
		Code:    0,
		Status:  "Accepted",
		Receipt: data,
	})

}
