package distro

import (
	"bytes"
	"encoding/json"
	"strconv"
	"io"
	"net/http"
	"time"

	"github.com/edgexfoundry/edgex-go/core/domain/models"
	"go.uber.org/zap"
)

type forumFormater struct {
}

type WriteObjectRequest struct {
	Mode struct {
		Value float64 `json:"value"`
		Origin int64 `json:"origin"`
	} `json:"mode"`
}

func (jsonTr forumFormater) Format(event *models.Event) []byte {
	logger.Info("**** FORUM FORMAT ****")
	if len(event.Readings) == 0 {
		logger.Error("No reading in event")
		return nil
	}
	v, err := strconv.ParseFloat(event.Readings[0].Value, 64)
	if err != nil {
		logger.Error("Error parsing reading value", zap.Error(err))
		return nil
	}
	req := WriteObjectRequest{}
	req.Mode.Value = v
	req.Mode.Origin = event.Readings[0].Origin
	b, err := json.Marshal(req)
	logger.Info("Forum data", zap.ByteString("data", b))
	if err != nil {
		logger.Error("Error parsing JSON", zap.Error(err))
		return nil
	}
	return b
}

func NewFORUMFormat() Formater {
	return forumFormater{}
}

type forumHttpSender struct {
	httpSender
}

func NewFORUMHTTPSender(addr models.Addressable) Sender {
	result := forumHttpSender{}
	sender, ok := NewHTTPSender(addr).(httpSender)
	if ok {
		result.httpSender = sender
	}
	return result
}

func put(url string, contentType string, data io.Reader) (resp *http.Response, err error) {
	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodPut, url, data)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

func (sender forumHttpSender) Send(event *models.Event, data []byte) {
	now := time.Now().UnixNano() / 1000000
	delta := now - event.Readings[0].Origin
	logger.Info("time delta", zap.Int64("delta", delta))
	if sender.method == http.MethodPut {
		if len(event.Readings) == 0 {
			logger.Error("No reading in event")
		}
		url := sender.url + "/" + event.Readings[0].Device + "-" + event.Readings[0].Name
		response, err := put(url, mimeTypeJSON, bytes.NewReader(data))
		if err != nil {
			logger.Error("Error: ", zap.Error(err))
			return
		}
		defer response.Body.Close()
		logger.Info("Response: ", zap.String("status", response.Status))
		logger.Info("Sent put data: ", zap.ByteString("data", data))
	} else {
		sender.httpSender.Send(event, data)
	}
}
