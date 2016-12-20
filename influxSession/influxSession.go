package influxSession

import (
    "github.com/influxdata/influxdb/client/v2"
    log "github.com/Sirupsen/logrus"
    "fmt"
    "time"
    "encoding/json"
)

type InfluxSession interface {
    QueryDB(cmd string) (res []client.Result, err error)
    AddPoint(series string, value float64) error
    IncrementSeries(series string, value float64) (float64, error)
    GetLastValue(series string) (float64, error)
    InitDb() error
    Close()
}

type influxSession struct {
    database string
    client client.Client
}

func NewInfluxSession(host, database, username, password string) InfluxSession {

    url := "http://" + host + "/" + database
    log.WithField("url", url).
        WithField("username", username).
        WithField("password", "*hidden*").
        Debug("Creating InfluxDB client...")

    client, err := client.NewHTTPClient(client.HTTPConfig{
        Addr: url,
        Username: username,
        Password: password,
    })
    if err != nil {
        log.Fatal(err)
    }

    obj := influxSession{
        database: database,
        client: client,
    }

    return obj
}

func (i influxSession) Close() {
    log.Debug("InfluxSession: Closing")
    i.client.Close()
}

func (i influxSession) QueryDB(cmd string) (res []client.Result, err error) {

    log.WithField("cmd", cmd).Debug("InfluxSession: Querying")

    q := client.Query{
        Command:  cmd,
        Database: i.database,
    }
    if response, err := i.client.Query(q); err == nil {
        if response.Error() != nil {
            return res, response.Error()
        }
        res = response.Results
    } else {
        return res, err
    }
    return res, nil
}

func (i influxSession) AddPoint(series string, value float64) error {

    log.WithField("series", series).WithField("value", value).Info("InfluxSession: Adding new point")

    // Create a new point batch
    bp, err := client.NewBatchPoints(client.BatchPointsConfig{
        Database:  i.database,
        Precision: "ms",
    })
    if err != nil {
        return err
    }

    // Create a point and add to batch
    tags := map[string]string{}
    fields := map[string]interface{}{
        "value": value,
    }
    pt, err := client.NewPoint(series, tags, fields, time.Now())
    if err != nil {
        return err
    }

    bp.AddPoint(pt)

    // Write the batch
    err = i.client.Write(bp)
    if err != nil {
        return err
    }

    return nil
}

func (i influxSession) GetLastValue(series string) (float64, error) {

    log.WithField("series", series).Info("InfluxSession: Getting last value")

    query := fmt.Sprintf(`SELECT LAST(value) FROM "%s"`, series)

    res, err := i.QueryDB(query)
    if err != nil {
        return 0, err
    }

    if len(res) != 1 || len(res[0].Series) != 1 {
        return 0, nil
    }

    return res[0].Series[0].Values[0][1].(json.Number).Float64()
}

func (i influxSession) IncrementSeries(series string, value float64) (float64, error) {

    log.WithField("series", series).WithField("value", value).Info("InfluxSession: Increment")

    newValue, err := i.GetLastValue(series)
    if err != nil {
        return 0, err
    }
    newValue += value
    err = i.AddPoint(series, newValue)
    if err != nil {
        return 0, err
    }

    return newValue, nil
}

func (i influxSession) InitDb() error {
    log.WithField("database", i.database).Debug("InfluxSession: Creating database if not exists...")
    _, err := i.QueryDB("CREATE DATABASE " + i.database)
    if err != nil {
        return err
    }
    return nil
}