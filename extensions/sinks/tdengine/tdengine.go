// Copyright 2021-2023 EMQ Technologies Co., Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tdengine

import (
	"database/sql"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"

	_ "github.com/taosdata/driver-go/v3/taosRestful"
	_ "github.com/taosdata/driver-go/v3/taosWS"

	"github.com/lf-edge/ekuiper/pkg/api"

	"github.com/lf-edge/ekuiper/pkg/cast"
)

type TaosConfig struct {
	ProvideTs    bool     `json:"provideTs"`
	Port         int      `json:"port"`
	Host         string   `json:"host"`
	User         string   `json:"user"`
	Password     string   `json:"password"`
	Database     string   `json:"database"`
	Table        string   `json:"table"`
	TsFieldName  string   `json:"tsFieldName"`
	Fields       []string `json:"fields"`
	STable       string   `json:"sTable"`
	TagFields    []string `json:"tagFields"`
	DataTemplate string   `json:"dataTemplate"`
	DataField    string   `json:"dataField"`
}

type tdengineSink3 struct {
	cfg *TaosConfig
	cli *sql.DB
}

func (t *tdengineSink3) Open(ctx api.StreamContext) error {
	url := fmt.Sprintf(`%s:%s@http(%s)/%s`, t.cfg.User, t.cfg.Password, net.JoinHostPort(t.cfg.Host, strconv.Itoa(t.cfg.Port)), t.cfg.Database)
	taosCli, err := sql.Open("taosRestful", url)
	if err != nil {
		return err
	}
	t.cli = taosCli
	if err := t.cli.Ping(); err != nil {
		ctx.GetLogger().Errorf("tdengine3 sink connection failed, err:%v", err.Error())
		return err
	}
	ctx.GetLogger().Infof("tdengine3 sink connection success")
	return nil
}

func (t *tdengineSink3) Configure(props map[string]interface{}) error {
	t.cfg = &TaosConfig{
		Host:     "localhost",
		Port:     6041,
		User:     "root",
		Password: "taosdata",
	}
	err := cast.MapToStruct(props, t.cfg)
	if err != nil {
		return err
	}
	if t.cfg.Database == "" {
		return fmt.Errorf("property database is required")
	}
	if t.cfg.Table == "" {
		return fmt.Errorf("property table is required")
	}
	if t.cfg.TsFieldName == "" {
		return fmt.Errorf("property TsFieldName is required")
	}
	if t.cfg.STable != "" && len(t.cfg.TagFields) == 0 {
		return fmt.Errorf("property tagFields is required when sTable is set")
	}
	return nil
}

func (t *tdengineSink3) Collect(ctx api.StreamContext, data interface{}) error {
	switch v := data.(type) {
	case map[string]any:
		return t.collect(ctx, v)
	case []map[string]any:
		for _, m := range v {
			if err := t.collect(ctx, m); err != nil {
				return err
			}
		}
	case []interface{}:
		for _, d := range v {
			m, ok := d.(map[string]any)
			if !ok {
				return fmt.Errorf("unsupported type: %T", data)
			}
			if err := t.collect(ctx, m); err != nil {
				return err
			}
		}
	default: // never happen
		return fmt.Errorf("unsupported type: %T", data)
	}
	return nil
}

func (t *tdengineSink3) Close(ctx api.StreamContext) error {
	ctx.GetLogger().Infof("tdengine3 sink close")
	t.cli.Close()
	return nil
}

func (t *tdengineSink3) collect(ctx api.StreamContext, item map[string]any) error {
	sqlStr, sqlE := t.cfg.buildSql(ctx, item)
	if sqlE != nil {
		return fmt.Errorf("failed to build sql to tdengine3: %s", sqlE)
	}
	ctx.GetLogger().Infof("tdengine3 sink collect sql: %s", sqlStr)
	_, e := t.cli.Exec(sqlStr)
	if e != nil {
		ctx.GetLogger().Errorf("failed to exec sql to tdengine3: %s", e)
		return fmt.Errorf("failed to exec sql to tdengine3: %s", e)
	}
	return nil
}

func (cfg *TaosConfig) buildSql(ctx api.StreamContext, mapData map[string]any) (string, error) {
	var (
		table, sTable    string
		keys, vals, tags []string
		err              error
	)
	if 0 == len(mapData) {
		return "", fmt.Errorf("data is empty")
	}
	table, err = ctx.ParseTemplate(cfg.Table, mapData)
	if err != nil {
		ctx.GetLogger().Errorf("parse template for table %s error: %v", cfg.Table, err)
		return "", err
	}
	sTable, err = ctx.ParseTemplate(cfg.STable, mapData)
	if err != nil {
		ctx.GetLogger().Errorf("parse template for sTable %s error: %v", cfg.STable, err)
		return "", err
	}

	if cfg.ProvideTs {
		if v, ok := mapData[cfg.TsFieldName]; !ok {
			return "", fmt.Errorf("timestamp field not found : %s", cfg.TsFieldName)
		} else {
			keys = append(keys, cfg.TsFieldName)
			vals = append(vals, fmt.Sprintf(`%v`, v))
		}
	} else {
		vals = append(vals, "now")
		keys = append(keys, cfg.TsFieldName)
	}

	if len(cfg.TagFields) > 0 {
		for _, v := range cfg.TagFields {
			switch mapData[v].(type) {
			case string:
				tags = append(tags, fmt.Sprintf(`"%s"`, mapData[v]))
			default:
				tags = append(tags, fmt.Sprintf(`%v`, mapData[v]))
			}
		}
	}

	if len(cfg.Fields) != 0 {
		for _, k := range cfg.Fields {
			if k == cfg.TsFieldName {
				continue
			}
			if contains(cfg.TagFields, k) {
				continue
			}
			if v, ok := mapData[k]; ok {
				keys = append(keys, k)
				if reflect.String == reflect.TypeOf(v).Kind() {
					vals = append(vals, fmt.Sprintf(`"%v"`, v))
				} else {
					vals = append(vals, fmt.Sprintf(`%v`, v))
				}
			} else {
				return "", fmt.Errorf("field not found : %s", k)
			}
		}
	} else {
		for k, v := range mapData {
			if k == cfg.TsFieldName {
				continue
			}
			if contains(cfg.TagFields, k) {
				continue
			}
			keys = append(keys, k)
			if reflect.String == reflect.TypeOf(v).Kind() {
				vals = append(vals, fmt.Sprintf(`"%v"`, v))
			} else {
				vals = append(vals, fmt.Sprintf(`%v`, v))
			}
		}
	}

	sqlStr := fmt.Sprintf("INSERT INTO %s (%s)", table, strings.Join(keys, ","))
	if sTable != "" {
		sqlStr += " USING " + sTable
	}
	if len(tags) != 0 {
		sqlStr += " TAGS(" + strings.Join(tags, ",") + ")"
	}
	sqlStr += " values (" + strings.Join(vals, ",") + ")"
	return sqlStr, nil
}

func GetSink() api.Sink {
	return &tdengineSink3{}
}

func contains(slice []string, target string) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}
