<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class InsertForms extends Migration
{
    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function up()
    {

        $form = file_get_contents(dirname(__FILE__) . '/2020_07_20_172422_client_form.json');



        DB::insert(
            'INSERT INTO `forms` (`type`, `initiator_role_names`, `owner_role_names`, `form`) values (?, ?, ?, ?)',
            ['sign_up', '', '["client"]', $form]
        );
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
