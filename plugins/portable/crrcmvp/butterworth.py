# coding=utf-8
from typing import List, Any

from ekuiper import Function, Context
from scipy import signal


class Butterworth(Function):

    def __init__(self):
        pass

    def validate(self, args: List[Any]):
        if len(args) != 5:
            return "require 5 parameters"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        # todo: type validation
        return apply_butter_filter(args[0], args[1], args[2], args[3], args[4]).tolist()

    def is_aggregate(self):
        return False


def apply_butter_filter(data, sample_rate, filter_band, type=0, order=5):
    """
    :param data:输入信号数据，list
    :param sample_rate:采样率
    :param filter_band:滤波频率，若type=0或者1，此参数为int或者float；若type=2或者3，此参数为list,list里的值为下截止频率和上截止频率[low_band, high_band]
    :param type:0--低通滤波， 1--高通滤波， 3--带通滤波， 4--带阻滤波
    :param order:滤波阶数，默认5
    :return:返回滤波后的数据
    """
    type_meaning = ['lowpass', 'highpass', 'bandpass', 'bandstop']
    if type == 0 or type == 1:
        # filter_band_para = filter_band[0]
        Wn = 2.0 * filter_band / sample_rate
        b, a = signal.butter(order, Wn, type_meaning[type], analog=False, output='ba')
        return signal.filtfilt(b, a, data, axis=-1, padtype='odd', padlen=None)
    elif type == 2 or type == 3:
        if len(filter_band) != 2:
            print("带通或者带阻需要两个参数")
            return None
        wn1 = 2.0 * filter_band[0] / sample_rate
        wn2 = 2.0 * filter_band[1] / sample_rate
        if wn1 < wn2:
            b, a = signal.butter(order, [wn1, wn2], type_meaning[type], analog=False, output='ba')
        else:
            b, a = signal.butter(order, [wn2, wn1], type_meaning[type], analog=False, output='ba')
        return signal.filtfilt(b, a, data, axis=-1, padtype='odd', padlen=None)
    else:
        print("滤波类型错误")
        return None


butterworthIns = Butterworth()
