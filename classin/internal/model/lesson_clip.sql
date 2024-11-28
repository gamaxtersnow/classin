CREATE TABLE `lesson_clip`
(
    id                        int unsigned auto_increment comment '主键id',
    courseId                  int unsigned          default 0                 not null comment 'classin 课程id',
    classId                   int unsigned          default 0                 not null comment 'classin 课节id',
    sequenceNumber            int unsigned          default 0                 not null comment '片段序号',
    fileId                    varchar(255)          default ''                not null comment '文件id',
    fileOriginUrl             varchar(255)          default ''                not null comment '片段原始地址',
    objectKey                 varchar(255)          default ''                not null comment '片段保存地址',
    fileSize                  int unsigned          default 0                 not null comment '文件大小',
    fileStatus                int unsigned          default 0                 not null comment '0 已添加 1 同步中 2 同步完成',
    sourceType                int unsigned          default 0                 not null comment '1 手动抓取 2 消息推送',
    addTime                   int unsigned          default 0                 not null comment '添加时间',
    PRIMARY KEY(`id`),
    constraint uk_file_origin_url
        unique (fileOriginUrl)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;