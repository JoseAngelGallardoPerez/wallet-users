<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddAttributesTable extends Migration
{
    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('attributes', function (Blueprint $table) {
            $table->charset = 'utf8';
            $table->collation = 'utf8_general_ci';

            $table->increments('id');

            $table->string('name', 255);
            $table->string('slug', 255)->unique();
            $table->string('description', 255);

            $table->timestamps();
        });

        Schema::create('user_attribute_values', function (Blueprint $table) {
            $table->charset = 'utf8';
            $table->collation = 'utf8_general_ci';

            $table->string('user_id', 255);
            $table->unsignedInteger('attribute_id');
            $table->string('value', 255);

            $table->timestamps();

            $table->primary(['user_id', 'attribute_id']);
            $table->foreign('user_id')->references('uid')->on('users')->onDelete('cascade');
            $table->foreign('attribute_id')->references('id')->on('attributes')->onDelete('cascade');
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
