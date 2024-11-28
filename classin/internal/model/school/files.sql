create table `files` (
                         `id` int NOT NULL AUTO_INCREMENT,
                         `uuid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '唯一id',
                         `object_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'oss key',
                         `name`        varchar(255) default '' not null comment '文件名称',
                         `file_type` varchar(50) NOT NULL DEFAULT '' COMMENT '文件类型',
                         `creator_id`  bigint       default 0  not null comment '创建人ID',
                         `add_time` int NOT NULL DEFAULT '0' COMMENT '创建时间',
                         `update_time` int          default 0  not null comment '更新时间',
                         `status` int NOT NULL DEFAULT '0' COMMENT '文件状态 1 正常 2 删除',
                         `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '备注',
                         PRIMARY KEY (`id`),
                         UNIQUE KEY `idx_uuid` (`uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='导出文件表';