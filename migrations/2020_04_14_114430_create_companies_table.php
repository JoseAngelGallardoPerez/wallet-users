<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateCompaniesTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('companies', function (Blueprint $table) {
            $table->bigIncrements('id');
            $table->string('company_name');
            $table->string('company_type');
            $table->string('company_role');
            $table->string('director_first_name');
            $table->string('director_last_name');
            $table->timestamps();
        });

        Schema::table('users', function (Blueprint $table) {
            $table->dropColumn('company_name');
            $table->dropColumn('company_type');
            $table->dropColumn('company_role');
            $table->dropColumn('director_first_name');
            $table->dropColumn('director_last_name');
            $table->unsignedBigInteger('company_id')->nullable(true);
        });

        Schema::table('users', function (Blueprint $table) {
            $table->foreign('company_id')->references('id')->on('companies');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('companies');

        Schema::table('users', function (Blueprint $table) {
            $table->dropColumn('company_id');
            $table->string('company_name');
            $table->string('company_type');
            $table->string('company_role');
            $table->string('director_first_name');
            $table->string('director_last_name');
        });
    }
}
