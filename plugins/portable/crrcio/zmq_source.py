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
        self.topic = ""

    def configure(self, datasource: str, conf: dict):
        logging.info("configuring with datasource {} and conf {}".format(datasource, conf))
        if "server" in conf:
            self.server = conf["server"]
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
            m = {
                "count": count,
                "taskid":org[1],
                "HATID": org[2],
                "ADDRESS": org[3],
                "CHANNEL": org[4],
                "samplerate": org[5],
                "timestap": org[6],
                "length": org[7],
                "signal": org[8:]
            }
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
        data_len = int((len(data) - 26) / 4)
        s = struct.Struct('>HIH2BIdI' + str(data_len) + 'f')
        org = s.unpack(data)
        return org
