
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

from ekuiper import plugin, PluginConfig


from file_source import File
from zmq_source import Zmq
from ws_source import ws_zmq
from pic_source import pic
from file_oss_source import file_oss_source
from pic_oss_source import pic_oss_source

if __name__ == '__main__':
    c = PluginConfig("crrcio", {"crrc_zmq": lambda: Zmq(), "crrc_zmq_ws": lambda:ws_zmq(),"crrc_file": lambda: File(),"crrc_pic":lambda :pic(),"crrc_oss_file":lambda:file_oss_source(),"crrc_oss_pic":lambda: pic_oss_source()}, {},{})
    plugin.start(c)
