<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;

use \formance\stack\Utils\SpeakeasyMetadata;
class RemoveAccountFromPoolRequest
{
    /**
     * The account ID.
     * 
     * @var string $accountId
     */
	#[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=accountId')]
    public string $accountId;
    
    /**
     * The pool ID.
     * 
     * @var string $poolId
     */
	#[SpeakeasyMetadata('pathParam:style=simple,explode=false,name=poolId')]
    public string $poolId;
    
	public function __construct()
	{
		$this->accountId = "";
		$this->poolId = "";
	}
}
