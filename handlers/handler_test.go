package handlers

import (
	"encoding/json"
	"github.com/SaCavid/simple-task/models"
	"github.com/labstack/echo"
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

func (h *Server) benchmarkNotRegistered(e *echo.Echo, msg string) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(msg))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Source-type", "server")
	req.Header.Set("Authorization", "only_for_testing_benchmark")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

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
		h.benchmarkNotRegistered(e, string(d))
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
	//    amount:'{{floating(0, 100, 2, "0,0.00")}}',
	//    transactionId:'{{guid()}}'
	//  }
	//]
	randomMsg = `[
  {
    "state": "lose",
    "amount": 88.88,
    "transactionId": "673c4749-af15-47ef-a477-e19ecca5030b"
  },
  {
    "state": "win",
    "amount": 49.9,
    "transactionId": "b7fbe724-c6da-4a2f-b83b-fdf075b35a95"
  },
  {
    "state": "lose",
    "amount": 17.92,
    "transactionId": "55500f97-0ea4-481d-8dbd-fcf5af63211a"
  },
  {
    "state": "win",
    "amount": 67.83,
    "transactionId": "b0575727-31fa-4d22-ab70-237b50ea4d24"
  },
  {
    "state": "win",
    "amount": 25.25,
    "transactionId": "75042fcc-f4a0-4233-8edd-36fca953afed"
  },
  {
    "state": "win",
    "amount": 33.82,
    "transactionId": "77fe4255-7337-4449-8707-ba4f6b04196b"
  },
  {
    "state": "lose",
    "amount": 84.58,
    "transactionId": "a2098078-a26d-41c4-aa05-af0ce80c617e"
  },
  {
    "state": "lose",
    "amount": 9.3,
    "transactionId": "9578f263-a622-4094-aff3-6995b2a8fe67"
  },
  {
    "state": "lose",
    "amount": 92.16,
    "transactionId": "9339301f-b779-447c-af51-e5395fd5493c"
  },
  {
    "state": "win",
    "amount": 78.41,
    "transactionId": "8f50c96f-29ad-4831-8b20-f98abfd8f5f7"
  },
  {
    "state": "lose",
    "amount": 97.4,
    "transactionId": "673b3a76-a720-48f7-8a8f-c09cab145385"
  },
  {
    "state": "lose",
    "amount": 88.83,
    "transactionId": "d62257c4-17bc-4cf5-b7a6-b5e559623dd6"
  },
  {
    "state": "win",
    "amount": 77.69,
    "transactionId": "6483e36c-52dd-401e-8a07-cecbb5e4d464"
  },
  {
    "state": "lose",
    "amount": 36.62,
    "transactionId": "cb7a7148-1a1b-4cf0-9432-d947f91f62e5"
  },
  {
    "state": "lose",
    "amount": 0.22,
    "transactionId": "64260628-6613-4deb-8226-55b07b104354"
  },
  {
    "state": "win",
    "amount": 93.95,
    "transactionId": "533638dc-7459-4b23-9209-aec89f2ee4e0"
  },
  {
    "state": "win",
    "amount": 86.84,
    "transactionId": "63219934-763c-4270-bd65-73d9c4806d2c"
  },
  {
    "state": "lose",
    "amount": 45.57,
    "transactionId": "875573f5-d1a2-48a5-9c0e-fe68791a2eea"
  },
  {
    "state": "win",
    "amount": 95.29,
    "transactionId": "24dc00c1-f82d-4dfd-b0e1-c9a6ce676051"
  },
  {
    "state": "win",
    "amount": 55.24,
    "transactionId": "d42b90a6-ce84-4f5c-b421-0d211a54a9ba"
  },
  {
    "state": "lose",
    "amount": 34.06,
    "transactionId": "132f0639-6d65-4123-a06b-00de6ab741ac"
  },
  {
    "state": "win",
    "amount": 85.73,
    "transactionId": "00edf77a-f3e3-43ac-8838-36bb88b5c770"
  },
  {
    "state": "win",
    "amount": 88.51,
    "transactionId": "55d55334-48dd-436b-86b8-45256ee3030d"
  },
  {
    "state": "lose",
    "amount": 33.37,
    "transactionId": "141464ed-b59a-44cc-b616-f67b1e92bb09"
  },
  {
    "state": "win",
    "amount": 12.99,
    "transactionId": "de319589-6469-44b8-a9e4-f414aa9d4631"
  },
  {
    "state": "lose",
    "amount": 89.63,
    "transactionId": "b433b39d-7237-497b-8ad2-4ec6c8a18930"
  },
  {
    "state": "lose",
    "amount": 83.22,
    "transactionId": "03e542e2-fcdd-4105-94aa-6e5b5348ea13"
  },
  {
    "state": "win",
    "amount": 88.51,
    "transactionId": "eef6532b-540e-48e6-95e1-36ae6b2949ee"
  },
  {
    "state": "win",
    "amount": 14.63,
    "transactionId": "77cf7410-65cd-4321-b45b-a8ea72898519"
  },
  {
    "state": "lose",
    "amount": 67.43,
    "transactionId": "72b08d13-2435-4106-95b9-f857ff787290"
  },
  {
    "state": "lose",
    "amount": 37.29,
    "transactionId": "b7b3d14e-79a4-44d3-951e-67dad3d1e2bd"
  },
  {
    "state": "lose",
    "amount": 46.05,
    "transactionId": "21ccf313-1516-4403-bdd9-8a26e139535d"
  },
  {
    "state": "lose",
    "amount": 9.42,
    "transactionId": "702d2ab1-d4df-4a23-8c7f-cfdb4688a2fd"
  },
  {
    "state": "lose",
    "amount": 5.17,
    "transactionId": "af9f9c62-9cac-4aa1-8352-07ae3af30519"
  },
  {
    "state": "win",
    "amount": 54.82,
    "transactionId": "5715ca49-386d-47c9-9b84-6814f1ea90d7"
  },
  {
    "state": "win",
    "amount": 28.64,
    "transactionId": "60eb1f7a-5cdc-468b-b791-5db21ef83851"
  },
  {
    "state": "lose",
    "amount": 28.97,
    "transactionId": "78b8e495-3446-43d1-bd04-810eba1bacd8"
  },
  {
    "state": "win",
    "amount": 13.7,
    "transactionId": "8947698a-8849-42c7-ab53-65ee3fffa36b"
  },
  {
    "state": "win",
    "amount": 71.94,
    "transactionId": "717c9f81-b7ab-45dc-ae7f-217bb21533af"
  },
  {
    "state": "lose",
    "amount": 94.88,
    "transactionId": "2279404f-5f50-4272-b217-265e4c172a9f"
  },
  {
    "state": "win",
    "amount": 30.2,
    "transactionId": "3da8414f-6abf-43fd-aae5-29baedeb8e7a"
  },
  {
    "state": "lose",
    "amount": 72.66,
    "transactionId": "398fc897-b753-4af9-8623-ea3d6bd8451f"
  },
  {
    "state": "lose",
    "amount": 81.77,
    "transactionId": "e9a66bda-287b-431a-9b8d-88a895b7a3c7"
  },
  {
    "state": "lose",
    "amount": 19.6,
    "transactionId": "6a7365e4-ff9d-426f-81ad-6b530d00900e"
  },
  {
    "state": "lose",
    "amount": 49.86,
    "transactionId": "aa7065a2-53cb-4964-a6f1-35fb9a6b0e77"
  },
  {
    "state": "win",
    "amount": 45.57,
    "transactionId": "a083d3cf-3b12-434b-b8f7-8d15453b0231"
  },
  {
    "state": "lose",
    "amount": 4.77,
    "transactionId": "f91330bb-1f10-4579-aa5c-2cacd1a7adc8"
  },
  {
    "state": "win",
    "amount": 73.16,
    "transactionId": "1e5623c5-40ef-4a00-b657-d118d0c07d53"
  },
  {
    "state": "win",
    "amount": 57.92,
    "transactionId": "7381f78e-a910-450c-ae61-7b91c5da28a5"
  },
  {
    "state": "win",
    "amount": 69.76,
    "transactionId": "89d9be3e-acaa-42a0-90ef-8208a48b2b4a"
  },
  {
    "state": "win",
    "amount": 85.55,
    "transactionId": "79ff9671-939c-4861-b5df-f853704d1736"
  },
  {
    "state": "lose",
    "amount": 89.91,
    "transactionId": "98c9ec16-1c2e-4240-8e71-ded2d7489aab"
  },
  {
    "state": "lose",
    "amount": 45.46,
    "transactionId": "03afe266-dcf8-4cc5-a36c-ca797f634280"
  },
  {
    "state": "win",
    "amount": 40.3,
    "transactionId": "550a81e6-0780-4c51-9a11-aa2e5784f959"
  },
  {
    "state": "lose",
    "amount": 61.89,
    "transactionId": "9bf88b76-7483-46c6-aaf7-26e28eae388e"
  },
  {
    "state": "lose",
    "amount": 29.11,
    "transactionId": "7f8ae903-384f-4389-976f-7d9d79b49875"
  },
  {
    "state": "lose",
    "amount": 86.9,
    "transactionId": "b37d7810-be48-4c67-9729-c1489af5d13d"
  },
  {
    "state": "win",
    "amount": 89.38,
    "transactionId": "87e3634c-0269-4e8c-a6d7-dfefbebab732"
  },
  {
    "state": "lose",
    "amount": 40.13,
    "transactionId": "1022c377-4584-4ca4-8671-5154fe7175d6"
  },
  {
    "state": "win",
    "amount": 60.38,
    "transactionId": "29beea45-20f9-43ff-8798-5e41f98aee0f"
  },
  {
    "state": "lose",
    "amount": 61.91,
    "transactionId": "afe8b446-7508-4d22-b140-398c30d64884"
  },
  {
    "state": "lose",
    "amount": 17.36,
    "transactionId": "41d5d983-371f-4ad6-953b-c1507cffcba2"
  },
  {
    "state": "win",
    "amount": 42.22,
    "transactionId": "ac4465c9-eb78-4cc1-8317-d81a076bb07e"
  },
  {
    "state": "win",
    "amount": 67.96,
    "transactionId": "00dfc8ab-a590-468c-906a-176a60fd7992"
  },
  {
    "state": "win",
    "amount": 18.89,
    "transactionId": "fe30a94f-bbf5-46cd-bb03-2390ea12a81c"
  },
  {
    "state": "lose",
    "amount": 82.33,
    "transactionId": "3c27834f-18e7-4361-bb2e-55fd592f1a70"
  },
  {
    "state": "lose",
    "amount": 72.01,
    "transactionId": "445b90df-bfed-45b4-a1b1-17344e4f69d9"
  },
  {
    "state": "win",
    "amount": 15.55,
    "transactionId": "ca5b31ca-714f-43b9-ad7e-a7d945712bbb"
  },
  {
    "state": "lose",
    "amount": 65.98,
    "transactionId": "222fc118-d163-4979-975b-b6d34a1fc705"
  },
  {
    "state": "lose",
    "amount": 12.85,
    "transactionId": "c5d0c6e5-ab2e-498b-8e1c-565a7d1d5ce7"
  },
  {
    "state": "win",
    "amount": 56.41,
    "transactionId": "e355992b-0e0d-4e14-b021-4345b8c24aed"
  },
  {
    "state": "lose",
    "amount": 64.07,
    "transactionId": "7f8b7afe-25ed-4f95-9672-91adbd17a1c2"
  },
  {
    "state": "lose",
    "amount": 90.13,
    "transactionId": "b923ac2d-378d-47a3-868f-2a69d3212185"
  },
  {
    "state": "win",
    "amount": 81.55,
    "transactionId": "9c4d0a3e-4178-4de0-a5bc-9ccb85c82182"
  },
  {
    "state": "win",
    "amount": 91.23,
    "transactionId": "a150cebe-1b21-4c95-bcf9-9005b373a186"
  },
  {
    "state": "win",
    "amount": 6.31,
    "transactionId": "af0c5c46-ee39-46ff-83b1-e14060bd2da8"
  },
  {
    "state": "win",
    "amount": 57.11,
    "transactionId": "6698946f-92df-450e-9fcd-dda511a317c5"
  },
  {
    "state": "lose",
    "amount": 77.1,
    "transactionId": "07785f0a-8fc6-4a82-aea9-7e6ab9f2cbe1"
  },
  {
    "state": "lose",
    "amount": 42.41,
    "transactionId": "4f48f474-e329-4356-9eb2-9443eae5cd04"
  },
  {
    "state": "win",
    "amount": 37.62,
    "transactionId": "9f8429e3-ba13-4f83-97d5-949e99901c78"
  },
  {
    "state": "lose",
    "amount": 71.99,
    "transactionId": "7de31d23-d22e-4b8c-bc7b-bea36466adb1"
  },
  {
    "state": "win",
    "amount": 13.81,
    "transactionId": "db5ca0de-ac8f-452b-b13d-631af61607ff"
  },
  {
    "state": "win",
    "amount": 36.75,
    "transactionId": "8c65d8e3-2e02-49ff-9140-65b238ec644b"
  },
  {
    "state": "win",
    "amount": 51.44,
    "transactionId": "292f057c-ad1d-4e90-8baf-1416cb218817"
  },
  {
    "state": "win",
    "amount": 93.36,
    "transactionId": "90bcbdba-a99e-4588-96bf-cbdd6166beb9"
  },
  {
    "state": "lose",
    "amount": 12.44,
    "transactionId": "2a981226-18f8-49d2-9871-f002b7572a36"
  },
  {
    "state": "win",
    "amount": 57.55,
    "transactionId": "cd85d59e-e281-4b47-b76b-70e37493b7de"
  },
  {
    "state": "lose",
    "amount": 1.37,
    "transactionId": "c2f4ebc5-d2a5-4264-867a-39dac51eea4e"
  },
  {
    "state": "win",
    "amount": 90.47,
    "transactionId": "1d3871d6-76ef-4659-9ba5-bc2f5d6d462d"
  },
  {
    "state": "lose",
    "amount": 36.4,
    "transactionId": "18719fda-3d4e-49bd-9f4c-39c8769281ab"
  },
  {
    "state": "lose",
    "amount": 88.62,
    "transactionId": "daddd3c3-68e1-4149-8a46-7fb25ec4af45"
  },
  {
    "state": "win",
    "amount": 8.61,
    "transactionId": "0f73ebd0-8836-461b-80ef-2ae1b9354fe5"
  },
  {
    "state": "lose",
    "amount": 54.1,
    "transactionId": "9c0c484a-b182-4bf7-8047-2b24b91b5198"
  },
  {
    "state": "win",
    "amount": 73.93,
    "transactionId": "b0f62138-53c2-4216-b91a-47e792996427"
  },
  {
    "state": "win",
    "amount": 43.58,
    "transactionId": "46de830d-4a48-4541-8495-c3ef9209f2dc"
  },
  {
    "state": "win",
    "amount": 44.98,
    "transactionId": "a02d9a09-423d-4113-8915-a4119bcc44dc"
  },
  {
    "state": "win",
    "amount": 27.17,
    "transactionId": "9f7c1d94-1122-48e4-8cbd-842fe77fbee0"
  },
  {
    "state": "win",
    "amount": 61.13,
    "transactionId": "55fc83d2-725d-40a0-b44b-6e3d07792a34"
  },
  {
    "state": "win",
    "amount": 47.41,
    "transactionId": "98f7f3ed-b4bf-4a0d-b219-09d1975e6515"
  },
  {
    "state": "win",
    "amount": 30.1,
    "transactionId": "57bad5d6-39f5-4774-ac69-6a6515215da8"
  },
  {
    "state": "win",
    "amount": 59.39,
    "transactionId": "84dbe99d-e798-4223-9bf6-fca15e723d5e"
  },
  {
    "state": "lose",
    "amount": 87,
    "transactionId": "21c813bd-a994-47ee-ac4d-46993552261e"
  },
  {
    "state": "win",
    "amount": 71.27,
    "transactionId": "cb191caf-79af-4e15-9759-bd8488354078"
  },
  {
    "state": "win",
    "amount": 61.14,
    "transactionId": "45506acd-1bad-4719-bf2e-a5cf2f93255d"
  },
  {
    "state": "win",
    "amount": 72.34,
    "transactionId": "4e8ebd98-024a-4d02-aef7-afb4a748781d"
  },
  {
    "state": "win",
    "amount": 27.01,
    "transactionId": "ead8c321-3e89-45c1-b6ef-578010607b49"
  },
  {
    "state": "win",
    "amount": 16.92,
    "transactionId": "6590dcff-5a4b-4b57-a329-14949b8e755d"
  },
  {
    "state": "win",
    "amount": 68.78,
    "transactionId": "3f77646e-3d07-410a-bdc9-b27ff63662ae"
  },
  {
    "state": "win",
    "amount": 35.66,
    "transactionId": "f3525548-2002-48b3-abb0-826375a648b6"
  },
  {
    "state": "lose",
    "amount": 36.34,
    "transactionId": "450083a3-3d15-4c85-9dd0-8a16d8ee76e8"
  },
  {
    "state": "lose",
    "amount": 11.54,
    "transactionId": "4bec1adc-a5d2-4601-8865-e6ae3806667e"
  },
  {
    "state": "lose",
    "amount": 58.65,
    "transactionId": "8942b8cc-4b3d-45e0-972b-cecd81bb0119"
  },
  {
    "state": "win",
    "amount": 41.08,
    "transactionId": "a02b7653-a614-4b3c-8a92-cdd1d974541f"
  },
  {
    "state": "win",
    "amount": 52.71,
    "transactionId": "7f365241-8545-4392-87b7-9c3c6983d983"
  },
  {
    "state": "win",
    "amount": 5.44,
    "transactionId": "a721baa4-46ad-484a-a67d-53f587c1d075"
  },
  {
    "state": "lose",
    "amount": 49.52,
    "transactionId": "f22a8f28-48c9-4f62-9438-5a00a7b336ff"
  },
  {
    "state": "lose",
    "amount": 14.01,
    "transactionId": "ee8a0d19-824e-49db-b904-079434a50dcc"
  },
  {
    "state": "win",
    "amount": 74.4,
    "transactionId": "a989ea8e-8f36-4ccb-92b0-a69f9805cb28"
  },
  {
    "state": "lose",
    "amount": 57.6,
    "transactionId": "c996af5c-e961-49d3-bd5a-8b86f49b9973"
  },
  {
    "state": "win",
    "amount": 84.44,
    "transactionId": "f1ea8d7a-ed08-4cd1-b000-25cca323d9e8"
  },
  {
    "state": "lose",
    "amount": 81.61,
    "transactionId": "1e70a746-6819-40b8-ba41-0fbc32009ea0"
  },
  {
    "state": "lose",
    "amount": 50.27,
    "transactionId": "3ddcd839-d029-40a8-bb1a-8e5a2d6f04be"
  },
  {
    "state": "win",
    "amount": 53.2,
    "transactionId": "a8af1fb7-62b0-4d96-b5ff-eab9d4964973"
  },
  {
    "state": "win",
    "amount": 10.88,
    "transactionId": "579a334d-13b3-4a5b-9dc9-7d035002fc7c"
  },
  {
    "state": "win",
    "amount": 8.63,
    "transactionId": "b8a5e4f5-a748-4828-96e6-8880435c7aa9"
  },
  {
    "state": "lose",
    "amount": 88.22,
    "transactionId": "ee0d5bb5-2d18-423b-99c7-4820d5e3ec77"
  },
  {
    "state": "lose",
    "amount": 80.97,
    "transactionId": "ac7a1d01-4959-44eb-9f82-61290fd2acd5"
  },
  {
    "state": "lose",
    "amount": 74.96,
    "transactionId": "d4938496-1497-40c9-afa3-763dd386ddcb"
  },
  {
    "state": "lose",
    "amount": 12.07,
    "transactionId": "b2a7ced7-38ee-45bf-8098-c816e090c32c"
  },
  {
    "state": "lose",
    "amount": 65.01,
    "transactionId": "b50e24b4-94f4-4e10-a22f-832f215fdd3e"
  },
  {
    "state": "win",
    "amount": 83.31,
    "transactionId": "e17b2848-20dc-4791-87e2-1bf24c740e3c"
  },
  {
    "state": "lose",
    "amount": 41.31,
    "transactionId": "58ff4847-d8df-443e-8b44-58804371a8ec"
  },
  {
    "state": "lose",
    "amount": 32.21,
    "transactionId": "e18a9640-67b0-436c-84cb-cf7a27f2edbe"
  },
  {
    "state": "win",
    "amount": 14.5,
    "transactionId": "0470c40b-1648-4023-ab00-03dc757851e2"
  },
  {
    "state": "win",
    "amount": 60.49,
    "transactionId": "914052ac-fbec-453e-a142-96bdbd131cdf"
  },
  {
    "state": "lose",
    "amount": 60.83,
    "transactionId": "cf12597d-fa2b-42dc-851b-cd05c21cf06a"
  },
  {
    "state": "lose",
    "amount": 43.12,
    "transactionId": "15fdd9de-1d00-41b5-bd29-fa581dd2f2bd"
  },
  {
    "state": "win",
    "amount": 45.55,
    "transactionId": "d887d5a3-6754-417e-af70-d156ef9feea7"
  },
  {
    "state": "lose",
    "amount": 10.37,
    "transactionId": "6a6659d5-60e2-4e56-a2ad-c74284762d47"
  },
  {
    "state": "win",
    "amount": 33.06,
    "transactionId": "ceac253a-ba41-4c1e-93c7-211ce4e5570a"
  },
  {
    "state": "lose",
    "amount": 33,
    "transactionId": "705e55bc-b6f9-423e-af58-dbf0ca741f2e"
  },
  {
    "state": "win",
    "amount": 28.72,
    "transactionId": "bc22d171-1919-445e-a197-3c0dbff7a45a"
  },
  {
    "state": "lose",
    "amount": 52.19,
    "transactionId": "a5ec6d45-e880-43e1-bcf3-1c32733f8c87"
  },
  {
    "state": "lose",
    "amount": 94.66,
    "transactionId": "a8d51691-ce3e-43ff-910d-4fee9d57c7d2"
  },
  {
    "state": "lose",
    "amount": 70.39,
    "transactionId": "23cd3844-db2b-4a35-9fe4-d2342c9bd06d"
  },
  {
    "state": "win",
    "amount": 67.27,
    "transactionId": "85f89250-9ef1-4959-a841-bb932c17344d"
  },
  {
    "state": "win",
    "amount": 83.74,
    "transactionId": "6deb5b08-8dc3-46d6-81f2-e273eb6e915a"
  },
  {
    "state": "win",
    "amount": 64.36,
    "transactionId": "1d529e6e-68ec-4a6b-8165-a772a4d1c837"
  },
  {
    "state": "lose",
    "amount": 38.12,
    "transactionId": "d1a55a0a-cc1a-44c4-898c-8061b19d4c4a"
  },
  {
    "state": "win",
    "amount": 92.79,
    "transactionId": "20eec209-df21-48f6-8d1d-1e94bc6d7066"
  },
  {
    "state": "lose",
    "amount": 28.13,
    "transactionId": "113b06af-383b-46bc-bd06-e64a57a31eec"
  },
  {
    "state": "lose",
    "amount": 12.02,
    "transactionId": "17a1939c-9285-4c1b-a71d-aa3c5686863b"
  },
  {
    "state": "win",
    "amount": 97.31,
    "transactionId": "a96cabb7-3929-4058-ae48-5bd9d29568fb"
  },
  {
    "state": "lose",
    "amount": 63.56,
    "transactionId": "0bfc0679-5462-4cbb-854d-3f176aa02ef7"
  },
  {
    "state": "lose",
    "amount": 41.21,
    "transactionId": "41062110-7ce7-4f49-9d08-ffbbc70e5a92"
  },
  {
    "state": "lose",
    "amount": 33.01,
    "transactionId": "9e5a5389-d906-4789-874d-821ec5a0ea02"
  },
  {
    "state": "win",
    "amount": 63.39,
    "transactionId": "1bcaf6ab-8d5c-4e76-ae28-648fa782471f"
  },
  {
    "state": "win",
    "amount": 17.84,
    "transactionId": "227e7f4e-7cda-4fee-a58d-85fd317e23d1"
  },
  {
    "state": "lose",
    "amount": 19.52,
    "transactionId": "c1c62a92-1aa1-4aa6-b599-b927ed6f64f3"
  },
  {
    "state": "win",
    "amount": 87.58,
    "transactionId": "424684a4-56f0-4157-ac43-8fa265e798bd"
  },
  {
    "state": "lose",
    "amount": 53.06,
    "transactionId": "4b7077ae-1d96-4b11-9846-9d62e1a519d0"
  },
  {
    "state": "lose",
    "amount": 52.06,
    "transactionId": "abc40e69-6c3f-4bd6-923b-022c6a011e87"
  },
  {
    "state": "win",
    "amount": 74.85,
    "transactionId": "76a13689-bcf1-4d3a-a8ec-cc45e07b9eea"
  },
  {
    "state": "lose",
    "amount": 27.5,
    "transactionId": "fca57d3b-1154-4b72-ac86-779aab53d97c"
  },
  {
    "state": "lose",
    "amount": 94.54,
    "transactionId": "1b6a12c7-804a-4a50-a35d-1ca612961779"
  },
  {
    "state": "lose",
    "amount": 34.53,
    "transactionId": "6a608c53-e0d1-43e1-95dd-c4516bf7012f"
  },
  {
    "state": "lose",
    "amount": 85.3,
    "transactionId": "e827a5bf-46d2-4450-bbf6-f3bf68c3b333"
  },
  {
    "state": "win",
    "amount": 49.94,
    "transactionId": "a45756f5-e92a-4e00-b29f-b8e5a6ddb0aa"
  },
  {
    "state": "win",
    "amount": 62.13,
    "transactionId": "62d4df3f-3331-44b0-91bd-aba5b1b0e3cb"
  },
  {
    "state": "lose",
    "amount": 95.28,
    "transactionId": "8be4cb2d-3f65-4302-a568-b2eea5c1b876"
  },
  {
    "state": "lose",
    "amount": 0.03,
    "transactionId": "b9fa1155-2bc5-4320-a756-96c9ea3b9e1a"
  },
  {
    "state": "lose",
    "amount": 65.45,
    "transactionId": "3fc59186-1616-4147-901d-7addf84dcd56"
  },
  {
    "state": "win",
    "amount": 1.95,
    "transactionId": "f74375ef-56ee-453a-b24f-c7cfe2531e00"
  },
  {
    "state": "win",
    "amount": 28.76,
    "transactionId": "9defdbed-0147-429a-91d3-311c8128895a"
  },
  {
    "state": "win",
    "amount": 28.18,
    "transactionId": "445ba58b-72e0-4621-a8c0-a8faf9e0d169"
  },
  {
    "state": "win",
    "amount": 90.23,
    "transactionId": "e6279daa-2b30-4e2b-be0e-b512c3decb95"
  },
  {
    "state": "win",
    "amount": 97.46,
    "transactionId": "0a3ad56c-2452-4dc2-aa97-5d16844d78b1"
  },
  {
    "state": "win",
    "amount": 45.88,
    "transactionId": "8b14d72b-a385-41fb-80ca-d011c2207322"
  },
  {
    "state": "lose",
    "amount": 68.27,
    "transactionId": "d47b4002-b61f-44f6-913b-e1c4ef4553dd"
  },
  {
    "state": "win",
    "amount": 17.2,
    "transactionId": "b30cd9ad-0841-469a-ac40-d62e6c07f953"
  },
  {
    "state": "win",
    "amount": 50.87,
    "transactionId": "c07beea7-e24b-49a8-b8be-2e34e8b3e8c0"
  },
  {
    "state": "win",
    "amount": 14.94,
    "transactionId": "8592ea03-a1ab-4d92-be6b-66e648e5d3d0"
  },
  {
    "state": "lose",
    "amount": 63.07,
    "transactionId": "5b370504-83f6-4974-b203-7f0ad9d91312"
  },
  {
    "state": "win",
    "amount": 76.9,
    "transactionId": "7be8f6e3-cead-4a38-b39a-766fb351888c"
  },
  {
    "state": "lose",
    "amount": 33.71,
    "transactionId": "2df0e05d-8f84-485f-a3b2-ff877e226e3a"
  },
  {
    "state": "win",
    "amount": 90.53,
    "transactionId": "47716d39-f2fd-4746-9465-cfdcf0544898"
  },
  {
    "state": "lose",
    "amount": 2.18,
    "transactionId": "a2e652df-c111-42c8-b89f-f0fab19a3c07"
  },
  {
    "state": "lose",
    "amount": 96.3,
    "transactionId": "21caae96-a316-4769-b275-9486e371b9cf"
  },
  {
    "state": "win",
    "amount": 45,
    "transactionId": "beb5153d-d2df-4d41-adcc-596cf0487cf2"
  },
  {
    "state": "win",
    "amount": 51,
    "transactionId": "2aea4b18-5a00-4bba-9c2b-fb8fd4de6133"
  },
  {
    "state": "win",
    "amount": 82.51,
    "transactionId": "4b352ea5-bb55-4f32-970a-c0bb69dd9324"
  },
  {
    "state": "lose",
    "amount": 22.53,
    "transactionId": "2057c12f-1f24-489e-b452-a73f1431e27e"
  },
  {
    "state": "win",
    "amount": 67.6,
    "transactionId": "24d1dff8-1cd0-46fa-8f55-99333b908558"
  },
  {
    "state": "lose",
    "amount": 62.15,
    "transactionId": "191123f8-e3bf-4810-bdc3-a1bb6c9f7562"
  },
  {
    "state": "win",
    "amount": 57.86,
    "transactionId": "2e9a2d6f-6e5a-4121-8c53-0ee3ad891a86"
  },
  {
    "state": "win",
    "amount": 6.98,
    "transactionId": "8e8dc835-c0d9-4b36-aa01-678ffd340283"
  },
  {
    "state": "win",
    "amount": 48.88,
    "transactionId": "b20eb646-d254-461b-a462-7519d75ed70a"
  },
  {
    "state": "win",
    "amount": 83.3,
    "transactionId": "db371fa4-d8a3-404a-b4c3-49275076b555"
  },
  {
    "state": "lose",
    "amount": 63.72,
    "transactionId": "2f26f4a4-d536-458f-8b8d-1e9b319e99f8"
  },
  {
    "state": "lose",
    "amount": 25.56,
    "transactionId": "713a4140-00ac-4f0c-a205-ac61ffb93f88"
  },
  {
    "state": "win",
    "amount": 76.41,
    "transactionId": "4f1a0c17-8bdc-497a-9ed4-f3561d66a5fe"
  },
  {
    "state": "win",
    "amount": 86.19,
    "transactionId": "0e3b25c0-611d-4c0f-8301-05b25044d059"
  },
  {
    "state": "lose",
    "amount": 46.2,
    "transactionId": "065d224b-64b2-42c0-b75c-62806c96a0b3"
  },
  {
    "state": "win",
    "amount": 99.26,
    "transactionId": "1aaff954-d91f-4a57-9f71-71736065e35a"
  },
  {
    "state": "lose",
    "amount": 30.55,
    "transactionId": "d22ae6d2-f195-4ac3-8d46-37d3f59c4d92"
  },
  {
    "state": "lose",
    "amount": 94.96,
    "transactionId": "184c55be-5b45-4ed7-830a-9a695d0ab8a0"
  },
  {
    "state": "win",
    "amount": 21.15,
    "transactionId": "024cd77c-47dd-407f-9a46-35159579a86a"
  },
  {
    "state": "lose",
    "amount": 7.42,
    "transactionId": "de91b637-ace8-479b-a1f8-36901f3be840"
  },
  {
    "state": "win",
    "amount": 5.04,
    "transactionId": "5599fe7b-2923-4abf-90b3-7572d0863d05"
  },
  {
    "state": "lose",
    "amount": 45.52,
    "transactionId": "0f10432d-2c64-4167-af8e-842d0c74c290"
  },
  {
    "state": "win",
    "amount": 25.32,
    "transactionId": "9017d6e3-8ecf-4b7a-a144-9a2069952194"
  },
  {
    "state": "win",
    "amount": 18.11,
    "transactionId": "d41917db-48e4-4f80-aec8-d3092243d1ef"
  },
  {
    "state": "win",
    "amount": 57.63,
    "transactionId": "899ee613-efdd-4ae8-bd54-f9a10e251c21"
  },
  {
    "state": "win",
    "amount": 89.41,
    "transactionId": "517678fb-0f3e-484b-9640-096de45c46e4"
  },
  {
    "state": "lose",
    "amount": 56.79,
    "transactionId": "72c0b167-ce01-4c1e-8c31-66871f340fc4"
  },
  {
    "state": "win",
    "amount": 20.86,
    "transactionId": "d70d9cc9-548a-426c-8cd6-e191fb895836"
  },
  {
    "state": "lose",
    "amount": 15.33,
    "transactionId": "f9af6f55-d394-4d73-af5e-8bed9d332591"
  },
  {
    "state": "win",
    "amount": 17.27,
    "transactionId": "68174cac-22e5-43d2-9f62-1bbc453450b9"
  },
  {
    "state": "lose",
    "amount": 7.99,
    "transactionId": "777af57b-c14d-46f9-9c5b-8b509bf50ffe"
  },
  {
    "state": "win",
    "amount": 63.46,
    "transactionId": "6aa87104-a824-4e51-9652-103a07cb7f15"
  },
  {
    "state": "win",
    "amount": 73.41,
    "transactionId": "8f825ea1-c596-4167-b18c-fe1552f4a3d3"
  },
  {
    "state": "win",
    "amount": 23.66,
    "transactionId": "4e705385-1965-4c35-974e-c7eaec4e1efa"
  },
  {
    "state": "win",
    "amount": 73.34,
    "transactionId": "db468b4d-c124-4ddc-86d1-588982773fc8"
  },
  {
    "state": "win",
    "amount": 17.86,
    "transactionId": "bff3d635-1bd8-4e38-9ab8-a581eeaebefa"
  },
  {
    "state": "win",
    "amount": 76.15,
    "transactionId": "84e181c8-609f-4b81-852d-5888519db9e9"
  },
  {
    "state": "win",
    "amount": 13.94,
    "transactionId": "9371e448-c3d5-4fc2-977e-dc5115409673"
  },
  {
    "state": "win",
    "amount": 19.75,
    "transactionId": "7b33e465-6551-425f-b75a-95c03c83d5e9"
  },
  {
    "state": "lose",
    "amount": 24.21,
    "transactionId": "52ee9d5a-9c82-4dd0-bf3b-32260cf196d4"
  },
  {
    "state": "win",
    "amount": 9.21,
    "transactionId": "f204d08c-c2dc-4449-b4c2-332e99bee8fc"
  },
  {
    "state": "win",
    "amount": 64.67,
    "transactionId": "2f8cee69-d5d9-49a1-b6df-c574e4c89a33"
  },
  {
    "state": "win",
    "amount": 51.12,
    "transactionId": "84a1f29c-0703-433e-ab14-c276efa8c52a"
  },
  {
    "state": "win",
    "amount": 95.65,
    "transactionId": "5c10e4be-cdc0-4eca-9484-bac0c6f5969d"
  },
  {
    "state": "win",
    "amount": 17.21,
    "transactionId": "c56aa915-586c-4794-b1e1-53c0ed0ee4dc"
  },
  {
    "state": "lose",
    "amount": 6.49,
    "transactionId": "3349804b-4c1c-4bc6-87d6-87c0a94c242d"
  },
  {
    "state": "win",
    "amount": 66.96,
    "transactionId": "85eac238-966b-46b4-9cdf-e79b4e736359"
  },
  {
    "state": "win",
    "amount": 64.49,
    "transactionId": "1cfad132-867a-4fab-b358-5cd00e54d6c0"
  },
  {
    "state": "win",
    "amount": 79.89,
    "transactionId": "3533e62d-13b0-4803-b8f4-3502ef4b8e08"
  },
  {
    "state": "lose",
    "amount": 18.53,
    "transactionId": "d3963e2c-5a8f-4746-9654-938432944111"
  },
  {
    "state": "lose",
    "amount": 94.05,
    "transactionId": "450e6ae9-95eb-4f01-acd3-a0a755fe3618"
  },
  {
    "state": "lose",
    "amount": 41.2,
    "transactionId": "5db03f0f-992a-41b0-99fb-84a79891eae1"
  },
  {
    "state": "win",
    "amount": 78,
    "transactionId": "39f489ef-43ba-4df1-836f-bb242b1bf40c"
  },
  {
    "state": "lose",
    "amount": 79.96,
    "transactionId": "cc4b8f1d-f6bc-45a2-b4b1-725e9d0a4a2b"
  },
  {
    "state": "lose",
    "amount": 68.98,
    "transactionId": "d093734f-3fdf-4758-a575-34ebd41e89d4"
  },
  {
    "state": "lose",
    "amount": 19.92,
    "transactionId": "f3a86bfe-84a4-4cfd-b13c-4b9729390e68"
  },
  {
    "state": "win",
    "amount": 97.98,
    "transactionId": "4f3d7019-205a-409b-9697-fe9ddd045559"
  },
  {
    "state": "lose",
    "amount": 92.86,
    "transactionId": "2c03ce73-7387-4ecb-8f06-d0b5b1783281"
  },
  {
    "state": "lose",
    "amount": 12.6,
    "transactionId": "5a7195e7-600a-40fc-a70b-5712f9fbdcbd"
  },
  {
    "state": "win",
    "amount": 37.04,
    "transactionId": "28ac8b36-a44a-4fdf-9893-8ebc76e06c62"
  },
  {
    "state": "lose",
    "amount": 80.3,
    "transactionId": "ee582b86-d226-40f3-a988-b011b6d14cd9"
  },
  {
    "state": "lose",
    "amount": 65.66,
    "transactionId": "826b1626-6524-405d-9acc-996dac6d6590"
  },
  {
    "state": "win",
    "amount": 9.83,
    "transactionId": "06a2c5ed-2f2b-486a-aa26-66d2ffde1091"
  },
  {
    "state": "lose",
    "amount": 11.71,
    "transactionId": "8c4f30ec-4522-4c31-8a5b-e2a1b633247c"
  },
  {
    "state": "lose",
    "amount": 81.32,
    "transactionId": "088ad334-74ea-4717-9bd1-9070f782f970"
  },
  {
    "state": "lose",
    "amount": 6.86,
    "transactionId": "49561730-0a4e-40a0-8718-05c4affa0995"
  },
  {
    "state": "win",
    "amount": 94.17,
    "transactionId": "e4e900e5-82b7-4a84-9762-525f4f2fe505"
  },
  {
    "state": "win",
    "amount": 83.72,
    "transactionId": "7c7a4e76-711e-4230-bd1f-14c49ec09152"
  },
  {
    "state": "lose",
    "amount": 26.46,
    "transactionId": "eec611c2-39d3-4232-8cce-d507e2369da7"
  },
  {
    "state": "lose",
    "amount": 17.03,
    "transactionId": "b52996a8-ea5a-46e2-bcdf-805a24a52169"
  },
  {
    "state": "lose",
    "amount": 18.02,
    "transactionId": "46cf25b6-1774-40a9-a1a2-c57a89ec8e28"
  },
  {
    "state": "win",
    "amount": 4.14,
    "transactionId": "15f7202e-0d9a-4621-849a-4dab515ef3e8"
  },
  {
    "state": "win",
    "amount": 65.8,
    "transactionId": "3100f84c-2cff-4b82-9d14-b5d0a062df1d"
  },
  {
    "state": "lose",
    "amount": 27.87,
    "transactionId": "569d6ff8-9738-46c3-be1f-3930cde755b7"
  },
  {
    "state": "lose",
    "amount": 87.07,
    "transactionId": "0ae92273-ff0d-4149-a543-67c5b19b3d58"
  },
  {
    "state": "win",
    "amount": 53.48,
    "transactionId": "d14a609c-361e-43ba-bb80-9cf3e3802f55"
  },
  {
    "state": "lose",
    "amount": 69.1,
    "transactionId": "00616920-a118-4d62-9dac-c2ab1593e284"
  },
  {
    "state": "lose",
    "amount": 66.54,
    "transactionId": "7f3901c4-684b-481d-aef6-5a69edbbdd19"
  },
  {
    "state": "win",
    "amount": 69.5,
    "transactionId": "420b81dc-89ae-49b9-a1d8-d63989f1ddd0"
  },
  {
    "state": "lose",
    "amount": 36.42,
    "transactionId": "1ed623b6-5511-473e-9581-ee5a0e78ba6f"
  },
  {
    "state": "lose",
    "amount": 4.61,
    "transactionId": "66e3524b-4d7f-45d5-8ecd-30c393c339eb"
  },
  {
    "state": "win",
    "amount": 68.03,
    "transactionId": "b9af2748-5fc1-46a4-ba5b-e5511735f29a"
  },
  {
    "state": "lose",
    "amount": 99.59,
    "transactionId": "cbe303ae-2918-4ba1-b680-f3231e0af5ae"
  },
  {
    "state": "win",
    "amount": 90.73,
    "transactionId": "5edd4f00-8dec-4b06-b6d1-bbb0826c5ee4"
  },
  {
    "state": "win",
    "amount": 1.57,
    "transactionId": "68e83c0e-7a8b-419d-a2a3-076f8a25d95c"
  },
  {
    "state": "win",
    "amount": 61.56,
    "transactionId": "e98ae7e3-7218-496b-a08b-b42c6b308f67"
  },
  {
    "state": "lose",
    "amount": 81.55,
    "transactionId": "7336ede6-f540-4915-a4fa-2c03a9aa93ff"
  },
  {
    "state": "lose",
    "amount": 99.39,
    "transactionId": "b88ae68b-03f6-41f3-8748-da472c112df9"
  },
  {
    "state": "win",
    "amount": 78.01,
    "transactionId": "f86390b5-99f1-4de2-906c-4f07bd062235"
  },
  {
    "state": "lose",
    "amount": 27.23,
    "transactionId": "c3bf8c5e-5ae4-4cba-91e3-85270d11dc5a"
  },
  {
    "state": "win",
    "amount": 10.67,
    "transactionId": "d3e977aa-a05d-4341-bb22-d90b8d04fe6c"
  },
  {
    "state": "lose",
    "amount": 76.95,
    "transactionId": "94a2d948-95da-4d33-9211-ee01e7a85dcf"
  },
  {
    "state": "lose",
    "amount": 16.9,
    "transactionId": "251d276f-fcd2-4fbd-bddb-76a71576749d"
  },
  {
    "state": "win",
    "amount": 46.59,
    "transactionId": "5f93f5fe-76bc-47e6-a1c2-0749e11036ac"
  },
  {
    "state": "lose",
    "amount": 80.99,
    "transactionId": "7046f689-10f0-40e1-8f6d-a5c26ccf5188"
  },
  {
    "state": "win",
    "amount": 97.32,
    "transactionId": "9f363300-1b7a-46a6-945e-7ae9661ca981"
  },
  {
    "state": "win",
    "amount": 55.75,
    "transactionId": "2047349a-5f52-490b-976c-ce5d8b5586d7"
  },
  {
    "state": "lose",
    "amount": 42.71,
    "transactionId": "4541d647-d1d6-4e7b-aa60-f2c7005aa882"
  },
  {
    "state": "lose",
    "amount": 74.51,
    "transactionId": "c3215e94-594b-42b5-a3c8-c52b427ae6ef"
  },
  {
    "state": "lose",
    "amount": 26.36,
    "transactionId": "f57af175-1d5f-4e80-9059-c2af4f90e443"
  },
  {
    "state": "lose",
    "amount": 59.74,
    "transactionId": "8ef5a8c6-2462-4e43-af21-bd5fa6e2d15c"
  },
  {
    "state": "lose",
    "amount": 38.84,
    "transactionId": "39d0492f-e768-47dc-a7fe-30bb0c5e4867"
  },
  {
    "state": "lose",
    "amount": 82.46,
    "transactionId": "5c6cfbc7-a1f2-4f9d-a3aa-97b1288610ef"
  },
  {
    "state": "win",
    "amount": 22.92,
    "transactionId": "8bde2d19-7381-4401-87b3-5bd2b9892d07"
  },
  {
    "state": "lose",
    "amount": 91.12,
    "transactionId": "15ae49f7-acc5-4f4d-8c97-aa4d23028ebc"
  },
  {
    "state": "lose",
    "amount": 42.75,
    "transactionId": "3522a7d9-3e0e-481d-b06a-5ab2623c3abb"
  },
  {
    "state": "win",
    "amount": 97.15,
    "transactionId": "161515c9-4211-4ff9-8706-5fab48132891"
  },
  {
    "state": "lose",
    "amount": 59.42,
    "transactionId": "7b371fed-7e95-4ce3-a5b1-5b5c4964ba77"
  },
  {
    "state": "lose",
    "amount": 10.38,
    "transactionId": "8bdcf252-be9c-49e8-838b-3648622f97b0"
  },
  {
    "state": "win",
    "amount": 79.72,
    "transactionId": "c08f76b1-aee8-41e2-be63-6e87dc63c116"
  },
  {
    "state": "win",
    "amount": 61.95,
    "transactionId": "c2e78000-2e53-4569-a4c4-913e0bd34745"
  },
  {
    "state": "lose",
    "amount": 93.58,
    "transactionId": "bbba507a-2f79-4007-b1c0-bfead94ef01d"
  },
  {
    "state": "win",
    "amount": 55.24,
    "transactionId": "c937ea32-3053-4955-bc2c-7912ed920397"
  },
  {
    "state": "win",
    "amount": 21.07,
    "transactionId": "ddb27b23-ed82-45ae-8047-6ccd268b7a11"
  },
  {
    "state": "lose",
    "amount": 50,
    "transactionId": "b2831967-39d3-4990-bcd9-4040fb05ee48"
  },
  {
    "state": "win",
    "amount": 4.1,
    "transactionId": "52f00fd7-b7b2-4f2e-a4a6-4f8b7103212c"
  },
  {
    "state": "lose",
    "amount": 79.95,
    "transactionId": "5fc60d2d-1654-4e22-b76e-0121e505bbec"
  },
  {
    "state": "lose",
    "amount": 69.11,
    "transactionId": "c257e138-48e2-4957-af05-b0523532ff99"
  },
  {
    "state": "win",
    "amount": 76.94,
    "transactionId": "3e07252a-7c08-4d50-834c-a6d51b01f71e"
  },
  {
    "state": "lose",
    "amount": 13.66,
    "transactionId": "a9b1440d-1267-43bd-be7a-cd0648d513cc"
  },
  {
    "state": "win",
    "amount": 85.3,
    "transactionId": "c2c3df89-b491-42c8-8acf-3ac5c84b1893"
  },
  {
    "state": "lose",
    "amount": 5.31,
    "transactionId": "f3344566-5820-44e8-8331-df3d7116dadb"
  },
  {
    "state": "win",
    "amount": 63.03,
    "transactionId": "a6665462-ecf0-430e-a7f8-89270d10cfd3"
  },
  {
    "state": "lose",
    "amount": 25.01,
    "transactionId": "e635d703-d04d-46f7-a1ae-d510736998c4"
  },
  {
    "state": "win",
    "amount": 66.12,
    "transactionId": "605d19d1-5ff3-4056-ae82-039503ab9f51"
  },
  {
    "state": "win",
    "amount": 65.66,
    "transactionId": "d0978e22-e81e-4283-be1b-76deff139ca1"
  },
  {
    "state": "win",
    "amount": 13.82,
    "transactionId": "2c534b46-8e34-4838-9fbc-551481d1f88c"
  },
  {
    "state": "win",
    "amount": 84.5,
    "transactionId": "945a42e8-2d81-4d85-a471-8e93e462d3ea"
  },
  {
    "state": "lose",
    "amount": 63.73,
    "transactionId": "70cce11f-3364-46f2-a66c-73de60761d4b"
  },
  {
    "state": "lose",
    "amount": 81.5,
    "transactionId": "f5e2dd57-226b-471e-924c-75dffb417f8c"
  },
  {
    "state": "lose",
    "amount": 81.39,
    "transactionId": "17c9e77a-de42-4923-8cc3-debb41d89f04"
  },
  {
    "state": "win",
    "amount": 69.24,
    "transactionId": "8b6d8a24-e14c-4331-8eff-8445ed93651f"
  },
  {
    "state": "lose",
    "amount": 82.48,
    "transactionId": "acfcd126-8590-44ac-a374-f32ab3b3a323"
  },
  {
    "state": "win",
    "amount": 36.34,
    "transactionId": "e00350c5-1f29-4e91-a56f-6eeaa8234d6c"
  },
  {
    "state": "win",
    "amount": 79.39,
    "transactionId": "7f9ef460-676d-4582-878f-25c31edc436b"
  },
  {
    "state": "win",
    "amount": 80.86,
    "transactionId": "00e4ea37-bc73-4f60-9413-3e67c36fec24"
  },
  {
    "state": "win",
    "amount": 17.96,
    "transactionId": "67bb649a-bd8b-4b90-ba98-2c4fddaf316d"
  },
  {
    "state": "lose",
    "amount": 15.05,
    "transactionId": "f29ddf86-0d6b-4d61-ac42-15c80adcfcf3"
  },
  {
    "state": "win",
    "amount": 30.64,
    "transactionId": "0c077aa4-8148-48e0-8ac5-abcafd819e2e"
  },
  {
    "state": "lose",
    "amount": 86.79,
    "transactionId": "b0f69668-74ac-4909-91c0-74c2df82b484"
  },
  {
    "state": "lose",
    "amount": 87.56,
    "transactionId": "99ddc265-5ca4-40bc-af4e-7874c6e23d3a"
  },
  {
    "state": "lose",
    "amount": 44.93,
    "transactionId": "a81068ab-7aa1-4c02-9f4e-6d4783b8f83a"
  },
  {
    "state": "lose",
    "amount": 2.59,
    "transactionId": "57646893-b478-4803-86ea-fa4acf4e158f"
  },
  {
    "state": "lose",
    "amount": 14.84,
    "transactionId": "5e546ee8-16bc-4375-8e3d-da20e9d683d5"
  },
  {
    "state": "win",
    "amount": 95.17,
    "transactionId": "dc2f7b99-3810-4de2-afea-17133769dc78"
  },
  {
    "state": "lose",
    "amount": 77.21,
    "transactionId": "b1947e03-d75c-4ffc-a694-0090098e3662"
  },
  {
    "state": "lose",
    "amount": 4.78,
    "transactionId": "19a8376d-95bf-41ce-b861-b27c5adb567b"
  },
  {
    "state": "win",
    "amount": 95.4,
    "transactionId": "858c2779-17a5-4675-96d5-b63372d6813d"
  },
  {
    "state": "win",
    "amount": 79.77,
    "transactionId": "a5807ccc-7784-4745-9357-ccfc3370fb36"
  },
  {
    "state": "win",
    "amount": 67.25,
    "transactionId": "666ed4e3-7fdf-4515-b116-e19dfca4f5df"
  },
  {
    "state": "lose",
    "amount": 93.65,
    "transactionId": "208a15d7-116f-4ff5-8c16-fe8afdc3255b"
  },
  {
    "state": "lose",
    "amount": 89.34,
    "transactionId": "d6b50090-78e7-4e16-a309-9ff251cf2fd5"
  },
  {
    "state": "win",
    "amount": 26.89,
    "transactionId": "22a7e458-2894-4180-94f6-7ea05e6fee86"
  },
  {
    "state": "lose",
    "amount": 50.73,
    "transactionId": "c64d5459-e9b6-41d3-bda1-567accdacc06"
  },
  {
    "state": "lose",
    "amount": 87.51,
    "transactionId": "f8ac0f86-f9a0-4825-8e95-8764243ceb5f"
  },
  {
    "state": "lose",
    "amount": 9.14,
    "transactionId": "389cdedd-9539-4b1e-9b7f-736f89e0f5e5"
  },
  {
    "state": "win",
    "amount": 77.07,
    "transactionId": "ccd36a0c-d559-427c-87f3-c63af69e916a"
  },
  {
    "state": "win",
    "amount": 29.98,
    "transactionId": "2b9a8a88-022d-47f7-b469-0bed4c5abb31"
  },
  {
    "state": "win",
    "amount": 67.92,
    "transactionId": "a01cb894-eff2-4803-84f2-165fd0872d6f"
  },
  {
    "state": "win",
    "amount": 4.17,
    "transactionId": "8818214d-b9dd-4f67-ab4b-c6fb14b07c64"
  },
  {
    "state": "win",
    "amount": 95.61,
    "transactionId": "d1d46d95-ac2f-4bd8-b7f2-524eb11e24f3"
  },
  {
    "state": "win",
    "amount": 67.39,
    "transactionId": "debf9b8f-891c-404c-845c-7751109b70b6"
  },
  {
    "state": "lose",
    "amount": 26.71,
    "transactionId": "52334988-8abc-4533-935e-7cb4b0065eab"
  },
  {
    "state": "lose",
    "amount": 30.28,
    "transactionId": "15723218-44d9-413b-ac69-1bd7829faf53"
  },
  {
    "state": "lose",
    "amount": 94.18,
    "transactionId": "a61b2935-f81d-4406-807a-f658290b249b"
  },
  {
    "state": "lose",
    "amount": 20.31,
    "transactionId": "521171ab-1165-46d7-88fe-745a460065f1"
  },
  {
    "state": "lose",
    "amount": 11.3,
    "transactionId": "ac66ca7b-9bc1-4fc7-a4d5-38bcef7ad623"
  },
  {
    "state": "lose",
    "amount": 36.58,
    "transactionId": "831c1048-2339-4158-a1ea-2a0b848dcdfc"
  },
  {
    "state": "lose",
    "amount": 62.02,
    "transactionId": "0b5966fd-c5ec-4c73-b810-55f716ff3fe5"
  },
  {
    "state": "win",
    "amount": 65.74,
    "transactionId": "63b13b81-6c26-4b30-85d8-fadb777f9412"
  },
  {
    "state": "win",
    "amount": 90.76,
    "transactionId": "dd496f3d-2e2e-4d3b-bd0f-74bbd4e3a414"
  },
  {
    "state": "lose",
    "amount": 98.63,
    "transactionId": "e44dd856-ca3a-4a5e-a696-963fb4d06186"
  },
  {
    "state": "win",
    "amount": 46.74,
    "transactionId": "3eb0c81b-bd81-4098-b639-1f837ba90942"
  },
  {
    "state": "lose",
    "amount": 91.37,
    "transactionId": "def8a630-4c37-4884-8a88-a2cdbd557878"
  },
  {
    "state": "lose",
    "amount": 63.93,
    "transactionId": "4d72e756-ce1c-4cf8-81eb-8a08c7de5072"
  },
  {
    "state": "lose",
    "amount": 8.21,
    "transactionId": "b1f9b4a2-c208-4c14-89f2-6c385c936311"
  },
  {
    "state": "win",
    "amount": 75.75,
    "transactionId": "a9f23fc6-95f5-487e-b2dd-a9393aed9ff7"
  },
  {
    "state": "win",
    "amount": 77.72,
    "transactionId": "6db5f835-b1a1-4d80-a681-1fc1c54e9adf"
  },
  {
    "state": "lose",
    "amount": 30.21,
    "transactionId": "120939df-6447-4aeb-97c8-a582a522d4a3"
  },
  {
    "state": "win",
    "amount": 99.6,
    "transactionId": "699772d8-6408-495f-a5aa-886b43847fe2"
  },
  {
    "state": "lose",
    "amount": 66.57,
    "transactionId": "2fc0298d-af12-4ada-803c-34c1d16b9cf6"
  },
  {
    "state": "lose",
    "amount": 2.87,
    "transactionId": "f5e1aa93-c7f7-44df-abf9-006f734912ee"
  },
  {
    "state": "win",
    "amount": 41.31,
    "transactionId": "c56b57ff-d41b-405f-86ab-f3830851ecea"
  },
  {
    "state": "lose",
    "amount": 74.32,
    "transactionId": "451c99e0-83f3-4c2d-82c1-f8066849cd2c"
  },
  {
    "state": "win",
    "amount": 33.46,
    "transactionId": "a0fdcb47-47b3-4d47-a952-1f29a045d901"
  },
  {
    "state": "win",
    "amount": 62.08,
    "transactionId": "63041df5-48a5-4b98-b7d8-47fa987b5046"
  },
  {
    "state": "lose",
    "amount": 6.43,
    "transactionId": "b0eda805-21f9-416e-aacb-3aad38df8225"
  },
  {
    "state": "win",
    "amount": 75.77,
    "transactionId": "64f8a8e0-3601-4fce-8672-687ee6a8f73b"
  },
  {
    "state": "lose",
    "amount": 4.29,
    "transactionId": "fea09afa-143b-4d07-9229-35a334dfe231"
  },
  {
    "state": "win",
    "amount": 6.54,
    "transactionId": "7fad0a8c-84f4-434d-9d41-8ac411874af8"
  },
  {
    "state": "win",
    "amount": 4.62,
    "transactionId": "4034d491-7f3d-4580-99d2-6e37e6f91a49"
  },
  {
    "state": "lose",
    "amount": 61.93,
    "transactionId": "9442bd0d-9a88-4721-b67d-7492ec55da47"
  },
  {
    "state": "lose",
    "amount": 65.95,
    "transactionId": "a80c616e-92a1-4af0-bd61-3ba679f0a3eb"
  },
  {
    "state": "lose",
    "amount": 55.54,
    "transactionId": "cb4441fe-ab83-4b5a-8636-5ece57ddb8b5"
  },
  {
    "state": "win",
    "amount": 7.41,
    "transactionId": "3cd26ecb-647e-4209-83bf-1ca84d249063"
  },
  {
    "state": "lose",
    "amount": 89.76,
    "transactionId": "46beb362-863d-4771-bb92-b828ccacd757"
  },
  {
    "state": "win",
    "amount": 65.18,
    "transactionId": "9b50065c-cd3b-4990-8ce3-f248acf0fdb2"
  },
  {
    "state": "win",
    "amount": 56.96,
    "transactionId": "42fcb4bd-d056-42b1-92d6-ec52f4689616"
  },
  {
    "state": "lose",
    "amount": 50.12,
    "transactionId": "ca075672-713b-42e7-88d0-5cce8edc24b1"
  },
  {
    "state": "win",
    "amount": 25.15,
    "transactionId": "d0fa1458-8432-4f74-9280-ea012cdd903c"
  },
  {
    "state": "lose",
    "amount": 86.76,
    "transactionId": "9ce85639-c95b-4268-876f-3e5dd22294cb"
  },
  {
    "state": "win",
    "amount": 1.71,
    "transactionId": "baf73b74-85a0-4f58-b646-5b98391531b7"
  },
  {
    "state": "lose",
    "amount": 43.31,
    "transactionId": "8ed639fe-4a01-452e-9916-94f5cb08a97f"
  },
  {
    "state": "lose",
    "amount": 92.29,
    "transactionId": "149c3be3-4df5-4248-9b9d-065986c9a9de"
  },
  {
    "state": "win",
    "amount": 26.4,
    "transactionId": "45ac249b-2e8e-493f-91be-8deb2fbf10f5"
  },
  {
    "state": "win",
    "amount": 95.18,
    "transactionId": "eae185bf-4a65-4019-a244-c58f6658ee4d"
  },
  {
    "state": "lose",
    "amount": 18.18,
    "transactionId": "cf91944a-94d8-44cf-9037-071616b4ad37"
  },
  {
    "state": "lose",
    "amount": 85.75,
    "transactionId": "814c2151-83e8-4948-b8f4-9bf09997ba4b"
  },
  {
    "state": "lose",
    "amount": 40.06,
    "transactionId": "8f106302-faac-4e82-b7e7-19088f33efe4"
  },
  {
    "state": "win",
    "amount": 81.66,
    "transactionId": "f961da26-e357-46f9-bef9-e545457fa753"
  },
  {
    "state": "win",
    "amount": 57.59,
    "transactionId": "477a5cb9-dcce-4782-8f89-42cc1da07f41"
  },
  {
    "state": "win",
    "amount": 91.96,
    "transactionId": "f28ab97a-be48-4dd1-9bce-d5f30c861ae2"
  },
  {
    "state": "lose",
    "amount": 36.51,
    "transactionId": "259972b4-b184-448e-8a71-c9dbf24672ea"
  },
  {
    "state": "win",
    "amount": 16.88,
    "transactionId": "84e6df2f-597a-45a8-9287-73311e3a3e7c"
  },
  {
    "state": "lose",
    "amount": 47.32,
    "transactionId": "1f62f75f-52b3-4c53-bb13-1faf16446c60"
  },
  {
    "state": "win",
    "amount": 0.75,
    "transactionId": "22456a11-30b9-4d37-90ca-4c118af4bbed"
  },
  {
    "state": "win",
    "amount": 90.17,
    "transactionId": "6931cfb7-aa9b-423b-8933-bd4a2d298c27"
  },
  {
    "state": "win",
    "amount": 67.82,
    "transactionId": "a4381227-2451-4545-bbaa-2fab2ba0b5e4"
  },
  {
    "state": "win",
    "amount": 96.88,
    "transactionId": "f309fc8e-bea5-4edd-9d43-ecd7dc483686"
  },
  {
    "state": "win",
    "amount": 52.68,
    "transactionId": "0a6fb879-b6cd-4914-991e-a4ec30982b08"
  },
  {
    "state": "lose",
    "amount": 11.74,
    "transactionId": "aa947e05-7278-4419-801c-86e9333fb9e8"
  },
  {
    "state": "win",
    "amount": 8.7,
    "transactionId": "031045a4-60b7-4efa-be38-84f85886f520"
  },
  {
    "state": "lose",
    "amount": 83.82,
    "transactionId": "44712620-3995-464d-8a79-bb535ca364f6"
  },
  {
    "state": "lose",
    "amount": 51.71,
    "transactionId": "69e597a8-225a-422b-8bb1-d672f8d39221"
  },
  {
    "state": "lose",
    "amount": 65.91,
    "transactionId": "cc97b057-3eec-4064-9e3d-921d3f5ee633"
  },
  {
    "state": "lose",
    "amount": 30.88,
    "transactionId": "ab6f6d6c-6012-4aac-aa86-50bb9ba3f0b3"
  },
  {
    "state": "win",
    "amount": 17.72,
    "transactionId": "bbed7ea3-114a-4095-b082-c6c757b3d4fa"
  },
  {
    "state": "lose",
    "amount": 70.02,
    "transactionId": "c4c07385-c0fe-4772-b8b2-9079e3c0de60"
  },
  {
    "state": "lose",
    "amount": 91.41,
    "transactionId": "e6e50f72-0654-4452-886b-47d838e29c9e"
  },
  {
    "state": "win",
    "amount": 83.34,
    "transactionId": "66a8e3cf-0b19-40bd-b2d9-f9d59394bc50"
  },
  {
    "state": "win",
    "amount": 45.6,
    "transactionId": "9815f03f-a06c-4142-93b4-f93722d5105d"
  },
  {
    "state": "win",
    "amount": 36.53,
    "transactionId": "eff0aac8-c1e0-4742-90c1-f224522df15e"
  },
  {
    "state": "lose",
    "amount": 97.21,
    "transactionId": "e015d8d5-13c1-4d1c-aeb3-f93aefacfa84"
  },
  {
    "state": "win",
    "amount": 73.32,
    "transactionId": "7093a04f-325d-4f9e-b7ac-9ff24b47117d"
  },
  {
    "state": "lose",
    "amount": 94.05,
    "transactionId": "078dd2f5-5038-4c52-ba65-09c8a84207af"
  },
  {
    "state": "win",
    "amount": 92,
    "transactionId": "def67d8e-ffe4-4a79-b91d-a9e0f240f02e"
  },
  {
    "state": "win",
    "amount": 84.24,
    "transactionId": "95399e2f-e995-4aff-a007-4bff8fab6391"
  },
  {
    "state": "lose",
    "amount": 56,
    "transactionId": "37e324fa-3322-4e5c-bac5-92f56d3b61d4"
  },
  {
    "state": "win",
    "amount": 50.89,
    "transactionId": "b2d5acfd-178c-4274-8cf2-b84fa324eb0c"
  },
  {
    "state": "win",
    "amount": 96.32,
    "transactionId": "9aa612f4-608e-49e6-8f86-7039d29f695e"
  },
  {
    "state": "win",
    "amount": 53.38,
    "transactionId": "ae9f5296-a022-4654-b7c5-a3f9ee810ea3"
  },
  {
    "state": "lose",
    "amount": 85.57,
    "transactionId": "7bc168f5-e685-4c9b-928a-5dc4920b12e9"
  },
  {
    "state": "lose",
    "amount": 65.54,
    "transactionId": "cc1a0133-f995-40d7-bfe8-8ac56f058517"
  },
  {
    "state": "lose",
    "amount": 55.24,
    "transactionId": "30d5c4d8-496f-4f71-8cff-20783c525d55"
  },
  {
    "state": "win",
    "amount": 88.52,
    "transactionId": "63cd986b-74f7-4435-b894-9654fd5db629"
  },
  {
    "state": "win",
    "amount": 2.78,
    "transactionId": "cf5f8571-64c8-4d54-869b-512aeddb976c"
  },
  {
    "state": "win",
    "amount": 72.09,
    "transactionId": "17d34772-732e-4ce4-9da5-83ce0690c45c"
  },
  {
    "state": "lose",
    "amount": 79.83,
    "transactionId": "ded6aa08-2dbb-4740-a11d-6c8bf8b083b8"
  },
  {
    "state": "win",
    "amount": 45.02,
    "transactionId": "bd83e8f1-a30f-4061-8189-ccd8198d0a02"
  },
  {
    "state": "win",
    "amount": 94.5,
    "transactionId": "b7c36fcb-4afa-4efb-9409-518693e1158c"
  },
  {
    "state": "win",
    "amount": 67.79,
    "transactionId": "66d32655-3e83-4df5-93e6-d33079bd4b76"
  },
  {
    "state": "win",
    "amount": 27.64,
    "transactionId": "7a636ed2-441d-4417-b4f8-67f000a75474"
  },
  {
    "state": "lose",
    "amount": 28.72,
    "transactionId": "49a05424-3338-4e10-9c99-306823704077"
  },
  {
    "state": "win",
    "amount": 91.27,
    "transactionId": "7698c9f6-57c6-43e4-b190-a17cd8ead930"
  },
  {
    "state": "win",
    "amount": 89.7,
    "transactionId": "f95b956e-434a-424b-aaed-f7a83882b3a9"
  },
  {
    "state": "win",
    "amount": 7.65,
    "transactionId": "418c6496-83c5-47fd-9c5d-f54f95c36c35"
  },
  {
    "state": "lose",
    "amount": 21.45,
    "transactionId": "b2ec7c0b-2333-4605-a380-8b5d80faa12b"
  },
  {
    "state": "win",
    "amount": 90.62,
    "transactionId": "c5513459-3f46-4d16-a0b4-829009f25dd5"
  },
  {
    "state": "lose",
    "amount": 81.46,
    "transactionId": "4ad8a80c-334b-443c-940e-23f8c96d20f1"
  },
  {
    "state": "lose",
    "amount": 27.75,
    "transactionId": "84e909d7-1e38-473f-adb7-c1c0731e80f4"
  },
  {
    "state": "lose",
    "amount": 88.3,
    "transactionId": "79a84e7d-cabd-45ac-89ea-61a15c484651"
  },
  {
    "state": "lose",
    "amount": 20.82,
    "transactionId": "50d22853-a711-43b4-a5f6-a80bdfe0df09"
  },
  {
    "state": "lose",
    "amount": 2.99,
    "transactionId": "f6a2c2e2-3ad2-4ab0-891a-522598e1ab07"
  },
  {
    "state": "lose",
    "amount": 36.9,
    "transactionId": "e02693c3-dbab-4538-83d4-413cf80b5479"
  },
  {
    "state": "win",
    "amount": 46.6,
    "transactionId": "15a7b652-658e-41d4-bf27-e68340489240"
  },
  {
    "state": "lose",
    "amount": 17.95,
    "transactionId": "555d8826-0f8a-45c8-8f57-3505433a42ce"
  },
  {
    "state": "lose",
    "amount": 56.66,
    "transactionId": "e60a74e8-7304-475e-b3bd-53cf46c33af0"
  },
  {
    "state": "win",
    "amount": 80.37,
    "transactionId": "b64eff1f-1210-43e0-8827-f8d78698052c"
  },
  {
    "state": "lose",
    "amount": 90.31,
    "transactionId": "1aa2aa41-63c8-4822-a60c-ba72903d10b1"
  },
  {
    "state": "win",
    "amount": 72.18,
    "transactionId": "f1401807-686e-4946-bff1-e40fe35a8a10"
  },
  {
    "state": "lose",
    "amount": 73.54,
    "transactionId": "de3c5b26-d99f-4686-8e75-7c2abcf4f7a9"
  },
  {
    "state": "win",
    "amount": 36.06,
    "transactionId": "1b96ae6a-e9cc-42ea-a33f-da87e7cd1bed"
  },
  {
    "state": "lose",
    "amount": 43.44,
    "transactionId": "473a6881-7c9f-44e8-a270-bc2090b5d692"
  },
  {
    "state": "lose",
    "amount": 86.77,
    "transactionId": "1abaae3e-c5ba-4639-992f-47b2bac80f68"
  },
  {
    "state": "lose",
    "amount": 33.33,
    "transactionId": "81262cb1-1354-42a7-b0a4-3e2d19529fa1"
  },
  {
    "state": "lose",
    "amount": 79.47,
    "transactionId": "14b66ddf-1105-433e-9ab6-95aa5b2457dd"
  },
  {
    "state": "lose",
    "amount": 33.99,
    "transactionId": "8b57662f-cc51-46cf-a47c-d8a23a5ea4e0"
  },
  {
    "state": "win",
    "amount": 67.36,
    "transactionId": "4bd3d049-b95e-46f6-b881-6f32efd08623"
  },
  {
    "state": "win",
    "amount": 75.66,
    "transactionId": "ad6a2f8f-a60b-44c5-afe8-d52120c2247f"
  },
  {
    "state": "win",
    "amount": 77.87,
    "transactionId": "d6c1c272-d036-44a3-8e20-ea4a140dfa4d"
  },
  {
    "state": "lose",
    "amount": 1.52,
    "transactionId": "1de8c68c-a388-4e6c-9080-6cfdae9b46f0"
  },
  {
    "state": "win",
    "amount": 3.31,
    "transactionId": "b082c966-a272-4b4a-a0ec-386c967a9852"
  },
  {
    "state": "win",
    "amount": 71.71,
    "transactionId": "8683cd8d-a8ad-4c5f-9eba-42205e411363"
  },
  {
    "state": "lose",
    "amount": 22.17,
    "transactionId": "4955f77b-8a99-4f97-a030-ed983fa5266a"
  },
  {
    "state": "win",
    "amount": 24.45,
    "transactionId": "fc2d3961-f99c-42a8-821f-7c8799cce910"
  },
  {
    "state": "win",
    "amount": 83.91,
    "transactionId": "7241f8fa-fecb-4d17-a2d9-b30ddde357cb"
  },
  {
    "state": "lose",
    "amount": 0.59,
    "transactionId": "acad9e35-e926-44b5-b36d-e97f30ef10de"
  },
  {
    "state": "lose",
    "amount": 26.06,
    "transactionId": "84009c92-cda8-4978-9769-6f4a0492e238"
  },
  {
    "state": "win",
    "amount": 9.53,
    "transactionId": "061dd7d4-9e76-40af-8f67-4eed63e670a3"
  },
  {
    "state": "lose",
    "amount": 51.12,
    "transactionId": "372802b2-a072-44fb-955c-62d14340241d"
  },
  {
    "state": "lose",
    "amount": 2.29,
    "transactionId": "9cdd5208-5ac6-4529-831c-f79e95d9fbe0"
  },
  {
    "state": "win",
    "amount": 43.88,
    "transactionId": "5e4f424f-5dea-4304-982f-33a14da7df0d"
  },
  {
    "state": "lose",
    "amount": 82.94,
    "transactionId": "b5964cba-4890-49ae-b86d-3bda5fbb4c74"
  },
  {
    "state": "lose",
    "amount": 28.45,
    "transactionId": "1eea68e8-a21b-4401-a5c9-69961ae4c0e3"
  },
  {
    "state": "lose",
    "amount": 97.08,
    "transactionId": "d7ab5d0d-7d17-49d9-abd6-0af8c509ca62"
  },
  {
    "state": "win",
    "amount": 38.61,
    "transactionId": "fd1754f6-2b34-4208-846f-67bcf3c6cde6"
  },
  {
    "state": "win",
    "amount": 79.89,
    "transactionId": "a4fa80ce-001e-4aca-9f8e-5ba137f13be0"
  },
  {
    "state": "win",
    "amount": 40.81,
    "transactionId": "a1fcf513-a18a-43c6-858c-d9a32d6c124a"
  },
  {
    "state": "lose",
    "amount": 73.91,
    "transactionId": "7c8cefda-e467-40a3-afaf-327f496e1056"
  },
  {
    "state": "win",
    "amount": 98.83,
    "transactionId": "9f691679-5a87-4ce9-886d-1c23835e5038"
  },
  {
    "state": "win",
    "amount": 93.3,
    "transactionId": "cf27e993-a40a-491d-be3f-0fafdf5cb9cb"
  },
  {
    "state": "win",
    "amount": 2.72,
    "transactionId": "be9a11cb-7896-4c30-b689-7ac983d06218"
  },
  {
    "state": "win",
    "amount": 47.02,
    "transactionId": "e09fcb30-29f6-40f2-8046-bd62bbb21505"
  },
  {
    "state": "lose",
    "amount": 24.55,
    "transactionId": "767d36ad-34ac-4e57-a26d-519b9be44e09"
  },
  {
    "state": "win",
    "amount": 8.74,
    "transactionId": "53521e97-0f42-4d95-8b10-dfce7a1adf03"
  },
  {
    "state": "lose",
    "amount": 5.62,
    "transactionId": "0a6b1a1d-c004-4e6c-8715-d0ea17ef7f12"
  },
  {
    "state": "win",
    "amount": 60.5,
    "transactionId": "1f4f51df-d465-4e14-838d-ae006b3cccc2"
  },
  {
    "state": "lose",
    "amount": 1.66,
    "transactionId": "3c0c8cf1-2747-4823-b6d3-d66b97e125e2"
  },
  {
    "state": "win",
    "amount": 92.47,
    "transactionId": "e931c9ab-4e44-4bbe-970e-e948136bf3a9"
  },
  {
    "state": "lose",
    "amount": 93.93,
    "transactionId": "60ab5a42-edff-45c2-bbf0-1abb663ab03b"
  },
  {
    "state": "lose",
    "amount": 43.89,
    "transactionId": "ded4ef28-9b2c-4e24-8450-d8b4645a7d82"
  },
  {
    "state": "lose",
    "amount": 28.18,
    "transactionId": "477c8cde-17f5-43f0-af5a-61dfe6201271"
  },
  {
    "state": "lose",
    "amount": 15.1,
    "transactionId": "65745b23-731f-4951-bc3d-8ec1f87fb24e"
  },
  {
    "state": "win",
    "amount": 25.47,
    "transactionId": "334f9562-ea22-427d-8c54-2f5ac2ceea1a"
  },
  {
    "state": "lose",
    "amount": 3.4,
    "transactionId": "0aa44ff1-fcfd-42d7-b540-305cd3fb4956"
  },
  {
    "state": "lose",
    "amount": 83.84,
    "transactionId": "822028eb-a7c6-46b7-b88a-4fb3efbf326d"
  },
  {
    "state": "lose",
    "amount": 59.91,
    "transactionId": "818f6b49-d856-41fe-b46c-7998b7b6de98"
  },
  {
    "state": "win",
    "amount": 48.56,
    "transactionId": "2e03b494-de32-49f3-8ea1-5dcf9196207b"
  },
  {
    "state": "lose",
    "amount": 98.49,
    "transactionId": "ac631c30-4f3c-485a-9473-5b93d9fee312"
  },
  {
    "state": "win",
    "amount": 9.48,
    "transactionId": "ed146ad5-cd1a-45f4-abe9-dc6bb6e62ba4"
  },
  {
    "state": "lose",
    "amount": 49.91,
    "transactionId": "c02284b7-4033-4199-99a3-1fe55b721190"
  },
  {
    "state": "lose",
    "amount": 2.13,
    "transactionId": "26b024cf-c8bf-49ba-871f-b023ac3db662"
  },
  {
    "state": "win",
    "amount": 80.89,
    "transactionId": "09bd364f-205d-446b-be47-24306262cfcb"
  },
  {
    "state": "win",
    "amount": 35.49,
    "transactionId": "73421aec-8c09-423a-bb7f-c681297e5575"
  },
  {
    "state": "lose",
    "amount": 0.64,
    "transactionId": "b7261eae-0374-48af-8468-ed5aa74859c5"
  },
  {
    "state": "win",
    "amount": 38.51,
    "transactionId": "080a52da-218a-49e9-b323-184b4f00f17e"
  },
  {
    "state": "win",
    "amount": 6.35,
    "transactionId": "ee157f18-d138-433e-8aa9-db80b489bb1b"
  },
  {
    "state": "lose",
    "amount": 96.05,
    "transactionId": "6284c812-e99a-48b5-9d5c-fbef74a41ef4"
  },
  {
    "state": "win",
    "amount": 38.83,
    "transactionId": "a55de45f-0168-4d14-84e5-a29ff3f91e09"
  },
  {
    "state": "lose",
    "amount": 70.63,
    "transactionId": "381af1fc-d5d0-42ad-89ba-731abc3bef98"
  },
  {
    "state": "lose",
    "amount": 24.12,
    "transactionId": "2791f116-48fd-4008-9893-4054681b019a"
  },
  {
    "state": "lose",
    "amount": 43.29,
    "transactionId": "fed381ea-f6b3-4651-8425-e24ddd0a5b82"
  },
  {
    "state": "lose",
    "amount": 8.51,
    "transactionId": "5b42d95d-3374-4fe4-974b-f8db116fb499"
  },
  {
    "state": "win",
    "amount": 43.83,
    "transactionId": "5fed4042-60e0-4ecf-b0d7-c9068a15ef99"
  },
  {
    "state": "lose",
    "amount": 35.07,
    "transactionId": "d575bdcf-f36c-4bd2-be12-9d84e1dec25e"
  },
  {
    "state": "win",
    "amount": 12.7,
    "transactionId": "6d65f3e9-1c70-4d69-bf26-8c3ad38faa36"
  },
  {
    "state": "win",
    "amount": 55.23,
    "transactionId": "a1fd211a-f266-4a28-980b-ac9a6a179b40"
  },
  {
    "state": "win",
    "amount": 23.22,
    "transactionId": "08001a95-b7f6-4996-ae4d-57f217f8ed22"
  },
  {
    "state": "win",
    "amount": 53.27,
    "transactionId": "62f045ef-b4cf-4fc3-972e-3c352c9464c6"
  },
  {
    "state": "lose",
    "amount": 8.03,
    "transactionId": "da22372e-6242-432a-a671-deafa603de7d"
  },
  {
    "state": "lose",
    "amount": 92.83,
    "transactionId": "977efc8b-3770-4045-8073-b0d7fefa691f"
  },
  {
    "state": "win",
    "amount": 42.67,
    "transactionId": "7b05041c-0caf-48dd-9794-8488996431e2"
  },
  {
    "state": "win",
    "amount": 81.35,
    "transactionId": "4d78f21b-8aea-4fff-9588-b15a2db9cc5a"
  },
  {
    "state": "lose",
    "amount": 49.15,
    "transactionId": "5f2b646c-3c8f-484a-9d0e-bc762ef74d05"
  },
  {
    "state": "lose",
    "amount": 28.98,
    "transactionId": "baa7011c-f2de-4a66-af0a-17a728ead23f"
  },
  {
    "state": "win",
    "amount": 16.98,
    "transactionId": "ca86ee94-7758-40ba-ae26-5b917c9ae4ee"
  },
  {
    "state": "lose",
    "amount": 92.61,
    "transactionId": "4d3bde1a-d208-447d-aa74-ef82fc381563"
  },
  {
    "state": "win",
    "amount": 83.27,
    "transactionId": "ddb2d3c4-5c62-4f97-a8fa-3e7976cbbaed"
  },
  {
    "state": "win",
    "amount": 60.67,
    "transactionId": "3f550f6f-12dc-42b7-81cf-8ebf8df7bd9a"
  },
  {
    "state": "win",
    "amount": 10.17,
    "transactionId": "adc7b75d-1f71-4dbe-b408-bb8a82d3354b"
  },
  {
    "state": "lose",
    "amount": 91.66,
    "transactionId": "5f05b0e3-0ca5-4e3a-87cc-ab656b099946"
  },
  {
    "state": "lose",
    "amount": 58.84,
    "transactionId": "cfd7163b-4c95-4dfb-9073-0d4ebfbc6fca"
  },
  {
    "state": "win",
    "amount": 29.54,
    "transactionId": "d3eadbc3-b963-4dc9-b8b4-f2c4a4f2ee43"
  },
  {
    "state": "win",
    "amount": 76.92,
    "transactionId": "f96a67f6-a1e4-4879-960c-b4ebb6ec7c4b"
  },
  {
    "state": "win",
    "amount": 78.1,
    "transactionId": "5092f6f9-b2b2-4839-8e53-fea51cfb21e0"
  },
  {
    "state": "lose",
    "amount": 95.82,
    "transactionId": "b8a9b9f5-f82e-4473-9803-2f5fb7bab3c4"
  },
  {
    "state": "win",
    "amount": 68.76,
    "transactionId": "48ecc401-df44-438c-a33a-cac41e758f50"
  },
  {
    "state": "lose",
    "amount": 46.45,
    "transactionId": "f21ff541-7d5e-45c2-b5e9-0ea6de148abd"
  },
  {
    "state": "lose",
    "amount": 75.75,
    "transactionId": "e5dceed8-c46d-4577-98b8-a57d3369b0e0"
  },
  {
    "state": "lose",
    "amount": 97.39,
    "transactionId": "6a1b1b5f-03d9-450b-a495-0886685beff2"
  },
  {
    "state": "win",
    "amount": 83.7,
    "transactionId": "aaaabccb-0050-4cae-a6f0-e830cf7aa7ed"
  },
  {
    "state": "win",
    "amount": 27.41,
    "transactionId": "86be3eb5-48fb-4dbf-a855-77f7ef072de0"
  },
  {
    "state": "win",
    "amount": 67.75,
    "transactionId": "32ef937e-0536-40a6-ade5-946b46422bed"
  },
  {
    "state": "lose",
    "amount": 51.8,
    "transactionId": "590c0083-4582-465d-ad2b-3b90df060f52"
  },
  {
    "state": "win",
    "amount": 9.59,
    "transactionId": "ad9e65e0-86c5-47da-b362-3bff7d500f8d"
  },
  {
    "state": "lose",
    "amount": 9.9,
    "transactionId": "8ed15f54-3b07-4e0f-b9e6-28ccd9304df0"
  },
  {
    "state": "lose",
    "amount": 68.66,
    "transactionId": "0cb10d0a-c3e7-433c-9660-be2e725b84a0"
  },
  {
    "state": "win",
    "amount": 78.17,
    "transactionId": "6fca614f-129e-451c-bc66-f9d326eba8f4"
  },
  {
    "state": "win",
    "amount": 68.84,
    "transactionId": "3c6cb894-6bb9-45f4-a7e6-be6bb97f9847"
  },
  {
    "state": "lose",
    "amount": 62.35,
    "transactionId": "e422a294-5d18-46d2-a92b-43c137b14798"
  },
  {
    "state": "win",
    "amount": 54.47,
    "transactionId": "7177200c-0a1c-4a10-a78b-b860287a85d6"
  },
  {
    "state": "win",
    "amount": 90.3,
    "transactionId": "ec71196f-19f7-4a35-9300-8ac77bd83044"
  },
  {
    "state": "win",
    "amount": 48.4,
    "transactionId": "aa0b93c4-0420-48e0-a11f-a9c763d7a3b2"
  },
  {
    "state": "lose",
    "amount": 5.07,
    "transactionId": "03119abe-4b5c-4b3a-8924-388c2cb6b805"
  },
  {
    "state": "lose",
    "amount": 2.21,
    "transactionId": "d5eb555a-5bd5-4550-9f9a-f668356699ef"
  },
  {
    "state": "win",
    "amount": 75.8,
    "transactionId": "5bff1509-be55-4565-a75b-49d6b4426e94"
  },
  {
    "state": "win",
    "amount": 44.2,
    "transactionId": "3072d2a1-15db-40ce-9dec-fc3faa3587be"
  },
  {
    "state": "win",
    "amount": 88.74,
    "transactionId": "fb328b72-0c7a-4d36-97f7-7629c41b0610"
  },
  {
    "state": "win",
    "amount": 66.71,
    "transactionId": "d3d3d65f-5381-4858-9b60-585a0aee3873"
  },
  {
    "state": "lose",
    "amount": 70.03,
    "transactionId": "61aa8cae-b7ef-4359-a214-9da0192120f0"
  },
  {
    "state": "win",
    "amount": 78.56,
    "transactionId": "b820c9af-b1b9-4362-9a5e-21022ad34396"
  },
  {
    "state": "win",
    "amount": 52.43,
    "transactionId": "c7502387-0c2b-48c1-b17b-b14173a2f2f2"
  },
  {
    "state": "lose",
    "amount": 25.18,
    "transactionId": "8fb0ba76-e657-407e-9a71-ed7f3b349dc9"
  },
  {
    "state": "win",
    "amount": 74.46,
    "transactionId": "bf7bc48a-7ee2-412e-8723-5b99cfc3c513"
  },
  {
    "state": "win",
    "amount": 13.73,
    "transactionId": "45b86265-d58e-4054-96af-c5cb0f556efb"
  },
  {
    "state": "win",
    "amount": 22.44,
    "transactionId": "750d6883-cebf-441c-8cc8-46ccc8e126d0"
  },
  {
    "state": "lose",
    "amount": 82.8,
    "transactionId": "f8bd3747-1a32-455f-a3c4-7290bf83d899"
  },
  {
    "state": "lose",
    "amount": 11.21,
    "transactionId": "70f595f2-4bb7-4083-a398-cc89a45a2072"
  },
  {
    "state": "lose",
    "amount": 56.38,
    "transactionId": "cf9b948e-c302-4df9-9497-6e65bf95033e"
  },
  {
    "state": "lose",
    "amount": 69.89,
    "transactionId": "451ebe35-6b0d-4497-b6f9-3fc6cbdf0e1c"
  },
  {
    "state": "win",
    "amount": 88.95,
    "transactionId": "c9c3ab3b-5824-4bce-984c-be059a9b4e30"
  },
  {
    "state": "lose",
    "amount": 0.11,
    "transactionId": "5f154190-7a68-48d5-a80f-e61dfb7a0901"
  },
  {
    "state": "win",
    "amount": 80.5,
    "transactionId": "aabfa109-9450-4519-baff-71b5b2a013d4"
  },
  {
    "state": "win",
    "amount": 46.22,
    "transactionId": "634b5a6b-1fec-4c93-a787-b3a34ba778c8"
  },
  {
    "state": "lose",
    "amount": 41.1,
    "transactionId": "7c092e05-63cb-4022-8889-2f037210888a"
  },
  {
    "state": "lose",
    "amount": 87.06,
    "transactionId": "f2ab6a8a-5f1c-40e3-8e61-fff2c4f955a2"
  },
  {
    "state": "lose",
    "amount": 89.57,
    "transactionId": "a707b6e4-ac4d-475b-93d1-60c46c1b6711"
  },
  {
    "state": "lose",
    "amount": 35.04,
    "transactionId": "0f6a441e-118f-4db8-b95f-2ce75feaeb03"
  },
  {
    "state": "win",
    "amount": 64.26,
    "transactionId": "1d94a1d3-1aae-4457-b916-1a0db4414500"
  },
  {
    "state": "lose",
    "amount": 22.08,
    "transactionId": "fcb60365-79f5-4dad-b1a6-37c209686bee"
  },
  {
    "state": "lose",
    "amount": 32.06,
    "transactionId": "b3082a58-7a2a-4ada-b7fe-39e21e9c8d09"
  },
  {
    "state": "lose",
    "amount": 5.49,
    "transactionId": "ea0c662f-53c3-4e6a-80fb-9ba045f266d8"
  },
  {
    "state": "lose",
    "amount": 77.93,
    "transactionId": "751f42ed-b6e2-4461-9b1e-072335cf31c1"
  },
  {
    "state": "win",
    "amount": 70.65,
    "transactionId": "01104d0a-a9bf-4486-b4dc-44c0d92cefdb"
  },
  {
    "state": "lose",
    "amount": 64.71,
    "transactionId": "cc830c0f-3aca-4f08-a92d-ddd9fe55732f"
  },
  {
    "state": "lose",
    "amount": 56.09,
    "transactionId": "630cca9b-2021-48dd-919a-e8b2eb981335"
  },
  {
    "state": "lose",
    "amount": 40.02,
    "transactionId": "81bfaeba-e4a2-46ac-afbe-66cee32ab033"
  },
  {
    "state": "win",
    "amount": 52.62,
    "transactionId": "d0fdeeb3-f49c-4170-a698-ef2c74b4bf3b"
  },
  {
    "state": "win",
    "amount": 92.37,
    "transactionId": "4c919216-2ce0-4342-bac1-daff4536c69e"
  },
  {
    "state": "win",
    "amount": 28.83,
    "transactionId": "b52b7fb5-2a38-4a2f-903a-13d642679fdf"
  },
  {
    "state": "lose",
    "amount": 72.15,
    "transactionId": "89e207ea-79fd-4fb3-9b6f-baf1102823f7"
  },
  {
    "state": "win",
    "amount": 19.59,
    "transactionId": "c9c7356c-5be4-46b8-8af0-aff739e3e79f"
  },
  {
    "state": "lose",
    "amount": 67.95,
    "transactionId": "ba9d355c-cfa1-45dc-b595-04385c5ec576"
  },
  {
    "state": "lose",
    "amount": 89.09,
    "transactionId": "ebadec94-0e7d-463d-827e-6eb53d5271fb"
  },
  {
    "state": "win",
    "amount": 27.4,
    "transactionId": "eaabe7e3-b696-4ef2-a8f7-3d764a9b001e"
  },
  {
    "state": "win",
    "amount": 30.93,
    "transactionId": "a1cbb346-5835-41b4-a39e-330bb9bb380a"
  },
  {
    "state": "win",
    "amount": 96.22,
    "transactionId": "80a364d6-eac1-4693-bd17-b718f1da7aaf"
  },
  {
    "state": "win",
    "amount": 30.79,
    "transactionId": "fe0d7c4b-d69b-4304-aea7-2ab8ce1319dc"
  },
  {
    "state": "lose",
    "amount": 32.09,
    "transactionId": "d73572d6-9d03-4d21-a07d-a116329f784b"
  },
  {
    "state": "win",
    "amount": 59.26,
    "transactionId": "4786286d-b297-4d7d-a17c-b3b0434baae0"
  },
  {
    "state": "win",
    "amount": 65.11,
    "transactionId": "d4e57fec-8b32-449a-9a83-4854cf87e03b"
  },
  {
    "state": "win",
    "amount": 61.21,
    "transactionId": "e97e5b6e-1850-447e-afc6-d664634dd0c1"
  },
  {
    "state": "win",
    "amount": 12.22,
    "transactionId": "e0abc1f2-bdde-44d3-b587-ba194a88578c"
  },
  {
    "state": "lose",
    "amount": 39.26,
    "transactionId": "d7214c12-63b3-4b46-b37e-29cca29f4d04"
  },
  {
    "state": "win",
    "amount": 19.11,
    "transactionId": "c5173b54-bbbd-4c02-8d6c-f1169df1911f"
  },
  {
    "state": "win",
    "amount": 46.56,
    "transactionId": "1f96d9b7-99cd-403a-a7f4-1a46f2bd9096"
  },
  {
    "state": "lose",
    "amount": 87.22,
    "transactionId": "5c1b30bf-601f-49eb-bd6d-ac99d746d0e4"
  },
  {
    "state": "win",
    "amount": 90.27,
    "transactionId": "b07463a0-f448-4698-a806-9f1d0de79dfa"
  },
  {
    "state": "win",
    "amount": 50.57,
    "transactionId": "7ccfa4ee-4cd6-4f73-a21c-b5b66706fe1b"
  },
  {
    "state": "lose",
    "amount": 46.53,
    "transactionId": "f8b6b012-5516-401c-9113-dfaa2cf99296"
  },
  {
    "state": "lose",
    "amount": 73.43,
    "transactionId": "4331fb2f-c85d-4945-8a2b-c1dc52a64468"
  },
  {
    "state": "lose",
    "amount": 65.81,
    "transactionId": "5ce8afb8-1b36-45b2-80a1-e025db09ca67"
  },
  {
    "state": "lose",
    "amount": 69.07,
    "transactionId": "0449bbc2-25de-48b1-af75-44828efc456b"
  },
  {
    "state": "win",
    "amount": 90.61,
    "transactionId": "652ffdea-90d8-4eea-a788-60d46dc49773"
  },
  {
    "state": "win",
    "amount": 54.69,
    "transactionId": "b788fb24-4283-4849-9497-234e483f0c68"
  },
  {
    "state": "lose",
    "amount": 79.66,
    "transactionId": "95ddd317-0444-4dd4-bbad-c96c7c51de3a"
  },
  {
    "state": "win",
    "amount": 86.2,
    "transactionId": "e8514ed3-b62d-4043-84d8-c43d79fcd49e"
  },
  {
    "state": "lose",
    "amount": 3.06,
    "transactionId": "ff029026-4590-4719-8ca3-987dfe4586d5"
  },
  {
    "state": "lose",
    "amount": 55.39,
    "transactionId": "184e7b45-6a4f-4fab-a6a4-175379025e5d"
  },
  {
    "state": "lose",
    "amount": 56.58,
    "transactionId": "6984670d-77d5-4913-95f0-418244c2bffa"
  },
  {
    "state": "lose",
    "amount": 31.33,
    "transactionId": "5fa7c774-4eef-4ebb-ae1b-e3ff053f678a"
  },
  {
    "state": "lose",
    "amount": 73.8,
    "transactionId": "e5425afb-b83a-47e0-a61f-25f43112b627"
  },
  {
    "state": "lose",
    "amount": 76.62,
    "transactionId": "a7982f15-2e3f-4b1e-99c7-6beaaba83b6f"
  },
  {
    "state": "lose",
    "amount": 19.37,
    "transactionId": "2b7af244-13a6-480d-9a00-787d15d89c8a"
  },
  {
    "state": "lose",
    "amount": 97.24,
    "transactionId": "44a70107-8c99-4145-a2c8-37267b14626d"
  },
  {
    "state": "lose",
    "amount": 55.75,
    "transactionId": "28e35620-6e4e-43e6-be94-30ec9c997feb"
  },
  {
    "state": "win",
    "amount": 40.14,
    "transactionId": "63bc7801-622e-4148-981a-129f11e98f2e"
  },
  {
    "state": "lose",
    "amount": 24.4,
    "transactionId": "9dbd37b4-3086-42f4-bbf9-90bc3aab2391"
  },
  {
    "state": "win",
    "amount": 27.69,
    "transactionId": "5e9cb4ff-2895-4db4-9863-f45ca1502f08"
  },
  {
    "state": "win",
    "amount": 66.36,
    "transactionId": "b6d39e5a-05a4-4b9e-8ce4-1258a034b1f7"
  },
  {
    "state": "lose",
    "amount": 34.06,
    "transactionId": "77bdc561-afac-4267-921d-e89e5158c7d9"
  },
  {
    "state": "win",
    "amount": 68.76,
    "transactionId": "f335d89a-0be7-4c75-b004-3f0dcdfef49d"
  },
  {
    "state": "win",
    "amount": 52.62,
    "transactionId": "1f899b24-c3af-46e8-8ee9-f9fbd70e5d2a"
  },
  {
    "state": "win",
    "amount": 3.8,
    "transactionId": "f3926aee-7910-406e-b6f6-f95c1273f02b"
  },
  {
    "state": "win",
    "amount": 9.78,
    "transactionId": "663581de-5a00-4e40-a5f2-18bc3d77b2a5"
  },
  {
    "state": "lose",
    "amount": 39.84,
    "transactionId": "9b2dcfca-3b89-4256-8e86-daed8eab7efb"
  },
  {
    "state": "lose",
    "amount": 55.05,
    "transactionId": "e549a1c3-7ec3-4fab-a0a8-9df6d068da1d"
  },
  {
    "state": "win",
    "amount": 23.64,
    "transactionId": "e19a4ae9-3016-41a5-9ac8-58fd740763bb"
  },
  {
    "state": "win",
    "amount": 70.11,
    "transactionId": "4229ef38-1230-4dc8-96e2-065474896766"
  },
  {
    "state": "lose",
    "amount": 29.48,
    "transactionId": "0a98bf80-0cbe-422a-af3d-93cdf2b84926"
  },
  {
    "state": "win",
    "amount": 72.58,
    "transactionId": "7971178d-9b7b-4942-9567-ab2895259ed2"
  },
  {
    "state": "win",
    "amount": 89.9,
    "transactionId": "7adf1bcb-d7c2-4f61-b918-4ed0d59cc0ec"
  },
  {
    "state": "lose",
    "amount": 98.33,
    "transactionId": "55bf98e2-9707-45cc-bf56-87a9b9266403"
  },
  {
    "state": "win",
    "amount": 53.61,
    "transactionId": "66f79322-ff39-4118-8f7d-37b415616bd6"
  },
  {
    "state": "lose",
    "amount": 46.05,
    "transactionId": "e9fd4784-916d-4be9-b606-c04c65596b9e"
  },
  {
    "state": "lose",
    "amount": 32.27,
    "transactionId": "85e0ec81-49c2-43c6-8214-fa37d4856ae1"
  },
  {
    "state": "lose",
    "amount": 36.68,
    "transactionId": "ff701a20-b679-492c-870a-d8513764c06d"
  },
  {
    "state": "lose",
    "amount": 1.35,
    "transactionId": "2a75c60b-5bc5-4068-8cde-b68ac5b815ee"
  },
  {
    "state": "win",
    "amount": 99.27,
    "transactionId": "0a53acf5-65cd-4189-8de3-70e0e634ae37"
  },
  {
    "state": "win",
    "amount": 88.48,
    "transactionId": "e340fda3-e8ca-4863-ae50-15ca7deb0d79"
  },
  {
    "state": "lose",
    "amount": 75.13,
    "transactionId": "e4214a37-d646-4c18-81f8-7c3480b171f5"
  },
  {
    "state": "win",
    "amount": 64.36,
    "transactionId": "efb8c327-095a-4b15-ae98-1e3a9b3b5e3c"
  },
  {
    "state": "lose",
    "amount": 27.49,
    "transactionId": "f9bd8f79-3815-4371-b905-d63715a3b846"
  },
  {
    "state": "win",
    "amount": 99.37,
    "transactionId": "3e69c964-b84e-4ecd-baf9-e63b17c071c2"
  },
  {
    "state": "win",
    "amount": 12.74,
    "transactionId": "98aad84e-f6c0-4294-a109-3bbc53f201eb"
  },
  {
    "state": "lose",
    "amount": 5.58,
    "transactionId": "96cd4762-af23-48bc-8f2b-12b936c89dd6"
  },
  {
    "state": "win",
    "amount": 63.15,
    "transactionId": "a64b99e4-11ec-45de-849e-a752f0efce06"
  },
  {
    "state": "win",
    "amount": 3.15,
    "transactionId": "52f56858-ccd0-4b72-a61a-7bbf62b5b0d3"
  },
  {
    "state": "win",
    "amount": 97.98,
    "transactionId": "2dcc3913-75b4-47b0-9710-a763aa1c5296"
  },
  {
    "state": "win",
    "amount": 66.84,
    "transactionId": "7a3e7777-5b59-46b0-ab33-a517b9a3b737"
  },
  {
    "state": "win",
    "amount": 21.51,
    "transactionId": "484593bd-5dc3-4ae7-ac5f-5b3e2a5f6d9a"
  },
  {
    "state": "win",
    "amount": 94.34,
    "transactionId": "0ba43536-22fa-4e5a-bb37-d7509b96e0e7"
  },
  {
    "state": "lose",
    "amount": 49.36,
    "transactionId": "f1ad3026-d7ce-4dcc-b8c9-2d7f5bd5cab0"
  },
  {
    "state": "win",
    "amount": 18.42,
    "transactionId": "ecaea5b3-9ae6-4b9a-9be2-7c7eec06ee34"
  },
  {
    "state": "win",
    "amount": 72.74,
    "transactionId": "aaf344f4-9c6a-4e59-8568-e49022718343"
  },
  {
    "state": "lose",
    "amount": 72.78,
    "transactionId": "db0d65b6-7f2c-4f1f-a7a3-dc4151e6819b"
  },
  {
    "state": "lose",
    "amount": 12.97,
    "transactionId": "1ee9feca-cc15-4843-9280-de4ac442bfb1"
  },
  {
    "state": "win",
    "amount": 40.63,
    "transactionId": "535440d9-38dd-43cb-a6f1-c64b3c6196d4"
  },
  {
    "state": "lose",
    "amount": 20.58,
    "transactionId": "9513de58-327d-4490-bb95-44e4e78b1f7a"
  },
  {
    "state": "lose",
    "amount": 33.39,
    "transactionId": "c35c44da-4c07-488c-8052-fe45ca737a5b"
  },
  {
    "state": "win",
    "amount": 71.05,
    "transactionId": "2eaf0c18-f2cd-425c-aa07-e6d2c397e7fa"
  },
  {
    "state": "win",
    "amount": 49.68,
    "transactionId": "bc361a6e-f183-4310-b19a-c40975e5b63d"
  },
  {
    "state": "lose",
    "amount": 84.59,
    "transactionId": "7705568b-4bc0-4f74-b9af-bec2a79c1149"
  },
  {
    "state": "lose",
    "amount": 53.53,
    "transactionId": "4a31841b-5ccb-48e9-bef9-79331d2cffee"
  },
  {
    "state": "win",
    "amount": 96.37,
    "transactionId": "970806df-7e13-47b2-8bf3-7cf66509d9df"
  },
  {
    "state": "win",
    "amount": 33.24,
    "transactionId": "ab61b35e-4e83-47fe-923d-dfac75de4afa"
  },
  {
    "state": "win",
    "amount": 76.94,
    "transactionId": "d3b2a2bc-dafb-4163-9e2e-dc6ee44e72e3"
  },
  {
    "state": "lose",
    "amount": 29.26,
    "transactionId": "35deac64-d45e-42e3-8606-19b39e059814"
  },
  {
    "state": "win",
    "amount": 13.23,
    "transactionId": "f9ed80d5-6d81-4a71-92f1-e887e65cd06c"
  },
  {
    "state": "lose",
    "amount": 1.18,
    "transactionId": "033ead23-0573-4fe3-b24b-becc4087cd7e"
  },
  {
    "state": "win",
    "amount": 94.63,
    "transactionId": "f5816114-93df-455e-a7f7-5f26edf3b4c1"
  },
  {
    "state": "lose",
    "amount": 24.72,
    "transactionId": "115a20c6-864a-43bb-a26e-c06c325b8ac8"
  },
  {
    "state": "win",
    "amount": 98.47,
    "transactionId": "0ffce26d-1862-492f-82c1-ea15d687177d"
  },
  {
    "state": "lose",
    "amount": 100,
    "transactionId": "e7bc944e-27ad-48f2-9f80-555f3dc7f3b8"
  },
  {
    "state": "win",
    "amount": 58.21,
    "transactionId": "4a75f18e-243b-4691-ab9d-cfdd55b990ec"
  },
  {
    "state": "lose",
    "amount": 88.46,
    "transactionId": "e826ac8f-2ed5-4745-8e96-a312b5366174"
  },
  {
    "state": "lose",
    "amount": 28.73,
    "transactionId": "3bd56a27-e365-45f3-97e3-282458055576"
  },
  {
    "state": "win",
    "amount": 19.11,
    "transactionId": "ff925cda-f338-45ba-a77f-8e40cc69330e"
  },
  {
    "state": "lose",
    "amount": 28.21,
    "transactionId": "23306e83-ea3d-4978-98d4-48edf5a5721d"
  },
  {
    "state": "lose",
    "amount": 31.22,
    "transactionId": "e32da796-a091-4053-abf5-61d182649da6"
  },
  {
    "state": "win",
    "amount": 16.43,
    "transactionId": "123569a1-4f40-439d-950e-5f811edcda06"
  },
  {
    "state": "win",
    "amount": 53.02,
    "transactionId": "1d9c2399-6d03-4440-b1d1-aa923af3f6a9"
  },
  {
    "state": "win",
    "amount": 58.69,
    "transactionId": "11f31390-b688-4671-9df7-c0518dde1014"
  },
  {
    "state": "lose",
    "amount": 61.86,
    "transactionId": "8c75a2ac-7e73-4f80-b15c-88d1d84dbd11"
  },
  {
    "state": "lose",
    "amount": 27.91,
    "transactionId": "af5e59cc-b6a1-49a4-9b0e-309cb657c3d3"
  },
  {
    "state": "lose",
    "amount": 26.57,
    "transactionId": "d1507b49-15c8-40e7-8923-8528ef8825da"
  },
  {
    "state": "lose",
    "amount": 72.37,
    "transactionId": "0c0532ae-1101-4daf-bfa2-e19fb19e487e"
  },
  {
    "state": "win",
    "amount": 20.35,
    "transactionId": "24a7d16f-7dbf-426d-a710-53559b5558cc"
  },
  {
    "state": "lose",
    "amount": 47.75,
    "transactionId": "f7cfc49b-d034-4ae8-bfea-17eee95aafdb"
  },
  {
    "state": "win",
    "amount": 81.42,
    "transactionId": "03e08a82-59a8-41a6-86c6-07c7c9e09eba"
  },
  {
    "state": "lose",
    "amount": 28.84,
    "transactionId": "a61cf4ce-5629-4b9c-bc7b-fc9d9fc2817c"
  },
  {
    "state": "win",
    "amount": 42.57,
    "transactionId": "008e45a3-f442-4e35-9ed6-d458c0188cb4"
  },
  {
    "state": "lose",
    "amount": 74.05,
    "transactionId": "10eb4416-97d9-4447-b2fa-1d0bf542d28d"
  },
  {
    "state": "win",
    "amount": 85.32,
    "transactionId": "cbc930f3-90d9-463e-8746-86afd00d6312"
  },
  {
    "state": "win",
    "amount": 20.62,
    "transactionId": "bb57907e-9d38-4727-99bd-381ffa54aa04"
  },
  {
    "state": "win",
    "amount": 33.98,
    "transactionId": "65982cbc-328e-4cef-be97-b23480f466aa"
  },
  {
    "state": "win",
    "amount": 67.45,
    "transactionId": "b0af4920-7c25-4d64-9c22-905eaeccbeb8"
  },
  {
    "state": "lose",
    "amount": 57.35,
    "transactionId": "cfef19cb-b255-4a34-a6cb-b6b448773c65"
  },
  {
    "state": "win",
    "amount": 6.05,
    "transactionId": "5ee585f4-b75d-44e9-bf4f-9467f221201f"
  },
  {
    "state": "win",
    "amount": 35.62,
    "transactionId": "2cfca440-4311-4044-a52c-6e1898f2b1ba"
  },
  {
    "state": "lose",
    "amount": 83.69,
    "transactionId": "ba49e664-4e96-4ed0-b419-f6a0f5975fb4"
  },
  {
    "state": "win",
    "amount": 88,
    "transactionId": "24fda39f-a779-4c31-91c3-1ff1245e0b2b"
  },
  {
    "state": "win",
    "amount": 48.94,
    "transactionId": "fa24e133-3ee7-42fd-84bb-bf49509f1c61"
  },
  {
    "state": "lose",
    "amount": 14.5,
    "transactionId": "c87d3db5-f7b5-4df2-bbd8-3a8f9747189f"
  },
  {
    "state": "win",
    "amount": 33.65,
    "transactionId": "68aef509-bd78-4c62-9838-c502b08e21b8"
  },
  {
    "state": "lose",
    "amount": 42.01,
    "transactionId": "79e5052f-184d-48ff-9885-44daf3c13e74"
  },
  {
    "state": "win",
    "amount": 42.27,
    "transactionId": "d244cd8f-6850-4382-8252-43bc19450b09"
  },
  {
    "state": "lose",
    "amount": 87.19,
    "transactionId": "45787470-c35e-42cf-bf46-2ad5388047b6"
  },
  {
    "state": "win",
    "amount": 70.65,
    "transactionId": "4f258a5a-ea39-46dc-9850-ab987da749f3"
  },
  {
    "state": "win",
    "amount": 1.13,
    "transactionId": "0d9a601a-21af-43c6-8515-c66f14b13c2b"
  },
  {
    "state": "win",
    "amount": 39.85,
    "transactionId": "9e820628-fc46-47f3-8629-b30973ef863e"
  },
  {
    "state": "lose",
    "amount": 40.86,
    "transactionId": "14b2e0c5-372a-4831-a14f-1b83a08adae8"
  },
  {
    "state": "win",
    "amount": 25.08,
    "transactionId": "7b76fe66-ad03-4d0c-9ec1-7f43399b1b8d"
  },
  {
    "state": "win",
    "amount": 6.46,
    "transactionId": "e4f1c197-fe61-4251-beb6-41309c7104dc"
  },
  {
    "state": "win",
    "amount": 56.86,
    "transactionId": "9a7b2cda-0a72-443f-be3b-9ff3ab5f9a1e"
  },
  {
    "state": "win",
    "amount": 92.96,
    "transactionId": "8c609f14-0f5d-4c02-b055-3c153c2b0554"
  },
  {
    "state": "win",
    "amount": 93.74,
    "transactionId": "40f782a4-fde7-4ce0-8c21-e75f7a443b92"
  },
  {
    "state": "win",
    "amount": 67.09,
    "transactionId": "dbdafbdc-c063-485f-bb22-01e39426568c"
  },
  {
    "state": "win",
    "amount": 87.73,
    "transactionId": "61296815-eadd-4d36-afc4-76213690fefc"
  },
  {
    "state": "lose",
    "amount": 49.2,
    "transactionId": "c760da1f-3d30-41ba-a96f-a4408f137912"
  },
  {
    "state": "lose",
    "amount": 8.78,
    "transactionId": "cf3d62a6-ffb7-4aa9-b61c-868be885f1fc"
  },
  {
    "state": "win",
    "amount": 32.07,
    "transactionId": "d85533f0-b992-4b8b-b7cd-6d6d53817dd3"
  },
  {
    "state": "lose",
    "amount": 10.18,
    "transactionId": "eb739d8e-80a9-42d4-a32a-eb661e84ba61"
  },
  {
    "state": "win",
    "amount": 16.14,
    "transactionId": "607a8c36-79d8-467f-b457-877f7974af15"
  },
  {
    "state": "lose",
    "amount": 75.03,
    "transactionId": "3223c08b-afae-4f81-aa84-6ae8cf81235c"
  },
  {
    "state": "lose",
    "amount": 28.28,
    "transactionId": "5c4248a7-3b93-4149-9274-49bf43252e8d"
  },
  {
    "state": "lose",
    "amount": 11.42,
    "transactionId": "b1b0202c-cc94-4c9b-9724-eee56b6931a0"
  },
  {
    "state": "lose",
    "amount": 66.08,
    "transactionId": "358aac20-af1b-401f-a546-d7348c143538"
  },
  {
    "state": "lose",
    "amount": 26.23,
    "transactionId": "3027e3af-8b00-483e-acd7-2e5be069800f"
  },
  {
    "state": "win",
    "amount": 6.07,
    "transactionId": "8bb82553-9507-4537-b403-5f65cd2ac0fd"
  },
  {
    "state": "lose",
    "amount": 3.61,
    "transactionId": "443b3905-7e16-45b0-bc04-6d3be298f3c7"
  },
  {
    "state": "lose",
    "amount": 2.06,
    "transactionId": "9ccacf46-55f3-43c1-b944-f91adbea3caf"
  },
  {
    "state": "win",
    "amount": 40.13,
    "transactionId": "2b080f9e-8c1e-4fd3-b7db-62930d457016"
  },
  {
    "state": "win",
    "amount": 10.35,
    "transactionId": "dbaf4cfe-54c3-4c43-8943-ddbbfd9b41be"
  },
  {
    "state": "lose",
    "amount": 83.65,
    "transactionId": "cff4d1e4-c9d4-4850-84c1-9daa69a94052"
  },
  {
    "state": "lose",
    "amount": 53.49,
    "transactionId": "04741d77-f6cf-4431-a6fd-086e25d3eb32"
  },
  {
    "state": "win",
    "amount": 65.68,
    "transactionId": "c7e5b198-72e9-4219-bca6-a5d345a3e44c"
  },
  {
    "state": "lose",
    "amount": 31.86,
    "transactionId": "7c5787c2-06f3-4e8f-9833-47d538999260"
  },
  {
    "state": "lose",
    "amount": 57.47,
    "transactionId": "a7610fb3-16ce-4f6f-9818-91ebd6f6a8b5"
  },
  {
    "state": "win",
    "amount": 76.01,
    "transactionId": "c1715d85-b730-415a-a6f0-d8bb6d47c716"
  },
  {
    "state": "win",
    "amount": 23.92,
    "transactionId": "2fbd4d6d-7408-40ff-a25d-87dcc0f83ddc"
  },
  {
    "state": "lose",
    "amount": 97.36,
    "transactionId": "57d6c0dd-c223-4631-bebe-a382a2aff1cc"
  },
  {
    "state": "lose",
    "amount": 68.34,
    "transactionId": "a7a0040c-cf4a-425b-b250-583ea09a2d7f"
  },
  {
    "state": "win",
    "amount": 30.15,
    "transactionId": "d6e7cae8-17db-464d-882a-c64c4fe67931"
  },
  {
    "state": "win",
    "amount": 76.97,
    "transactionId": "9843bc2c-3fbf-4faf-acf2-7085b430970e"
  },
  {
    "state": "win",
    "amount": 10.71,
    "transactionId": "bb38f518-6c7e-4f5b-ac46-990fe475ffb7"
  },
  {
    "state": "win",
    "amount": 82.24,
    "transactionId": "f3ad0f9a-4c2a-401a-b49a-4450af0bcb86"
  },
  {
    "state": "win",
    "amount": 46.24,
    "transactionId": "9b70a5e6-fe25-4bc9-9dd3-d3b808232a77"
  },
  {
    "state": "win",
    "amount": 30.67,
    "transactionId": "35b9f21a-08cd-4cf3-9fb9-123ac6e068d8"
  },
  {
    "state": "lose",
    "amount": 93.81,
    "transactionId": "ffda086b-b760-4d4b-bdf5-7c89122079e9"
  },
  {
    "state": "lose",
    "amount": 8.45,
    "transactionId": "e77b23ac-27ea-4d69-805b-728390362d25"
  },
  {
    "state": "win",
    "amount": 72.89,
    "transactionId": "59c4641a-553a-4064-97ab-2bd226985e00"
  },
  {
    "state": "win",
    "amount": 2.72,
    "transactionId": "cbaaad7e-495b-419b-a2ce-3f02ba401d99"
  },
  {
    "state": "lose",
    "amount": 7.37,
    "transactionId": "b1af5674-fd50-492a-ba34-e31e94779c49"
  },
  {
    "state": "win",
    "amount": 3.06,
    "transactionId": "52aedc19-5917-40d4-ad8b-a60d265dbed5"
  },
  {
    "state": "win",
    "amount": 27.85,
    "transactionId": "ab9222b8-1238-4570-a206-44d88511c6bf"
  },
  {
    "state": "lose",
    "amount": 71.45,
    "transactionId": "1367d73e-bc22-449e-885d-d60c5238c351"
  },
  {
    "state": "win",
    "amount": 33.86,
    "transactionId": "29c0a5c4-e95c-4336-94b6-89f7d530820e"
  },
  {
    "state": "win",
    "amount": 11.45,
    "transactionId": "1d476edb-0565-41a1-9dc4-e74c207a2ecd"
  },
  {
    "state": "win",
    "amount": 39.66,
    "transactionId": "82508d54-4d6e-4217-b347-7586b5dc0628"
  },
  {
    "state": "lose",
    "amount": 86.13,
    "transactionId": "a53a1852-00c6-416e-a838-59fb7e90c50a"
  },
  {
    "state": "win",
    "amount": 38.33,
    "transactionId": "8f081aed-9f7b-41b0-a8c1-ebf531ec5889"
  },
  {
    "state": "win",
    "amount": 90.43,
    "transactionId": "663102d9-92bf-4389-9a44-27d81bebf40c"
  },
  {
    "state": "win",
    "amount": 69.68,
    "transactionId": "aade50a1-89e2-4b5c-be88-0b24dc7f3d6d"
  },
  {
    "state": "lose",
    "amount": 71.45,
    "transactionId": "2eb99a5b-3efc-4054-9d8b-4ffa613855f8"
  },
  {
    "state": "lose",
    "amount": 70.9,
    "transactionId": "d1dc70e1-5f8b-4a5e-b05a-c7fc5afc7a84"
  },
  {
    "state": "lose",
    "amount": 50.67,
    "transactionId": "d5cbc2b6-2c6c-483f-b507-0c0dc51858a1"
  },
  {
    "state": "win",
    "amount": 2.46,
    "transactionId": "479b3ae6-543c-4e6c-af5e-d4bd7c931438"
  },
  {
    "state": "win",
    "amount": 7.22,
    "transactionId": "e360f442-9267-45a9-9732-51f755e97f88"
  },
  {
    "state": "lose",
    "amount": 2.54,
    "transactionId": "20cef485-8b0f-4555-9642-4307860ed94c"
  },
  {
    "state": "win",
    "amount": 46.07,
    "transactionId": "6ed1aba4-70ee-46e4-87c1-37db6a7430bd"
  },
  {
    "state": "lose",
    "amount": 24.89,
    "transactionId": "42558063-c0cf-4952-956a-2ee582ff328e"
  },
  {
    "state": "lose",
    "amount": 69.06,
    "transactionId": "fa1a6722-817a-4e49-a6b8-fd4977bd6bcf"
  },
  {
    "state": "win",
    "amount": 66.51,
    "transactionId": "f934d592-1745-4470-8013-eecc1265615d"
  },
  {
    "state": "lose",
    "amount": 10.67,
    "transactionId": "a69c4696-d734-4ba9-988b-bdc2a91e9a63"
  },
  {
    "state": "lose",
    "amount": 86.67,
    "transactionId": "55cb38a7-1600-4d5e-8ff6-9a2dd2103005"
  },
  {
    "state": "lose",
    "amount": 6.54,
    "transactionId": "63cfd4e8-2d7b-49bf-85b7-4a7e9badcf9f"
  },
  {
    "state": "win",
    "amount": 78.47,
    "transactionId": "a07a5bd2-4676-4137-87a5-0e634c7d4528"
  },
  {
    "state": "win",
    "amount": 69.8,
    "transactionId": "774fb0d5-1cb0-4531-a979-4bfbefb51ae8"
  },
  {
    "state": "win",
    "amount": 7.3,
    "transactionId": "f5bf5c2d-0fe4-4375-ae43-c99fe3eef5ff"
  },
  {
    "state": "win",
    "amount": 75.36,
    "transactionId": "12c29f8a-d664-4594-b86a-c6b3d555109f"
  },
  {
    "state": "win",
    "amount": 91.56,
    "transactionId": "a64233a2-8125-4eba-ba6e-70c8d1403ed9"
  },
  {
    "state": "win",
    "amount": 13.47,
    "transactionId": "f8358871-3571-49ff-911b-5ba6f305c76f"
  },
  {
    "state": "win",
    "amount": 71.99,
    "transactionId": "ceef2792-44fe-48e1-873d-6f48aa27fc8c"
  },
  {
    "state": "lose",
    "amount": 43.09,
    "transactionId": "5c2d915c-44d7-462c-b505-ac77273a5be4"
  },
  {
    "state": "lose",
    "amount": 52.95,
    "transactionId": "39dd139e-668a-4c4d-b6e2-12291cee0343"
  },
  {
    "state": "win",
    "amount": 87.25,
    "transactionId": "5498de5d-4ba5-4d3e-847c-e9a60c9be14a"
  },
  {
    "state": "win",
    "amount": 37.17,
    "transactionId": "35c26d00-90fa-4ae1-9555-a46bc218c14f"
  },
  {
    "state": "lose",
    "amount": 75.39,
    "transactionId": "d521a8df-3ea1-4984-bc01-6ba8f33d462d"
  },
  {
    "state": "win",
    "amount": 88.93,
    "transactionId": "4d21b19a-18f3-4240-86c0-bf46acf93293"
  },
  {
    "state": "win",
    "amount": 4.06,
    "transactionId": "180195e5-a914-429e-a32f-cdf2039a0c23"
  },
  {
    "state": "win",
    "amount": 93.19,
    "transactionId": "af171093-6b46-4799-90f6-1864b3f4d023"
  },
  {
    "state": "lose",
    "amount": 4.34,
    "transactionId": "daf7d93c-ac7f-4fe5-8465-8832117fa54b"
  },
  {
    "state": "lose",
    "amount": 66.96,
    "transactionId": "b1cb64f1-f419-4e25-b407-b3bec357ba81"
  },
  {
    "state": "win",
    "amount": 0.92,
    "transactionId": "e44047c6-e7f5-45de-bb58-25411caf98c4"
  },
  {
    "state": "lose",
    "amount": 66.72,
    "transactionId": "ec727f2c-165a-4cd5-8a3e-930eb6cf83af"
  },
  {
    "state": "lose",
    "amount": 86.68,
    "transactionId": "12218f53-d8a3-42a0-af69-1dbfcff3dcc7"
  },
  {
    "state": "lose",
    "amount": 66.84,
    "transactionId": "f753fc90-7c18-4164-a114-b0c20e28fc4e"
  },
  {
    "state": "win",
    "amount": 74.19,
    "transactionId": "6a9f9507-426b-46f6-9226-c96a853d6f13"
  },
  {
    "state": "win",
    "amount": 90.25,
    "transactionId": "ea0cc16c-7685-4c51-8f37-cf24a0ca9bc7"
  },
  {
    "state": "win",
    "amount": 79.14,
    "transactionId": "ab417608-9a60-4585-9afa-914a014698fa"
  },
  {
    "state": "win",
    "amount": 46.94,
    "transactionId": "676b3cc4-f9a8-47e4-915b-cbbabd987cf0"
  },
  {
    "state": "lose",
    "amount": 88.13,
    "transactionId": "4ffd7a1f-03a4-47c2-9d38-30fd25ef07e4"
  },
  {
    "state": "lose",
    "amount": 77.56,
    "transactionId": "321a949c-c996-4277-a194-b72edd1a95ce"
  },
  {
    "state": "win",
    "amount": 79.44,
    "transactionId": "ef3ea8c0-7706-4694-b82c-53c2e57757c6"
  },
  {
    "state": "lose",
    "amount": 22.3,
    "transactionId": "f940773c-8e3d-4143-8041-46c9872ff8e3"
  },
  {
    "state": "win",
    "amount": 41.43,
    "transactionId": "2744dda6-50ed-4660-8893-3e98c45b6de3"
  },
  {
    "state": "win",
    "amount": 42.21,
    "transactionId": "19e3b607-c839-4448-aa68-f9f674a814ae"
  },
  {
    "state": "lose",
    "amount": 98.68,
    "transactionId": "bef13feb-5d51-40c2-9d09-7fc2a07e193b"
  },
  {
    "state": "lose",
    "amount": 40.87,
    "transactionId": "53219b38-a794-464e-9699-66ab5ad6f1a3"
  },
  {
    "state": "lose",
    "amount": 28.89,
    "transactionId": "5d1c0a8d-6809-48e0-bfe9-7dfef8cf928b"
  },
  {
    "state": "lose",
    "amount": 11.49,
    "transactionId": "af9814f0-c5e8-4834-9546-fe141f4539fb"
  },
  {
    "state": "lose",
    "amount": 32.95,
    "transactionId": "e9f625ba-3d64-4757-b6d4-d38fac1b7956"
  },
  {
    "state": "lose",
    "amount": 65.88,
    "transactionId": "496f936d-d2c4-466f-862a-d20838917203"
  },
  {
    "state": "win",
    "amount": 31.44,
    "transactionId": "f029d778-0b87-4f15-bb53-cec19def676d"
  },
  {
    "state": "win",
    "amount": 37.84,
    "transactionId": "33293948-b6ae-46a1-bcd1-acbc8643ae21"
  },
  {
    "state": "win",
    "amount": 42.63,
    "transactionId": "468cae6b-138e-4e23-ade4-ec68020cc027"
  },
  {
    "state": "lose",
    "amount": 27.88,
    "transactionId": "fbba64a7-6496-42b3-878a-a06296ad5feb"
  },
  {
    "state": "lose",
    "amount": 95.67,
    "transactionId": "caffbaff-eb6d-442e-a119-24d0555c2045"
  },
  {
    "state": "lose",
    "amount": 2.86,
    "transactionId": "8c45ce70-fa01-49fc-b2bf-f2ee8a6be0ad"
  },
  {
    "state": "win",
    "amount": 86.46,
    "transactionId": "f10a9995-bdae-43b4-81d8-a96f09c73dcd"
  },
  {
    "state": "win",
    "amount": 93.29,
    "transactionId": "af85368a-f35d-4a6a-94d2-af69786a90e5"
  },
  {
    "state": "lose",
    "amount": 5.73,
    "transactionId": "9a92cbb8-ed9b-4860-b788-bb911a7feefe"
  },
  {
    "state": "lose",
    "amount": 76.1,
    "transactionId": "1d83ff5b-ca97-4e82-bbe9-44be718023ac"
  },
  {
    "state": "win",
    "amount": 89.2,
    "transactionId": "49808df4-78e8-428f-b4d9-3bb950cfbdec"
  },
  {
    "state": "win",
    "amount": 53.12,
    "transactionId": "03434314-8e24-4e22-bc1f-2546da77438d"
  },
  {
    "state": "win",
    "amount": 96.58,
    "transactionId": "66dcfdad-2443-47c2-88f4-6dc92d877be6"
  },
  {
    "state": "lose",
    "amount": 74.35,
    "transactionId": "f393e7b2-1c9b-4d8f-85ed-da4440f3fd1a"
  },
  {
    "state": "win",
    "amount": 81.7,
    "transactionId": "af8ec365-50d5-4124-9609-0ea147c847f7"
  },
  {
    "state": "win",
    "amount": 98.18,
    "transactionId": "3a25395c-95ab-48c3-9f1d-388daa933693"
  },
  {
    "state": "lose",
    "amount": 54.68,
    "transactionId": "cfb8d469-cadc-4281-bcf1-9d0dda6f80d8"
  },
  {
    "state": "lose",
    "amount": 93.48,
    "transactionId": "7682cf90-3351-4046-b1e5-4d445739940f"
  },
  {
    "state": "lose",
    "amount": 91.63,
    "transactionId": "1759175d-df52-4800-ab70-c0eb08e57bff"
  },
  {
    "state": "win",
    "amount": 23.01,
    "transactionId": "d30d0d6e-6834-4bfe-ba31-e0bcaa910a7a"
  },
  {
    "state": "lose",
    "amount": 92.06,
    "transactionId": "256afa87-c7ee-4a8f-a667-719248ee1277"
  },
  {
    "state": "lose",
    "amount": 13.92,
    "transactionId": "0d6a7222-5065-4fd2-b743-0196fa96ed6b"
  },
  {
    "state": "lose",
    "amount": 46.65,
    "transactionId": "28ba71d1-7879-495d-b134-91d0029800b1"
  },
  {
    "state": "win",
    "amount": 55.47,
    "transactionId": "35efca6c-158f-4877-acbf-ed3b6625cc46"
  },
  {
    "state": "win",
    "amount": 55.7,
    "transactionId": "3bcfaf80-4379-4cf5-8f9c-37c89e9b2687"
  },
  {
    "state": "win",
    "amount": 71.61,
    "transactionId": "b1a79633-65b3-4e2f-bd13-5a20616e7aa5"
  },
  {
    "state": "lose",
    "amount": 13.08,
    "transactionId": "6f316b26-29e7-49f5-a1f8-44975d359f98"
  },
  {
    "state": "lose",
    "amount": 9.79,
    "transactionId": "7fc2ad37-c8b0-4b33-b655-e8d325c0b5f0"
  },
  {
    "state": "lose",
    "amount": 52.44,
    "transactionId": "ac24ad31-2dea-4771-8ab9-c20efbf66a6d"
  },
  {
    "state": "lose",
    "amount": 81.88,
    "transactionId": "22793c11-afc3-4694-9e1e-fbdd6b337a93"
  },
  {
    "state": "win",
    "amount": 65.53,
    "transactionId": "e8270494-f50b-4019-b320-f5fab766fcb4"
  },
  {
    "state": "lose",
    "amount": 36.56,
    "transactionId": "ff73b027-bad4-4d87-ba08-2cb3aaa90aae"
  },
  {
    "state": "lose",
    "amount": 13.05,
    "transactionId": "688a224d-0a3c-4188-a31b-b322b84e77b3"
  },
  {
    "state": "win",
    "amount": 27.93,
    "transactionId": "23d16be9-f4cc-4fd6-a101-b8bc41215b7a"
  },
  {
    "state": "lose",
    "amount": 57.59,
    "transactionId": "28372986-96ce-461b-b212-2a8d441e1728"
  },
  {
    "state": "win",
    "amount": 38.34,
    "transactionId": "5a8ae034-eeaf-480d-956e-4bf137e88eb7"
  },
  {
    "state": "lose",
    "amount": 73.68,
    "transactionId": "91ec4c3a-92c6-496e-be6f-0029c52fe7c7"
  },
  {
    "state": "win",
    "amount": 0.35,
    "transactionId": "2f172be2-be01-438b-9b0b-63a357c5fcd7"
  },
  {
    "state": "lose",
    "amount": 76.03,
    "transactionId": "e1118701-bf95-40cc-ba5c-2da464c0be8c"
  },
  {
    "state": "lose",
    "amount": 99.59,
    "transactionId": "5b2a33be-b7dd-4007-9c54-55c795d51847"
  },
  {
    "state": "lose",
    "amount": 2.35,
    "transactionId": "a8381b70-a300-4021-a79d-dfb893dfccff"
  },
  {
    "state": "lose",
    "amount": 77.43,
    "transactionId": "e6669199-4623-4d63-a5eb-2b946e71c0c8"
  },
  {
    "state": "win",
    "amount": 67.6,
    "transactionId": "cfb4bab9-4f69-4363-80ab-703e500a7a2d"
  },
  {
    "state": "win",
    "amount": 67.72,
    "transactionId": "5cb0aeb2-8623-45e5-8584-8399453c73e7"
  },
  {
    "state": "lose",
    "amount": 91.81,
    "transactionId": "c132149c-efb5-4f75-b00b-5901f29b7ec2"
  },
  {
    "state": "win",
    "amount": 82.33,
    "transactionId": "c3f1a48e-3e0e-416e-ad22-54f4e3dbb31f"
  },
  {
    "state": "win",
    "amount": 0.4,
    "transactionId": "6296ed61-eea5-4d29-b39f-1c8619d43dcc"
  },
  {
    "state": "win",
    "amount": 82.35,
    "transactionId": "c2141acc-d12a-4cad-87f4-7484313cef43"
  },
  {
    "state": "lose",
    "amount": 3.74,
    "transactionId": "488c03c2-e0b9-452b-848b-435d359fcff2"
  },
  {
    "state": "win",
    "amount": 5.35,
    "transactionId": "0f3a3782-39fd-4f03-ad7d-efb623f12ffd"
  },
  {
    "state": "lose",
    "amount": 18.92,
    "transactionId": "f007749c-f874-48c6-aef3-a640d52eecb4"
  },
  {
    "state": "win",
    "amount": 65,
    "transactionId": "2bf97e93-f38b-4222-9382-477ba0091fc8"
  },
  {
    "state": "win",
    "amount": 74.94,
    "transactionId": "9fe92907-8da4-4766-bf79-3828efdf2333"
  },
  {
    "state": "lose",
    "amount": 24.64,
    "transactionId": "0be0c53d-d54b-4f27-9e0b-74063ff6e4a8"
  },
  {
    "state": "lose",
    "amount": 21.06,
    "transactionId": "3bef82fe-f375-4107-a06f-4cda5cec0dc8"
  },
  {
    "state": "lose",
    "amount": 64.41,
    "transactionId": "59746532-5715-4c2f-9514-d3f64539f051"
  },
  {
    "state": "win",
    "amount": 4.38,
    "transactionId": "d84de2ce-0bf6-4923-b980-42da342d2037"
  },
  {
    "state": "lose",
    "amount": 22.77,
    "transactionId": "bdf29ab1-6a99-4769-8dfb-088ba28a543f"
  },
  {
    "state": "win",
    "amount": 47.15,
    "transactionId": "81ff7d7b-388b-48c8-b15e-7070682d5658"
  },
  {
    "state": "win",
    "amount": 34.08,
    "transactionId": "5f04ed50-29e0-4f79-bcc7-26e076a9b6f0"
  },
  {
    "state": "win",
    "amount": 7.42,
    "transactionId": "6aac4daf-8194-4af2-b10e-34d9b39d602c"
  },
  {
    "state": "win",
    "amount": 27.47,
    "transactionId": "54c87821-7770-4a3b-8c2a-dc46087b34b6"
  },
  {
    "state": "lose",
    "amount": 77.67,
    "transactionId": "4b4cc31e-17e5-4a49-8545-1eb7d11031f5"
  },
  {
    "state": "win",
    "amount": 87.68,
    "transactionId": "827e2992-ddf0-4cc9-bd88-66b0eec2e49e"
  },
  {
    "state": "win",
    "amount": 6.33,
    "transactionId": "e1a404fb-35c5-4146-8acd-ab79a5e85660"
  },
  {
    "state": "win",
    "amount": 20.68,
    "transactionId": "3b921b46-d3f8-440f-b49b-3ba7325b6ec3"
  },
  {
    "state": "win",
    "amount": 73.18,
    "transactionId": "803a33e8-7d3d-4170-beef-6190d5c03fb2"
  },
  {
    "state": "lose",
    "amount": 50.06,
    "transactionId": "767011d0-4d1e-4840-821a-3cd8df29c55d"
  },
  {
    "state": "lose",
    "amount": 12.29,
    "transactionId": "0c992d4f-e881-46d3-a4d8-362ad64c97eb"
  },
  {
    "state": "lose",
    "amount": 59.74,
    "transactionId": "6b2df149-39c0-46ba-99fa-80dc1e984542"
  },
  {
    "state": "win",
    "amount": 84.22,
    "transactionId": "684a2ca1-08d2-45a7-acf4-90eeb2593a0e"
  },
  {
    "state": "lose",
    "amount": 46.46,
    "transactionId": "377ee754-0222-43a3-a162-b2a5230175e9"
  },
  {
    "state": "win",
    "amount": 41.19,
    "transactionId": "dda4e25e-4d69-45b6-8f8f-d5b5f7762288"
  },
  {
    "state": "lose",
    "amount": 88.9,
    "transactionId": "1cd584d5-a7d5-4d51-ba44-65aa48592427"
  },
  {
    "state": "win",
    "amount": 30.95,
    "transactionId": "2187c728-8c1d-454f-9547-c1a5edc7cbc5"
  },
  {
    "state": "lose",
    "amount": 32.53,
    "transactionId": "4e9d541f-38aa-45e2-929e-ae12ee8e453a"
  },
  {
    "state": "lose",
    "amount": 58.18,
    "transactionId": "f69151d7-8ca0-4f79-85f0-53188d6d7752"
  },
  {
    "state": "lose",
    "amount": 27.94,
    "transactionId": "feab2085-92c1-46ba-8fe9-046e948211ee"
  },
  {
    "state": "lose",
    "amount": 54.04,
    "transactionId": "5b1a12c5-587a-4640-ac5e-928c5417dc8a"
  },
  {
    "state": "win",
    "amount": 26.02,
    "transactionId": "0f9fbd9e-f308-4804-bb81-646f617cd2da"
  },
  {
    "state": "win",
    "amount": 28.55,
    "transactionId": "130814ba-36cd-480b-8558-0453df4e050c"
  },
  {
    "state": "lose",
    "amount": 99.42,
    "transactionId": "3d06a7b4-7fa3-4894-a24f-f5017f4dde09"
  },
  {
    "state": "lose",
    "amount": 9.18,
    "transactionId": "35a3bb0b-00e3-4001-81eb-686bb2d1ea80"
  },
  {
    "state": "lose",
    "amount": 69.84,
    "transactionId": "7d9a39ae-d497-4879-8783-cc4f7319b4a3"
  },
  {
    "state": "win",
    "amount": 71.07,
    "transactionId": "58efa8a3-e9d0-42b6-88b8-8a9bc262c60b"
  },
  {
    "state": "lose",
    "amount": 82.88,
    "transactionId": "490f8e83-dac2-4611-90d7-bff74195d635"
  },
  {
    "state": "win",
    "amount": 11.67,
    "transactionId": "73b7c804-feba-4b36-b436-1066fdc0038d"
  },
  {
    "state": "win",
    "amount": 74.76,
    "transactionId": "89350742-22f7-491a-8033-0b6db7082abf"
  },
  {
    "state": "win",
    "amount": 61.21,
    "transactionId": "e5c5be3f-97fd-4a90-82e2-c0527da808ba"
  },
  {
    "state": "win",
    "amount": 60.51,
    "transactionId": "50f88594-ef12-48b4-8509-0bf9aa9f9b83"
  },
  {
    "state": "lose",
    "amount": 87.82,
    "transactionId": "a03818a5-fd93-42fb-8b67-a058c87ec1ff"
  },
  {
    "state": "win",
    "amount": 97.15,
    "transactionId": "6e55502a-1e81-4bbd-b2fe-a986483415f3"
  },
  {
    "state": "lose",
    "amount": 72.86,
    "transactionId": "5a0894ab-c5df-43aa-8056-566493a7627f"
  },
  {
    "state": "lose",
    "amount": 86.8,
    "transactionId": "d85a5df1-11d2-4c2d-b449-44a8d51e3644"
  },
  {
    "state": "win",
    "amount": 13.83,
    "transactionId": "4a38fa3a-ef38-4af8-bb34-84c82cc3ad93"
  },
  {
    "state": "win",
    "amount": 60.43,
    "transactionId": "b656dd8c-46bc-45d1-8155-e6f8e87c1db1"
  },
  {
    "state": "win",
    "amount": 51.23,
    "transactionId": "1268c2be-f107-4704-bcc3-b542d916723d"
  },
  {
    "state": "lose",
    "amount": 35.6,
    "transactionId": "32b371fe-da85-4621-9cc9-195bfa4c6bed"
  },
  {
    "state": "win",
    "amount": 95.97,
    "transactionId": "5f805bb3-c37a-4958-bed9-d8593de83562"
  },
  {
    "state": "win",
    "amount": 7.66,
    "transactionId": "b667e11a-6b2f-45f6-8105-c9340b93303d"
  },
  {
    "state": "win",
    "amount": 28.29,
    "transactionId": "2a713913-efc3-4ec3-8ac8-94cfd4e96c69"
  },
  {
    "state": "lose",
    "amount": 76.64,
    "transactionId": "bf610e0d-b780-4aec-bb27-fb926d96951d"
  },
  {
    "state": "win",
    "amount": 45.07,
    "transactionId": "526e7582-a45e-45cd-84f8-937c07bde394"
  },
  {
    "state": "win",
    "amount": 89.22,
    "transactionId": "e95f1451-d128-4f2c-8d10-99c821bff090"
  },
  {
    "state": "win",
    "amount": 61.87,
    "transactionId": "953535a8-2389-4121-b07e-aa4f1e5131bc"
  },
  {
    "state": "win",
    "amount": 12.6,
    "transactionId": "23f64f33-5197-4d37-b46d-2a089dc3fe73"
  },
  {
    "state": "win",
    "amount": 3.33,
    "transactionId": "f669e8c8-8784-4e48-bd50-8d8e80074ccb"
  },
  {
    "state": "win",
    "amount": 26.73,
    "transactionId": "0c5c417a-0047-4f95-a37e-51caf67302a2"
  },
  {
    "state": "win",
    "amount": 52.24,
    "transactionId": "dd5e1942-9093-4137-9dbd-0f0cf6fc05e4"
  },
  {
    "state": "win",
    "amount": 37.75,
    "transactionId": "b6a4c5f4-c0ef-4267-8fd5-66a55f3505ee"
  },
  {
    "state": "win",
    "amount": 6.78,
    "transactionId": "ceef0d14-3816-4fe3-8655-aae1a1389fed"
  },
  {
    "state": "lose",
    "amount": 64.61,
    "transactionId": "044c6f5e-4872-42d9-a7ce-1d86124c4a6d"
  },
  {
    "state": "win",
    "amount": 95.67,
    "transactionId": "caa534ea-771a-4d2d-8370-97aa0a85e6fe"
  },
  {
    "state": "lose",
    "amount": 87.1,
    "transactionId": "a6a1f207-b6af-42bb-82f6-808815dff79a"
  },
  {
    "state": "win",
    "amount": 15.33,
    "transactionId": "64e65fe4-0d84-415b-8200-ed80362a2e30"
  },
  {
    "state": "win",
    "amount": 35.26,
    "transactionId": "aab753ea-1462-4bdf-b413-b0d4e27a38b6"
  },
  {
    "state": "win",
    "amount": 34.06,
    "transactionId": "1c7283f4-3982-4222-8106-b0bce295601a"
  },
  {
    "state": "win",
    "amount": 19.85,
    "transactionId": "3dfb663a-8127-4d10-a90a-feaf2064d61a"
  },
  {
    "state": "win",
    "amount": 1.24,
    "transactionId": "e68e13bb-4c43-47e3-9353-3c3f22b8dcbd"
  },
  {
    "state": "win",
    "amount": 62.8,
    "transactionId": "60f4c7cf-da2c-43c3-9d19-f2451123d358"
  },
  {
    "state": "win",
    "amount": 77.62,
    "transactionId": "d41e6f86-225b-4821-9b9b-7e3897099cbf"
  },
  {
    "state": "lose",
    "amount": 59.64,
    "transactionId": "9cc577b8-40f9-4d12-a59e-f8f3de11df4d"
  },
  {
    "state": "win",
    "amount": 51.75,
    "transactionId": "7b903272-9c96-4516-81f8-35180999e3cc"
  },
  {
    "state": "win",
    "amount": 99.02,
    "transactionId": "ffcb224b-b872-4d62-a475-018642da6a02"
  },
  {
    "state": "win",
    "amount": 38.21,
    "transactionId": "b55cc614-ad9f-454f-a655-0492e6a7320c"
  },
  {
    "state": "lose",
    "amount": 84.83,
    "transactionId": "f9820e6f-63cb-4708-97ba-b71f7e2d1efe"
  },
  {
    "state": "lose",
    "amount": 97.89,
    "transactionId": "b7a0c5a9-be97-4edd-b671-7463356ec3a9"
  },
  {
    "state": "lose",
    "amount": 39.73,
    "transactionId": "52fc20ee-14d6-47c5-9d3e-5abf42cad7a7"
  },
  {
    "state": "win",
    "amount": 63.86,
    "transactionId": "028d8953-6ed8-4126-94fd-3a1e85f4e1a9"
  },
  {
    "state": "lose",
    "amount": 35.24,
    "transactionId": "ded166a0-59bb-46e6-9112-eb4c9921a456"
  },
  {
    "state": "lose",
    "amount": 11.12,
    "transactionId": "de9c3298-0a9a-43f1-9bd2-b72071d819e4"
  },
  {
    "state": "win",
    "amount": 65.97,
    "transactionId": "68d963a7-373f-4992-a299-c193b8203cf1"
  },
  {
    "state": "lose",
    "amount": 76.68,
    "transactionId": "a1f36ef9-55b3-4b88-9423-7f3332bab114"
  },
  {
    "state": "win",
    "amount": 39.5,
    "transactionId": "468c29c8-a4f5-4162-84a6-f7432e3a3248"
  },
  {
    "state": "win",
    "amount": 77.29,
    "transactionId": "8fe04782-306d-46b6-bebd-51af28df9086"
  },
  {
    "state": "win",
    "amount": 57.04,
    "transactionId": "47689d0c-b183-4715-a406-da2e2da012a4"
  },
  {
    "state": "win",
    "amount": 40.66,
    "transactionId": "6c17323f-104f-43c6-917a-bfe53aca46e4"
  },
  {
    "state": "lose",
    "amount": 21.14,
    "transactionId": "5103fcba-3a16-4f22-aa6f-a6dd15062587"
  },
  {
    "state": "lose",
    "amount": 91.57,
    "transactionId": "bb8bf674-c351-4494-a333-8e855fd4cba2"
  },
  {
    "state": "win",
    "amount": 28.81,
    "transactionId": "f06a2119-600f-4ae2-bc7f-7e6bdc3cc812"
  },
  {
    "state": "win",
    "amount": 79.54,
    "transactionId": "72857d01-34d3-4cd6-b188-6336064acbe5"
  },
  {
    "state": "lose",
    "amount": 54.71,
    "transactionId": "07bd54c0-1f88-4739-ae1e-2ac1bb76b534"
  },
  {
    "state": "win",
    "amount": 94.38,
    "transactionId": "bb8da45e-de1a-4dc6-bbb1-b98e9fe18171"
  },
  {
    "state": "win",
    "amount": 32.36,
    "transactionId": "9e076acb-d89c-4239-9a19-da5f6d54912b"
  },
  {
    "state": "win",
    "amount": 79.72,
    "transactionId": "6189be0e-c754-4a0d-acbf-9a414d3ca362"
  },
  {
    "state": "lose",
    "amount": 6,
    "transactionId": "7056e5ec-3520-470a-8862-da31abf45f32"
  },
  {
    "state": "lose",
    "amount": 68.69,
    "transactionId": "63ac1988-6610-4a10-9316-e96c2be9ae78"
  },
  {
    "state": "win",
    "amount": 51.99,
    "transactionId": "ea9deb30-3464-46d5-9990-bef85e7823eb"
  },
  {
    "state": "win",
    "amount": 6.96,
    "transactionId": "77dab3dd-7252-4b71-809d-b683c3caec4d"
  },
  {
    "state": "win",
    "amount": 1.07,
    "transactionId": "c92a55b2-1ed1-4098-95ca-0043090122a4"
  },
  {
    "state": "win",
    "amount": 13.69,
    "transactionId": "203a6b0f-9a2c-4cd4-8e49-4f6f253a2512"
  },
  {
    "state": "lose",
    "amount": 32.55,
    "transactionId": "0502362f-cda9-4582-be3f-24ec407591b7"
  },
  {
    "state": "lose",
    "amount": 15.42,
    "transactionId": "ae4e20cd-721b-4400-8ce4-20b3b2c5dc84"
  },
  {
    "state": "lose",
    "amount": 99.99,
    "transactionId": "7f14b0aa-ae01-4d33-b766-f1c33036d461"
  },
  {
    "state": "lose",
    "amount": 6.5,
    "transactionId": "cc840f55-bad2-4124-b358-8580b9478b0b"
  },
  {
    "state": "lose",
    "amount": 46.91,
    "transactionId": "8c2e21f1-a098-4247-b606-0b7b6b35f4ee"
  },
  {
    "state": "win",
    "amount": 42.38,
    "transactionId": "ee33425f-594a-428e-a42f-c13e254e3ae7"
  },
  {
    "state": "lose",
    "amount": 52.22,
    "transactionId": "f64a51dd-5ba0-49ea-875d-4cb958961d21"
  },
  {
    "state": "lose",
    "amount": 46.59,
    "transactionId": "34370643-bb99-4a51-bec4-1e9235c1a8fe"
  },
  {
    "state": "lose",
    "amount": 4.5,
    "transactionId": "d8406106-d328-4fa3-ae36-5d838cb8431b"
  },
  {
    "state": "win",
    "amount": 22.79,
    "transactionId": "5b3fd6c9-d300-4165-825f-57d277504545"
  },
  {
    "state": "lose",
    "amount": 44.04,
    "transactionId": "8605b7c1-5948-4130-b318-21f307a9b0cc"
  },
  {
    "state": "lose",
    "amount": 10.2,
    "transactionId": "4207e6dc-0b17-4992-85dd-ea09250827be"
  },
  {
    "state": "lose",
    "amount": 77.29,
    "transactionId": "b78890ca-d714-4cac-86ee-33322cc650c6"
  },
  {
    "state": "win",
    "amount": 7.06,
    "transactionId": "45740a98-99d0-40bb-99dc-db4bda0acb83"
  },
  {
    "state": "lose",
    "amount": 75.93,
    "transactionId": "8c515777-4ad6-4fd8-a7da-bdd2d6b38279"
  },
  {
    "state": "lose",
    "amount": 72.35,
    "transactionId": "3576b845-7704-4287-8122-f5db143e8ab1"
  },
  {
    "state": "lose",
    "amount": 89.6,
    "transactionId": "a7159d10-adb3-425e-8309-fb73ad51f2d9"
  }
]`
)
