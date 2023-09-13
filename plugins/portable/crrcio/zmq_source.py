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
import struct

import zmq
from ekuiper import Source, Context


class Zmq(Source):
    def __init__(self):
        self.running = True
        self.s = None
        self.server = "tcp://127.0.0.1:5000"
        self.address = 0
        self.channels = []
        self.channelDatas = []
        self.topic = ""

    def configure(self, datasource: str, conf: dict):
        logging.info("configuring with datasource {} and conf {}".format(datasource, conf))
        if "server" in conf:
            self.server = conf["server"]
        if "channels" not in conf:
            logging.error("not found channels in configuration".format(datasource, conf))
        if "address" not in conf:
            logging.error("not found address in configuration".format(datasource, conf))
        if "address" in conf:
            self.address = conf["address"]
        if "channels" in conf:
            self.channels = conf["channels"]
        for chanId in self.channels:
            self.channelDatas.append({
                "taskid": 0,
                "HATID": 0,
                "ADDRESS": self.address,
                "CHANNEL": chanId,
                "samplerate": 0,
                "timestamp": 0,
                "length": 0,
                "signal": [0, 0]
            })

        self.topic = datasource

    # noinspection PyTypeChecker
    def open(self, ctx: Context):
        print("opening")
        context = zmq.Context()
        socket = context.socket(zmq.SUB)
        socket.connect(self.server)
        socket.setsockopt_string(zmq.SUBSCRIBE, '')
        print("connected to ", self.server)

        count = 0
        while self.running:
            data = socket.recv()
            org = self.data_unpack(data)
            # print("=======================================================================")
            # print("package count:{}".format(count))
            # print("HATID={}, ADDRESS={}, CHANNEL={}, timestap={}, length={}".
            #       format(org[0], org[1], org[2], org[3], org[4]))
            # print("data={}".format(org[5:]))
            count += 1

            # chan = org[4]
            # for chanData in self.channelDatas:
            #     if chanData["CHANNEL"] == chan:
            #         chanData["taskid"] = org[1]
            #         chanData["HATID"] = org[2]
            #         chanData["ADDRESS"] = org[3]
            #         chanData["CHANNEL"] = org[4]
            #         chanData["samplerate"] = org[5]
            #         chanData["timestamp"] = org[6]
            #         chanData["length"] = org[7]
            #         chanData["signal"] = org[8:]

            chan = org["CHANNEL"]
            addr = org["ADDRESS"]
            for chanData in self.channelDatas:
                if chanData["ADDRESS"] == addr and chanData["CHANNEL"] == chan:
                    chanData["taskid"] = org["taskid"]
                    chanData["HATID"] = org["HATID"]
                    chanData["ADDRESS"] = org["ADDRESS"]
                    chanData["CHANNEL"] = org["CHANNEL"]
                    chanData["samplerate"] = org["samplerate"]
                    chanData["timestamp"] = org["timestamp"]
                    chanData["length"] = org["length"]
                    chanData["signal"] = org["signal"]
            m = {
                "count": count,
                "data": self.channelDatas
            }

            # print("data={}".format(m))

            ctx.emit(m, None)
        print("closed")

    def close(self, ctx: Context):
        print("closing")
        self.running = False

    def data_unpack(self, data):
        """
        将二进制数据按照相应的格式解包
        :param data: 二进制格式数据
        :return: 解包后的数据
        """
        # data_len = int((len(data) - 26) / 4)
        # s = struct.Struct('>HIH2BIdI' + str(data_len) + 'f')
        # org = s.unpack(data)

        data_len = int((len(data) - 26) / 4)
        s1 = struct.Struct('>HIH2BIdI')
        org = s1.unpack(data[:26])
        m = {
            "taskid": org[1],
            "HATID": org[2],
            "ADDRESS": org[3],
            "CHANNEL": org[4],
            "samplerate": org[5],
            # "timestamp": org[3],
            "length": org[7],
        }
        t = struct.unpack("d", data[14:22])
        m["timestamp"] = t[0]
        dataStruct = struct.Struct(str(data_len) + 'f')
        dataUnpack = dataStruct.unpack(data[26:])
        m["signal"] = dataUnpack

        return m
