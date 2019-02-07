package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"strconv"
	"sumwhere/factory"
	"sumwhere/models"
	"sumwhere/utils"
)

type PurchaseController struct {
}

func (p PurchaseController) Init(g *echo.Group) {
	g.POST("/purchase/receipt", p.ReceiptValidate)
	g.GET("/purchase/list", p.GetAllProducts)
	g.GET("/purchase/info", p.GetProductInfo)
	g.GET("/purchase/history/:page", p.GetPurchaseHistory)
}

func (p PurchaseController) ReceiptValidate(e echo.Context) error {
	receipt := e.Request().PostFormValue("receiptdata")
	identifier := e.Request().PostFormValue("identifier")

	result := make(map[string]string)
	result["receipt-data"] = receipt

	data := convertReceiptFromServer(result)
	for data := range data {
		if data.Status == 0 {

			models.User{}.UpdateKeyWithHistory(e, identifier)

			return utils.ReturnApiSucc(e, http.StatusOK, true)
		} else {
			return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorNotFound, errors.New("에러"))
		}
	}

	return utils.ReturnApiSucc(e, http.StatusOK, "")
}

func convertReceiptFromServer(receipt map[string]string) <-chan models.Receipt {
	out := make(chan models.Receipt)

	go func() {
		jsonData, err := json.Marshal(receipt)
		if err != nil {
			return
		}

		resp, err := http.Post("https://sandbox.itunes.apple.com/verifyReceipt", "Application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return
		}
		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return
		}

		var receiptModel models.Receipt
		json.Unmarshal(respBody, &receiptModel)
		if err != nil {
			return
		}
		out <- receiptModel
		close(out)
	}()

	return out
}

// 모든 상품 리스트
func (p PurchaseController) GetAllProducts(e echo.Context) error {
	var keyModels []models.PurchaseProduct
	factory.DB(e.Request().Context()).Find(&keyModels)

	return utils.ReturnApiSucc(e, http.StatusOK, keyModels)
}

// 스토어 상품 정보
func (p PurchaseController) GetProductInfo(e echo.Context) error {
	productName := e.QueryParam("productName")
	var product models.PurchaseProduct
	result, err := factory.DB(e.Request().Context()).Where("product_name = ?", productName).Get(&product)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorNotFound, err)
	}

	if !result {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorNotFound, errors.New("아이템이 없습니다."))
	}
	return utils.ReturnApiSucc(e, http.StatusOK, product)
}

// 키 구매 히스토리
func (p PurchaseController) GetPurchaseHistory(e echo.Context) error {

	page, err := strconv.Atoi(e.Param("page"))
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusOK, utils.ApiErrorParameter, errors.New("숫자가 아닙니다."))
	}

	users := e.Get("user").(*jwt.Token)
	claims := users.Claims.(*models.JwtCustomClaims)

	result, err := models.PurchaseHistory{}.GetUserHistory(e.Request().Context(), claims.Id, page)
	if err != nil {
		return utils.ReturnApiFail(e, http.StatusBadRequest, utils.ApiErrorNotFound, err)
	}

	return utils.ReturnApiSucc(e, http.StatusOK, result)

}
