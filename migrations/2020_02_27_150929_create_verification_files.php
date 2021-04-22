<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class CreateVerificationFiles extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        DB::beginTransaction();
        try {
            Schema::create('verification_files', function (Blueprint $table) {
                $table->increments('id');
                $table->unsignedInteger('verification_id');
                $table->unsignedInteger('file_id')->nullable(false);
                $table->timestamps();

                $table->foreign('verification_id')->references('id')->on('verifications')->onDelete('cascade');
            });

            Schema::table('verifications', function (Blueprint $table) {
                $table->dropColumn('file_id');
            });
        } catch (\Throwable $e) {
            DB::rollBack();
            throw $e;
        }
        DB::commit();
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        DB::beginTransaction();
        try {
            Schema::dropIfExists('verification_files');

            Schema::table('verifications', function (Blueprint $table) {
                $table->unsignedInteger('file_id')->nullable(false);
            });
        } catch (\Throwable $e) {
            DB::rollBack();
            throw $e;
        }
        DB::commit();
    }
}
