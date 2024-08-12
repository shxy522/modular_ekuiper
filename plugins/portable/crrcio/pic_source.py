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

import base64
class pic(Source):
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
            with open(self.file, 'rb') as image_file:
                encoded_string = base64.b64encode(image_file.read()).decode()
        except Exception as e:
            print("stop reading for ex:", e)
            image_file.close()
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
        print(m)
        print("file source closed")
        image_file.close()

    def close(self, ctx: Context):
        print("closing")
        self.running = False
# 这里把base64写入文件测试用的
# import base64
# from PIL import Image
# from io import BytesIO
# def display_decoded_image(base64_string):
#     with open("2.txt", 'w') as file:
#         file.write(base64_string)
#     decoded_image = base64.b64decode(base64_string)
#     image = Image.open(BytesIO(decoded_image))
#     image.show()