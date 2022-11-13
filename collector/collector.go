package collector

import (
	"fmt"
	"time"

	"github.com/influxdata/influxdb-client-go/api"

	"github.com/ahmedalkabir/kolibrio-infrastrucutre/config"

	influxdb2 "github.com/influxdata/influxdb-client-go"
)

var (
	client influxdb2.Client
	write  api.WriteApi
)

//Collector this package it will be used
// as medium between influxdb and web server
type Point struct {
	Table  string
	Tags   map[string]string
	Fields map[string]interface{}
}

func init() {
	uri := fmt.Sprintf("http://%s:%s", config.InfluxConf.Address, config.InfluxConf.Port)
	// without token for now
	client = influxdb2.NewClient(uri, "")

	write = client.WriteApi("asas", "kolibrio")
}

//NewPoint ...
func NewPoint(table string, tags map[string]string, fields map[string]interface{}) *Point {
	return &Point{table, tags, fields}
}

func (p *Point) Write() {
	pt := influxdb2.NewPoint(p.Table, p.Tags, p.Fields, time.Now())
	// write it immediately
	write.WritePoint(pt)
}
