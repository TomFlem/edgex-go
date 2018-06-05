package distro

import (
	"bytes"
	"encoding/json"
	"encoding/gob"
	"strconv"
	"io"
	"net/http"
	"math"
	"time"

	"github.com/edgexfoundry/edgex-go/core/domain/models"
	"go.uber.org/zap"
)

type forumFormatter struct {
}

type WriteObjectRequest struct {
	Mode struct {
		Value float64 `json:"value"`
		Origin int64 `json:"origin"`
	} `json:"mode"`
}

type RequestSequence struct {
	Reqs [][]byte
}

func (jsonTr forumFormatter) Format(event *models.Event) []byte {
	logger.Info("**** FORUM FORMAT ****")
	if len(event.Readings) == 0 {
		logger.Error("No reading in event")
		return nil
	}
	rs := RequestSequence{}
	rs.Reqs = make([][]byte, len(event.Readings))
	for i := 0; i < len(event.Readings); i++ {
		v, err := strconv.ParseFloat(event.Readings[i].Value, 64)
		if err != nil {
			logger.Error("Error parsing reading value", zap.Error(err))
			return nil
		}
		req := WriteObjectRequest{}
		req.Mode.Value = v
		req.Mode.Origin = event.Readings[i].Origin
		if math.IsNaN(req.Mode.Value) {
			req.Mode.Value = 0;
		}
		b, err := json.Marshal(req)
		logger.Info("Forum data", zap.ByteString("data", b))
		if err != nil {
			logger.Error("Error marshaling JSON", zap.Error(err))
			return nil
		}
		rs.Reqs[i] = b
	}
	buf := new(bytes.Buffer)
	err := gob.NewEncoder(buf).Encode(rs)
	if err != nil {
			logger.Error("Error encoding readings", zap.Error(err))
			return nil
	}
	return buf.Bytes()
}

func NewFORUMFormat() Formatter {
	return forumFormatter{}
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
		buf := bytes.NewBuffer(data)
		rs := RequestSequence{}
		err := gob.NewDecoder(buf).Decode(&rs)
		if err != nil {
			logger.Error("Error decoding: ", zap.Error(err))
			return
		}
		for i := 0; i < len(event.Readings); i++ {
			url := sender.url + "/" + event.Readings[i].Device + "-" + event.Readings[i].Name
			response, err := put(url, mimeTypeJSON, bytes.NewReader(rs.Reqs[i]))
			if err != nil {
				logger.Error("Error: ", zap.Error(err))
				return
			}
			defer response.Body.Close()
			logger.Info("Response: ", zap.String("status", response.Status))
			logger.Info("Sent put data: ", zap.ByteString("data", data))
		}
	} else {
		sender.httpSender.Send(event, data)
	}
}
