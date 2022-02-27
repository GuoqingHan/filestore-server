# filestore-server
## 新知识
### 分库分表
1、水平分表
每个上传文件的hash值的后两位来切分，分别对应数据库的tbl_00, tbl_01, ..., tbl_ff，共计256张table
2、垂直分表

### 数据库
1. 预编译方式，创建sql操作，防止sql注入
2. replace into 跟 insert into 功能类似，不同点在于：replace into 首先尝试插入数据到表中， 
(1) 如果发现表中已经有此行数据（根据主键或者唯一索引判断）则先删除此行数据，然后插入新的数据。 否则，直接插入新数据。
(2) 要注意的是：插入数据的表必须有主键或者是唯一索引！否则的话，replace into 会直接插入数据，这将导致表中出现重复的数据
3. ON UPDATE CURRENT_TIMESTAMP:自动更新时间戳，例如：`update_at` datetime default NOW() on update current_timestamp() COMMENT '更新日期'。

### go写法
1、数据表字段值类型转换为string：
if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
    return true
}
2、上传文件：
file, head, err := r.FormFile("file")
// 在服务器创建文件，用来存储客户发送过来的文件
newFile, err := os.Create("/tmp/" + head.Filename)
// 内存文件，拷贝到新的buffer中
filesize, err = io.Copy(newFile, file)
// 如果为了计算SHA1值，newFile.Seek(0, 0) // !!! 移动到文件的起始位置
3、前端下载文件，后端的代码如下：
// 下载文件需修的header
w.Header().Set("Content-type", "application/octect-stream")
w.Header().Set("Content-Descrption", "attachment;filename=\""+fm.FileName+"\"")
w.Write(data)

## 环境配置
### github访问慢，go get xx不下来
nslookup github.global.ssl.fastly.net
nslookup github.com
在/etc/hosts添加ip–>域名的映射，刷新DNS缓存便可。


结构体成员类型是interface，表示可以是任何类型

