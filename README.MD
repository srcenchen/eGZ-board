# eGZ-Board
用于显示教学信息的壁纸软件，功能包括高考倒计时、公告栏、以及能够提示当前课程的课程表，包括当前天气信息。并且提供班牌版本
## 教程
### 启动
首先启动main.exe服务端，在到DreamScene2中启动壁纸软件 这个软件是壁纸接管软件
使用软件打开 http://serverip:5020
随后会要求选择班级 要求绑定，如需解除绑定 连续点击进度条10下即可
### 管理
http://serverip:5020/setting
课程表模板：[课程表模板.xlsx](images%2F%E8%AF%BE%E7%A8%8B%E8%A1%A8.xlsx)
## 班牌使用
搭建好服务器后 班牌地址为 http://serverip:5020/outside 本地址为班牌特殊优化
## 截图
![s1.png](images%2Fs1.png)
![s2.png](images%2Fs2.png)
## 为啥要做成另起一个服务端？
本项目旨在为整个学校服务，并且为班牌使用，此外班牌为Android系统，必须另起服务。

使用服务器端可以方便管理，一个班级为一个单位，大屏幕与班牌共用一套数据。
## 鸣谢
本项目是对ClassBoard的拙劣模仿，在其基础上实现了多班级管理，天气预报，提供班牌版本
[@Candlest](https://github.com/Candlest/)