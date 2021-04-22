<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class InsertRoles extends Migration
{
    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function up()
    {
        $sql = "INSERT INTO `roles` (`slug`, `name`) VALUES
           ('root', 'Root'),
           ('admin', 'Admin'),
           ('client', 'Client');";

        app("db")->getPdo()->exec($sql);
    }

    /**
     * Run the migrations.
     *
     * @return void
     */
    public function down()
    {
    }
}
