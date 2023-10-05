drop table if exists tokens cascade;

CREATE TABLE `tokens` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_id` int(11) unsigned NOT NULL,
    `name` varchar(255) NOT NULL,
    `email` varchar(255) NOT NULL,
    `token` varchar(255) NOT NULL,
    `token_hash` varbinary(255) DEFAULT NULL,
    `created_at` datetime NOT NULL DEFAULT current_timestamp(),
    `updated_at` datetime NOT NULL DEFAULT current_timestamp(),
    `expiry` datetime NOT NULL,
    PRIMARY KEY (`id`),
    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE cascade ON DELETE cascade
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8mb4;

drop table if exists oauth_clients cascade;

CREATE TABLE `oauth_clients` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_id` int(11) unsigned NOT NULL,
    `name` varchar(255) NOT NULL,
    `secret` varchar(255) NOT NULL,
    `token` varchar(255) NOT NULL,
    `revoked` int(11) NOT NULL,
    `created_at` datetime NOT NULL DEFAULT current_timestamp(),
    `updated_at` datetime NOT NULL DEFAULT current_timestamp(),
    PRIMARY KEY (`id`),
    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE cascade ON DELETE cascade
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8mb4;
