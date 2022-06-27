from typing import List, Any

from ekuiper import Function, Context
from scipy import signal


class Datasmooth(Function):

    def __init__(self):
        pass

    def validate(self, args: List[Any]):
        if len(args) != 2:
            return "require 2 parameters"
        return ""

    def exec(self, args: List[Any], ctx: Context):
        # todo: type validation
        return savitzky_golay(args[0], args[1]).tolist()

    def is_aggregate(self):
        return False


def savitzky_golay(data, degree=2):
    """
    多项式拟合平滑法
    :param data:用户输入的数据
    :param degree:平滑程度，int，degree范围0-4，分别代表弱，较弱，适中，较强，较强，数值越大，平滑效果约强，默认为2
    :return:
    """
    # 弱，较弱，适中，较强，较强
    n = len(data)
    percent = [0.02, 0.05, 0.1, 0.2, 0.3]
    window_length = []
    polyorder = [5, 4, 3, 2, 1]
    for p in percent:
        wl = int(p * n)
        if wl % 2 == 0:
            wl += 1
        window_length.append(wl)
    if window_length[degree] < polyorder[degree]:
        window_length[degree] = polyorder[degree] + 1

    y_smooth = signal.savgol_filter(data, window_length[degree], polyorder[degree])
    return y_smooth


datasmoothIns = Datasmooth()
