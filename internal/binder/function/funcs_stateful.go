// Copyright 2023 EMQ Technologies Co., Ltd.
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

package function

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/lf-edge/ekuiper/internal/compressor"
	"github.com/lf-edge/ekuiper/internal/ossuploader"
	"github.com/lf-edge/ekuiper/pkg/api"
	"github.com/lf-edge/ekuiper/pkg/ast"
	"github.com/lf-edge/ekuiper/pkg/cast"
	"github.com/lf-edge/ekuiper/pkg/message"
)

type compressFunc struct {
	compressType string
	compressor   message.Compressor
}

func (c *compressFunc) Validate(args []interface{}) error {
	var eargs []ast.Expr
	for _, arg := range args {
		if t, ok := arg.(ast.Expr); ok {
			eargs = append(eargs, t)
		} else {
			// should never happen
			return fmt.Errorf("receive invalid arg %v", arg)
		}
	}
	return ValidateTwoStrArg(nil, eargs)
}

func (c *compressFunc) Exec(args []interface{}, ctx api.FunctionContext) (interface{}, bool) {
	if args[0] == nil {
		return nil, true
	}
	arg0, err := cast.ToBytes(args[0], cast.CONVERT_SAMEKIND)
	if err != nil {
		return fmt.Errorf("require string or bytea parameter, but got %v", args[0]), false
	}
	arg1 := cast.ToStringAlways(args[1])
	if c.compressor != nil {
		if c.compressType != arg1 {
			return fmt.Errorf("compress type must be consistent, previous %s, now %s", c.compressType, arg1), false
		}
	} else {
		ctx.GetLogger().Infof("creating compressor %s", arg1)
		c.compressor, err = compressor.GetCompressor(arg1)
		if err != nil {
			return err, false
		}
		c.compressType = arg1
	}
	r, e := c.compressor.Compress(arg0)
	if e != nil {
		return e, false
	}
	return r, true
}

func (c *compressFunc) IsAggregate() bool {
	return false
}

type decompressFunc struct {
	compressType string
	decompressor message.Decompressor
}

func (d *decompressFunc) Validate(args []interface{}) error {
	var eargs []ast.Expr
	for _, arg := range args {
		if t, ok := arg.(ast.Expr); ok {
			eargs = append(eargs, t)
		} else {
			// should never happen
			return fmt.Errorf("receive invalid arg %v", arg)
		}
	}
	return ValidateTwoStrArg(nil, eargs)
}

func (d *decompressFunc) Exec(args []interface{}, ctx api.FunctionContext) (interface{}, bool) {
	if args[0] == nil {
		return nil, true
	}
	arg0, err := cast.ToBytes(args[0], cast.CONVERT_SAMEKIND)
	if err != nil {
		return fmt.Errorf("require string or bytea parameter, but got %v", args[0]), false
	}
	arg1 := cast.ToStringAlways(args[1])
	if d.decompressor != nil {
		if d.compressType != arg1 {
			return fmt.Errorf("decompress type must be consistent, previous %s, now %s", d.compressType, arg1), false
		}
	} else {
		ctx.GetLogger().Infof("creating decompressor %s", arg1)
		d.decompressor, err = compressor.GetDecompressor(arg1)
		if err != nil {
			return err, false
		}
		d.compressType = arg1
	}
	r, e := d.decompressor.Decompress(arg0)
	if e != nil {
		return e, false
	}
	return r, true
}

func (d *decompressFunc) IsAggregate() bool {
	return false
}

type ossUploaderFunc struct {
	uploader *ossuploader.AliyunOss
}

func (c *ossUploaderFunc) Validate(args []interface{}) error {
	if len(args) != 5 {
		return fmt.Errorf("The arguments should be at five .")
	}
	// for i, a := range args {
	// 	if ast.IsNumericArg(a) || ast.IsTimeArg(a) || ast.IsBooleanArg(a) {
	// 		return ProduceErrInfo(i, "string")
	// 	}
	// }
	return nil
}

func (c *ossUploaderFunc) Exec(args []interface{}, ctx api.FunctionContext) (interface{}, bool) {
	endpoint, ok := args[1].(string)
	if !ok {
		return fmt.Errorf("oss function is missing property endpoint"), false
	}
	accessKeyId, ok := args[2].(string)
	if !ok {
		return fmt.Errorf("oss function is missing property accessKeyId"), false
	}
	accessKeySecret, ok := args[3].(string)
	if !ok {
		return fmt.Errorf("oss function is missing property accessKeySecret"), false
	}
	bucketName, ok := args[4].(string)
	if !ok {
		return fmt.Errorf("oss function is missing property bucketName"), false
	}

	if c.uploader == nil {
		c.uploader = ossuploader.GetUploader(endpoint, accessKeyId, accessKeySecret, bucketName)
	}

	// 流程名模块名时间
	tU := time.Now().Unix()
	// flowName 流程名
	var flowName string = ctx.GetRuleId()
	// modularName 模块名
	var modularName string = ctx.GetOpId()
	// objectName 文件名为流程名模块名时间
	var objectName string = flowName + "_" + modularName + "_" + strconv.FormatInt(tU, 10)
	user := make(map[string]string)
	// 这里全部存储是因为aliyun不支持url下载
	user["endpoint"] = endpoint
	user["accessKeyId"] = accessKeyId
	user["accessKeySecret"] = accessKeySecret
	user["bucketName"] = bucketName
	user["objectName"] = objectName
	jsonStr, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("The arguments should be at five ."), false
	}
	errOS := c.uploader.UploadString(objectName, args[0].(string))
	if errOS != nil {
		return fmt.Errorf("Input oss data error"), false
	}
	return string(jsonStr), true
}

func (c *ossUploaderFunc) IsAggregate() bool {
	return false
}
