<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Utils\SpeakeasyMetadata;
class GetBalanceRequest
{
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=balanceName')]
    public string $balanceName;

    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=id')]
    public string $id;

    public function __construct()
    {
        $this->balanceName = '';
        $this->id = '';
    }
}