CREATE DATABASE IF NOT EXISTS mini_video CHARACTER SET UTF8mb4 COLLATE utf8mb4_bin;
CREATE DATABASE IF NOT EXISTS mini_video_dblog CHARACTER SET UTF8mb4 COLLATE utf8mb4_bin;

use mini_video;
# 用户
CREATE TABLE IF NOT EXISTS `user` (
    `openId`        VARCHAR(64) NOT NULL PRIMARY KEY,
    `uid`           INT(11) NOT NULL,
    `name`          VARCHAR(64),
    `vipTime`       BIGINT NOT NULL,
    `version`       VARCHAR(32),
    `password`      VARCHAR(32),
    `data`          MEDIUMBLOB,
    `clientIP`      VARCHAR(64),
    `createTime`    BIGINT COMMENT '创建时间',
    `device`        VARCHAR(64),
     INDEX uidIndex(uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

# id 生成器
CREATE TABLE IF NOT EXISTS `uid_generator`(
    `uid`       BIGINT AUTO_INCREMENT PRIMARY KEY
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
ALTER TABLE `uid_generator` AUTO_INCREMENT=10000;

# 支付
CREATE TABLE IF NOT EXISTS `pay_order` (
    `id` VARCHAR(64) NOT NULL COMMENT '订单id',
    `itemId` INT NOT NULL,
    `itemName` VARCHAR(128) NOT NULL COMMENT '商品名字',
    `itemDes` VARCHAR(128) NOT NULL COMMENT '商品描述',
    `itemPrice` INT NOT NULL COMMENT '商品价钱',
    `token` VARCHAR(64) NOT NULL COMMENT '交易的token号',
    `tokenParam` VARCHAR(128) NOT NULL COMMENT '参数',
    `createTime` INT NOT NULL COMMENT '创建时间',
    `uid` INT NOT NULL COMMENT '玩家id',
    `platform` INT NOT NULL DEFAULT '0' COMMENT '平台（android，ios）',
    `payStatus` INT DEFAULT '0' COMMENT '支付状态，0创建订单',

    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

# 会话
CREATE TABLE IF NOT EXISTS `user_session` (
    `uid` INT NOT NULL COMMENT '用户id',
    `sessionId` VARCHAR(64) NOT NULL COMMENT '会话id',
    `expire` BIGINT COMMENT '过期时间',

    PRIMARY KEY (`uid`),
    INDEX sessionIdIndex(sessionId)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

# CREATE TABLE IF NOT EXISTS `user_session` (
    `uid` INT NOT NULL COMMENT '用户id',
    `sessionId` VARCHAR(64) NOT NULL COMMENT '会话id',
    `expire` BIGINT COMMENT '过期时间',

    PRIMARY KEY (`uid`),
    INDEX sessionIdIndex(sessionId)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

# 兑换码配置
CREATE TABLE IF NOT EXISTS `cdkey_config` (
    `cdkey` VARCHAR(64) NOT NULL COMMENT '兑换码',
    `num` INT NOT NULL COMMENT '兑换码数量',
    `cdkeyType` INT NOT NULL COMMENT '兑换码类型',
    `items` VARCHAR(256) COMMENT '兑换码奖励内容',
    `createTime` BIGINT COMMENT '创建时间',
    `expireTime` BIGINT COMMENT '过期时间',

    PRIMARY KEY (`cdkey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

# 兑换码记录 (量大考虑分库分表)
CREATE TABLE IF NOT EXISTS `user_cdkey_use` (
    `index` INT AUTO_INCREMENT PRIMARY KEY COMMENT '序号',
    `uid` INT NOT NULL COMMENT '用户',
    `cdkey` VARCHAR(64) NOT NULL COMMENT '兑换码',
    `createTime` BIGINT COMMENT '时间',

    INDEX uidIdIndex(uid),
    INDEX codeIdIndex(cdkey)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

# 玩家收藏 (量大考虑分库分表)
CREATE TABLE IF NOT EXISTS `user_video_collect` (
    `index` INT AUTO_INCREMENT PRIMARY KEY COMMENT '序号',
    `uid` INT NOT NULL COMMENT '用户',
    `vid` INT NOT NULL COMMENT '视频id',
    `eid` INT NOT NULL COMMENT '集数',
    `createTime` BIGINT COMMENT '时间',

    INDEX uidIdIndex(uid),
    INDEX vidIndex(vid),
    INDEX eidIndex(eid)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

# 内容配置
CREATE TABLE IF NOT EXISTS `video_content` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT 'id',
    `name` VARCHAR(256) NOT NULL COMMENT '名称',
    `data` MEDIUMTEXT NOT NULL COMMENT '播放列表',
    `cover` VARCHAR(256) NOT NULL COMMENT '封面',
    `total` INT NOT NULL COMMENT '集数',
    `desc` VARCHAR(256) COMMENT '描述',
    `label` VARCHAR(512) COMMENT '标签'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
