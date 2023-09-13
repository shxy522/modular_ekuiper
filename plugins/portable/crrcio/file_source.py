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


class File(Source):
    def __init__(self):
        self.running = True
        self.file = 'etc/data.txt'
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

    # noinspection PyTypeChecker
    def open(self, ctx: Context):
        print("opening file source", self.file)
        f = open(self.file)
        print("start reading file for length", self.length)
        count = 0
        while self.running:
            try:
                i = 0
                signal = []
                while i < self.length:
                    line = f.readline()
                    if len(line) == 0:
                        raise Exception("end of file")
                    line = line.strip('\n')
                    try:
                        signal.append(float(line))
                    except Exception:
                        print("err reading", i, "line:", line)
                    i = i + 1
            except Exception as e:
                print("stop reading for ex:", e)
                break

            for chanData in self.channelDatas:
                chanData["signal"] = signal
            m = {
                "count": count,
                "data": self.channelDatas
            }
            ctx.emit(m, None)
            count = count + 1
            time.sleep(self.interval)
        print("file source closed")
        f.close()

    def close(self, ctx: Context):
        print("closing")
        self.running = False
