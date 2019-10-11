package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/siddontang/go-log/log"
	"io/ioutil"
	"net/http"
	"strconv"
)

type gbeOrder struct {
	Id    int64           `json:"id"`
	Price decimal.Decimal `json:"price"`
	Side  string          `json:"side"`
}

func placeOrder(gbeToken string, productId string, size string, price string, funds string, side string, orderType string) (*gbeOrder, error) {
	params := map[string]interface{}{}
	params["productId"] = productId
	params["side"] = side
	params["type"] = orderType
	params["price"], _ = strconv.ParseFloat(price, 10)
	params["size"], _ = strconv.ParseFloat(size, 10)
	params["funds"], _ = strconv.ParseFloat(funds, 10)

	log.Infof("new order : %+v", params)

	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(fmt.Sprintf("%v/api/orders/?token=%v", gitBitExAddr, gbeToken), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(productId + " " + resp.Status + " " + string(buf))
	}

	var order gbeOrder
	err = json.Unmarshal(buf, &order)
	return &order, err
}

func cancelOrders(gbeToken, productId string, side string) error {
	log.Infof("cancel all : %v %v", productId, side)
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%v/api/orders?productId=%v&side=%v&token=%v",
		gitBitExAddr, productId, side, gbeToken), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(buf))
	}
	return nil
}

func cancelOrder(gbeToken string, orderId int64) error {
	log.Infof("cancellOrder %v", orderId)
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%v/api/orders/%v?token=%v",
		gitBitExAddr, orderId, gbeToken), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(buf))
	}
	return nil
}

func cancelOrderByClientOid(gbeToken, clientOid string) error {
	log.Infof("cancellOrderByClientOid %v", clientOid)
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%v/api/orders/client:%v?token=%v",
		gitBitExAddr, clientOid, gbeToken), nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		buf, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(buf))
	}
	return nil
}

func getToken(email, password string) (string, error) {
	params := map[string]interface{}{}
	params["email"] = email
	params["password"] = password

	data, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(fmt.Sprintf("%v/api/users/token", gitBitExAddr), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(email + " : " + string(buf))
	}
	return string(buf), nil
}
