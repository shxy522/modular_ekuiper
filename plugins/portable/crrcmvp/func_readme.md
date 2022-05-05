## 函数说明

### 函数分类

#### 基本函数（ekuiper自带函数）

#### 信号处理工具箱

* 信号预处理

  * 降采样

    dataPreprocess.clean_drop

  * 归一化

    dataPreprocess.clean_normal

  * 排序

    dataPreprocess.clean_sort

  * 降维

    dataPreprocess.dimension_reduction

  * 平滑处理

    * 多项式平滑，dataSmooth.savitzky_golay

  * 特征统计

    * 峰谷值统计，peakVallyCal.peakVallydet

* 滤波

  * IIR滤波器
    * 巴特沃斯滤波（butterworth.apply_butter_filter)
    * 切比雪夫滤波  (chebyshev.apply_chebyshev_filter)
  * FIR滤波器

* 频谱分析

  * 幅值谱

    fftAmp.fftTrans

  * 功率谱密度

    fftPower.fftPSD

### 函数调用说明

main.py文件中通过构造模拟数据，进行函数调用说明