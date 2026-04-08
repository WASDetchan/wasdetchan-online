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
	"github.com/jackc/pgx/v5/pgtype"
)

type InvalidData struct{}

func (d InvalidData) Error() string {
	return "the qr-code or receipt string contained invalid data"
}

type OtherError struct{ code int }

func (d OtherError) Error() string {
	return fmt.Sprintf("An error has occured (code %v). Try again later.", d.code)
}

type ReceiptItem struct {
	Name     string `json:"name"`
	Nds      int64  `json:"nds"`
	Price    int64  `json:"price"`
	Quantity int64  `json:"quantity"`
	Sum      int64  `json:"sum"`
}

type Receipt struct {
	User               string        `json:"user"`
	RetailPlaceAddress string        `json:"retailPlaceAddress"`
	RetailPlace        string        `json:"retailPlace"`
	UserInn            string        `json:"userInn"`
	TicketDate         time.Time     `json:"ticketDate"`
	RequestNumber      int64         `json:"requestNumber"`
	ShiftNumber        int64         `json:"shiftNumber"`
	Operator           int64         `json:"operator"`
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

	var data responseData
	// err = json.NewDecoder(resp.Body).Decode(&data)
	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Printf("API response: %s", body)
		return nil, "", fmt.Errorf("error decoding response: %v", err)
	}

	if data.Code == 0 {
		return nil, "", InvalidData{}
	}

	if data.Code > 1 {
		return nil, "", OtherError{data.Code}
	}

	return data.Data.Json, data.Data.Html, nil
}

type receiptPostError struct {
	Code          int
	ErrorString   string
	AcceptedCount int
}

type receiptPostOk struct {
	Code    int      `json:"code"`
	IsNew   int      `json:"isNew"`
	Receipt *Receipt `json:"receipt"`
}

func HandlePostReceipt(c *gin.Context) {
	qrraw := c.PostForm("qrraw")

	if qrraw == "" {
		c.JSON(http.StatusBadRequest, receiptPostError{
			Code:        1,
			ErrorString: InvalidData{}.Error(),
		})
		return
	}

	data, _, err := GetReceiptFromString(qrraw)
	if err != nil {
		switch err := err.(type) {
		case OtherError:
			c.JSON(http.StatusInternalServerError, receiptPostError{
				Code:        err.code,
				ErrorString: err.Error(),
			})
		case InvalidData:
			c.JSON(http.StatusBadRequest, receiptPostError{
				Code:        1,
				ErrorString: err.Error(),
			})
		default:
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		return
	}

	user, _ := c.Get(auth.UserKey{})
	queries, _ := c.Get(repository.QueriesKey{})
	q := queries.(*repository.Queries)
	user_id := user.(*repository.User).ID

	_, err = q.CreateReceipt(context.Background(), repository.CreateReceiptParams{
		UserID: user_id,
		Fpd:    data.FiscalSign,
		Total:  data.TotalSum,
		Time: pgtype.Timestamp{
			Time:  data.TicketDate,
			Valid: true,
		},
		Place: data.RetailPlace,
	})

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, receiptPostOk{
		0,
		1,
		data,
	})
}
