
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

from butterworth import butterworthIns
from chebyshev import chebyshevIns
from dataPreprocess import cleanDropIns, cleanNormalIns, cleanSortIns, dimensionReductionIns
from dataSmooth import datasmoothIns
from fftAmp import fftampIns
from fftPower import fftpowerIns
from file_source import File
from peakVallyCal import peakVallydetIns
from plot import PlotSink
from statistic import listAverageIns, listMaxIns, listMinIns, listAmpIns, listRmsIns, \
    listVarianceIns, listStandardDeviationIns
from zmq_source import Zmq

if __name__ == '__main__':
    c = PluginConfig("crrcmvp", {"crrc_zmq": lambda: Zmq(), "crrc_file": lambda: File()}, {"plot": lambda: PlotSink()},
                     {"butterworth": lambda: butterworthIns, "chebyshev": lambda: chebyshevIns, "datasmooth": lambda: datasmoothIns, "fftpower": lambda: fftpowerIns, "peakvallydet": lambda: peakVallydetIns, "fftamp": lambda: fftampIns,"listaverage": lambda: listAverageIns, "listmin": lambda: listMinIns, "listmax": lambda: listMaxIns, "listamp": lambda: listAmpIns, "listrms": lambda: listRmsIns, "listvariance": lambda: listVarianceIns, "liststandarddeviation": lambda: listStandardDeviationIns, "cleandrop": lambda: cleanDropIns,"cleannormal": lambda: cleanNormalIns,"cleansort": lambda: cleanSortIns,"dimensionreduction": lambda: dimensionReductionIns,})
    plugin.start(c)
