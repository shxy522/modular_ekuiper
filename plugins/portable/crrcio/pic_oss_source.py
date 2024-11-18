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
import base64

class pic_oss_source(Source):
    def __init__(self):
        self.file = ''
        self.fileType="jpeg"
        self.running = True

    def configure(self, datasource: str, conf: dict):
        logging.info(
            "configuring with datasource {} and conf {}".format(datasource, conf))
        if 'file' in conf:
            self.file = conf["file"]
        else:
            logging.error("not found file in configuration".format(datasource, conf))
        if 'fileType' in conf:
            self.fileType = conf["fileType"]

    # noinspection PyTypeChecker
    def open(self, ctx: Context):
        print("opening file source", self.file)
        try:
            encoded_string = self.fetch_and_encode_image(self.file)
        except Exception as e:
            print("stop reading for ex:", e)
            return
        m={
            "taskid": 0,
            "HATID": 0,
            "CHANNEL": 0,
            "samplerate": 0,
            "timestamp": 0,
            "length": 0,
            "signal": encoded_string,
            "fileType":self.fileType,
        }
        ctx.emit(m, None)
        # print(m)
        # print("file source closed")

    def fetch_and_encode_image(self,url):
        # 从网络上获取图片
        response = requests.get(url)
        response.raise_for_status()  # 检查请求是否成功

        # 将图片数据转换为 base64 编码
        image_base64 = base64.b64encode(response.content).decode('utf-8')
        return image_base64
    def close(self, ctx: Context):
        print("closing")
        self.running = False

# pic1=pic_minos_source()
# pic1.configure("pic",{"file":"http://localhost:8000/1.jpg"})
# pic1.open(Context())
