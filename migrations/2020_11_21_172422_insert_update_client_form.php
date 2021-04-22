<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class InsertUpdateClientForm extends Migration
{
    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function up()
    {
        $form = '{
                   "fields": [
                     {
                       "name": "uid",
                       "type": "string"
                     },
                     {
                       "name": "firstName",
                       "type": "string",
                       "validators": [
                         {
                           "name": "required",
                           "param": ""
                         },
                         {
                           "name": "max",
                           "param": "255"
                         }
                       ]
                     },
                     {
                       "name": "lastName",
                       "type": "string",
                       "validators": [
                         {
                           "name": "required",
                           "param": ""
                         },
                         {
                           "name": "max",
                           "param": "255"
                         }
                       ]
                     },
                     {
                       "name": "phoneNumber",
                       "type": "string",
                       "validators": [
                         {
                           "name": "required"
                         },
                         {
                           "name": "uniquePhoneNumber",
                           "param": "Uid"
                         }
                       ]
                     },
                     {
                       "name": "dateOfBirth",
                       "type": "string",
                       "validators": [
                         {
                           "name": "omitempty"
                         },
                         {
                           "name": "dayBeforeNow"
                         }
                       ]
                     },
                     {
                       "name": "documentPersonalId",
                       "type": "string",
                       "validators": [
                         {
                           "name": "omitempty"
                         },
                         {
                           "name": "max",
                           "param": "255"
                         }
                       ]
                     },
                     {
                       "name": "companyDetails",
                       "type": "object",
                       "validators": [
                         {
                           "name": "omitempty"
                         },
                         {
                           "name": "dive"
                         }
                       ],
                       "children": [
                         {
                           "name": "id",
                           "type": "int",
                           "validators": [
                             {
                               "name": "omitempty"
                             }
                           ]
                         },
                         {
                           "name": "companyName",
                           "type": "string",
                           "validators": [
                             {
                               "name": "max",
                               "param": "255"
                             }
                           ]
                         }
                       ]
                     },
                     {
                       "name": "physicalAddresses",
                       "type": "array",
                       "validators": [
                         {
                           "name": "omitempty"
                         },
                         {
                           "name": "max",
                           "param": "1"
                         },
                         {
                           "name": "dive"
                         }
                       ],
                       "children": [
                         {
                           "name": "id",
                           "type": "int",
                           "validators": [
                             {
                               "name": "omitempty"
                             }
                           ]
                         },
                         {
                           "name": "countryIsoTwo",
                           "type": "string",
                           "validators": [
                             {
                               "name": "required"
                             },
                             {
                               "name": "existCountry"
                             }
                           ]
                         },
                         {
                           "name": "zipCode",
                           "type": "string",
                           "validators": [
                             {
                               "name": "required"
                             },
                             {
                               "name": "max",
                               "param": "45"
                             }
                           ]
                         },
                         {
                           "name": "address",
                           "type": "string",
                           "validators": [
                             {
                               "name": "required"
                             },
                             {
                               "name": "max",
                               "param": "255"
                             }
                           ]
                         },
                         {
                           "name": "addressSecondLine",
                           "type": "string",
                           "validators": [
                             {
                               "name": "required"
                             },
                             {
                               "name": "max",
                               "param": "255"
                             }
                           ]
                         },
                         {
                           "name": "city",
                           "type": "string",
                           "validators": [
                             {
                               "name": "required"
                             },
                             {
                               "name": "max",
                               "param": "45"
                             }
                           ]
                         },
                         {
                           "name": "region",
                           "type": "string",
                           "validators": [
                             {
                               "name": "required"
                             },
                             {
                               "name": "max",
                               "param": "255"
                             }
                           ]
                         }
                       ]
                     },
                     {
                       "name": "mailingAddresses",
                       "type": "array",
                       "validators": [
                         {
                           "name": "omitempty"
                         },
                         {
                           "name": "max",
                           "param": "1"
                         },
                         {
                           "name": "dive"
                         }
                       ],
                       "children": [
                         {
                           "name": "id",
                           "type": "int",
                           "validators": [
                             {
                               "name": "omitempty"
                             }
                           ]
                         },
                         {
                           "name": "countryIsoTwo",
                           "type": "string",
                           "validators": [
                             {
                               "name": "required"
                             },
                             {
                               "name": "existCountry"
                             }
                           ]
                         },
                         {
                           "name": "zipCode",
                           "type": "string",
                           "validators": [
                             {
                               "name": "required"
                             },
                             {
                               "name": "max",
                               "param": "45"
                             }
                           ]
                         },
                         {
                           "name": "address",
                           "type": "string",
                           "validators": [
                             {
                               "name": "required"
                             },
                             {
                               "name": "max",
                               "param": "255"
                             }
                           ]
                         },
                         {
                           "name": "addressSecondLine",
                           "type": "string",
                           "validators": [
                             {
                               "name": "required"
                             },
                             {
                               "name": "max",
                               "param": "255"
                             }
                           ]
                         },
                         {
                           "name": "city",
                           "type": "string",
                           "validators": [
                             {
                               "name": "required"
                             },
                             {
                               "name": "max",
                               "param": "45"
                             }
                           ]
                         },
                         {
                           "name": "region",
                           "type": "string",
                           "validators": [
                             {
                               "name": "required"
                             },
                             {
                               "name": "max",
                               "param": "255"
                             }
                           ]
                         }
                       ]
                     }
                   ]
                 }';

        DB::insert(
            'INSERT INTO `forms` (`type`, `initiator_role_names`, `owner_role_names`, `form`) values (?, ?, ?, ?)',
            ['update', '["client"]', '["client"]', $form]
        );

        $formByAdmin = '{
                          "fields": [
                            {
                              "name": "uid",
                              "type": "string"
                            },
                            {
                              "name": "email",
                              "type": "string",
                              "validators": [
                                {
                                  "name": "required"
                                },
                                {
                                  "name": "email"
                                },
                                {
                                  "name": "max",
                                  "param": "255"
                                },
                                {
                                  "name": "uniqueEmail",
                                  "param": "Uid"
                                }
                              ]
                            },
                            {
                              "name": "status",
                              "type": "string",
                              "validators": [
                                {
                                  "name": "required"
                                },
                                {
                                  "name": "oneof",
                                  "param": "pending active blocked dormant"
                                }
                              ]
                            },
                            {
                              "name": "firstName",
                              "type": "string",
                              "validators": [
                                {
                                  "name": "required",
                                  "param": ""
                                },
                                {
                                  "name": "max",
                                  "param": "255"
                                }
                              ]
                            },
                            {
                              "name": "lastName",
                              "type": "string",
                              "validators": [
                                {
                                  "name": "required",
                                  "param": ""
                                },
                                {
                                  "name": "max",
                                  "param": "255"
                                }
                              ]
                            },
                            {
                              "name": "phoneNumber",
                              "type": "string",
                              "validators": [
                                {
                                  "name": "required"
                                },
                                {
                                  "name": "uniquePhoneNumber",
                                  "param": "Uid"
                                }
                              ]
                            },
                            {
                              "name": "dateOfBirth",
                              "type": "string",
                              "validators": [
                                {
                                  "name": "omitempty"
                                },
                                {
                                  "name": "dayBeforeNow"
                                }
                              ]
                            },
                            {
                              "name": "documentType",
                              "type": "stringPointer",
                              "validators": [
                                {
                                  "name": "omitempty"
                                }
                              ]
                            },
                            {
                              "name": "documentPersonalId",
                              "type": "string",
                              "validators": [
                                {
                                  "name": "omitempty"
                                },
                                {
                                  "name": "max",
                                  "param": "255"
                                }
                              ]
                            },
                            {
                              "name": "userGroupId",
                              "type": "intPointer",
                              "validators": [
                                {
                                  "name": "omitempty"
                                }
                              ]
                            },
                            {
                              "name": "companyDetails",
                              "type": "object",
                              "validators": [
                                {
                                  "name": "omitempty"
                                },
                                {
                                  "name": "dive"
                                }
                              ],
                              "children": [
                                {
                                  "name": "id",
                                  "type": "int",
                                  "validators": [
                                    {
                                      "name": "omitempty"
                                    }
                                  ]
                                },
                                {
                                  "name": "companyName",
                                  "type": "string",
                                  "validators": [
                                    {
                                      "name": "max",
                                      "param": "255"
                                    }
                                  ]
                                },
                                {
                                  "name": "companyType",
                                  "type": "string",
                                  "validators": [
                                    {
                                      "name": "max",
                                      "param": "255"
                                    }
                                  ]
                                }
                              ]
                            },
                            {
                              "name": "physicalAddresses",
                              "type": "array",
                              "validators": [
                                {
                                  "name": "omitempty"
                                },
                                {
                                  "name": "max",
                                  "param": "1"
                                },
                                {
                                  "name": "dive"
                                }
                              ],
                              "children": [
                                {
                                  "name": "id",
                                  "type": "int",
                                  "validators": [
                                    {
                                      "name": "omitempty"
                                    }
                                  ]
                                },
                                {
                                  "name": "countryIsoTwo",
                                  "type": "string",
                                  "validators": [
                                    {
                                      "name": "required"
                                    },
                                    {
                                      "name": "existCountry"
                                    }
                                  ]
                                },
                                {
                                  "name": "zipCode",
                                  "type": "string",
                                  "validators": [
                                    {
                                      "name": "required"
                                    },
                                    {
                                      "name": "max",
                                      "param": "45"
                                    }
                                  ]
                                },
                                {
                                  "name": "address",
                                  "type": "string",
                                  "validators": [
                                    {
                                      "name": "required"
                                    },
                                    {
                                      "name": "max",
                                      "param": "255"
                                    }
                                  ]
                                },
                                {
                                  "name": "addressSecondLine",
                                  "type": "string",
                                  "validators": [
                                    {
                                      "name": "required"
                                    },
                                    {
                                      "name": "max",
                                      "param": "255"
                                    }
                                  ]
                                },
                                {
                                  "name": "city",
                                  "type": "string",
                                  "validators": [
                                    {
                                      "name": "required"
                                    },
                                    {
                                      "name": "max",
                                      "param": "45"
                                    }
                                  ]
                                },
                                {
                                  "name": "region",
                                  "type": "string",
                                  "validators": [
                                    {
                                      "name": "required"
                                    },
                                    {
                                      "name": "max",
                                      "param": "255"
                                    }
                                  ]
                                }
                              ]
                            },
                            {
                              "name": "mailingAddresses",
                              "type": "array",
                              "validators": [
                                {
                                  "name": "omitempty"
                                },
                                {
                                  "name": "max",
                                  "param": "1"
                                },
                                {
                                  "name": "dive"
                                }
                              ],
                              "children": [
                                {
                                  "name": "id",
                                  "type": "int",
                                  "validators": [
                                    {
                                      "name": "omitempty"
                                    }
                                  ]
                                },
                                {
                                  "name": "countryIsoTwo",
                                  "type": "string",
                                  "validators": [
                                    {
                                      "name": "required"
                                    },
                                    {
                                      "name": "existCountry"
                                    }
                                  ]
                                },
                                {
                                  "name": "zipCode",
                                  "type": "string",
                                  "validators": [
                                    {
                                      "name": "required"
                                    },
                                    {
                                      "name": "max",
                                      "param": "45"
                                    }
                                  ]
                                },
                                {
                                  "name": "address",
                                  "type": "string",
                                  "validators": [
                                    {
                                      "name": "required"
                                    },
                                    {
                                      "name": "max",
                                      "param": "255"
                                    }
                                  ]
                                },
                                {
                                  "name": "addressSecondLine",
                                  "type": "string",
                                  "validators": [
                                    {
                                      "name": "required"
                                    },
                                    {
                                      "name": "max",
                                      "param": "255"
                                    }
                                  ]
                                },
                                {
                                  "name": "city",
                                  "type": "string",
                                  "validators": [
                                    {
                                      "name": "required"
                                    },
                                    {
                                      "name": "max",
                                      "param": "45"
                                    }
                                  ]
                                },
                                {
                                  "name": "region",
                                  "type": "string",
                                  "validators": [
                                    {
                                      "name": "required"
                                    },
                                    {
                                      "name": "max",
                                      "param": "255"
                                    }
                                  ]
                                }
                              ]
                            }
                          ]
                        }';

        DB::insert(
            'INSERT INTO `forms` (`type`, `initiator_role_names`, `owner_role_names`, `form`) values (?, ?, ?, ?)',
            ['update', '["root","admin"]', '["client"]', $formByAdmin]
        );

        $formAdminByAdmin = '{
                               "fields": [
                                 {
                                   "name": "uid",
                                   "type": "string"
                                 },
                                 {
                                   "name": "username",
                                   "type": "string",
                                   "validators": [
                                     {
                                       "name": "required"
                                     },
                                     {
                                       "name": "max",
                                       "param": "255"
                                     },
                                     {
                                       "name": "uniqueUsername",
                                       "param": "Uid"
                                     }
                                   ]
                                 },
                                 {
                                   "name": "email",
                                   "type": "string",
                                   "validators": [
                                     {
                                       "name": "required"
                                     },
                                     {
                                       "name": "email"
                                     },
                                     {
                                       "name": "max",
                                       "param": "255"
                                     },
                                     {
                                       "name": "uniqueEmail",
                                       "param": "Uid"
                                     }
                                   ]
                                 },
                                 {
                                   "name": "status",
                                   "type": "string",
                                   "validators": [
                                     {
                                       "name": "required"
                                     },
                                     {
                                       "name": "oneof",
                                       "param": "pending active blocked dormant"
                                     }
                                   ]
                                 },
                                 {
                                   "name": "firstName",
                                   "type": "string",
                                   "validators": [
                                     {
                                       "name": "required",
                                       "param": ""
                                     },
                                     {
                                       "name": "max",
                                       "param": "255"
                                     }
                                   ]
                                 },
                                 {
                                   "name": "lastName",
                                   "type": "string",
                                   "validators": [
                                     {
                                       "name": "required",
                                       "param": ""
                                     },
                                     {
                                       "name": "max",
                                       "param": "255"
                                     }
                                   ]
                                 },
                                 {
                                   "name": "phoneNumber",
                                   "type": "string",
                                   "validators": [
                                     {
                                       "name": "required"
                                     },
                                     {
                                       "name": "uniquePhoneNumber",
                                       "param": "Uid"
                                     }
                                   ]
                                 },
                                 {
                                   "name": "position",
                                   "type": "string",
                                   "validators": [
                                     {
                                       "name": "omitempty"
                                     },
                                     {
                                       "name": "max",
                                       "param": "255"
                                     }
                                   ]
                                 },
                                 {
                                   "name": "classId",
                                   "type": "intPointer"
                                 },
                                 {
                                   "name": "internalNotes",
                                   "type": "string",
                                   "validators": [
                                     {
                                       "name": "omitempty"
                                     },
                                     {
                                       "name": "max",
                                       "param": "255"
                                     }
                                   ]
                                 }
                               ]
                             }';

        DB::insert(
            'INSERT INTO `forms` (`type`, `initiator_role_names`, `owner_role_names`, `form`) values (?, ?, ?, ?)',
            ['update', '["root","admin"]', '["admin"]', $formAdminByAdmin]
        );

        DB::insert(
            'INSERT INTO `forms` (`type`, `initiator_role_names`, `owner_role_names`, `form`) values (?, ?, ?, ?)',
            ['update', '["root"]', '["root"]', $formAdminByAdmin]
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
