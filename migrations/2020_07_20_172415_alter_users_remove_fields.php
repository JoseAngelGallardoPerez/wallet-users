<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AlterUsersRemoveFields extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function down()
    {
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function up()
    {
        DB::getDoctrineSchemaManager()->getDatabasePlatform()->registerDoctrineTypeMapping('enum', 'string');
        Schema::table('users', function (Blueprint $table) {
            // add
            $table->date('date_of_birth')->nullable();

            // rename
            $table->renameColumn('country_of_residence_iso2', 'country_of_residence_iso_two');
            $table->renameColumn('country_of_citizenship_iso2', 'country_of_citizenship_iso_two');

            // change
            $table->boolean('is_email_confirmed')->default(0)->change();

            // drop
            $table->dropColumn('date_of_birth_year');
            $table->dropColumn('date_of_birth_month');
            $table->dropColumn('date_of_birth_day');

            $table->dropColumn('pa_zip_postal_code');
            $table->dropColumn('pa_address');
            $table->dropColumn('pa_address_2nd_line');
            $table->dropColumn('pa_city');
            $table->dropColumn('pa_country_iso2');
            $table->dropColumn('pa_state_prov_region');

            $table->dropColumn('ma_zip_postal_code');
            $table->dropColumn('ma_state_prov_region');
            $table->dropColumn('ma_phone_number');
            $table->dropColumn('ma_name');
            $table->dropColumn('ma_country_iso2');
            $table->dropColumn('ma_city');
            $table->dropColumn('ma_address');
            $table->dropColumn('ma_address_2nd_line');
            $table->dropColumn('ma_as_physical');

            $table->dropColumn('bo_full_name');
            $table->dropColumn('bo_phone_number');
            $table->dropColumn('bo_date_of_birth_year');
            $table->dropColumn('bo_date_of_birth_month');
            $table->dropColumn('bo_date_of_birth_day');
            $table->dropColumn('bo_document_personal_id');
            $table->dropColumn('bo_document_type');
            $table->dropColumn('bo_address');
            $table->dropColumn('bo_relationship');

            // foreign keys
            $table->foreign('role_name')->references('slug')->on('roles')->onDelete('set null');
        });
    }
}
