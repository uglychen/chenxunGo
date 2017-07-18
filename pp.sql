#用户表userInfo 用户id 微信openID 昵称 头像地址 性别 年龄 金币 现金
CREATE TABLE IF NOT EXISTS `ppserver`.`userInfo`(
   `userId`     int UNSIGNED AUTO_INCREMENT NOT null,
   `wechatId`   VARCHAR(100)    DEFAULT null unique,
   `nickName`   VARCHAR(100)    DEFAULT null,
   `headUrl`    VARCHAR(1000)   DEFAULT null,
   `gender`     VARCHAR(15)     DEFAULT null,
   `age`        TINYINT         DEFAULT null,
   `gold`       int             DEFAULT null,
   `cash`       int             DEFAULT null,
   PRIMARY KEY (`userId`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `ppserver`.`userInfo` ( `wchatId`, `nickName`, `headUrl`, 
    `gender`, `age`, `gold`, `cash`)
     VALUES ('chen_123456', 'chen', 'http://', 'male', '12', '0', '0');


#任务内容表
CREATE TABLE IF NOT EXISTS `ppserver`.`taskTable`(
   `taskId`      VARCHAR(100)    NOT null,
   `taskType`    VARCHAR(100)    DEFAULT null,
   `taskName`    VARCHAR(100)    DEFAULT null,
   `taskContent` VARCHAR(256)    DEFAULT null,
   `taskStatus`  VARCHAR(20)     DEFAULT null,
   `validTime`   TIMESTAMP       DEFAULT CURRENT_TIMESTAMP,
   `addGoldCoin` int             DEFAULT null,
   `addCash`     int             DEFAULT null,
   PRIMARY KEY ( `taskId` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

#任务事件表
CREATE TABLE IF NOT EXISTS `ppserver`.`eventTable`(
   `userId`         INT UNSIGNED      NOT null,
   `taskId`         VARCHAR(100)      not null,
   `createTime`     TIMESTAMP         DEFAULT CURRENT_TIMESTAMP,
   `endTime`        TIMESTAMP,        DEFAULT CURRENT_TIMESTAMP,
   `taskStatus`     VARCHAR(20)       DEFAULT null,
   `addGoldCoin`    int               DEFAULT null,
   `addCash`        int               DEFAULT null,
   `eventId`        INT UNSIGNED      DEFAULT null,
   PRIMARY KEY ( `userId` ,`taskId`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;

#金币现金表
CREATE TABLE IF NOT EXISTS `ppserver`.`goldAndCash`(
   `userId`         INT UNSIGNED      NOT null,
   `eventId`        INT UNSIGNED      DEFAULT null,
   `taskType`       VARCHAR(100)      DEFAULT null,
   `eventtime`      TIMESTAMP         DEFAULT CURRENT_TIMESTAMP,
   `currentGold`    int               DEFAULT null,
   `currentBalance` int               DEFAULT null,
   `addGoldCoin`    int               DEFAULT null,
   `addCash`        int               DEFAULT null,
   `glodBalance`    int               DEFAULT null,
   `cashBalance`    int               DEFAULT null,
   PRIMARY KEY ( `userId` )
)ENGINE=InnoDB DEFAULT CHARSET=utf8;




