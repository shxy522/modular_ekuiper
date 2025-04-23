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
import struct

import zmq
from ekuiper import Source, Context


class Zmq(Source):
    def __init__(self):
        self.running = True
        self.server = "tcp://127.0.0.1:5000"
        self.topic = ""
        self.prefix_cpu = "10000000"
        self.channels = []
        self.channelDatas = []
        # self.context = None
        # self.socket = None

    def configure(self, datasource: str, conf: dict):
        # logging.info("configuring with datasource {} and conf {}".format(datasource, conf))
        if "server" in conf:
            self.server = conf["server"]
            self.topic = datasource
        if "topic" in conf:
            self.topic = conf["topic"]
        if "prefix_cpu" in conf:
            self.prefix_cpu = conf["prefix_cpu"]
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
        socket.setsockopt_string(zmq.SUBSCRIBE, self.topic)
        print("connected to ", self.server)
        # print("channelDatas", self.channelDatas)
        count = 0
        while self.running:
            if (self.topic != ""):
                topic = socket.recv_string()
            data = socket.recv()
            org = self.decode_data(data)
            # print("解包后的数据------------->", org)
            for d in org:
                for chanData in self.channelDatas:
                    if chanData["ADDRESS"] == d["ADDRESS"] and chanData["CHANNEL"] == d["CHANNEL"]:
                        chanData["taskid"] = d["taskid"]
                        chanData["HATID"] = d["HATID"]
                        chanData["ADDRESS"] = d["ADDRESS"]
                        chanData["CHANNEL"] = d["CHANNEL"]
                        chanData["samplerate"] = d["samplerate"]
                        chanData["timestamp"] = d["timestamp"]
                        chanData["length"] = d["length"]
                        chanData["signal"] = d["signal"]
            m = {
                "count": count,
                "data": self.channelDatas
            }
            print("count", count)
            count += 1
            ctx.emit(m, None)
        print("closed")

    def close(self, ctx: Context):
        self.running = False
        # if self.socket:
        #     self.socket.close()
        # if self.context:
        #     self.context.term()
        print("ZMQ连接已关闭")
    def decode_data(self, data):
        # Ensure the data length is correct
        if len(data) < 42:
            raise ValueError("Data length is too short to decode")

        # Initialize an index to keep track of the current position in the data
        index = 0

        # Decode the fixed header
        header = data[index:index + 2]
        index += 2

        if header != b'\xa5\x5a':
            raise ValueError("Invalid header")

        # Decode taskid (4 bytes, little-endian)
        taskid = struct.unpack_from('<I', data, index)[0]
        index += 4

        # Decode Msgid (2 bytes, little-endian)
        msgid = struct.unpack_from('<H', data, index)[0]
        index += 2

        # Decode hatid (2 bytes, little-endian)
        hatid = struct.unpack_from('<H', data, index)[0]
        index += 2

        # Decode address (1 byte)
        address = struct.unpack_from('<Q', data, index)[0]
        index += 8

        # Decode cpuid (4 bytes, little-endian)
        # Uint32解码
        cpuid = struct.unpack_from('<I', data, index)[0]
        cpuid = self.prefix_cpu + f"{cpuid:08X}"
        index += 4

        # Decode actual_scan_rate (4 bytes, little-endian float)
        actual_scan_rate = struct.unpack_from('<f', data, index)[0]
        index += 4

        # Decode chmask (1 byte)
        chmask = struct.unpack_from('<B', data, index)[0]
        index += 1

        # Decode timestamp (8 bytes, little-endian)
        timestamp = struct.unpack_from('<Q', data, index)[0]
        index += 8
        # timestamp = datetime.datetime.fromtimestamp(timestamp / 1000.0)

        # Decode bufsize (4 bytes, little-endian)
        bufsize = struct.unpack_from('<I', data, index)[0]
        index += 4
        value = []

        while index < len(data):
            unpack = struct.unpack_from('>f', data, index)[0]
            value.append(unpack)
            index += 4
        # 放入class加这个
        # if "CRRCWS" in self.topic:
        #     chmask &= 0b00000111
        # chmask &= 0b00000111
        chmask_str = '{:08b}'.format(chmask)
        num_of_ones = chmask_str.count('1')
        count_chmask = 0
        data = []
        # print("ch: ", chmask_str)
        for index, i in enumerate(chmask_str):
            if i == '1':
                base_data = {
                    "taskid": taskid,
                    'msgid': msgid,
                    "HATID": hatid,
                    "ADDRESS": address,
                    'CPUID': cpuid,
                    "CHANNEL": 0,
                    "samplerate": actual_scan_rate,
                    "timestamp": timestamp,
                    "length": bufsize,
                    "signal": [],
                    "ch": 0,
                    'addr': address,
                }
                base_data['CHANNEL'] = 7 - index
                base_data['ch'] = 7 - index
                base_data['signal'] = value[count_chmask::num_of_ones]
                count_chmask += 1
                data.append(base_data)
        return data

