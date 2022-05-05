import json
import os.path
import time
from typing import Any

from ekuiper import Sink, Context
import matplotlib.pyplot as plt


class PlotSink(Sink):

    def __init__(self):
        self.x = ''
        self.y = ''
        self.filedir = '.'
        self.title = ''

    def configure(self, conf: dict):
        if 'x' in conf:
            self.x = conf['x']
        if 'y' in conf:
            self.y = conf['y']
        else:
            raise Exception('configuration y is missing')
        if 'filedir' in conf:
            self.filedir = conf['filedir']
        if 'title' in conf:
            self.title = conf['title']

    def open(self, ctx: Context):
        print('open print sink: ', ctx)

    def collect(self, ctx: Context, data: Any):
        msg = json.loads(data)
        if self.x != '':
            plt.plot(msg[self.x], msg[self.y])
        else:
            plt.plot(msg[self.y])
        plt.title(self.title)
        plt.savefig(os.path.join(self.filedir, str(time.time())))

    def close(self, ctx: Context):
        print("closing plot sink")
