package driver

import (
	"fmt"
	"hcc/piano/lib/influxdb"
	"hcc/piano/lib/logger"
	"hcc/piano/model"
	"strconv"

	"github.com/influxdata/influxdb1-client/models"
)

func GetInfluxData(args map[string]interface{}) (interface{}, error) {

	metric, metricOk := args["metric"].(string)
	subMetric, subMetricOk := args["subMetric"].(string)
	period, periodOk := args["period"].(string)
	aggregateType, aggregateTypeOk := args["aggregateType"].(string)
	duration, durationOk := args["duration"].(string)
	uuid, uuidOk := args["uuid"].(string)

	if !metricOk || !subMetricOk || !periodOk || !aggregateTypeOk || !durationOk || !uuidOk {
		return nil, nil
	}

	var telegraf model.Telegraf
	var series []model.Series
	var s model.Series

	queryResult, err := influxdb.Influx.ReadMetric(metric, subMetric, period, aggregateType, duration, uuid)
	if err != nil {
		return nil, nil
	}

	logger.Logger.Println("queryResult : " + fmt.Sprintf("%v", queryResult))

	telegraf.UUID = fmt.Sprintf("%v", queryResult.(models.Row).Tags["host"])

	dataLength := len(queryResult.(models.Row).Values)

	if dataLength < 10 {
		for j := 0; j < 10-dataLength; j++ {
			s.Time = j
			s.Value = 0
			series = append(series, s)
		}
	}

	for i := 0; i < dataLength; i++ {
		s.Time = 10 - dataLength + i

		switch metric {
		case "cpu":
			valueStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i][1])
			valueFloat, _ := strconv.ParseFloat(valueStr, 64)
			s.Value = int(valueFloat * 100)
			break
		case "mem":
			valueStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i][1])
			valueFloat, _ := strconv.ParseFloat(valueStr, 64)
			s.Value = int(valueFloat * 1)
			break
		case "net":
		case "disk":
			valueStr := fmt.Sprintf("%v", queryResult.(models.Row).Values[i][1])
			valueFloat, _ := strconv.ParseFloat(valueStr, 64)
			s.Value = int(valueFloat * 1)
			break

		}

		series = append(series, s)
	}
	telegraf.Series = series
	telegraf.Metric = metric
	telegraf.SubMetric = subMetric

	return telegraf, nil
}
