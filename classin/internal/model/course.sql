create table course
(
    id              int unsigned auto_increment comment '主键id',
    uniqueId        varchar(255) default ''  not null comment '唯一id',
    courseId        int unsigned default 0   not null comment 'classin 课程id',
    classId         int unsigned default 0   not null comment 'classin 课节id',
    courseName      varchar(255) default ''  not null comment 'classin 课程名称',
    className       varchar(255) default ''  not null comment 'classin 课节名称',
    teacherName     varchar(255) default ''  not null comment '老师姓名',
    teacherMobile   varchar(255) default ''  not null comment '老师手机号',
    courseStartTime int unsigned default '0' not null comment '课程开始时间',
    courseEndTime   int unsigned default '0' not null comment '课程结束时间',
    syncStatus      int unsigned default 0   not null comment '0 已添加 1 同步中 2 同步完成',
    sourceType      int unsigned default 0   not null comment '1 手动抓取 2 消息推送',
    addTime         int unsigned default 0   not null comment '添加时间',
    seatNum         int unsigned default 0   not null comment '上课人数',
    isDc            int unsigned default 0   not null comment '是否双摄',
    isHd            int unsigned default 0   not null comment '是否高清',
    PRIMARY KEY(`id`),
    constraint uk_unique_id
        unique (uniqueId)
);