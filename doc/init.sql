-- pes_cloud_vpp.c_tenant_place definition

/*
CREATE TABLE `c_tenant_place` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `pid` int(11) DEFAULT NULL COMMENT '父id',
  `name` varchar(255) COLLATE utf8_general_ci NOT NULL,
  `code` varchar(255) COLLATE utf8_general_ci DEFAULT NULL,
  `remarks` varchar(255) COLLATE utf8_general_ci DEFAULT NULL,
  `num` int(10) DEFAULT '0' COMMENT '排序',
  `tenant_id` bigint(20) DEFAULT NULL,
  `place_flag` int(10) DEFAULT '0' COMMENT '是否场所节点（可以关联设备和组织）',
  `longitude` double NOT NULL DEFAULT '0' COMMENT '经度，叶节点非空',
  `latitude` double NOT NULL DEFAULT '0' COMMENT '纬度，叶节点非空',
  `type` varchar(4) COLLATE utf8_general_ci DEFAULT NULL COMMENT '场所类型，56类列管场所，叶节点非空；缺省''''，表示未知',
  `create_time` timestamp DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci AUTO_INCREMENT=120128;
*/

CREATE DATABASE `pscc` /*!40100 DEFAULT CHARACTER SET utf8 */

CREATE TABLE `ps_cc_park` (
 `park_id` BIGINT (20) NOT NULL AUTO_INCREMENT COMMENT '园区ID',
 `name` VARCHAR (20) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '园区名称',
 `business_id` BIGINT (20) NOT NULL COMMENT '商户ID',
 `opertor_id` BIGINT (20) NULL DEFAULT NULL COMMENT '操作人员ID',
 `status` TINYINT (4) NULL DEFAULT 0 COMMENT '状态：0：正常，1删除',
  `ctime` timestamp NULL DEFAULT CURRENT_TIMESTAMP  COMMENT '创建时间',
  `utime` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
 PRIMARY KEY (`park_id`)
) ENGINE = INNODB DEFAULT CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '园区表' AUTO_INCREMENT = 1 ROW_FORMAT = DYNAMIC;

