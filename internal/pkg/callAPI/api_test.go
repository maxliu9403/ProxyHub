package callAPi

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestCallAPI(t *testing.T) {
	p := map[string]interface{}{
		"Id": 113,
	}
	var resp interface{}
	err := CallAPI(context.TODO(), "url", p, &resp,
		SetHeader(map[string]string{"Content-Type": "application/json"}), SetTimeout(1*time.Second), SetMethod(HTTPPost))
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(resp)
}
