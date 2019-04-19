// Copyright (c) 2018-2019 The Cybavo developers
// All Rights Reserved.
// NOTICE: All information contained herein is, and remains
// the property of Cybavo and its suppliers,
// if any. The intellectual and technical concepts contained
// herein are proprietary to Cybavo
// Dissemination of this information or reproduction of this materia
// is strictly forbidden unless prior written permission is obtained
// from Cybavo.

package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/cybavo/SOFA_MOCK_SERVER/api"
	"github.com/cybavo/SOFA_MOCK_SERVER/models"
)

type OuterController struct {
	beego.Controller
}

func (c *OuterController) AbortWithError(status int, err error) {
	resp := api.ErrorCodeResponse{
		ErrMsg:  err.Error(),
		ErrCode: status,
	}
	c.Data["json"] = resp
	c.Abort(strconv.Itoa(status))
}

// @Title Set API token
// @router /wallets/:wallet_id/apitoken [post]
func (c *OuterController) SetAPIToken() {
	defer c.ServeJSON()

	walletID, err := strconv.ParseInt(c.Ctx.Input.Param(":wallet_id"), 10, 64)
	if err != nil {
		logs.Error("Invalid walled ID =>", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var request api.SetAPICodeRequest
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	apiCodeParams := models.APICode{
		APICode:   request.APICode,
		ApiSecret: request.ApiSecret,
		WalletID:  walletID,
	}
	err = models.SetAPICode(&apiCodeParams)
	if err != nil {
		logs.Error("SetAPICode failed", err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	response := &api.CommonResponse{
		Result: 1,
	}
	c.Data["json"] = response
}

// @Title Create deposit wallet addresses
// @router /wallets/:wallet_id/addresses [post]
func (c *OuterController) CreateDepositWalletAddresses() {
	defer c.ServeJSON()

	walletID, err := strconv.ParseInt(c.Ctx.Input.Param(":wallet_id"), 10, 64)
	if err != nil {
		logs.Error("Invalid wallet ID =>", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}

	var request api.CreateDepositWalletAddressesRequest
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	resp, err := api.CreateDepositWalletAddresses(walletID, &request)
	if err != nil {
		logs.Error("CreateDepositWalletAddresses failed", err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	var walletAddresses []models.DepositWalletAddress
	for _, address := range resp.Addresses {
		walletAddresses = append(walletAddresses, models.DepositWalletAddress{
			Address:  address,
			WalletID: walletID,
		})
	}
	_, err = models.AddNewWalletAddresses(walletAddresses)
	if err != nil {
		logs.Error("AddNewWalletAddresses failed", err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.Data["json"] = resp
}

// @Title Get deposit wallet addresses
// @router /wallets/:wallet_id/addresses [get]
func (c *OuterController) GetDepositWalletAddresses() {
	defer c.ServeJSON()

	walletID, err := strconv.ParseInt(c.Ctx.Input.Param(":wallet_id"), 10, 64)
	if err != nil {
		logs.Error("Invalid wallet ID =>", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	startIndex, _ := c.GetInt("start_index", 0)
	requestNumber, _ := c.GetInt("request_number", 1000)

	resp, err := api.GetDepositWalletAddresses(walletID, startIndex, requestNumber)
	if err != nil {
		logs.Error("GetDepositWalletAddresses failed", err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Data["json"] = resp
}

// @Title Callback
// @router /wallets/callback [post]
func (c *OuterController) Callback() {
	var request api.CallbackRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	logs.Debug("Callback => %s\n%#v", c.Ctx.Input.RequestBody, request)

	c.Ctx.WriteString("OK")
}

// @Title Resend Callback
// @router /wallets/:wallet_id/callback/resend [post]
func (c *OuterController) CallbackResend() {
	defer c.ServeJSON()

	walletID, err := strconv.ParseInt(c.Ctx.Input.Param(":wallet_id"), 10, 64)
	if err != nil {
		logs.Error("Invalid wallet ID =>", err)
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	var request api.CallbackResendRequest
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	resp, err := api.ResendCallback(walletID, &request)
	if err != nil {
		logs.Error("ResendCallback failed", err)
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.Data["json"] = resp
}
