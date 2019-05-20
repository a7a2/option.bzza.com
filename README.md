# option.bzza.com


一份2018年写的


外汇、指数、期货、加密数字货币、贵金属、CFD的二次期权交易系统


配合 https://github.com/a7a2/bzza.com 使用方可正常运作


还需要自行导入各种交易对的历史数据到KDB+数据库（否则没有历史数据）

在helpers/historyDataIn 目录下有简单导入例子，数据源自己找

主目录的config.yaml是配置文件


记得安装好redis、kdb+、postgresql，

其中kdb+启动9999端口，启动命令q.exe -p 9999

redis默认端口或自行修改配置文件

postgresql安装好后按照配置文件的设置新建一个数据库，对应账号密码

配置文件中的wsORwss项目就是另外一个项目的api接口websocket地址


仅用于学习，切莫用于商业用途！