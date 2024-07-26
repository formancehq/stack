<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Utils\SpeakeasyMetadata;
class GetTransferInitiationRequest
{
    /**
     * The transfer ID.
     *
     * @var string $transferId
     */
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=transferId')]
    public string $transferId;

    public function __construct()
    {
        $this->transferId = '';
    }
}