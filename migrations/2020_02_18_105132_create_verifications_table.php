<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateVerificationsTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('verifications', function (Blueprint $table) {
            $table->charset = 'utf8';
            $table->collation = 'utf8_general_ci';
            $table->increments('id');
            $table->enum('status', ['pending','progress','approved','cancelled'])->default('pending')->nullable(false);
            $table->enum('type', ['personal_id', 'personal_photo', 'credit_rating'])->default('personal_id')->nullable(false);
            $table->string('user_uid', 255)->nullable(false);
            $table->unsignedInteger('file_id')->nullable(false);
            $table->timestamps();
            $table->foreign('user_uid')->references('uid')->on('users')->onDelete('cascade');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('verifications');
    }
}
