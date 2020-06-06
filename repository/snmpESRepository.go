package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/sirupsen/logrus"
	"os"
)

//go:generate mockgen -destination=../mocks/repository/snmpESRepositoryMock.go -package=repository ntm-backend/repository ISnmpESRepository

type ISnmpESRepository interface {
	SnmpSearchQuery(query map[string]interface{}) (*esapi.Response, error)
}

type SnmpESRepository struct {
	Logger   *logrus.Logger
	ESClient *elasticsearch.Client
}

func (ses SnmpESRepository) SnmpSearchQuery(query map[string]interface{}) (*esapi.Response, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		ses.Logger.Println("Error encoding query: ", err)
		return nil, err
	} else {
		res, err := ses.ESClient.Search(
			ses.ESClient.Search.WithContext(context.Background()),
			ses.ESClient.Search.WithIndex(os.Getenv("SNMP_ELASTICSEARCH_INDEX")),
			ses.ESClient.Search.WithBody(&buf),
			ses.ESClient.Search.WithTrackTotalHits(true),
			ses.ESClient.Search.WithPretty(),
		)
		if err != nil {
			return &esapi.Response{}, err
		}
		return res, nil
	}
}
