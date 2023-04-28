# coding=utf-8
from typing import List, Any

from ekuiper.function import Function
from ekuiper.runtime.context import Context
from scipy import signal as sg


def chebyshev(signal, sample_rate, ftype, filter_band, order=5, rp=3):
    """
    切比雪夫滤波
    :param data:用户输入信号数据
    :param fs:采样率
    :param ftype:滤波类型，0--低通滤波， 1--高通滤波， 3--带通滤波， 4--带阻滤波
    :param freqs: 滤波频率，若type=0或者1，此参数为int或者float；若type=2或者3，此参数为list,list里的值为下截止频率和上截止频率[low_band, high_band]
    :param order:阶数
    :param rp:增益
    :return:滤波后的数据
    """
    type_meaning = ['lowpass', 'highpass', 'bandpass', 'bandstop']
    nyq = 0.5 * sample_rate
    if ftype == 0:
        cut = filter_band / nyq
        b, a = sg.cheby1(order, rp, cut, btype=type_meaning[ftype])
    elif ftype == 1:
        assert len(filter_band) == 1
        cut = filter_band / nyq
        b, a = sg.cheby1(order, rp, cut, btype=type_meaning[ftype])
    elif ftype == 2:
        assert len(filter_band) == 2
        lowcut, highcut = filter_band[0] / nyq, filter_band[1] / nyq
        b, a = sg.cheby1(order, rp, [lowcut, highcut], btype=type_meaning[ftype])
    elif ftype == 3:
        assert len(filter_band) == 3
        lowcut, highcut = filter_band[0] / nyq, filter_band[1] / nyq
        b, a = sg.cheby1(order, rp, [lowcut, highcut], btype=type_meaning[ftype])
    else:
        print("filter type wrong")
        return []
    filtered = sg.lfilter(b, a, signal)
    return filtered.tolist()

class Chebyshev(Function):

    def __init__(self):
        pass

    def validate(self, args: List[Any]):
        if len(args) != 6:
            return "require 6 parameters"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        # todo: type validation
        return chebyshev(args[0], args[1], args[2], args[3], args[4], args[5])

    def is_aggregate(self):
        return False


chebyshevIns = Chebyshev()
