<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AlterUsersChangeRestriction extends Migration
{
    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('users', function (Blueprint $table) {
            $table->dropForeign(['company_id']);

            $table->foreign('company_id')->references('id')->on('companies')->onDelete('set null');
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
