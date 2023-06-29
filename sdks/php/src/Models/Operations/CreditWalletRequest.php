<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use \formance\stack\Utils\SpeakeasyMetadata;
class CreditWalletRequest
{
	#[SpeakeasyMetadata('request:mediaType=application/json')]
    public ?\formance\stack\Models\Shared\CreditWalletRequest $creditWalletRequest = null;
    
	#[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=id')]
    public string $id;
    
	public function __construct()
	{
		$this->creditWalletRequest = null;
		$this->id = "";
	}
}
