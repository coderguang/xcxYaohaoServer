/*
Navicat MySQL Data Transfer

Source Server         : 47.107.177.155
Source Server Version : 50505
Source Host           : 47.107.177.155:3306
Source Database       : xcx_template

Target Server Type    : MYSQL
Target Server Version : 50505
File Encoding         : 65001

Date: 2019-08-28 20:06:08
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for xcx_card_data_template
-- ----------------------------
DROP TABLE IF EXISTS `xcx_card_data_template`;
CREATE TABLE `xcx_card_data_template` (
  `type` int(11) NOT NULL DEFAULT 0 COMMENT '类型：1个人，2公司',
  `card_type` int(11) NOT NULL DEFAULT 0 COMMENT '车牌类型：1普通牌，2节能牌',
  `code` varchar(100) NOT NULL,
  `name` varchar(920) NOT NULL,
  `time` varchar(11) NOT NULL,
  `tips` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`code`,`name`),
  KEY `id` (`code`) USING BTREE,
  KEY `name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for xcx_download_data_template
-- ----------------------------
DROP TABLE IF EXISTS `xcx_download_data_template`;
CREATE TABLE `xcx_download_data_template` (
  `url` varchar(900) NOT NULL DEFAULT '',
  `status` varchar(255) NOT NULL DEFAULT '',
  `title` varchar(1000) NOT NULL DEFAULT '',
  `tips` varchar(1000) NOT NULL DEFAULT '',
  PRIMARY KEY (`url`),
  KEY `url` (`url`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
