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


class ws_zmq(Source):
    def __init__(self):
        self.running = True
        self.server = "tcp://192.168.31.183:5555"
        self.topic = "CRRCWS/e45f01bb48da/31/139"

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
        socket.setsockopt_string(zmq.SUBSCRIBE, self.topic)
        print("connected to ", self.server)

        count = 0
        while self.running:
            if(self.topic != ""):
                topic=socket.recv_string()
            data = socket.recv()
            org = self.decode_data(data)
            for item in org: 
                m = {
                    "count": count,
                    "data": item,
                }
                ctx.emit(m, None)
            count += 1
        print("closed")

    def close(self, ctx: Context):
        print("closing")
        self.running = False

    def decode_data(self, data):
        # Ensure the data length is correct
        if len(data) < 42:
            raise ValueError("Data length is too short to decode")

        # Initialize an index to keep track of the current position in the data
        index = 0

        # Decode the fixed header
        header = data[index:index+2]
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
        address = struct.unpack_from('B', data, index)[0]
        index += 1

        # Skip 7 bytes of padding
        index += 7

        # Decode cpuid (4 bytes, little-endian)
        cpuid = struct.unpack_from('<I', data, index)[0]
        index += 4

        # Decode actual_scan_rate (4 bytes, little-endian float)
        actual_scan_rate = struct.unpack_from('<f', data, index)[0]
        index += 4

        # Decode chmask (1 byte)
        chmask = struct.unpack_from('B', data, index)[0]
        index += 1

        # Decode timestamp (8 bytes, little-endian)
        timestamp = struct.unpack_from('<Q', data, index)[0]
        index += 8
        # timestamp = datetime.datetime.fromtimestamp(timestamp / 1000.0)

        # Decode bufsize (4 bytes, little-endian)
        bufsize = struct.unpack_from('<I', data, index)[0]
        index += 4
        value=[]
        base_data={
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
            "ch":0,
        }
        for i in range(bufsize):
            # Decode data (4 bytes, little-endian)
            value.append(struct.unpack_from('<f', data, index)[0])
            index += 4
        # chmask转换为二进制0和1的字符串
        chmask_str = '{:08b}'.format(chmask)
        num_of_ones = chmask_str.count('1')
        count_chmask=0
        data=[]
        for index,i in  enumerate(chmask_str):
            if i == '1':
                base_data['CHANNEL'] = index
                base_data['ch'] = index
                base_data['signal'] = value[count_chmask::num_of_ones]
                count_chmask+=1
                data.append(base_data)
        return data
zmq1 = ws_zmq()
zmq1.configure("", {
    "server": "tcp://192.168.31.183:5555",
    "address": 7
})
zmq1.open(Context())
