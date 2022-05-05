#coding=utf-8
from typing import List, Any

import numpy as np
from ekuiper import Function, Context
from scipy.fftpack import fft

class FFTAmp(Function):

    def __init__(self):
        pass

    def validate(self, args: List[Any]):
        if len(args) != 3:
            return "require 3 parameters"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        # todo: type validation
        return fftTrans(args[0], args[1], args[2])

    def is_aggregate(self):
        return False

def fftTrans(signal, sample_rate, type=0):
    """
    计算傅立叶变换（幅值谱）
    :param signal: 用户输入的数据，
    :param sample_rate: 采样频率
    :param type: 谱类型，0--单边谱，1--双边谱
    :return:返回1--频率list，返回2--幅值list
    """
    N = len(signal)
    # sample_interval = 1.0 / sample_rate  # 采样间隔
    # signal_len = N * sample_interval  # 信号长度
    # t = np.arange(0, signal_len, sample_interval)
    fft_data = fft(signal)
    # 这里幅值要进行一定的处理，才能得到与真实的信号幅值相对应
    fft_amp_double = np.array(np.abs(fft_data) / N * 2)  # 用于计算双边谱
    #direct = fft_amp_double[0]
    fft_amp_double[0] = 0
    freq_double = sample_rate * np.array(range(0, N)) / float(N)
    if type == 1:
        return [freq_double.tolist(), fft_amp_double.tolist()]
    else:
        N_2 = int(N / 2)
        fft_amp_single = fft_amp_double[0:N_2]
        freq_single = sample_rate * np.array(range(0, int(N / 2))) / float(N)
        return [freq_single.tolist(), fft_amp_single.tolist()]

fftampIns = FFTAmp()