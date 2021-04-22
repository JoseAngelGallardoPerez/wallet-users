<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddAddressesTable extends Migration
{
    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('addresses', function (Blueprint $table) {
            $table->charset = 'utf8';
            $table->collation = 'utf8_general_ci';

            $table->increments('id');
            $table->enum('type', ['physical','mailing']);
            $table->string('user_id', 255);
            $table->string('country_iso_two', 2)->default("");
            $table->string('region', 255)->default("");
            $table->string('city', 255)->default("");
            $table->string('zip_code', 100)->default("");
            $table->string('address', 255)->default("");
            $table->string('address_second_line', 255)->default("");
            $table->string('name', 255)->default("");
            $table->string('phone_number', 255)->default("");
            $table->string('description', 255)->default("");

            $table->decimal('latitude', 19, 15)->nullable();
            $table->decimal('longitude', 19, 15)->nullable();

            $table->timestamps();

            $table->foreign('user_id')->references('uid')->on('users')->onDelete('cascade');
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
