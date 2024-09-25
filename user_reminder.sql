/*
 Navicat Premium Data Transfer

 Source Server         : Demo
 Source Server Type    : MySQL
 Source Server Version : 80028
 Source Host           : localhost:3306
 Source Schema         : demo

 Target Server Type    : MySQL
 Target Server Version : 80028
 File Encoding         : 65001

 Date: 26/09/2024 04:08:23
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for user_reminder
-- ----------------------------
DROP TABLE IF EXISTS `user_reminder`;
CREATE TABLE `user_reminder`  (
  `id` int(0) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `creator_id` int(0) NULL DEFAULT NULL COMMENT '创建者id',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '提醒内容',
  `reminder_at` datetime(0) NULL DEFAULT NULL COMMENT '提醒时间',
  `send_type` tinyint(1) NULL DEFAULT NULL COMMENT '发送方式(1电话2邮箱)',
  `contact_info` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '联系方式(电话/邮箱)',
  `deleted` tinyint(1) NULL DEFAULT 0 COMMENT '逻辑删除(0否1是)',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户内容提醒' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_reminder
-- ----------------------------
INSERT INTO `user_reminder` VALUES (2, 2, '前往大学城南', '2024-09-26 02:30:00', 2, '837425169@qq.com', 1);
INSERT INTO `user_reminder` VALUES (3, 2, '前往天河体育中心', '2024-09-25 16:24:00', 2, '837425169@qq.com', 1);
INSERT INTO `user_reminder` VALUES (4, 1, '前往南村万博', '2024-09-26 02:31:29', 2, '837425169@qq.com', 0);
INSERT INTO `user_reminder` VALUES (5, 3, '前往客村', '2024-09-26 10:39:00', 2, '837425169@qq.com', 0);
INSERT INTO `user_reminder` VALUES (7, 3, '前往广州塔', '2024-09-26 10:39:00', 2, '837425169@qq.com', 0);
INSERT INTO `user_reminder` VALUES (8, 3, '前往广州塔', '2024-09-26 10:39:00', 2, '837425169@qq.com', 0);
INSERT INTO `user_reminder` VALUES (9, 3, '前往广州塔', '2024-09-26 03:35:00', 2, '837425169@qq.com', 0);

SET FOREIGN_KEY_CHECKS = 1;
