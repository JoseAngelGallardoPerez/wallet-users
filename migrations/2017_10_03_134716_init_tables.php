<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class InitTables extends Migration
{
    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
    }

    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        // skip the migration if there are another migrations
        // It means this migration was already applied
        $migrations = DB::select('SELECT * FROM migrations LIMIT 1');
        if (!empty($migrations)) {
            return;
        }
        $oldMigrationTable = DB::select("SHOW TABLES LIKE 'schema_migrations'");
        if (!empty($oldMigrationTable)) {
            return;
        }

        DB::beginTransaction();

        try {
            app("db")->getPdo()->exec($this->getSql());
        } catch (\Throwable $e) {
            DB::rollBack();
            throw $e;
        }

        DB::commit();
    }

    private function getSql()
    {
        return <<<SQL
            CREATE TABLE `blocked_ips` (
              `id` int(11) NOT NULL COMMENT 'Primary Key: unique ID for IP addresses.',
              `ip` int(10) UNSIGNED NOT NULL COMMENT 'IP address',
              `created_at` timestamp NULL DEFAULT NULL COMMENT 'Timestamp added',
              `blocked_until` timestamp NULL DEFAULT NULL
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Stores blocked IP addresses.';

            CREATE TABLE `confirmation_codes` (
              `id` int(11) NOT NULL,
              `user_uid` varchar(255) NOT NULL,
              `subject` varchar(255) NOT NULL,
              `code` varchar(255) NOT NULL,
              `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
              `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
              `expires_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

            CREATE TABLE `devices` (
              `id` int(11) NOT NULL COMMENT 'The device id',
              `uid` varchar(255) NOT NULL COMMENT 'The user ID',
              `pin` varchar(45) DEFAULT NULL,
              `push_token` varchar(255) NOT NULL,
              `os_type` enum('ios','android') DEFAULT NULL,
              `created_at` timestamp NULL DEFAULT NULL,
              `updated_at` timestamp NULL DEFAULT NULL
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

            CREATE TABLE `fail_auth_attempts` (
              `id` int(11) NOT NULL COMMENT 'Primary Key: Unique ID',
              `ip` int(10) UNSIGNED DEFAULT NULL COMMENT 'The IP address of the visitor that attempted to auth',
              `uid` varchar(255) DEFAULT NULL,
              `created_at` timestamp NULL DEFAULT NULL COMMENT 'Timestamp of the failed attempt'
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Stores site fail auth attempts';

            CREATE TABLE `registration_requests` (
              `id` int(11) NOT NULL,
              `uid` varchar(255) DEFAULT NULL,
              `status` enum('accepted','pending','canceled') NOT NULL,
              `cancellation_reason` text,
              `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
              `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
              `username` varchar(255) NOT NULL,
              `first_name` varchar(255) DEFAULT NULL,
              `last_name` varchar(255) DEFAULT NULL
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

            CREATE TABLE `schema_migrations` (
              `version` bigint(20) NOT NULL,
              `dirty` tinyint(1) NOT NULL
            ) ENGINE=InnoDB DEFAULT CHARSET=latin1;

            INSERT INTO `schema_migrations` (`version`, `dirty`) VALUES
            (20190627024114, 0);

            CREATE TABLE `security_questions` (
              `sqid` int(11) NOT NULL COMMENT 'The security question ID',
              `question` varchar(255) DEFAULT NULL COMMENT 'The text of the question',
              `uid` varchar(255) NOT NULL DEFAULT '0' COMMENT '0 for questions available system-wide, or the owning uid for custom per-user questions.'
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Contains possible security questions';

            INSERT INTO `security_questions` (`sqid`, `question`, `uid`) VALUES
            (1, 'Your city of birth?', '0'),
            (2, 'Name of your last school?', '0'),
            (3, 'Name of your first pet?', '0');

            CREATE TABLE `security_questions_answers` (
              `aid` int(11) NOT NULL,
              `sqid` int(11) DEFAULT NULL COMMENT 'The security question ID',
              `uid` varchar(255) DEFAULT NULL COMMENT 'The user ID',
              `answer` varchar(255) DEFAULT NULL,
              `created_at` timestamp NULL DEFAULT NULL,
              `updated_at` timestamp NULL DEFAULT NULL
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Contains users security question answers.';

            CREATE TABLE `security_questions_incorrect` (
              `aid` int(11) NOT NULL COMMENT 'Unique attempt ID',
              `sqid` int(11) DEFAULT NULL COMMENT 'The security question ID',
              `uid` varchar(255) DEFAULT NULL COMMENT 'The user ID',
              `ip` int(10) UNSIGNED DEFAULT NULL COMMENT 'The IP address of the visitor that attempted to answer the question as the user',
              `created_at` timestamp NULL DEFAULT NULL COMMENT 'Timestamp of the failed attempt'
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Tracks incorrect answer attempts by IP';

            CREATE TABLE `tokens` (
              `id` int(11) NOT NULL,
              `user_uid` varchar(255) NOT NULL,
              `refresh_token_id` int(11) DEFAULT NULL,
              `subject` varchar(255) NOT NULL,
              `signed_string` text NOT NULL,
              `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
              `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

            CREATE TABLE `users` (
              `uid` varchar(255) NOT NULL COMMENT 'The UID is AWS Cognito User Id which is unique and never reassignable to another user.',
              `email` varchar(255) DEFAULT NULL COMMENT 'User email',
              `username` varchar(255) NOT NULL COMMENT 'Unique user name.',
              `password` varchar(255) NOT NULL,
              `first_name` varchar(255) DEFAULT NULL,
              `last_name` varchar(255) DEFAULT NULL,
              `phone_number` varchar(255) DEFAULT NULL,
              `sms_phone_number` varchar(45) DEFAULT NULL,
              `company_name` varchar(255) DEFAULT NULL,
              `is_corporate` tinyint(1) NOT NULL DEFAULT '0',
              `created_at` timestamp NULL DEFAULT NULL,
              `updated_at` timestamp NULL DEFAULT NULL,
              `last_login_at` timestamp NULL DEFAULT NULL,
              `last_login_ip` varchar(45) DEFAULT NULL,
              `role_name` varchar(55) DEFAULT NULL,
              `user_group_id` int(11) UNSIGNED DEFAULT NULL,
              `status` enum('active','pending','blocked','dormant','canceled') DEFAULT NULL,
              `date_of_birth_year` int(5) DEFAULT NULL,
              `date_of_birth_month` tinyint(4) DEFAULT NULL,
              `date_of_birth_day` tinyint(4) DEFAULT NULL,
              `document_type` enum('passport','driver-license','gov-issued-photo-id') DEFAULT NULL,
              `document_personal_id` varchar(255) DEFAULT NULL,
              `country_of_residence_iso2` varchar(2) DEFAULT NULL,
              `country_of_citizenship_iso2` varchar(2) DEFAULT NULL,
              `pa_zip_postal_code` varchar(45) DEFAULT NULL,
              `pa_address` varchar(255) DEFAULT NULL,
              `pa_address_2nd_line` varchar(255) DEFAULT NULL,
              `pa_city` varchar(45) DEFAULT NULL,
              `pa_country_iso2` varchar(2) DEFAULT NULL,
              `pa_state_prov_region` varchar(255) DEFAULT NULL,
              `position` varchar(255) DEFAULT NULL,
              `internal_notes` text,
              `class_id` varchar(50) DEFAULT NULL,
              `ma_zip_postal_code` varchar(45) DEFAULT NULL,
              `ma_state_prov_region` varchar(255) DEFAULT NULL,
              `ma_phone_number` varchar(45) DEFAULT NULL,
              `ma_name` varchar(255) DEFAULT NULL,
              `ma_country_iso2` varchar(2) DEFAULT NULL,
              `ma_address` varchar(255) DEFAULT NULL,
              `ma_address_2nd_line` varchar(255) DEFAULT NULL,
              `ma_city` varchar(45) DEFAULT NULL,
              `ma_as_physical` tinyint(1) NOT NULL DEFAULT '0',
              `bo_full_name` varchar(255) DEFAULT NULL,
              `bo_phone_number` varchar(45) DEFAULT NULL,
              `bo_date_of_birth_year` int(5) DEFAULT NULL,
              `bo_date_of_birth_month` tinyint(4) DEFAULT NULL,
              `bo_date_of_birth_day` tinyint(4) DEFAULT NULL,
              `bo_document_personal_id` varchar(255) DEFAULT NULL,
              `bo_document_type` enum('passport','driver-license','gov-issued-photo-id') DEFAULT NULL,
              `bo_address` varchar(255) DEFAULT NULL,
              `bo_relationship` varchar(255) DEFAULT NULL,
              `home_phone_number` varchar(45) DEFAULT NULL,
              `office_phone_number` varchar(45) DEFAULT NULL,
              `fax` varchar(45) DEFAULT NULL,
              `blocked_until` timestamp NULL DEFAULT NULL,
              `last_acted_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
              `challenge_name` varchar(45) DEFAULT NULL
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

            CREATE TABLE `users_accesslog` (
              `alid` int(11) NOT NULL COMMENT 'Primary key',
              `uid` varchar(255) NOT NULL COMMENT 'User id',
              `ip` int(10) UNSIGNED NOT NULL COMMENT 'IP address from which the user came',
              `created_at` timestamp NULL DEFAULT NULL COMMENT 'User login date'
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Users access logs';

            CREATE TABLE `user_groups` (
              `id` int(11) UNSIGNED NOT NULL,
              `name` varchar(255) DEFAULT NULL,
              `description` varchar(255) NOT NULL,
              `created_at` timestamp NULL DEFAULT NULL,
              `updated_at` timestamp NULL DEFAULT NULL
            ) ENGINE=InnoDB DEFAULT CHARSET=utf8;


            ALTER TABLE `blocked_ips`
              ADD PRIMARY KEY (`id`);

            ALTER TABLE `confirmation_codes`
              ADD PRIMARY KEY (`id`),
              ADD UNIQUE KEY `code_UNIQUE` (`code`),
              ADD KEY `uid_fk` (`user_uid`);

            ALTER TABLE `devices`
              ADD PRIMARY KEY (`id`),
              ADD UNIQUE KEY `id_UNIQUE` (`id`),
              ADD UNIQUE KEY `push_token_UNIQUE` (`push_token`);

            ALTER TABLE `fail_auth_attempts`
              ADD PRIMARY KEY (`id`);

            ALTER TABLE `registration_requests`
              ADD PRIMARY KEY (`id`),
              ADD UNIQUE KEY `uid_UNIQUE` (`uid`);

            ALTER TABLE `schema_migrations`
              ADD PRIMARY KEY (`version`);

            ALTER TABLE `security_questions`
              ADD PRIMARY KEY (`sqid`),
              ADD KEY `uid` (`uid`);

            ALTER TABLE `security_questions_answers`
              ADD PRIMARY KEY (`aid`),
              ADD KEY `FC_security_questions_answers_user_idx` (`uid`),
              ADD KEY `FK_security_questions_answers_questions_idx` (`sqid`);

            ALTER TABLE `security_questions_incorrect`
              ADD PRIMARY KEY (`aid`),
              ADD KEY `FK_security_questions_incorrect_questions_idx` (`sqid`),
              ADD KEY `FK_security_questions_incorrect_users_idx` (`uid`);

            ALTER TABLE `tokens`
              ADD PRIMARY KEY (`id`),
              ADD KEY `uid_fk` (`user_uid`),
              ADD KEY `refresh_token_fk` (`user_uid`),
              ADD KEY `FK_refresh_token_tokens` (`refresh_token_id`);

            ALTER TABLE `users`
              ADD PRIMARY KEY (`uid`),
              ADD UNIQUE KEY `username` (`username`),
              ADD UNIQUE KEY `uix_users_uid` (`uid`),
              ADD UNIQUE KEY `uix_users_email` (`email`),
              ADD KEY `user_group_index` (`user_group_id`);

            ALTER TABLE `users_accesslog`
              ADD PRIMARY KEY (`alid`),
              ADD KEY `FK_users_accesslog_users_idx` (`uid`);

            ALTER TABLE `user_groups`
              ADD PRIMARY KEY (`id`);


            ALTER TABLE `blocked_ips`
              MODIFY `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Primary Key: unique ID for IP addresses.';

            ALTER TABLE `confirmation_codes`
              MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;

            ALTER TABLE `devices`
              MODIFY `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'The device id', AUTO_INCREMENT=1;

            ALTER TABLE `fail_auth_attempts`
              MODIFY `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Primary Key: Unique ID', AUTO_INCREMENT=1;

            ALTER TABLE `registration_requests`
              MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;

            ALTER TABLE `security_questions_answers`
              MODIFY `aid` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;

            ALTER TABLE `security_questions_incorrect`
              MODIFY `aid` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Unique attempt ID';

            ALTER TABLE `tokens`
              MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;

            ALTER TABLE `users_accesslog`
              MODIFY `alid` int(11) NOT NULL AUTO_INCREMENT COMMENT 'Primary key', AUTO_INCREMENT=1;

            ALTER TABLE `user_groups`
              MODIFY `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=1;


            ALTER TABLE `confirmation_codes`
              ADD CONSTRAINT `FK_confirmation_codes_users` FOREIGN KEY (`user_uid`) REFERENCES `users` (`uid`) ON DELETE CASCADE ON UPDATE CASCADE;

            ALTER TABLE `security_questions_answers`
              ADD CONSTRAINT `FK_security_questions_answers_questions` FOREIGN KEY (`sqid`) REFERENCES `security_questions` (`sqid`) ON DELETE NO ACTION ON UPDATE NO ACTION,
              ADD CONSTRAINT `FK_security_questions_answers_users` FOREIGN KEY (`uid`) REFERENCES `users` (`uid`) ON DELETE CASCADE ON UPDATE CASCADE;

            ALTER TABLE `security_questions_incorrect`
              ADD CONSTRAINT `FK_security_questions_incorrect_questions` FOREIGN KEY (`sqid`) REFERENCES `security_questions` (`sqid`) ON DELETE CASCADE ON UPDATE CASCADE,
              ADD CONSTRAINT `FK_security_questions_incorrect_users` FOREIGN KEY (`uid`) REFERENCES `users` (`uid`) ON DELETE CASCADE ON UPDATE CASCADE;

            ALTER TABLE `tokens`
              ADD CONSTRAINT `FK_refresh_token_tokens` FOREIGN KEY (`refresh_token_id`) REFERENCES `tokens` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
              ADD CONSTRAINT `FK_tokens_users` FOREIGN KEY (`user_uid`) REFERENCES `users` (`uid`) ON DELETE CASCADE ON UPDATE CASCADE;

            ALTER TABLE `users`
              ADD CONSTRAINT `FK_users_user_groups` FOREIGN KEY (`user_group_id`) REFERENCES `user_groups` (`id`);

            ALTER TABLE `users_accesslog`
              ADD CONSTRAINT `FK_users_accesslog_users` FOREIGN KEY (`uid`) REFERENCES `users` (`uid`) ON DELETE CASCADE ON UPDATE CASCADE;
SQL;
    }
}
