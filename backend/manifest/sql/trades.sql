CREATE TABLE IF NOT EXISTS `trades` (
                                        `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                                        `tx_hash` varchar(66) NOT NULL,
    `user_address` varchar(42) NOT NULL,
    `token_in` varchar(42) NOT NULL,
    `token_out` varchar(42) NOT NULL,
    `amount_in` varchar(78) NOT NULL,
    `amount_out` varchar(78) NOT NULL,
    `status` varchar(20) NOT NULL,
    `block_number` bigint unsigned DEFAULT NULL,
    `confirmations` int unsigned DEFAULT '0',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_tx_hash` (`tx_hash`),
    KEY `idx_user_address` (`user_address`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;