package handlers

import (
	"encoding/json"
	"github.com/SaCavid/simple-task/models"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	msg1                   = `{"state": "win", "amount": "10.15", "transactionId": "Same identification 1"}`
	msg2                   = `{"state": "win", "amount": "10.15", "transactionId": "Same identification 2"}`
	msg3                   = `{"state": "win", "amount": "10.15", "transactionId": "Same identification 3"}`
	errorUsedTransactionId = `{"state": "win", "amount": "10.15", "transactionId": "Same identification 1"}`
	errorNoTransactionId   = `{"state": "win", "amount": "10.15", "transactionId": ""}`
	errorNoState           = `{"state": "", "amount": "10.15", "transactionId": "Some identification"}`
	errorState             = `{"state": "error-state", "amount": "10.15", "transactionId": "Some identification 4"}`
	errorNullAmount        = `{"state": "win", "amount": "", "transactionId": "Some identification 5"}`
	winMsg                 = `{"state": "win", "amount": "27.99", "transactionId": "Some identification 10"}`
	loseMsg                = `{"state": "lose", "amount": "12.33", "transactionId": "Some identification 11"}`
	negativeMsg            = `{"state": "lose", "amount": "107.99", "transactionId": "Some identification 12"}`
)

type msg struct {
	state         string
	amount        float64
	transactionId string
}

func TestServer_Handler(t *testing.T) {
	h := &Server{
		TransactionIds: make(map[string]string, 0),
		UserBalances:   make(map[string]models.Balance, 0),
	}

	e := echo.New()

	h.badRequest(e)
	h.notAcceptableSourceType(e)
	h.noState(e)
	h.wrongState(e)

	h.noTransactionId(e)
	h.noAmount(e)
	h.sameTransactionId(e)
	h.notLogged(e)

	h.notRegistered(e)
	h.userWin(e)
	h.userLose(e)
	h.errorNegativeBalance(e)
}

func (h *Server) badRequest(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("message"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "not-source")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing bad request. Expected Code: 400. Got:", err.Error())
	}

}

func (h *Server) notAcceptableSourceType(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(msg1))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "not-source")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing not acceptable source type. Expected Code: 400. Got:", err.Error())
	}

}

func (h *Server) noState(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(errorNoState))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing null state. Expected Code: 400. Got:", err.Error())
	}

}

func (h *Server) wrongState(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(errorState))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing wrong state. Expected Code: 400. Got:", err.Error())
	}

}

func (h *Server) noTransactionId(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(errorNoTransactionId))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing null transaction Id. Expected Code: 400. Got:", err.Error())
	}

}

func (h *Server) noAmount(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(errorNullAmount))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing null amount. Expected Code: 400. Got:", err.Error())
	}

}

func (h *Server) sameTransactionId(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(errorUsedTransactionId))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing repeat transaction id. Expected Code: 406. Got:", err.Error())
	}
}

func (h *Server) notLogged(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(msg2))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing not logged. Expected Code: 403. Got:", err.Error())
	}
}

func (h *Server) notRegistered(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(msg3))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	req.Header.Set("Authorization", "not-registered-id")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing not registered. Expected Code: 400. Got:", err.Error())
	}
}

func (h *Server) userWin(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(winMsg))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	req.Header.Set("Authorization", "registered-id")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	b := models.Balance{
		Amount: 0,
		Saved:  true,
	}

	h.UserBalances["registered-id"] = b

	err := h.Handler(c)
	if err == nil {
		log.Println("Testing win state. Expected Code: 201. Got:", rec.Code, rec.Body.String())
	}
}

func (h *Server) userLose(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(loseMsg))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	req.Header.Set("Authorization", "registered-id")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err == nil {
		log.Println("Testing lose state. Expected Code: 201. Got:", rec.Code, rec.Body.String())
	}
}

func (h *Server) errorNegativeBalance(e *echo.Echo) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(negativeMsg))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	req.Header.Set("Authorization", "registered-id")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Handler(c)
	if err != nil {
		log.Println("Testing negative balance. Expected Code: 400. Got:", err.Error())
	}
}

func (h *Server) benchmarkRegistered(e *echo.Echo, msg string) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(msg))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	req.Header.Set("Authorization", "registered-id")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	b := models.Balance{
		Amount: 0,
		Saved:  true,
	}

	h.UserBalances["registered-id"] = b

	h.Handler(c)
}

// Benchmark 1000 messages
func BenchmarkServer_Handler(b *testing.B) {

	h := &Server{
		TransactionIds: make(map[string]string, 0),
		UserBalances:   make(map[string]models.Balance, 0),
	}

	e := echo.New()

	m := make([]msg, 0)

	err := json.Unmarshal([]byte(randomMsg), &m) // random generated 1000 messages
	if err != nil {
		log.Println(err)

	}

	// run the function N times
	N := len(m)
	for i := 0; i < N; i++ {
		d, _ := json.Marshal(m[i])
		h.benchmarkRegistered(e, string(d))
	}
}

