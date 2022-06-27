# coding=utf-8
from typing import List, Any

import numpy as np
import math

from ekuiper.function import Function
from ekuiper.runtime.context import Context
from sklearn.cluster import KMeans


def clean_drop(data, compression_ratio):
    """
    数据降采样
    :param data: 用户输入的数据
    :param compression_ratio: 压缩比例
    :return: 将采样后的数据
    """
    if int(compression_ratio) <= 0:
        raise Exception("降采样压缩比例必须为大于0的整数")
    after = []
    i = 0
    while i < len(data):
        after.append(data[i])
        i = i + int(compression_ratio)
    return after


def clean_normal(array, regression_type):
    """
    归一化函数
    :param data:用户输入数据
    :param regression_type: 回归类型0: 最大最小值标准化, 1: Z_Score标准化, 2:sigmoid
    :return: 数组格式输出数据
    """
    dict_type = {0: "MaxMinNormalization", 1: "Z_ScoreNormalization", 2: "sigmoid"}
    regression_type = int(regression_type)
    after = []
    if regression_type == 0:
        max = np.max(array)
        min = np.min(array)
        for i in array:
            x = (i - min) / (max - min)
            after.append(x)
    if regression_type == 1:
        mu = np.average(array)
        sigma = np.std(array)
        for i in array:
            x = (i - mu) / sigma
            after.append(x)
    if regression_type == 2:
        k = np.ceil(np.log10(np.max(abs(np.array(array)))))
        after_array = np.array(array) / (10 ** k)
        after = list(after_array)

    return after


def clean_sort(array, range, sort_type):
    """
    排序函数
    :param array:用户输入数据
    :param range: 取之范围，list格式，[min,max]
    :param sort_type: 排序类型，0--正序，1--负序
    :return: 数组格式输出数据
    """
    # 0:正序 1：负序
    dict_type = {0: False, 1: True}
    sort_type = int(sort_type)
    after = []
    # 取范围
    for value in array:
        if range[0] <= value <= range[1]:
            after.append(value)
    # 排序
    after.sort(reverse=dict_type[sort_type])
    return after


def dimension_reduction(array, compression_ratio, method):
    """
    数据降维处理
    :param array: 用户输入数据
    :param compression_ratio: 压缩比例
    :param dimension_reduction: 降维方式, 0: 均值, 1: 总和, 2: 最大, 3: 最小， 4：K-means
    :return: 数组格式输出数据
    """
    dimension_reduction = int(method)
    compression_ratio = int(compression_ratio)
    # 处理数据

    length = len(array)
    n = int(length / compression_ratio)
    # 拆成n段
    after = []
    if method == 4:
        line = np.array(array).reshape(-1, 1)  # 变成列
        km = KMeans(n_clusters=n).fit(line)  # 聚类中心n
        list_all = km.cluster_centers_.reshape(1, -1)
        after = list(list_all[0])  # 取出list
        return after
    else:
        for value in range(n):
            list_array = array[math.floor(value / n * length):math.floor((value + 1) / n * length)]
            list_part = list(list_array)
            # 对每一段分析
            # {0: '均值', 1: '总和', 2: '最大', 3: '最小'}
            if dimension_reduction == 0:
                after.append(np.mean(list_part))
            if dimension_reduction == 1:
                after.append(np.sum(list_part))
            if dimension_reduction == 2:
                after.append(np.max(list_part))
            if dimension_reduction == 3:
                after.append(np.min(list_part))
        return after

class CleanDrop(Function):

    def __init__(self):
        pass

    def validate(self, args: List[Any]):
        if len(args) != 2:
            return "require 2 parameters"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        return clean_drop(args[0], args[1])

    def is_aggregate(self):
        return False


class CleanNormal(Function):

    def __init__(self):
        pass

    def validate(self, args: List[Any]):
        if len(args) != 2:
            return "require 2 parameters"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        return clean_normal(args[0], args[1])

    def is_aggregate(self):
        return False


class CleanSort(Function):

    def __init__(self):
        pass

    def validate(self, args: List[Any]):
        if len(args) != 4:
            return "require 4 parameters"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        return clean_sort(args[0], [args[1], args[2]], args[3])

    def is_aggregate(self):
        return False



class DimensionReduction(Function):

    def __init__(self):
        pass

    def validate(self, args: List[Any]):
        if len(args) != 3:
            return "require 3 parameters"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        return dimension_reduction(args[0], args[1], args[2])

    def is_aggregate(self):
        return False


cleanDropIns = CleanDrop()
cleanNormalIns = CleanNormal()
cleanSortIns = CleanSort()
dimensionReductionIns = DimensionReduction()