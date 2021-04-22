<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class AlterDevicesAddDeviceId extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        DB::getDoctrineSchemaManager()->getDatabasePlatform()->registerDoctrineTypeMapping('enum', 'string');
        DB::table('devices')->delete();
        Schema::table('devices', function (Blueprint $table) {
            $table->dropIndex('id_UNIQUE'); // duplicated index
            $table->string('id', 255)->change();
            $table->string('pin', 255)->nullable(false)->change();
        });

        Schema::table('users', function (Blueprint $table) {
            $table->dropIndex('uix_users_uid'); // duplicated index
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
    }
}