var (
	//random generated 1000 messages
	// https://www.json-generator.com/
	//
	// generating script
	//
	// [
	//  '{{repeat(1000,1000)}}',
	//  {
	//    state:'{{random("win", "lose")}}',
	//    amount:'{{floating(1000, 4000, 2, "0,0.00")}}',
	//    transactionId:'{{guid()}}'
	//  }
	//]
	randomMsg = `[
  {
    "state": "lose",
    "amount": "3,148.19",
    "transactionId": "3b11e8c4-469e-4725-aa7b-7a69e27c02a2"
  },
  {
    "state": "win",
    "amount": "1,723.36",
    "transactionId": "391bf6da-8db2-42f9-9e3c-499161727543"
  },
  {
    "state": "lose",
    "amount": "2,212.79",
    "transactionId": "fe37191c-76a2-4e5b-b164-c0c6d82eeb6c"
  },
  {
    "state": "lose",
    "amount": "3,713.08",
    "transactionId": "4b6e46cf-1715-48bb-93e6-481776a0d65b"
  },
  {
    "state": "lose",
    "amount": "3,225.07",
    "transactionId": "62275fd5-bb43-4d8e-8676-d39025377868"
  },
  {
    "state": "win",
    "amount": "1,973.18",
    "transactionId": "cf35155c-e21c-42da-822c-38fbdaf00e0e"
  },
  {
    "state": "lose",
    "amount": "3,563.27",
    "transactionId": "d69569f2-a3d4-446a-b9be-2198dbb1af26"
  },
  {
    "state": "win",
    "amount": "1,751.92",
    "transactionId": "88281923-746b-4bd6-a48d-d0a6ed032e5d"
  },
  {
    "state": "win",
    "amount": "1,069.99",
    "transactionId": "04afa048-a92f-424e-87b4-1880cfb636f4"
  },
  {
    "state": "win",
    "amount": "2,116.21",
    "transactionId": "09a905b7-3e6c-4606-a6dc-22d9913b790b"
  },
  {
    "state": "win",
    "amount": "2,758.58",
    "transactionId": "715ebd55-187f-4461-b544-1f73fece7165"
  },
  {
    "state": "win",
    "amount": "3,710.50",
    "transactionId": "4f2e5770-2629-4d72-aaa7-c8510b185532"
  },
  {
    "state": "win",
    "amount": "1,349.68",
    "transactionId": "b4d259fc-b622-4fa8-bcc9-7ee6965ed8b3"
  },
  {
    "state": "lose",
    "amount": "3,491.43",
    "transactionId": "d44b0dda-4045-4303-b4e1-a8b4b467300c"
  },
  {
    "state": "lose",
    "amount": "2,755.66",
    "transactionId": "8b02f708-1c76-44c0-9004-582cd72ad001"
  },
  {
    "state": "lose",
    "amount": "2,112.77",
    "transactionId": "d022f977-b235-47bd-a5a3-bf716633a77d"
  },
  {
    "state": "win",
    "amount": "1,407.81",
    "transactionId": "80e50b4d-16de-4528-99ba-847d944aa26d"
  },
  {
    "state": "win",
    "amount": "3,363.47",
    "transactionId": "239b1ea5-95ae-4c14-bb73-123ff64cea3c"
  },
  {
    "state": "win",
    "amount": "3,995.49",
    "transactionId": "bd015d8f-7723-4139-b479-6df559b602ce"
  },
  {
    "state": "win",
    "amount": "1,732.94",
    "transactionId": "0667e49a-61bb-46e4-9d57-11ec25bc3ddb"
  },
  {
    "state": "win",
    "amount": "3,153.38",
    "transactionId": "4749c578-13cd-4038-9f8e-76657a1fe66c"
  },
  {
    "state": "win",
    "amount": "2,317.07",
    "transactionId": "07737b74-a850-454b-b91b-b722b225024a"
  },
  {
    "state": "win",
    "amount": "3,136.70",
    "transactionId": "ac30586e-ecfd-4a0b-be97-c4dd4947a43b"
  },
  {
    "state": "win",
    "amount": "1,941.03",
    "transactionId": "242e152c-4731-48be-8b9f-ddabf3bfeb42"
  },
  {
    "state": "lose",
    "amount": "3,010.83",
    "transactionId": "f2b594ac-8c83-4d76-9e71-5db04fbb2dda"
  },
  {
    "state": "win",
    "amount": "1,622.73",
    "transactionId": "fa7e2cea-ad06-411d-92cc-8a90578f9442"
  },
  {
    "state": "lose",
    "amount": "1,721.25",
    "transactionId": "3f2cbb3b-881a-4714-b2e1-cb0798ee492e"
  },
  {
    "state": "lose",
    "amount": "2,179.69",
    "transactionId": "b48c0beb-bb57-45a0-8092-0375414af0cc"
  },
  {
    "state": "win",
    "amount": "3,099.80",
    "transactionId": "7563a0d8-c918-42f5-8dea-ba21b875b6f2"
  },
  {
    "state": "lose",
    "amount": "1,232.15",
    "transactionId": "59877510-15b1-4797-be14-8f5a706b997e"
  },
  {
    "state": "lose",
    "amount": "1,013.74",
    "transactionId": "90555418-75b7-4aaf-aa58-6e8412440f62"
  },
  {
    "state": "win",
    "amount": "3,405.33",
    "transactionId": "31096f22-8090-4803-bcfe-b03d9486c923"
  },
  {
    "state": "win",
    "amount": "1,578.75",
    "transactionId": "e6051f9a-74bb-4cc3-9027-36d1ed4e1899"
  },
  {
    "state": "lose",
    "amount": "1,966.32",
    "transactionId": "412d0336-acd8-45ae-96ae-9da8c8a37f21"
  },
  {
    "state": "lose",
    "amount": "2,670.51",
    "transactionId": "510842d9-4a6c-4a16-9fb8-ef14d8c4e1fd"
  },
  {
    "state": "win",
    "amount": "2,722.11",
    "transactionId": "f74c03f4-8f4e-401d-8a0a-9ace28dac987"
  },
  {
    "state": "win",
    "amount": "1,001.01",
    "transactionId": "ab9aec5b-b31b-4211-a9d9-a4b0aa86546d"
  },
  {
    "state": "win",
    "amount": "1,719.95",
    "transactionId": "d5b35297-b99c-4199-857b-ae0be2ae8680"
  },
  {
    "state": "lose",
    "amount": "3,791.80",
    "transactionId": "8f438b22-90f8-49c0-ad84-15f5b6479cb7"
  },
  {
    "state": "win",
    "amount": "1,215.52",
    "transactionId": "d13c62f7-3dd0-4f86-99e1-3e0ddd3d5f96"
  },
  {
    "state": "lose",
    "amount": "3,368.19",
    "transactionId": "03d6f837-8a5a-4b84-9e88-117a95b143d0"
  },
  {
    "state": "win",
    "amount": "2,055.36",
    "transactionId": "85c664b5-057e-4494-b41c-47ae24af0593"
  },
  {
    "state": "win",
    "amount": "3,990.26",
    "transactionId": "0501d474-59f5-47c2-91ac-243d013cf18c"
  },
  {
    "state": "win",
    "amount": "2,171.97",
    "transactionId": "fef3fac6-b55f-49e2-94ed-63b83bc834d0"
  },
  {
    "state": "lose",
    "amount": "2,843.38",
    "transactionId": "1da4b5cd-872c-44fc-bfe0-a0dd942dbb2b"
  },
  {
    "state": "win",
    "amount": "2,370.62",
    "transactionId": "8b767547-ae9c-47a6-b2b5-1ec812efbdca"
  },
  {
    "state": "lose",
    "amount": "3,609.64",
    "transactionId": "bf853514-7a22-40d2-8a77-eb227f32e97f"
  },
  {
    "state": "lose",
    "amount": "1,092.89",
    "transactionId": "1105c87a-a41c-49f7-8cf8-81e4509de12b"
  },
  {
    "state": "lose",
    "amount": "3,938.66",
    "transactionId": "a6089c00-f3cf-45e4-8711-d05b91e686a2"
  },
  {
    "state": "win",
    "amount": "1,808.77",
    "transactionId": "2187a89a-fd92-414f-8fd1-20feba19cf4e"
  },
  {
    "state": "lose",
    "amount": "1,530.43",
    "transactionId": "785e15c2-7742-4256-8a6d-5762b26f6ec0"
  },
  {
    "state": "lose",
    "amount": "1,607.85",
    "transactionId": "f8f6a2a2-4b72-4cbf-bb52-eea4a16ef613"
  },
  {
    "state": "win",
    "amount": "2,481.22",
    "transactionId": "127e373b-8513-4055-b960-c781482b23d6"
  },
  {
    "state": "lose",
    "amount": "2,145.31",
    "transactionId": "c0a83516-c8f4-4fe0-a500-cc55e8f7d624"
  },
  {
    "state": "lose",
    "amount": "2,408.33",
    "transactionId": "db52b25a-ba08-4f3e-a149-bfe2ef40793c"
  },
  {
    "state": "win",
    "amount": "1,459.49",
    "transactionId": "1765b251-059c-475f-9444-966ee8d887b7"
  },
  {
    "state": "lose",
    "amount": "1,191.40",
    "transactionId": "a5ca9af6-7563-4644-974c-9b222b3ee30c"
  },
  {
    "state": "lose",
    "amount": "3,790.04",
    "transactionId": "97118bc5-7919-40b4-b268-df94876ca5ed"
  },
  {
    "state": "win",
    "amount": "2,413.30",
    "transactionId": "733d3252-83f0-4b8e-816a-9a52c8e850ce"
  },
  {
    "state": "lose",
    "amount": "2,790.03",
    "transactionId": "84a869e6-6fb7-416b-825e-d6a9bc23ea85"
  },
  {
    "state": "lose",
    "amount": "2,515.22",
    "transactionId": "1dc06b14-2158-4906-b59b-0c96dbbff604"
  },
  {
    "state": "lose",
    "amount": "1,204.82",
    "transactionId": "807292b8-5714-4652-b2ed-4ccc529a402e"
  },
  {
    "state": "lose",
    "amount": "1,049.31",
    "transactionId": "8e8afbe8-4cfe-49eb-af50-90f92ee98c29"
  },
  {
    "state": "win",
    "amount": "3,514.77",
    "transactionId": "ad4b4967-b7f6-4fb0-910a-076de0ff552d"
  },
  {
    "state": "lose",
    "amount": "3,829.55",
    "transactionId": "df18af07-6cf7-4a0f-a12d-7b6f3d2e4481"
  },
  {
    "state": "win",
    "amount": "3,225.45",
    "transactionId": "34ff2d90-8c59-42fd-9d34-b69ba462e945"
  },
  {
    "state": "lose",
    "amount": "1,508.27",
    "transactionId": "85fc7172-ba43-4db0-984f-9cdc39e3268b"
  },
  {
    "state": "lose",
    "amount": "2,513.82",
    "transactionId": "73b872ec-9ef4-4ab6-8806-4df5f70ed609"
  },
  {
    "state": "win",
    "amount": "2,643.85",
    "transactionId": "dc60eed8-6050-444b-b1f7-f023a18a1482"
  },
  {
    "state": "lose",
    "amount": "1,404.68",
    "transactionId": "77ec2ecc-9fa6-4f97-a259-57c96d48438d"
  },
  {
    "state": "lose",
    "amount": "1,100.17",
    "transactionId": "ce390ec6-64b2-407e-b5d6-57ab6773fa36"
  },
  {
    "state": "win",
    "amount": "3,529.86",
    "transactionId": "f03cb919-e497-4d96-a76a-0fafd63af08f"
  },
  {
    "state": "win",
    "amount": "2,320.20",
    "transactionId": "368b299f-d739-418b-9eba-e24b32fff465"
  },
  {
    "state": "lose",
    "amount": "2,929.56",
    "transactionId": "3a9b910f-2357-4e36-ab10-98b545687a32"
  },
  {
    "state": "lose",
    "amount": "1,421.99",
    "transactionId": "441d1c8b-b78a-42a1-85b9-376b91413142"
  },
  {
    "state": "win",
    "amount": "2,284.56",
    "transactionId": "01d07c54-2a6d-40b8-95b3-d46e41ee5e62"
  },
  {
    "state": "win",
    "amount": "3,373.66",
    "transactionId": "e30550af-b2dc-4b2f-a39a-6271691a80d7"
  },
  {
    "state": "lose",
    "amount": "3,504.58",
    "transactionId": "5ee7735d-a5ec-43b8-9baf-7bc389d3ffb8"
  },
  {
    "state": "win",
    "amount": "2,844.09",
    "transactionId": "f0d9bc1b-0b42-4539-a7ee-bf297e8d1544"
  },
  {
    "state": "win",
    "amount": "1,975.64",
    "transactionId": "39fc6881-74de-4003-8f41-0f7b562a2cb7"
  },
  {
    "state": "win",
    "amount": "2,930.46",
    "transactionId": "c423a4df-edd3-45cf-8b71-2d0150aaac0b"
  },
  {
    "state": "lose",
    "amount": "3,954.81",
    "transactionId": "a1c1c428-99f5-4338-85ce-4b82527e2c73"
  },
  {
    "state": "win",
    "amount": "1,542.13",
    "transactionId": "638ce1a1-618b-4788-8d3d-4c9c9baec3f5"
  },
  {
    "state": "win",
    "amount": "3,676.51",
    "transactionId": "656bcc65-eb19-4362-9655-fe7be99c1a50"
  },
  {
    "state": "win",
    "amount": "1,756.94",
    "transactionId": "da372ccf-b5a5-4669-8477-c98be4a2d9da"
  },
  {
    "state": "win",
    "amount": "2,783.69",
    "transactionId": "f089c11d-cd30-4f92-bf63-0894f775abb1"
  },
  {
    "state": "win",
    "amount": "1,601.16",
    "transactionId": "d6eb3120-accc-4555-a3db-d32c56bfdc83"
  },
  {
    "state": "lose",
    "amount": "2,236.39",
    "transactionId": "bc0fe566-8f2f-4224-8d46-1e7f8f79bffe"
  },
  {
    "state": "win",
    "amount": "2,639.77",
    "transactionId": "c7c14d58-938e-48b8-9ec4-b25dbe7251e8"
  },
  {
    "state": "lose",
    "amount": "3,481.44",
    "transactionId": "99b9df8c-d08c-4926-9021-114c8bb49c55"
  },
  {
    "state": "lose",
    "amount": "2,312.61",
    "transactionId": "3b664bb8-d844-47a7-bea9-b0c102a02c2f"
  },
  {
    "state": "win",
    "amount": "3,415.15",
    "transactionId": "1b46a5bd-61da-482f-828c-e4484638fa01"
  },
  {
    "state": "lose",
    "amount": "2,934.62",
    "transactionId": "39f59370-629f-430a-91a8-be78ef53b329"
  },
  {
    "state": "win",
    "amount": "3,444.55",
    "transactionId": "691a46a8-2295-457f-8dd4-23ca7bcef77f"
  },
  {
    "state": "win",
    "amount": "2,086.48",
    "transactionId": "406118cb-b6ac-49b0-82b3-05484a29fda3"
  },
  {
    "state": "lose",
    "amount": "3,464.03",
    "transactionId": "f0387936-c677-46df-8750-9c8875a64e6e"
  },
  {
    "state": "win",
    "amount": "2,570.85",
    "transactionId": "bd052740-694f-4ae1-aeeb-54b715ed2270"
  },
  {
    "state": "win",
    "amount": "3,969.85",
    "transactionId": "aeb46f5a-f4a7-4e64-9742-a1f6fadd7dca"
  },
  {
    "state": "lose",
    "amount": "2,992.78",
    "transactionId": "55c8a68e-8850-4c74-9666-5420ab461901"
  },
  {
    "state": "lose",
    "amount": "2,218.94",
    "transactionId": "147f6227-af8b-4e51-b5ef-e9ca9bdff3bf"
  },
  {
    "state": "lose",
    "amount": "2,660.39",
    "transactionId": "3386aafc-fd30-47ca-b533-e6aa0e1d7b38"
  },
  {
    "state": "lose",
    "amount": "3,674.76",
    "transactionId": "df5efdd3-a381-4579-8bba-b95300068160"
  },
  {
    "state": "lose",
    "amount": "1,316.58",
    "transactionId": "6dbd68e1-ac42-4989-bd52-5778247330f4"
  },
  {
    "state": "win",
    "amount": "2,838.02",
    "transactionId": "83cc6c03-9b10-4c09-ae89-5bc423fc9df1"
  },
  {
    "state": "lose",
    "amount": "2,952.21",
    "transactionId": "86996207-9141-423a-9778-ab197d25669d"
  },
  {
    "state": "win",
    "amount": "1,787.48",
    "transactionId": "39aaf42c-e95c-44c5-949d-89110d8fbbe8"
  },
  {
    "state": "win",
    "amount": "2,792.35",
    "transactionId": "e2c7a799-fb87-466b-85f7-2a42d4d6e1cb"
  },
  {
    "state": "lose",
    "amount": "1,252.78",
    "transactionId": "b8fb699f-a846-49f2-b02a-39df2d7b3ad4"
  },
  {
    "state": "win",
    "amount": "1,304.15",
    "transactionId": "644c9262-346a-4ef6-868e-c78005537702"
  },
  {
    "state": "win",
    "amount": "3,507.29",
    "transactionId": "08d76121-e73d-452b-a205-9060fab986c6"
  },
  {
    "state": "lose",
    "amount": "1,404.39",
    "transactionId": "5dae538e-52fa-49ff-983d-984c2fe6e673"
  },
  {
    "state": "win",
    "amount": "1,714.62",
    "transactionId": "335fba2c-fa75-49dc-be04-2bfde7c92239"
  },
  {
    "state": "lose",
    "amount": "1,312.55",
    "transactionId": "287774c7-2fd3-47e8-b822-cf39ede16f1e"
  },
  {
    "state": "lose",
    "amount": "3,516.26",
    "transactionId": "48e96418-45e8-49d2-83ff-859f78b3eee3"
  },
  {
    "state": "win",
    "amount": "3,401.11",
    "transactionId": "73c2505b-95c4-463c-9e79-c0868b5d2586"
  },
  {
    "state": "win",
    "amount": "1,231.69",
    "transactionId": "b8b9e771-590a-4135-b928-805e6a181aa9"
  },
  {
    "state": "win",
    "amount": "2,407.34",
    "transactionId": "174d3d57-5850-4122-9d6d-b344e9de1bf6"
  },
  {
    "state": "win",
    "amount": "2,399.28",
    "transactionId": "2ea6d078-1347-4ae8-afdc-997bba8c6d20"
  },
  {
    "state": "win",
    "amount": "2,407.41",
    "transactionId": "7d217018-e8d8-4e00-848d-3d468d584c25"
  },
  {
    "state": "win",
    "amount": "1,123.07",
    "transactionId": "f261b4c0-9103-4d71-acbd-055d08d5473b"
  },
  {
    "state": "win",
    "amount": "2,324.32",
    "transactionId": "5dc8b975-64d8-4038-883b-b35d8f0187e5"
  },
  {
    "state": "win",
    "amount": "2,246.34",
    "transactionId": "8889ba22-f458-40ae-85f1-6655006cf5b7"
  },
  {
    "state": "lose",
    "amount": "3,758.33",
    "transactionId": "7d29ab5f-d02c-4277-8190-006c41c939d4"
  },
  {
    "state": "win",
    "amount": "1,113.14",
    "transactionId": "662ac09f-5926-4776-b2a2-53d33f22163d"
  },
  {
    "state": "lose",
    "amount": "3,912.12",
    "transactionId": "67604ef0-75b1-4156-96e7-b19beb484cef"
  },
  {
    "state": "lose",
    "amount": "2,572.25",
    "transactionId": "bb69b59d-57ce-43ff-923a-84daac211572"
  },
  {
    "state": "win",
    "amount": "2,886.05",
    "transactionId": "eed7cced-357d-4b29-b1a5-8afc6f1c2a00"
  },
  {
    "state": "lose",
    "amount": "1,170.15",
    "transactionId": "4c7c3aac-7ba1-459e-ae94-3c5b67ac40ec"
  },
  {
    "state": "lose",
    "amount": "2,152.07",
    "transactionId": "4c7d2450-8962-48e0-b015-dd02da2c19fe"
  },
  {
    "state": "win",
    "amount": "3,258.45",
    "transactionId": "7d9dd45c-02a6-434b-918b-f6e5d78cc518"
  },
  {
    "state": "lose",
    "amount": "1,859.82",
    "transactionId": "cc918b3b-d471-457c-98e3-9c05c927e335"
  },
  {
    "state": "win",
    "amount": "2,676.20",
    "transactionId": "dbb39afa-aedb-4377-a029-d0bdadcdaf90"
  },
  {
    "state": "lose",
    "amount": "2,192.06",
    "transactionId": "18c16a5f-6964-4d00-9416-e9b6a12c5d71"
  },
  {
    "state": "win",
    "amount": "1,704.63",
    "transactionId": "cec5f3a1-a491-4204-87d4-914a4d22fbb8"
  },
  {
    "state": "win",
    "amount": "2,672.65",
    "transactionId": "8c117cd4-d194-41e4-aef8-3d52f527b076"
  },
  {
    "state": "win",
    "amount": "3,404.25",
    "transactionId": "3987339f-37c6-4497-b761-2d814c332ac7"
  },
  {
    "state": "win",
    "amount": "1,496.61",
    "transactionId": "5c37f0c9-35cd-4a21-80be-ffd3a5e0621b"
  },
  {
    "state": "lose",
    "amount": "2,936.03",
    "transactionId": "bf6203e2-ee31-4e8d-8141-aa6df35be012"
  },
  {
    "state": "win",
    "amount": "1,474.07",
    "transactionId": "426c4d00-5678-4a14-92eb-2a68feaaba93"
  },
  {
    "state": "lose",
    "amount": "3,084.02",
    "transactionId": "10de65a0-d0e3-4b3e-bf15-93865cc0e9b9"
  },
  {
    "state": "lose",
    "amount": "1,361.79",
    "transactionId": "0e76ddea-e659-449a-bd04-06678db737fc"
  },
  {
    "state": "lose",
    "amount": "1,434.73",
    "transactionId": "2570d95c-4957-4c30-aa4b-53274190a8d2"
  },
  {
    "state": "lose",
    "amount": "1,903.40",
    "transactionId": "bc3f0b14-d0ee-4266-9551-75580936dac5"
  },
  {
    "state": "lose",
    "amount": "1,790.91",
    "transactionId": "ab48ebd9-0473-43dc-914d-b1ffc69fb479"
  },
  {
    "state": "lose",
    "amount": "1,576.48",
    "transactionId": "55bd914f-f23c-4f76-a347-5eea8f495cde"
  },
  {
    "state": "win",
    "amount": "3,046.32",
    "transactionId": "c215ef3c-7c20-4ec7-8fc1-32d3380ac705"
  },
  {
    "state": "win",
    "amount": "1,904.30",
    "transactionId": "cffa048a-136e-446c-b0f7-6cdf294965b0"
  },
  {
    "state": "lose",
    "amount": "2,698.52",
    "transactionId": "c5a06b4a-9bb8-4cad-8389-9385d0d922e7"
  },
  {
    "state": "win",
    "amount": "1,044.09",
    "transactionId": "89733668-6935-4e78-82c7-d973b63616a4"
  },
  {
    "state": "lose",
    "amount": "2,046.86",
    "transactionId": "063eb6ca-0d74-4fb6-93bc-2ab3e6991fd0"
  },
  {
    "state": "lose",
    "amount": "3,974.83",
    "transactionId": "5a4695d7-d034-4d07-b28d-098eee2e749f"
  },
  {
    "state": "win",
    "amount": "2,607.09",
    "transactionId": "deb3c96f-2ebf-4d0a-9a58-39d614d74450"
  },
  {
    "state": "lose",
    "amount": "1,583.98",
    "transactionId": "bb35f4a0-3c49-4e02-9203-0e3e34f310a4"
  },
  {
    "state": "win",
    "amount": "3,215.16",
    "transactionId": "15cf325a-06d8-4ccd-8e0e-b78a98881b82"
  },
  {
    "state": "lose",
    "amount": "3,941.12",
    "transactionId": "bd358468-5b61-4b73-8fcb-c3f871262998"
  },
  {
    "state": "win",
    "amount": "2,978.38",
    "transactionId": "66925bd1-6a8b-42f6-af0e-f02c04dcd493"
  },
  {
    "state": "win",
    "amount": "3,426.21",
    "transactionId": "1d34e7bb-676e-4dc4-9ef3-16d449378260"
  },
  {
    "state": "win",
    "amount": "1,524.98",
    "transactionId": "0af9a539-2e11-4838-839c-36b32e8999a9"
  },
  {
    "state": "win",
    "amount": "1,485.64",
    "transactionId": "12e83029-63a6-4420-a3b4-7e8761239cbb"
  },
  {
    "state": "win",
    "amount": "1,191.27",
    "transactionId": "bf643549-a7b9-4828-8ee3-feb77f723e69"
  },
  {
    "state": "lose",
    "amount": "3,001.36",
    "transactionId": "1bc47d0f-db46-4f6e-a729-e6e0cd65e09d"
  },
  {
    "state": "win",
    "amount": "3,807.11",
    "transactionId": "82c6aec2-15af-40b1-b7b3-b77b19e4ae7a"
  },
  {
    "state": "lose",
    "amount": "1,686.76",
    "transactionId": "5d6e2c2f-4a80-4cc2-95e5-5f0a56b8f7f7"
  },
  {
    "state": "win",
    "amount": "2,774.12",
    "transactionId": "01f6a914-91c8-44d9-8190-3256efc3691a"
  },
  {
    "state": "lose",
    "amount": "2,136.85",
    "transactionId": "5693a882-da42-4ccf-872c-ea53d79d6718"
  },
  {
    "state": "win",
    "amount": "1,858.07",
    "transactionId": "746b41e0-1f45-4efa-9a89-491ec485b9f3"
  },
  {
    "state": "lose",
    "amount": "1,913.76",
    "transactionId": "39b924fe-5496-44f8-aeaa-29fbda43613c"
  },
  {
    "state": "win",
    "amount": "2,561.05",
    "transactionId": "468fb340-bdbb-42f0-81ad-1b394a203a01"
  },
  {
    "state": "win",
    "amount": "1,508.56",
    "transactionId": "ff858763-5237-4154-82c6-616d19dbca65"
  },
  {
    "state": "lose",
    "amount": "3,173.10",
    "transactionId": "48989ba5-d675-46fa-bac7-3f5758d51835"
  },
  {
    "state": "win",
    "amount": "3,692.50",
    "transactionId": "76cdd42b-1e0a-4863-b8e6-cda97c13adc6"
  },
  {
    "state": "win",
    "amount": "1,103.57",
    "transactionId": "62600a4f-80e8-49cf-9ca7-698b725fea53"
  },
  {
    "state": "win",
    "amount": "2,691.53",
    "transactionId": "16f1a4b5-9847-45cd-8c7a-b856c0a3abdf"
  },
  {
    "state": "win",
    "amount": "1,341.87",
    "transactionId": "3e470d54-47fa-488f-a026-fc5f1ce36066"
  },
  {
    "state": "win",
    "amount": "3,759.69",
    "transactionId": "031f0dd7-cc2e-47b6-8d54-d0140cbc3f5d"
  },
  {
    "state": "win",
    "amount": "3,652.91",
    "transactionId": "3f333b84-d0da-4264-8393-281c5099ee96"
  },
  {
    "state": "win",
    "amount": "3,342.29",
    "transactionId": "e29c4d65-0b68-4844-88ee-a3002b6f7550"
  },
  {
    "state": "win",
    "amount": "3,143.30",
    "transactionId": "22c68eac-e5fa-4b26-ba5c-f3e6f4ae141c"
  },
  {
    "state": "win",
    "amount": "1,031.65",
    "transactionId": "da8395ef-2659-469c-922c-5827b9fcb2ce"
  },
  {
    "state": "win",
    "amount": "2,928.14",
    "transactionId": "09303cf6-5588-4e38-8a53-772f46cd7ad5"
  },
  {
    "state": "lose",
    "amount": "1,459.19",
    "transactionId": "4d724515-2f42-4fb8-8560-9a89c07d1cbc"
  },
  {
    "state": "win",
    "amount": "1,851.30",
    "transactionId": "948ea555-4f60-451c-a003-a6f2645aa3cb"
  },
  {
    "state": "lose",
    "amount": "2,628.90",
    "transactionId": "f4e24b3f-8202-4d16-9648-9df01cfbd4c9"
  },
  {
    "state": "win",
    "amount": "2,186.05",
    "transactionId": "03d48846-7741-496b-bd55-8e7761cc07c0"
  },
  {
    "state": "win",
    "amount": "3,668.35",
    "transactionId": "85b171b6-3211-4216-a5af-07c9be141c38"
  },
  {
    "state": "win",
    "amount": "1,088.75",
    "transactionId": "1fe65ead-a10c-41a4-91cf-19e18af2ce3a"
  },
  {
    "state": "lose",
    "amount": "3,390.00",
    "transactionId": "ffaceacd-84be-4761-87c2-0d68b4dbe897"
  },
  {
    "state": "lose",
    "amount": "1,498.91",
    "transactionId": "f9704c7e-e15d-4033-b64f-ecd05ecea660"
  },
  {
    "state": "win",
    "amount": "3,357.27",
    "transactionId": "b07c0f46-b6e0-4783-89dd-351b2d7e0d0f"
  },
  {
    "state": "win",
    "amount": "2,919.71",
    "transactionId": "f370f207-a9ea-4987-a7aa-0ce7638548d7"
  },
  {
    "state": "lose",
    "amount": "2,215.14",
    "transactionId": "b0bc24b4-3856-4436-a7c3-386eac032891"
  },
  {
    "state": "lose",
    "amount": "1,481.22",
    "transactionId": "9d778b83-28d4-4673-ae33-1b24f31a6007"
  },
  {
    "state": "lose",
    "amount": "1,534.59",
    "transactionId": "ce9fa24c-0048-485f-a793-13b46c71b6ed"
  },
  {
    "state": "lose",
    "amount": "2,652.47",
    "transactionId": "009d903d-660f-461a-907c-9036ec0e8520"
  },
  {
    "state": "win",
    "amount": "3,607.56",
    "transactionId": "456af84b-21f8-4510-b63c-f361a471d501"
  },
  {
    "state": "lose",
    "amount": "3,957.67",
    "transactionId": "7980a5af-9e00-4896-b7af-c63c26352570"
  },
  {
    "state": "lose",
    "amount": "2,064.56",
    "transactionId": "a0af3437-1428-42bd-be9e-0849d657580f"
  },
  {
    "state": "win",
    "amount": "3,494.96",
    "transactionId": "d9dafa6f-bef4-4c43-a068-254dcf2c482c"
  },
  {
    "state": "lose",
    "amount": "2,526.04",
    "transactionId": "ddb9118a-58b9-4c6b-a810-6fce29fa6e32"
  },
  {
    "state": "win",
    "amount": "3,696.83",
    "transactionId": "897c39f8-b8f2-43be-b50a-a7123ce13ef2"
  },
  {
    "state": "lose",
    "amount": "2,859.72",
    "transactionId": "efbad8ce-943c-4127-8ce8-bb277ef60175"
  },
  {
    "state": "win",
    "amount": "2,290.25",
    "transactionId": "87d5957a-64ad-41b6-b59b-f31406e48df8"
  },
  {
    "state": "win",
    "amount": "2,197.28",
    "transactionId": "bf637cbb-d586-4c4f-b20a-d859b4818d68"
  },
  {
    "state": "win",
    "amount": "1,094.61",
    "transactionId": "d6319364-a5c5-4bed-ba22-b06f27bd2b53"
  },
  {
    "state": "lose",
    "amount": "1,349.24",
    "transactionId": "7f2519bc-d7e7-4873-a135-9396a1bb959b"
  },
  {
    "state": "lose",
    "amount": "2,525.50",
    "transactionId": "5c09a60e-ab33-440e-8056-3f4d2ccba322"
  },
  {
    "state": "lose",
    "amount": "3,278.45",
    "transactionId": "be5fc622-10e8-4f8e-badb-f875db6f8ef0"
  },
  {
    "state": "lose",
    "amount": "3,991.85",
    "transactionId": "eafc6047-c69f-4471-8ec2-9f61a4657d76"
  },
  {
    "state": "lose",
    "amount": "2,941.45",
    "transactionId": "66a8a071-b84a-42a6-a2b8-aab20b92ff84"
  },
  {
    "state": "win",
    "amount": "1,708.48",
    "transactionId": "d6ad9ce8-9bdf-4c5e-97fa-2f6fbf45ca44"
  },
  {
    "state": "win",
    "amount": "2,732.75",
    "transactionId": "c6eb3eb7-c8b8-402e-92f1-59c9f4132447"
  },
  {
    "state": "win",
    "amount": "2,344.86",
    "transactionId": "507578cc-8aaf-4180-b306-2f38d8cf03b6"
  },
  {
    "state": "lose",
    "amount": "3,662.18",
    "transactionId": "4b624dcf-2237-445b-b75a-e7f55f2685f9"
  },
  {
    "state": "win",
    "amount": "3,922.55",
    "transactionId": "2bcbe630-c3fd-440e-ad85-a584eb871fa7"
  },
  {
    "state": "win",
    "amount": "2,040.79",
    "transactionId": "74d9e0f9-d484-4f09-8b06-807217384fd6"
  },
  {
    "state": "win",
    "amount": "1,762.69",
    "transactionId": "dbdd6452-d318-45cb-9bc6-efa512d5015a"
  },
  {
    "state": "lose",
    "amount": "3,301.85",
    "transactionId": "4f5df840-c5b4-4cf1-81cb-9bd6c5e7961a"
  },
  {
    "state": "win",
    "amount": "2,517.29",
    "transactionId": "a6a34519-f82b-44c2-bf10-40b725f75893"
  },
  {
    "state": "win",
    "amount": "3,569.36",
    "transactionId": "62819d63-f814-46cc-9dcc-e3d134bfdfda"
  },
  {
    "state": "win",
    "amount": "1,583.86",
    "transactionId": "f97cd405-342d-4f41-93ac-6b934cce8a03"
  },
  {
    "state": "lose",
    "amount": "2,697.42",
    "transactionId": "cb2758aa-f77d-4609-b5a4-31ba66c41bf2"
  },
  {
    "state": "win",
    "amount": "2,948.66",
    "transactionId": "451abebb-1e53-498a-b773-8f8dc2586843"
  },
  {
    "state": "lose",
    "amount": "1,904.35",
    "transactionId": "c3398331-d5e4-4cf5-82eb-a06af259cb44"
  },
  {
    "state": "lose",
    "amount": "2,750.49",
    "transactionId": "830046cd-17f2-4f63-bc44-fa6ba6c0a52f"
  },
  {
    "state": "win",
    "amount": "2,710.79",
    "transactionId": "8795085b-5b0d-4ecb-bca8-599e57932318"
  },
  {
    "state": "win",
    "amount": "3,061.35",
    "transactionId": "14380ef1-247a-4668-b10d-779f14f82be0"
  },
  {
    "state": "lose",
    "amount": "1,063.42",
    "transactionId": "6c4dd2a8-8879-4498-9f8f-d1b310235282"
  },
  {
    "state": "lose",
    "amount": "2,917.64",
    "transactionId": "a8173c48-75eb-4448-b7e3-9ae2dc069a85"
  },
  {
    "state": "win",
    "amount": "1,409.37",
    "transactionId": "c481581e-a0e2-495a-bcfb-c6c62a0cef87"
  },
  {
    "state": "win",
    "amount": "2,085.56",
    "transactionId": "9e722b6d-c109-4ac1-83a5-bc2c137d2f11"
  },
  {
    "state": "lose",
    "amount": "1,992.45",
    "transactionId": "9ef44939-1e6f-45ef-b150-2d98d00d10a4"
  },
  {
    "state": "lose",
    "amount": "1,960.14",
    "transactionId": "685101e5-a7a8-4e53-be37-29bc9e7831e2"
  },
  {
    "state": "win",
    "amount": "3,636.71",
    "transactionId": "c93bd7a6-8271-40a6-8152-37380775162b"
  },
  {
    "state": "lose",
    "amount": "1,103.26",
    "transactionId": "4f12edec-071c-4c3a-aed9-afa07fc002a7"
  },
  {
    "state": "win",
    "amount": "2,146.04",
    "transactionId": "fd184102-b5ad-4fbe-861c-277a792e3c0f"
  },
  {
    "state": "lose",
    "amount": "1,727.81",
    "transactionId": "619c9a31-c3c6-479e-ab4b-b0ab7a4776e8"
  },
  {
    "state": "win",
    "amount": "2,133.19",
    "transactionId": "ff20cbee-fecd-409e-822f-99b31d0b217c"
  },
  {
    "state": "win",
    "amount": "1,447.92",
    "transactionId": "1a79b7f6-95d0-43e8-bb16-0fa93dd7d59b"
  },
  {
    "state": "lose",
    "amount": "3,991.76",
    "transactionId": "e1451fa4-d9c6-4bcf-8a28-b4c5223847e2"
  },
  {
    "state": "win",
    "amount": "1,289.29",
    "transactionId": "d3d7e140-a52b-4f8c-aa06-b05c82f548c1"
  },
  {
    "state": "lose",
    "amount": "1,405.79",
    "transactionId": "ad5dc84b-bc4e-4cf0-94fd-321396d08ea5"
  },
  {
    "state": "lose",
    "amount": "2,748.78",
    "transactionId": "a976da25-3185-4fa3-a24a-1f5e70cbbbe8"
  },
  {
    "state": "win",
    "amount": "1,139.30",
    "transactionId": "3365a2d6-8166-457c-9f70-c88473814331"
  },
  {
    "state": "lose",
    "amount": "2,285.77",
    "transactionId": "1a801ab9-22e6-4dfa-bd64-b83ef0e25f75"
  },
  {
    "state": "win",
    "amount": "1,444.90",
    "transactionId": "0e919aac-1e81-4063-993a-70691c0a6e4c"
  },
  {
    "state": "lose",
    "amount": "1,634.59",
    "transactionId": "433ee35c-52d3-480a-9103-9743eeaf9924"
  },
  {
    "state": "win",
    "amount": "2,147.84",
    "transactionId": "54dad892-4902-4581-a89a-c62f97be4639"
  },
  {
    "state": "lose",
    "amount": "2,395.75",
    "transactionId": "28b09d8d-08e2-4234-be4c-a05a710ec39b"
  },
  {
    "state": "lose",
    "amount": "2,786.29",
    "transactionId": "973abca6-c432-46ff-8c7a-3d4a305f4f4f"
  },
  {
    "state": "win",
    "amount": "3,872.54",
    "transactionId": "93f0a22c-5b1d-41db-9b76-5ea86a4e377b"
  },
  {
    "state": "lose",
    "amount": "2,214.67",
    "transactionId": "87c53664-f30f-46a1-a411-3f62c7d8dbd3"
  },
  {
    "state": "lose",
    "amount": "1,995.93",
    "transactionId": "c3f28083-f100-4c8d-ab02-f500706d0ae2"
  },
  {
    "state": "lose",
    "amount": "1,828.50",
    "transactionId": "7369390f-3690-4370-b67d-0daed45fcd21"
  },
  {
    "state": "lose",
    "amount": "2,466.09",
    "transactionId": "1a928d61-bf47-4125-bb37-f6c16e46cf67"
  },
  {
    "state": "lose",
    "amount": "1,535.02",
    "transactionId": "4d455f39-0dc8-4370-ac19-e2fd0a5b7b51"
  },
  {
    "state": "lose",
    "amount": "3,426.47",
    "transactionId": "c42ce82f-9e06-47e7-89fd-2409de0d5324"
  },
  {
    "state": "win",
    "amount": "2,293.06",
    "transactionId": "0452907e-7e28-4c87-8b86-568219995b19"
  },
  {
    "state": "win",
    "amount": "3,120.54",
    "transactionId": "5de10509-a0b3-4376-a507-945d25a5a695"
  },
  {
    "state": "win",
    "amount": "3,755.56",
    "transactionId": "a5db216e-86cc-4530-9694-538021e52f5e"
  },
  {
    "state": "lose",
    "amount": "2,404.56",
    "transactionId": "97a5d505-c4de-47bd-bfa8-b0833bdc27c4"
  },
  {
    "state": "win",
    "amount": "1,638.10",
    "transactionId": "4464fefb-fee6-418f-9ca1-7afc67e03912"
  },
  {
    "state": "lose",
    "amount": "3,402.51",
    "transactionId": "89ed3a82-2446-4266-b9ac-9b042f431c64"
  },
  {
    "state": "lose",
    "amount": "3,206.06",
    "transactionId": "a6d6b427-0425-49ce-8120-42ba01b99085"
  },
  {
    "state": "win",
    "amount": "3,729.57",
    "transactionId": "a891de31-ef70-4d75-b5ee-8fa9f6795873"
  },
  {
    "state": "lose",
    "amount": "2,075.44",
    "transactionId": "dc5e7459-eeec-44db-af83-071b91f93ab0"
  },
  {
    "state": "lose",
    "amount": "3,057.17",
    "transactionId": "365039bf-987b-48a5-8bf7-af43b451ef36"
  },
  {
    "state": "lose",
    "amount": "1,967.54",
    "transactionId": "1d87d951-822c-400a-b2d0-63edba88d49c"
  },
  {
    "state": "lose",
    "amount": "3,964.16",
    "transactionId": "60d28b12-9c85-451a-971e-b1d5f3ac5fdf"
  },
  {
    "state": "lose",
    "amount": "2,378.38",
    "transactionId": "d4a289c9-5082-4022-9f62-596c445d660b"
  },
  {
    "state": "win",
    "amount": "3,383.95",
    "transactionId": "f7f68476-8dd5-4ef8-b8d6-47e5416c5b8d"
  },
  {
    "state": "win",
    "amount": "2,068.46",
    "transactionId": "ee2c643c-1a49-4404-b9df-1d8a75a1eb86"
  },
  {
    "state": "lose",
    "amount": "1,038.12",
    "transactionId": "4088d3bf-4388-4ae2-94e9-445fb37bdb7d"
  },
  {
    "state": "win",
    "amount": "3,229.14",
    "transactionId": "50968bca-5bea-4166-b0f9-3ab2b497fd44"
  },
  {
    "state": "win",
    "amount": "3,222.38",
    "transactionId": "93c05026-8bd3-4698-a26b-e9b633f17688"
  },
  {
    "state": "lose",
    "amount": "1,178.39",
    "transactionId": "d638346c-c5e3-4a88-8a32-41ced720cbf2"
  },
  {
    "state": "win",
    "amount": "2,774.62",
    "transactionId": "0d40a778-631d-4957-bce8-ef4c0dcdf853"
  },
  {
    "state": "win",
    "amount": "3,196.69",
    "transactionId": "18971761-3c17-4d78-a23c-501b889c272b"
  },
  {
    "state": "lose",
    "amount": "3,277.61",
    "transactionId": "3052e88e-b35f-4b36-aa73-c39c217b576c"
  },
  {
    "state": "win",
    "amount": "1,765.82",
    "transactionId": "be788141-4b3b-430f-b387-74109f6e9362"
  },
  {
    "state": "lose",
    "amount": "1,006.33",
    "transactionId": "176bec08-6ce8-45f1-9817-4779c3b6e6b8"
  },
  {
    "state": "win",
    "amount": "3,049.67",
    "transactionId": "d504d039-03b9-4ce7-b6d2-e05a899696b1"
  },
  {
    "state": "lose",
    "amount": "3,536.04",
    "transactionId": "c9308881-9ea1-47ee-a02e-90d556f9064f"
  },
  {
    "state": "win",
    "amount": "3,195.87",
    "transactionId": "df5efbae-5451-4004-b263-2bc6b2ac3a51"
  },
  {
    "state": "lose",
    "amount": "3,669.70",
    "transactionId": "e3c0a6d9-95ff-4481-ad97-66d1d0f8792e"
  },
  {
    "state": "win",
    "amount": "3,347.13",
    "transactionId": "a0d73556-ccd1-47e7-bb07-aa8d2144b2a1"
  },
  {
    "state": "win",
    "amount": "1,355.30",
    "transactionId": "1757b56b-fe5a-472a-91df-26eb8a2b571e"
  },
  {
    "state": "win",
    "amount": "1,273.66",
    "transactionId": "78552319-d58b-4a6c-99b8-d7ada61a5af3"
  },
  {
    "state": "win",
    "amount": "3,995.34",
    "transactionId": "7504ec01-54ad-4750-a2ef-4b68f4c589f1"
  },
  {
    "state": "lose",
    "amount": "2,008.99",
    "transactionId": "d27d81c9-6c3b-4903-bd54-1b08efa23af8"
  },
  {
    "state": "win",
    "amount": "3,366.76",
    "transactionId": "6202a4c8-9a0d-45c3-b0e0-a0c424537eed"
  },
  {
    "state": "lose",
    "amount": "1,437.83",
    "transactionId": "f8357976-bd1a-4171-8b0f-1423ece7adf3"
  },
  {
    "state": "win",
    "amount": "3,109.70",
    "transactionId": "71943d12-f108-4bf8-a9a8-48ad3637367a"
  },
  {
    "state": "win",
    "amount": "1,709.48",
    "transactionId": "dd9da8e1-678e-4232-a23c-162d8d4f7885"
  },
  {
    "state": "win",
    "amount": "3,249.02",
    "transactionId": "c63e2c96-eb38-44e7-a1d3-70492d40c2e8"
  },
  {
    "state": "lose",
    "amount": "3,967.80",
    "transactionId": "3dcb7281-482d-4ab9-86f9-38d3c3b347ca"
  },
  {
    "state": "lose",
    "amount": "3,300.00",
    "transactionId": "723559c0-16f5-4ce7-9cb8-d736506d2c9c"
  },
  {
    "state": "win",
    "amount": "3,249.10",
    "transactionId": "d702b953-9917-4657-8dfc-0d1cd01c478b"
  },
  {
    "state": "win",
    "amount": "2,497.16",
    "transactionId": "cef3f73b-b86f-47ce-a0cb-ad4ce762a07b"
  },
  {
    "state": "win",
    "amount": "2,162.27",
    "transactionId": "950a1f2b-9a22-423a-9656-24d1e2acdf43"
  },
  {
    "state": "lose",
    "amount": "1,243.09",
    "transactionId": "4b31fbac-f4df-48c3-ae63-b18ad9e9d7c8"
  },
  {
    "state": "lose",
    "amount": "1,690.26",
    "transactionId": "e6f7a280-f5db-43fd-ab8d-74482a9e0092"
  },
  {
    "state": "win",
    "amount": "1,157.54",
    "transactionId": "591d3a9b-9473-49c1-b182-b317ad94f77f"
  },
  {
    "state": "lose",
    "amount": "3,896.56",
    "transactionId": "fe57cb88-4a82-4245-aa97-56dd14d15f63"
  },
  {
    "state": "lose",
    "amount": "3,304.25",
    "transactionId": "0242821f-7936-48c1-bf42-a602f833e05e"
  },
  {
    "state": "win",
    "amount": "2,949.33",
    "transactionId": "b3176b9e-7684-4cab-9e5d-4a5b3d45636e"
  },
  {
    "state": "lose",
    "amount": "2,755.71",
    "transactionId": "3152fb41-3df6-4da3-a371-dcf227c34c35"
  },
  {
    "state": "win",
    "amount": "3,439.61",
    "transactionId": "29a6053c-2973-4b17-8198-90a04c97a22e"
  },
  {
    "state": "win",
    "amount": "1,184.52",
    "transactionId": "78dd87dd-9dda-4622-9b76-de52686011df"
  },
  {
    "state": "win",
    "amount": "3,462.80",
    "transactionId": "b91cae3a-fe58-485d-80ea-6d00eb337989"
  },
  {
    "state": "lose",
    "amount": "2,946.98",
    "transactionId": "f56f31c4-1752-48f4-84f1-038c2f1e1205"
  },
  {
    "state": "win",
    "amount": "2,573.18",
    "transactionId": "6d6d3df6-3108-4f59-8460-5bbde257768f"
  },
  {
    "state": "lose",
    "amount": "3,461.73",
    "transactionId": "2ce9f3c9-ef79-4441-a82e-ad2aedc5971c"
  },
  {
    "state": "lose",
    "amount": "1,177.81",
    "transactionId": "7e23e99d-1858-49bc-92cb-3769ca33bdbd"
  },
  {
    "state": "win",
    "amount": "2,258.68",
    "transactionId": "2d297b55-238b-4815-a964-2b45d7b0d3dd"
  },
  {
    "state": "win",
    "amount": "3,084.04",
    "transactionId": "f969782a-90cd-4446-9e2c-08d935194e99"
  },
  {
    "state": "lose",
    "amount": "3,221.28",
    "transactionId": "7ca269f1-bf57-4a9f-9d3e-74cedf218dca"
  },
  {
    "state": "lose",
    "amount": "1,425.86",
    "transactionId": "771bc48c-a68d-480b-a4d0-478fe5c3ab82"
  },
  {
    "state": "win",
    "amount": "3,174.70",
    "transactionId": "0c1b9c84-8db1-49f8-bbff-0304801dab82"
  },
  {
    "state": "lose",
    "amount": "1,198.30",
    "transactionId": "c0cc3faa-19c6-42dd-96a7-36de27904b21"
  },
  {
    "state": "lose",
    "amount": "1,153.50",
    "transactionId": "612350df-38a7-42f4-9ded-2ce93fd4e492"
  },
  {
    "state": "lose",
    "amount": "1,841.77",
    "transactionId": "1d80a1b6-532b-4897-b17c-cdbbff2130fd"
  },
  {
    "state": "win",
    "amount": "1,087.88",
    "transactionId": "e773bd38-5aee-4234-bdae-d16e389be0c3"
  },
  {
    "state": "win",
    "amount": "1,650.38",
    "transactionId": "632d7ddf-7ba8-4661-a4df-720cf4c38f74"
  },
  {
    "state": "lose",
    "amount": "1,545.88",
    "transactionId": "b8658783-18be-4528-9273-ac77c0391d90"
  },
  {
    "state": "lose",
    "amount": "2,741.95",
    "transactionId": "0d610288-a7f4-4ee1-8110-fd1e26b14a67"
  },
  {
    "state": "lose",
    "amount": "1,107.96",
    "transactionId": "786d9dd1-b164-4922-9686-02c493ab5164"
  },
  {
    "state": "win",
    "amount": "3,282.91",
    "transactionId": "c962afbe-7484-4e5a-84b6-c6d75bce80d4"
  },
  {
    "state": "lose",
    "amount": "3,803.92",
    "transactionId": "3370e3b1-5d0b-40f0-9e80-5de6a4f1b0b8"
  },
  {
    "state": "lose",
    "amount": "1,970.78",
    "transactionId": "3aa3d3a1-be8c-4c24-b992-fd0fc3810b4e"
  },
  {
    "state": "win",
    "amount": "3,631.22",
    "transactionId": "a66e2cf6-8040-4634-8d1d-222810175816"
  },
  {
    "state": "lose",
    "amount": "1,176.27",
    "transactionId": "a4988e96-eebb-45db-b291-f71d41b59602"
  },
  {
    "state": "win",
    "amount": "1,817.80",
    "transactionId": "403edc21-3919-49e7-9db1-4facd6dfde72"
  },
  {
    "state": "win",
    "amount": "3,912.62",
    "transactionId": "f67a5b63-982c-441a-ab06-7b9813b2fb46"
  },
  {
    "state": "lose",
    "amount": "3,232.59",
    "transactionId": "de26bea9-eae4-4270-a685-308493b603bf"
  },
  {
    "state": "win",
    "amount": "2,730.46",
    "transactionId": "a941731f-b9aa-4bac-89e4-3eacdbecfeca"
  },
  {
    "state": "win",
    "amount": "3,130.55",
    "transactionId": "f11aad6a-b863-4d80-bcdc-8c4a3ec5efbd"
  },
  {
    "state": "win",
    "amount": "1,135.84",
    "transactionId": "5ba05c75-1676-4a0c-ae93-97ac06699f61"
  },
  {
    "state": "lose",
    "amount": "2,713.57",
    "transactionId": "0be6caa9-0a5d-4b35-9b6e-3840d28df434"
  },
  {
    "state": "win",
    "amount": "3,763.67",
    "transactionId": "5041bdf7-4f30-4db0-9dc9-676c9d56951c"
  },
  {
    "state": "win",
    "amount": "2,946.65",
    "transactionId": "64b20c3a-0191-4a9a-a86a-a91b63f140e1"
  },
  {
    "state": "lose",
    "amount": "2,521.65",
    "transactionId": "be7ed77c-3ab9-446f-b182-0923175ee36d"
  },
  {
    "state": "win",
    "amount": "2,411.97",
    "transactionId": "813e8024-3683-4482-9dcb-c0f8097f8661"
  },
  {
    "state": "win",
    "amount": "1,589.83",
    "transactionId": "d7a593ee-3fa0-4d0f-928a-277ba26dc759"
  },
  {
    "state": "lose",
    "amount": "1,979.57",
    "transactionId": "10ead6fd-5974-4a31-86cc-58640f29a7ac"
  },
  {
    "state": "lose",
    "amount": "1,249.99",
    "transactionId": "576c73ae-6da6-473a-ab7e-2714605682cf"
  },
  {
    "state": "win",
    "amount": "3,812.37",
    "transactionId": "4f477da0-cabe-46f5-b811-a5e0cfc5bdba"
  },
  {
    "state": "win",
    "amount": "1,788.45",
    "transactionId": "c4c4f9f5-df29-4809-a551-dfc94d7b4100"
  },
  {
    "state": "lose",
    "amount": "3,887.46",
    "transactionId": "58bd6900-af12-4ad6-9f52-5f46687ff2e0"
  },
  {
    "state": "lose",
    "amount": "3,386.98",
    "transactionId": "e0f46b69-117b-4705-97eb-9730e019d1e6"
  },
  {
    "state": "win",
    "amount": "1,488.14",
    "transactionId": "a35914b4-af98-432a-9f79-cec2026f0ec7"
  },
  {
    "state": "lose",
    "amount": "1,209.56",
    "transactionId": "5773af93-92f7-4df1-a84f-2a180746758a"
  },
  {
    "state": "lose",
    "amount": "2,191.61",
    "transactionId": "92243237-0b96-483d-af8c-7bb7834095b4"
  },
  {
    "state": "win",
    "amount": "1,975.25",
    "transactionId": "bdcf5b6f-fea6-4d31-af1c-03030bd635c1"
  },
  {
    "state": "win",
    "amount": "2,155.50",
    "transactionId": "245b6431-e994-418d-8543-e3215f2ea535"
  },
  {
    "state": "lose",
    "amount": "1,910.10",
    "transactionId": "e4e9d377-b71f-4cca-902c-0271db537040"
  },
  {
    "state": "lose",
    "amount": "3,778.62",
    "transactionId": "3033d490-675b-4bdd-b404-414ac183c804"
  },
  {
    "state": "lose",
    "amount": "3,873.88",
    "transactionId": "ebe5b229-24f8-4427-83c0-4cb2448183e9"
  },
  {
    "state": "win",
    "amount": "3,634.12",
    "transactionId": "4977633b-3249-428a-8a56-08f382130be1"
  },
  {
    "state": "win",
    "amount": "3,072.10",
    "transactionId": "d46db59f-cddc-4ae0-929c-38bbef10d317"
  },
  {
    "state": "lose",
    "amount": "3,087.52",
    "transactionId": "d24ba3a6-0a20-4c1f-96f6-c5519ead585f"
  },
  {
    "state": "win",
    "amount": "2,069.86",
    "transactionId": "dda8829c-50f6-4dd0-94eb-9be92b4ae525"
  },
  {
    "state": "win",
    "amount": "2,502.46",
    "transactionId": "98617d90-74ab-4343-80d6-c83354fdfec7"
  },
  {
    "state": "lose",
    "amount": "1,973.16",
    "transactionId": "aaf97201-07c7-4d24-a3ce-085665dd5e2a"
  },
  {
    "state": "lose",
    "amount": "1,709.62",
    "transactionId": "189aa732-d4df-4807-b36a-20f60743a92a"
  },
  {
    "state": "lose",
    "amount": "2,050.46",
    "transactionId": "1d4162f8-64de-4a05-8eea-297d58a9b94b"
  },
  {
    "state": "lose",
    "amount": "2,535.65",
    "transactionId": "26e6dc0d-8de5-4628-9223-848654105d8a"
  },
  {
    "state": "win",
    "amount": "1,852.30",
    "transactionId": "a36877e6-873b-44a9-aead-fdde81086f20"
  },
  {
    "state": "win",
    "amount": "1,470.11",
    "transactionId": "d96a1e9f-5b99-4d22-9670-63cb94bf1f19"
  },
  {
    "state": "lose",
    "amount": "3,996.63",
    "transactionId": "d832e6bb-86b9-49a7-8aec-18a9ece79215"
  },
  {
    "state": "lose",
    "amount": "1,581.96",
    "transactionId": "3de811b4-601a-4086-9050-7dd6894ec207"
  },
  {
    "state": "lose",
    "amount": "3,035.02",
    "transactionId": "bd9b285c-f8fa-480a-9406-32d6d3f0133b"
  },
  {
    "state": "win",
    "amount": "1,057.65",
    "transactionId": "b6b71c16-9d44-46e6-a210-216d7a3a21bf"
  },
  {
    "state": "lose",
    "amount": "2,218.16",
    "transactionId": "ba788532-8c89-4305-965c-933b41ed917a"
  },
  {
    "state": "lose",
    "amount": "2,115.23",
    "transactionId": "9c86b428-b002-45cb-b813-596e056298a9"
  },
  {
    "state": "win",
    "amount": "3,601.22",
    "transactionId": "44b04910-0c07-4ad2-8a35-26aee47bb865"
  },
  {
    "state": "win",
    "amount": "1,788.01",
    "transactionId": "6fb7fcab-ac4e-4ade-993e-2775fe07a0a1"
  },
  {
    "state": "win",
    "amount": "3,486.52",
    "transactionId": "f4460c7b-c5d3-4279-a069-8a8e92cd354b"
  },
  {
    "state": "lose",
    "amount": "3,090.31",
    "transactionId": "1b20b8fe-e23f-4f7c-a4d8-8036b8c033db"
  },
  {
    "state": "lose",
    "amount": "2,434.70",
    "transactionId": "eeb5a576-ca7b-464d-b7c5-eb7310583236"
  },
  {
    "state": "lose",
    "amount": "3,830.51",
    "transactionId": "5905e188-4127-4f20-b07a-640e0bb1128b"
  },
  {
    "state": "lose",
    "amount": "3,240.67",
    "transactionId": "d4474cde-ebce-40ae-9393-c1de33bbfd72"
  },
  {
    "state": "lose",
    "amount": "3,851.40",
    "transactionId": "78c5753e-9b22-4b5e-a467-4f15c0cc908c"
  },
  {
    "state": "win",
    "amount": "1,468.22",
    "transactionId": "b2ad9c98-d0b8-415e-91f8-fbad18d5c557"
  },
  {
    "state": "lose",
    "amount": "1,164.98",
    "transactionId": "7d1d858a-fd6a-43f9-a971-8e0ef18d6f83"
  },
  {
    "state": "lose",
    "amount": "2,114.84",
    "transactionId": "6fd39996-dcc9-4629-b03b-5e85ce40693c"
  },
  {
    "state": "win",
    "amount": "3,096.02",
    "transactionId": "38e61253-273d-487e-abd2-9c638c837c91"
  },
  {
    "state": "lose",
    "amount": "3,639.85",
    "transactionId": "78b82b01-d50c-4ae7-ab54-1a4cb44ef166"
  },
  {
    "state": "lose",
    "amount": "2,765.69",
    "transactionId": "a805280c-08a7-49f6-b494-76ade67574be"
  },
  {
    "state": "win",
    "amount": "1,842.52",
    "transactionId": "84eef08e-2d81-4c4c-a606-edc6233409fb"
  },
  {
    "state": "lose",
    "amount": "3,393.87",
    "transactionId": "3bc227b8-9d8f-461d-b18f-8da1b47fc0cd"
  },
  {
    "state": "win",
    "amount": "2,215.14",
    "transactionId": "46549646-e457-4f66-9efc-27f83d2893e0"
  },
  {
    "state": "win",
    "amount": "3,940.13",
    "transactionId": "d1df3a78-9c23-4fe0-845b-eda59c8288c4"
  },
  {
    "state": "win",
    "amount": "3,555.29",
    "transactionId": "900f87c9-a644-4b72-b291-3860f9d8d614"
  },
  {
    "state": "win",
    "amount": "2,070.34",
    "transactionId": "752c97de-e79f-4e97-9016-424c34fd848b"
  },
  {
    "state": "lose",
    "amount": "1,492.86",
    "transactionId": "d5e680bf-f20a-438f-a347-3c9400207866"
  },
  {
    "state": "win",
    "amount": "2,928.85",
    "transactionId": "0ef53c46-cca4-41e7-af3f-6e2836895d63"
  },
  {
    "state": "lose",
    "amount": "2,050.31",
    "transactionId": "1136e347-8385-4b52-8a36-2331b3d8b199"
  },
  {
    "state": "lose",
    "amount": "2,988.87",
    "transactionId": "93a17a48-2a23-4a0f-a2f0-5f1d241dca72"
  },
  {
    "state": "win",
    "amount": "2,984.88",
    "transactionId": "62bb270e-f66f-4ed5-9a19-27e7149b88b7"
  },
  {
    "state": "lose",
    "amount": "2,321.83",
    "transactionId": "2e0a12cd-31e2-46ef-bbc7-1826f05e8369"
  },
  {
    "state": "win",
    "amount": "1,729.61",
    "transactionId": "ab84bd22-bbac-40ce-855f-cee749f74941"
  },
  {
    "state": "lose",
    "amount": "3,924.11",
    "transactionId": "4f71ff57-ada4-4593-ba50-36d0c7ded168"
  },
  {
    "state": "win",
    "amount": "3,908.49",
    "transactionId": "65688958-da84-4d4e-b27c-b9bae920873a"
  },
  {
    "state": "lose",
    "amount": "2,312.41",
    "transactionId": "9389f0ae-33f3-4a93-9e7c-c77371cd5d82"
  },
  {
    "state": "win",
    "amount": "2,198.50",
    "transactionId": "847ecaed-f9a5-4cca-b576-23f59d57a156"
  },
  {
    "state": "win",
    "amount": "2,923.09",
    "transactionId": "02a33a15-cc5d-4a7e-a14e-739098e36cdc"
  },
  {
    "state": "win",
    "amount": "3,811.59",
    "transactionId": "d56eb98c-81fe-49b3-80cc-1f0454d33d5c"
  },
  {
    "state": "lose",
    "amount": "2,373.97",
    "transactionId": "0ac0f016-fc4e-4953-a207-3f4797d92d14"
  },
  {
    "state": "lose",
    "amount": "3,765.61",
    "transactionId": "1964c035-b143-46e9-a37f-20499800f8cf"
  },
  {
    "state": "win",
    "amount": "2,576.65",
    "transactionId": "5fdc63c6-9410-4db0-882b-bce2b2190b01"
  },
  {
    "state": "lose",
    "amount": "1,097.83",
    "transactionId": "8db67331-dc3b-403a-9beb-3fc9e0834c02"
  },
  {
    "state": "win",
    "amount": "2,599.28",
    "transactionId": "417b54b1-32ad-4b8a-814e-89c99f19d5cc"
  },
  {
    "state": "lose",
    "amount": "2,511.32",
    "transactionId": "e9394964-052e-4319-b448-df648aaa8e60"
  },
  {
    "state": "lose",
    "amount": "1,184.71",
    "transactionId": "541165c5-2d3a-4f03-b8ae-7d2586771c81"
  },
  {
    "state": "win",
    "amount": "1,031.75",
    "transactionId": "487125e6-0ecc-4827-8f60-8ec763e59c80"
  },
  {
    "state": "win",
    "amount": "2,568.94",
    "transactionId": "c05c0e31-37c1-4bc7-a198-c9e89407e196"
  },
  {
    "state": "lose",
    "amount": "3,114.08",
    "transactionId": "5d3d6e60-c346-4c23-809e-114763a20950"
  },
  {
    "state": "win",
    "amount": "2,463.49",
    "transactionId": "ec9955a1-06b5-4c5e-8415-2c92fc3ea07c"
  },
  {
    "state": "lose",
    "amount": "1,453.17",
    "transactionId": "ba28ea99-4191-47be-8c69-eb9ddc711f50"
  },
  {
    "state": "win",
    "amount": "3,350.95",
    "transactionId": "9d76dddf-ac24-47b4-b9dd-8a3b9ba4f745"
  },
  {
    "state": "lose",
    "amount": "1,685.84",
    "transactionId": "725d20b1-4da1-47cb-823e-9783998b6a26"
  },
  {
    "state": "win",
    "amount": "1,837.49",
    "transactionId": "ddf6cefb-d0d4-442d-af15-378f059b3709"
  },
  {
    "state": "lose",
    "amount": "1,425.23",
    "transactionId": "a67b2f42-12df-4395-85bb-6e06d10719cd"
  },
  {
    "state": "win",
    "amount": "1,324.82",
    "transactionId": "4628e5a4-4137-4a18-95f1-655ee9fdb163"
  },
  {
    "state": "win",
    "amount": "2,841.34",
    "transactionId": "4a53388c-03bf-4ace-a416-51d41da2a11a"
  },
  {
    "state": "lose",
    "amount": "1,251.41",
    "transactionId": "3f5ad090-6a42-43ef-a0d6-59975f80f57f"
  },
  {
    "state": "win",
    "amount": "3,563.16",
    "transactionId": "86c90a83-73fe-432b-a487-218f1c912d8d"
  },
  {
    "state": "win",
    "amount": "3,729.57",
    "transactionId": "7dcec45a-bfda-44c1-b673-93bd156b3ea7"
  },
  {
    "state": "lose",
    "amount": "3,998.20",
    "transactionId": "7a2f2229-da65-4b09-9a27-0740e3a24e2d"
  },
  {
    "state": "lose",
    "amount": "1,437.44",
    "transactionId": "fd8dccb0-c036-42e9-9e81-acd122cd3800"
  },
  {
    "state": "win",
    "amount": "2,449.46",
    "transactionId": "a7fce126-4723-408d-bdc7-e3a4e54a83db"
  },
  {
    "state": "win",
    "amount": "1,355.96",
    "transactionId": "cd6ad5f3-008d-4ba9-810f-18b1f1133348"
  },
  {
    "state": "lose",
    "amount": "1,105.88",
    "transactionId": "25940266-81dc-4646-9923-e40313915cd8"
  },
  {
    "state": "win",
    "amount": "1,060.28",
    "transactionId": "411dc689-8f9e-4f01-9e03-82c9f0bd7d28"
  },
  {
    "state": "win",
    "amount": "1,077.64",
    "transactionId": "4b6958d9-a6e1-4d02-a6f0-f2ae6ab3f9b0"
  },
  {
    "state": "lose",
    "amount": "2,102.70",
    "transactionId": "2671418c-51a3-474c-923f-588a2d12c2dd"
  },
  {
    "state": "win",
    "amount": "1,859.04",
    "transactionId": "a9457cd1-e572-4c49-a2c4-ac2bf2d6797a"
  },
  {
    "state": "win",
    "amount": "2,807.38",
    "transactionId": "941c27f6-007b-4939-89db-6b99f1391ba2"
  },
  {
    "state": "lose",
    "amount": "3,073.27",
    "transactionId": "3e4f863d-fb95-4e70-8372-a1193e909d00"
  },
  {
    "state": "lose",
    "amount": "1,729.34",
    "transactionId": "b12458d0-d675-4e75-bee0-027c7d8d589c"
  },
  {
    "state": "win",
    "amount": "2,473.89",
    "transactionId": "e0dcf7be-c5f0-4a42-b16c-b73e74946373"
  },
  {
    "state": "lose",
    "amount": "2,202.02",
    "transactionId": "e88f1ec2-98fe-438c-84db-a16e432dd767"
  },
  {
    "state": "win",
    "amount": "2,117.00",
    "transactionId": "6edbf8b9-9860-40fc-a089-1d2b837343a4"
  },
  {
    "state": "lose",
    "amount": "3,021.21",
    "transactionId": "8d584fb2-e190-4251-baa2-2421e6cf8305"
  },
  {
    "state": "lose",
    "amount": "2,806.37",
    "transactionId": "da447b4e-b2e7-4429-a02f-02e64c86d9d3"
  },
  {
    "state": "lose",
    "amount": "2,398.37",
    "transactionId": "eb8b1536-d6dc-4caa-ad54-88a292a56dbd"
  },
  {
    "state": "win",
    "amount": "1,945.64",
    "transactionId": "2bb4803c-dd7b-44ff-bb18-ab72a9a1bc8d"
  },
  {
    "state": "win",
    "amount": "1,623.35",
    "transactionId": "f27c7bf6-ac1e-45d4-b02d-3d6b866e34b2"
  },
  {
    "state": "lose",
    "amount": "2,609.10",
    "transactionId": "614663aa-bf3c-4f35-880b-461937a9fabc"
  },
  {
    "state": "lose",
    "amount": "1,896.85",
    "transactionId": "5b638764-ea75-4cb0-9400-0c99c30962a6"
  },
  {
    "state": "win",
    "amount": "1,652.88",
    "transactionId": "97042082-7fb1-4348-9c8d-cbd2ab459a8c"
  },
  {
    "state": "lose",
    "amount": "1,623.79",
    "transactionId": "d785ef96-00e4-4e4f-8538-69aa6d13cb06"
  },
  {
    "state": "lose",
    "amount": "2,577.08",
    "transactionId": "800a7699-e30e-4b94-8d00-95d5a64982a9"
  },
  {
    "state": "win",
    "amount": "1,192.80",
    "transactionId": "5037df13-362b-4bc0-af3e-9a04eebabf99"
  },
  {
    "state": "lose",
    "amount": "1,368.45",
    "transactionId": "3932e358-bb27-45d9-ba72-9a66af13ddf1"
  },
  {
    "state": "win",
    "amount": "3,847.19",
    "transactionId": "392b86c3-d220-4146-89ed-80d73bcb469c"
  },
  {
    "state": "lose",
    "amount": "3,643.03",
    "transactionId": "46cc203c-63a7-4e33-9832-003fb3ec20dd"
  },
  {
    "state": "lose",
    "amount": "2,992.05",
    "transactionId": "12857cc6-74bf-4549-99c4-4e9a5f036e26"
  },
  {
    "state": "lose",
    "amount": "3,226.53",
    "transactionId": "3c2eb14a-29f4-4f9a-919e-15a7802795eb"
  },
  {
    "state": "lose",
    "amount": "3,048.83",
    "transactionId": "98eb2f73-a11a-410a-9fcd-018911b33e47"
  },
  {
    "state": "win",
    "amount": "3,783.32",
    "transactionId": "bc9fb698-0393-47bc-935c-db26e4b17ec1"
  },
  {
    "state": "lose",
    "amount": "2,765.66",
    "transactionId": "0c482997-1315-4945-9fda-b00031b5bcf2"
  },
  {
    "state": "win",
    "amount": "2,615.12",
    "transactionId": "56f86e09-b2a6-448c-b63c-fcf04320fd9a"
  },
  {
    "state": "win",
    "amount": "1,848.53",
    "transactionId": "e7207ab3-530b-4602-8487-672b14c05eaf"
  },
  {
    "state": "lose",
    "amount": "2,045.70",
    "transactionId": "cdaa3c48-3c50-4327-a010-e9d350f9c3de"
  },
  {
    "state": "lose",
    "amount": "2,185.23",
    "transactionId": "8f352ad4-ee26-421c-b1c6-d5d8e8368dde"
  },
  {
    "state": "lose",
    "amount": "1,168.74",
    "transactionId": "87e18bfb-5c75-414e-b68c-ee8d2918509e"
  },
  {
    "state": "lose",
    "amount": "2,417.68",
    "transactionId": "1d249bea-f1a0-40c8-9459-0969b1d0c0fd"
  },
  {
    "state": "win",
    "amount": "3,367.07",
    "transactionId": "9407624e-7ca1-460e-a9ef-80af0cf1f354"
  },
  {
    "state": "lose",
    "amount": "1,369.45",
    "transactionId": "73c37ef3-3101-4426-af06-156bd5e642f0"
  },
  {
    "state": "lose",
    "amount": "1,914.15",
    "transactionId": "90e6ee69-f2c3-4a2c-8360-88ea2e7b59e4"
  },
  {
    "state": "lose",
    "amount": "1,914.77",
    "transactionId": "955aa712-665c-4fd1-801c-427ba8815e8f"
  },
  {
    "state": "lose",
    "amount": "2,036.86",
    "transactionId": "e3c21bf0-3377-4582-84e9-0e48f1927c10"
  },
  {
    "state": "win",
    "amount": "3,172.49",
    "transactionId": "f6c85314-6b41-4406-af48-51951c2f232e"
  },
  {
    "state": "lose",
    "amount": "3,566.89",
    "transactionId": "c6097043-5ea8-46f1-a655-4888b461f726"
  },
  {
    "state": "win",
    "amount": "3,705.97",
    "transactionId": "24958c3b-d624-4af5-8234-7c494a5e28f6"
  },
  {
    "state": "win",
    "amount": "1,661.47",
    "transactionId": "d6115641-5efe-4158-9ed9-7c689c082d87"
  },
  {
    "state": "win",
    "amount": "3,617.10",
    "transactionId": "7c750c6a-96c6-49e7-b801-2e7490dead93"
  },
  {
    "state": "lose",
    "amount": "2,206.59",
    "transactionId": "c31062db-983b-405b-9470-0841e12496f9"
  },
  {
    "state": "win",
    "amount": "1,151.17",
    "transactionId": "1b0eefba-ecd6-44c6-bd4c-df4b6d275c5d"
  },
  {
    "state": "lose",
    "amount": "2,280.73",
    "transactionId": "c805f9f5-5fa0-42a8-8256-cd9b2b0e8c26"
  },
  {
    "state": "win",
    "amount": "2,769.11",
    "transactionId": "6a611a41-c430-4047-960d-7e23e93ec123"
  },
  {
    "state": "win",
    "amount": "3,063.11",
    "transactionId": "1df4c5b7-4199-49c4-bd52-093e45b74c8b"
  },
  {
    "state": "win",
    "amount": "1,630.28",
    "transactionId": "2468192d-1737-4d66-acc7-af9881692861"
  },
  {
    "state": "win",
    "amount": "1,835.97",
    "transactionId": "0fd75544-01f1-4e7f-90a6-960d7d304bf4"
  },
  {
    "state": "win",
    "amount": "1,486.82",
    "transactionId": "904523b6-bc14-4a4f-87c9-feeb32a5fc92"
  },
  {
    "state": "lose",
    "amount": "3,057.62",
    "transactionId": "fa811a4c-9475-4da7-ac06-466ee4b1284f"
  },
  {
    "state": "win",
    "amount": "2,102.82",
    "transactionId": "540d4511-6144-4ad4-9e8d-bbc5fb84f103"
  },
  {
    "state": "win",
    "amount": "2,666.76",
    "transactionId": "565a4849-50dd-46ec-8beb-811aeb42edfa"
  },
  {
    "state": "win",
    "amount": "3,499.96",
    "transactionId": "91cfaf93-1d35-4c96-95e2-f4be4c406373"
  },
  {
    "state": "win",
    "amount": "3,536.55",
    "transactionId": "d408f1e4-3223-4476-bc0a-09cac955d522"
  },
  {
    "state": "lose",
    "amount": "3,305.15",
    "transactionId": "12b0ac96-41c1-4cc9-914e-491d45695853"
  },
  {
    "state": "lose",
    "amount": "2,035.64",
    "transactionId": "9c252624-87a8-4cdd-a6ad-2b53514314c7"
  },
  {
    "state": "lose",
    "amount": "1,785.13",
    "transactionId": "e30f51d5-f1d7-40c5-95fc-0821fb7f5ae5"
  },
  {
    "state": "lose",
    "amount": "1,927.86",
    "transactionId": "bf9fdb11-e801-4602-b5c4-6abeffc9a13c"
  },
  {
    "state": "win",
    "amount": "3,278.27",
    "transactionId": "5002bb9d-df30-4f77-a969-b2272f8e572e"
  },
  {
    "state": "win",
    "amount": "1,133.89",
    "transactionId": "e597c97c-58d4-4bc6-8bb9-2b99925bc058"
  },
  {
    "state": "win",
    "amount": "2,066.59",
    "transactionId": "9e035fd9-38c9-4415-a9d1-2edad857fb57"
  },
  {
    "state": "lose",
    "amount": "1,130.49",
    "transactionId": "fa5fbcc1-f675-4d8e-9c68-20d838e429a1"
  },
  {
    "state": "lose",
    "amount": "3,459.84",
    "transactionId": "6990462f-7aae-4ec0-85ac-02555f2854f0"
  },
  {
    "state": "win",
    "amount": "1,073.42",
    "transactionId": "fb5e51e6-e301-424a-93b6-06c2dd93ba5b"
  },
  {
    "state": "lose",
    "amount": "3,741.95",
    "transactionId": "e34e14cb-885e-4b19-a0a7-b7e54e0a7a00"
  },
  {
    "state": "lose",
    "amount": "3,541.82",
    "transactionId": "e6e75a79-7067-42fd-a7fc-d9a370949d83"
  },
  {
    "state": "win",
    "amount": "1,298.21",
    "transactionId": "0a96e2a4-d185-425c-99d3-c732ae115e32"
  },
  {
    "state": "win",
    "amount": "3,893.65",
    "transactionId": "51199ce0-891b-424f-9922-66d8c5d016c8"
  },
  {
    "state": "win",
    "amount": "2,849.78",
    "transactionId": "aacfb30b-c5e8-4f1a-b32d-155a22b4841b"
  },
  {
    "state": "lose",
    "amount": "3,138.64",
    "transactionId": "551f50f2-f450-459c-9cfb-b1b231b9d623"
  },
  {
    "state": "win",
    "amount": "1,054.32",
    "transactionId": "b3df4500-2f08-4a33-899d-f1db5da56a96"
  },
  {
    "state": "win",
    "amount": "2,836.19",
    "transactionId": "ecaa546e-d73c-4afb-85c8-0dc30b922f49"
  },
  {
    "state": "win",
    "amount": "2,404.55",
    "transactionId": "e3b3372e-4a04-4f9a-9a52-fd384c2b854f"
  },
  {
    "state": "win",
    "amount": "2,472.35",
    "transactionId": "d67757da-81c2-45fc-92b1-22f8df4e8b10"
  },
  {
    "state": "win",
    "amount": "2,054.30",
    "transactionId": "fbb200b9-55af-4338-838b-8b8c457201f9"
  },
  {
    "state": "lose",
    "amount": "3,892.90",
    "transactionId": "912d5325-f843-4198-944c-379423fe2941"
  },
  {
    "state": "win",
    "amount": "1,004.07",
    "transactionId": "b08ab281-205d-4b02-a304-f8011ec7f961"
  },
  {
    "state": "win",
    "amount": "3,912.23",
    "transactionId": "16e44736-7cce-4a9d-9631-4b84ce73ae22"
  },
  {
    "state": "lose",
    "amount": "2,394.09",
    "transactionId": "30746c7a-2b0b-4ab0-8543-fa3862806a4c"
  },
  {
    "state": "lose",
    "amount": "1,358.63",
    "transactionId": "93b03bc3-b208-416e-a7b2-93beb59a6478"
  },
  {
    "state": "lose",
    "amount": "3,156.06",
    "transactionId": "4ad90f29-075b-44d0-8e1c-198a1f75927e"
  },
  {
    "state": "lose",
    "amount": "2,576.32",
    "transactionId": "cbdadb5a-b601-4096-9794-adaddf2c54f9"
  },
  {
    "state": "win",
    "amount": "3,538.00",
    "transactionId": "5ffb9775-2a0a-4935-898c-b9ab5c9d977c"
  },
  {
    "state": "win",
    "amount": "1,960.91",
    "transactionId": "284d2754-7021-4f86-ab5d-739adf7cb8df"
  },
  {
    "state": "win",
    "amount": "1,199.72",
    "transactionId": "956bd498-6e27-40f2-a55b-f34b1fee23fb"
  },
  {
    "state": "lose",
    "amount": "3,626.71",
    "transactionId": "79dce787-10a5-40f5-8903-baff13601b8f"
  },
  {
    "state": "win",
    "amount": "3,144.17",
    "transactionId": "0ae56192-b928-4589-98a5-eff23fa54204"
  },
  {
    "state": "lose",
    "amount": "1,266.13",
    "transactionId": "4e5fdf21-5e7e-4c9b-8f15-559f18886769"
  },
  {
    "state": "win",
    "amount": "2,201.63",
    "transactionId": "e66d1043-c7b1-46d9-bc21-7cc958391459"
  },
  {
    "state": "lose",
    "amount": "2,329.78",
    "transactionId": "b386dc65-6bcc-4bc9-b6a4-d911b413b1e0"
  },
  {
    "state": "win",
    "amount": "3,174.95",
    "transactionId": "0a0b9e29-f3aa-4624-88dd-2aac33b866a4"
  },
  {
    "state": "win",
    "amount": "1,595.73",
    "transactionId": "cdc31d24-e82d-408f-843f-d7ffdf44609b"
  },
  {
    "state": "lose",
    "amount": "3,504.07",
    "transactionId": "5a66f0ae-3861-4106-894a-4f877746c038"
  },
  {
    "state": "lose",
    "amount": "1,446.90",
    "transactionId": "26374c91-4c6b-4523-a453-21dbc99f8bfc"
  },
  {
    "state": "lose",
    "amount": "2,448.05",
    "transactionId": "8403ccbc-6bfd-47ff-a222-26366eba223a"
  },
  {
    "state": "lose",
    "amount": "2,943.35",
    "transactionId": "9333616b-1330-40b6-aac6-3eac9cd9e97a"
  },
  {
    "state": "win",
    "amount": "1,591.67",
    "transactionId": "02089b36-1e13-4e68-bfd2-fc1a1de5a1ee"
  },
  {
    "state": "win",
    "amount": "1,753.08",
    "transactionId": "7134fe7d-5153-404b-834a-bd2f680dcc70"
  },
  {
    "state": "win",
    "amount": "1,332.97",
    "transactionId": "a3dcf5a7-421e-4407-9c58-e6d91925e7a5"
  },
  {
    "state": "lose",
    "amount": "3,866.98",
    "transactionId": "10157147-56cf-45d9-b5c4-f6ca077ab45b"
  },
  {
    "state": "lose",
    "amount": "1,064.44",
    "transactionId": "eeb61275-ce8d-40d4-bcda-34371f077983"
  },
  {
    "state": "lose",
    "amount": "3,414.89",
    "transactionId": "0952f8f2-90bb-4d4a-a866-d972bf8e057b"
  },
  {
    "state": "win",
    "amount": "2,423.03",
    "transactionId": "267d46e8-a8f2-4905-b586-d1e4feedd15b"
  },
  {
    "state": "win",
    "amount": "1,896.91",
    "transactionId": "7fc15001-5f86-49f3-8157-72569ce49ecc"
  },
  {
    "state": "win",
    "amount": "1,550.53",
    "transactionId": "aff517e9-7975-4f81-aea6-caaa6f7f8ef6"
  },
  {
    "state": "lose",
    "amount": "1,497.35",
    "transactionId": "744ead75-01a7-446b-881e-3bedf54a27d1"
  },
  {
    "state": "lose",
    "amount": "3,037.08",
    "transactionId": "3760065e-990b-4a5e-9962-a7da2ef1b115"
  },
  {
    "state": "lose",
    "amount": "1,946.76",
    "transactionId": "5e919130-70d1-4cad-b721-e968f3e40d0f"
  },
  {
    "state": "lose",
    "amount": "2,872.53",
    "transactionId": "e39f6e5a-f4db-4b77-8c3a-5ee46f0dd162"
  },
  {
    "state": "win",
    "amount": "1,495.55",
    "transactionId": "dba2fac1-21d7-4948-b4cb-83fe7e61353a"
  },
  {
    "state": "win",
    "amount": "2,320.70",
    "transactionId": "9cf9a728-a676-4e28-b576-924a3f76e947"
  },
  {
    "state": "lose",
    "amount": "3,514.27",
    "transactionId": "7bde7768-889a-4ff2-b92e-5f4f579dac02"
  },
  {
    "state": "win",
    "amount": "2,101.43",
    "transactionId": "2794bf32-a052-4947-bfb3-4e810236fd8f"
  },
  {
    "state": "win",
    "amount": "2,476.77",
    "transactionId": "195bffc1-5129-4068-940a-7f13a50d3237"
  },
  {
    "state": "win",
    "amount": "3,907.13",
    "transactionId": "3483ee43-7e78-4853-83c5-e0e4d2353c63"
  },
  {
    "state": "lose",
    "amount": "3,181.09",
    "transactionId": "8f779e91-fb1c-4482-a00d-90f447ec1527"
  },
  {
    "state": "win",
    "amount": "1,942.30",
    "transactionId": "389a5043-dcd3-48ec-bbca-add4f83f9057"
  },
  {
    "state": "win",
    "amount": "3,090.80",
    "transactionId": "c25c3ddd-e935-4862-8eeb-b2e57424f4e2"
  },
  {
    "state": "win",
    "amount": "2,406.97",
    "transactionId": "49777721-1d81-4928-8f71-357a80a806a5"
  },
  {
    "state": "lose",
    "amount": "1,938.58",
    "transactionId": "37d832d6-0e4a-4be2-b5aa-82031474e278"
  },
  {
    "state": "lose",
    "amount": "1,600.58",
    "transactionId": "6a0126ba-9a8e-45c7-9635-df9e22c94872"
  },
  {
    "state": "win",
    "amount": "3,529.14",
    "transactionId": "d1531cd6-cc80-4d48-a139-70cda2abf641"
  },
  {
    "state": "lose",
    "amount": "2,338.11",
    "transactionId": "d0548ff1-be6a-4b9e-9873-874ac7f69175"
  },
  {
    "state": "lose",
    "amount": "1,939.45",
    "transactionId": "df2b57df-45a7-4efa-801d-86ca52c186a9"
  },
  {
    "state": "lose",
    "amount": "1,673.41",
    "transactionId": "fadddf49-aa5b-4ab1-90b4-87f2dc4f6faa"
  },
  {
    "state": "lose",
    "amount": "2,739.74",
    "transactionId": "3f2c53ed-a214-4cc3-b063-e11bea42c281"
  },
  {
    "state": "win",
    "amount": "1,778.58",
    "transactionId": "c9820ddd-35bd-43ac-9177-b78af66c58cc"
  },
  {
    "state": "win",
    "amount": "3,264.55",
    "transactionId": "12a4a199-3019-4905-b0ed-01fa8a9716fd"
  },
  {
    "state": "lose",
    "amount": "3,732.49",
    "transactionId": "121df20f-7fdf-4f27-b76d-4e2aae0e54bf"
  },
  {
    "state": "lose",
    "amount": "3,816.91",
    "transactionId": "97f234ff-142d-44ad-a026-0fbc43f297c9"
  },
  {
    "state": "lose",
    "amount": "3,261.87",
    "transactionId": "d2455f91-c348-4939-a612-32e0e1889d51"
  },
  {
    "state": "win",
    "amount": "2,562.04",
    "transactionId": "16b38abb-dd59-45aa-ac00-0a4849bf99aa"
  },
  {
    "state": "lose",
    "amount": "3,301.22",
    "transactionId": "1dd5c7b3-be76-44e4-8421-d5f17cf8a008"
  },
  {
    "state": "lose",
    "amount": "1,182.01",
    "transactionId": "ad373399-c6e3-4438-a84f-236e6957d591"
  },
  {
    "state": "win",
    "amount": "1,001.78",
    "transactionId": "c69c9089-5168-4815-b408-b4edf2735ace"
  },
  {
    "state": "lose",
    "amount": "2,824.77",
    "transactionId": "2804a0b3-9157-40a7-99d8-ba2ae4584981"
  },
  {
    "state": "win",
    "amount": "2,753.19",
    "transactionId": "7e6e5ff7-4fd3-4481-99c0-fc668800feee"
  },
  {
    "state": "win",
    "amount": "1,347.95",
    "transactionId": "a6411ff4-03d8-4482-b0b0-3b74bc0bf22a"
  },
  {
    "state": "win",
    "amount": "1,963.87",
    "transactionId": "50251568-9a6a-45ee-ad51-34edb0fdcc9b"
  },
  {
    "state": "lose",
    "amount": "2,820.23",
    "transactionId": "34d346f5-fd16-4ce8-a8d4-296d88ed5e23"
  },
  {
    "state": "lose",
    "amount": "2,798.91",
    "transactionId": "c1dd332e-cb0b-43b2-b3b6-34f7ee4614dc"
  },
  {
    "state": "lose",
    "amount": "1,688.39",
    "transactionId": "7bb280e1-559a-4c81-b8c2-e87730709809"
  },
  {
    "state": "win",
    "amount": "2,598.13",
    "transactionId": "12952ffb-02af-4ef6-b3d3-886dfa8d47a6"
  },
  {
    "state": "win",
    "amount": "1,108.61",
    "transactionId": "3961aa1d-3ca5-4e64-af40-b6c0b88f9133"
  },
  {
    "state": "lose",
    "amount": "1,607.99",
    "transactionId": "1e9d3797-1f48-44df-b1d2-edadd692dfe1"
  },
  {
    "state": "win",
    "amount": "3,867.74",
    "transactionId": "df6da640-e4b9-4dfc-80da-4904d7769de9"
  },
  {
    "state": "win",
    "amount": "2,629.33",
    "transactionId": "826fa0a2-5a42-4c26-bc85-9ccd54e1bae2"
  },
  {
    "state": "lose",
    "amount": "1,453.07",
    "transactionId": "847d85ac-4364-4d6a-a15d-c3be018da0c0"
  },
  {
    "state": "lose",
    "amount": "2,544.12",
    "transactionId": "bee79212-ce61-4b13-9a44-6367b7694228"
  },
  {
    "state": "win",
    "amount": "1,456.75",
    "transactionId": "6fd06656-3fd1-4602-bcc0-03134366a825"
  },
  {
    "state": "win",
    "amount": "1,023.86",
    "transactionId": "c85771f9-909b-4614-ab8b-fb9fba990895"
  },
  {
    "state": "win",
    "amount": "1,907.09",
    "transactionId": "3b2df1b1-05d6-47d4-b3d9-776d0367268f"
  },
  {
    "state": "lose",
    "amount": "2,967.70",
    "transactionId": "fffe421e-a44b-44c4-88eb-2cd52d3dbcb3"
  },
  {
    "state": "win",
    "amount": "2,521.59",
    "transactionId": "0b96fa3f-c2db-4be8-bfbb-0787c6ce6472"
  },
  {
    "state": "lose",
    "amount": "2,425.98",
    "transactionId": "91f69c41-fcf2-4444-b359-b4ed2f4bf267"
  },
  {
    "state": "lose",
    "amount": "2,419.22",
    "transactionId": "dac6f64b-0f98-434a-91fa-bf55923af202"
  },
  {
    "state": "lose",
    "amount": "1,861.63",
    "transactionId": "d82eddcc-d132-44a0-9884-e0345b9a17ce"
  },
  {
    "state": "lose",
    "amount": "3,231.48",
    "transactionId": "61659e24-f0b2-4e96-b557-2959c043865d"
  },
  {
    "state": "lose",
    "amount": "3,520.78",
    "transactionId": "4d91f96f-298a-4d3f-a64d-ea3b4296c873"
  },
  {
    "state": "win",
    "amount": "1,142.01",
    "transactionId": "01d6ee1c-16c1-43ab-88b1-4063d0c76163"
  },
  {
    "state": "win",
    "amount": "3,306.97",
    "transactionId": "03c7c5a5-22ef-4b1f-a2ed-50ae0221599a"
  },
  {
    "state": "lose",
    "amount": "3,033.41",
    "transactionId": "5023558c-f69c-4b6d-bbe7-76d70dcf6a0c"
  },
  {
    "state": "lose",
    "amount": "2,366.73",
    "transactionId": "fccd3647-2a71-4af1-8551-1e0b381f1780"
  },
  {
    "state": "win",
    "amount": "2,993.38",
    "transactionId": "42fe4e0d-fcae-48d0-a7c9-d72d8b41f902"
  },
  {
    "state": "win",
    "amount": "2,779.23",
    "transactionId": "0f4f662e-ece1-4b5f-a070-b29e580a84fa"
  },
  {
    "state": "lose",
    "amount": "1,720.41",
    "transactionId": "f551fad2-de98-430e-8f27-efce61a36ada"
  },
  {
    "state": "lose",
    "amount": "2,821.70",
    "transactionId": "2647097d-db4a-4073-a9fe-9c4c23123909"
  },
  {
    "state": "win",
    "amount": "1,518.73",
    "transactionId": "b9e6d4dc-c556-4da6-a7a2-8e8d987a510e"
  },
  {
    "state": "lose",
    "amount": "1,690.80",
    "transactionId": "ff849504-b9df-4dc1-a030-754d135c7992"
  },
  {
    "state": "win",
    "amount": "1,662.21",
    "transactionId": "51f7ef86-9f11-4cfd-ae59-ceed0728ae75"
  },
  {
    "state": "lose",
    "amount": "2,821.02",
    "transactionId": "187bd18c-7c62-4e33-a1cf-9db647b6ce18"
  },
  {
    "state": "lose",
    "amount": "2,082.56",
    "transactionId": "cbd234f8-d21d-4884-9fb6-8f1ac16934f8"
  },
  {
    "state": "win",
    "amount": "2,151.13",
    "transactionId": "9650d216-2d9a-4755-bf0e-5ab07092ffa9"
  },
  {
    "state": "lose",
    "amount": "1,712.35",
    "transactionId": "65d43b84-ea83-4c12-a469-02dec9c5b19b"
  },
  {
    "state": "lose",
    "amount": "2,838.27",
    "transactionId": "8b6d43de-f56b-4fea-9bc5-8fd2ddfc17a7"
  },
  {
    "state": "win",
    "amount": "3,598.10",
    "transactionId": "f637beb7-7288-49a6-9b55-1a3b932d17b5"
  },
  {
    "state": "lose",
    "amount": "1,246.02",
    "transactionId": "9ffe598a-d16b-4601-9419-72e0aeb2a3bd"
  },
  {
    "state": "win",
    "amount": "2,718.75",
    "transactionId": "1be3613d-583c-4859-b3c7-f158aaf6de53"
  },
  {
    "state": "lose",
    "amount": "1,409.42",
    "transactionId": "45937e25-cc62-4b72-9c1d-3dc8e78a8efc"
  },
  {
    "state": "lose",
    "amount": "2,619.01",
    "transactionId": "6ac2f8fb-acbc-488d-a254-687f8aec78e9"
  },
  {
    "state": "lose",
    "amount": "2,661.61",
    "transactionId": "d6938aa3-484b-4a09-afb6-82af2f4e9be5"
  },
  {
    "state": "lose",
    "amount": "3,342.23",
    "transactionId": "dd6a9302-1048-4be2-a2f8-2b1c52f2f28f"
  },
  {
    "state": "win",
    "amount": "3,636.01",
    "transactionId": "04e9a165-0587-4a74-a9ee-061adfa078fc"
  },
  {
    "state": "win",
    "amount": "3,704.28",
    "transactionId": "7b642f06-d02d-4613-93bb-0386db83028b"
  },
  {
    "state": "lose",
    "amount": "1,352.29",
    "transactionId": "98c44153-c1b7-430c-9384-841ad9d35128"
  },
  {
    "state": "win",
    "amount": "2,979.51",
    "transactionId": "db2c347e-9726-493b-9065-98498902ea9d"
  },
  {
    "state": "win",
    "amount": "2,218.30",
    "transactionId": "638b6d60-55dc-4d88-a32f-4eb17c322d08"
  },
  {
    "state": "win",
    "amount": "2,113.64",
    "transactionId": "eb4310a8-8675-419f-90d2-4d1603667fa6"
  },
  {
    "state": "win",
    "amount": "2,093.26",
    "transactionId": "3312dbf5-2a55-434f-b18b-a1fc0c9c221f"
  },
  {
    "state": "lose",
    "amount": "2,797.07",
    "transactionId": "c26a708b-187c-4b21-8443-65bbfeeed1a5"
  },
  {
    "state": "win",
    "amount": "1,022.30",
    "transactionId": "587be099-9d00-4df3-ad9f-e8e239e66dd8"
  },
  {
    "state": "lose",
    "amount": "3,527.52",
    "transactionId": "2a03850d-72cf-4ad0-8718-e7e05d4fa49f"
  },
  {
    "state": "lose",
    "amount": "3,837.34",
    "transactionId": "daccd890-64f7-42df-8316-9305c9e9826e"
  },
  {
    "state": "win",
    "amount": "1,783.37",
    "transactionId": "cc2e1928-4d14-441e-936a-02c0b063f6b9"
  },
  {
    "state": "win",
    "amount": "2,853.45",
    "transactionId": "f0ced5af-75f4-4dad-bf57-532cac42fa75"
  },
  {
    "state": "win",
    "amount": "1,960.02",
    "transactionId": "086ca1b3-c69e-4424-ab29-10a29517d7c7"
  },
  {
    "state": "lose",
    "amount": "2,922.28",
    "transactionId": "97924010-f653-4689-9c5b-cf9d25351c71"
  },
  {
    "state": "win",
    "amount": "2,956.71",
    "transactionId": "ea13f1b6-c057-4790-b1b7-f697f1292607"
  },
  {
    "state": "win",
    "amount": "3,659.68",
    "transactionId": "1be7e9f9-612f-40b0-9c21-0e951bc9086a"
  },
  {
    "state": "win",
    "amount": "1,348.94",
    "transactionId": "c5296553-2834-4906-af43-99e3792f3388"
  },
  {
    "state": "win",
    "amount": "3,883.76",
    "transactionId": "cd67ab51-f1cd-4675-8baf-7d31842e8713"
  },
  {
    "state": "lose",
    "amount": "3,944.07",
    "transactionId": "a55924c4-d786-4150-bafa-d4ffc1c5f74b"
  },
  {
    "state": "lose",
    "amount": "3,688.52",
    "transactionId": "71a7a97e-5cde-420c-94d5-dc7541966ea1"
  },
  {
    "state": "win",
    "amount": "3,852.32",
    "transactionId": "f6e7ecca-161f-45b0-b471-7419892c4a7c"
  },
  {
    "state": "lose",
    "amount": "2,497.24",
    "transactionId": "2543848c-8121-4754-a1c1-1cd1f7c6f1da"
  },
  {
    "state": "lose",
    "amount": "2,978.07",
    "transactionId": "2d6333c0-e388-4323-9844-b42b4ed700d7"
  },
  {
    "state": "win",
    "amount": "3,512.21",
    "transactionId": "b562a887-d200-4976-8b7c-31c128b3e6be"
  },
  {
    "state": "win",
    "amount": "3,591.78",
    "transactionId": "45d47d82-132e-40f3-8e3e-6046294fd3cc"
  },
  {
    "state": "lose",
    "amount": "3,492.69",
    "transactionId": "9ae5852f-29cb-4d3a-bb99-bd580e765150"
  },
  {
    "state": "lose",
    "amount": "2,000.60",
    "transactionId": "1a79bf1d-57e5-4810-bf86-35853a48b744"
  },
  {
    "state": "lose",
    "amount": "2,225.96",
    "transactionId": "768a4efc-f530-462e-9f90-9854fb29651a"
  },
  {
    "state": "lose",
    "amount": "2,702.24",
    "transactionId": "f6fd7676-1c5c-4f8c-864b-84a862670ca6"
  },
  {
    "state": "lose",
    "amount": "2,579.74",
    "transactionId": "263fbd18-ef62-45f8-b251-3ff54cfa62c7"
  },
  {
    "state": "lose",
    "amount": "3,894.05",
    "transactionId": "883bc590-f258-48b4-b63f-7661defeba81"
  },
  {
    "state": "lose",
    "amount": "2,605.04",
    "transactionId": "44da7e8d-77fa-4372-ab3f-73aebc462118"
  },
  {
    "state": "lose",
    "amount": "2,987.03",
    "transactionId": "17d5da56-d021-4a5b-9871-bb80e362e7cb"
  },
  {
    "state": "win",
    "amount": "3,344.21",
    "transactionId": "8ec88493-2ab9-47f1-9714-0edf6292d51c"
  },
  {
    "state": "lose",
    "amount": "2,511.11",
    "transactionId": "e771970e-8c4a-48b9-a74e-23dabbdd29c5"
  },
  {
    "state": "win",
    "amount": "3,222.41",
    "transactionId": "17ad6dec-49db-4a9b-9386-185a5d37a295"
  },
  {
    "state": "win",
    "amount": "3,758.16",
    "transactionId": "9799c097-5e80-45f0-a3ee-7e47a3da359b"
  },
  {
    "state": "lose",
    "amount": "3,737.38",
    "transactionId": "bddd7b3a-24fb-43c0-8751-fd8b265554b6"
  },
  {
    "state": "win",
    "amount": "2,296.34",
    "transactionId": "3a066b7e-aee2-46d4-9f10-98e7bb95870e"
  },
  {
    "state": "lose",
    "amount": "3,408.50",
    "transactionId": "f9faaa2d-ad80-47ae-94cb-9db3897f4508"
  },
  {
    "state": "win",
    "amount": "2,202.37",
    "transactionId": "013964af-df16-4577-8238-e5981b7b907b"
  },
  {
    "state": "win",
    "amount": "2,467.27",
    "transactionId": "7860364c-2679-413f-a6e8-8c2e5e062bba"
  },
  {
    "state": "win",
    "amount": "1,109.85",
    "transactionId": "71888d22-1a25-4908-a8e9-27583b77bea0"
  },
  {
    "state": "lose",
    "amount": "1,650.06",
    "transactionId": "9f6060f3-36e7-4a91-a588-05fb41d1a631"
  },
  {
    "state": "lose",
    "amount": "1,530.26",
    "transactionId": "ed8392d5-72f3-42ad-96bb-d7957064a3e4"
  },
  {
    "state": "lose",
    "amount": "2,706.42",
    "transactionId": "4a6ec1e5-2b40-441b-804e-8b85e7a45343"
  },
  {
    "state": "lose",
    "amount": "1,079.19",
    "transactionId": "2b5b14da-093b-41f9-914b-80dddc2a992b"
  },
  {
    "state": "win",
    "amount": "2,996.12",
    "transactionId": "9241bb7f-fa11-47e0-8685-7e38568a9574"
  },
  {
    "state": "win",
    "amount": "3,365.13",
    "transactionId": "3c0f64c2-2515-476f-ab9b-e840ba6856d3"
  },
  {
    "state": "lose",
    "amount": "2,049.94",
    "transactionId": "ed8d3b30-15bc-44f5-b07a-0b18a9a54308"
  },
  {
    "state": "lose",
    "amount": "2,173.45",
    "transactionId": "202e2e3b-6437-4255-a500-441f403d0f03"
  },
  {
    "state": "lose",
    "amount": "1,793.48",
    "transactionId": "9e8927f2-beb9-44ef-909e-1bd3ae98e1f1"
  },
  {
    "state": "lose",
    "amount": "3,990.63",
    "transactionId": "8259e020-a780-45c1-a542-6e6840c8526e"
  },
  {
    "state": "win",
    "amount": "2,120.78",
    "transactionId": "af9b5602-c442-4e49-9385-12d671521d46"
  },
  {
    "state": "win",
    "amount": "1,524.59",
    "transactionId": "39542c8a-e49f-4304-916e-cb89b5e66e15"
  },
  {
    "state": "win",
    "amount": "3,826.69",
    "transactionId": "35e58924-d0b0-4492-a7e5-8cb54cd3a57b"
  },
  {
    "state": "win",
    "amount": "2,382.89",
    "transactionId": "7bf98a47-6999-45c0-9752-7dd5eeebad3b"
  },
  {
    "state": "win",
    "amount": "2,559.52",
    "transactionId": "2f52ee6d-9aa6-4593-ad3d-faacd5afcdf2"
  },
  {
    "state": "win",
    "amount": "1,624.66",
    "transactionId": "82e866e3-3a2d-47ec-9f5b-43585e413156"
  },
  {
    "state": "lose",
    "amount": "3,417.70",
    "transactionId": "ba223a5f-214c-48f5-bbd1-865b125860a2"
  },
  {
    "state": "lose",
    "amount": "1,359.84",
    "transactionId": "ee1ac178-f4a2-4d33-b4b8-5357810c3168"
  },
  {
    "state": "win",
    "amount": "2,738.08",
    "transactionId": "e689806a-b36e-4931-afa3-130b8970f8d3"
  },
  {
    "state": "lose",
    "amount": "1,687.70",
    "transactionId": "9940525a-ba1a-44a9-acd6-2f840c65997b"
  },
  {
    "state": "lose",
    "amount": "1,902.45",
    "transactionId": "96ca731d-25ff-4541-ab0b-b23025d804b1"
  },
  {
    "state": "win",
    "amount": "2,196.79",
    "transactionId": "37e8d2af-d1b0-4734-b90c-729d00fd3a80"
  },
  {
    "state": "lose",
    "amount": "2,385.13",
    "transactionId": "dc57fe58-e3b4-4d28-863b-fe403ca58cb2"
  },
  {
    "state": "lose",
    "amount": "2,410.06",
    "transactionId": "d9a41f4e-02f9-481c-8d60-06a0d90ff6c6"
  },
  {
    "state": "win",
    "amount": "3,798.58",
    "transactionId": "9ef15c2d-3e59-47bb-aa8f-63bac66ec8a3"
  },
  {
    "state": "lose",
    "amount": "1,831.45",
    "transactionId": "59d84220-6371-4502-aef9-f1d2a6cf9b62"
  },
  {
    "state": "lose",
    "amount": "3,530.22",
    "transactionId": "643a15bf-fb0a-45af-954a-c57e993708c5"
  },
  {
    "state": "lose",
    "amount": "1,651.79",
    "transactionId": "8f681c93-e6d3-4588-9f38-4aca41c67dac"
  },
  {
    "state": "lose",
    "amount": "1,467.17",
    "transactionId": "e7e164fa-12c1-4ce2-9265-484a9479911f"
  },
  {
    "state": "lose",
    "amount": "3,623.59",
    "transactionId": "ecab7f73-e43b-40a5-ae9e-6add065617d0"
  },
  {
    "state": "win",
    "amount": "2,749.45",
    "transactionId": "6305f562-5095-4318-9de7-b6dc6ebb5ae1"
  },
  {
    "state": "win",
    "amount": "3,188.34",
    "transactionId": "39786720-8335-4057-8565-a62395b3e5e3"
  },
  {
    "state": "win",
    "amount": "1,195.69",
    "transactionId": "25f0f089-016b-482c-a66b-17048112c672"
  },
  {
    "state": "win",
    "amount": "3,006.14",
    "transactionId": "c26baa5e-a002-4770-850e-a0eace6f09b0"
  },
  {
    "state": "lose",
    "amount": "3,233.04",
    "transactionId": "06263124-8350-42f1-99ab-cf07c0f2910c"
  },
  {
    "state": "lose",
    "amount": "3,321.16",
    "transactionId": "648aea6f-2bc5-433b-959e-e0ced2fedff6"
  },
  {
    "state": "lose",
    "amount": "1,836.65",
    "transactionId": "3b6a9ebe-1f51-473d-957a-a924ac7e882f"
  },
  {
    "state": "win",
    "amount": "3,709.27",
    "transactionId": "d0271c2c-f1e1-40ad-a928-b4d9009886c8"
  },
  {
    "state": "lose",
    "amount": "1,471.85",
    "transactionId": "d867b563-abe8-4920-8387-27f942d1bee5"
  },
  {
    "state": "lose",
    "amount": "2,785.82",
    "transactionId": "7ddfa4a1-3b67-4133-82bd-2ee8e78bb5fe"
  },
  {
    "state": "lose",
    "amount": "2,302.01",
    "transactionId": "008a8fae-de62-4162-b4e4-1b802aa4fcc2"
  },
  {
    "state": "lose",
    "amount": "2,195.31",
    "transactionId": "1d6439ef-df3a-4f03-b26a-a7e72547c111"
  },
  {
    "state": "lose",
    "amount": "3,867.35",
    "transactionId": "919c4a0f-3256-44b2-a45c-977e4edef8e0"
  },
  {
    "state": "win",
    "amount": "3,492.25",
    "transactionId": "4ffb0a47-cd39-4b2c-a3c8-6601fa7cedde"
  },
  {
    "state": "lose",
    "amount": "1,581.12",
    "transactionId": "611e774a-c828-40ec-b60e-2aa483ed4a3d"
  },
  {
    "state": "lose",
    "amount": "3,572.93",
    "transactionId": "d844eb6a-f0e7-4394-a122-d3aa5eaf64a0"
  },
  {
    "state": "win",
    "amount": "1,846.99",
    "transactionId": "154e311a-96c7-4ace-b9b3-bf2c59579c8d"
  },
  {
    "state": "lose",
    "amount": "2,051.90",
    "transactionId": "0a4b2ac9-2c53-41ba-aef9-15951bb4d720"
  },
  {
    "state": "win",
    "amount": "3,183.49",
    "transactionId": "16c121bb-f1c7-4974-9a04-aa5eb0836678"
  },
  {
    "state": "lose",
    "amount": "1,948.53",
    "transactionId": "ddd3990a-764b-44d4-88d6-2405801811d8"
  },
  {
    "state": "win",
    "amount": "2,692.82",
    "transactionId": "05414c59-9711-4a53-b693-0e89226bca64"
  },
  {
    "state": "win",
    "amount": "2,775.08",
    "transactionId": "103d30fc-2a95-4cdb-be34-735f4018b351"
  },
  {
    "state": "lose",
    "amount": "1,996.29",
    "transactionId": "43a0cd85-7ff7-4ece-9918-d5fe18c63edf"
  },
  {
    "state": "lose",
    "amount": "1,459.79",
    "transactionId": "5eb85453-3979-4f16-95ca-a2f2c5d43535"
  },
  {
    "state": "win",
    "amount": "2,487.52",
    "transactionId": "9fbb548b-7081-47fc-9d11-3e04293a63cb"
  },
  {
    "state": "win",
    "amount": "3,273.09",
    "transactionId": "86ba08ea-a5a4-46ae-a6e2-033c6bfd6031"
  },
  {
    "state": "lose",
    "amount": "2,534.69",
    "transactionId": "75c05e64-9470-422a-a73b-fc5d6261c82b"
  },
  {
    "state": "win",
    "amount": "1,668.84",
    "transactionId": "b368bd48-a9af-4588-84cc-ff9ac2a929c1"
  },
  {
    "state": "win",
    "amount": "3,451.88",
    "transactionId": "61ad7ee4-4120-4a38-bd98-4984a01391d6"
  },
  {
    "state": "lose",
    "amount": "2,701.98",
    "transactionId": "a33ba455-abce-4b16-8671-b3b6574562e0"
  },
  {
    "state": "lose",
    "amount": "1,953.50",
    "transactionId": "6c7bf9cd-ec23-4323-be5a-2814cb054e81"
  },
  {
    "state": "lose",
    "amount": "3,250.53",
    "transactionId": "edea6d3b-6587-4da3-81ad-9198b2691a63"
  },
  {
    "state": "lose",
    "amount": "1,689.65",
    "transactionId": "16e815de-5b4a-4f3b-9b0d-df65298a1054"
  },
  {
    "state": "win",
    "amount": "1,904.04",
    "transactionId": "8c3e5fab-f22d-49f3-89a1-0d0fe4aa58a6"
  },
  {
    "state": "lose",
    "amount": "3,916.81",
    "transactionId": "b599cf1a-f846-4486-b013-cb000bc8b9b0"
  },
  {
    "state": "win",
    "amount": "1,577.57",
    "transactionId": "fb72e784-1c99-477a-89be-626e159da82a"
  },
  {
    "state": "win",
    "amount": "1,561.26",
    "transactionId": "255557fe-989b-4193-9e5a-7a5fba422a10"
  },
  {
    "state": "lose",
    "amount": "3,835.29",
    "transactionId": "65e22f68-5b80-4387-a1fe-5abe6c9073e4"
  },
  {
    "state": "win",
    "amount": "3,011.64",
    "transactionId": "02cb32d5-6a11-4ec1-984f-4d4a2e0d9641"
  },
  {
    "state": "win",
    "amount": "3,816.03",
    "transactionId": "82ff55cf-7c92-4115-93c3-ca7190d83685"
  },
  {
    "state": "win",
    "amount": "3,370.84",
    "transactionId": "66be4a1f-4389-4e82-8fe7-7956ec42e631"
  },
  {
    "state": "win",
    "amount": "1,533.11",
    "transactionId": "ee5aab9c-99d9-4d95-940c-8419db7e46e4"
  },
  {
    "state": "win",
    "amount": "3,462.55",
    "transactionId": "b64fdab1-be4a-4d28-b1eb-5fe919acf3ea"
  },
  {
    "state": "win",
    "amount": "1,177.58",
    "transactionId": "e827a95c-a438-4774-9d00-01de9767f350"
  },
  {
    "state": "win",
    "amount": "1,177.72",
    "transactionId": "28fb914d-d1eb-487d-b422-bcfb923dcf65"
  },
  {
    "state": "lose",
    "amount": "2,855.59",
    "transactionId": "480d9dae-17ec-423b-90fb-4fffad90b07f"
  },
  {
    "state": "lose",
    "amount": "1,072.78",
    "transactionId": "2b56f280-7d99-4ce7-8f1b-b326cdac03fe"
  },
  {
    "state": "win",
    "amount": "2,058.03",
    "transactionId": "b578e4a8-f25b-4c40-9921-0394aaabeca2"
  },
  {
    "state": "lose",
    "amount": "3,557.00",
    "transactionId": "312916f6-308e-4f96-9c18-8181a504c103"
  },
  {
    "state": "win",
    "amount": "2,914.86",
    "transactionId": "ff71348e-1e5a-4005-867e-4968574db15e"
  },
  {
    "state": "lose",
    "amount": "2,278.79",
    "transactionId": "03713947-f2d7-4db1-b876-addf6c5bd99e"
  },
  {
    "state": "win",
    "amount": "1,002.17",
    "transactionId": "96b6c8b9-68ee-4959-beac-b44cfcba1477"
  },
  {
    "state": "lose",
    "amount": "1,713.45",
    "transactionId": "19bf0fe6-08c5-456a-bf7b-b62ea86655fb"
  },
  {
    "state": "lose",
    "amount": "2,738.25",
    "transactionId": "597e3c21-f1aa-4d44-ad5b-a40e4d65f708"
  },
  {
    "state": "lose",
    "amount": "3,373.11",
    "transactionId": "467e5a4e-2b18-408a-8654-ce1f8d775ca8"
  },
  {
    "state": "lose",
    "amount": "3,822.79",
    "transactionId": "a4c3b812-75d4-4839-b7aa-f1ff8a493508"
  },
  {
    "state": "win",
    "amount": "2,880.57",
    "transactionId": "9ed61133-f876-4142-a71c-e3a7d6358085"
  },
  {
    "state": "win",
    "amount": "2,540.95",
    "transactionId": "64e9c170-fd16-49e0-b951-8130a83ebf64"
  },
  {
    "state": "lose",
    "amount": "2,038.34",
    "transactionId": "49d31519-8b02-46ac-b5cc-9f630d09e0cc"
  },
  {
    "state": "lose",
    "amount": "3,462.95",
    "transactionId": "a62b6548-a85a-4119-9b1e-27d660ae229e"
  },
  {
    "state": "win",
    "amount": "3,686.53",
    "transactionId": "8dd4beb0-783f-4085-b3d8-70925e10eb38"
  },
  {
    "state": "win",
    "amount": "2,676.90",
    "transactionId": "c3ed9cac-9b9d-425f-8ce5-e4afb200fd1d"
  },
  {
    "state": "lose",
    "amount": "1,723.13",
    "transactionId": "c8f9037f-9d89-4c44-aa94-2ea1ba937580"
  },
  {
    "state": "win",
    "amount": "1,000.85",
    "transactionId": "f1961e3c-fac0-4e27-8a08-390f6a4f7522"
  },
  {
    "state": "win",
    "amount": "2,043.12",
    "transactionId": "a86fa69f-9570-43a8-aa9f-c2e325e9fe6f"
  },
  {
    "state": "win",
    "amount": "2,220.89",
    "transactionId": "34a3ac2d-73db-4275-851a-4cd2d4d2d70d"
  },
  {
    "state": "win",
    "amount": "2,255.83",
    "transactionId": "ea5067f4-1713-4015-a4c8-654948e96297"
  },
  {
    "state": "lose",
    "amount": "3,391.12",
    "transactionId": "97bc322f-7e9e-480c-9004-dc17bb2214a7"
  },
  {
    "state": "lose",
    "amount": "3,742.73",
    "transactionId": "ca68bcdb-b743-46b2-b4e0-c60b67528660"
  },
  {
    "state": "win",
    "amount": "1,244.42",
    "transactionId": "64c0c251-1191-4625-aaf5-28da1fe433c0"
  },
  {
    "state": "win",
    "amount": "3,825.29",
    "transactionId": "35fbec9c-6623-49bd-a0a2-fe3495da27b9"
  },
  {
    "state": "lose",
    "amount": "2,955.25",
    "transactionId": "91d8becc-c57e-42c8-b6c0-006decb2839e"
  },
  {
    "state": "lose",
    "amount": "1,414.24",
    "transactionId": "48571d86-8cfe-44a5-b585-d24bf522ca7c"
  },
  {
    "state": "win",
    "amount": "2,284.36",
    "transactionId": "997dc060-5d16-4398-8d88-b34013c30a8b"
  },
  {
    "state": "lose",
    "amount": "1,978.42",
    "transactionId": "adf725ec-981c-40b9-aac7-0d5d71d28a94"
  },
  {
    "state": "lose",
    "amount": "1,215.71",
    "transactionId": "cbf1b18c-0d82-42b6-be97-8ede723b6bbf"
  },
  {
    "state": "win",
    "amount": "3,234.89",
    "transactionId": "de86dbe7-d379-4c2a-bce3-20a27cef6c40"
  },
  {
    "state": "lose",
    "amount": "1,070.09",
    "transactionId": "6616c35e-9a04-4d0e-9e03-db32332087c5"
  },
  {
    "state": "win",
    "amount": "1,106.34",
    "transactionId": "12499a4d-fbe1-484c-8451-ef51a12855de"
  },
  {
    "state": "win",
    "amount": "3,786.47",
    "transactionId": "491b6078-f3c2-493c-aabf-7ddd14f2caf5"
  },
  {
    "state": "win",
    "amount": "3,135.08",
    "transactionId": "10d39323-30cb-4899-bf72-0ddd3aa37801"
  },
  {
    "state": "lose",
    "amount": "1,965.49",
    "transactionId": "0072406c-a164-4c37-a8de-e60f19fafa44"
  },
  {
    "state": "lose",
    "amount": "3,039.16",
    "transactionId": "3912ed9b-522f-4e18-87e3-f2fff95dd5a9"
  },
  {
    "state": "win",
    "amount": "1,435.40",
    "transactionId": "b22e2cbf-9259-409d-aa46-66b9f64cd114"
  },
  {
    "state": "win",
    "amount": "2,028.18",
    "transactionId": "c657b949-9195-4387-adf6-b0b941a10f8a"
  },
  {
    "state": "win",
    "amount": "3,269.87",
    "transactionId": "f41edffb-68e0-4304-9c9b-1848677a5dd2"
  },
  {
    "state": "lose",
    "amount": "1,412.19",
    "transactionId": "3d2ea80e-d078-41c9-8fcc-d38b5dfa5038"
  },
  {
    "state": "lose",
    "amount": "2,926.00",
    "transactionId": "4c9e7253-e40d-4b27-83aa-7c6bb393ed74"
  },
  {
    "state": "win",
    "amount": "1,121.17",
    "transactionId": "fb16f3e3-7199-4f0b-978f-5586e8b0672d"
  },
  {
    "state": "lose",
    "amount": "3,493.24",
    "transactionId": "92b93fcb-f448-4da9-b562-8e20a2a31e9b"
  },
  {
    "state": "win",
    "amount": "1,446.07",
    "transactionId": "5d70cd66-a662-4d09-ac39-8c8223b698bf"
  },
  {
    "state": "win",
    "amount": "3,498.63",
    "transactionId": "d67f9d03-3259-4348-824c-7fbc74ffb001"
  },
  {
    "state": "lose",
    "amount": "1,829.29",
    "transactionId": "a0eac237-11f2-4756-a788-b0cb4c13875d"
  },
  {
    "state": "lose",
    "amount": "1,953.42",
    "transactionId": "a90a1e33-f891-4a7f-92d1-d98cf42f8097"
  },
  {
    "state": "lose",
    "amount": "3,967.28",
    "transactionId": "8d5e0a52-9cf1-4355-a310-e045140ffeb1"
  },
  {
    "state": "lose",
    "amount": "1,704.85",
    "transactionId": "09efd2e5-a667-4a6b-8f4e-94c271837bcf"
  },
  {
    "state": "win",
    "amount": "1,741.69",
    "transactionId": "8ecd2273-e6ff-4426-8ede-5ab2aadc4d2d"
  },
  {
    "state": "lose",
    "amount": "2,339.02",
    "transactionId": "476ef16e-d5a1-4b2b-96f3-ead44f6d256c"
  },
  {
    "state": "lose",
    "amount": "3,784.40",
    "transactionId": "85f3b6f7-411a-489f-aa72-fb8a5a2d75cc"
  },
  {
    "state": "win",
    "amount": "1,417.40",
    "transactionId": "18fd3c0f-fdbc-4b2e-adef-297eed9545e9"
  },
  {
    "state": "lose",
    "amount": "3,829.47",
    "transactionId": "253c4139-80ae-4de0-acc6-aad10e33ddb9"
  },
  {
    "state": "lose",
    "amount": "1,354.61",
    "transactionId": "d759cc9c-8407-40ae-bc1f-8b04963393dd"
  },
  {
    "state": "lose",
    "amount": "2,559.60",
    "transactionId": "e6a3fd28-9afd-475b-bf46-ba397d1fbe0a"
  },
  {
    "state": "win",
    "amount": "2,861.25",
    "transactionId": "ed2efb15-2714-43a2-bb69-9e44aab56065"
  },
  {
    "state": "lose",
    "amount": "3,423.47",
    "transactionId": "e058893a-7f16-4a82-aa09-9a6408e0f3b3"
  },
  {
    "state": "win",
    "amount": "2,631.36",
    "transactionId": "8f430573-ba1e-4027-bfc2-94aae93a6e71"
  },
  {
    "state": "win",
    "amount": "3,430.33",
    "transactionId": "474d2461-b544-4532-a831-1992d5b5adf3"
  },
  {
    "state": "lose",
    "amount": "2,679.48",
    "transactionId": "7986497c-7c01-40bc-8700-bbf1cc30efac"
  },
  {
    "state": "win",
    "amount": "3,757.57",
    "transactionId": "2882103f-1aeb-4ae3-9e0c-4fb56f3823a9"
  },
  {
    "state": "lose",
    "amount": "3,613.47",
    "transactionId": "8f55dd38-42d8-4977-b26b-5fd0f43b562f"
  },
  {
    "state": "lose",
    "amount": "1,564.72",
    "transactionId": "221fe24a-81cb-4fad-bcc1-ccfc7de30635"
  },
  {
    "state": "lose",
    "amount": "2,611.18",
    "transactionId": "00145bf0-09fc-4004-91f1-7f6957537271"
  },
  {
    "state": "win",
    "amount": "2,466.63",
    "transactionId": "a052e434-9c22-4d7d-ae79-12f22fdec934"
  },
  {
    "state": "lose",
    "amount": "3,045.47",
    "transactionId": "37f5d8b6-35bc-4970-ad74-f4c57bec7c0f"
  },
  {
    "state": "lose",
    "amount": "3,675.72",
    "transactionId": "66918058-ea47-47a3-834d-33e8b448b708"
  },
  {
    "state": "win",
    "amount": "2,364.16",
    "transactionId": "562a86af-3fa1-452f-92b1-81a0c3b9bfac"
  },
  {
    "state": "lose",
    "amount": "2,508.25",
    "transactionId": "e25d7b25-88be-4e35-b6d8-01c8ae2c79be"
  },
  {
    "state": "lose",
    "amount": "2,621.33",
    "transactionId": "55689a68-d82f-4aa4-ac06-5680f704ff34"
  },
  {
    "state": "lose",
    "amount": "1,674.60",
    "transactionId": "1e02dfb6-53ea-43c9-b566-41f2929909c8"
  },
  {
    "state": "win",
    "amount": "1,464.60",
    "transactionId": "6104e689-4b0b-4449-aad6-468e56ac5017"
  },
  {
    "state": "win",
    "amount": "1,621.77",
    "transactionId": "74e05a99-74cc-4605-9c86-a11d08d06b20"
  },
  {
    "state": "lose",
    "amount": "2,968.14",
    "transactionId": "d6f5c1f2-bf33-4829-8e6e-100baa773af6"
  },
  {
    "state": "lose",
    "amount": "1,322.85",
    "transactionId": "48bf0ef1-2253-4c0f-bf8a-5428e077fb4e"
  },
  {
    "state": "win",
    "amount": "3,958.54",
    "transactionId": "3d7236ee-c9b6-494f-9425-b7e652b38ed0"
  },
  {
    "state": "lose",
    "amount": "3,064.19",
    "transactionId": "07054cfd-ad13-4389-aabe-3c95c22d3108"
  },
  {
    "state": "lose",
    "amount": "3,358.55",
    "transactionId": "25508ac9-4b7d-4d0c-8ae0-dc13487ca967"
  },
  {
    "state": "win",
    "amount": "2,700.78",
    "transactionId": "072c0d09-6bd2-49b9-b309-cdf77d77e612"
  },
  {
    "state": "lose",
    "amount": "3,793.99",
    "transactionId": "11692651-10f6-4338-bdc9-aea69c54f506"
  },
  {
    "state": "win",
    "amount": "1,073.90",
    "transactionId": "c07bd755-3237-40bf-a440-a93529dfe01d"
  },
  {
    "state": "lose",
    "amount": "2,384.33",
    "transactionId": "d6c7ccd5-c35f-4ac3-9d41-682eb6d4b12a"
  },
  {
    "state": "win",
    "amount": "3,494.12",
    "transactionId": "71f7c480-ea99-4dc1-9b49-306e368fc98c"
  },
  {
    "state": "win",
    "amount": "3,067.76",
    "transactionId": "db6a37bf-24f3-4a34-92d2-cfa2b902f756"
  },
  {
    "state": "win",
    "amount": "2,221.30",
    "transactionId": "b25b1be6-5820-459b-b3b7-d2265f9a78b9"
  },
  {
    "state": "lose",
    "amount": "3,378.09",
    "transactionId": "6b010701-5e42-4f5f-8cc6-dd531d90bdb5"
  },
  {
    "state": "win",
    "amount": "2,954.10",
    "transactionId": "7f2d6ab4-fa92-40ac-b9cc-e1055e53ae36"
  },
  {
    "state": "win",
    "amount": "1,584.94",
    "transactionId": "88bf7e6c-a8bd-453a-bf68-86d73bfa46a4"
  },
  {
    "state": "lose",
    "amount": "1,048.06",
    "transactionId": "7ceecdd4-4d83-47ab-b7ca-c18c4fe03de0"
  },
  {
    "state": "lose",
    "amount": "2,944.04",
    "transactionId": "dd2f9aae-afe2-4aad-b568-ec22412a63a6"
  },
  {
    "state": "win",
    "amount": "1,419.36",
    "transactionId": "8bea3799-fde1-416b-aad3-01908a03f15d"
  },
  {
    "state": "win",
    "amount": "1,747.21",
    "transactionId": "0d98804d-6579-4ef4-856e-672afa5c6323"
  },
  {
    "state": "win",
    "amount": "1,571.09",
    "transactionId": "e2bb1a87-ec6b-407b-8e1c-e96d5e424f85"
  },
  {
    "state": "win",
    "amount": "1,039.97",
    "transactionId": "a8874ac7-06ed-46ce-ae36-b560396412f5"
  },
  {
    "state": "win",
    "amount": "3,522.43",
    "transactionId": "834c52bf-4800-4bed-bd78-a30f856c6a9c"
  },
  {
    "state": "lose",
    "amount": "1,382.71",
    "transactionId": "ad8d8b96-ab34-450b-95de-541db1bd5050"
  },
  {
    "state": "lose",
    "amount": "1,302.67",
    "transactionId": "3e2883e8-eff1-458d-9342-f3ecd9144b6d"
  },
  {
    "state": "win",
    "amount": "3,897.98",
    "transactionId": "fbe92c3c-7adf-41c6-9f50-b45b336e3be7"
  },
  {
    "state": "lose",
    "amount": "3,058.27",
    "transactionId": "adccf876-2092-4eb9-9d21-4cdc1b1a2441"
  },
  {
    "state": "win",
    "amount": "3,240.77",
    "transactionId": "4745bf10-eead-420a-a91a-a178ccffa2fb"
  },
  {
    "state": "lose",
    "amount": "1,823.39",
    "transactionId": "4d68424c-277c-47dc-9929-96988f7bc293"
  },
  {
    "state": "win",
    "amount": "1,174.91",
    "transactionId": "46b825a3-9e31-4fa4-a577-f16234ff3598"
  },
  {
    "state": "lose",
    "amount": "1,461.74",
    "transactionId": "c250ecaf-9f74-46d8-8ab8-d1dfdbdf040e"
  },
  {
    "state": "lose",
    "amount": "2,120.80",
    "transactionId": "96a9f89a-5414-440a-ac50-62823c6973a1"
  },
  {
    "state": "lose",
    "amount": "1,223.30",
    "transactionId": "9946e883-7974-40ae-b88e-729496c893f5"
  },
  {
    "state": "lose",
    "amount": "1,059.70",
    "transactionId": "033b1fba-b95e-4ab3-bb25-f6e30b1fc68a"
  },
  {
    "state": "win",
    "amount": "2,763.44",
    "transactionId": "3ab690ec-1b3a-4dbc-966a-a36f879bdc3f"
  },
  {
    "state": "win",
    "amount": "1,957.55",
    "transactionId": "71f71a0f-ef6a-46ea-9cb1-f7be9754577e"
  },
  {
    "state": "win",
    "amount": "1,086.38",
    "transactionId": "27cc8f23-2ef7-4ecd-b3eb-3303f8f68201"
  },
  {
    "state": "lose",
    "amount": "2,328.29",
    "transactionId": "0e973fde-91be-472f-9987-58dd5c8fba31"
  },
  {
    "state": "win",
    "amount": "1,526.94",
    "transactionId": "2d198674-53df-4149-88bd-c462c798ad23"
  },
  {
    "state": "lose",
    "amount": "2,335.75",
    "transactionId": "d779afdb-1bdd-4e6e-a11a-13396542016f"
  },
  {
    "state": "lose",
    "amount": "3,524.93",
    "transactionId": "56dfed62-4202-4aa7-9d65-60e3eced38e1"
  },
  {
    "state": "lose",
    "amount": "2,455.46",
    "transactionId": "f1e8e16c-2d4e-4006-89f2-c15bf1bcbae7"
  },
  {
    "state": "win",
    "amount": "1,249.36",
    "transactionId": "79df2da0-bf75-4196-98c6-566b76bad05b"
  },
  {
    "state": "win",
    "amount": "2,323.60",
    "transactionId": "0421d57b-f8bc-4437-aca8-845cf13b532e"
  },
  {
    "state": "lose",
    "amount": "2,711.80",
    "transactionId": "e3d99866-7f6d-45ea-be4e-eba2daaa4650"
  },
  {
    "state": "win",
    "amount": "2,051.61",
    "transactionId": "0869daef-6955-4e3d-948a-ec0238476489"
  },
  {
    "state": "win",
    "amount": "3,228.46",
    "transactionId": "56e3d5bc-7789-4ec8-ae19-4c33d122e127"
  },
  {
    "state": "win",
    "amount": "1,822.29",
    "transactionId": "6230d233-f3d0-4051-8637-3d91d976d9c4"
  },
  {
    "state": "lose",
    "amount": "1,222.70",
    "transactionId": "23063dee-4f8e-45cb-a1e6-f872a3e19d8e"
  },
  {
    "state": "lose",
    "amount": "2,005.20",
    "transactionId": "13b2423c-570a-47d9-b2d7-ed29eec2b44b"
  },
  {
    "state": "lose",
    "amount": "1,178.77",
    "transactionId": "9a979226-983e-44a7-822e-b3a3a7443853"
  },
  {
    "state": "win",
    "amount": "3,484.88",
    "transactionId": "4fa16e67-3050-4a4b-bb44-2c7c13b50155"
  },
  {
    "state": "lose",
    "amount": "2,487.93",
    "transactionId": "c5398671-ff5a-4b8f-949f-bde65fed2b4d"
  },
  {
    "state": "lose",
    "amount": "1,546.56",
    "transactionId": "0552964d-e271-4278-85c3-b18ba7d3d60a"
  },
  {
    "state": "win",
    "amount": "2,239.65",
    "transactionId": "09a75502-621a-4ce0-b206-ff50b52675d8"
  },
  {
    "state": "lose",
    "amount": "1,246.94",
    "transactionId": "4e1b21d4-3046-4221-b2f6-a3462285988f"
  },
  {
    "state": "win",
    "amount": "2,589.46",
    "transactionId": "730d5618-d86e-498f-886c-e03ce4eb019c"
  },
  {
    "state": "win",
    "amount": "2,772.83",
    "transactionId": "9878ecb1-4687-4a54-ac45-b59e201162f2"
  },
  {
    "state": "lose",
    "amount": "1,562.45",
    "transactionId": "0dfae029-5b1d-4d2e-b0e3-d247fe373097"
  },
  {
    "state": "lose",
    "amount": "1,410.93",
    "transactionId": "1495aa10-148b-46b5-b10e-b9f443c59a91"
  },
  {
    "state": "win",
    "amount": "1,815.91",
    "transactionId": "358a4007-baa3-4257-ac11-0b0719ce0b28"
  },
  {
    "state": "win",
    "amount": "1,867.18",
    "transactionId": "79140ad9-0c7c-4337-bc0a-1d8f24424e66"
  },
  {
    "state": "lose",
    "amount": "2,621.33",
    "transactionId": "8a3f0cba-5951-4ed9-821c-14e722381c0c"
  },
  {
    "state": "win",
    "amount": "1,308.42",
    "transactionId": "1097ba38-5143-4753-8b4c-d6d1b5bacb42"
  },
  {
    "state": "lose",
    "amount": "2,695.87",
    "transactionId": "46774cab-55f4-4b77-b0e1-2818ed940b33"
  },
  {
    "state": "win",
    "amount": "2,005.75",
    "transactionId": "0857c095-f250-4a5d-891f-11433825ebaf"
  },
  {
    "state": "win",
    "amount": "1,848.99",
    "transactionId": "5b32ea0c-d4d3-43d5-8f5e-347d89db85c9"
  },
  {
    "state": "win",
    "amount": "2,135.32",
    "transactionId": "b3bfba59-f7f2-4edf-b129-198525fae56f"
  },
  {
    "state": "win",
    "amount": "1,421.80",
    "transactionId": "8fb97343-a291-4436-a879-d405a8d21fa9"
  },
  {
    "state": "win",
    "amount": "1,865.92",
    "transactionId": "6f0b50e7-de9c-4785-8a4f-da501107d233"
  },
  {
    "state": "lose",
    "amount": "2,463.36",
    "transactionId": "678a982e-c06c-4057-b667-7dae81591218"
  },
  {
    "state": "win",
    "amount": "1,336.16",
    "transactionId": "6e01fc60-9c7a-4181-81b6-2e46c20f6335"
  },
  {
    "state": "lose",
    "amount": "1,340.75",
    "transactionId": "88d3ff94-07e2-4e9a-8ac8-c0d1df56b2a6"
  },
  {
    "state": "win",
    "amount": "3,629.51",
    "transactionId": "37e6cfcf-a6f2-40fb-8796-8c3267009862"
  },
  {
    "state": "win",
    "amount": "2,432.06",
    "transactionId": "ea5cb057-365b-4275-aec4-e6be066f5508"
  },
  {
    "state": "lose",
    "amount": "1,858.81",
    "transactionId": "d3613bef-98e9-4335-97d3-cb87b12068d0"
  },
  {
    "state": "win",
    "amount": "3,576.50",
    "transactionId": "c237a834-c1be-422f-bdc0-23cd673eab59"
  },
  {
    "state": "win",
    "amount": "3,481.35",
    "transactionId": "47482402-1b8c-4d90-b686-c51b7fca0566"
  },
  {
    "state": "win",
    "amount": "1,837.43",
    "transactionId": "d757bd8c-d0b9-4583-af95-b510df3e58c4"
  },
  {
    "state": "lose",
    "amount": "1,942.57",
    "transactionId": "349b1d1b-1cac-4244-b05c-d88d40cbc0c7"
  },
  {
    "state": "lose",
    "amount": "1,209.36",
    "transactionId": "ff5e8924-a309-4ae1-86bf-ea3edc45e822"
  },
  {
    "state": "lose",
    "amount": "3,091.93",
    "transactionId": "c62ce925-60bf-435b-a88f-5cc941a111cf"
  },
  {
    "state": "lose",
    "amount": "1,983.23",
    "transactionId": "5f37928b-9622-420f-9db9-8455ef5c6589"
  },
  {
    "state": "win",
    "amount": "1,236.22",
    "transactionId": "0a44b616-736e-4ed7-90d5-ee7faf0fe468"
  },
  {
    "state": "lose",
    "amount": "2,337.06",
    "transactionId": "ed6014cb-f4fc-43cb-bfdd-36796c643467"
  },
  {
    "state": "win",
    "amount": "3,678.33",
    "transactionId": "af9e05ff-02e2-4266-85cf-1662d0dbeffb"
  },
  {
    "state": "win",
    "amount": "2,900.52",
    "transactionId": "d2b28e99-b731-4ced-9b1c-edbbabda4491"
  },
  {
    "state": "lose",
    "amount": "1,353.32",
    "transactionId": "5a6181bd-be1c-4f5b-9a3c-4d693d95403d"
  },
  {
    "state": "win",
    "amount": "1,221.58",
    "transactionId": "f86f286e-b7d6-41a5-81a0-cba4692a58e3"
  },
  {
    "state": "win",
    "amount": "1,821.21",
    "transactionId": "9d6689df-e0f5-4192-b482-4bdb9a279407"
  },
  {
    "state": "lose",
    "amount": "2,874.09",
    "transactionId": "cb17e847-2e2a-4d8e-90ed-bf3881f2474e"
  },
  {
    "state": "lose",
    "amount": "1,004.14",
    "transactionId": "04538bf8-3293-42ad-a294-99d1be43338a"
  },
  {
    "state": "win",
    "amount": "2,200.60",
    "transactionId": "a82f40b6-4b7c-4904-b6d8-a024f17e20c5"
  },
  {
    "state": "win",
    "amount": "3,984.78",
    "transactionId": "08a19804-1e24-4a8b-9a91-dfa27ed71c00"
  },
  {
    "state": "lose",
    "amount": "1,659.32",
    "transactionId": "09d384e0-7822-4104-8887-8fee7317770a"
  },
  {
    "state": "win",
    "amount": "2,055.83",
    "transactionId": "efe42aee-3ae1-4bf3-86b4-932f849a4aad"
  },
  {
    "state": "lose",
    "amount": "3,353.23",
    "transactionId": "1d44c1d7-5f63-4e9e-8599-238e99df0332"
  },
  {
    "state": "lose",
    "amount": "1,900.13",
    "transactionId": "02dad3dd-77f5-422d-bad2-33f5041b841b"
  },
  {
    "state": "win",
    "amount": "3,433.74",
    "transactionId": "654cadf3-1b1a-4e51-a039-b4cdc865504c"
  },
  {
    "state": "win",
    "amount": "3,808.87",
    "transactionId": "6a6065e3-d639-46c4-a57c-1d4aaacae989"
  },
  {
    "state": "win",
    "amount": "1,721.49",
    "transactionId": "f251345b-256a-427d-bb50-165a0a637ad1"
  },
  {
    "state": "win",
    "amount": "1,068.70",
    "transactionId": "8af18d85-4a62-4b8e-a4a7-354a3dbec23f"
  },
  {
    "state": "win",
    "amount": "2,974.22",
    "transactionId": "3b01bd50-2e5d-4fc3-807e-62c74d702cec"
  },
  {
    "state": "lose",
    "amount": "2,572.16",
    "transactionId": "d4755faa-a206-4486-8e60-46bdecf4fe35"
  },
  {
    "state": "win",
    "amount": "2,015.50",
    "transactionId": "208698ed-743c-4a61-8f10-192c3e696388"
  },
  {
    "state": "lose",
    "amount": "3,727.05",
    "transactionId": "18e60066-e2d4-40df-89aa-93e31f1e05f7"
  },
  {
    "state": "win",
    "amount": "3,381.57",
    "transactionId": "fcc2e4b5-5063-4585-997d-2b50539ce5e4"
  },
  {
    "state": "lose",
    "amount": "3,789.68",
    "transactionId": "863515e0-9195-4ca9-82b0-30426017b93f"
  },
  {
    "state": "lose",
    "amount": "3,336.91",
    "transactionId": "59d05f99-a908-49a1-8c6a-731a24583f30"
  },
  {
    "state": "lose",
    "amount": "2,570.67",
    "transactionId": "963381b6-e612-4092-86e2-64f8bde8a57f"
  },
  {
    "state": "lose",
    "amount": "1,174.87",
    "transactionId": "0a5c157e-7855-40f0-b598-5026dfd74ac2"
  },
  {
    "state": "lose",
    "amount": "3,991.90",
    "transactionId": "93c7b17a-ceab-4645-afaf-4be73b09956a"
  },
  {
    "state": "lose",
    "amount": "3,917.49",
    "transactionId": "ee14ca0a-0a70-4f8c-8959-00c19dafedd3"
  },
  {
    "state": "win",
    "amount": "1,007.89",
    "transactionId": "699e94c9-f10f-4d77-bce0-15d47a2e3c42"
  },
  {
    "state": "lose",
    "amount": "2,378.01",
    "transactionId": "c5b49d1d-a13e-4ac3-9830-1ca5934d0596"
  },
  {
    "state": "lose",
    "amount": "2,402.25",
    "transactionId": "b67c140f-7889-443d-baa2-64910ce1fb7d"
  },
  {
    "state": "win",
    "amount": "1,448.87",
    "transactionId": "30df3d2e-2f5b-4131-92e7-6412587de21c"
  },
  {
    "state": "win",
    "amount": "3,593.13",
    "transactionId": "5602045d-5797-4d0b-b09b-7af328de67e0"
  },
  {
    "state": "win",
    "amount": "3,682.78",
    "transactionId": "89b89e83-68ac-4082-b484-a11ab640b20b"
  },
  {
    "state": "win",
    "amount": "2,464.02",
    "transactionId": "cdd3124a-c7bd-4b52-860c-8459aa5534d6"
  },
  {
    "state": "win",
    "amount": "1,168.63",
    "transactionId": "b0bd28dd-8e97-485f-9fa2-7df2963285c2"
  },
  {
    "state": "win",
    "amount": "1,488.55",
    "transactionId": "0b6fd745-57d8-44c2-97bc-4c4408aaf529"
  },
  {
    "state": "win",
    "amount": "2,656.30",
    "transactionId": "4905ef7f-3bf6-4dd0-8be8-6c79a28c505b"
  },
  {
    "state": "lose",
    "amount": "3,618.17",
    "transactionId": "ef897f1f-b6e2-44f1-9c21-3656c876dda1"
  },
  {
    "state": "win",
    "amount": "2,365.74",
    "transactionId": "05983a65-71c5-496e-9e80-0adcc8fa8698"
  },
  {
    "state": "win",
    "amount": "1,124.29",
    "transactionId": "c6175721-62d6-4c7f-bb4d-c4f88da8e87b"
  },
  {
    "state": "lose",
    "amount": "3,903.26",
    "transactionId": "15a852a7-c7fa-40be-ab29-a82e51609341"
  },
  {
    "state": "lose",
    "amount": "3,399.14",
    "transactionId": "8d59744d-6c9c-44b1-b364-4a9d9caa0554"
  },
  {
    "state": "win",
    "amount": "2,315.61",
    "transactionId": "31cfa1e7-7c81-4fb2-9d77-70980938f774"
  },
  {
    "state": "win",
    "amount": "1,633.83",
    "transactionId": "d73e44e3-83cd-4825-aaa2-8532cbae64e5"
  },
  {
    "state": "lose",
    "amount": "1,934.29",
    "transactionId": "f7e935cf-54cc-43f8-9f45-3b4a7060196d"
  },
  {
    "state": "win",
    "amount": "1,154.41",
    "transactionId": "b34fa756-68f3-47fe-b201-c86f556a2bf3"
  },
  {
    "state": "lose",
    "amount": "2,256.86",
    "transactionId": "42e72065-f418-4830-9a6c-ad483483c3a5"
  },
  {
    "state": "lose",
    "amount": "3,944.81",
    "transactionId": "ff7ebb2c-4202-408f-b67d-430e4cf456f6"
  },
  {
    "state": "lose",
    "amount": "3,530.52",
    "transactionId": "13240775-a967-4329-990c-0e170803d84c"
  },
  {
    "state": "lose",
    "amount": "1,992.68",
    "transactionId": "0b87428f-d87a-4cb5-8c4b-befd9e45a5c9"
  },
  {
    "state": "lose",
    "amount": "1,939.82",
    "transactionId": "02d5c526-737c-4ff1-9b01-d222eaf3a42f"
  },
  {
    "state": "lose",
    "amount": "3,062.56",
    "transactionId": "a2809741-a8ee-48e4-b0ce-c6e4d301d979"
  },
  {
    "state": "win",
    "amount": "3,296.03",
    "transactionId": "fa0a9c61-d7b6-4b5a-87dc-490ce5873d7a"
  },
  {
    "state": "lose",
    "amount": "2,715.19",
    "transactionId": "0d9dc67f-1d38-44e0-86ba-e937202925a3"
  },
  {
    "state": "lose",
    "amount": "2,032.92",
    "transactionId": "3cea60a5-b45b-4e66-b71c-a48ff8908169"
  },
  {
    "state": "win",
    "amount": "1,340.27",
    "transactionId": "b26a60e5-8003-48e8-866e-4c2af0cb9bde"
  },
  {
    "state": "lose",
    "amount": "2,733.76",
    "transactionId": "bf9dbc5d-77d8-4ba1-94fb-8bbea8d374c8"
  },
  {
    "state": "win",
    "amount": "1,897.97",
    "transactionId": "4fbdc982-c99c-4034-8ba2-bccd7e6ebdbc"
  },
  {
    "state": "lose",
    "amount": "2,309.45",
    "transactionId": "bcb0b323-3170-44c5-872f-b9ae8b49c287"
  },
  {
    "state": "lose",
    "amount": "1,922.92",
    "transactionId": "60ef25ed-e363-42c7-948d-550510c1087b"
  },
  {
    "state": "lose",
    "amount": "2,978.94",
    "transactionId": "d750284e-ccd6-4d17-8713-97b6c29297cd"
  },
  {
    "state": "lose",
    "amount": "1,377.21",
    "transactionId": "ffae2d7e-78bd-4491-92d6-4aec781ef82e"
  },
  {
    "state": "win",
    "amount": "1,846.11",
    "transactionId": "aafbf4c6-3ab1-4481-94b1-c9f259f140d2"
  },
  {
    "state": "win",
    "amount": "1,782.25",
    "transactionId": "d1a53f2f-d404-4378-8944-388c6fff815f"
  },
  {
    "state": "lose",
    "amount": "3,976.49",
    "transactionId": "cc81e686-8b04-4f2e-a0e8-bad4113d0447"
  },
  {
    "state": "lose",
    "amount": "2,076.41",
    "transactionId": "79b27267-0fde-42aa-9507-1917f7e7c3e5"
  },
  {
    "state": "win",
    "amount": "2,573.10",
    "transactionId": "a1520b2f-1052-42ca-bc88-f63dad084ec1"
  },
  {
    "state": "win",
    "amount": "1,136.93",
    "transactionId": "58122416-a6ae-4a40-bd84-54db39002608"
  },
  {
    "state": "win",
    "amount": "1,930.49",
    "transactionId": "9d4b62ef-5ea9-4ff5-80c5-5ad03200d32c"
  },
  {
    "state": "lose",
    "amount": "3,846.51",
    "transactionId": "90c746c9-20fd-4660-b8b8-9be7cd5150b7"
  },
  {
    "state": "lose",
    "amount": "2,373.74",
    "transactionId": "c081ed14-084e-4f86-a53e-2b5d43f1808d"
  },
  {
    "state": "lose",
    "amount": "2,609.23",
    "transactionId": "6065474e-2d55-4043-b671-4d371e7bedca"
  },
  {
    "state": "lose",
    "amount": "2,631.41",
    "transactionId": "16b2367a-6929-4ad0-8886-dd3c9b005915"
  },
  {
    "state": "win",
    "amount": "2,937.96",
    "transactionId": "f92f0813-adb3-4943-8b5f-6f7a8f2cabda"
  },
  {
    "state": "lose",
    "amount": "1,782.85",
    "transactionId": "9a682045-2486-4260-95da-d8a720b957e7"
  },
  {
    "state": "win",
    "amount": "3,282.84",
    "transactionId": "29deb347-9837-4a94-af3b-560beba795aa"
  },
  {
    "state": "lose",
    "amount": "3,175.94",
    "transactionId": "52381862-fc0a-4d0d-a0cc-c82e40658efe"
  },
  {
    "state": "win",
    "amount": "1,464.44",
    "transactionId": "b81e00e2-a36d-4494-8176-f2fcb461d670"
  },
  {
    "state": "lose",
    "amount": "3,749.14",
    "transactionId": "d120bf0e-3d39-45d8-ad5a-052137d58c6f"
  },
  {
    "state": "lose",
    "amount": "2,807.57",
    "transactionId": "af195b3d-fe74-45af-88d7-f947d850227e"
  },
  {
    "state": "win",
    "amount": "1,692.78",
    "transactionId": "c1343069-615a-4830-8315-3e60a1806223"
  },
  {
    "state": "lose",
    "amount": "1,274.46",
    "transactionId": "3836195b-909d-45e6-8a62-172acffc937c"
  },
  {
    "state": "lose",
    "amount": "1,079.33",
    "transactionId": "bd13e34d-2756-4029-b972-e4522b603abd"
  },
  {
    "state": "win",
    "amount": "2,410.11",
    "transactionId": "ab75a3a0-a3da-4025-9040-7745a51b8433"
  },
  {
    "state": "win",
    "amount": "2,529.69",
    "transactionId": "503c9048-52ff-4bca-93ba-a1ba5b1788f6"
  },
  {
    "state": "lose",
    "amount": "1,045.97",
    "transactionId": "2cf026f0-a79a-483e-9ad7-2da6f8760e93"
  },
  {
    "state": "win",
    "amount": "1,426.72",
    "transactionId": "fae460b9-3db4-442b-9a4b-c821d238824c"
  },
  {
    "state": "win",
    "amount": "1,134.81",
    "transactionId": "40598e9e-1229-4db3-a97c-125f1d912af2"
  },
  {
    "state": "win",
    "amount": "3,046.12",
    "transactionId": "ac5b1219-96e7-4b56-a656-5e346d9fc382"
  },
  {
    "state": "lose",
    "amount": "3,258.05",
    "transactionId": "1e849df6-2238-4049-a740-f7531152bb1f"
  },
  {
    "state": "lose",
    "amount": "1,001.66",
    "transactionId": "8bea0fcd-4c50-49b2-b829-9386c3cee9e7"
  },
  {
    "state": "lose",
    "amount": "2,137.53",
    "transactionId": "a23a2d99-fd06-482e-bd9d-f1dbfdb4a3f2"
  },
  {
    "state": "win",
    "amount": "3,118.33",
    "transactionId": "f3ab6793-3ea5-42e8-a9d0-9c5ee535035b"
  },
  {
    "state": "win",
    "amount": "2,696.22",
    "transactionId": "6488ce02-abbc-4bc6-96f1-e2bd05a7b8b6"
  },
  {
    "state": "lose",
    "amount": "1,348.23",
    "transactionId": "58264e08-b36f-46b9-8fee-37f06b310f77"
  },
  {
    "state": "lose",
    "amount": "3,293.98",
    "transactionId": "9569a6d6-b724-46c0-9f03-a7cf1dec36ff"
  },
  {
    "state": "lose",
    "amount": "3,787.62",
    "transactionId": "66984576-d823-48f6-8e74-21600fa625e7"
  },
  {
    "state": "win",
    "amount": "2,218.22",
    "transactionId": "8d0fb591-6203-4053-802a-9fb4d955808a"
  },
  {
    "state": "win",
    "amount": "1,463.86",
    "transactionId": "09c403d5-6619-422a-8e60-e4ac8c3d6e0a"
  },
  {
    "state": "lose",
    "amount": "1,541.92",
    "transactionId": "0bb833dd-c861-4922-a1f0-55dd851b713f"
  },
  {
    "state": "lose",
    "amount": "1,483.24",
    "transactionId": "9fa5a042-6c1e-46c0-a5f7-231caefad456"
  }
]`
)
