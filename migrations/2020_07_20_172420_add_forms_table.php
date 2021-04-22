<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFormsTable extends Migration
{
    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('forms', function (Blueprint $table) {
            $table->charset = 'utf8';
            $table->collation = 'utf8_general_ci';

            $table->increments('id');
            $table->string('type', 255);
            $table->text('initiator_role_names');
            $table->text('owner_role_names');
            $table->text('form');
        });
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