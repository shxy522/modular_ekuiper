# -*- coding: utf-8 -*-

import math
from typing import List, Any

import numpy as np
from ekuiper.function import Function
from ekuiper.runtime.context import Context


def split_list_by_n(list_collection, n):
    """
    将集合均分，每份n个元素
    :param list_collection:
    :param n:
    :return:返回的结果为评分后的每份可迭代对象
    """
    for i in range(0, len(list_collection), n):
        yield list_collection[i: i + n]


def get_average(signal, calCount):
    """
    平均值
    """
    signal_split = split_list_by_n(signal, calCount)
    return [np.mean(a) for a in signal_split]


def get_max(signal, calCount):
    """
    最大值
    """
    signal_split = split_list_by_n(signal, calCount)
    return [np.max(a) for a in signal_split]


def get_min(signal, calCount):
    """
    最小值
    """
    signal_split = split_list_by_n(signal, calCount)
    return [np.min(a) for a in signal_split]


def get_amp(signal, calCount):
    """
    幅值
    """
    signal_split = split_list_by_n(signal, calCount)
    return [np.max(a) - np.min(a) for a in signal_split]


def get_variance(signal):
    """
    方差 反映一个数据集的离散程度
    """
    average = get_average(signal)
    return sum([(x - average) ** 2 for x in signal]) / len(signal)


def get_standard_deviation(signal):
    """
    标准差 == 均方差 反映一个数据集的离散程度
    """
    variance = get_variance(signal)
    return math.sqrt(variance)


def get_rms(signal, calCount):
    """
    均方根值 反映的是有效值而不是平均值
    """
    signal_split = split_list_by_n(signal, calCount)
    list_split = list(signal_split)
    r = []
    for s in list_split:
        rms = math.sqrt(sum([x ** 2 for x in s]) / len(s))
        r.append(rms)
    return r


class ListAverage(Function):

    def __init__(self):
        pass

    def validate(self, args: List[Any]):
        if len(args) != 2:
            return "require 2 parameters"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        # todo: type validation
        return get_average(args[0], args[1])

    def is_aggregate(self):
        return False


class ListMax(Function):

    def __init__(self):
        pass

    def validate(self, args: List[Any]):
        if len(args) != 2:
            return "require 2 parameters"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        return get_max(args[0], args[1])

    def is_aggregate(self):
        return False

class ListMin(Function):

    def __init__(self):
        pass

    def validate(self, args: List[Any]):
        if len(args) != 2:
            return "require 2 parameters"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        return get_min(args[0], args[1])

    def is_aggregate(self):
        return False


class ListAmp(Function):

    def __init__(self):
        pass

    def validate(self, args: List[Any]):
        if len(args) != 2:
            return "require 2 parameters"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        return get_amp(args[0], args[1])

    def is_aggregate(self):
        return False

class ListRMS(Function):

    def __init__(self):
        pass

    def validate(self, args: List[Any]):
        if len(args) != 2:
            return "require 2 parameters"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        return get_rms(args[0], args[1])

    def is_aggregate(self):
        return False

class ListVariance(Function):

    def __init__(self):
        pass

    def validate(self, args: List[Any]):
        if len(args) != 1:
            return "require 1 parameter"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        return get_variance(args[0])

    def is_aggregate(self):
        return False


class ListStandardDeviation(Function):

    def __init__(self):
        pass

    def validate(self, args: List[Any]):
        if len(args) != 1:
            return "require 1 parameter"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        return get_standard_deviation(args[0])

    def is_aggregate(self):
        return False


listAverageIns = ListAverage()
listMaxIns = ListMax()
listMinIns = ListMin()
listAmpIns = ListAmp()
listRmsIns = ListRMS()
listVarianceIns = ListVariance()
listStandardDeviationIns = ListStandardDeviation()