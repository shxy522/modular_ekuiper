import faulthandler
import ctypes
from math import *
import numpy as np

faulthandler.enable()

lib_fftPower = ctypes.cdll.LoadLibrary('plugins/portable/so/lib_fftPower.so')
# lib_fftPower.so调用说明
# 调用lib_fftPower.main_fftPower函数


def fftPower(fs, data):
    # Data initialization
    # 将python类型转为c语言里的类型
    c_fs = ctypes.c_double(fs)
    lenth = len(data)
    c_data = (ctypes.c_double * lenth)(*data)
    c_f = (ctypes.c_double * lenth)()
    c_power = (ctypes.c_double * lenth)()
    # Calling function library
    '''
    main_fftPower参数：
    para1：fs--采样率，输入，type--double
    para2：length--数据长度，输入，type--int
    para3：data--数据，输入，type--list数组
    para4：f--频率值，输出，type--list数组
    para5：power--功率值，输出，type--list数组
    '''
    # 调用示例如下
    lib_fftPower.main_fftPower(c_fs, ctypes.c_int(lenth), c_data, c_f, c_power)
    # Keep two decimal places
    for i in range(0, lenth):
        c_f[i] = round(c_f[i], 2)
        c_power[i] = round(c_power[i], 2)
    return [c_f.tolist(), c_power.tolist()]


# 按间距中的绿色按钮以运行脚本。
if __name__ == '__main__':
    fs = 100
    gap = 1 / fs
    t = 0
    py_data = []

    for i in range(0, 1000):
        t = gap * i
        py_data.append(1.3 * sin(2 * pi * 15 * t) + 1.7 * sin(2 * pi * 40 * (t - 2)) + 2.5 * np.random.randn(1)[0])

    f, power = fftPower(fs, py_data)

    for i in range(0, len(py_data)):
        print('f = ' + str(f[i]) + ' \tpower = ' + str(power[i]))
