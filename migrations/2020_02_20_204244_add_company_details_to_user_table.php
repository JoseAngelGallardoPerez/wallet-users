<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddCompanyDetailsToUserTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('users', function (Blueprint $table) {
            $table->string('company_type')->nullable($value = true)->after('company_name');
            $table->string('company_role')->nullable($value = true)->after('company_type');
            $table->string('director_first_name')->nullable($value = true)->after('company_role');
            $table->string('director_last_name')->nullable($value = true)->after('director_first_name');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('users', function (Blueprint $table) {
            $table->dropColumn('company_type');
            $table->dropColumn('company_role');
            $table->dropColumn('director_first_name');
            $table->dropColumn('director_last_name');
        });
    }
}
