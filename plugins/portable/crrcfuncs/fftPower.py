# coding=utf-8
import numpy as np
from typing import List, Any

from ekuiper import Function, Context
from scipy.fftpack import fft


class FFTPower(Function):

    def __init__(self):
        pass

    def validate(self, args: List[Any]):
        if len(args) != 5:
            return "require 5 parameters"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        # todo: type validation
        return fftPSD(args[0], args[1], args[2], args[3], args[4])

    def is_aggregate(self):
        return False


def fftPSD(signal, sample_rate, unit=0, method=0, type=0):
    """
    计算傅立叶变换（幅值谱）
    :param signal: 用户输入的数据
    :param sample_rate: 采样频率
    :param unit: 单位，0--（W/Hz）；1--dB
    :param method: 方法，0--直接法；1--间接法
    :param type: 谱类型，0--单边谱，1--双边谱
    :return:返回1--频率list，返回2--功率谱list
    """
    N = len(signal)
    # sample_interval = 1.0 / sample_rate  # 采样间隔
    # signal_len = N * sample_interval  # 信号长度
    # t = np.arange(0, signal_len, sample_interval)
    if type == 0:
        freq = sample_rate * np.array(range(0, int(N / 2))) / float(N)
    else:
        freq = sample_rate * np.array(range(0, int(N))) / float(N)
    if method == 0:
        if type == 0:
            fft_data = fft(signal)[0:N // 2]
        else:
            fft_data = fft(signal)
        fy = np.abs(fft_data)
        ps = fy ** 2 / N
        ps[0] = 0.1
        if unit == 1:
            ps = 20 * np.log10(ps)
        return [freq.tolist(), ps.tolist()]
    if method == 1:
        cor_x = np.correlate(signal, signal, 'same')
        if type == 0:
            cor_X = fft(cor_x, N)[:N // 2]
        else:
            cor_X = fft(cor_x, N)
        ps_cor = np.abs(cor_X)
        ps_cor = ps_cor / np.max(ps_cor)
        ps_cor[0] = 0.1
        if unit == 1:
            ps_cor = 20 * np.log10(ps_cor)
        return [freq.tolist(), ps_cor.tolist()]


fftpowerIns = FFTPower()