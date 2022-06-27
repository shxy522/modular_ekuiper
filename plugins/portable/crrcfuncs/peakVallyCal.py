# coding=utf-8
# from math import sin
# from matplotlib import pylab
from typing import List, Any

from ekuiper import Function, Context
from pylab import *

class PeakVallydet(Function):

    def __init__(self):
        pass

    def validate(self, args: List[Any]):
        if len(args) != 2:
            return "require 2 parameters"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        # todo: type validation
        return peakVallydet(args[0], args[1])

    def is_aggregate(self):
        return False


def peakVallydet(data, pthresh):
    """
    峰谷值检测
    :param data:用户输入数据
    :param thresh:阈值百分比，取值为0-1，float，表示以最值的百分比作为判断波峰波谷
    :return:返回每个波峰波谷的差值
    """
    maxValue = max(data)
    minValue = min(data)
    thresh = maxValue - pthresh * maxValue
    treshVally = minValue + pthresh * abs(minValue)
    maxthresh = []
    minthresh = []
    peaks = []
    valleys = []
    peak_valleys_value = []

    t = array(range(len(data)))
    v = [*zip(t, data)]
    for x, y in v:
        if y > thresh:
            maxthresh.append((x, y))
        elif y < treshVally:
            minthresh.append((x, y))

    for x, y in maxthresh:
        try:
            if (v[x - 1][1] < y) & (v[x + 1][1] < y):
                peaks.append(y)
        except Exception:
            pass

    for x, y in minthresh:
        try:
            if (v[x - 1][1] > y) & (v[x + 1][1] > y):
                valleys.append(y)
        except Exception:
            pass
    for p, k in zip(peaks, valleys):
        peak_valleys_value.append(p-k)
    return peak_valleys_value

# t = array(range(100))
# series = 6.3 * sin(t) + 4.7 * cos(2 * t) - 3.5 * sin(1.2 * t)
# # print(peakdet(series, 0.95))
# #
# #arr = zip(t, series)
# # thresh = 0.95
# #
# p = peakVallydet(series, 0.6)
# print(p)
# for p, k in zip(peaks, valleys):
#     print(p-k)

# scatter([x for x, y in peaks], [y for x, y in peaks], color = 'red')
# scatter([x for x, y in valleys], [y for x, y in valleys], color = 'blue')
# plot(t, 100 * [thresh], color='green', linestyle='--', dashes=(5, 3))
# plot(t, 100 * [-thresh], color='green', linestyle='--', dashes=(5, 3))
# plot(t, series, 'k')
# show()


peakVallydetIns = PeakVallydet()