CREATE TABLE `ps_cc_building` (
  `building_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '楼栋ID',
  `name` varchar(20) NOT NULL COMMENT '名称',
  `business_id` bigint(20) NOT NULL COMMENT '商户ID',
  `opertor_id` bigint(20) DEFAULT NULL COMMENT '操作人员ID',
  `park_id` bigint(20) NOT NULL COMMENT '园区ID',
  `status` tinyint(4) DEFAULT '0' COMMENT '状态：0：正常，1删除',
  `ctime` timestamp NULL DEFAULT CURRENT_TIMESTAMP  COMMENT '创建时间',
  `utime` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`building_id`),
  KEY `idx_park_id` (`park_id`) USING BTREE,
  CONSTRAINT  FOREIGN KEY (`park_id`) REFERENCES `ps_cc_park` (`park_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='楼栋表';


CREATE TABLE `ps_cc_unit` (
  `unit_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '单元ID',
  `name` varchar(20) NOT NULL COMMENT '名称',
  `business_id` bigint(20) NOT NULL COMMENT '商户ID',
  `opertor_id` bigint(20) DEFAULT NULL COMMENT '操作人员ID',
  `building_id` bigint(20) NOT NULL COMMENT '楼栋ID',
  `status` tinyint(4) DEFAULT '0' COMMENT '状态：0：正常，1删除',
  `ctime` timestamp NULL DEFAULT CURRENT_TIMESTAMP  COMMENT '创建时间',
  `utime` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`unit_id`),
  KEY `idx_building_id` (`building_id`) USING BTREE,
  CONSTRAINT  FOREIGN KEY (`building_id`) REFERENCES `ps_cc_building` (`building_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='单元表';


CREATE TABLE `ps_cc_room` (
  `room_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '房间ID',
  `name` varchar(20) NOT NULL COMMENT '名称',
  `alias` varchar(20) DEFAULT NULL COMMENT '别名称',
  `business_id` bigint(20) NOT NULL COMMENT '商户ID',
  `opertor_id` bigint(20) DEFAULT NULL COMMENT '操作人员ID',
  `unit_id` bigint(20) NOT NULL COMMENT '单元ID',
  `area` decimal(10,2) NOT NULL COMMENT '面积m2',
  `status` tinyint(4) DEFAULT 0 COMMENT '状态：0：正常，1删除,2：待租，3：已租，4:自住，5:未知',
  `remarks` varchar(100) NOT NULL COMMENT '备注',
  `ctime` timestamp NULL DEFAULT CURRENT_TIMESTAMP  COMMENT '创建时间',
  `utime` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`room_id`),
  KEY `idx_unit_id` (`unit_id`) USING BTREE,
  CONSTRAINT  FOREIGN KEY (`unit_id`) REFERENCES `ps_cc_unit` (`unit_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='房间表';

CREATE TABLE `ps_cc_room_bind` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '绑定ID',
  `name` varchar(20) DEFAULT NULL COMMENT '名称',
  `room_id` bigint(20) NOT NULL COMMENT '房间ID',
  `phone_number` varchar(16) DEFAULT NULL COMMENT '转接手机号',
  `business_id` bigint(20) NOT NULL COMMENT '商户ID',
  `opertor_id` bigint(20) DEFAULT NULL COMMENT '操作人员ID',
  `status` tinyint(4) DEFAULT '0' COMMENT '状态：0：正常，1删除',
  `ctime` timestamp NULL DEFAULT CURRENT_TIMESTAMP  COMMENT '创建时间',
  `utime` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_room_id` (`room_id`) USING BTREE,
  CONSTRAINT  FOREIGN KEY (`room_id`) REFERENCES `ps_cc_room` (`room_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='房屋绑定信息';

CREATE TABLE `ps_cc_user` (
  `user_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `name` varchar(20) DEFAULT NULL COMMENT '名称',
  `room_id` bigint(20) NOT NULL COMMENT '房间ID',
  `user_type` tinyint(4) DEFAULT '0' COMMENT '人员类型：0：业主、1：家人、2：代理人、3：租客、4：临时客人',
  `phone_number` varchar(16) DEFAULT NULL COMMENT '手机号',
  `business_id` bigint(20) NOT NULL COMMENT '商户ID',
  `opertor_id` bigint(20) DEFAULT NULL COMMENT '操作人员ID',
  `status` tinyint(4) DEFAULT '0' COMMENT '状态：0：正常，1删除',
  `avatar` varchar(250)  DEFAULT "" COMMENT '头像URL',
  `ctime` timestamp NULL DEFAULT CURRENT_TIMESTAMP  COMMENT '创建时间',
  `utime` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`user_id`),
  KEY `idx_room_id` (`room_id`) USING BTREE,
  CONSTRAINT  FOREIGN KEY (`room_id`) REFERENCES `ps_cc_room` (`room_id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='人员表';

CREATE TABLE `ps_cc_device_bind` (
  `sn` bigint(20) NOT NULL  COMMENT '设备码',
  `business_id` bigint(20) NOT NULL COMMENT '商户ID',
  `opertor_id` bigint(20) DEFAULT NULL COMMENT '操作人员ID',
  `park_id` BIGINT (20) NOT NULL DEFAULT 0 COMMENT '园区ID',
  `building_id` bigint(20) NOT NULL DEFAULT 0 COMMENT '楼栋ID',
  `unit_id` bigint(20) NOT NULL DEFAULT 0  COMMENT '单元ID',
  `room_id` bigint(20) NOT NULL DEFAULT 0  COMMENT '房间ID',
  `status` tinyint(4) DEFAULT '0' COMMENT '状态：0：正常，1删除',
  `ctime` timestamp NULL DEFAULT CURRENT_TIMESTAMP  COMMENT '创建时间',
  `utime` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`sn`),
  KEY `idx_pid_bid_uid_rid` (`park_id`,`building_id`,`unit_id`,`room_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='设备绑定表';