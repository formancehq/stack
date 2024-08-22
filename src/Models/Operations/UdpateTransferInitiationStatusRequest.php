<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use formance\stack\Models\Shared;
use formance\stack\Utils\SpeakeasyMetadata;
class UdpateTransferInitiationStatusRequest
{
    /**
     *
     * @var Shared\UpdateTransferInitiationStatusRequest $updateTransferInitiationStatusRequest
     */
    #[SpeakeasyMetadata('request:mediaType=application/json')]
    public Shared\UpdateTransferInitiationStatusRequest $updateTransferInitiationStatusRequest;

    /**
     * The transfer ID.
     *
     * @var string $transferId
     */
    #[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=transferId')]
    public string $transferId;

    /**
     * @param  ?Shared\UpdateTransferInitiationStatusRequest  $updateTransferInitiationStatusRequest
     * @param  ?string  $transferId
     */
    public function __construct(?Shared\UpdateTransferInitiationStatusRequest $updateTransferInitiationStatusRequest = null, ?string $transferId = null)
    {
        $this->updateTransferInitiationStatusRequest = $updateTransferInitiationStatusRequest;
        $this->transferId = $transferId;
    }
}