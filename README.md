# filestore-server
## 新知识
### 分库分表
1、水平分表
每个上传文件的hash值的后两位来切分，分别对应数据库的tbl_00, tbl_01, ..., tbl_ff，共计256张table
2、垂直分表

预编译方式，创建sql操作，防止sql注入

## mysql表
### tbl_file
create table `tbl_file`(
    `id` int(11) not null auto_increment,
    `file_sha1` char(40) not null default '' comment '文件hash',
    `file_name` varchar(256) not null default '' comment '文件名',
    `file_size` bigint(20) default '0' comment '文件大小',
    `file_addr` varchar(1024) not null default '' comment '创建日期',
    `create_at` datetime default NOW() comment '创建日期',
    `update_at` datetime default NOW() on update current_timestamp() comment '更新日期',
    `status` int(11) not null default '0' comment '状态(可用/禁用/已删除等状)',
    `ext1` int(11) default '0' comment '备用字段1',
    `ext2` int(11) default '0' comment '备用字段2',
    primary key (`id`),
    unique key `idx_file_hash` (`file_sha1`),
    key `idx_status` (`status`)
) engine=InnoDB default charset=utf8;