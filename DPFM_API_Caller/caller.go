package dpfm_api_caller

import (
	"context"
	dpfm_api_input_reader "data-platform-api-district-exconf-rmq-kube/DPFM_API_Input_Reader"
	dpfm_api_output_formatter "data-platform-api-district-exconf-rmq-kube/DPFM_API_Output_Formatter"
	"encoding/json"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	database "github.com/latonaio/golang-mysql-network-connector"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client-for-data-platform"
)

type ExistenceConf struct {
	ctx context.Context
	db  *database.Mysql
	l   *logger.Logger
}

func NewExistenceConf(ctx context.Context, db *database.Mysql, l *logger.Logger) *ExistenceConf {
	return &ExistenceConf{
		ctx: ctx,
		db:  db,
		l:   l,
	}
}

func (e *ExistenceConf) Conf(msg rabbitmq.RabbitmqMessage) interface{} {
	var ret interface{}
	ret = map[string]interface{}{
		"ExistenceConf": false,
	}
	input := make(map[string]interface{})
	err := json.Unmarshal(msg.Raw(), &input)
	if err != nil {
		return ret
	}

	_, ok := input["District"]
	if ok {
		input := &dpfm_api_input_reader.SDC{}
		err = json.Unmarshal(msg.Raw(), input)
		ret = e.confDistrict(input)
		goto endProcess
	}

endProcess:
	if err != nil {
		e.l.Error(err)
	}
	return ret
}

func (e *ExistenceConf) confDistrict(input *dpfm_api_input_reader.SDC) *dpfm_api_output_formatter.District {
	exconf := dpfm_api_output_formatter.District{
		ExistenceConf: false,
	}
	if input.District.District == nil {
		return &exconf
	}
	if input.District.Country == nil {
		return &exconf
	}
	exconf = dpfm_api_output_formatter.District{
		District:      *input.District.District,
		Country:       *input.District.Country,
		ExistenceConf: false,
	}

	rows, err := e.db.Query(
		`SELECT District
		FROM DataPlatformMastersAndTransactionsMysqlKube.data_platform_district_district_data 
		WHERE (District, Country) = (?, ?);`, exconf.District, exconf.Country,
	)
	if err != nil {
		e.l.Error(err)
		return &exconf
	}

	exconf.ExistenceConf = rows.Next()
	return &exconf
}
