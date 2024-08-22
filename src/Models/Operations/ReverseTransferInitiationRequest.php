<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Models\Shared;
use formance\stack\Utils\SpeakeasyMetadata;
class ReverseTransferInitiationRequest
{
    /**
     *
     * @var Shared\ReverseTransferInitiationRequest $reverseTransferInitiationRequest
     */
    #[SpeakeasyMetadata('request:mediaType=application/json')]
    public Shared\ReverseTransferInitiationRequest $reverseTransferInitiationRequest;

    /**
     * The transfer ID.
     *
     * @var string $transferId
     */
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=transferId')]
    public string $transferId;

    /**
     * @param  ?Shared\ReverseTransferInitiationRequest  $reverseTransferInitiationRequest
     * @param  ?string  $transferId
     */
    public function __construct(?Shared\ReverseTransferInitiationRequest $reverseTransferInitiationRequest = null, ?string $transferId = null)
    {
        $this->reverseTransferInitiationRequest = $reverseTransferInitiationRequest;
        $this->transferId = $transferId;
    }
}