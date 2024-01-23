<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use \formance\stack\Utils\SpeakeasyMetadata;
class ReverseTransferInitiationRequest
{
	#[SpeakeasyMetadata('request:mediaType=application/json')]
    public \formance\stack\Models\Shared\ReverseTransferInitiationRequest $reverseTransferInitiationRequest;
    
    /**
     * The transfer ID.
     * 
     * @var string $transferId
     */
	#[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=transferId')]
    public string $transferId;
    
	public function __construct()
	{
		$this->reverseTransferInitiationRequest = new \formance\stack\Models\Shared\ReverseTransferInitiationRequest();
		$this->transferId = "";
	}
}
