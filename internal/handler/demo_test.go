/*
@Date: 2022/4/22 12:26
@Author: max.liu
@File : demo_test
@Desc:
*/

package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxliu9403/go-template/internal/common"
	"github.com/maxliu9403/go-template/internal/logic/demo"
)

func TestGetDetail(t *testing.T) {
	params, _ := json.Marshal(&demo.IDParams{
		Test: common.Test{true},
		ID:   1,
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/?Action=GetDetail", bytes.NewBuffer(params))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(userHeader, userHeaderValue)
	req.Header.Set(orgHeader, orgHeaderValue)
	setupRouter().ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatal(w.Code)
	}

	res := common.Response{}
	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Fatal(err.Error())
	}

	respBytes, _ := json.Marshal(res.DataSet)
	t.Logf("message is %s, data is %s", res.Message, string(respBytes))
}
