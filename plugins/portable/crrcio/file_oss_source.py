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


class file_oss_source(Source):
    def __init__(self):
        self.running = True
        self.file = ''
        self.length = 1024
        self.samplerate = 100
        self.hatId = 322
        self.address = 1
        self.channel = 1
        self.interval = 1
        self.channelDatas = []

    def configure(self, datasource: str, conf: dict):
        logging.info(
            "configuring with datasource {} and conf {}".format(datasource, conf))
        if 'file' in conf:
            self.file = conf["file"]
        if 'length' in conf:
            self.length = conf["length"]
        if 'samplerate' in conf:
            self.samplerate = conf["samplerate"]
        if 'hatId' in conf:
            self.hatId = conf['hatId']
        if 'address' in conf:
            self.address = conf['address']
        if 'channel' in conf:
            self.channel = conf['channel']
        if 'interval' in conf:
            self.interval = conf['interval']
        self.channelDatas.append({
            "taskid": 0,
            "HATID": 0,
            "ADDRESS": self.address,
            "CHANNEL": 0,
            "samplerate": 0,
            "timestamp": 0,
            "length": 0,
            "signal": [0, 0]
        })
    def open(self, ctx: Context):
        response = requests.get(self.file, stream=True)
        response.raise_for_status()  # 检查请求是否成功
        count = 0
        buffer=[]
        line_count = 0
        for line in response.iter_lines():
            if line:
                buffer.append(float(line))
                line_count += 1
                if line_count == self.length:
                    line_count = 0
                    for chanData in self.channelDatas:
                        chanData["signal"] = buffer
                    m = {
                        "count": count,
                        "data": self.channelDatas
                    }
                    ctx.emit(m, None)
                    print(m)
                    buffer = []
                    count = count + 1
                    # time.sleep(self.interval)
        if buffer:
            signal = []
            for line in buffer:
                signal.append(int(line))
            buffer = []
            line_count = 0
            for chanData in self.channelDatas:
                chanData["signal"] = signal
            m = {
                "count": count,
                "data": self.channelDatas
            }
            ctx.emit(m, None)
            print(m)
            count = count + 1
        print("file source closed")

    def close(self, ctx: Context):
        print("closing")
        self.running = False


# file1=File_minos_source()
# file1.configure("file",{"file":"http://localhost:8000/2.csv","length":1024,"samplerate":100,"hatId":322,"address":1,"channel":1,"interval":1})
# file1.open(Context())
