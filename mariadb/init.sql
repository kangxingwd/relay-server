
CREATE DATABASE IF NOT EXISTS webdb DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;


CREATE TABLE `account` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `created` datetime DEFAULT CURRENT_TIMESTAMP,
  `modified` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

INSERT INTO `webdb`.`account` (`id`, `name`, `password`) VALUES ('1', 'admin', 'admin');

DROP TABLE IF EXISTS `fb_attachment`;
CREATE TABLE `fb_attachment` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `fid` int(10) NOT NULL COMMENT '对应的feedback表的id',
  `type` varchar(255) DEFAULT NULL COMMENT '附件类型',
  `filename` varchar(255) DEFAULT NULL COMMENT '附件名',
  `path` varchar(255) DEFAULT NULL COMMENT '附件存储路径',
  `size` int(10) DEFAULT NULL COMMENT '附件大小  单位：字节',
  `ext` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `feedback`;
CREATE TABLE `feedback` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `device_id` varchar(255) DEFAULT NULL COMMENT '反馈的设备id',
  `product_class` varchar(255) DEFAULT NULL COMMENT '设备的产品型号',
  `ac_version` varchar(255) DEFAULT NULL COMMENT '反馈提交时设备的软件版本',
  `account` varchar(255) DEFAULT NULL COMMENT '提交反馈的账号',
  `app_version` varchar(255) DEFAULT NULL COMMENT '反馈提交时APP的软件版本',
  `type` tinyint(8) DEFAULT NULL COMMENT '反馈类型 1-功能故障  2-产品建议 3-其它',
  `content` varchar(1024) CHARACTER SET utf8 DEFAULT NULL COMMENT '反馈内容',
  `contact_info` varchar(255) DEFAULT '' COMMENT '联系方式',
  `commit_time` varchar(255) DEFAULT NULL COMMENT '提交时间',
  `status` tinyint(8) DEFAULT NULL COMMENT '处理状态 1-未处理 2-已处理 3-已解决',
  `handle_way` tinyint(8) DEFAULT NULL COMMENT '处理方式 1-接受 2-忽略',
  `handle_desc` varchar(1024) CHARACTER SET utf8 DEFAULT NULL COMMENT '处理说明',
  `handle_time` varchar(255) DEFAULT NULL COMMENT '最后处理时间',
  `ext` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4;
