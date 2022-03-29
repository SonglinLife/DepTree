# DepTree

`DepTree`根据xlab实验室数据开发的依赖解析工具，可以实现包括组件依赖解析、依赖树的构建以及可视化保存、组件历史版本号的拉取以及多种指标计算。

## 使用说明

### clone代码到本地

`git clone git@jihulab.com:x-laber/deptree.git`

`cd deptree/bin/`

`chmod 700 deptree`

### config.json配置文件

**请确保您拥有config.json文件，并把该文件放于bin目录下，该文件用于配置数据库认证参数。**

```json
{
   "tunnelIp": "":,
   "tunnelPasswd": "",
   "dbAddr": "",
   "dbUser": "",
   "dbPassword": ""
}
```

如果您没有持有该文件，请加入xlab社区和管理员联系。

### 参数说明

##### 查询组件的MTTU或者外部依赖度

`deptree -t <type> -n <name> -v <verison> `

`-t`：1代表查mttu， 2表示查外部依赖度。

`-n`：组件名称。

`-v`：组件版本号。

示例：`./deptree -t 2 -n webpack -v 5.58.1`

结果：`webpack@5.58.1 external:25.876045`

其中外部依赖度的定义如下：

![image-20220329195658375](https://image-2021-wu.oss-cn-beijing.aliyuncs.com/blogs/picturesimage-20220329195658375.png)

MTTU的定义如下：

![image-20220329195715073](https://image-2021-wu.oss-cn-beijing.aliyuncs.com/blogs/picturesimage-20220329195715073.png)

其中，$T_\text{update\_dependecy }$ 是指这个依赖项新版本发布的时间，$T_\text{dependecy\_release}$ 是指在该项目中，将依赖更新到新版本的时间。

#### 拉取组件的历史版本

deptree支持查询组件的所有历史版本号，以及该版本的发行时间

示例：`./deptree -t 3 -n react`

结果保存在`react_all_version.json`中：

```json
{
    "Version": "0.0.1",
    "Name": "react",
    "Time": "2011-10-26T17:46:22.746Z"
},
{
    "Version": "0.0.2",
    "Name": "react",
    "Time": "2011-10-28T22:40:36.115Z"
}
```

#### 构建组件的依赖树

`deptree`根据`semver`规则依赖版本号进行解析，`semver`的解析规则可以查看官网说明：[语义化版本 2.0.0](https://semver.org/lang/zh-CN/)

示例：`./deptree -t 4 -n react -v 15.6.2 `

在命令行中会打印出依赖树的简单可视化：

![image-20220329191325361](https://image-2021-wu.oss-cn-beijing.aliyuncs.com/blogs/picturesimage-20220329191325361.png)

在`react@15.6.2DepTree.json`中保存了`react`组件的依赖树：

```json
{
    "name": "react",
    "version": "15.6.2",
    "time": "2017-09-26T00:10:25.817Z",
    "Index": 1,
    "parentIndex": 0
}
```

`metadata`的格式如上，每个组件简单都拥有唯一的`Index`，并且拥有一个`parentIndex`字段指明其父节点的`Index`，其中根节点的`Index`为1，其`parentIndex`为0。

## 数据集

在项目的`dataset`目录下，提供了npm组件采样数据`1.5w_npm_meta_2021.json`：

```json
{
    "name": "@prazdevs/eslint-config-javascript",
    "releases": [
        {
            "version": "1.0.0",
            "time": "2021-07-04T00:14:56.638Z"
        },
        {
            "version": "1.0.1",
            "time": "2021-07-04T23:30:38.52Z"
        },
        {
            "version": "1.0.2",
            "time": "2021-07-05T00:11:57.004Z"
        }    
	]
}
```

该采样数据包括了1.5万个npm组件2021年发布的所有版本号和发行时间。

同时`dataset`目录下，还有`mttu.json`和`external.json`两个json文件，分别对应这1.5万个组件的mttu和外部依赖度。

`external.json`:

```json
{
    "name": "@prazdevs/eslint-config-javascript",
    "verion": "1.1.1",
    "time": "2021-07-06T20:37:09.143Z",
    "externaldep": 14.6848955
}
```

`mttu.json`:

```json
{
    "name": "@prazdevs/eslint-config-javascript",
    "verion": "2.1.1",
    "time": "2021-07-10T00:57:05.443Z",
    "mttu": 106.55556
}
```

## 采样

`deptree`并没有将采样和批量计算指标功能编译到可执行文件中，但您可以在项目子文件夹`sample`中找到采样数据的相关代码。

