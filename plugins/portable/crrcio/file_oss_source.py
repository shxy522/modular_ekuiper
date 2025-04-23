#  Copyright 2022 EMQ Technologies Co., Ltd.
#
#  Licensed under the Apache License, Version 2.0 (the "License");
#  you may not use this file except in compliance with the License.
#  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
#  limitations under the License.

import logging
import time

from ekuiper import Source, Context
import requests
import io
from copy import deepcopy
from typing import List, Optional

class file_oss_source(Source):
    def __init__(self):
        self.running = True
        self.file = ''
        self.length = 1024
        self.samplerate = 100
        self.hatId = 322
        self.address = 1
        self.delimiter = ','
        self.channel: Optional[List[int]] = None  # 用户可见的1-based通道号
        self.interval = 1
        self.channelDatas = []

    def configure(self, datasource: str, conf: dict):
        # logging.info(
        #     "configuring with datasource {} and conf {}".format(datasource, conf))
        if 'file' in conf:
            self.file = conf["file"]
        if 'length' in conf:
            self.length = conf["length"]
        self.delimiter = ','
        if 'samplerate' in conf:
            self.samplerate = conf["samplerate"]
        if 'hatId' in conf:
            self.hatId = conf['hatId']
        if 'address' in conf:
            self.address = conf['address']
        if 'channel' in conf:
            self.channel = conf['channel']
        else:
            self.channel = None

        if 'interval' in conf:
            self.interval = conf['interval']

        if self.channel is None:
            # 需要延迟到open阶段获取实际通道数
            pass
        else:
            for ch in self.channel:
                self.channelDatas.append({
                    "taskid": 0,
                    "HATID": 0,
                    "ADDRESS": 0,
                    "CHANNEL": [ch],  # 注意保持列表结构
                    "samplerate": self.samplerate,
                    "timestamp": 0,
                    "length": self.length,
                    "signal": []
                })


    def open(self, ctx: Context):
        """支持1-based通道号选择的多列数据流式读取（单次请求）"""
        logging.basicConfig(level=logging.INFO)
        try:
            with requests.get(self.file, stream=True, timeout=10) as response:
                response.raise_for_status()
                lines = response.iter_lines()

                # 获取有效第一行以确定列数
                first_line = None
                for line_bytes in lines:  # 逐行找第一个非空行
                    line = line_bytes.decode('utf-8-sig').strip()
                    if line:
                        first_line = line
                        break
                if not first_line:
                    raise ValueError("CSV文件无有效数据")

                total_csv_columns = len(first_line.split(self.delimiter))
                # 确定实际读取的列索引
                if self.channel is None:
                    self.channel = list(range(1, total_csv_columns + 1))
                    for ch in self.channel:
                        self.channelDatas.append({
                            "taskid": 0,
                            "HATID": 0,
                            "ADDRESS": 0,
                            "CHANNEL": [ch],
                            "samplerate": self.samplerate,
                            "timestamp": 0,
                            "length": self.length,
                            "signal": []
                        })
                    column_indices = list(range(total_csv_columns))
                else:
                    column_indices = [c - 1 for c in self.channel]
                    if max(column_indices) >= total_csv_columns:
                        raise ValueError(f"CSV列数不足，最大通道号应为{total_csv_columns}")

                num_channels = len(column_indices)
                buffers = [[] for _ in range(num_channels)]
                count = 0

                # 处理第一行数据
                parts = first_line.split(self.delimiter)
                if len(parts) >= total_csv_columns:
                    try:
                        row_data = [float(parts[i]) for i in column_indices]
                        for i in range(num_channels):
                            buffers[i].append(row_data[i])
                    except (IndexError, ValueError) as e:
                        logging.error(f"首行数据错误: {str(e)}")

                # 继续处理剩余行（包括后续实时流）
                for line_bytes in lines:  # 直接从迭代器继续读取剩余行
                    if not self.running:
                        break
                    line = line_bytes.decode().strip()
                    if not line:
                        continue
                    parts = line.split(self.delimiter)
                    if len(parts) < total_csv_columns:
                        logging.warning(f"列数不足，跳过行: {line}")
                        continue
                    try:
                        row_data = [float(parts[i]) for i in column_indices]
                    except Exception as e:
                        logging.error(f"数据转换失败: {str(e)} - 行内容: {line}")
                        continue

                    # 填充缓冲区
                    for i in range(num_channels):
                        buffers[i].append(row_data[i])

                    # 批量发送逻辑（保持原有）
                    if len(buffers[0]) >= self.length:
                        timestamp = int(time.time() * 1000)
                        current_batch = []
                        for i in range(num_channels):
                            # 复用channelDatas中的模板
                            data_template = deepcopy(self.channelDatas[i])
                            data_template.update({
                                "timestamp": timestamp,
                                "signal": buffers[i][:self.length].copy(),
                                "length": self.length
                            })
                            current_batch.append(data_template)
                        ctx.emit({"count": count, "data": current_batch}, None)
                        logging.info(
                            f"Sent batch {count} | Channels: {self.channel} | "
                            f"Samples: {self.length}"
                        )
                            # try:
                            #     # 先发送数据，再记录日志
                            #     ctx.emit({"count": count, "data": [data_template]}, None)
                            #     logging.info(
                            #         f"Sent remaining data - count={count}, channel={current_ch}, "
                            #         f"samples={self.length}"
                            #     )
                            # except Exception as e:
                            #     logging.error(
                            #         f"Failed to send remaining data - count={count}, channel={current_ch}: {str(e)}"
                            #     )
                            # 清空已发送数据
                        for buf in buffers:
                            del buf[:self.length]
                        count += 1
                        time.sleep(self.interval)  # 控制发送频率

                # # 处理剩余数据（保持原有）
                # if buffers and len(buffers[0]) > 0:
                #     if num_channels == 1:
                #         self.channelDatas[0]["signal"] = deepcopy(buffers[0])
                #     else:
                #         self.channelDatas[0]["signal"] = [deepcopy(buf) for buf in buffers]
                #     ctx.emit({"count": count, "data": deepcopy(self.channelDatas)}, None)
                # 处理剩余数据
                if any(len(buf) > 0 for buf in buffers):
                    timestamp = int(time.time() * 1000)
                    last_batch = []
                    for i in range(num_channels):
                        if buffers[i]:
                            data_template = deepcopy(self.channelDatas[i])
                            data_template.update({
                                        "timestamp": timestamp,
                                        "signal": buffers[i].copy(),
                                        "length": len(buffers[i])
                                    })
                            last_batch.append(data_template)

                    ctx.emit({"count": count, "data": last_batch}, None)
                    logging.info(
                        f"Sent final batch | Channels: {self.channel} | "
                        f"Samples: {sum(len(buf) for buf in buffers)}"
                    )
                # timestamp = int(time.time() * 1000)
                # for i in range(num_channels):
                #     if buffers[i]:
                #         current_ch = self.channelDatas[i]["CHANNEL"][0]
                #         data_template = deepcopy(self.channelDatas[i])
                #         data_template.update({
                #                 "timestamp": timestamp,
                #                 "signal": buffers[i].copy(),
                #                 "length": len(buffers[i])
                #             })
                #         try:
                #             # 先发送数据，再记录日志
                #             ctx.emit({"count": count, "data": [data_template]}, None)
                #             logging.info(
                #                 f"Sent remaining data - count={count}, channel={current_ch}, "
                #                 f"samples={len(buffers[i])}"
                #             )
                #         except Exception as e:
                #             logging.error(
                #                 f"Failed to send remaining data - count={count}, channel={current_ch}: {str(e)}"
                #             )
                #         # ctx.emit({"count": count, "data": [data_template]}, None)
                #         # count += 1
        except Exception as e:
            logging.error(f"运行时错误: {str(e)}")
            raise
        finally:
            logging.info("数据源已关闭")

    def close(self, ctx: Context):
        print("closing")
        self.running = False


# file1=File_minos_source()
# file1.configure("file",{"file":"http://localhost:8000/2.csv","length":1024,"samplerate":100,"hatId":322,"address":1,"channel":1,"interval":1})
# file1.open(Context())
