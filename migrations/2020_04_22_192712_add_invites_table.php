<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddInvitesTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('invites', function (Blueprint $table) {
            $table->charset = 'utf8';
            $table->collation = 'utf8_general_ci';
            $table->increments('id');
            $table->string('code', 255)->unique();
            $table->string('to', 255)->nullable(false);
            $table->unsignedInteger('uses')->nullable(false)->default(0);
            $table->unsignedInteger('max_usages')->nullable(false)->default(0);
            $table->string('user_uid', 255)->nullable(false);
            $table->timestamp('expires_at')->nullable(true);
            $table->timestamps();
            $table->foreign('user_uid')->references('uid')->on('users')->onDelete('cascade');
            $table->index('code');
            $table->index('user_uid');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('invites', function (Blueprint $table)
        {
            $table->dropIndex('user_uid');
            $table->dropIndex('code');
        });

        Schema::dropIfExists('invites');
    }
}
