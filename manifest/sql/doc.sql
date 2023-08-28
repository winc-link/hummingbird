-- -------------------------------------------------------------
-- TablePlus 3.11.0(352)
--
-- https://tableplus.com/
--
-- Database: core.db
-- Generation Time: 2023-08-18 16:34:02.6110
-- -------------------------------------------------------------


DROP TABLE IF EXISTS "doc";
CREATE TABLE `doc` (`created` integer,`modified` integer,`id` text NOT NULL,`name` text,`sort` integer,`jump_link` text,PRIMARY KEY (`id`));

INSERT INTO "doc" ("created", "modified", "id", "name", "sort", "jump_link") VALUES
('0', '0', '11695073', 'SDK参考', '2', 'https://doc.hummingbird.winc-link.com/guide/sdk/go.html'),
('0', '0', '52093426', 'API参考', '3', 'https://doc.hummingbird.winc-link.com/guide/api/call.html'),
('0', '0', '71610269', '产品概述', '1', 'https://doc.hummingbird.winc-link.com/guide/product/introduce.html');
