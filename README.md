# MIT6.824 Solutions
### 2019/11/29
1. 这是MIT分布式系统课的实验答案仓库,目前完成Lab1的部分
2. 计划在三个月内完成所有Lab
3. 如果你也在学习这个课程，那么建议先熟悉一下go，对做实验感觉会很有帮助

### 0. 程序大体流程

​	这个简单的map-reduce库的入口是`Sequential()`或者`Distributed()`函数，输入的参数包括需要处理的文件、自己定义的`map`和`reduce`函数以及reduce task的数量

​	入口函数首先会创建一个`master`，对于分布式处理的`master`，它会启动一个RPC server，并且等待workers注册，然后根据`schedule()`函数（part3、4）中定义的规则为workers分配任务

​	对于`map worker`，它会调用`doMap()`（part1），这个函数中会创建中间文件、调用用户定义的`map()`函数（part2）、将`map`函数返回的中间结果哈希到对应的中间文件中；总之，在`map worker`端除`map()`函数定义了以外的所有事情都由它来处理

​	对于`reduce worker`，在所有中间文件生成完毕后调用`doReduce()`（part1），这个函数将中间文件中的每一个键值对输入到用户定义的`reduce()`然后将处理结果存到文件中

​	所有任务完成后`master`会调用`merge()`将所有结果文件归并在一起形成最终结果

​	注意：这只是一个简化版的map-reduce库，它用多线程模拟多结点，也没有使用诸如HDFS之类的分布式文件系统，中间文件、结果文件都存在本地文件系统中，有一些与论文描述以及HADOOP中实现的map-reduce不吻合的地方

### 1. Lab1总结

#### 1.1 part1

​	这个实验分为part1-4，其中part1要求我们实现调用用户定义的`map()`和`reduce()`函数的函数`doMap()`和`doReduce()`，其中`doMap()`的大体思路为：

- 读取输入文件
- 根据`nReduce`创建中间文件
- 将文件输入到`map()`中
- 将返回结果一条一条的哈希到对应的中间文件中

`doReduce()`的思路为：

- 读取`nReduce`个文件
- 将这些文件拼到一起，并且排序
- 对文件中每一个键值对调用`reduce()`，并将结果存到文件中

#### 1.2 part2

​	part2要求我们实现经典例子word-count的`map()`和`reduce()`，相对简单，注意一下文件读取细节、换行符处理就行

#### 1.3 part3 & part4

 	这一部分要求我们实现`schedule()`，这个函数的作用就是为workers分配任务，程序的大体思路为：

- 确定阶段为`map`还是`reduce`
- 填充参数`arg`，根据阶段不同填充有细微区别
- 进入循环，对每一个任务分配一个worker，这里建立了一个任务队列，失败的任务会重新进入队列，当成功数等于`nReduce`结束循环
- 这里顺便完成了part4，即关于worker执行失败的部分，失败的任务重新进入任务队列即可

#### 1.4 总结

​	Lab1以及map-reduce论文结合可以有效提高对这个计算模型的理解，总体来说，还存在以下几点疑问：

- 关于shuffle的详细过程
- 关于combine函数
- 几次排序的意义