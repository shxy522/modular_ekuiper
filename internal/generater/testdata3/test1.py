#  导入库
import cv2
import fastdeploy as fd
import time
import numpy as np

BASEPATH = "{{.BASEPATH}}"


def init_model():
    # 配置runtime，加载模型
    runtime_option = fd.RuntimeOption()

    # 配置后端，使用GPU和TensorRT
    # runtime_option.use_gpu()
    runtime_option.use_ort_backend()

    runtime_option.enable_pinned_memory()
    # 载入模型，分别为几个模型文件的路径
    model = fd.vision.detection.PicoDet(
        BASEPATH + "model.pdmodel",
        BASEPATH + "model.pdiparams",
        BASEPATH + "infer_cfg.yml", runtime_option=runtime_option)
    ## 输入随机数预热网络
    for i in range(10):
        result = model.predict(np.random.rand(640, 640, 3).astype("uint8"))
    return model


def picodet(model, path=BASEPATH + "1.jpg"):
    # 预测图片检测结果
    ## 读取图片
    im = cv2.imread(path)

    kaishi = time.time()
    result = model.predict(im)  # 预测
    print(result)  # 输出预测结果
    print('infer time:', '{:.3f}'.format(time.time() - kaishi), 'S')  # 输出预测所耗费时间

    # 预测结果可视化
    vis_im = fd.vision.vis_detection(im, result, score_threshold=0.5)  # 阈值为0.5

    cv2.imwrite(BASEPATH + "visualized_result.jpg", vis_im)
    print("Visualized result saved")
    outpath = BASEPATH + "visualized_result.jpg"

    return outpath


def paddle():
    model = init_model()
    picodet(path=BASEPATH + "1.jpg", model=model)