# ws=ws_source()
# ws.configure('111',{"server":"tcp://127.0.0.1:5000",'topic':"5555"})
# data=b'\xa5Zu\x00\x00\x00\x05\x00E\x01\x01\x00\x00\x00\x00\x00\x00\x00\xf5O\xba\x93\x00\xc0zD\x01\xab\xd4c_\x94\x01\x00\x00d\x00\x00\x00@\x14AI>w\x89\x18<\x8c\x8c\x99\xbf\x7f\x05\x8d?)\x8bf?\xd7\xf5\x96?\x0e\xe6\xeb?0f`?*\xa2\x99>k\x12:?\x15"Z\xbf\x17\x19\xb0\xbf6\x1a\xf9?]E\x80@\x16\x1f\xeb?\xafsD\xbf(\xdc\xad\xbf\xf7\xb6\xb6\xbfL\x12\xe1?\xabRJ?\xec\xea<?v\xd2\xc8?\x1a\xe6!\xbf\x0c\xf9\x1b?\x06}\x15?\xe3\xa5\x07@\x04\x98\xc2@\x1fG6>\xc8\x1c\x9d\xbf\x89S\x96?N\xa0;=\xa2F\x01>\xc5\xee7?@J\xbc?\xe6#3\xbf;?5?\xcc\x1eC?fv\xc4\xbe\xe4\t\xb6=\xc2\xaeC\xbf\xc8\xe5\x03?R\r\xb8\xbd\xfaK\xbc\xbe\xc1s\r\xbe!Kb?\xa4wP\xbe\x8di-?\xdf V>\xe85\x1a?}\x0e7>P\xe5g?\xe67%?\x08\x83\x99?By"@#\xdf\xd8?\x1d\x14\x88?r&5\xbf\xe8\x99\xc9>W\xc0a\xbe\x80\xa2\x89?N\xf0\x01?U+p\xbe\x1e-\xab?i\xbc_?F6e\xbf\xb0\xd2\xa6\xbfY\xa0\xf2?\xac\xb9C>\xc4\xfe\xe6\xbf\xc2\xa9\x94?t|~?Z\'\xc9?\x9a\x06\xf6?u\x93\xb2\xbf\x0e\xff\xa0>E\x0e\x14>\xd9\x8f\xd5\xbf\x89\xdf/\xbf\x8a\xba\x8f\xbfu\xacg\xbe\xba\xe7\xd8?\xf13\x18=\x8eT\x9e@\x18l<>\x83?\x11\xbe\x00\xe3!?\xc4S\xf8?\xb1Q\xe5?\xff\x9c\x89?\x81U\x0e=\xf4\x89\xbc?\xa11\xb6\xbe\xdf]"?\x82D^?\xad\x1c\xfa\xbe7\x1bg?\xeda\xe4?\xcc\xbd\xce\xbe\x1cO\n\xbf\x97YP'
# print(ws.decode_data(data))
# ws

