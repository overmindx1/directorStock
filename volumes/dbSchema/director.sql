-- --------------------------------------------------------
-- 主機:                           127.0.0.1
-- 伺服器版本:                        10.5.10-MariaDB-1:10.5.10+maria~focal - mariadb.org binary distribution
-- 伺服器作業系統:                      debian-linux-gnu
-- HeidiSQL 版本:                  10.3.0.5771
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;


-- 傾印 stocks 的資料庫結構
CREATE DATABASE IF NOT EXISTS `stocks` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */;
USE `stocks`;

-- 傾印  資料表 stocks.director_selection 結構
CREATE TABLE IF NOT EXISTS `director_selection` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `stock_id` varchar(10) DEFAULT '0' COMMENT '股票代號',
  `stock_name` varchar(16) DEFAULT '0' COMMENT '股票名稱',
  `avg_cost` varchar(15) DEFAULT '0' COMMENT '平均成本',
  `shares` varchar(7) DEFAULT '0' COMMENT '股數',
  `on_list` tinyint(1) DEFAULT 1 COMMENT '是否在會長存股上',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=334 DEFAULT CHARSET=utf8mb4 COMMENT='會長的選股';

-- 取消選取資料匯出。

-- 傾印  資料表 stocks.stock_items 結構
CREATE TABLE IF NOT EXISTS `stock_items` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `stock_id` varchar(10) DEFAULT NULL COMMENT '股票代號',
  `stock_name` varchar(16) DEFAULT NULL COMMENT '股票名稱',
  `open_price` varchar(8) DEFAULT NULL COMMENT '開盤價',
  `close_price` varchar(8) DEFAULT NULL COMMENT '收盤價',
  `up_down` varchar(1) DEFAULT NULL COMMENT '漲或跌',
  `up_down_count` varchar(7) DEFAULT NULL COMMENT '漲跌幅',
  `date` varchar(50) DEFAULT NULL COMMENT '更新的日期',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '建立時間',
  `updated_at` timestamp NULL DEFAULT NULL COMMENT '更新時間',
  PRIMARY KEY (`id`),
  UNIQUE KEY `stock_id` (`stock_id`)
) ENGINE=InnoDB AUTO_INCREMENT=31988 DEFAULT CHARSET=utf8mb4 COMMENT='股票大全';

-- 取消選取資料匯出。

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
